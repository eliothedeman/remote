package remote

import "github.com/boltdb/bolt"

var (
	d bolt.DB
)

// DB provides an interface to a boltdb
type DB interface {
	Close() error
	View(func(tx Tx) error) error
	Update(func(tx Tx) error) error
}

// Open returns a new view to a database
func Open(path string, options *bolt.Options) (DB, error) {
	return newLocalClient(path, options)
}