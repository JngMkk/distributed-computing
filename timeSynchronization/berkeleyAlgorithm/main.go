package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	portMaster    = ":8000"
	portFollower1 = ":8001"
	portFollower2 = ":8002"
	serverIP      = "127.0.0.1"
)

func getTimeNow() string {
	return time.Now().Format(time.RFC3339Nano)
}

func startServer(port string, handleConnection func(net.Conn)) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Failed to start server on %s: %v\n", port, err)
		return
	}
	defer listener.Close()
	fmt.Printf("Server started on %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection on %s: %v\n", port, err)
			continue
		}
		go handleConnection(conn)
	}
}

func masterServer() {
	followers := []string{serverIP + portFollower1, serverIP + portFollower2}
	timeDiffs := make([]time.Duration, len(followers))

	// 팔로워 노드들 시간 수집
	for i, follower := range followers {
		conn, err := net.Dial("tcp", follower)
		if err != nil {
			fmt.Printf("Failed to connect to follower %s: %v\n", follower, err)
			continue
		}

		followerTimeStr, _ := bufio.NewReader(conn).ReadString('\n')
		followerTime, _ := time.Parse(time.RFC3339Nano, strings.TrimSpace(followerTimeStr))
		masterTime := time.Now()

		timeDiffs[i] = masterTime.Sub(followerTime)
		conn.Close()
	}

	// 시간 차 평균 계산
	var total time.Duration
	for _, diff := range timeDiffs {
		total += diff
	}
	averageDiff := total / time.Duration(len(timeDiffs))

	// 조정값 전송
	for _, follower := range followers {
		conn, err := net.Dial("tcp", follower)
		if err != nil {
			fmt.Printf("Failed to connect to adjust time for follower %s: %v\n", follower, err)
			continue
		}
		fmt.Fprintf(conn, "%v\n", averageDiff)
		conn.Close()
	}
}

func followerServer(port string) {
	handleConnection := func(conn net.Conn) {
		defer conn.Close()

		// 현재 시간 전송
		fmt.Fprintln(conn, getTimeNow())

		adjustmentStr, _ := bufio.NewReader(conn).ReadString('\n')
		adjustment, _ := time.ParseDuration(strings.TrimSpace(adjustmentStr))
		adjustedTime := time.Now().Add(adjustment)
		fmt.Printf("Adjusted time on %s: %v\n", port, adjustedTime)
	}

	startServer(port, handleConnection)
}

func main() {
	go followerServer(portFollower1)
	go followerServer(portFollower2)
	time.Sleep(1 * time.Second)
	masterServer()
}
