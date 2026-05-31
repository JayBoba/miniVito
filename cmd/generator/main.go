package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func main() {
	baseURL := "http://localhost:8080"
	client := &http.Client{Timeout: 5 * time.Second}

	for {
		login := fmt.Sprintf("user_%d", time.Now().UnixNano())
		password := "testpass123"

		user := User{
			Login:    login,
			Password: password,
		}

		body, err := json.Marshal(user)
		if err != nil {
			fmt.Printf("[Generator] JSON error: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		respReg, err := client.Post(baseURL+"/register", "application/json", bytes.NewBuffer(body))
		if err == nil {
			fmt.Printf("[Generator] Register %s: Status %d\n", login, respReg.StatusCode)
			respReg.Body.Close()
		} else {
			fmt.Printf("[Generator] Register error: %v\n", err)
		}

		time.Sleep(500 * time.Millisecond)

		respLog, err := client.Post(baseURL+"/login", "application/json", bytes.NewBuffer(body))
		if err == nil {
			fmt.Printf("[Generator] Login %s: Status %d\n", login, respLog.StatusCode)
			respLog.Body.Close()
		} else {
			fmt.Printf("[Generator] Login error: %v\n", err)
		}

		time.Sleep(1 * time.Second)
	}
}