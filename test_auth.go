package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

type DebugResponse struct {
	TokenUserID string `json:"token_user_id"`
	Username    string `json:"username"`
	Message     string `json:"message"`
}

func main() {
	baseURL := "https://asia-southeast2-awangga.cloudfunctions.net/wechat"
	
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run test_auth.go <username> <password>")
		return
	}
	
	username := os.Args[1]
	password := os.Args[2]

	// 1. Login untuk mendapatkan token
	fmt.Printf("=== Testing Login ===\n")
	loginReq := LoginRequest{
		Username: username,
		Password: password,
	}
	
	loginData, _ := json.Marshal(loginReq)
	resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(loginData))
	if err != nil {
		fmt.Printf("Error login: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Login Response Status: %d\n", resp.StatusCode)
	fmt.Printf("Login Response Body: %s\n", string(body))
	
	var loginResp LoginResponse
	json.Unmarshal(body, &loginResp)
	
	if loginResp.Token == "" {
		fmt.Printf("Login failed: %s\n", loginResp.Error)
		return
	}
	
	token := loginResp.Token
	fmt.Printf("Token obtained: %s\n", token[:50]+"...")
	
	// 2. Test debug endpoint untuk melihat informasi token
	fmt.Printf("\n=== Testing Debug Token ===\n")
	req, _ := http.NewRequest("GET", baseURL+"/debug/token", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Error debug token: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Debug Token Response Status: %d\n", resp.StatusCode)
	fmt.Printf("Debug Token Response Body: %s\n", string(body))
	
	var debugResp DebugResponse
	json.Unmarshal(body, &debugResp)
	
	if debugResp.TokenUserID == "" {
		fmt.Printf("Debug failed\n")
		return
	}
	
	userID := debugResp.TokenUserID
	fmt.Printf("User ID from token: %s\n", userID)
	
	// 3. Test update profile dengan ID yang sama
	fmt.Printf("\n=== Testing Update Profile ===\n")
	updateData := map[string]interface{}{
		"email":    "test@example.com",
		"fullname": "Test User Updated",
	}
	
	updateJSON, _ := json.Marshal(updateData)
	req, _ = http.NewRequest("PUT", baseURL+"/user/"+userID, bytes.NewBuffer(updateJSON))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Error update profile: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Update Profile Response Status: %d\n", resp.StatusCode)
	fmt.Printf("Update Profile Response Body: %s\n", string(body))
}
