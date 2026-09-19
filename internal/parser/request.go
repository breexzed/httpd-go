package parser

import (
	"bufio"
	"fmt"
	"strings"
	"time"
)

type Request struct {
	Method   string
	Target   string
	Addr     string
	Version  string
	Headers  map[string][]string
	Duration time.Duration
}

func ReadRequest(reader *bufio.Reader, remoteAddr string) (*Request, error) {
	start := time.Now()

	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(line, "\r\n") {
		return nil, fmt.Errorf("malformed line, missing CRLF: %q", line)
	}
	line = line[:len(line)-2]

	method, target, version, err := ParseRequestLine(line)
	if err != nil {
		return nil, err
	}

	headers, err := ParseHeaders(reader)
	if err != nil {
		return nil, err
	}

	return &Request{
		Method:   method,
		Target:   target,
		Addr:     remoteAddr,
		Version:  version,
		Headers:  headers,
		Duration: time.Since(start),
	}, nil
}
