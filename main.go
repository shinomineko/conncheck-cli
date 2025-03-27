package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var connTimeout time.Duration

func main() {
	connType, exists := os.LookupEnv("CONN_TYPE")
	if !exists {
		connType = "tcp"
	}

	destHost, exists := os.LookupEnv("DEST_HOST")
	if !exists {
		destHost = "localhost"
	}

	destPort, exists := os.LookupEnv("DEST_PORT")
	if !exists {
		destPort = "80"
	}

	timeoutStr, exists := os.LookupEnv("CONN_TIMEOUT")
	if !exists {
		connTimeout = 10 * time.Second
	} else {
		timeoutInt, err := strconv.Atoi(timeoutStr)
		if err != nil {
			fmt.Printf("invalid timeout value: %s\n", timeoutStr)
			os.Exit(1)
		}
		connTimeout = time.Duration(timeoutInt) * time.Second
	}

	addr := fmt.Sprintf("%s:%s", destHost, destPort)
	fmt.Printf("testing %s connection to %s\n", connType, addr)

	switch connType {
	case "tcp":
		testTCP(addr)
	case "udp":
		testUDP(addr)
	default:
		fmt.Printf("unsupported connection type: %s\n", connType)
		os.Exit(1)
	}
}

func testTCP(addr string) {
	conn, err := net.DialTimeout("tcp", addr, connTimeout)
	if err != nil {
		fmt.Printf("connection failed: %s\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("connection successful\n")
}

func testUDP(addr string) {
	conn, err := net.DialTimeout("udp", addr, connTimeout)
	if err != nil {
		fmt.Printf("connection failed: %s\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// udp doesn't establish connection, so just test if we can create socket
	fmt.Printf("socket creation successful\n")
}
