

package main

#comment added

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
)

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	maxConn := getenvInt("MAX_CONN", 5)
	limit := make(chan struct{}, maxConn)

	listener, err := net.Listen("tcp", ":"+port)
	var wg sync.WaitGroup
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		fmt.Println("\nShutting down...")
		listener.Close()
	}()

	fmt.Println("Server running on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		limit <- struct{}{}
		wg.Add(1)
		go handleConnection(conn, limit, &wg)
	}

	wg.Wait()
	fmt.Println("Shutdown complete")
}
