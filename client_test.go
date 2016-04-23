package remote

import "testing"

func TestLocalOpen(t *testing.T) {
	db, err := Open("/tmp/test.db", nil)
	if err != nil {
		t.Error(err)
	}

	err = db.Update(func(tx Tx) error {
		b, verr := tx.CreateBucketIfNotExists([]byte("hello"))
		if verr != nil {
			return verr
		}

		return b.Put([]byte("world"), []byte("world"))
	})
	if err != nil {
		t.Error(err)
	}

	db.Close()
}
