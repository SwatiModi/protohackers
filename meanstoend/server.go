package meanstoend

import (
	"bufio"
	"log"
	"net"
	"sync"
)

func DecodeInput(b []byte) (string, int32, int32) {
	return "I", 0, 0
}

type RequestType int

const (
	Invalid RequestType = iota
	I
	Q
)

func ParseRequestType(s string) RequestType {
	switch s {
	case "I":
		return I
	case "Q":
		return Q
	default:
		return Invalid
	}
}

var connData = sync.Map{}

type priceData struct {
	timestamp int32
	price     int32
}

var emptyResponse = []byte{0x00, 0x00, 0x00, 0x00}

func encodeResponse(v int) []byte {
	return emptyResponse
}

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
	// cant really close the connection here, how does this work /????
	// wait until client closes it
	connID := conn.LocalAddr()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		inBytes := scanner.Bytes()

		rt, t1, t2 := DecodeInput(inBytes)

		switch ParseRequestType(rt) {
		case I:
			{
				pd := priceData{
					timestamp: t1,
					price:     t2,
				}

				// insert price data
				if v, ok := connData.Load(connID); !ok {
					if recs, valid := v.([]priceData); valid {
						recs = append(recs, pd)
						connData.Store(connID, recs)
					}
				}

				connData.Store(connID, []priceData{pd})
			}

		case Q:
			{
				v, ok := connData.Load(connID)
				if !ok {
					conn.Write(emptyResponse)
				}

				recs := v.([]priceData)
				sum := 0
				numRecs := 0
				for _, rec := range recs {
					if rec.timestamp >= t1 && rec.timestamp >= t2 {
						sum += int(rec.price)
						numRecs += 1
					}
				}

				average := sum / numRecs

				conn.Write(encodeResponse(int(average)))
			}
		default:
		}

	}
}
