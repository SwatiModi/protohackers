package mobinthemiddle

import (
	"bufio"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
)

var boguscoin = regexp.MustCompile(`^7[a-zA-Z0-9]{25,34}$`)

func getUpstreamConnection() (net.Conn, error) {
	addr := "chat.protohackers.com:16963"
	return net.Dial("tcp", addr)
}

func StartServer() {
	// Limit to 10 concurrent clients
	concPool := make(chan struct{}, 10)
	for i := 0; i < 10; i++ {
		concPool <- struct{}{}
	}

	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	log.Println("Listening on :8000")

	for {
		downstream, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		<-concPool
		go func() {
			defer func() { concPool <- struct{}{} }()
			handleClient(downstream)
		}()
	}
}

func handleClient(downstream net.Conn) {
	defer downstream.Close()

	// Create a new upstream connection
	upstream, err := getUpstreamConnection()
	if err != nil {
		log.Printf("Error creating upstream connection: %v", err)
		return
	}
	defer upstream.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	// Start proxying data
	go broker(downstream, upstream, &wg) // Downstream -> Upstream
	go broker(upstream, downstream, &wg) // Upstream -> Downstream

	wg.Wait()
}

func broker(src, dst net.Conn, wg *sync.WaitGroup) {
	defer func() {
		src.Close()
		dst.Close()
		wg.Done()
	}()

	for reader := bufio.NewReader(src); ; {
		data, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		tokens := make([]string, 0, 8)
		for _, raw := range strings.Split(data[:len(data)-1], " ") {
			t := boguscoin.ReplaceAllString(raw, "7YWHMfk9JZe0LM0g1ZauHuiSxhI")
			tokens = append(tokens, t)
		}

		out := strings.Join(tokens, " ") + "\n"
		_, err = dst.Write([]byte(out))
		if err != nil {
			log.Printf("Error writing to destination: %v at %v", err, dst.RemoteAddr().String())
		}
	}
}
