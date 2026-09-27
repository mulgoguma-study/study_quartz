## 텐서의 개념

**기본 이해:**

- 스칼라(0차원) → 벡터(1차원) → 행렬(2차원) → **텐서(N차원)**
- 텐서는 다차원 배열의 일반화된 형태입니다

```
스칼라: 5
벡터: [1, 2, 3]
행렬: [[1, 2], [3, 4]]
3D 텐서: [[[1, 2], [3, 4]], [[5, 6], [7, 8]]]
```

## Shape (형태)
텐서의 각 차원 크기를 나타냅니다.

go

```go
// 이미지 예시 (NCHW 형식)
shape := []int64{1, 3, 224, 224}
// N: Batch size (1개 이미지)
// C: Channels (RGB 3채널)
// H: Height (224픽셀)
// W: Width (224픽셀)

// 실제 데이터 크기 = 1 × 3 × 224 × 224 = 150,528 개 요소
```

**실제 모델 예시:**

go

```go
// 얼굴 인식 모델
input_shape := [4]int64{1, 3, 112, 112}  // 112x112 RGB 이미지
output_shape := [2]int64{1, 512}         // 512차원 얼굴 특징 벡터

// 텍스트 임베딩 모델
input_shape := [2]int64{1, 128}          // 최대 128 토큰
output_shape := [2]int64{1, 768}         // 768차원 임베딩

// 객체 탐지 모델
input_shape := [4]int64{1, 3, 640, 640}  // 640x640 이미지
output_shape := [3]int64{1, 25200, 85}   // 25200개 박스, 각 85개 값
```

### DType (데이터 타입)

텐서 요소의 데이터 타입입니다.

go

```go
// ONNX Runtime에서 지원하는 주요 타입
const (
    TensorElementDataTypeFloat    = 1  // float32
    TensorElementDataTypeUint8    = 2  // uint8
    TensorElementDataTypeInt8     = 3  // int8
    TensorElementDataTypeUint16   = 4  // uint16
    TensorElementDataTypeInt16    = 5  // int16
    TensorElementDataTypeInt32    = 6  // int32
    TensorElementDataTypeInt64    = 7  // int64
    TensorElementDataTypeDouble   = 11 // float64
)
```

**실제 사용 예시:**

go

```go
// 이미지는 보통 float32 (정규화된 픽셀 값 0.0~1.0)
inputData := make([]float32, 1*3*224*224)
dtype := TensorElementDataTypeFloat

// 토큰 ID는 int64
tokenIds := make([]int64, 1*128)
dtype := TensorElementDataTypeInt64

// 양자화 모델은 uint8 또는 int8 사용
quantizedData := make([]uint8, 1*3*224*224)
dtype := TensorElementDataTypeUint8
```

### Go 코드 예시

go

```go
type ModelConfig struct {
    // 입력 텐서 정보
    InputName  string
    InputShape []int64      // [1, 3, 224, 224]
    InputType  int          // float32
    
    // 출력 텐서 정보
    OutputName  string
    OutputShape []int64     // [1, 1000] - 1000개 클래스
    OutputType  int         // float32
}

// 실제 텐서 생성
func createInputTensor(image []byte, config ModelConfig) (*onnxruntime.Tensor, error) {
    // 1. 이미지를 float32 배열로 변환
    floatData := preprocessImage(image, config.InputShape)
    
    // 2. 텐서 생성
    tensor, err := onnxruntime.NewTensor(
        config.InputShape,  // shape
        floatData,          // data
    )
    
    return tensor, err
}

// 전처리 예시
func preprocessImage(image []byte, shape []int64) []float32 {
    h, w := int(shape[2]), int(shape[3])
    
    // Resize → Normalize → HWC to CHW 변환
    floatData := make([]float32, 1*3*h*w)
    
    // RGB 각 채널별로 정규화
    for c := 0; c < 3; c++ {
        for y := 0; y < h; y++ {
            for x := 0; x < w; x++ {
                idx := c*h*w + y*w + x
                pixel := getPixel(image, x, y, c)
                floatData[idx] = (float32(pixel) - 127.5) / 127.5
            }
        }
    }
    
    return floatData
}
```

### 모델마다 다른 이유

go

```go
// ResNet (이미지 분류)
Input:  [1, 3, 224, 224] float32
Output: [1, 1000] float32

// YOLO (객체 탐지)
Input:  [1, 3, 640, 640] float32
Output: [1, 25200, 85] float32  // bbox 좌표 + 클래스 확률

// BERT (텍스트)
Input:  [1, 512] int64           // 토큰 IDs
Output: [1, 512, 768] float32    // 각 토큰의 임베딩

// FaceNet (얼굴 인식)
Input:  [1, 3, 160, 160] float32
Output: [1, 128] float32         // 128차원 얼굴 벡터
```

**핵심은**: shape와 dtype만 설정 파일이나 config로 관리하면, 나머지 Pool/Session/Buffer 구조는 재사용 가능합니다!