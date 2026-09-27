# Phase 3: Image Processing (MNIST)

## Goal
Work with higher-dimensional tensors (`NCHW` format) commonly used in computer vision.
Simulate image input and interpret classification results.

## Model Details
*   **Model**: `mnist.onnx` (Simple CNN).
*   **Input**: `input` tensor of shape `[1, 1, 28, 28]`.
    *   1: Batch Size
    *   1: Channel (Grayscale)
    *   28: Height
    *   28: Width
*   **Output**: `output` tensor of shape `[1, 10]` (Log Softmax scores for digits 0-9).

## Code Highlights
*   Manually creating a flattened slice `[]float32` representing the 2D image.
*   Drawing "ASCII art" to visualize the input vector.
*   Finding the argmax of the output to determine the predicted digit.

## How to Run
```bash
cd phase3
go run main.go
```
