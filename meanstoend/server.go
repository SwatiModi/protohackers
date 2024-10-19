package meanstoend

import (
	"log"
	"net"
)

func StartServer() {

	concPool := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		concPool <- true
	}

	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Println("start server", err)
	}

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
	log.Println("connID", conn.LocalAddr().String(), "local", conn.RemoteAddr().String())
	conn.Write([]byte("connID" + conn.LocalAddr().String() + "local" + conn.RemoteAddr().String() + "\n"))
}

func hexToDecimal() {

}
