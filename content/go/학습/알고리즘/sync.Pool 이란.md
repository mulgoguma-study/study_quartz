---
aliases:
  - sync.Pool 이란?
---
# Go의 sync.Pool 완벽 가이드

## sync.Pool이란?

`sync.Pool`은 Go 표준 라이브러리에서 제공하는 **임시 객체 풀(temporary object pool)**입니다. 자주 할당하고 해제되는 객체들을 재사용함으로써 **가비지 컬렉션(GC) 압력을 줄이고 성능을 향상**시킵니다.

결국 `sync.Pool`은 “재활용(객체 풀링)”이고 `map`은 “재활용 도구가 아니라, 그냥 키-값 저장용 자료구조”야.

### sync.Pool은?

- 자주 생성·파괴되는 **임시 객체를 재사용**해서 할당/GC 부담을 줄이려고 만든 구조체.
    
- `Get()`으로 꺼내 쓰고, 사용 후 상태 초기화(`Reset()` 등) 한 뒤 `Put()`으로 다시 넣어두면, 다음에 또 그 객체를 쓸 수 있어.
    
- 런타임/GC가 풀 안의 객체를 언제든 다 버릴 수 있으니까 “있으면 재사용, 없으면 새로 할당”하는 캐시 느낌.
    

그래서 “재활용/재사용”이라는 말이 딱 맞는 기능.


#### map은?

- `map[K]V`는 **연관 배열(해시맵)**일 뿐이고, 목적은 “키로 값을 찾는다”야.
    
- 거기에 포인터나 구조체를 넣어두고 프로그램 로직에서 그걸 다시 꺼내 쓰면, “결과적으로는” 네가 그 객체를 계속 쓰는 거라 일종의 재사용처럼 보일 수는 있어.
    
- 하지만:
    
    - map은 “임시 객체 풀” 같은 정책이 없고
        
    - GC가 map이 들고 있는 값들을 “살아있는 객체”라고 보고 계속 유지해 버림.
        
    - sync.Pool처럼 “GC가 알아서 비워도 되는 캐시”가 아니라, 너가 지우지 않는 한 계속 붙잡고 있는 저장소.
        

그래서 **map = 재활용 메커니즘이 아니라, 그냥 보관소**라고 보는 게 맞다.


##### 한 줄 정리

- “재활용/풀링 기능”을 의미하는 건 **sync.Pool**이고,
    
- **map은 원래 목적이 풀링이 아니라서, 재활용이 필요하면 sync.Pool을 쓰고, map은 캐시/상태 저장에 쓴다** 정도로 정리하면 돼.



## 핵심 개념

### 언제 사용하나?

- 같은 타입의 객체를 **반복적으로 생성/삭제**하는 경우
- 객체 생성 비용이 큰 경우 (예: 큰 버퍼, 복잡한 구조체)
- **성능이 중요한 고빈도 작업**에서

### 주의사항

- Pool에 저장된 객체는 **언제든지 자동으로 제거될 수 있음** (GC 발생 시)
- 영구 저장용이 아닌 **임시 캐시**로만 사용
- Pool에서 꺼낸 객체는 **상태를 초기화**하고 사용해야 함

## 기본 사용법

```go
package main

import (
    "fmt"
    "sync"
)

// 1. Pool 생성
var bufferPool = sync.Pool{
    New: func() interface{} {
        // Pool이 비었을 때 새 객체 생성
        return make([]byte, 1024)
    },
}

func main() {
    // 2. Pool에서 객체 가져오기
    buffer := bufferPool.Get().([]byte)
    
    // 3. 객체 사용
    copy(buffer, []byte("Hello, Pool!"))
    fmt.Println(string(buffer[:12]))
    
    // 4. 사용 후 Pool에 반환 (재사용을 위해)
    bufferPool.Put(buffer)
}
```

## 실전 예제

### 예제 1: HTTP 요청 처리에서 버퍼 재사용

```go
var requestPool = sync.Pool{
    New: func() interface{} {
        return &bytes.Buffer{}
    },
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Pool에서 버퍼 가져오기
    buf := requestPool.Get().(*bytes.Buffer)
    buf.Reset() // 이전 데이터 초기화 (중요!)
    
    // 버퍼 사용
    buf.WriteString("Response: ")
    buf.WriteString(r.URL.Path)
    
    w.Write(buf.Bytes())
    
    // 사용 후 반환
    requestPool.Put(buf)
}
```

