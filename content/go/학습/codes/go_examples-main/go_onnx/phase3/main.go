package main

import (
	"fmt"
	"log"
	"math"
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
		log.Fatalf("Failed to initialize ORT: %v", err)
	}
	defer CleanupORT()

	modelPath := filepath.Join("..", "models", "mnist.onnx")

	// Prepare data: 1x1x28x28
	data := make([]float32, 28*28)
	// Draw a synthetic '7'
	for x := 5; x < 25; x++ {
		data[5*28+x] = 1.0
	}
	for i := 0; i < 20; i++ {
		y := 5 + i
		x := 24 - i
		if y < 28 && x >= 0 {
			data[y*28+x] = 1.0
		}
	}

	fmt.Println("Input Image (ASCII Art):")
	for y := 0; y < 28; y++ {
		for x := 0; x < 28; x++ {
			if data[y*28+x] > 0.5 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}

	predictedClass, rawOutput, err := predictMNIST(modelPath, data)
	if err != nil {
		log.Fatalf("Prediction failed: %v", err)
	}

	fmt.Printf("\nPredicted Digit: %d\n", predictedClass)
	fmt.Printf("Raw Output: %v\n", rawOutput)
}

func predictMNIST(modelPath string, data []float32) (int, []float32, error) {
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return 0, nil, fmt.Errorf("model file not found at %s", modelPath)
	}

	session, err := ort.NewAdvancedSession(modelPath,
		[]string{"input"},
		[]string{"output"},
		nil)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Destroy()

	// Create Tensor
	shape := ort.NewShape(1, 1, 28, 28)
	inputTensor, err := ort.NewTensor(shape, data)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create tensor: %w", err)
	}
	defer inputTensor.Destroy()

	outputs, err := session.Run([]ort.Value{inputTensor})
	if err != nil {
		return 0, nil, fmt.Errorf("inference failed: %w", err)
	}
	defer func() {
		for _, t := range outputs {
			t.Destroy()
		}
	}()

	outputTensor := outputs[0].(*ort.Tensor[float32])
	outputDataOriginal := outputTensor.GetData()

	// Copy data
	outputData := make([]float32, len(outputDataOriginal))
	copy(outputData, outputDataOriginal)

	maxVal := float32(-math.MaxFloat32)
	maxIdx := -1
	for i, val := range outputData {
		if val > maxVal {
			maxVal = val
			maxIdx = i
		}
	}

	return maxIdx, outputData, nil
}

func getSharedLibPath() string {
	if runtime.GOOS == "windows" {
		return "onnxruntime.dll"
	} else if runtime.GOOS == "darwin" {
		return "libonnxruntime.dylib"
	}
	return "libonnxruntime.so"
}
