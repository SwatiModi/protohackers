package mobinthemiddle

import (
	"bufio"
	"io"
	"log"
	"net"
	"regexp"
	"sync"
)

func RewriteAddresses(data string) string {
	expr := regexp.MustCompile(`(^|\s)7[a-zA-Z0-9]{26,35}($|\s)`)
	tonysAddress := "7YWHMfk9JZe0LM0g1ZauHuiSxhI"

	return expr.ReplaceAllStringFunc(data, func(match string) string {
		// Preserve leading/trailing space
		if match[0] == ' ' {
			return " " + tonysAddress
		}
		if match[len(match)-1] == ' ' {
			return tonysAddress + " "
		}
		return tonysAddress
	})
}

func initUpstreamConnection() (net.Conn, error) {
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
	upstream, err := initUpstreamConnection()
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

func broker(src, dest net.Conn, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	reader := bufio.NewReader(src)
	writer := bufio.NewWriter(dest)

	for {
		// Read data
		data, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from source: %v", err)
			}
			break
		}

		// Rewrite Boguscoin addresses
		data = RewriteAddresses(data)

		// Write data
		_, err = writer.WriteString(data)
		if err != nil {
			log.Printf("Error writing to destination: %v", err)
			break
		}
		writer.Flush()
	}
}
