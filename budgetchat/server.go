package budgetchat

import (
	"bufio"
	"context"
	"errors"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
)

type ClientState int

const (
	ClientConnected = iota // 0
	ClientJoined
	ClientLeft
)

type Message struct {
	name string
	msg  string
}

var clientList sync.Map // username -> user

type Client struct {
	name    string
	msgChan chan Message
	c       net.Conn
	state   ClientState
}

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
	ctx := context.TODO()
	addr := conn.RemoteAddr().String()
	log.Printf("new connection from %v", conn.RemoteAddr())

	// ask for username on connection
	conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))

	client := Client{
		name:    "",
		msgChan: make(chan Message, 10),
		c:       conn,
		state:   ClientConnected,
	}

	scanner := bufio.NewScanner(conn)

	// check msgs from user, publish to all users
	for scanner.Scan() {
		inTxt := scanner.Text()

		// skip longer msgs
		if len(inTxt) > 1000 {
			continue
		}

		switch client.state {
		case ClientConnected:
			nm := strings.TrimSpace(inTxt)
			if err := validateUsername(nm, addr); err != nil {
				log.Printf("-----invalid name %v", err)
				goto DISCONNECT
			}

			log.Printf("user %v has joined", nm)

			client.name = nm
			client.state = ClientJoined

			// name is valid so welcome msg and store user
			clientList.Store(addr, client)
			go receivedChatUpdates(ctx, client)

			broadcastMsg(Message{nm, "* " + nm + " has entered the room\n"})

			introMsg := "* The room contains:"
			clientList.Range(func(key, value interface{}) bool {
				cl := value.(Client)
				if cl.name != client.name {
					introMsg += " " + cl.name
				}
				return true
			})

			conn.Write([]byte(introMsg + "\n"))
		case ClientJoined:
			log.Printf("broadcasting msg from %v : %v", client.name, inTxt)
			uMsg := "[" + client.name + "] " + inTxt + "\n"
			broadcastMsg(Message{client.name, uMsg})
		}
	}

	clientList.Delete(addr)
	client.state = ClientLeft

DISCONNECT:
	conn.Close()
	ctx.Done()
	if len(client.name) > 0 {
		broadcastMsg(Message{client.name, "* " + client.name + " has left the room\n"})
	}
	log.Printf("closed connection (%v) - (%v)", addr, client.name)
}

func broadcastMsg(msg Message) {
	clientList.Range(func(key, value interface{}) bool {
		userS := value.(Client)

		if userS.c != nil && userS.name != msg.name {
			log.Printf("sending msg to %v", userS.name)
			userS.msgChan <- msg
		}
		return true
	})
}

func receivedChatUpdates(ctx context.Context, u Client) {
	for {
		select {
		case <-ctx.Done():
			log.Println("stop receiving updates", u.name)
			return
		case mg := <-u.msgChan:
			if mg.name != u.name {
				u.c.Write([]byte(mg.msg))
			}
		}
	}
}

func validateUsername(name string, addr string) error {
	n := len(name)
	log.Println("validateUsername", n)
	if n < 1 || n > 16 {
		return errors.New("username must be between 1 and 16 characters")
	}

	// strip trailing newline and carrier return chars from username
	expr, _ := regexp.Compile("^[a-zA-Z0-9_]*$")
	if expr != nil && !expr.MatchString(name) {
		return errors.New("username must be alphanumeric")
	}

	_, ok := clientList.Load(addr)
	if ok {
		return errors.New("username already taken")
	}

	return nil
}
