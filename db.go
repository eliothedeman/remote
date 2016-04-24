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
	Path() string
	View(func(tx Tx) error) error
	Update(func(tx Tx) error) error
	Stats() bolt.Stats
}

// Open returns a new view to a database.
func Open(path string, options *bolt.Options) (DB, error) {
	if strings.HasPrefix(path, "tcp://") {

		// don't include the "tcp" part of the address
		return dialRemoteClient(path[6:])
	}
	return newLocalClient(path, options)
}
