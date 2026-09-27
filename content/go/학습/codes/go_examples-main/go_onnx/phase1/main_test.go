package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInference(t *testing.T) {
	// Initialize environment for tests
	if err := InitializeORT(); err != nil {
		t.Fatalf("Failed to initialize ORT: %v", err)
	}
	defer CleanupORT()

	modelPath := filepath.Join("..", "models", "simple_add.onnx")

	// Skip if model doesn't exist (user hasn't run generation script)
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skipf("Model file %s not found, skipping inference test", modelPath)
	}

	tests := []struct {
		a, b     float32
		expected float32
	}{
		{1.0, 2.0, 3.0},
		{1.5, 2.5, 4.0},
		{-1.0, 1.0, 0.0},
	}

	for _, tc := range tests {
		result, err := runInference(modelPath, tc.a, tc.b)
		if err != nil {
			t.Errorf("runInference(%v, %v) failed: %v", tc.a, tc.b, err)
			continue
		}
		if result != tc.expected {
			t.Errorf("runInference(%v, %v) = %v; want %v", tc.a, tc.b, result, tc.expected)
		}
	}
}
