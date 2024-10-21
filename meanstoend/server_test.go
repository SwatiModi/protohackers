package meanstoend_test

import (
	"io"
	"net"
	"protohackers/meanstoend"
	"testing"

	"github.com/go-playground/assert/v2"
)

var insertInput = []byte{0x49, 0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x65}
var queryInput = []byte{0x51, 0x00, 0x00, 0x03, 0xe8, 0x00, 0x01, 0x86, 0xa0}
var invalidInput = []byte{0x55, 0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x65}
var outofTimeRangeInput = []byte{0x51, 0x00, 0x00, 0x00, 0xe8, 0x00, 0x01, 0x00, 0xa0}

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

		// t.Log(resp)

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

func TestDecoding(t *testing.T) {
	t.Run("decode insert request", func(t *testing.T) {

		rt, timestamp, price := meanstoend.DecodeInput(insertInput)

		assert.Equal(t, meanstoend.I, rt)
		assert.Equal(t, uint32(12345), timestamp)
		assert.Equal(t, uint32(101), price)
	})

	t.Run("decode query request", func(t *testing.T) {

		rt, minTime, maxTime := meanstoend.DecodeInput(queryInput)

		assert.Equal(t, meanstoend.Q, rt)
		assert.Equal(t, uint32(1000), minTime)
		assert.Equal(t, uint32(100000), maxTime)
	})

	t.Run("invalid request type", func(t *testing.T) {
		rt, _, _ := meanstoend.DecodeInput(invalidInput)

		assert.Equal(t, meanstoend.Invalid, rt)
	})
}
