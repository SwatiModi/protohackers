package unusualdbprogram

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var store = sync.Map{}
var RequestType int

const (
	InsertRequest = iota
	RetreiveRequest
)

func StartServer() {
	// add version to kv-store
	store.Store("version", "0.0.1")

	pc, err := net.ListenPacket("udp", ":8000")
	if err != nil {
		log.Fatal("start server", err)
	}

	log.Printf("listening on port %v", 8000)
	buff := make([]byte, 1024)
	timeout := time.Second * 5

	for {
		n, addr, err := pc.ReadFrom(buff)
		if err != nil {
			log.Println("accept conn", err)
		}

		log.Printf("received %v bytes from %v", n, addr.String())

		// set timeout deadline to udp write
		deadline := time.Now().Add(timeout)

		if err := pc.SetWriteDeadline(deadline); err != nil {
			log.Println("set write deadline", err)
		}

		reqType, key, value := getKeyValue(string(buff[:n]))
		if reqType == InsertRequest {
			store.Store(key, value)
		} else {
			val, ok := store.Load(key)
			if !ok {
				pc.WriteTo([]byte("Key not found"), addr)
			} else {
				resp := fmt.Sprintf("%v=%v", key, val)
				pc.WriteTo([]byte(resp), addr)
			}
		}
	}
}

func getKeyValue(in string) (int, string, string) {
	parts := strings.SplitN(in, "=", 2)
	if len(parts) < 2 {
		return RetreiveRequest, in, ""
	}

	return InsertRequest, parts[0], parts[1]
}
