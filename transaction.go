package remote

import "github.com/boltdb/bolt"

var (
	t bolt.Tx
	b bolt.Bucket
)

// Tx is a transaction
type Tx interface {
	Bucket(name []byte) Bucket
	Commit() error
	CreateBucket(name []byte) (Bucket, error)
	CreateBucketIfNotExists(name []byte) (Bucket, error)
	DB() DB
	DeleteBucket(name []byte) error
	OnCommit(h func())
	Rollback() error
}

// RTx is a local transaction.
type RTx struct {
	r              *RClient
	contextID      uint64
	commitHandlers []func()
}

// OnCommit adds a handler to be called when a transaction is commited.
func (r *RTx) OnCommit(h func()) {
	if r.commitHandlers == nil {
		r.commitHandlers = []func(){h}
		return
	}
	r.commitHandlers = append(r.commitHandlers, h)
}

// Bucket returns the bucket with the given name.
func (r *RTx) Bucket(name []byte) Bucket {
	req := &BucketRequest{}
	req.Key = name
	resp := &BucketResponse{}
	req.ContextID = r.contextID
	err := r.r.call("srv.Bucket", req, resp)
	if err != nil {
		return nil
	}
	b := &RBucket{}
	b.tx = r
	b.r = r.r
	b.id = resp.BucketID
	b.parent = resp.BucketContextID
	return b
}

// DB returns the database that this transaction is from.
func (r *RTx) DB() DB {
	return r.r
}

// Commit this transaction.
func (r *RTx) Commit() error {
	err := r.r.commit(r.contextID)
	for i := range r.commitHandlers {
		r.commitHandlers[i]()
	}
	return err
}

// Rollback this transaction.
func (r *RTx) Rollback() error {
	return r.r.rollback(r.contextID)
}

// CreateBucket creates and returns a new bucket.
func (r *RTx) CreateBucket(name []byte) (Bucket, error) {
	req := &BucketRequest{}
	req.Key = name
	resp := &BucketResponse{}
	req.ContextID = r.contextID
	err := r.r.call("srv.CreateBucket", req, resp)
	b := &RBucket{}
	b.tx = r
	b.r = r.r
	b.id = resp.BucketID
	b.parent = resp.BucketContextID
	return b, err
}

// CreateBucketIfNotExists creates and returns a new bucket.
func (r *RTx) CreateBucketIfNotExists(name []byte) (Bucket, error) {
	req := &BucketRequest{}
	req.Key = name
	resp := &BucketResponse{}
	req.ContextID = r.contextID
	err := r.r.call("srv.CreateBucketIfNotExists", req, resp)
	b := &RBucket{}
	b.tx = r
	b.r = r.r
	b.id = resp.BucketID
	b.parent = resp.BucketContextID
	return b, err
}

// DeleteBucket creates and returns a new bucket.
func (r *RTx) DeleteBucket(name []byte) error {
	req := &BucketRequest{}
	req.Key = name
	req.ContextID = r.contextID
	resp := &BucketResponse{}
	return r.r.call("srv.DeleteBucket", req, resp)
}

// LTx is a local transaction.
type LTx struct {
	tx *bolt.Tx
	db DB
}

// Bucket returns the bucket with the given name
func (l *LTx) Bucket(name []byte) Bucket {
	b := l.tx.Bucket(name)
	if b == nil {
		return nil
	}

	return &LBucket{
		b:  b,
		tx: l,
	}

}

// OnCommit adds a handelr to be called after commit is called.
func (l *LTx) OnCommit(h func()) {
	l.tx.OnCommit(h)
}

// Commit this transaction.
func (l *LTx) Commit() error {
	return l.tx.Commit()
}

// Rollback this transaction.
func (l *LTx) Rollback() error {
	return l.tx.Rollback()
}

// DB returns the database that this transaction is from.
func (l *LTx) DB() DB {
	return l.db
}

// CreateBucket creates and returns a new bucket.
func (l *LTx) CreateBucket(name []byte) (Bucket, error) {
	b, err := l.tx.CreateBucket(name)
	if err != nil {
		return nil, err
	}

	return &LBucket{
		b:  b,
		tx: l,
	}, nil
}

// CreateBucketIfNotExists creates and returns a new bucket.
func (l *LTx) CreateBucketIfNotExists(name []byte) (Bucket, error) {
	b, err := l.tx.CreateBucketIfNotExists(name)
	if err != nil {
		return nil, err
	}

	return &LBucket{
		b:  b,
		tx: l,
	}, nil
}

// DeleteBucket creates and returns a new bucket.
func (l *LTx) DeleteBucket(name []byte) error {
	return l.tx.DeleteBucket(name)
}
