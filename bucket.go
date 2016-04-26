package remote

import (
	"log"

	"github.com/boltdb/bolt"
)

// Bucket is an interface into a bolt bucket.
type Bucket interface {
	Bucket(name []byte) Bucket
	CreateBucket(name []byte) (Bucket, error)
	CreateBucketIfNotExists(name []byte) (Bucket, error)
	DeleteBucket(name []byte) error
	Get(key []byte) []byte
	Put(key, value []byte) error
	ForEach(func(k, v []byte) error) error
	NextSequence() (uint64, error)
	Stats() bolt.BucketStats
	Tx() Tx
	Writeable() bool
}

// RBucket is a remote bucket.
type RBucket struct {
	r      *RClient
	tx     *RTx
	id     uint64
	parent uint64
}

// Writeable will always be true.
func (r *RBucket) Writeable() bool {
	return true
}

// ForEach k,v pair
// TODO this needs pagination.
func (r *RBucket) ForEach(f func(k, v []byte) error) error {
	req := &BucketForEachRequest{}
	resp := &BucketForEachResponse{}
	req.BucketID = r.id
	req.ContextID = r.parent

	err := r.r.call("srv.BucketForEachStart", req, nil)
	if err != nil {
		return err
	}
	for {
		err = r.r.call("srv.BucketForEachNext", req, resp)
		if err != nil {
			return r.r.call("srv.BucketForEachStop", req, nil)
		}

		err = f(resp.Key, resp.Value)
		if err != nil {
			return r.r.call("srv.BucketForEachStop", req, nil)
		}

		if resp.Index >= resp.Size-1 {
			break
		}
	}

	return nil
}

// Tx returns the transaction this bucket is a part of.
func (r *RBucket) Tx() Tx {
	return r.tx
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
	b.tx = r.tx
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
	b.tx = r.tx
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
	b.tx = r.tx
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

// Stats returns the stats of a bucket.
func (r *RBucket) Stats() bolt.BucketStats {
	resp := &BucketStatsResponse{}
	err := r.r.call("srv.BucketStats", Empty{}, resp)
	if err != nil {
		log.Println(err)
	}

	return resp.BucketStats
}

// NextSequence returns the next unique id for this bucket.
func (r *RBucket) NextSequence() (uint64, error) {
	var resp uint64
	req := &NextSequenceRequest{}
	req.BucketID = r.id
	req.ContextID = r.parent

	err := r.r.call("srv.NextSequence", req, &resp)
	return resp, err
}

// LBucket is a local bucket
type LBucket struct {
	tx *LTx
	b  *bolt.Bucket
}

// ForEach a func on every k,v.
func (l *LBucket) ForEach(f func(k, v []byte) error) error {
	return l.b.ForEach(f)
}

// Writeable is this a read only bucket?
func (l *LBucket) Writeable() bool {
	return l.b.Writable()
}

// Tx returns the parent transaction of the bucket.
func (l *LBucket) Tx() Tx {
	return l.tx
}

// Stats returns the stats about this bucket.
func (l *LBucket) Stats() bolt.BucketStats {
	return l.b.Stats()
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

// NextSequence inserts the given value at the givne key.
func (l *LBucket) NextSequence() (uint64, error) {
	return l.b.NextSequence()
}
