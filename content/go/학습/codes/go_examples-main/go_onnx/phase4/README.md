# Phase 4: Serving (Web API)

## Goal
Wrap the inference logic in a Go HTTP server.
Understand how to handle concurrency and request parsing for ML services.

## Architecture
*   **Init**: Load the model **ONCE** at startup. The `session` object is thread-safe.
*   **Handler**:
    *   Accept JSON POST requests.
    *   Create Input Tensor from JSON data.
    *   Run Inference.
    *   Return JSON response.

## How to Run

1.  Start the server:
    ```bash
    cd phase4
    go run main.go
    ```

2.  Send a request (in another terminal):
    ```bash
    # PowerShell
    Invoke-RestMethod -Method Post -Uri "http://localhost:8080/predict" `
      -Body '{"features": [5.1, 3.5, 1.4, 0.2]}' `
      -ContentType "application/json"
    ```
    
    Or using `curl`:
    ```bash
    curl -X POST http://localhost:8080/predict \
      -H "Content-Type: application/json" \
      -d '{"features": [5.1, 3.5, 1.4, 0.2]}'
    ```

## Expected Response
```json
{
    "class_id": 0,
    "class_name": "Setosa",
    "probabilities": [1.0, 0.0, 0.0]
}
```
