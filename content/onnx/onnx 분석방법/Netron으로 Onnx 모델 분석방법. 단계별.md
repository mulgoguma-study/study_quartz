Netron으로 ONNX 모델을 분석하는 방법을 단계별로 알려드리겠습니다!

## 1단계: Netron 기본 UI 이해하기

### 주요 영역

```
┌─────────────────────────────────────┐
│  [Model Properties]  (상단)          │  ← 모델 전체 정보
├─────────────────────────────────────┤
│                                     │
│  [Graph Visualization] (중앙)       │  ← 노드/레이어 시각화
│       ↓                             │
│     [Node]                          │
│       ↓                             │
│     [Node]                          │
│                                     │
├─────────────────────────────────────┤
│  [Node Properties]  (측면 패널)      │  ← 선택한 노드 상세정보
└─────────────────────────────────────┘
```

## 2단계: 모델 입출력 확인 (가장 중요!)

### 모델 최상단 보기

```
┌──────────────────────────────────┐
│ MODEL PROPERTIES                 │
├──────────────────────────────────┤
│ Format: ONNX v1.13              │
│ Producer: pytorch v2.0.0         │
│                                  │
│ INPUTS (여기 주목!)               │
│ └─ input                         │
│    └─ float32[1,3,224,224]      │  ← shape와 dtype
│                                  │
│ OUTPUTS (여기도 주목!)            │
│ └─ output                        │
│    └─ float32[1,1000]           │
└──────────────────────────────────┘
```

**Go 코드로 변환:**

```go
type ModelConfig struct {
    InputName:  "input",
    InputShape: []int64{1, 3, 224, 224},
    InputType:  TensorElementDataTypeFloat,
    
    OutputName:  "output",
    OutputShape: []int64{1, 1000},
    OutputType:  TensorElementDataTypeFloat,
}
```

## 3단계: 그래프 읽는 법

### 기본 구조

```
Input Tensor
    ↓
[Conv] ← 여기 클릭하면 측면에 상세정보
    ↓
[BatchNorm]
    ↓
[ReLU]
    ↓
[MaxPool]
    ↓
... (반복)
    ↓
[Dense/Gemm]
    ↓
Output Tensor
```

### 노드 클릭 시 보이는 정보

```
┌─────────────────────────────────┐
│ NODE PROPERTIES                 │
├─────────────────────────────────┤
│ Type: Conv                      │  ← 연산 타입
│ Name: Conv_0                    │
│                                 │
│ ATTRIBUTES                      │  ← 중요!
│ ├─ kernel_shape: [3, 3]        │
│ ├─ strides: [1, 1]             │
│ ├─ pads: [1, 1, 1, 1]          │
│ └─ dilations: [1, 1]           │
│                                 │
│ INPUTS                          │
│ ├─ input: float32[1,3,224,224] │
│ └─ weight: float32[64,3,3,3]   │  ← 가중치 shape
│                                 │
│ OUTPUTS                         │
│ └─ output: float32[1,64,224,224]│
└─────────────────────────────────┘
```

## 4단계: 실전 분석 예제

### 예제 1: 이미지 분류 모델 (ResNet)

```
[INPUT] float32[1,3,224,224]
   ↓
[Conv] kernel=7x7, stride=2
   ↓  float32[1,64,112,112]
[BatchNorm]
   ↓
[ReLU]
   ↓
[MaxPool] kernel=3x3, stride=2
   ↓  float32[1,64,56,56]
[ResBlock] × 3
   ↓
[ResBlock] × 4
   ↓
[ResBlock] × 6
   ↓
[GlobalAveragePool]
   ↓  float32[1,2048]
[Gemm/Dense] (2048 → 1000)
   ↓
[OUTPUT] float32[1,1000]
```

**분석 포인트:**

1. ✅ Input shape: `[1,3,224,224]` → 224×224 RGB 이미지 필요
2. ✅ Output shape: `[1,1000]` → 1000개 클래스 분류
3. ✅ 전처리: ImageNet 정규화 필요 (mean/std)

### 예제 2: 객체 탐지 모델 (YOLO)

```
[INPUT] float32[1,3,640,640]
   ↓
[Backbone: CSPDarknet]
   ↓
[Neck: PANet]
   ↓
[Head: Detection layers]
   ↓
   ├─ [Small objects]  → float32[1,80,80,85]
   ├─ [Medium objects] → float32[1,40,40,85]
   └─ [Large objects]  → float32[1,20,20,85]
   ↓
[Concat/Reshape]
   ↓
[OUTPUT] float32[1,25200,85]
```

**85의 의미:**

