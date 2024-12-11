package budgetchat_test

import (
	"bufio"
	"log"
	"net"
	"protohackers/budgetchat"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestStartServer(t *testing.T) {
	t.Run("StartServer", func(t *testing.T) {
		go budgetchat.StartServer()

		// send a connection request to the server : user 1
		// wait for server msg asking for username : user 1
		// send username : user 1
		// verify welcome msg : user 1

		// send a connection request to the server : user 2
		// wait for server msg asking for username : user 2
		// send username : user 2
		// verify welcome msg : user 2
		// send chat msg : user 2
		// verify chat msg :	user 2

		// code
		conn1, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fail()
		}

		t.Log("connected to server")

		reader := bufio.NewReader(conn1)
		msg, _ := reader.ReadString('\n')

		assert.Equal(t, "Welcome to budgetchat! What shall I call you?\n", msg)

		if _, err := conn1.Write([]byte("user1\n")); err != nil {
			t.Error("write", err)
		}

		msg2, _ := reader.ReadString('\n')
		assert.Equal(t, "* The room contains:\n", msg2)

		log.Println("connecting user2")
		conn2, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fail()
		}

		reader2 := bufio.NewReader(conn2)

		msg3, _ := reader2.ReadString('\n')
		assert.Equal(t, "Welcome to budgetchat! What shall I call you?\n", msg3)

		if _, err := conn2.Write([]byte("user2\n")); err != nil {
			t.Error("write", err)
		}

		msg4, _ := reader2.ReadString('\n')
		assert.Equal(t, "* The room contains: user1\n", msg4)

		if _, err := conn2.Write([]byte("Hello user1\n")); err != nil {
			t.Error("write", err)
		}

		msg5, _ := reader.ReadString('\n')
		assert.Equal(t, "* user2 has entered the room\n", msg5)

		msg6, _ := reader.ReadString('\n')
		assert.Equal(t, "[user2] Hello user1\n", msg6)

		conn3, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			t.Fail()
		}

		reader3 := bufio.NewReader(conn3)
		msg7, _ := reader3.ReadString('\n')

		assert.Equal(t, "Welcome to budgetchat! What shall I call you?\n", msg7)

		if _, err := conn3.Write([]byte("userreuw8723rbdjsay8iyiohwnd1\n")); err != nil {
			t.Error("write", err)
		}

		time.Sleep(1 * time.Second)
	})
}
