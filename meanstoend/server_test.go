package meanstoend_test

import (
	"bufio"
	"net"
	"protohackers/meanstoend"
	"testing"
)

func TestServer(t *testing.T) {
	t.Run("basic case test", func(t *testing.T) {
		go meanstoend.StartServer()

		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fail()
		}

		if _, err := conn.Write([]byte("hello")); err != nil {
			t.Error("write", err)
		}

		resp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.FailNow()
		}

		t.Log(resp)

	})
}
