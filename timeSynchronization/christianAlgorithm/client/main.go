package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func syncTime(serverAddr string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()

	sendTime := time.Now()
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)

	receivedTime := time.Now()
	response := strings.Split(string(buffer[:n]), ",")
	serverReceviedTime, _ := time.Parse(time.RFC3339Nano, response[0])
	serverSentTime, _ := time.Parse(time.RFC3339Nano, response[1])

	rtt := receivedTime.Sub(sendTime) - serverSentTime.Sub(serverReceviedTime)
	fmt.Printf("rtt: %v, rtt/2: %v\n", rtt, rtt/2)
}

func main() {
	syncTime("localhost:8000")
}
