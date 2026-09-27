네! Netron보다 더 쉽고 실용적인 방법들이 있습니다.

## 1. Python으로 모델 정보 추출 (가장 추천!)

### 방법 A: onnx 라이브러리 사용

```python
# inspect_model.py
import onnx

model = onnx.load("model.onnx")

# ===== 입력 정보 =====
print("=== INPUTS ===")
for input in model.graph.input:
    print(f"Name: {input.name}")
    print(f"Type: {input.type.tensor_type.elem_type}")
    shape = []
    for dim in input.type.tensor_type.shape.dim:
        shape.append(dim.dim_value if dim.dim_value > 0 else -1)
    print(f"Shape: {shape}")
    print()

# ===== 출력 정보 =====
print("=== OUTPUTS ===")
for output in model.graph.output:
    print(f"Name: {output.name}")
    shape = []
    for dim in output.type.tensor_type.shape.dim:
        shape.append(dim.dim_value if dim.dim_value > 0 else -1)
    print(f"Shape: {shape}")
    print()

# ===== 전체 레이어 개수 =====
print(f"Total nodes: {len(model.graph.node)}")
print(f"Total inputs: {len(model.graph.input)}")
print(f"Total outputs: {len(model.graph.output)}")
```

**실행:**

```bash
pip install onnx
python inspect_model.py

# 출력 예시:
# === INPUTS ===
# Name: images
# Type: 1 (float32)
# Shape: [1, 3, 640, 640]
#
# === OUTPUTS ===
# Name: output0
# Shape: [1, 25200, 85]
```

### 방법 B: ONNX Runtime으로 직접 확인

```python
# quick_check.py
import onnxruntime as ort
import numpy as np

session = ort.InferenceSession("model.onnx")

# 입력 정보
print("=== INPUTS ===")
for inp in session.get_inputs():
    print(f"{inp.name}: {inp.type} {inp.shape}")

# 출력 정보
print("\n=== OUTPUTS ===")
for out in session.get_outputs():
    print(f"{out.name}: {out.type} {out.shape}")

# 실제 추론 테스트 (더미 데이터)
print("\n=== TEST INFERENCE ===")
dummy_input = np.random.randn(1, 3, 640, 640).astype(np.float32)
input_name = session.get_inputs()[0].name
output = session.run(None, {input_name: dummy_input})
print(f"Output shape: {output[0].shape}")
```

**이게 제일 간단하고 실용적입니다!**

---

## 2. CLI 도구 사용

### onnx-tool (가볍고 빠름)

```bash
pip install onnx-tool

# 모델 요약
onnx-tool model.onnx

# 출력 예시:
# Input:
#   images: float32[1,3,640,640]
# Output:
#   output0: float32[1,25200,85]
# 
# Operators: 270
# Parameters: 7.2M
```

### polygraphy (NVIDIA 도구)

```bash
pip install polygraphy

# 상세 정보
polygraphy inspect model model.onnx

# 더미 데이터로 실행 테스트
polygraphy run model.onnx --onnxrt
```

---

## 3. 자동 Go 코드 생성 스크립트

### 최고의 방법: Python → Go config 자동 생성

```python
# gen_go_config.py
import onnx
import sys

def generate_go_config(model_path):
    model = onnx.load(model_path)
    
    # 입력 정보 추출
    input_info = model.graph.input[0]
    input_name = input_info.name
    input_shape = [dim.dim_value for dim in input_info.type.tensor_type.shape.dim]
    
    # 출력 정보 추출
    output_info = model.graph.output[0]
    output_name = output_info.name
    output_shape = [dim.dim_value for dim in output_info.type.tensor_type.shape.dim]
    
    # Go 코드 생성
    go_code = f'''package config

type ModelConfig struct {{
    ModelPath string
    
    // Input configuration
    InputName  string
    InputShape []int64
    
    // Output configuration
    OutputName  string
    OutputShape []int64
}}

func NewDefaultConfig() *ModelConfig {{
    return &ModelConfig{{
        ModelPath:   "{model_path}",
        InputName:   "{input_name}",
        InputShape:  []int64{{{', '.join(map(str, input_shape))}}},
        OutputName:  "{output_name}",
        OutputShape: []int64{{{', '.join(map(str, output_shape))}}},
    }}
}}
'''
    
    # 파일로 저장
    with open('model_config.go', 'w') as f:
        f.write(go_code)
    
    print("✅ Generated: model_config.go")
    print(f"   Input:  {input_name} {input_shape}")
    print(f"   Output: {output_name} {output_shape}")

if __name__ == "__main__":
    generate_go_config(sys.argv[1])
```

**사용법:**

```bash
python gen_go_config.py yolov8.onnx

# 자동으로 model_config.go 생성됨!
```

