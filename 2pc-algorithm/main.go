package main

import (
	"fmt"
	"net"
	"time"
)

const (
	PREPARE  = "PREPARE"
	COMMIT   = "COMMIT"
	ROLLBACK = "ROLLBACK"
	ACK      = "ACK"
)

func coordinate(participants []string) {
	// Phase 1: Prepare
	allAck := true
	for _, p := range participants {
		conn, err := net.Dial("tcp", p)
		if err != nil {
			fmt.Printf("Failed to connect to %s: %v\n", p, err)
			allAck = false
			break
		}
		defer conn.Close()

		conn.Write([]byte(PREPARE))
		buffer := make([]byte, 1024)
		n, _ := conn.Read(buffer)
		if string(buffer[:n]) != ACK {
			allAck = false
			break
		}
	}

	// Phase 2: Commit/Rollback
	var msg string
	if allAck {
		msg = COMMIT
		fmt.Println("All participants prepared, sending commit...")
	} else {
		msg = ROLLBACK
		fmt.Println("Some participants not prepared, sending rollback...")
	}

	for _, p := range participants {
		conn, err := net.Dial("tcp", p)
		if err != nil {
			fmt.Printf("Failed to connect to %s: %v\n", p, err)
			continue
		}
		defer conn.Close()

		conn.Write([]byte(msg))
	}
}

func participate(id int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8000+id))
	if err != nil {
		fmt.Printf("Failed to start participant %d: %v\n", id, err)
		return
	}
	defer listener.Close()

	fmt.Printf("Participant %d started\n", id)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Participant %d failed to accept connection: %v\n", id, err)
			continue
		}

		go func(id int, conn net.Conn) {
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Printf("Participant %d failed to read message: %v\n", id, err)
				return
			}

			msg := string(buffer[:n])
			if msg == PREPARE {
				fmt.Printf("Participant %d received prepare\n", id)
				conn.Write([]byte(ACK))
			} else if msg == COMMIT {
				fmt.Printf("Participant %d transaction commit...\n", id)
			} else if msg == ROLLBACK {
				fmt.Printf("Participant %d transaction rollback...\n", id)
			}
		}(id, conn)
	}
}

func main() {
	go participate(0)
	go participate(1)
	time.Sleep(1 * time.Second)

	participants := []string{"localhost:8000", "localhost:8001"}
	coordinate(participants)
}
