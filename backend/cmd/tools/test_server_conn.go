package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("http://localhost:8080/api/products")
	if err != nil {
		fmt.Println("Error connecting to http://localhost:8080:", err)
		return
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d, Response length: %d\n", resp.StatusCode, len(data))
	if len(data) > 200 {
		fmt.Printf("Snippet: %s...\n", string(data[:200]))
	} else {
		fmt.Printf("Response: %s\n", string(data))
	}
}
