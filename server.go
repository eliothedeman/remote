package remote

//go:generate msgp -tests=false
//msgp:ignore PingRequest PingResponse Server Context DBStatsResponse BucketStatsResponse
//go:generate codecgen -o values.generated.go server.go server_gen.go

import (
	"errors"
	"io"
	"net"
	"net/rpc"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/ugorji/go/codec"
)

var (
	errTransactionNotFound = errors.New("Transaction not found")
	errContextNotFound     = errors.New("Context not found")
	errBucketNotFound      = errors.New("Bucket not found")
)

func rpcHandle() codec.Handle {
	return &codec.MsgpackHandle{}
}

type parent interface {
	Bucket(key []byte) *bolt.Bucket
	CreateBucket(key []byte) (*bolt.Bucket, error)
	CreateBucketIfNotExists(key []byte) (*bolt.Bucket, error)
	DeleteBucket(key []byte) error
}

// A Server runs a remote view to a boltdb.
type Server struct {
	db           *bolt.DB
	contexts     map[uint64]*Context
	cCount       uint64
	sync.RWMutex `msg:"-"`
}

// OpenServer opens a bolt db and and wraps it in a server.
func OpenServer(path string) (*Server, error) {
	db, err := bolt.Open(path, 0777, nil)
	if err != nil {
		return nil, err
	}

	return &Server{
		db:       db,
		contexts: make(map[uint64]*Context),
	}, nil
}

// ServeTCP will serve the rpc via a tcp socket
func (s *Server) ServeTCP(addr string) error {
	l, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}

	var conn net.Conn
	for {
		conn, err = l.Accept()
		if err != nil {
			break
		}
		go s.ServeConn(conn)
	}

	return err
}

// ServeConn will serve the rpc via a io.ReadWriteCloser
func (s *Server) ServeConn(conn io.ReadWriteCloser) error {
	srv := rpc.NewServer()
	srv.RegisterName("srv", s)
	cod := codec.GoRpc.ServerCodec(conn, rpcHandle())
	srv.ServeCodec(cod)
	return nil
}

func (s *Server) getContext(id uint64) *Context {
	s.RLock()
	c, ok := s.contexts[id]
	s.RUnlock()

	if !ok {
		return nil
	}

	return c
}

func (s *Server) newContext(p parent) *Context {
	s.Lock()
	s.cCount++
	c := &Context{
		id:      s.cCount,
		parent:  p,
		buckets: make(map[uint64]*bolt.Bucket),
	}
	s.contexts[s.cCount] = c
	s.Unlock()
	return c
}

func (s *Server) closeContext(c *Context) {
	s.Lock()
	c.Lock()
	delete(s.contexts, c.id)
	if c.forEachCancel != nil {
		close(c.forEachCancel)
		c.forEachCancel = nil
	}
	if c.forEachOut != nil {
		close(c.forEachOut)
		c.forEachOut = nil
	}
	c.Unlock()
	s.Unlock()
}

// A Context holds information about a transaction.
type Context struct {
	sync.RWMutex  `msg:"-"`
	id            uint64
	tx            *bolt.Tx
	parent        parent
	bCount        uint64
	buckets       map[uint64]*bolt.Bucket
	forEachOut    chan *BucketForEachResponse
	forEachCancel chan struct{}
}

func (c *Context) getBucket(id uint64) *bolt.Bucket {
	c.RLock()
	b, ok := c.buckets[id]
	c.RUnlock()
	if !ok {
		return nil
	}
	return b
}

func (c *Context) bucket(name []byte) (*bolt.Bucket, uint64) {
	var id uint64
	c.Lock()
	b := c.parent.Bucket(name)
	if b != nil {
		c.bCount++
		id = c.bCount
		c.buckets[id] = b
	}
	c.Unlock()

	return b, id
}

func (c *Context) createBucket(name []byte) (*bolt.Bucket, uint64, error) {
	var id uint64
	c.Lock()
	b, err := c.parent.CreateBucket(name)
	if err == nil {
		c.bCount++
		id = c.bCount
		c.buckets[id] = b
	}
	c.Unlock()

	return b, id, err
}