```
[x, y, w, h, confidence, class1, class2, ..., class80]
 ←─ 4 ───→ ←─ 1 ─→     ←────────── 80 ──────────→
```

## 5단계: 실용적인 분석 체크리스트

### ✅ 필수 확인 사항

```
□ Input 정보
  └─ Name: ___________
  └─ Shape: [_, _, _, _]
  └─ Type: float32/int64/...
  
□ Output 정보
  └─ Name: ___________
  └─ Shape: [_, _, _, _]
  └─ Type: float32/...
  
□ 동적 차원 확인
  └─ batch_size가 -1 또는 'N'인가?
  └─ 가변 길이 입력인가?
  
□ 특수 전처리 요구사항
  └─ 정규화 범위: [0,1] or [-1,1]?
  └─ RGB or BGR?
  └─ Channel order: CHW or HWC?
```

### Go 코드 작성 템플릿

```go
// Netron에서 확인한 정보를 바로 코드로
type ModelMetadata struct {
    // ===== Netron 상단에서 확인 =====
    InputName:  "images",              // Input 노드 이름
    InputShape: []int64{1, 3, 640, 640}, // Input shape
    InputType:  C.ONNX_TENSOR_ELEMENT_DATA_TYPE_FLOAT,
    
    OutputName:  "output0",            // Output 노드 이름
    OutputShape: []int64{1, 25200, 85}, // Output shape
    OutputType:  C.ONNX_TENSOR_ELEMENT_DATA_TYPE_FLOAT,
    
    // ===== 모델 문서나 학습 코드에서 확인 =====
    Mean: []float32{0.0, 0.0, 0.0},
    Std:  []float32{255.0, 255.0, 255.0},
    
    // ===== 그래프 분석으로 추론 =====
    RequiresBGR: false,  // OpenCV 모델이면 true
    RequiresCHW: true,   // 대부분 true
}
```

## 6단계: 복잡한 모델 읽는 법

### 다중 입력/출력 모델

```
┌──────────────────────────────────┐
│ INPUTS                           │
│ ├─ input_ids: int64[1,512]      │  ← BERT 토큰
│ ├─ attention_mask: int64[1,512] │
│ └─ token_type_ids: int64[1,512] │
│                                  │
│ OUTPUTS                          │
│ ├─ last_hidden_state: [1,512,768]│
│ └─ pooler_output: [1,768]       │
└──────────────────────────────────┘
```

**Go 코드:**

```go
type BERTInputs struct {
    InputIDs      []int64  // shape: [1, 512]
    AttentionMask []int64  // shape: [1, 512]
    TokenTypeIDs  []int64  // shape: [1, 512]
}

type BERTOutputs struct {
    LastHiddenState []float32  // shape: [1, 512, 768]
    PoolerOutput    []float32  // shape: [1, 768]
}
```

## 7단계: 디버깅 팁

### Shape 불일치 디버깅

```go
// Netron에서 본 shape
expected := []int64{1, 3, 224, 224}

// 실제 입력 데이터
actualSize := len(imageData)  // 150528

// 검증
expectedSize := 1 * 3 * 224 * 224  // 150528
if actualSize != expectedSize {
    log.Fatalf("Shape mismatch: expected %d, got %d", 
               expectedSize, actualSize)
}
```

### 중간 레이어 출력 확인

Netron에서 중간 노드 이름을 찾아서:

```go
// 특정 레이어까지만 실행
outputNames := []string{
    "Conv_0_output",  // Netron에서 확인한 중간 노드
    "final_output",
}
```

## 8단계: 학습 로드맵

### Week 1: 기초

- [ ] 간단한 분류 모델 (ResNet) 분석
- [ ] Input/Output shape 확인 연습
- [ ] 전처리 파이프라인 구현

### Week 2: 중급

- [ ] 객체 탐지 모델 (YOLO) 분석
- [ ] 다중 출력 이해
- [ ] 후처리 로직 구현

### Week 3: 고급

- [ ] Transformer 모델 (BERT) 분석
- [ ] 다중 입력 처리
- [ ] 동적 shape 처리

## 실전 연습 문제

**모델 다운로드 후 분석해보기:**

```bash
# 1. 간단한 모델
wget https://github.com/onnx/models/raw/main/vision/classification/resnet/model/resnet50-v1-7.onnx

# 2. Netron에서 열기
# 3. 다음 정보 찾기:
#    - Input name & shape
#    - Output name & shape
#    - 첫 Conv layer의 kernel size
#    - 마지막 Dense layer의 출력 클래스 수
```

이렇게 실제 모델로 연습하면서 Netron 보는 눈이 생깁니다!