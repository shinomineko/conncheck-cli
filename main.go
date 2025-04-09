package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	connTimeout time.Duration
	httpsVerify bool
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "environment variables:\n")
		fmt.Fprintf(os.Stderr, "  CONN_TYPE     connection type: tcp, udp, http, https (default: tcp)\n")
		fmt.Fprintf(os.Stderr, "  DEST_HOST     destination hostname or IP (default: localhost)\n")
		fmt.Fprintf(os.Stderr, "  DEST_PORT     destination port number (default: 80)\n")
		fmt.Fprintf(os.Stderr, "  CONN_TIMEOUT  connection timeout in seconds (default: 5)\n")
		fmt.Fprintf(os.Stderr, "  HTTPS_VERIFY  verify HTTPS certificates (default: true)\n")
		fmt.Fprintf(os.Stderr, "exit codes:\n")
		fmt.Fprintf(os.Stderr, "  0  connection successful\n")
		fmt.Fprintf(os.Stderr, "  1  connection failed\n")
		flag.PrintDefaults()
	}

	flag.Parse()

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
		connTimeout = 5 * time.Second
	} else {
		timeoutInt, err := strconv.Atoi(timeoutStr)
		if err != nil {
			fmt.Printf("invalid timeout value: %s\n", timeoutStr)
			os.Exit(1)
		}
		connTimeout = time.Duration(timeoutInt) * time.Second
	}

	httpsVerifyStr, exists := os.LookupEnv("HTTPS_VERIFY")
	if !exists {
		httpsVerify = true
	} else {
		var err error
		httpsVerify, err = strconv.ParseBool(httpsVerifyStr)
		if err != nil {
			fmt.Printf("invalid value: %s\n", httpsVerifyStr)
			os.Exit(1)
		}
	}

	addr := fmt.Sprintf("%s:%s", destHost, destPort)
	fmt.Printf("testing %s connection to %s\n", connType, addr)

	switch connType {
	case "tcp":
		testTCP(addr)
	case "udp":
		testUDP(addr)
	case "http":
		testHTTP(addr, false)
	case "https":
		testHTTP(addr, true)
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

func testHTTP(addr string, isHTTPS bool) {
	var url string
	if isHTTPS {
		url = "https://" + addr
	} else {
		url = "http://" + addr
	}

	client := &http.Client{
		Timeout: connTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !httpsVerify,
			},
		},
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		fmt.Printf("failed to create request: %s\n", err)
		os.Exit(1)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("connection failed: %s\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	fmt.Printf("connection successful: %s %s\n", resp.Proto, resp.Status)
}
