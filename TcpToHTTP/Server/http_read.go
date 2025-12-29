package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type HTTPrequest struct {
	Method  string
	Path    string
	Version string
	headers map[string]string
	Body    string
}

func readHTTP(conn net.Conn) (*HTTPrequest, error) {
	fmt.Println("Client is :", conn.RemoteAddr())

	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 1024)

	for {
		n, err := conn.Read(tmp)
		if err != nil {
			return nil, err
		}
		buf = append(buf, tmp[:n]...)
		if strings.Contains(string(buf), "\r\n\r\n") {
			break
		}
	}

	const maxHeaderSize = 8 * 1024
	if len(buf) > maxHeaderSize {
		return nil, fmt.Errorf("headers too large")
	}

	parts := strings.SplitN(string(buf), "\r\n\r\n", 2)
	headersParts := parts[0]
	initialBody := ""
	if len(parts) == 2 {
		initialBody = parts[1]
	}

	lines := strings.Split(headersParts, "\r\n")
	line0 := strings.Split(lines[0], " ")

	Met, Pa, Ver := "", "", ""
	if len(line0) == 3 {
		Met = line0[0]
		Pa = line0[1]
		Ver = line0[2]
	}

	headers := make(map[string]string)
	for _, l := range lines[1:] {
		if l == "" {
			break
		}
		dual := strings.SplitN(l, ":", 2)
		if len(dual) == 2 {
			headers[strings.ToLower(strings.TrimSpace(dual[0]))] = strings.TrimSpace(dual[1])
		}
	}

	length := 0
	if val, ok := headers["content-length"]; ok {
		l, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("invalid Content-Length header")
		}
		length = l
	}

	bodyBytes := []byte(initialBody)
	for len(bodyBytes) < length {
		n, err := conn.Read(tmp)
		if err != nil {
			return nil, fmt.Errorf("error reading body: %v", err)
		}
		if n == 0 {
			break
		}
		bodyBytes = append(bodyBytes, tmp[:n]...)
	}
	
	if len(bodyBytes) > length {
		bodyBytes = bodyBytes[:length]
	}
	
	bodyParts := string(bodyBytes)

	req := &HTTPrequest{
		Method:  Met,
		Path:    Pa,
		Version: Ver,
		headers: headers,
		Body:    bodyParts,
	}
	return req, nil
}
