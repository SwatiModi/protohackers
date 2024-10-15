package primetime

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
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

	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Println("read message", err)
		return
	}

	n, err := parseRequest(msg)
	if err != nil {
		if writeResp([]byte("hello world\n"), conn) != nil {
			log.Println("write response", err)
		}
		return
	}

	var respBytes []byte
	if isPrime(n) {
		respBytes = []byte("{\"method\":\"isPrime\", \"prime\":true}\n")
	} else {
		respBytes = []byte("{\"method\":\"isPrime\", \"prime\":false}\n")
	}

	if err := writeResp(respBytes, conn); err != nil {
		log.Println("write response", err)
	}
}

func parseRequest(msg string) (int, error) {
	var n int = -1
	var methodExists bool
	msg = strings.TrimSpace(msg)

	if len(msg) < 2 || msg[0] != '{' || msg[len(msg)-1] != '}' {
		return 0, errors.New("invalid brackets")
	}

	kvs := strings.Split(msg[1:len(msg)-1], ",")
	for _, kv := range kvs {
		kv = strings.TrimSpace(kv)
		keyValue := strings.Split(kv, ":")
		if len(keyValue) != 2 {
			return 0, errors.New("malformed key-value pair")
		}
		key, value := keyValue[0], keyValue[1]

		if key == `"method"` && value == `"isPrime"` {
			methodExists = true
			continue
		} else if key == `"number"` {
			number, err := strconv.Atoi(value)
			if err != nil {
				return 0, err
			}
			n = number
		}
	}
	if methodExists && n >= 0 {
		return n, nil
	}

	return 0, errors.New("unknown error")
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

func writeResp(bytes []byte, conn net.Conn) error {
	if _, err := conn.Write(bytes); err != nil {
		log.Println("conn write", err)
		return err
	}

	return nil
}
