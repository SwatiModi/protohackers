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
}