func (c *Context) createBucketIfNotExists(name []byte) (*bolt.Bucket, uint64, error) {
	var id uint64
	c.Lock()
	b, err := c.parent.CreateBucketIfNotExists(name)
	if err == nil {
		c.bCount++
		id = c.bCount
		c.buckets[id] = b
	}
	c.Unlock()

	return b, id, err
}

// PingRequest is a test ping request
type PingRequest struct {
	T time.Time
}

// PingResponse is a test ping result
type PingResponse struct {
	To   time.Duration
	From time.Time
}

// RoundTrip gives the time the request took round trip
func (p *PingResponse) RoundTrip() time.Duration {
	return p.To + time.Now().Sub(p.From)
}

// Ping the remote server.
func (s *Server) Ping(req PingRequest, resp *PingResponse) error {
	resp.From = time.Now()
	resp.To = resp.From.Sub(req.T)
	return nil
}

// Empty requests/responses are for functions that do not requre inputs.
type Empty struct {
}

// DBStatsResponse contains boltdb stats.
type DBStatsResponse struct {
	bolt.Stats
}

// DBStats returns database level stats
func (s *Server) DBStats(Empty, resp *DBStatsResponse) error {
	stats := s.db.Stats()
	resp.Stats = stats
	return nil
}

// BucketStatsRequest contains stats about a bucket.
type BucketStatsRequest struct {
	ContextID uint64
	BucketID  uint64
}

// BucketStatsResponse contains stats about a bucket.
type BucketStatsResponse struct {
	bolt.BucketStats
}

// BucketStats returns the stats about a bucket.
func (s *Server) BucketStats(req *BucketStatsRequest, resp *BucketStatsResponse) error {
	c := s.getContext(req.ContextID)
	if c == nil {
		logrus.WithFields(logrus.Fields{
			"context_id": req.ContextID,
		}).Error(errContextNotFound)

		return errContextNotFound
	}

	b := c.getBucket(req.BucketID)
	if b == nil {
		logrus.WithFields(logrus.Fields{
			"context_id": req.ContextID,
			"bucket_id":  req.BucketID,
		}).Error(errBucketNotFound)
		return errBucketNotFound
	}

	resp.BucketStats = b.Stats()

	return nil
}

// BeginTransactionRequest is the response to StartTransaction
type BeginTransactionRequest struct {
	Writable bool
}

// BeginTransactionResponse is the response to StartTransaction
type BeginTransactionResponse struct {
	ContextID uint64
}

// BeginTransaction creates a new transaction with the given mode.
func (s *Server) BeginTransaction(req *BeginTransactionRequest, resp *BeginTransactionResponse) error {
	c := s.newContext(nil)

	tx, err := s.db.Begin(req.Writable)
	c.tx = tx
	c.parent = tx
	resp.ContextID = c.id
	if err != nil {
		s.closeContext(c)
	}

	logrus.WithFields(logrus.Fields{
		"writeable":  req.Writable,
		"context_id": c.id,
	}).Info("Starting transaction")

	return err
}

// CommitTransactionResponse contains the stats for the transaction that was closed.
type CommitTransactionResponse struct {
}

// CommitTransaction creates a new transaction with the given mode.
func (s *Server) CommitTransaction(contextID uint64, c *CommitTransactionResponse) error {
	ctx := s.getContext(contextID)
	if ctx == nil {
		return errors.New("Context not found")

	}
	logrus.WithFields(logrus.Fields{
		"context_id": ctx.id,
	}).Info("Commiting transaction")

	s.closeContext(ctx)
	return ctx.tx.Commit()
}

// RollbackTransactionResponse contains the stats for the transaction that was closed.
type RollbackTransactionResponse struct {
}

// RollbackTransaction creates a new transaction with the given mode.
func (s *Server) RollbackTransaction(contextID uint64, r *RollbackTransactionResponse) error {
	ctx := s.getContext(contextID)
	if ctx == nil {
		return errors.New("Transaction not found")
	}

	logrus.WithFields(logrus.Fields{
		"context_id": ctx.id,
	}).Error("Rolling back transaction")

	s.closeContext(ctx)
	return ctx.tx.Rollback()
}

