# Phase 1: Basics of ONNX in Go

## Goal
Load a very simple model (`simple_add.onnx`) that adds two numbers.
Understand the `onnxruntime_go` session, tensor creation, and data retrieval.

## Code Overview
The code in `main.go`:
1.  Initializes the ONNX runtime environment.
2.  Loads the `.onnx` model file.
3.  Creates Input Tensors for `X` and `Y`.
4.  Runs the session.
5.  Extracts the result from output tensor `Z`.

## How to Run

```bash
cd phase1
go run main.go
```

Expected Output:
```
Input: 1.5 + 2.5
Output (from ONNX): 4
```
