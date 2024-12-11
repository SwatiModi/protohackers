package budgetchat

import (
	"log"
	"net"
	"regexp"
	"sync"
)

func StartServer() {
	concPool := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		concPool <- true
	}

	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Println("start server", err)
	}

	log.Printf("listening on port %v", 8000)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept request", err)
		}

		<-concPool
		go func(conn net.Conn) {
			defer func() {
				concPool <- true
			}()
			handleRequest(conn)
		}(conn)
	}
}

func handleRequest(conn net.Conn) {
	log.Printf("new connection from %v", conn.RemoteAddr())
	addr := conn.RemoteAddr()
	var username string

	defer func() {
		users.Delete(username)
		conn.Close()
		log.Printf("closed connection (%v)", addr)
	}()

	// ask for username on connection
	conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))

	// receive username, validate it, if alright, store user and read chat msgs from users
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("read username", err)
		return
	}

	if n < 1 || n > 16 {
		conn.Write([]byte("Username must be between 1 and 16 characters.\n"))
		return
	}

	if expr, err := regexp.Compile("^[a-zA-Z0-9_]*$"); err != nil {
		log.Println("compile regex", err)
		return
	} else if expr != nil && !expr.Match(buf[:n]) {
		conn.Write([]byte("Username must be alphanumeric.\n"))
		return
	}

	username = string(buf[:n])
	_, ok := users.Load(username)
	if ok {
		conn.Write([]byte("Username already taken."))
		return
	}

	uStruct := user{username, make(chan message), conn}
	// name is valid so welcome msg and store user
	users.Store(username, uStruct)
	go receivedChatUpdates(username, uStruct)

	broadcastMsg(username, message{username, "* " + username + " has entered the room\n"})

	introMsg := "* The room contains:"
	users.Range(func(key, value interface{}) bool {
		if key.(string) != username {
			introMsg += " " + key.(string)
		}
		return true
	})

	conn.Write([]byte(introMsg + "\n"))

	// check msgs from user, publish to all users
	for {
		var buf = make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("read msg", err)
			break
		}

		if n < 1 || n > 1000 {
			log.Println("msg too long ", n)
			continue
		}

		log.Println("broadcasting msg from ", username)
		uMsg := "[" + username + "] " + string(buf[:n]) + "\n"
		broadcastMsg(username, message{username, uMsg})
	}
}

type message struct {
	username string
	msg      string
}

var users sync.Map // username -> user

type user struct {
	name    string
	msgChan chan message
	conn    net.Conn
}

func broadcastMsg(username string, msg message) {
	users.Range(func(key, value interface{}) bool {
		userS := value.(user)

		if userS.conn != nil && userS.name != username {
			log.Printf("sending msg to %v", key)
			userS.msgChan <- msg
		}
		return true
	})
}

func receivedChatUpdates(username string, u user) {
	for {
		m := <-u.msgChan
		if m.username != username {
			u.conn.Write([]byte(m.msg))
		}
	}
}