// BucketRequest contains the stats for the transaction that was closed.
type BucketRequest struct {
	ContextID uint64
	Key       []byte
}

// BucketResponse contains the stats for the transaction that was closed.
type BucketResponse struct {
	BucketID        uint64
	BucketContextID uint64
}

// Bucket creates a new transaction with the given mode.
func (s *Server) Bucket(req BucketRequest, resp *BucketResponse) error {
	c := s.getContext(req.ContextID)
	if c == nil {
		return errors.New("Transaction not found")
	}

	b, id := c.bucket(req.Key)

	if b == nil {
		return errors.New("Bucket not found")
	}

	resp.BucketContextID = c.id
	resp.BucketID = id
	c = s.newContext(b)

	return nil
}

// DeleteBucket removes the bucket with the given key.
func (s *Server) DeleteBucket(req BucketRequest, resp *BucketResponse) error {
	c := s.getContext(req.ContextID)
	if c == nil {
		return errors.New("Transaction not found")
	}

	// TODO this needs to clear this bucket out of the context
	err := c.parent.DeleteBucket(req.Key)
	return err
}

// CreateBucket creates a new bucket with the given name.
func (s *Server) CreateBucket(req BucketRequest, resp *BucketResponse) error {
	c := s.getContext(req.ContextID)
	if c == nil {
		logrus.WithFields(logrus.Fields{
			"context_id": req.ContextID,
		}).Error(errContextNotFound)
		return errContextNotFound
	}

	b, id, err := c.createBucket(req.Key)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"context_id": req.ContextID,
			"key":        string(req.Key),
		}).Error(err)
		return err
	}

	resp.BucketContextID = c.id
	resp.BucketID = id
	c = s.newContext(b)

	return nil
}

// CreateBucketIfNotExists creates a new bucket if it doesn't already exist. Returns the bucket regardless.
func (s *Server) CreateBucketIfNotExists(req BucketRequest, resp *BucketResponse) error {
	c := s.getContext(req.ContextID)
	if c == nil {
		logrus.WithFields(logrus.Fields{
			"context_id": req.ContextID,
		}).Error(errContextNotFound)
		return errContextNotFound
	}

	b, id, err := c.createBucketIfNotExists(req.Key)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"context_id": req.ContextID,
			"key":        string(req.Key),
		}).Error(err)
		return err
	}

	resp.BucketContextID = c.id
	resp.BucketID = id
	c = s.newContext(b)

	return nil
}

// GetReqeust has get request data.
type GetReqeust struct {
	BucketID  uint64
	ContextID uint64
	Key       []byte
}

// GetResponse has get request data.
type GetResponse struct {
	Val []byte
}

func cnf(id uint64) error {
	logrus.WithFields(logrus.Fields{
		"context_id": id,
	}).Error(errContextNotFound)
	return errContextNotFound
}

func bnf(c, b uint64) error {
	logrus.WithFields(logrus.Fields{
		"context_id": c,
		"bucket_id":  b,
	}).Error(errBucketNotFound)
	return errBucketNotFound
}

// Get returns the value stored at the given key.
func (s *Server) Get(req *GetReqeust, resp *GetResponse) error {
	c := s.getContext(req.ContextID)
	if c == nil {
		return cnf(req.ContextID)
	}
	b := c.getBucket(req.BucketID)
	if b == nil {
		return bnf(req.ContextID, req.BucketID)
	}

	resp.Val = b.Get(req.Key)
	logrus.WithFields(logrus.Fields{
		"context_id": c,
		"bucket_id":  b,
		"key":        string(req.Key),
	}).Debug("Get")

	return nil
}

// PutReqeust has get request data.
type PutReqeust struct {
	BucketID  uint64
	ContextID uint64
	Key       []byte
	Val       []byte
}

// PutResponse has get request data.
type PutResponse struct {
}

