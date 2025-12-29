package main

import "strings"
import "os"
import "encoding/base64"
import "strconv"


func readStaticFile(candidates ...string) string {
	for _, p := range candidates {
		if b, err := os.ReadFile(p); err == nil {
			return string(b)
		}
	}
	return ""
}

func urlDecode(s string) string {
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+2 < len(s) {
			if val, err := strconv.ParseUint(s[i+1:i+3], 16, 8); err == nil {
				result = append(result, byte(val))
				i += 2
				continue
			}
		}
		result = append(result, s[i])
	}
	return string(result)
}

func routeRequest(req *HTTPrequest) (int, string) {
	if req.Method == "OPTIONS" {
		return 200, ""
	}
	switch true {

	case req.Path == "/" || req.Path == "/index.html":
		body := readStaticFile(
			"Server/static/index.html",
			"static/index.html",
			"UI/index.html",
			"../UI/index.html",
		)
		return 200, body

	case req.Path == "/styles.css":
		body := readStaticFile(
			"Server/static/styles.css",
			"static/styles.css",
			"UI/styles.css",
			"../UI/styles.css",
		)
		return 200, body


	case req.Path == "/favicon.ico":
		return 200, ""

	case req.Path == "/files":
		switch req.Method {
		case "GET":
			files, err := listFiles()
			if err != nil {
				return 500, "Error listing files"
			}
			jsonResp, _ := toJSON(files)
			return 200, jsonResp
		default:
			return 405, "Method Not Allowed"
		}

	case strings.HasPrefix(req.Path, "/files/"):
		filename := strings.TrimPrefix(req.Path, "/files/")
		if idx := strings.Index(filename, "?"); idx != -1 {
			filename = filename[:idx]
		}
		filename = urlDecode(filename)
		
		switch req.Method {
		case "GET":
			if strings.Contains(req.Path, "?download=true") {
				return 200, ""
			}
			data, err := downloadFile(filename)
			if err != nil {
				return 404, "File not found"
			}
			encoded := base64.StdEncoding.EncodeToString(data)
			jsonResp := `{"content":"` + encoded + `","filename":"` + filename + `"}`
			return 200, jsonResp
			
		case "DELETE":
			err := deleteFile(filename)
			if err != nil {
				return 404, "File not found"
			}
			return 204, ""
			
		default:
			return 405, "Method Not Allowed"
		}

	case req.Path == "/upload":
		if req.Method == "POST" {
			var input struct {
				Filename string `json:"filename"`
				Content  string `json:"content"`
			}
			if err := parseJSONBody(req.Body, &input); err != nil {
				return 400, `{"error":"Invalid JSON body: ` + err.Error() + `"}`
			}
			if input.Filename == "" {
				return 400, `{"error":"Filename is required"}`
			}
			if input.Content == "" {
				return 400, `{"error":"File content is required"}`
			}
			if err := uploadFile(input.Filename, input.Content); err != nil {
				return 500, `{"error":"Error uploading file: ` + err.Error() + `"}`
			}
			return 201, `{"message":"File uploaded successfully"}`
		}
		return 405, `{"error":"Method Not Allowed"}`

	default:
		return 404, "Not Found"
	}
}