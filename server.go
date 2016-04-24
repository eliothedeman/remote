package remote

//go:generate msgp -tests=false
//msgp:ignore PingRequest PingResponse Server Context DBStatsResponse BucketStatsResponse
//go:generate codecgen -o values.generated.go server.go server_gen.go

import (
	"errors"
	"io"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/ugorji/go/codec"
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
	cod := codec.MsgpackSpecRpc.ServerCodec(conn, rpcHandle())
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
	delete(s.contexts, c.id)
	s.Unlock()
}

// A Context holds information about a transaction.
type Context struct {
	sync.RWMutex `msg:"-"`
	id           uint64
	tx           *bolt.Tx
	parent       parent
	bCount       uint64
	buckets      map[uint64]*bolt.Bucket
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
		return errors.New("Context not found.")
	}

	b := c.getBucket(req.BucketID)
	if b == nil {
		return errors.New("Bucket not found")
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

	log.Println("Created transaction with id", c.id)

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

	log.Println("Commiting transaction", contextID)
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

	log.Println("Rolling back transaction", contextID)

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
		return errors.New("Transaction not found")
	}

	b, id, err := c.createBucket(req.Key)
	if err != nil {
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
		return errors.New("Transaction not found")
	}

	b, id, err := c.createBucketIfNotExists(req.Key)
	if err != nil {
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

// Get returns the value stored at the given key.
func (s *Server) Get(req *GetReqeust, resp *GetResponse) error {
	c := s.getContext(req.ContextID)
	b := c.getBucket(req.BucketID)

	if b == nil {
		return errors.New("Bucket not found.")
	}

	resp.Val = b.Get(req.Key)

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
	b := c.getBucket(req.BucketID)

	if b == nil {
		return errors.New("Bucket not found.")
	}

	return b.Put(req.Key, req.Val)
}

// BucketForEachRequest gives context for the requst.
type BucketForEachRequest struct {
	ContextID uint64
	BucketID  uint64
}

// BucketForEachResponse is a response of every key value pair in a bucket.
type BucketForEachResponse struct {
	Keys   [][]byte
	Values [][]byte
}

// BucketForEach runs back all k,v pairs so that a function can be run over them.
func (s *Server) BucketForEach(req *BucketForEachRequest, resp *BucketForEachResponse) error {
	c := s.getContext(req.ContextID)
	b := c.getBucket(req.BucketID)

	stats := b.Stats()
	resp.Keys = make([][]byte, 0, stats.KeyN)
	resp.Values = make([][]byte, 0, stats.KeyN)

	return b.ForEach(func(k, v []byte) error {
		resp.Keys = append(resp.Keys, k)
		resp.Values = append(resp.Values, v)
		return nil
	})
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
