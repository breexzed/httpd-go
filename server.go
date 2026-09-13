package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"unicode"
)

func main() {
	ln, err := net.Listen("tcp", ":2000")
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on :2000")
	for {
		conn, err := ln.Accept()
		if err != nil {

			panic(err)
		}
		fmt.Println("Accepted new connection")

		go func() {
			err := handleConnection(conn)
			if err != nil {
				fmt.Println("connection error:", err)
			}
		}()
	}
}

func handleConnection(conn net.Conn) error {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	if !strings.HasSuffix(line, "\r\n") {

		return fmt.Errorf("malformed line, missing CRLF: %q", line)
	}
	line = line[:len(line)-2]

	method, target, version, err := parseRequestLine(line)
	if err != nil {
		return err
	}
	fmt.Printf("parsed request: method=%s target=%s version=%s\n", method, target, version)

	response := "HTTP/1.1 200 OK\r\nContent-Length: 2\r\n\r\nOK"
	_, err = conn.Write([]byte(response))
	if err != nil {
		return err
	}

	return nil
}

func parseRequestLine(line string) (method, target, version string, err error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("malformed line: %q", line)
	}
	method, target, version = parts[0], parts[1], parts[2]

	if !strings.HasPrefix(version, "HTTP/") {
		return "", "", "", fmt.Errorf("malformed line: %q", line)
	}
	versionNumber := strings.TrimPrefix(version, "HTTP/")
	if len(versionNumber) != 3 || versionNumber[1] != '.' ||
		!unicode.IsDigit(rune(versionNumber[0])) || !unicode.IsDigit(rune(versionNumber[2])) {
		return "", "", "", fmt.Errorf("malformed version: %q", version)
	}

	switch method {
	case "GET":
		// supported
	default:
		return "", "", "", fmt.Errorf("unsupported method: %q", method)
	}

	return method, target, version, nil
}