### 예제 2: JSON 인코더 재사용

```go
var encoderPool = sync.Pool{
    New: func() interface{} {
        return json.NewEncoder(nil)
    },
}

func encodeJSON(data interface{}) ([]byte, error) {
    buf := &bytes.Buffer{}
    
    encoder := encoderPool.Get().(*json.Encoder)
    encoder.Reset(buf) // 새 writer로 리셋
    
    err := encoder.Encode(data)
    
    encoderPool.Put(encoder)
    
    return buf.Bytes(), err
}
```

### 예제 3: 구조체 재사용

```go
type Worker struct {
    ID      int
    Data    []byte
    Results []string
}

var workerPool = sync.Pool{
    New: func() interface{} {
        return &Worker{
            Data:    make([]byte, 0, 4096),
            Results: make([]string, 0, 100),
        }
    },
}

func processTask(taskID int) {
    worker := workerPool.Get().(*Worker)
    
    // 상태 초기화
    worker.ID = taskID
    worker.Data = worker.Data[:0]       // 슬라이스 길이만 0으로
    worker.Results = worker.Results[:0]
    
    // 작업 수행
    worker.Data = append(worker.Data, "processing"...)
    worker.Results = append(worker.Results, "done")
    
    // 반환
    workerPool.Put(worker)
}
```

## 성능 비교

```go
// Pool 없이 (매번 새로 생성)
func BenchmarkWithoutPool(b *testing.B) {
    for i := 0; i < b.N; i++ {
        buf := make([]byte, 1024)
        _ = buf
    }
}

// Pool 사용
var pool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func BenchmarkWithPool(b *testing.B) {
    for i := 0; i < b.N; i++ {
        buf := pool.Get().([]byte)
        pool.Put(buf)
    }
}
```

**결과**: Pool 사용 시 메모리 할당이 크게 줄어들고, GC 부담이 감소합니다.

## 주요 메서드

| 메서드      | 설명                                      |
| -------- | --------------------------------------- |
| `Get()`  | Pool에서 객체를 가져옴. Pool이 비어있으면 `New` 함수 호출 |
| `Put(x)` | 객체를 Pool에 반환. `nil`은 무시됨                |
| `New`    | Pool이 비었을 때 새 객체를 생성하는 함수 (필드)          |

## 모범 사례

### ✅ 좋은 예

```go
// 1. 타입 단언을 명확하게
buf := pool.Get().(*bytes.Buffer)

// 2. 반드시 상태 초기화
buf.Reset()

// 3. defer로 안전하게 반환
defer pool.Put(buf)

// 4. 큰 객체나 자주 생성되는 객체에 사용
var bigStructPool = sync.Pool{
    New: func() interface{} {
        return &BigStruct{
            Data: make([]byte, 1<<20), // 1MB
        }
    },
}
```

### ❌ 나쁜 예

```go
// 1. 연결이나 파일 핸들 같은 리소스는 Pool에 넣지 않기
// (GC가 언제든 제거할 수 있음)
var connPool = sync.Pool{ // 잘못됨!
    New: func() interface{} {
        return openDBConnection()
    },
}

// 2. Put 없이 Get만 하기 (재사용 안 됨)
obj := pool.Get()
// pool.Put(obj) 호출 안 함 - 메모리 낭비

// 3. 초기화 없이 사용 (이전 데이터가 남아있을 수 있음)
buf := pool.Get().(*bytes.Buffer)
// buf.Reset() 안 함 - 버그 발생 가능
```

## 실제 사용 사례

Go 표준 라이브러리에서도 `sync.Pool`을 활발히 사용합니다:

- `fmt` 패키지: 포맷팅 버퍼 재사용
- `encoding/json`: 인코더/디코더 재사용
- `net/http`: HTTP 요청/응답 처리 시 버퍼 재사용

## 요약

**sync.Pool은 "임시 객체 재활용통"**입니다. 자주 만들었다 버리는 객체가 있다면 Pool에 넣어두고 재사용하세요. 단, Pool에 넣은 건 언제든 사라질 수 있으니 중요한 데이터는 절대 넣지 마세요. 사용 전엔 꼭 초기화하고, 사용 후엔 꼭 반환하세요!