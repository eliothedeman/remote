package remote

import (
	"strings"

	"github.com/boltdb/bolt"
)

var (
	d bolt.DB
)

// DB provides an interface to a boltdb
type DB interface {
	Begin(writeable bool) (Tx, error)
	Close() error
	View(func(tx Tx) error) error
	Update(func(tx Tx) error) error
}

// Open returns a new view to a database.
func Open(path string, options *bolt.Options) (DB, error) {
	if strings.HasPrefix(path, "tcp://") {
		return dialRemoteClient(path)
	}
	return newLocalClient(path, options)
}
