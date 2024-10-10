package smoketest

import (
	"bufio"
	"io"
	"log"
	"net"
)

func StartServer() {
	// support 5 simultaneous requests
	concPool := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		concPool <- true
	}

	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("listening on port 8000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept connection", err)
		}

		<-concPool
		go handleRequest(conn)
		concPool <- true
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println("read message", err)
		}

		if _, err := conn.Write([]byte(msg)); err != nil {
			log.Println("conn write", err)
			return
		}
	}
}
