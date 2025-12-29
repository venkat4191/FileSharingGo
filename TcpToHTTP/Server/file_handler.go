package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const sharedDir = "./shared"

func initSharedDir() {
	if _, err := os.Stat(sharedDir); os.IsNotExist(err) {
		os.MkdirAll(sharedDir, 0755)
	}
}

type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func listFiles() ([]FileInfo, error) {
	initSharedDir()
	files, err := os.ReadDir(sharedDir)
	if err != nil {
		return nil, err
	}

	var fileList []FileInfo
	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				continue
			}
			fileList = append(fileList, FileInfo{
				Name: file.Name(),
				Size: info.Size(),
			})
		}
	}
	return fileList, nil
}

func uploadFile(filename string, content string) error {
	initSharedDir()
	
	if content == "" {
		return fmt.Errorf("content is empty")
	}
	
	var data []byte
	var err error
	
	data, err = base64.StdEncoding.DecodeString(content)
	if err != nil {
		return fmt.Errorf("failed to decode base64 content: %v", err)
	}

	filename = filepath.Base(filename)
	if filename == "" || filename == "." || filename == ".." {
		return fmt.Errorf("invalid filename")
	}
	
	filePath := filepath.Join(sharedDir, filename)
	
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	
	return nil
}

func downloadFile(filename string) ([]byte, error) {
	initSharedDir()
	filename = filepath.Base(filename)
	filePath := filepath.Join(sharedDir, filename)
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("file not found")
	}
	return data, nil
}

func deleteFile(filename string) error {
	initSharedDir()
	filename = filepath.Base(filename)
	filePath := filepath.Join(sharedDir, filename)
	return os.Remove(filePath)
}

func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}

