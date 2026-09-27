# Phase 2 Guide: Structured Data

## Handling Multiple Outputs
Real-world models often return multiple pieces of information. For a classifier, this is typically the **predicted class** and the **confidence scores** (probabilities).

When you call `session.Run()`, you get a slice of `ort.Value`. You must know the order of outputs or check the model metadata. In our example (and typical `skl2onnx` exports), the order is `[Label, Probabilities]`.

```go
outputs, _ := session.Run(...)
labelTensor := outputs[0].(*ort.Tensor[int64])
probTensor := outputs[1].(*ort.Tensor[float32])
```

**Note**: You must cast the generic `ort.Tensor` to the correct type (`int64`, `float32`) matching the model's output definition. If you cast incorrectly, Go will panic.

## Batching
Although we used a batch size of 1 (`[1, 4]`), ONNX models are usually optimized for batch processing. You could pass 10 samples at once by creating a tensor of shape `[10, 4]` and providing a flat slice of 40 floats. The output would then be `[10] labels` and `[10, 3] probabilities`.

## Dealing with Maps (ZipMap)
By default, some ONNX converters export probabilities as a "Sequence of Maps". This is complex to handle in C/Go runtimes. It is generated recommended to disable `ZipMap` during export (as we did in the Python script) to get purely tensor outputs.
