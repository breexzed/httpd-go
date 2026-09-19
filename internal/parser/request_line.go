package parser

import (
	"fmt"
	"strings"
	"unicode"
)

func ParseRequestLine(line string) (method, target, version string, err error) {

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
