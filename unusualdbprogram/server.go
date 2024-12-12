package unusualdbprogram

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
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

	log.Printf("Listening on port %v", 8000)
	for {
		buff := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buff)
		if err != nil {
			log.Printf("Error reading from connection: %v", err)
			continue
		}

		handleRequest(pc, addr, buff[:n])
	}
}

func handleRequest(pc net.PacketConn, addr net.Addr, data []byte) {
	reqType, key, value := getKeyValue(string(data))

	switch reqType {
	case InsertRequest:
		if key == "version" {
			log.Printf("Cannot insert version")
			return
		}
		store.Store(key, value)
		log.Printf("Inserted: %s=%s", key, value)
	case RetreiveRequest:
		val, ok := store.Load(key)
		if !ok {
			// no resp for key not found
			log.Printf("Key not found: %s", key)
		} else {
			resp := fmt.Sprintf("%s=%s", key, val)
			pc.WriteTo([]byte(resp), addr)
			log.Printf("Retrieved: %s", resp)
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
