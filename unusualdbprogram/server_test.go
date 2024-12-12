package unusualdbprogram_test

import (
	"net"
	"protohackers/unusualdbprogram"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestServer(t *testing.T) {
	// start server
	go unusualdbprogram.StartServer()

	t.Run("version check", func(t *testing.T) {
		// send version request
		pc, err := net.Dial("udp", ":8000")
		if err != nil {
			t.Fatal("dial", err)
		}

		_, err = pc.Write([]byte("version"))
		if err != nil {
			t.Fatal("write", err)
		}

		// read response
		buff := make([]byte, 1024)
		n, err := pc.Read(buff)
		if err != nil {
			t.Fatal("read", err)
		}

		assert.Equal(t, "version=0.0.1", string(buff[:n]))
	})

	t.Run("insert and retrieve", func(t *testing.T) {
		// send insert request
		pc, err := net.Dial("udp", ":8000")
		if err != nil {
			t.Fatal("dial", err)
		}

		_, err = pc.Write([]byte("foo=bar"))
		if err != nil {
			t.Fatal("write", err)
		}

		// send retrieve request
		_, err = pc.Write([]byte("foo"))
		if err != nil {
			t.Fatal("write", err)
		}

		buff := make([]byte, 1024)
		// read response
		n, err := pc.Read(buff)
		if err != nil {
			t.Fatal("read", err)
		}

		assert.Equal(t, "foo=bar", string(buff[:n]))
	})

	t.Run("insert and retrieve 2", func(t *testing.T) {
		// send insert request
		pc, err := net.Dial("udp", ":8000")
		if err != nil {
			t.Fatal("dial", err)
		}

		_, err = pc.Write([]byte("foo=bar=baz"))
		if err != nil {
			t.Fatal("write", err)
		}

		// send retrieve request
		_, err = pc.Write([]byte("foo"))
		if err != nil {
			t.Fatal("write", err)
		}

		buff := make([]byte, 1024)
		// read response
		n, err := pc.Read(buff)
		if err != nil {
			t.Fatal("read", err)
		}

		assert.Equal(t, "foo=bar=baz", string(buff[:n]))
	})

	t.Run("insert and retrieve 3", func(t *testing.T) {
		// send insert request
		pc, err := net.Dial("udp", ":8000")
		if err != nil {
			t.Fatal("dial", err)
		}

		_, err = pc.Write([]byte("foo="))
		if err != nil {
			t.Fatal("write", err)
		}

		// send retrieve request
		_, err = pc.Write([]byte("foo"))
		if err != nil {
			t.Fatal("write", err)
		}

		buff := make([]byte, 1024)
		// read response
		n, err := pc.Read(buff)
		if err != nil {
			t.Fatal("read", err)
		}

		assert.Equal(t, "foo=", string(buff[:n]))
	})

	t.Run("insert and retrieve 4", func(t *testing.T) {
		// send insert request
		pc, err := net.Dial("udp", ":8000")
		if err != nil {
			t.Fatal("dial", err)
		}

		_, err = pc.Write([]byte("foo==="))
		if err != nil {
			t.Fatal("write", err)
		}

		// send retrieve request
		_, err = pc.Write([]byte("foo"))
		if err != nil {
			t.Fatal("write", err)
		}

		buff := make([]byte, 1024)
		// read response
		n, err := pc.Read(buff)
		if err != nil {
			t.Fatal("read", err)
		}

		assert.Equal(t, "foo===", string(buff[:n]))
	})

	t.Run("insert and retrieve 5", func(t *testing.T) {
		// send insert request
		pc, err := net.Dial("udp", ":8000")
		if err != nil {
			t.Fatal("dial", err)
		}

		_, err = pc.Write([]byte("=foo"))
		if err != nil {
			t.Fatal("write", err)
		}

		// send retrieve request
		_, err = pc.Write([]byte(""))
		if err != nil {
			t.Fatal("write", err)
		}

		buff := make([]byte, 1024)
		// read response
		n, err := pc.Read(buff)
		if err != nil {
			t.Fatal("read", err)
		}

		assert.Equal(t, "=foo", string(buff[:n]))
	})

	t.Run("insert and retrieve 6", func(t *testing.T) {
		// send insert request
		pc, err := net.Dial("udp", ":8000")
		if err != nil {
			t.Fatal("dial", err)
		}

		_, err = pc.Write([]byte("message=Hello,world!"))
		if err != nil {
			t.Fatal("write", err)
		}

		// send retrieve request
		_, err = pc.Write([]byte("message"))
		if err != nil {
			t.Fatal("write", err)
		}

		buff := make([]byte, 1024)
		// read response
		n, err := pc.Read(buff)
		if err != nil {
			t.Fatal("read", err)
		}

		assert.Equal(t, "message=Hello,world!", string(buff[:n]))
	})
}
