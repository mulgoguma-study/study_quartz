# Phase 2: Classification (Iris)

## Goal
Perform classification on tabular data using the classic Iris dataset.
Learn to handle multiple outputs (Label and Probabilities).

## Model Details
*   **Model**: `iris.onnx` (Random Forest Classifier).
*   **Input**: `float_input` (Shape `[1, 4]`).
*   **Outputs**:
    *   `output_label`: Int64 tensor containing the predicted class ID (0, 1, or 2).
    *   `output_probability`: Float32 tensor containing probabilities for each class.

## Code Highlights
*   Creating a tensor with shape `[1, 4]` to represent a batch of 1 sample with 4 features.
*   Parsing multiple outputs from the `session.Run` return slice.
*   Mapping integer labels to human-readable strings.

## How to Run
```bash
cd phase2
go run main.go
```
