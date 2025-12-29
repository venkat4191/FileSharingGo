package main

import (
    "net"
    "os"
    "strings"
)

func serveStatic(conn net.Conn, path string) bool {
    if !strings.HasPrefix(path, "/static/") {
        return false
    }

    filePath := "." + path
    data, err := os.ReadFile(filePath)
    if err != nil {
    writeHTTP(conn, 404, "File not found", true, "text/plain")
    return true
    }

    writeHTTP(conn, 200, string(data), true, "text/plain")
    return true
}
