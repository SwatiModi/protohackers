package primetime

import (
	"bufio"
	"encoding/json"
	"log"
	"math"
	"net"
)

var primeMethod = "isPrime"

type request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

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
			continue
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

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		inbytes := scanner.Bytes()

		var req = request{}
		var respBytes []byte

		// malformed request case
		if err := json.Unmarshal(inbytes, &req); err != nil || req.Method != "isPrime" || req.Number == nil {
			respBytes = []byte("malformed request\n")
		} else {
			var resp = response{
				Method: primeMethod,
			}
			n := *req.Number

			if n == math.Trunc(n) && isPrime(int(n)) {
				resp.Prime = true
			} else {
				resp.Prime = false
			}

			if bytes, err := json.Marshal(resp); err != nil {
				log.Println("failed to unmarshal")
			} else {
				respBytes = append(bytes, 10)
			}
		}

		if _, err := conn.Write(respBytes); err != nil {
			log.Println("write resp", err)
		}
	}
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
