package meanstoend

import (
	"io"
	"log"
	"net"
	"sync"
)

func DecodeInput(b []byte) (int, uint32, uint32) {
	typeByte := b[0]
	tsBytes := b[1:5]
	priceBytes := b[5:]

	rt := ParseRequestType(typeByte)
	if rt == Invalid {
		return 0, 0, 0
	}

	log.Println("rt ", rt)

	timestamp := convertToDecimal(tsBytes)
	price := convertToDecimal(priceBytes)

	return rt, timestamp, price
}

func convertToDecimal(b []byte) uint32 {
	// handle two's complement
	if b[0]&0x80 != 0 { // if first bit is set, 0x80 is a special number in binary (10000000).
		for i := range b {
			b[i] = ^b[i] // flip bits
		}
		val := (bigEndian(b))
		return -val
	}

	return bigEndian(b)
}

func bigEndian(b []byte) uint32 {

	// mathematically we find the big endian by multiplying by powers of 256,
	// but to do this more efficiently we do byte shifting, 24, 16, 8, 0
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

const (
	Invalid = 0
	I       = 73
	Q       = 81
)

func ParseRequestType(s byte) int {
	switch s {
	case 73:
		return I // ASCII FOR I : 73
	case 81:
		return Q // ASCII FOR Q : 81
	default:
		return Invalid
	}
}

var connData = sync.Map{}

type priceData struct {
	timestamp uint32
	price     uint32
}

var emptyResponse = []byte{0x00, 0x00, 0x00, 0x00}

func encodeResponse(v int32) []byte {
	if v == 0 {
		return emptyResponse
	}

	res := make([]byte, 4)
	res[0] = byte(v >> 24)
	res[1] = byte(v >> 16)
	res[2] = byte(v >> 8)
	res[3] = byte(v)

	return res
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
	log.Println("connID", connID)

	for {
		inBytes := make([]byte, 9)
		if _, err := io.ReadFull(conn, inBytes); err != nil {
			log.Println(err, " errrrrrrrrr ")
			if err == io.EOF {
				_, ok := connData.LoadAndDelete(connID)
				log.Println("CONN closed, clearning data", ok)
				break
			}

			if err == io.ErrUnexpectedEOF {
				log.Println("undefined behavior")
				conn.Close()
				break
			}

			log.Println("failed to read", err)
			continue
		}

		rt, t1, t2 := DecodeInput(inBytes)

		switch rt {
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

				sum := 0
				numRecs := 0

				if v != nil {
					recs := v.([]priceData)

					for _, rec := range recs {
						if rec.timestamp >= t1 && rec.timestamp <= t2 {
							sum += int(rec.price)
							numRecs += 1
						}
					}
				}

				var average int32
				if numRecs > 0 {
					average = int32(sum / numRecs)
				}

				conn.Write(encodeResponse(average))
				log.Println("return response", average)
			}
		default:
		}

	}
}
