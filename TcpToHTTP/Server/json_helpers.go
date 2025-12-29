package main

import (
	"encoding/json"
	"fmt"
)

func parseJSONBody(body string, v interface{}) error {
	if body == "" {
		return fmt.Errorf("empty request body")
	}
	return json.Unmarshal([]byte(body), v)
}


func toJSON(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
