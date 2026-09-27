package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	ort "github.com/yalue/onnxruntime_go"
)

// Global session for simplicity in this example
var session *ort.AdvancedSession
var ortInitialized bool

// initORT setups the environment and returns cleanup function
func initORT() (func(), error) {
	if ortInitialized {
		return func() {}, nil
	}
	ort.SetSharedLibraryPath(getSharedLibPath())
	if err := ort.InitializeEnvironment(); err != nil {
		return nil, err
	}
	ortInitialized = true
	return func() {
		ort.DestroyEnvironment()
		ortInitialized = false
	}, nil
}

// loadModel loads the global session
func loadModel(modelPath string) error {
	var err error
	session, err = ort.NewAdvancedSession(modelPath,
		[]string{"float_input"},
		[]string{"output_label", "output_probability"},
		nil)
	return err
}

type PredictionRequest struct {
	Features []float32 `json:"features"`
}

type PredictionResponse struct {
	ClassID       int64     `json:"class_id"`
	ClassName     string    `json:"class_name"`
	Probabilities []float32 `json:"probabilities"`
}

func main() {
	cleanup, err := initORT()
	if err != nil {
		log.Fatalf("Failed to initialize ORT: %v", err)
	}
	defer cleanup()

	modelPath := filepath.Join("..", "models", "iris.onnx")
	if err := loadModel(modelPath); err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}
	defer session.Destroy()

	http.HandleFunc("/predict", predictHandler)
	fmt.Println("Server listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func predictHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PredictionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if len(req.Features) != 4 {
		http.Error(w, "Expected 4 features", http.StatusBadRequest)
		return
	}

	shape := ort.NewShape(1, 4)
	inputTensor, err := ort.NewTensor(shape, req.Features)
	if err != nil {
		http.Error(w, "Failed to create tensor", http.StatusInternalServerError)
		return
	}
	defer inputTensor.Destroy()

	outputs, err := session.Run([]ort.Value{inputTensor})
	if err != nil {
		http.Error(w, "Inference failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		for _, t := range outputs {
			t.Destroy()
		}
	}()

	labelTensor := outputs[0].(*ort.Tensor[int64])
	labels := labelTensor.GetData()
	classID := labels[0]

	probTensor := outputs[1].(*ort.Tensor[float32])
	probs := probTensor.GetData()

	classNames := []string{"Setosa", "Versicolor", "Virginica"}
	className := "Unknown"
	if classID >= 0 && int(classID) < len(classNames) {
		className = classNames[classID]
	}

	resp := PredictionResponse{
		ClassID:       classID,
		ClassName:     className,
		Probabilities: probs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func getSharedLibPath() string {
	if runtime.GOOS == "windows" {
		return "onnxruntime.dll"
	} else if runtime.GOOS == "darwin" {
		return "libonnxruntime.dylib"
	}
	return "libonnxruntime.so"
}
