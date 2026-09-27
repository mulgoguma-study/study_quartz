# Go ONNX Learning Journey

This repository contains a step-by-step guide to learning how to use ONNX models with Go.

## Prerequisites

1.  **Go**: Installed on your system.
    *   **Important**: This project uses `CGO`. You must have a C compiler installed (e.g., **MinGW-w64** on Windows, `gcc` on Linux/Mac).
2.  **ONNX Runtime Shared Library**:
    *   You **MUST** download the ONNX Runtime shared library (`onnxruntime.dll` for Windows, `libonnxruntime.so` for Linux/Mac) for your architecture.
    *   Download from: [https://github.com/microsoft/onnxruntime/releases](https://github.com/microsoft/onnxruntime/releases) (Look for `onnxruntime-win-x64-1.x.x.zip` or similar).
    *   Extract the zip and place the `onnxruntime.dll` (or `.so`/`.dylib`) in the root of this project or in your system path (e.g., `C:\Windows\System32` or add to `$PATH`).
    *   **Phase Folders**: Tests run inside phase folders, so verify if `.dll` needs to be accessible there (or simpler, put it in system PATH).
3.  **Python** (for generating models):
    *   We use Python to generate the initial ONNX models used in the examples.
    *   Install dependencies:
        ```bash
        pip install onnx numpy scikit-learn skl2onnx torch
        ```

## Setup

Before running any Go code, generate the ONNX models:

```bash
python utils/generate_models.py
```

This will create:
*   `models/simple_add.onnx`
*   `models/iris.onnx`
*   `models/mnist.onnx`

## Running Tests

We have included unit tests for each phase to verify the logic.
Running tests requires the ONNX models to be generated first.

```bash
go test -v ./...
```

*Note: If models are missing, tests will be skipped automatically.*

## Phases

### [Phase 1: Basics](phase1/)
Learn how to load a simple model and run inference with basic scalar/vector inputs.

### [Phase 2: Structured Data (Iris)](phase2/)
Handle structured tabular data, dealing with classifiers, and interpreting outputs (labels and probabilities).

### [Phase 3: Image Processing (MNIST)](phase3/)
Working with higher-dimensional tensors (images), preprocessing image data in Go, and running an image classification model.

### [Phase 4: Serving (Web API)](phase4/)
Wrap your ONNX inference logic into a production-ready HTTP web service.
