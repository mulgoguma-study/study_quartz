# Phase 1 Guide: The Foundation

## What is ONNX?
ONNX (Open Neural Network Exchange) is an open standard format for representing machine learning models. It allows models trained in frameworks like PyTorch, TensorFlow, or Scikit-Learn to be run in a variety of environments, including Go.

## The ONNX Runtime
To run ONNX models, we use the **ONNX Runtime (ORT)**. This is a high-performance engine (written in C++) that executes the model. 
In Go, we use a wrapper library (bindings) to talk to this C++ engine.

### Key Concepts in Go

1.  **Environment**: 
    You must initialize the ORT environment once. This sets up threading and logging.
    ```go
    ort.InitializeEnvironment()
    defer ort.DestroyEnvironment()
    ```

2.  **Session**:
    A session loads a specific model and holds its state. You run inference through a session.
    ```go
    session, err := ort.NewAdvancedSession("model.onnx", inputNames, outputNames, options)
    ```

3.  **Tensors**:
    Data is passed to/from the model as Tensors. A tensor is a multi-dimensional array (like a NumPy array or PyTorch tensor).
    In `onnxruntime_go`, you create tensors from Go slices.
    - **Shape**: The dimensions of the tensor (e.g., `[1, 3, 224, 224]` for an image batch).
    - **Data**: A flat slice of data (e.g., `[]float32`).
    
    ```go
    shape := ort.NewShape(1, 4) // Batch size 1, 4 features
    data := []float32{1.0, 2.0, 3.0, 4.0}
    tensor, _ := ort.NewTensor(shape, data)
    ```

4.  **Inference**:
    You call `session.Run()` with a list of input tensors. It returns a list of output tensors.
    Make sure to `Destroy()` tensors when done to free C++ memory!

## Common Pitfalls
*   **Shared Library**: The most common error is "DLL not found" or "Shared library not loaded". Ensure `onnxruntime.dll` (Windows) or `libonnxruntime.so` (Linux) is in the correct path.
*   **Tensor Shapes**: If your model expects `[1, 10]` but you provide `[10]`, it might fail. Always check the model's expected input shape (you can use Netron.app to view `.onnx` files).
*   **Data Types**: If the model expects `Float32`, passing `Float64` will cause a crash or error.
