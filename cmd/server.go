package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/breexzed/httpd-go/internal/parser"
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
			defer conn.Close() //ensure connections closes after it's being handles. finally I get to use this stuff in the real world. very useful.
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
	// Parse request line & headers via internal package
	req, err := parser.ReadRequest(reader, conn.RemoteAddr().String())
	if err != nil {
		return err
	}

	// 2. Access exported struct fields
	fmt.Printf("parsed request:\n"+
		"method: %s\n"+
		"target: %s\n"+
		"addr: %s\n"+
		"version: %s\n"+
		"headers: %v\n"+
		"duration: %v\n",
		req.Method, req.Target, req.Addr, req.Version, req.Headers, req.Duration)

	// 3. Send response
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 2\r\n\r\n" +
		"OK"

	_, err = conn.Write([]byte(response))
	return err
}
