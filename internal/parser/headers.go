package parser

import (
	"bufio"
	"fmt"
	"strings"
)

func ParseHeaders(reader *bufio.Reader) (map[string][]string, error) {
	m := make(map[string][]string)

	const maxHeaders = 100
	headerCount := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if line == "\r\n" {
			break
		}

		headerCount++
		if headerCount > maxHeaders {
			return nil, fmt.Errorf("header limit reached: maximum allowed: %d", maxHeaders)
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed line: %q", line)
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.TrimSpace(parts[1])

		m[key] = append(m[key], val)
	}
	return m, nil
}
