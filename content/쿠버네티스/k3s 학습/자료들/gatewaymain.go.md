---
aliases:
  - gateway/main.go
  - gateway/main.go
---
```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Service string      `json:"service"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
}

func main() {
	userServiceURL := getEnv("USER_SERVICE_URL", "http://user-service:8081")
	productServiceURL := getEnv("PRODUCT_SERVICE_URL", "http://product-service:8082")
	port := getEnv("PORT", "8080")

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/users", proxyHandler(userServiceURL, "/users"))
	http.HandleFunc("/products", proxyHandler(productServiceURL, "/products"))
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"service": "API Gateway",
			"version": "1.0.0",
			"endpoints": "/users, /products, /health",
		})
	})

	log.Printf("Gateway starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func proxyHandler(serviceURL, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{Timeout: 5 * time.Second}
		
		resp, err := client.Get(serviceURL + path)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(Response{
				Service: "gateway",
				Error:   fmt.Sprintf("Service unavailable: %v", err),
			})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
```