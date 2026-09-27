# Phase 4 Guide: Building an ML Microservice

## Thread Safety and Performance
One of the biggest advantages of `onnxruntime` is that the `Session` object is **thread-safe**.
This means you can calculate inference separate requests on different goroutines without needing a mutex lock around the `Run()` call (internally, the C++ library handles this).

### Optimization Tips
1.  **Memory Reuse**: In this example, we create a new `ort.Tensor` and `ort.Shape` for every request. In extremely high-load systems, this might generate GC pressure. You might consider using a `sync.Pool` for reusable buffers if inputs are fixed size.
2.  **Batching**: If you receive high traffic, "dynamic batching" is a common pattern. You collect requests for 10-50ms and send them to the model as a batch (Batch Size > 1). This leverages vectorization efficiently (checking 10 items is cheaper than 1 check * 10 times).

## JSON vs Protobuf
We used JSON for simplicity (`encoding/json`). In production internal services, gRPC with Protobuf is often preferred for speed and type safety, especially when sending large tensors (images/audio).

## Deployment
*   **Docker**: When containerizing this, ensure you include the `onnxruntime.so` / `.dll` in the image. Alpine Linux often requires `glibc` compatibility layers or a specific build of ONNX Runtime; usually, an Ubuntu/Debian base image is easier to start with.
