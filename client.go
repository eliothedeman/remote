package remote

import (
	"io"
	"net"
	"net/rpc"
	"time"

	"github.com/boltdb/bolt"
)

// RClient is a remote view to the remote server.
type RClient struct {
	conn io.ReadWriteCloser
	c    *rpc.Client
}

func (r *RClient) call(name string, args, resp interface{}) error {
	return r.c.Call(name, args, resp)
}

func dialRemoteClient(host string) (*RClient, error) {
	c, err := net.DialTimeout("tcp4", host, time.Second*1)

	return &RClient{
		conn: c,
		c:    rpc.NewClient(c),
	}, err
}
func (r *RClient) genTx() *RTx {

	return &RTx{
		r: r,
	}
}

func (r *RClient) begin(write bool) (*RTx, error) {
	resp := &BeginTransactionResponse{}
	req := BeginTransactionRequest{}
	req.Writable = write

	err := r.call("srv.BeginTransaction", req, resp)

	return &RTx{
		r:         r,
		contextID: resp.ContextID,
	}, err
}

// Begin a transaction.
func (r *RClient) Begin(writeable bool) (Tx, error) {
	return r.begin(writeable)
}

func (r *RClient) commit(id uint64) error {
	resp := &CommitTransactionResponse{}
	return r.call("srv.CommitTransaction", id, resp)
}

func (r *RClient) rollback(id uint64) error {
	resp := &RollbackTransactionResponse{}
	return r.call("srv.RollbackTransaction", id, resp)
}

// IsReadOnly or not?
func (r *RClient) IsReadOnly() bool {
	return false
}

// Close the connection to the database.
func (r *RClient) Close() error {
	return r.c.Close()
}

// View opens a read only transaction to the database.
func (r *RClient) View(f func(tx Tx) error) error {
	t, err := r.begin(false)
	if err != nil {
		return err
	}

	return f(t)
}

// Update opens a read/write transaction to the database.
func (r *RClient) Update(f func(tx Tx) error) error {
	t, err := r.begin(true)
	if err != nil {
		return err
	}

	err = f(t)
	if err != nil {
		rErr := r.rollback(t.contextID)
		if rErr != nil {
			return rErr
		}
		return err
	}
	return r.commit(t.contextID)
}

// LClient is a local view to a boltdb.
type LClient struct {
	*bolt.DB
}

func newLocalClient(path string, options *bolt.Options) (*LClient, error) {
	db, err := bolt.Open(path, 0777, options)
	return &LClient{
		DB: db,
	}, err
}

// Begin a transaction.
func (l *LClient) Begin(writeable bool) (Tx, error) {
	tx, err := l.DB.Begin(writeable)
	return l.genTx(tx), err
}

// Close the underlying database.
func (l *LClient) Close() error {
	return l.DB.Close()
}

// View opens the database for read only operation.
func (l *LClient) View(f func(Tx) error) error {
	return l.DB.View(func(t *bolt.Tx) error {
		tx := l.genTx(t)
		return f(tx)
	})
}

// Update opens the database for read/write operations.
func (l *LClient) Update(f func(Tx) error) error {
	return l.DB.Update(func(t *bolt.Tx) error {
		tx := l.genTx(t)
		return f(tx)
	})
}

func (l *LClient) genTx(tx *bolt.Tx) Tx {
	return &LTx{
		tx: tx,
		db: l,
	}
}
