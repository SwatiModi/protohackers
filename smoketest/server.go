package smoketest

import (
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
		go func(conn net.Conn) {
			defer func() {
				concPool <- true
			}()
			handleRequest(conn)
		}(conn)

	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	if _, err := io.Copy(conn, conn); err != nil {
		return
	}
}
