package meanstoend_test

import (
	"io"
	"net"
	"protohackers/meanstoend"
	"testing"

	"github.com/go-playground/assert/v2"
)

var insertInput = []byte{0x49, 0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x65}
var queryInput = []byte{0x51, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x40, 0x00}
var invalidInput = []byte{0x55, 0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x65}

var input1 = []byte{0x49, 0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x65} // I 12345 101
var input2 = []byte{0x49, 0x00, 0x00, 0x30, 0x3a, 0x00, 0x00, 0x00, 0x66} // I 12346 102
var input3 = []byte{0x49, 0x00, 0x00, 0x30, 0x3b, 0x00, 0x00, 0x00, 0x64} // I 12347 100
var input4 = []byte{0x49, 0x00, 0x00, 0xa0, 0x00, 0x00, 0x00, 0x00, 0x05} // I 40960 5

// var outofTimeRangeInput = []byte{0x51, 0x00, 0x00, 0x00, 0xe8, 0x00, 0x01, 0x00, 0xa0}

func TestServer(t *testing.T) {
	go meanstoend.StartServer()

	t.Run("invalid insert op", func(t *testing.T) {
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fail()
		}

		if _, err := conn.Write([]byte("foo=bar")); err != nil {
			t.Error("write", err)
		}

		conn.Close()
	})

	t.Run("valid insert op2", func(t *testing.T) {
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fail()
		}

		if _, err := conn.Write(insertInput); err != nil {
			t.Error("write", err)
		}

		conn.Close()
	})

	t.Run("valid query op", func(t *testing.T) {
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fail()
		}

		if _, err := conn.Write(input1); err != nil {
			t.Error("write", err)
		}
		if _, err := conn.Write(input2); err != nil {
			t.Error("write", err)
		}
		if _, err := conn.Write(input3); err != nil {
			t.Error("write", err)
		}
		if _, err := conn.Write(input4); err != nil {
			t.Error("write", err)
		}

		if _, err := conn.Write(queryInput); err != nil {
			t.Error("write", err)
		}

		resp := make([]byte, 4)
		if _, err := io.ReadFull(conn, resp); err != nil {
			t.FailNow()
		}

		t.Log(resp)

		assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x65}, resp)
		conn.Close()
	})

	t.Run("invalid query op", func(t *testing.T) {
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fail()
		}

		if _, err := conn.Write(invalidInput); err != nil {
			t.Error("write", err)
		}

		conn.Close()
	})

	t.Run("insert with conn a, query with conn b", func(t *testing.T) {
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fail()
		}

		if _, err := conn.Write(insertInput); err != nil {
			t.Error("write", err)
		}

		if _, err := conn.Write(queryInput); err != nil {
			t.Error("write", err)
		}

		resp := make([]byte, 4)
		if _, err := io.ReadFull(conn, resp); err != nil {
			t.FailNow()
		}

		assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x65}, resp)

		conn2, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fail()
		}

		if _, err := conn2.Write(queryInput); err != nil {
			t.Error("write", err)
		}

		resp2 := make([]byte, 4)
		if _, err := io.ReadFull(conn2, resp2); err != nil {
			t.FailNow()
		}

		assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00}, resp2)

		conn.Close()
		conn2.Close()
	})

	t.Run("out of range query, no results", func(t *testing.T) {

	})
}
