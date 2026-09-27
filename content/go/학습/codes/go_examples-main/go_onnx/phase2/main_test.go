package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPredictIris(t *testing.T) {
	if err := InitializeORT(); err != nil {
		t.Fatalf("Failed to initialize ORT: %v", err)
	}
	defer CleanupORT()

	modelPath := filepath.Join("..", "models", "iris.onnx")

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skipf("Model file %s not found, skipping inference test", modelPath)
	}

	// Test Setosa case
	input := []float32{5.1, 3.5, 1.4, 0.2}
	label, probs, err := predictIris(modelPath, input)
	if err != nil {
		t.Fatalf("predictIris failed: %v", err)
	}

	if label != 0 {
		t.Errorf("Expected label 0 (Setosa), got %d", label)
	}
	if len(probs) != 3 {
		t.Errorf("Expected 3 probabilities, got %d", len(probs))
	}
	if probs[0] < 0.5 {
		t.Errorf("Expected high probability for class 0, got %v", probs[0])
	}
}
