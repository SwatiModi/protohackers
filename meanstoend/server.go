package meanstoend

import (
	"encoding/binary"
	"io"
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
	addr := conn.RemoteAddr()

	defer func() {
		conn.Close()
		log.Printf("closed connection (%v)", addr)
	}()

	connData := make(map[int32]int32)
	buf := make([]byte, 9)

	for {
		if _, err := io.ReadFull(conn, buf); err == io.EOF {
			break
		} else if err != nil {
			log.Printf("%v (%v)", err, addr)
			break
		}

		t1 := int32(binary.BigEndian.Uint32(buf[1:5]))
		t2 := int32(binary.BigEndian.Uint32(buf[5:]))

		switch buf[0] {
		case 'I':
			connData[t1] = t2
			log.Printf("insert: %v %v (%v)", t1, t2, addr)

		case 'Q':
			{
				sum := 0
				n := 0

				for ts, price := range connData {
					if ts >= t1 && ts <= t2 {
						sum += int(price)
						n += 1
					}
				}

				var average int32
				if n > 0 {
					average = int32(sum / n)
				}

				out := make([]byte, 4)
				binary.BigEndian.PutUint32(out, uint32(average))

				if _, err := conn.Write(out); err != nil {
					log.Printf("%v (%v)", err, addr)
				} else {
					log.Printf("query: %v %v ⇒ %v (%v)", t1, t2, out, addr)
				}
			}
		default:
			log.Printf("received invalid input  { %v } ", buf[0])
		}
	}
}