---

## 4. 웹 기반 간단 뷰어

### Gradio 대시보드 만들기

```python
# model_viewer.py
import gradio as gr
import onnxruntime as ort
import numpy as np
from PIL import Image

def analyze_model(model_file):
    session = ort.InferenceSession(model_file.name)
    
    inputs = []
    for inp in session.get_inputs():
        inputs.append(f"📥 {inp.name}: {inp.type} {inp.shape}")
    
    outputs = []
    for out in session.get_outputs():
        outputs.append(f"📤 {out.name}: {out.type} {out.shape}")
    
    return "\n".join(inputs), "\n".join(outputs)

def test_inference(model_file, image):
    session = ort.InferenceSession(model_file.name)
    
    # 이미지 전처리
    img = image.resize((640, 640))
    img_array = np.array(img).transpose(2, 0, 1)  # HWC → CHW
    img_array = img_array.astype(np.float32) / 255.0
    img_array = np.expand_dims(img_array, 0)
    
    # 추론
    input_name = session.get_inputs()[0].name
    output = session.run(None, {input_name: img_array})
    
    return f"Output shape: {output[0].shape}\nSample values: {output[0].flatten()[:10]}"

with gr.Blocks() as demo:
    gr.Markdown("# ONNX Model Analyzer")
    
    with gr.Tab("Model Info"):
        model_upload = gr.File(label="Upload ONNX Model")
        analyze_btn = gr.Button("Analyze")
        input_info = gr.Textbox(label="Inputs", lines=5)
        output_info = gr.Textbox(label="Outputs", lines=5)
        
        analyze_btn.click(analyze_model, model_upload, [input_info, output_info])
    
    with gr.Tab("Test Inference"):
        model_upload2 = gr.File(label="Upload ONNX Model")
        image_input = gr.Image(type="pil")
        test_btn = gr.Button("Run Inference")
        result_output = gr.Textbox(label="Result", lines=10)
        
        test_btn.click(test_inference, [model_upload2, image_input], result_output)

demo.launch()
```

**실행:**

```bash
pip install gradio onnxruntime pillow
python model_viewer.py

# 브라우저에서 http://localhost:7860 열림
# 드래그 앤 드롭으로 모델 분석!
```

---

## 5. VSCode Extension

### ONNX Viewer Extension

```
1. VSCode 열기
2. Extensions → 검색: "ONNX"
3. "ONNX Viewer" 설치
4. .onnx 파일 클릭 → 자동으로 그래프 표시
```

**장점:**

- 코드 편집하면서 바로 확인
- 별도 앱 필요 없음
- 가볍고 빠름

---

## 6. 완전 자동화 워크플로우

### Makefile로 통합

```makefile
# Makefile
.PHONY: analyze generate test

MODEL ?= model.onnx

# 1단계: 모델 분석
analyze:
	@python -c "import onnxruntime as ort; \
	sess = ort.InferenceSession('$(MODEL)'); \
	print('INPUT:', sess.get_inputs()[0].name, sess.get_inputs()[0].shape); \
	print('OUTPUT:', sess.get_outputs()[0].name, sess.get_outputs()[0].shape)"

# 2단계: Go 설정 파일 생성
generate:
	@python gen_go_config.py $(MODEL)

# 3단계: 더미 데이터로 테스트
test:
	@python -c "import onnxruntime as ort; import numpy as np; \
	sess = ort.InferenceSession('$(MODEL)'); \
	inp = sess.get_inputs()[0]; \
	dummy = np.random.randn(*[d if d > 0 else 1 for d in inp.shape]).astype(np.float32); \
	out = sess.run(None, {inp.name: dummy}); \
	print('✅ Test passed! Output shape:', out[0].shape)"

# 모든 단계 실행
all: analyze generate test
```

**사용법:**

```bash
# 한 번에 모든 분석 완료
make all MODEL=yolov8.onnx

# 출력:
# INPUT: images [1, 3, 640, 640]
# OUTPUT: output0 [1, 25200, 85]
# ✅ Generated: model_config.go
# ✅ Test passed! Output shape: (1, 25200, 85)
```

---

## 추천 워크플로우

```bash
# 1. 빠른 확인 (5초)
python -c "import onnxruntime as ort; sess = ort.InferenceSession('model.onnx'); print(sess.get_inputs()[0].shape)"

# 2. 상세 분석 필요 시 (10초)
python inspect_model.py

# 3. Go 코드 생성 (15초)
python gen_go_config.py model.onnx

# 4. 추론 테스트
python quick_check.py
```

**결론**: Netron은 시각화용으로만 쓰고, 실제 개발은 **Python 스크립트 + 자동 생성**이 훨씬 효율적입니다!