package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPredictMNIST(t *testing.T) {
	if err := InitializeORT(); err != nil {
		t.Fatalf("Failed to initialize ORT: %v", err)
	}
	defer CleanupORT()

	modelPath := filepath.Join("..", "models", "mnist.onnx")
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skipf("Model file %s not found, skipping inference test", modelPath)
	}

	// Create a blank image
	data := make([]float32, 28*28)

	// Predict
	digit, probs, err := predictMNIST(modelPath, data)
	if err != nil {
		t.Fatalf("predictMNIST failed: %v", err)
	}

	if digit < 0 || digit > 9 {
		t.Errorf("Expected digit between 0-9, got %d", digit)
	}
	if len(probs) != 10 {
		t.Errorf("Expected 10 probabilities, got %d", len(probs))
	}
}
