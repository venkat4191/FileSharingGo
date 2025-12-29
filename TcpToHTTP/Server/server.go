package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func handleConnection(conn net.Conn, limit chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	defer func() { <-limit }()

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		req, err := readHTTP(conn)
		if err != nil {
			return
		}

		fmt.Printf("[%s] %s %s\n", conn.RemoteAddr(), req.Method, req.Path)

		status, body := routeRequest(req)
		keepAlive := strings.ToLower(req.headers["connection"]) != "close"
		
		contentType := "text/plain"
		isBinaryDownload := strings.HasPrefix(req.Path, "/files/") && strings.Contains(req.Path, "?download=true")
		
		if isBinaryDownload {
			filename := strings.TrimPrefix(req.Path, "/files/")
			if idx := strings.Index(filename, "?"); idx != -1 {
				filename = filename[:idx]
			}
			filename = urlDecode(filename)
			contentType = getContentType(filename)
			
			data, err := downloadFile(filename)
			if err == nil {
				writeHTTPBinary(conn, 200, data, keepAlive, contentType)
			} else {
				writeHTTP(conn, 404, "File not found", keepAlive, "text/plain")
			}
		} else {
			if req.Path == "/" || req.Path == "/index.html" {
				contentType = "text/html"
			} else if req.Path == "/styles.css" {
				contentType = "text/css"
			} else if strings.HasPrefix(req.Path, "/files") || req.Path == "/upload" {
				contentType = "application/json"
			}

			writeHTTP(
				conn,
				status,
				body,
				keepAlive,
				contentType,
			)
		}
		

		if !keepAlive {
			return
		}
	}
}