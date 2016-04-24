package remote

import (
	"io"
	"net/rpc"
	"os"
	"testing"
	"time"

	"github.com/eliothedeman/randutil"
	"github.com/ugorji/go/codec"
)

type testReadWriteCloser struct {
	io.Reader
	io.Writer
}

func (t *testReadWriteCloser) Close() error {
	return nil
}

func newRWC(r io.Reader, w io.Writer) io.ReadWriteCloser {
	return &testReadWriteCloser{
		Reader: r,
		Writer: w,
	}
}

func testConn() (src, dst io.ReadWriteCloser) {
	a, b := io.Pipe()
	c, d := io.Pipe()

	return newRWC(a, d), newRWC(c, b)
}

func run(f func(*Server, *RClient)) {
	s, err := OpenServer("/tmp/test.db")
	if err != nil {
		panic(err)
	}

	src, dst := testConn()

	go s.ServeConn(dst)
	cod := codec.GoRpc.ClientCodec(src, &codec.MsgpackHandle{})

	c := &RClient{
		conn: src,
		c:    rpc.NewClientWithCodec(cod),
	}

	f(s, c)
	os.Remove(s.db.Path())
}

func TestServe(t *testing.T) {
	run(func(s *Server, c *RClient) {
		resp := &PingResponse{}
		c.call("srv.Ping", &PingRequest{T: time.Now()}, resp)

		if resp.RoundTrip() > time.Millisecond {
			t.Fail()
		}
	})
}

func TestView(t *testing.T) {
	run(func(s *Server, c *RClient) {
		err := c.View(func(tx Tx) error {
			return nil
		})
		if err != nil {
			t.Error(err)
		}
	})
}

func TestCreateBucket(t *testing.T) {
	run(func(s *Server, c *RClient) {
		err := c.Update(func(tx Tx) error {
			b, err := tx.CreateBucket([]byte("hello"))
			if b == nil {
				t.Fail()
			}
			if b == nil {
				t.Fail()
			}
			_, err = tx.CreateBucket([]byte("hello"))
			if err == nil {
				t.Fail()
			}
			return nil
		})
		if err != nil {
			t.Error(err)
		}
	})
}

func TestCreateBucketIfNotExists(t *testing.T) {
	run(func(s *Server, c *RClient) {
		err := c.Update(func(tx Tx) error {
			b, err := tx.CreateBucket([]byte("hello"))
			if b == nil {
				t.Fail()
			}
			if b == nil {
				t.Fail()
			}
			_, err = tx.CreateBucketIfNotExists([]byte("hello"))
			if err != nil {
				t.Error(err)
			}

			return nil
		})
		if err != nil {
			t.Error(err)
		}
	})
}

func TestGetPut(t *testing.T) {
	run(func(s *Server, c *RClient) {
		err := c.Update(func(tx Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("hello"))
			if err != nil {
				t.Error(err)
				return err
			}

			err = b.Put([]byte("hello"), []byte("world"))
			if err != nil {
				t.Error(err)
				return err
			}

			out := b.Get([]byte("hello"))
			if out == nil {
				t.Fail()
				return nil
			}
			if string(out) != "world" {
				t.Error(string(out))
			}
			return nil

		})
		if err != nil {
			t.Error(err)
		}
	})
}
func TestBucketForEach(t *testing.T) {
	run(func(s *Server, c *RClient) {
		err := c.Update(func(tx Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("hello"))
			if err != nil {
				t.Error(err)
				return err
			}

			for i := 0; i < 100; i++ {
				b.Put(randutil.Bytes(10), randutil.Bytes(10))
			}
			return nil
		})

		if err != nil {
			t.Error(err)
		}
		err = c.View(func(tx Tx) error {
			b := tx.Bucket([]byte("hello"))

			var count int
			err = b.ForEach(func(k, v []byte) error {
				count++
				return nil
			})

			if count != 100 {
				t.FailNow()
			}

			return err

		})

		if err != nil {
			t.Error(err)
		}
	})
}

func TestCommit(t *testing.T) {
	run(func(s *Server, c *RClient) {
		err := c.Update(func(tx Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("hello"))
			if err != nil {
				t.Error(err)
				return err
			}

			err = b.Put([]byte("hello"), []byte("world"))
			if err != nil {
				t.Error(err)
				return err
			}
			return nil
		})
		if err != nil {
			t.Error(err)

		}
		err = c.View(func(tx Tx) error {
			b := tx.Bucket([]byte("hello"))
			if err != nil {
				t.Error(err)
				return err
			}

			out := b.Get([]byte("hello"))
			if string(out) != "world" {
				t.Fail()
			}

			return nil
		})
		if err != nil {
			t.Error(err)

		}

	})
}

func BenchmarkPut(b *testing.B) {
	b.ReportAllocs()
	run(func(s *Server, c *RClient) {
		err := c.Update(func(tx Tx) error {
			bt, err := tx.CreateBucketIfNotExists([]byte("hello"))

			for i := 0; i < b.N; i++ {
				bt.Put([]byte("hello"), []byte("world"))
			}
			return err
		})
		if err != nil {
			b.Error(err)
		}
	})

}
