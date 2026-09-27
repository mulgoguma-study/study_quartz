package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	ort "github.com/yalue/onnxruntime_go"
)

var ortInitialized bool

func InitializeORT() error {
	if ortInitialized {
		return nil
	}
	ort.SetSharedLibraryPath(getSharedLibPath())
	err := ort.InitializeEnvironment()
	if err != nil {
		return err
	}
	ortInitialized = true
	return nil
}

func CleanupORT() {
	if ortInitialized {
		ort.DestroyEnvironment()
		ortInitialized = false
	}
}

func main() {
	if err := InitializeORT(); err != nil {
		log.Fatalf("Failed to initialize ONNX Runtime: %v", err)
	}
	defer CleanupORT()

	modelPath := filepath.Join("..", "models", "simple_add.onnx")
	result, err := runInference(modelPath, 1.5, 2.5)
	if err != nil {
		log.Fatalf("Inference failed: %v", err)
	}

	fmt.Printf("Input: 1.5 + 2.5\n")
	fmt.Printf("Output (from ONNX): %v\n", result)
}

func runInference(modelPath string, a, b float32) (float32, error) {
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return 0, fmt.Errorf("model file not found at %s", modelPath)
	}

	session, err := ort.NewAdvancedSession(modelPath,
		[]string{"X", "Y"},
		[]string{"Z"},
		nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Destroy()

	inputA := []float32{a}
	inputB := []float32{b}

	tensorA, err := ort.NewTensor(ort.NewShape(1), inputA)
	if err != nil {
		return 0, fmt.Errorf("failed to create tensor A: %w", err)
	}
	defer tensorA.Destroy()

	tensorB, err := ort.NewTensor(ort.NewShape(1), inputB)
	if err != nil {
		return 0, fmt.Errorf("failed to create tensor B: %w", err)
	}
	defer tensorB.Destroy()

	outputTensors, err := session.Run([]ort.Value{tensorA, tensorB})
	if err != nil {
		return 0, fmt.Errorf("failed to run inference: %w", err)
	}
	defer func() {
		for _, t := range outputTensors {
			t.Destroy()
		}
	}()

	outputC := outputTensors[0].(*ort.Tensor[float32])
	outputData := outputC.GetData()

	// Copy data to return regular float32
	return outputData[0], nil
}

func getSharedLibPath() string {
	if runtime.GOOS == "windows" {
		return "onnxruntime.dll"
	} else if runtime.GOOS == "darwin" {
		return "libonnxruntime.dylib"
	}
	return "libonnxruntime.so"
}
