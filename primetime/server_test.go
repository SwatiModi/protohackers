package primetime_test

import (
	"bufio"
	"net"
	"protohackers/primetime"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestSmokeServer(t *testing.T) {
	t.Run("basic case : not prime number", func(t *testing.T) {
		// start server
		go primetime.StartServer()

		// make a request to localhost 8000
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fatal("failed to start server")
		}

		defer conn.Close()

		msg := "{\"method\":\"isPrime\",\"number\":123}\n"
		expectedMsg := "{\"method\":\"isPrime\", \"prime\":false}\n"

		if _, err := conn.Write([]byte(msg)); err != nil {
			t.FailNow()
		}

		resp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.FailNow()
		}

		assert.Equal(t, expectedMsg, resp)
	})

	t.Run("basic case : not prime number 2", func(t *testing.T) {
		// make a request to localhost 8000
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fatal("failed to start server")
		}

		defer conn.Close()

		msg := "{\"method\":\"isPrime\",\"number\":46}\n"
		expectedMsg := "{\"method\":\"isPrime\", \"prime\":false}\n"

		if _, err := conn.Write([]byte(msg)); err != nil {
			t.FailNow()
		}

		resp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.FailNow()
		}

		assert.Equal(t, expectedMsg, resp)
	})

	t.Run("basic case : prime number", func(t *testing.T) {
		// make a request to localhost 8000
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fatal("failed to start server")
		}

		defer conn.Close()

		msg := "{\"number\":43,  \"method\":\"isPrime\"}\n"
		expectedMsg := "{\"method\":\"isPrime\", \"prime\":true}\n"

		if _, err := conn.Write([]byte(msg)); err != nil {
			t.FailNow()
		}

		resp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.FailNow()
		}

		assert.Equal(t, expectedMsg, resp)
	})

	t.Run("basic case : prime number 2", func(t *testing.T) {
		// make a request to localhost 8000
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fatal("failed to start server")
		}

		defer conn.Close()

		msg := "{\"method\":\"isPrime\",\"number\":431}\n"
		expectedMsg := "{\"method\":\"isPrime\", \"prime\":true}\n"

		if _, err := conn.Write([]byte(msg)); err != nil {
			t.FailNow()
		}

		resp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.FailNow()
		}

		assert.Equal(t, expectedMsg, resp)
	})

	t.Run("basic case : malformed request", func(t *testing.T) {
		// make a request to localhost 8000
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fatal("failed to start server")
		}

		defer conn.Close()

		msg := "{\"method\":\"isPrime\",\"numbers\":47}\n"
		expectedMsg := "{\"method\":\"isPrime\", \"prime\":true}\n"

		if _, err := conn.Write([]byte(msg)); err != nil {
			t.FailNow()
		}

		resp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.FailNow()
		}

		assert.NotEqual(t, expectedMsg, resp)
	})
}
