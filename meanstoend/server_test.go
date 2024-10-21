package meanstoend_test

import (
	"bufio"
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

	t.Run("valid insert op", func(t *testing.T) {
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fail()
		}

		if _, err := conn.Write(insertInput); err != nil {
			t.Error("write", err)
		}

		resp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.FailNow()
		}

		t.Log(resp)

		assert.Equal(t, "", resp)
	})

	t.Run("valid query op", func(t *testing.T) {
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fail()
		}

		if _, err := conn.Write(insertInput); err != nil {
			t.Error("write", err)
		}

		resp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.FailNow()
		}

		if _, err := conn.Write(insertInput); err != nil {
			t.Error("write", err)
		}

		resp2, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.FailNow()
		}

		t.Log(resp)

		assert.Equal(t, "", resp2)
	})

	t.Run("invalid query op", func(t *testing.T) {

	})

	t.Run("insert with conn a, query with conn a", func(t *testing.T) {

	})

	t.Run("insert with conn a, query with conn b", func(t *testing.T) {

	})

	t.Run("out of range query, no results", func(t *testing.T) {

	})
}

func TestDecoding(t *testing.T) {
	t.Run("decode insert request", func(t *testing.T) {

		rType, timestamp, price := meanstoend.DecodeInput(insertInput)

		rt := meanstoend.ParseRequestType(rType)
		t.Log(rt, timestamp, price)

		assert.Equal(t, meanstoend.RequestType(1), rt)
		assert.Equal(t, 12345, timestamp)
		assert.Equal(t, 101, price)
	})

	t.Run("decode query request", func(t *testing.T) {

		rType, minTime, maxTime := meanstoend.DecodeInput(queryInput)

		rt := meanstoend.ParseRequestType(rType)
		t.Log(rt, minTime, maxTime)

		assert.Equal(t, meanstoend.RequestType(2), rt)
		assert.Equal(t, 1000, minTime)
		assert.Equal(t, 10000, maxTime)
	})

	t.Run("invalid request type", func(t *testing.T) {

		rType, minTime, maxTime := meanstoend.DecodeInput(invalidInput)

		rt := meanstoend.ParseRequestType(rType)
		t.Log(rt, minTime, maxTime)

		assert.Equal(t, meanstoend.RequestType(0), rt)
	})
}
