package remote

import (
	"errors"
	"io"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"

	"github.com/boltdb/bolt"
)

type parent interface {
	Bucket(key []byte) *bolt.Bucket
	CreateBucket(key []byte) (*bolt.Bucket, error)
	CreateBucketIfNotExists(key []byte) (*bolt.Bucket, error)
	DeleteBucket(key []byte) error
}

// A Server runs a remote view to a boltdb.
type Server struct {
	db       *bolt.DB
	contexts map[uint64]*Context
	cCount   uint64
	sync.RWMutex
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

	srv := rpc.NewServer()

	srv.RegisterName("srv", s)
	srv.Accept(l)
	return nil
}

// ServeConn will serve the rpc via a io.ReadWriteCloser
func (s *Server) ServeConn(conn io.ReadWriteCloser) error {
	srv := rpc.NewServer()
	srv.RegisterName("srv", s)
	srv.ServeConn(conn)
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
	sync.RWMutex
	id      uint64
	tx      *bolt.Tx
	parent  parent
	bCount  uint64
	buckets map[uint64]*bolt.Bucket
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

// DeleteBucket creates a new transaction with the given mode.
func (s *Server) DeleteBucket(req BucketRequest, resp *BucketResponse) error {
	c := s.getContext(req.ContextID)
	if c == nil {
		return errors.New("Transaction not found")
	}

	// TODO this needs to clear this bucket out of the context
	err := c.parent.DeleteBucket(req.Key)
	return err
}

// CreateBucket creates a new transaction with the given mode.
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

// CreateBucketIfNotExists creates a new transaction with the given mode.
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

// Get creates a new transaction with the given mode.
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

// Put creates a new transaction with the given mode.
func (s *Server) Put(req *PutReqeust, resp *PutResponse) error {
	c := s.getContext(req.ContextID)
	b := c.getBucket(req.BucketID)

	if b == nil {
		return errors.New("Bucket not found.")
	}

	return b.Put(req.Key, req.Val)
}
