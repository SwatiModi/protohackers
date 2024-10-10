package smoketest

import (
	"bufio"
	"net"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestSmokeServer(t *testing.T) {
	t.Run("basic case", func(t *testing.T) {
		// start server
		go StartServer()

		// make a request to localhost 8000
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fatal("failed to start server")
		}

		defer conn.Close()

		msg := "hello World!\n"
		if _, err := conn.Write([]byte(msg)); err != nil {
			t.FailNow()
		}

		resp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.FailNow()
		}

		assert.Equal(t, msg, resp)
	})
}
