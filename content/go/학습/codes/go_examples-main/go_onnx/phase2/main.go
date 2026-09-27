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
		log.Fatalf("Failed to initialize ORT: %v", err)
	}
	defer CleanupORT()

	modelPath := filepath.Join("..", "models", "iris.onnx")
	inputData := []float32{5.1, 3.5, 1.4, 0.2}

	predictedClass, probs, err := predictIris(modelPath, inputData)
	if err != nil {
		log.Fatalf("Prediction failed: %v", err)
	}

	fmt.Printf("Input Features: %v\n", inputData)
	fmt.Printf("Predicted Class ID: %d\n", predictedClass)
	fmt.Printf("Class Probabilities: %v\n", probs)

	classNames := []string{"Setosa", "Versicolor", "Virginica"}
	if predictedClass >= 0 && int(predictedClass) < len(classNames) {
		fmt.Printf("Predicted Species: %s (Confidence: %.2f%%)\n",
			classNames[predictedClass], probs[predictedClass]*100)
	}
}

func predictIris(modelPath string, inputData []float32) (int64, []float32, error) {
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return 0, nil, fmt.Errorf("model file not found at %s", modelPath)
	}

	session, err := ort.NewAdvancedSession(modelPath,
		[]string{"float_input"},
		[]string{"output_label", "output_probability"},
		nil)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Destroy()

	tensorShape := ort.NewShape(1, 4)
	inputTensor, err := ort.NewTensor(tensorShape, inputData)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create input tensor: %w", err)
	}
	defer inputTensor.Destroy()

	outputs, err := session.Run([]ort.Value{inputTensor})
	if err != nil {
		return 0, nil, fmt.Errorf("error running inference: %w", err)
	}
	defer func() {
		for _, t := range outputs {
			t.Destroy()
		}
	}()

	labelTensor := outputs[0].(*ort.Tensor[int64])
	labels := labelTensor.GetData()

	probTensor := outputs[1].(*ort.Tensor[float32])
	probsOriginal := probTensor.GetData()

	probsCopy := make([]float32, len(probsOriginal))
	copy(probsCopy, probsOriginal)

	return labels[0], probsCopy, nil
}

func getSharedLibPath() string {
	if runtime.GOOS == "windows" {
		return "onnxruntime.dll"
	} else if runtime.GOOS == "darwin" {
		return "libonnxruntime.dylib"
	}
	return "libonnxruntime.so"
}
