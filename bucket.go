package remote

import "github.com/boltdb/bolt"

// Bucket is an interface into a bolt bucket.
type Bucket interface {
	Bucket(name []byte) Bucket
	CreateBucket(name []byte) (Bucket, error)
	CreateBucketIfNotExists(name []byte) (Bucket, error)
	DeleteBucket(name []byte) error
	Get(key []byte) []byte
	Put(key, value []byte) error
}

// RBucket is a remote bucket.
type RBucket struct {
	r      *RClient
	id     uint64
	parent uint64
}

// Bucket returns the bucket with the given name
func (r *RBucket) Bucket(name []byte) Bucket {
	req := &BucketRequest{}
	req.Key = name
	resp := &BucketResponse{}
	req.ContextID = r.parent
	err := r.r.call("srv.Bucket", req, resp)
	if err != nil {
		return nil
	}
	b := &RBucket{}
	b.r = r.r
	b.id = resp.BucketID
	b.parent = resp.BucketContextID
	return b
}

// CreateBucket creates and returns a new bucket.
func (r *RBucket) CreateBucket(name []byte) (Bucket, error) {
	req := &BucketRequest{}
	req.Key = name
	resp := &BucketResponse{}
	req.ContextID = r.parent
	err := r.r.call("srv.CreateBucket", req, resp)
	b := &RBucket{}
	b.r = r.r
	b.id = resp.BucketID
	b.parent = resp.BucketContextID
	return b, err
}

// CreateBucketIfNotExists creates and returns a new bucket.
func (r *RBucket) CreateBucketIfNotExists(name []byte) (Bucket, error) {
	req := &BucketRequest{}
	req.Key = name
	resp := &BucketResponse{}
	req.ContextID = r.parent
	err := r.r.call("srv.CreateBucketIfNotExists", req, resp)
	b := &RBucket{}
	b.r = r.r
	b.id = resp.BucketID
	b.parent = resp.BucketContextID
	return b, err
}

// DeleteBucket creates and returns a new bucket.
func (r *RBucket) DeleteBucket(name []byte) error {
	req := &BucketRequest{}
	req.Key = name
	req.ContextID = r.parent
	resp := &BucketResponse{}
	return r.r.call("srv.DeleteBucket", req, resp)
}

// Get returns the value of the given key.
func (r *RBucket) Get(key []byte) []byte {
	req := &GetReqeust{}
	resp := &GetResponse{}
	req.Key = key
	req.BucketID = r.id
	req.ContextID = r.parent
	err := r.r.call("srv.Get", req, resp)
	if err != nil {
		return nil
	}
	return resp.Val
}

// Put inserts the given value at the given key.
func (r *RBucket) Put(key, value []byte) error {
	req := &PutReqeust{}
	resp := &PutResponse{}
	req.Key = key
	req.Val = value
	req.BucketID = r.id
	req.ContextID = r.parent
	return r.r.call("srv.Put", req, resp)
}

// LBucket is a local bucket
type LBucket struct {
	b *bolt.Bucket
}

// Bucket returns the bucket with the given name if it exists.
func (l *LBucket) Bucket(name []byte) Bucket {
	return &LBucket{
		b: l.b.Bucket(name),
	}
}

// CreateBucket creats a new bucket with the given name.
func (l *LBucket) CreateBucket(name []byte) (Bucket, error) {
	b, err := l.b.CreateBucket(name)

	return &LBucket{
		b: b,
	}, err
}

// CreateBucketIfNotExists creats a new bucket with the given name.
func (l *LBucket) CreateBucketIfNotExists(name []byte) (Bucket, error) {
	b, err := l.b.CreateBucketIfNotExists(name)

	return &LBucket{
		b: b,
	}, err
}

// DeleteBucket creats a new bucket with the given name.
func (l *LBucket) DeleteBucket(name []byte) error {
	return l.b.DeleteBucket(name)
}

// Get returns the value stored at the given key.
func (l *LBucket) Get(key []byte) []byte {
	return l.b.Get(key)
}

// Put inserts the given value at the givne key.
func (l *LBucket) Put(key, value []byte) error {
	return l.b.Put(key, value)
}
