package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestPredictHandler(t *testing.T) {
	cleanup, err := initORT()
	if err != nil {
		t.Fatalf("Failed to initialize ORT: %v", err)
	}
	defer cleanup()

	modelPath := filepath.Join("..", "models", "iris.onnx")
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skipf("Model file %s not found, skipping handler test", modelPath)
	}

	if err := loadModel(modelPath); err != nil {
		t.Fatalf("Failed to load model: %v", err)
	}
	defer session.Destroy()

	// Prepare Request
	reqBody := PredictionRequest{
		Features: []float32{5.1, 3.5, 1.4, 0.2},
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/predict", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	// Run Handler
	predictHandler(w, req)

	// Check Response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var predResp PredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&predResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if predResp.ClassID != 0 {
		t.Errorf("Expected class 0, got %d", predResp.ClassID)
	}
	if predResp.ClassName != "Setosa" {
		t.Errorf("Expected Setosa, got %s", predResp.ClassName)
	}
}
