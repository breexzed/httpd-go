package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
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

		//wrapped handleConnection in a goroutine this specific way
		//so I can handle errors for any connection that fails without it affecting the others
		//initially was go handleConnection(conn)
		go func() {
			defer conn.Close()

			reader := bufio.NewReader(conn)

			for {
				err := handleConnection(conn, reader)
				if err != nil {
					if !errors.Is(err, io.EOF) {
						fmt.Println("Error handling connection:", err)
					}
					break
				}
			}
		}()
	}
}

func handleConnection(conn net.Conn, reader *bufio.Reader) error {
	//used a struct cause I like reading elegant code, not really necessary
	//probably not the right way to do this but going to let this one slide
	type Request struct {
		method   string
		target   string
		addr     string
		version  string
		duration int
	}

	//ensure connections closes after it's being handles
	//finally I get to use this stuff in the real world. very useful.

	start := time.Now()
	ipAddr := conn.RemoteAddr().String()

	//without bufio I wouldn't have been able to read off conn(which is a byte)
	//Off is the word here cause this implies I'm reading its content on the fly
	//which is exactly what I went on t do with it

	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	if !strings.HasSuffix(line, "\r\n") { //impressive strings.Method from the golang engineers

		return fmt.Errorf("malformed line, missing CRLF: %q", line)
	}
	line = line[:len(line)-2]

	method, target, version, err := parseRequestLine(line)
	if err != nil {
		return err
	}

	r := Request{
		method:   method,
		target:   target,
		addr:     ipAddr,
		version:  version,
		duration: int(time.Since(start).Milliseconds()),
	}
	fmt.Printf("parsed request:\n"+
		"method: %s\n"+
		"target: %s\n"+
		"addr: %s\n"+
		"version: %s\n"+
		"duration: %d\n",
		r.method, r.target, r.addr, r.version, r.duration)

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 2\r\n\r\n" +
		"OK"

	//It was not clear to me how this even is possible but basically what it does is
	//it asserts that the only value that matters here is error, so we deal with it here,
	//as well as in this function usage back in main wrapped in a goroutine
	_, err = conn.Write([]byte(response))
	if err != nil {
		return err
	}

	return nil
}

func parseRequestLine(line string) (method, target, version string, err error) {

	//learnt a ton about the string package writing this function
	//not so much but enough that I can now think about how to use it
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
