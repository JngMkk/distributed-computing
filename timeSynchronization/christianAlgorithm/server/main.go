package main

import (
	"fmt"
	"net"
	"time"
)

func start(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Time server started on %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte(fmt.Sprintf("%s,%s", getTimeNow(), getTimeNow())))
}

func getTimeNow() string {
	return time.Now().Format(time.RFC3339Nano)
}

func main() {
	start("8000")
}
