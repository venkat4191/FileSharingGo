package main

import (
	"net"
	"strconv"
)

func writeHTTP(
    conn net.Conn,
    status int,
    body string,
    keepAlive bool,
    contentType string,
) {
	writeHTTPBinary(conn, status, []byte(body), keepAlive, contentType)
}

func writeHTTPBinary(
    conn net.Conn,
    status int,
    body []byte,
    keepAlive bool,
    contentType string,
) {

	connectionHeader := "close"
	if keepAlive {
		connectionHeader = "keep-alive"
	}

	statusText := "OK"
	switch status {
	case 200:
		statusText = "OK"
	case 201:
		statusText = "Created"
	case 204:
		statusText = "No Content"
	case 400:
		statusText = "Bad Request"
	case 404:
		statusText = "Not Found"
	case 405:
		statusText = "Method Not Allowed"
	}

	header := "HTTP/1.1 " + strconv.Itoa(status) + " " + statusText + "\r\n" +
		"Content-Type: " + contentType + "\r\n" +
		"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
		"Connection: " + connectionHeader + "\r\n" +
		"\r\n"

	conn.Write([]byte(header))
	conn.Write(body)
}