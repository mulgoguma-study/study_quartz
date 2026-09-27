---
aliases:
  - python -> go onnx 자동화
---
## 자동 Go 코드 생성 스크립트

### 최고의 방법: Python → Go config 자동 생성

python

```python
#!/usr/bin/env python3
"""
Advanced ONNX to Go Code Generator
실무에 필요한 모든 정보를 추출하여 Go 코드 자동 생성
"""

import onnx
import onnxruntime as ort
import numpy as np
import sys
import json
from pathlib import Path

class ONNXAnalyzer:
    def __init__(self, model_path):
        self.model_path = model_path
        self.model = onnx.load(model_path)
        self.session = ort.InferenceSession(model_path)
        
    def get_dtype_mapping(self, onnx_type_str):
        """ONNX dtype을 Go/C 타입으로 변환"""
        # ONNX Runtime의 타입 문자열 처리
        type_mapping = {
            'tensor(float)': (1, "float32", "C.ONNX_TENSOR_ELEMENT_DATA_TYPE_FLOAT", "[]float32"),
            'tensor(uint8)': (2, "uint8", "C.ONNX_TENSOR_ELEMENT_DATA_TYPE_UINT8", "[]uint8"),
            'tensor(int8)': (3, "int8", "C.ONNX_TENSOR_ELEMENT_DATA_TYPE_INT8", "[]int8"),
            'tensor(int32)': (6, "int32", "C.ONNX_TENSOR_ELEMENT_DATA_TYPE_INT32", "[]int32"),
            'tensor(int64)': (7, "int64", "C.ONNX_TENSOR_ELEMENT_DATA_TYPE_INT64", "[]int64"),
            'tensor(double)': (11, "float64", "C.ONNX_TENSOR_ELEMENT_DATA_TYPE_DOUBLE", "[]float64"),
        }
        result = type_mapping.get(onnx_type_str, (0, "unknown", "unknown", "[]byte"))
        return result[1], result[2], result[3]  # dtype_name, dtype_const, go_slice
    
    def analyze_inputs(self):
        """모든 입력 분석"""
        inputs = []
        for inp in self.session.get_inputs():
            shape = inp.shape
            # 동적 차원 처리
            processed_shape = []
            has_dynamic = False
            for dim in shape:
                if dim is None or isinstance(dim, str) or dim == -1:
                    processed_shape.append(-1)
                    has_dynamic = True
                else:
                    processed_shape.append(int(dim))
            
            dtype_name, dtype_const, go_slice = self.get_dtype_mapping(inp.type)
            
            # 총 크기 계산 (동적 차원은 1로 가정)
            total_size = 1
            for dim in processed_shape:
                if dim > 0:
                    total_size *= dim
            
            inputs.append({
                'name': inp.name,
                'shape': processed_shape,
                'dtype': dtype_name,
                'dtype_const': dtype_const,
                'go_type': go_slice,
                'total_size': total_size,
                'has_dynamic': has_dynamic,
            })
        return inputs
    
    def analyze_outputs(self):
        """모든 출력 분석"""
        outputs = []
        for out in self.session.get_outputs():
            shape = out.shape
            processed_shape = []
            has_dynamic = False
            for dim in shape:
                if dim is None or isinstance(dim, str) or dim == -1:
                    processed_shape.append(-1)
                    has_dynamic = True
                else:
                    processed_shape.append(int(dim))
            
            dtype_name, dtype_const, go_slice = self.get_dtype_mapping(out.type)
            
            total_size = 1
            for dim in processed_shape:
                if dim > 0:
                    total_size *= dim
            
            outputs.append({
                'name': out.name,
                'shape': processed_shape,
                'dtype': dtype_name,
                'dtype_const': dtype_const,
                'go_type': go_slice,
                'total_size': total_size,
                'has_dynamic': has_dynamic,
            })
        return outputs
    
    def analyze_metadata(self):
        """모델 메타데이터 분석"""
        meta = self.model.metadata_props
        metadata = {}
        for prop in meta:
            metadata[prop.key] = prop.value
        
        # 모델 통계
        num_nodes = len(self.model.graph.node)
        
        # Opset 버전
        opset_version = self.model.opset_import[0].version if self.model.opset_import else 0
        
        return {
            'producer_name': self.model.producer_name,
            'producer_version': self.model.producer_version,
            'model_version': self.model.model_version,
            'ir_version': self.model.ir_version,
            'opset_version': opset_version,
            'num_nodes': num_nodes,
            'custom_metadata': metadata,
        }
    
    def test_inference(self):
        """실제 추론 테스트로 정확한 shape 확인"""
        try:
            # 더미 입력 생성
            dummy_inputs = {}
            for inp in self.session.get_inputs():
                shape = []
                for dim in inp.shape:
                    if dim is None or isinstance(dim, str) or dim == -1:
                        shape.append(1)  # 동적 차원은 1로
                    else:
                        shape.append(int(dim))
                
                if 'float' in inp.type:
                    dummy_inputs[inp.name] = np.random.randn(*shape).astype(np.float32)
                elif 'int64' in inp.type:
                    dummy_inputs[inp.name] = np.random.randint(0, 100, shape).astype(np.int64)
                else:
                    dummy_inputs[inp.name] = np.random.randn(*shape).astype(np.float32)
            
            # 추론 실행
            outputs = self.session.run(None, dummy_inputs)
            
            # 실제 출력 shape 기록
            actual_shapes = [out.shape for out in outputs]
            return True, actual_shapes
        except Exception as e:
            return False, str(e)
    
    def detect_model_type(self):
        """모델 타입 추측"""
        inputs = self.analyze_inputs()
        outputs = self.analyze_outputs()
        
        inp_shape = inputs[0]['shape']
        out_shape = outputs[0]['shape']
        
        # 이미지 분류
        if len(inp_shape) == 4 and inp_shape[1] == 3 and len(out_shape) == 2:
            return "image_classification"
        
        # 객체 탐지
        if len(inp_shape) == 4 and len(out_shape) == 3 and out_shape[-1] > 80:
            return "object_detection"
        
        # 텍스트/임베딩
        if 'int64' in inputs[0]['dtype']:
            return "text_model"
        
        # 세그멘테이션
        if len(out_shape) == 4 and out_shape[1] > 1:
            return "segmentation"
        
        return "unknown"

def generate_go_code(analyzer):
    """완전한 Go 코드 생성"""
    inputs = analyzer.analyze_inputs()
    outputs = analyzer.analyze_outputs()
    metadata = analyzer.analyze_metadata()
    model_type = analyzer.detect_model_type()
    test_success, test_result = analyzer.test_inference()
    
    model_name = Path(analyzer.model_path).stem.replace('-', '_').replace('.', '_')
    
    # Go 코드 생성
    from datetime import datetime
    code = f'''package config

/*
Auto-generated from ONNX model: {analyzer.model_path}
Generated on: {datetime.now().strftime("%Y-%m-%d %H:%M:%S")}

Model Information:
- Type: {model_type}
- Producer: {metadata['producer_name']} {metadata['producer_version']}
- Opset Version: {metadata['opset_version']}
- Total Nodes: {metadata['num_nodes']}
- Test Inference: {'✅ SUCCESS' if test_success else '❌ FAILED'}
*/

/*
#cgo LDFLAGS: -lonnxruntime
#include <onnxruntime_c_api.h>
*/
import "C"
import "fmt"  // fmt import 추가

// {model_name.upper()}_CONFIG
type {model_name.title().replace('_', '')}Config struct {{
    ModelPath string
    
    // Input Configuration
'''
    
    # 입력 정보 생성
    for i, inp in enumerate(inputs):
        code += f'''    Input{i}Name    string
    Input{i}Shape   []int64
    Input{i}Type    C.enum_ONNXTensorElementDataType
    Input{i}Size    int  // Total elements: {inp['total_size']}
'''
        if inp['has_dynamic']:
            code += f'''    Input{i}HasDynamic bool  // ⚠️  Dynamic shape detected
'''
    
    # 출력 정보 생성
    code += '''    
    // Output Configuration
'''
    for i, out in enumerate(outputs):
        code += f'''    Output{i}Name   string
    Output{i}Shape  []int64
    Output{i}Type   C.enum_ONNXTensorElementDataType
    Output{i}Size   int  // Total elements: {out['total_size']}
'''
        if out['has_dynamic']:
            code += f'''    Output{i}HasDynamic bool  // ⚠️  Dynamic shape detected
'''
    
    # 전처리 힌트
    code += f'''    
    // Preprocessing hints (model type: {model_type})
'''
    
    if model_type == "image_classification":
        code += '''    RequiresNormalization bool
    Mean []float32  // e.g., [0.485, 0.456, 0.406] for ImageNet
    Std  []float32  // e.g., [0.229, 0.224, 0.225] for ImageNet
    RequiresBGR bool
    RequiresCHW bool  // Channel-first format
'''
    elif model_type == "object_detection":
        code += '''    ConfidenceThreshold float32
    NMSThreshold float32
    MaxDetections int
    RequiresBGR bool
'''
    elif model_type == "text_model":
        code += '''    MaxSequenceLength int
    VocabSize int
    PaddingTokenID int64
'''
    
    code += '''}

// NewDefaultConfig creates a default configuration
'''
    code += f'''func NewDefault{model_name.title().replace('_', '')}Config() *{model_name.title().replace('_', '')}Config {{
    return &{model_name.title().replace('_', '')}Config{{
        ModelPath: "{analyzer.model_path}",
'''
    
    # 입력 기본값
    for i, inp in enumerate(inputs):
        code += f'''        
        Input{i}Name:  "{inp['name']}",
        Input{i}Shape: []int64{{{', '.join(map(str, inp['shape']))}}},
        Input{i}Type:  {inp['dtype_const']},
        Input{i}Size:  {inp['total_size']},
'''
        if inp['has_dynamic']:
            code += f'''        Input{i}HasDynamic: true,
'''
    
    # 출력 기본값
    for i, out in enumerate(outputs):
        code += f'''        
        Output{i}Name:  "{out['name']}",
        Output{i}Shape: []int64{{{', '.join(map(str, out['shape']))}}},
        Output{i}Type:  {out['dtype_const']},
        Output{i}Size:  {out['total_size']},
'''
        if out['has_dynamic']:
            code += f'''        Output{i}HasDynamic: true,
'''
    
    # 전처리 기본값
    if model_type == "image_classification":
        code += '''        
        RequiresNormalization: true,
        Mean: []float32{0.485, 0.456, 0.406},  // ImageNet default
        Std:  []float32{0.229, 0.224, 0.225},
        RequiresBGR: false,
        RequiresCHW: true,
'''
    elif model_type == "object_detection":
        code += '''        
        ConfidenceThreshold: 0.25,
        NMSThreshold: 0.45,
        MaxDetections: 100,
        RequiresBGR: false,
'''
    elif model_type == "text_model":
        inp_shape = inputs[0]['shape']
        code += f'''        
        MaxSequenceLength: {inp_shape[1] if len(inp_shape) > 1 else 512},
        VocabSize: 30000,  // TODO: verify
        PaddingTokenID: 0,
'''
    
    code += '''    }
}

// Validate checks if the configuration is valid
'''
    code += f'''func (c *{model_name.title().replace('_', '')}Config) Validate() error {{
    if c.ModelPath == "" {{
        return fmt.Errorf("model path is empty")
    }}
'''
    
    # 동적 shape 검증
    for i, inp in enumerate(inputs):
        if inp['has_dynamic']:
            code += f'''    
    if c.Input{i}HasDynamic {{
        // ⚠️  Dynamic shape - validation skipped
    }}
'''
    
    code += '''    
    return nil
}
'''
    
    return code

def generate_test_code(analyzer):
    """테스트 코드 생성"""
    inputs = analyzer.analyze_inputs()
    model_name = Path(analyzer.model_path).stem.replace('-', '_').replace('.', '_')
    
    test_code = f'''package config

import (
    "testing"
)

func TestDefaultConfig(t *testing.T) {{
    config := NewDefault{model_name.title().replace('_', '')}Config()
    
    if err := config.Validate(); err != nil {{
        t.Fatalf("Config validation failed: %v", err)
    }}
    
    // Verify input shapes
'''
    
    for i, inp in enumerate(inputs):
        test_code += f'''    expectedInput{i}Size := {inp['total_size']}
    if config.Input{i}Size != expectedInput{i}Size {{
        t.Errorf("Input{i} size mismatch: expected %d, got %d", 
                 expectedInput{i}Size, config.Input{i}Size)
    }}
'''
    
    test_code += '''}
'''
    
    return test_code

def generate_json_schema(analyzer):
    """JSON 스키마 생성 (설정 파일용)"""
    inputs = analyzer.analyze_inputs()
    outputs = analyzer.analyze_outputs()
    metadata = analyzer.analyze_metadata()
    
    schema = {
        "model_info": {
            "path": analyzer.model_path,
            "type": analyzer.detect_model_type(),
            "producer": metadata['producer_name'],
            "opset_version": metadata['opset_version'],
        },
        "inputs": [],
        "outputs": [],
    }
    
    for inp in inputs:
        schema["inputs"].append({
            "name": inp['name'],
            "shape": inp['shape'],
            "dtype": inp['dtype'],
            "size": inp['total_size'],
            "has_dynamic": inp['has_dynamic'],
        })
    
    for out in outputs:
        schema["outputs"].append({
            "name": out['name'],
            "shape": out['shape'],
            "dtype": out['dtype'],
            "size": out['total_size'],
            "has_dynamic": out['has_dynamic'],
        })
    
    return json.dumps(schema, indent=2)

def main():
    if len(sys.argv) < 2:
        print("Usage: python advanced_onnx_to_go.py <model.onnx>")
        sys.exit(1)
    
    model_path = sys.argv[1]
    
    print(f"🔍 Analyzing ONNX model: {model_path}")
    print("=" * 60)
    
    analyzer = ONNXAnalyzer(model_path)
    
    # 분석 정보 출력
    inputs = analyzer.analyze_inputs()
    outputs = analyzer.analyze_outputs()
    metadata = analyzer.analyze_metadata()
    model_type = analyzer.detect_model_type()
    
    print(f"\n📊 Model Type: {model_type}")
    print(f"📦 Producer: {metadata['producer_name']} {metadata['producer_version']}")
    print(f"🔢 Total Nodes: {metadata['num_nodes']}")
    
    print(f"\n📥 INPUTS ({len(inputs)}):")
    for i, inp in enumerate(inputs):
        dynamic = " (DYNAMIC)" if inp['has_dynamic'] else ""
        print(f"  [{i}] {inp['name']}: {inp['dtype']}{inp['shape']}{dynamic}")
        print(f"      → Total size: {inp['total_size']:,} elements")
    
    print(f"\n📤 OUTPUTS ({len(outputs)}):")
    for i, out in enumerate(outputs):
        dynamic = " (DYNAMIC)" if out['has_dynamic'] else ""
        print(f"  [{i}] {out['name']}: {out['dtype']}{out['shape']}{dynamic}")
        print(f"      → Total size: {out['total_size']:,} elements")
    
    # 추론 테스트
    print(f"\n🧪 Testing inference...")
    test_success, test_result = analyzer.test_inference()
    if test_success:
        print(f"   ✅ SUCCESS - Actual output shapes: {test_result}")
    else:
        print(f"   ❌ FAILED - {test_result}")
    
    # 파일 생성
    model_name = Path(model_path).stem.replace('-', '_')
    
    print(f"\n📝 Generating Go code...")
    
    # 1. Config 파일
    config_code = generate_go_code(analyzer)
    config_file = f"{model_name}_config.go"
    with open(config_file, 'w') as f:
        f.write(config_code)
    print(f"   ✅ {config_file}")
    
    # 2. Test 파일
    test_code = generate_test_code(analyzer)
    test_file = f"{model_name}_config_test.go"
    with open(test_file, 'w') as f:
        f.write(test_code)
    print(f"   ✅ {test_file}")
    
    # 3. JSON 스키마
    json_schema = generate_json_schema(analyzer)
    json_file = f"{model_name}_schema.json"
    with open(json_file, 'w') as f:
        f.write(json_schema)
    print(f"   ✅ {json_file}")
    
    print(f"\n🎉 Done! Generated {3} files.")
    print(f"\nNext steps:")
    print(f"  1. Review {config_file}")
    print(f"  2. Adjust preprocessing parameters if needed")
    print(f"  3. Run: go test -v")

if __name__ == "__main__":
    main()
```

**사용법:**

bash

```bash
python gen_go_config.py yolov8.onnx

# 자동으로 model_config.go 생성됨!
```