// Put inserts the given value at the given key.
func (s *Server) Put(req *PutReqeust, resp *PutResponse) error {
	c := s.getContext(req.ContextID)
	if c == nil {
		return cnf(req.ContextID)
	}
	b := c.getBucket(req.BucketID)
	if b == nil {
		return bnf(req.ContextID, req.BucketID)
	}
	logrus.WithFields(logrus.Fields{
		"context_id": c,
		"bucket_id":  b,
		"key":        string(req.Key),
	}).Debug("Put")

	return b.Put(req.Key, req.Val)
}

// BucketForEachRequest gives context for the requst.
type BucketForEachRequest struct {
	ContextID uint64
	BucketID  uint64
}

// BucketForEachResponse is a response of every key value pair in a bucket.
type BucketForEachResponse struct {
	Key         []byte
	Value       []byte
	Index, Size uint64
}

// BucketForEachStart runs back all k,v pairs so that a function can be run over them.
func (s *Server) BucketForEachStart(req *BucketForEachRequest, resp *BucketForEachResponse) error {
	c := s.getContext(req.ContextID)
	b := c.getBucket(req.BucketID)
	logrus.WithFields(logrus.Fields{
		"context_id": req.ContextID,
		"bucket_id":  req.BucketID,
		"state":      "start",
	}).Info("ForEach")

	c.Lock()
	c.forEachOut = make(chan *BucketForEachResponse)
	c.forEachCancel = make(chan struct{})
	c.Unlock()

	go func() {
		var stopped bool
		stats := b.Stats()
		size := uint64(stats.KeyN)
		var index uint64
		b.ForEach(func(k, v []byte) error {

			// if we are stopped we don't need to do anything else.
			if stopped {
				return nil
			}
			x := &BucketForEachResponse{}
			x.Key = k
			x.Value = v
			x.Size = size
			x.Index = index
			select {
			case <-c.forEachCancel:
				stopped = true
			case c.forEachOut <- x:
				index++
			}

			return nil
		})
		logrus.WithFields(logrus.Fields{
			"context_id": req.ContextID,
			"bucket_id":  req.BucketID,
			"state":      "complete",
		}).Info("BucketForEach")

		// clean up
		c.Lock()
		close(c.forEachCancel)
		close(c.forEachOut)
		c.forEachCancel = nil
		c.forEachOut = nil
		c.Unlock()

	}()

	return nil

}

// BucketForEachNext returns the next key value pair in the iteration.
func (s *Server) BucketForEachNext(req *BucketForEachRequest, resp *BucketForEachResponse) error {

	c := s.getContext(req.ContextID)
	x := <-c.forEachOut
	resp.Index = x.Index
	resp.Size = x.Size
	resp.Key = x.Key
	resp.Value = x.Value

	logrus.WithFields(logrus.Fields{
		"context_id": req.ContextID,
		"bucket_id":  req.BucketID,
		"index":      resp.Index,
		"size":       resp.Size,
		"state":      "next",
	}).Debug("BucketForEach")
	return nil
}

// BucketForEachStop will stop the for-each loop.
func (s *Server) BucketForEachStop(req *BucketForEachRequest, resp *Empty) error {
	c := s.getContext(req.ContextID)
	c.forEachCancel <- struct{}{}
	logrus.WithFields(logrus.Fields{
		"context_id": req.ContextID,
		"bucket_id":  req.BucketID,
	}).Info("BucketForEachStop")

	return nil
}

// TransactionSize returns the size of the database from the view of the current transaction.
func (s *Server) TransactionSize(contextID uint64, size *uint64) error {
	c := s.getContext(contextID)
	if c == nil {
		return cnf(contextID)
	}

	*size = uint64(c.tx.Size())

	return nil
}

// NextSequenceRequest gives context for the next sequence call.
type NextSequenceRequest struct {
	ContextID uint64
	BucketID  uint64
}

// NextSequence will return the next unique ID for this bucket.
func (s *Server) NextSequence(req *NextSequenceRequest, resp *uint64) error {
	c := s.getContext(req.ContextID)
	b := c.getBucket(req.BucketID)

	id, err := b.NextSequence()

	*resp = id
	return err
}
