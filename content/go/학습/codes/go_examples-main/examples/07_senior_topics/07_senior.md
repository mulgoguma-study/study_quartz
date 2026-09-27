# Go 시니어 개발자 면접 가이드

시니어 Go 개발자가 알아야 할 심화 주제들을 정리한 문서입니다.

---

## 목차

1. [Slice vs Map 내부 구조](#1-slice-vs-map-내부-구조)
2. [메모리 누수 방지](#2-메모리-누수-방지)
3. [동시성 심화](#3-동시성-심화)
4. [Interface 내부 구조 (iface, eface, itab)](#4-interface-내부-구조)
5. [Channel 내부 구조 (Buffered vs Unbuffered)](#5-channel-내부-구조)
6. [에러 처리 철학과 래핑 전략](#6-에러-처리-철학과-래핑-전략)
7. [map[T]struct{} 사용 이유](#7-maptstruct-사용-이유)
8. [Redis Data Types](#8-redis-data-types)
9. [캐싱 전략 (Look-aside vs Write-through)](#9-캐싱-전략)
10. [gRPC와 Protocol Buffers](#10-grpc와-protocol-buffers)

---

## 1. Slice vs Map 내부 구조

### 면접 질문
> "Slice와 Map의 내부 구조 차이점과 메모리 특성을 설명해주세요"

### Slice 내부 구조 (SliceHeader)

```go
// runtime/slice.go의 실제 구조
type slice struct {
    array unsafe.Pointer  // 8 bytes: 실제 배열을 가리키는 포인터
    len   int             // 8 bytes: 현재 요소 개수
    cap   int             // 8 bytes: 용량 (할당된 공간)
}
// 총 24바이트 헤더
```

**핵심 포인트:**
- Slice는 "값"으로 전달되지만, 내부 array 포인터가 같은 배열을 참조
- 함수에 slice를 전달해도 원본 배열이 수정됨

```go
// 슬라이스가 같은 배열을 공유하는 예시
original := []int{1, 2, 3, 4, 5}
sliced := original[1:3]  // [2, 3]

sliced[0] = 100
fmt.Println(original)  // [1, 100, 3, 4, 5] - 원본도 변경됨!
```

### Slice 용량 증가 전략

**Go 1.18+ 슬라이스 증가 전략:**
- `cap < 256`: 용량을 2배로 증가
- `cap >= 256`: 용량을 1.25배 + 192 증가 (점진적 증가)

```go
var s []int
prevCap := 0

for i := 0; i < 2000; i++ {
    s = append(s, i)
    if cap(s) != prevCap {
        fmt.Printf("len=%4d, cap: %4d → %4d\n", len(s), prevCap, cap(s))
        prevCap = cap(s)
    }
}
```

### Map 내부 구조 (hmap, bmap)

```go
// runtime/map.go의 실제 구조
type hmap struct {
    count     int            // 현재 요소 개수
    flags     uint8          // 상태 플래그
    B         uint8          // 버킷 개수 = 2^B
    noverflow uint16         // overflow 버킷 개수
    hash0     uint32         // 해시 시드 (랜덤)
    buckets   unsafe.Pointer // 버킷 배열 포인터
    oldbuckets unsafe.Pointer // 확장 시 이전 버킷
    nevacuate  uintptr       // 이전 진행 상황
    extra     *mapextra      // 추가 정보
}

// 버킷 구조
type bmap struct {
    tophash [8]uint8  // 각 키의 해시 상위 8비트
    // 이후 8개의 key, 8개의 value가 연속 배치
    // overflow *bmap (암시적)
}
```

**핵심 포인트:**
- 각 버킷은 8개의 key-value 쌍 저장
- tophash로 빠른 비교 (full hash 계산 전에 필터링)
- Load factor 6.5 (평균적으로 버킷당 6.5개 요소)

### 메모리 특성 비교

| 특성 | Slice | Map |
|-----|-------|-----|
| 헤더 크기 | 24 bytes | 8 bytes (포인터) |
| 함수 전달 | 헤더 복사, 배열 공유 | 포인터 복사 |
| nil 상태 | nil slice 가능 | nil map 가능 |
| 요소 주소 획득 | `&s[i]` 가능 | **불가능** (rehash로 이동 가능) |
| 메모리 해제 | nil 할당 또는 재슬라이싱 | delete() 또는 nil 할당 |

### 시니어급 최적화 기법

```go
// 1. 용량 미리 할당
good := make([]int, 0, 10000)  // 재할당 방지

// 2. Map 용량 힌트
m := make(map[string]int, 10000)  // 초기 버킷 할당

// 3. 슬라이스 재사용 ([:0] 트릭)
buffer := make([]byte, 0, 1024)
for i := 0; i < 3; i++ {
    buffer = buffer[:0]  // 길이만 0으로, 용량 유지
    buffer = append(buffer, []byte("reused data")...)
}
```

---

## 2. 메모리 누수 방지

### 면접 질문
> "Go에서 메모리 누수가 발생하는 패턴과 방지법을 설명해주세요"

### 1. Goroutine 누수 (가장 흔한 패턴)

**원인:**
- 채널 수신 대기 중 송신자가 사라짐
- 채널 송신 대기 중 수신자가 사라짐
- 무한 루프에서 탈출 조건 없음
- context 취소 미처리

```go
// 나쁜 예: Goroutine이 영원히 대기
func BadGoroutineLeak() {
    ch := make(chan int)
    go func() {
        val := <-ch  // 영원히 대기 - 누수!
        fmt.Println(val)
    }()
    // ch에 값을 보내지 않고 함수 종료
}

// 좋은 예: context로 취소 가능하게
func GoodGoroutineWithContext(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-ctx.Done():
            fmt.Println("goroutine 정상 종료")
            return
        }
    }()
}
```

### 2. time.Ticker / time.Timer 누수

```go
// 나쁜 예
func BadTickerLeak() {
    ticker := time.NewTicker(1 * time.Second)
    // ticker.Stop() 호출 안 함 - 누수!
    for i := 0; i < 3; i++ {
        <-ticker.C
    }
}

// 좋은 예
func GoodTickerUsage(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()  // 반드시 Stop!

    for {
        select {
        case <-ticker.C:
            fmt.Println("tick")
        case <-ctx.Done():
            return
        }
    }
}
```

**time.After 주의사항:**
```go
// 나쁜 예: 루프에서 time.After - 매번 새 Timer 생성
for {
    select {
    case <-ch:
        // 처리
    case <-time.After(5 * time.Second):  // 매 루프마다 새 Timer!
        // 타임아웃
    }
}

// 좋은 예: Timer 재사용
timer := time.NewTimer(5 * time.Second)
defer timer.Stop()
```

### 3. HTTP Response Body 누수

```go
// 나쁜 예
func BadHTTPUsage() error {
    resp, err := http.Get("https://example.com")
    if err != nil {
        return err
    }
    // resp.Body.Close() 안 함 - 연결 누수!
    return nil
}

// 좋은 예
func GoodHTTPUsage() error {
    resp, err := http.Get("https://example.com")
    if err != nil {
        return err
    }
    defer resp.Body.Close()  // 반드시 Close!

    // Body를 완전히 읽어야 연결 재사용 가능
    _, _ = io.Copy(io.Discard, resp.Body)
    return nil
}
```

### 4. Slice 메모리 누수

```go
// 패턴 1: 큰 배열의 작은 슬라이스
bigData := make([]byte, 1024*1024)  // 1MB

// 나쁜 예: 전체 1MB가 유지됨
smallSlice := bigData[:10]

// 좋은 예: 복사하여 새 슬라이스 생성
smallSlice := make([]byte, 10)
copy(smallSlice, bigData[:10])

// 패턴 2: 포인터 슬라이스 축소
items := make([]*Item, 100)
// 나쁜 예: 뒤 90개 포인터가 여전히 배열에 존재
items = items[:10]

// 좋은 예: nil로 초기화 후 축소
for i := 10; i < len(items); i++ {
    items[i] = nil  // GC가 수집 가능하게 함
}
items = items[:10]
```

### 5. Map 메모리 특성

```go
// Map은 축소되지 않음!
m := make(map[int]int)
for i := 0; i < 100000; i++ {
    m[i] = i
}
for i := 0; i < 99999; i++ {
    delete(m, i)
}
// m의 버킷 메모리는 그대로!

// 해결책: 새 map으로 복사
newMap := make(map[int]int, len(m))
for k, v := range m {
    newMap[k] = v
}
m = newMap
```

### 메모리 누수 방지 체크리스트

- [ ] **Goroutine**: context 또는 done 채널로 종료 가능한가?
- [ ] **Channel**: Unbuffered 채널의 상대방이 항상 존재하는가?
- [ ] **Timer/Ticker**: Stop()을 호출하는가?
- [ ] **HTTP**: Response.Body를 Close하는가?
- [ ] **Slice**: 큰 배열의 작은 슬라이스를 유지하지 않는가?
- [ ] **Map**: 대량 삭제 후 재생성하는가?
- [ ] **defer**: 루프 안 defer를 함수로 분리했는가?

---

## 3. 동시성 심화

### 면접 질문
> "Go의 동시성 모델을 설명하고, 실제 프로덕션에서의 주의점은?"

### G-M-P 스케줄링 모델

```
G (Goroutine): 경량 스레드, 약 2KB 스택으로 시작
M (Machine): OS 스레드, 실제 실행 단위
P (Processor): 논리적 프로세서, GOMAXPROCS 개수만큼 존재
```

**작동 원리:**
1. 각 P는 Local Run Queue (LRQ) 보유
2. G는 P의 LRQ에 추가됨
3. M은 P에 연결되어 G를 실행
4. Work Stealing: 다른 P의 LRQ에서 G를 훔쳐옴

**스케줄링 시점:**
- `go` 키워드로 새 goroutine 생성
- 채널 송/수신 블록
- syscall 호출
- `runtime.Gosched()` 명시적 양보
- GC 실행

### 채널 vs Mutex 선택 기준

| 상황 | 권장 |
|-----|------|
| 데이터 소유권 전달 | 채널 |
| 작업 분배 (Fan-out/Fan-in) | 채널 |
| 이벤트 통지 | 채널 |
| 공유 상태 보호 (캐시, 맵) | Mutex |
| 간단한 카운터 | Atomic |
| 성능이 중요한 상황 | Mutex (채널보다 빠름) |

```go
// Atomic 카운터 (가장 빠름)
type AtomicCounter struct {
    count int64
}

func (c *AtomicCounter) Inc() {
    atomic.AddInt64(&c.count, 1)
}

func (c *AtomicCounter) Value() int64 {
    return atomic.LoadInt64(&c.count)
}
```

### sync.RWMutex - 읽기/쓰기 분리

```go
type ConcurrentCache struct {
    mu   sync.RWMutex
    data map[string]interface{}
}

func (c *ConcurrentCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()  // 읽기 락 - 다수 동시 접근 가능
    defer c.mu.RUnlock()
    v, ok := c.data[key]
    return v, ok
}

func (c *ConcurrentCache) Set(key string, value interface{}) {
    c.mu.Lock()  // 쓰기 락 - 독점적 접근
    defer c.mu.Unlock()
    c.data[key] = value
}
```

### sync.Once - 단 한 번만 실행

```go
var (
    dbInstance *Database
    dbOnce     sync.Once
)

func GetDatabase() *Database {
    dbOnce.Do(func() {
        fmt.Println("데이터베이스 연결 초기화 (한 번만 실행)")
        dbInstance = &Database{conn: "connected"}
    })
    return dbInstance
}
```

### 동시성 패턴: Pipeline

```go
// 파이프라인 구성
gen := func(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

sq := func(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

// 사용
c := gen(1, 2, 3, 4, 5)
c = sq(c)
for result := range c {
    fmt.Println(result)
}
```

### 동시성 패턴: Fan-out/Fan-in

```go
// Fan-out: 하나의 입력을 여러 goroutine으로 분배
// Fan-in: 여러 채널의 결과를 하나로 합침

jobs := make(chan int, 10)
results := make(chan int, 10)

// Fan-out: 3개의 worker
var wg sync.WaitGroup
for w := 1; w <= 3; w++ {
    wg.Add(1)
    go func(workerID int) {
        defer wg.Done()
        for j := range jobs {
            results <- j * 2
        }
    }(w)
}

// Fan-in
go func() {
    wg.Wait()
    close(results)
}()
```

---

## 4. Interface 내부 구조

### 면접 질문
> "Go interface의 내부 구조를 설명해주세요. iface와 eface의 차이는?"

### eface vs iface 구조

```
eface (empty interface = interface{} = any):
┌─────────────────────────────────────┐
│ _type *_type    (8 bytes)           │  → 타입 정보
│ data  unsafe.Pointer (8 bytes)      │  → 실제 값
└─────────────────────────────────────┘
총 16 bytes

iface (non-empty interface):
┌─────────────────────────────────────┐
│ tab  *itab     (8 bytes)            │  → 인터페이스 테이블
│ data unsafe.Pointer (8 bytes)       │  → 실제 값
└─────────────────────────────────────┘
총 16 bytes
```

**핵심 차이:**
- `eface`: 메서드가 없으므로 `_type`만 필요
- `iface`: 메서드 디스패치를 위해 `itab` 필요

### itab (Interface Table) 구조

```go
type itab struct {
    inter *interfacetype  // 인터페이스 타입 정보
    _type *_type          // 구체적 타입 정보
    hash  uint32          // _type.hash 복사 (빠른 타입 assertion)
    _     [4]byte         // 패딩
    fun   [1]uintptr      // 메서드 주소 배열 (가변 크기)
}
```

**itab의 역할:**
1. 구체적 타입과 인터페이스 타입의 연결
2. 메서드 주소 저장 (메서드 디스패치)
3. 동일한 (인터페이스, 타입) 쌍은 전역 캐시에 저장

### nil 인터페이스 주의점

```go
var nilInterface interface{} = nil
var nilPointer *int = nil
var interfaceWithNil interface{} = nilPointer

fmt.Println(nilInterface == nil)      // true
fmt.Println(interfaceWithNil == nil)  // false! (타입 정보가 있음)
```

### interface{} 사용의 비용

1. **Boxing/Unboxing**: 큰 값은 힙 할당 후 포인터 저장
2. **타입 정보 유지**: 모든 interface{}는 `_type` 포인터 필요
3. **런타임 타입 체크**: Type assertion 비용
4. **최적화 방해**: 컴파일러 인라이닝 제한

```go
// 제네릭으로 대체하여 박싱 방지 (Go 1.18+)
func Process[T any](v T) T {
    return v  // 박싱 없음
}
```

---

## 5. Channel 내부 구조

### 면접 질문
> "Buffered와 Unbuffered 채널의 내부 동작 차이를 설명해주세요"

### hchan 구조체

```go
type hchan struct {
    qcount   uint           // 현재 버퍼에 있는 요소 수
    dataqsiz uint           // 버퍼 크기 (make의 두 번째 인자)
    buf      unsafe.Pointer // 버퍼 (Ring Buffer)
    elemsize uint16         // 요소 크기
    closed   uint32         // 채널 닫힘 여부
    elemtype *_type         // 요소 타입
    sendx    uint           // 송신 인덱스
    recvx    uint           // 수신 인덱스
    recvq    waitq          // 수신 대기 goroutine 큐
    sendq    waitq          // 송신 대기 goroutine 큐
    lock     mutex          // 채널 락
}
```

### Unbuffered Channel (dataqsiz = 0)

- 버퍼 없음 (`buf = nil`)
- 송신과 수신이 직접 동기화됨 (Synchronous)
- **직접 전달**: 데이터를 송신자 스택에서 수신자 스택으로 직접 복사

```go
ch := make(chan int)  // unbuffered

go func() {
    val := <-ch  // recvq에 추가되어 대기
    fmt.Println(val)
}()

ch <- 42  // 수신자가 있으므로 직접 전달
```

### Buffered Channel (dataqsiz > 0)

- Ring Buffer로 구현
- 버퍼가 가득 차지 않으면 송신 즉시 완료
- 버퍼가 비어있지 않으면 수신 즉시 완료

```
Ring Buffer 동작 (버퍼 크기 4):
초기:   [_][_][_][_]  sendx=0, recvx=0, qcount=0
송신1:  [A][_][_][_]  sendx=1, recvx=0, qcount=1
송신2:  [A][B][_][_]  sendx=2, recvx=0, qcount=2
수신1:  [_][B][_][_]  sendx=2, recvx=1, qcount=1
송신3:  [_][B][C][_]  sendx=3, recvx=1, qcount=2
송신4:  [_][B][C][D]  sendx=0, recvx=1, qcount=3  (wrap around)
```

### 성능 특성

| 특성 | Unbuffered | Buffered |
|-----|-----------|----------|
| 동기화 | 보장 (Rendezvous) | 비동기 가능 |
| 데이터 전달 | 직접 전달 | 버퍼 복사 |
| goroutine 전환 | 매번 발생 | 버퍼 여유 시 없음 |

### 채널 선택 가이드라인

**Unbuffered 사용:**
- 동기화 지점 필요
- 작업 완료 신호 (done 채널)
- 리소스 접근 제어

**Buffered 사용:**
- 생산자-소비자 속도 차이 완화
- 버스트 트래픽 처리
- Worker Pool
- 처리량 최적화

```go
// 동기화 지점 (Unbuffered)
done := make(chan struct{})
go func() {
    // 작업 수행
    close(done)
}()
<-done

// 세마포어 (Buffered)
sem := make(chan struct{}, maxConcurrent)
sem <- struct{}{}  // 슬롯 획득
<-sem              // 슬롯 반환
```

---

## 6. 에러 처리 철학과 래핑 전략

### 면접 질문
> "왜 Go는 Exception이 없고 Error를 반환하나요?"
> "시니어로서 에러 래핑(Wrapping) 전략은?"

### Go가 Exception 대신 Error를 선택한 이유

1. **명시적 제어 흐름**: 에러 경로가 코드에서 보임
2. **로컬 처리 강제**: 각 호출 지점에서 처리 결정
3. **성능**: Exception의 스택 unwinding 비용 없음
4. **단순성**: try-catch-finally 없음

```go
// Go 스타일 - 명시적
func ProcessFile() error {
    file, err := os.Open("data.txt")
    if err != nil {
        return fmt.Errorf("파일 열기 실패: %w", err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return fmt.Errorf("파일 읽기 실패: %w", err)
    }
    // ...
    return nil
}
```

### panic과 recover

**panic의 올바른 사용:**
- 프로그램이 계속될 수 없는 심각한 오류
- 프로그래머의 실수 (버그)
- 초기화 실패

**panic을 사용하지 말아야 할 때:**
- 예상 가능한 에러 (파일 없음, 네트워크 타임아웃)
- 일반적인 비즈니스 로직 에러

### 에러 래핑 (Go 1.13+)

```go
// %w 동사: 체인 유지
err := fmt.Errorf("GetUserProfile(id=%s): %w", id, originalErr)

// %v 동사: 체인 끊김 (보안상 민감한 정보 숨기기)
err := fmt.Errorf("operation failed: %v", internalErr)
```

### errors.Is vs errors.As

```go
// errors.Is: Sentinel Error (값) 비교
var ErrNotFound = errors.New("not found")

if errors.Is(err, ErrNotFound) {
    // 404 처리
}

// errors.As: Custom Error Type 추출
type ValidationError struct {
    Field   string
    Message string
}

var validErr *ValidationError
if errors.As(err, &validErr) {
    fmt.Printf("검증 실패: %s\n", validErr.Field)
}
```

### 래핑 전략 원칙

1. **컨텍스트는 행동을 설명** (동사로 시작)
   - 좋음: `"querying user: %w"`
   - 나쁨: `"user error: %w"`

2. **중복 피하기**: 새로운 정보가 있을 때만 래핑

3. **공개 API 경계에서 정리**: 내부 구현 세부사항 숨기기

4. **에러는 한 번만 처리**: 로깅하고 반환하면 중복

```go
// 나쁜 예: 과도한 래핑
err = fmt.Errorf("handler error: %w",
    fmt.Errorf("service error: %w",
        fmt.Errorf("repository error: %w", dbErr)))

// 좋은 예: 의미 있는 컨텍스트
err = fmt.Errorf("finding user %s: %w", userID, dbErr)
```

---

## 7. map[T]struct{} 사용 이유

### 면접 질문
> "Go에서 Set을 구현할 때 map[T]struct{}를 사용하는 이유는?"

### struct{} (Empty Struct) 특성

```go
var empty struct{}
fmt.Println(unsafe.Sizeof(empty))  // 0 bytes

var s1, s2, s3 struct{}
fmt.Printf("%p %p %p\n", &s1, &s2, &s3)  // 동일한 주소 (zerobase)
```

### map[T]struct{} vs map[T]bool

```go
// map[string]bool의 문제
boolSet := make(map[string]bool)
boolSet["apple"] = true
boolSet["banana"] = false  // 이게 뭘 의미하지?
// boolSet["cherry"]는 false - 없는 건가?

// map[string]struct{} - 명확함
structSet := make(map[string]struct{})
structSet["apple"] = struct{}{}
_, exists := structSet["apple"]  // 존재 여부만 확인
```

**메모리 차이:**
- `bool`: 각 값에 1 byte (+ padding)
- `struct{}`: 값에 0 byte

### Generic Set 구현 (Go 1.18+)

```go
type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
    return make(Set[T])
}

func (s Set[T]) Add(item T) {
    s[item] = struct{}{}
}

func (s Set[T]) Contains(item T) bool {
    _, exists := s[item]
    return exists
}

func (s Set[T]) Union(other Set[T]) Set[T] {
    result := NewSet[T]()
    for item := range s {
        result.Add(item)
    }
    for item := range other {
        result.Add(item)
    }
    return result
}
```

### 실제 사용 사례

```go
// 1. 중복 제거
func RemoveDuplicates(items []string) []string {
    seen := make(map[string]struct{})
    result := make([]string, 0, len(items))
    for _, item := range items {
        if _, exists := seen[item]; !exists {
            seen[item] = struct{}{}
            result = append(result, item)
        }
    }
    return result
}

// 2. 신호 채널
done := make(chan struct{})
close(done)  // 완료 신호

// 3. 세마포어
type Semaphore chan struct{}
sem := make(Semaphore, maxConcurrent)
sem <- struct{}{}  // Acquire
<-sem              // Release
```

---

## 8. Redis Data Types

### 면접 질문
> "Redis의 주요 데이터 타입과 각각의 활용 사례를 설명해주세요"

### 데이터 타입 요약

| 타입 | 특징 | 활용 사례 |
|-----|------|----------|
| String | 단순 키-값, 최대 512MB | 캐시, 카운터, 분산 락 |
| Hash | 필드-값 쌍 | 사용자 프로필, 세션, 쇼핑 카트 |
| List | 양방향 연결 리스트 | 메시지 큐, 최근 활동 |
| Set | 중복 없는 집합 | 태그, 팔로워, 온라인 사용자 |
| Sorted Set | 점수 기반 정렬 | 리더보드, 타임라인 |

### 1. String

```bash
SET user:123 '{"name":"John"}'
GET user:123
INCR page:views          # Atomic 카운터
SETNX lock:resource 1    # 분산 락
```

**활용:** 단순 캐시, 카운터, 분산 락, Rate Limiting

### 2. Hash

```bash
HSET user:profile:123 name "John" email "john@ex.com"
HGET user:profile:123 name
HGETALL user:profile:123
HINCRBY cart:user123 product456 2  # 수량 증가
```

**활용:** 객체 저장 (부분 업데이트 가능), 세션, 쇼핑 카트

### 3. List

```bash
LPUSH queue:jobs "job1"
RPOP queue:jobs           # FIFO 큐
LPUSH activity:user123 "logged in"
LTRIM activity:user123 0 99  # 최근 100개만 유지
```

**활용:** 메시지 큐, 최근 활동 로그

### 4. Set

```bash
SADD post:123:tags "go" "redis" "backend"
SMEMBERS post:123:tags
SISMEMBER online_users "user123"
SINTER user:123:interests user:456:interests  # 공통 관심사
```

**활용:** 태그, 팔로워/팔로잉, 온라인 사용자

### 5. Sorted Set (가장 강력)

```bash
ZADD leaderboard 100 "player1" 200 "player2"
ZREVRANGE leaderboard 0 9 WITHSCORES  # Top 10
ZINCRBY leaderboard 50 "player1"       # 점수 증가
ZRANK leaderboard "player1"            # 순위 조회
```

**활용:**
- **리더보드**: 점수 = 게임 스코어
- **타임라인**: 점수 = 타임스탬프
- **지연 작업 큐**: 점수 = 실행 시간

```go
// 리더보드 구현
func (l *Leaderboard) UpdateScore(gameID, playerID string, score float64) error {
    return redis.ZAdd(ctx, "leaderboard:"+gameID,
        &redis.Z{Score: score, Member: playerID})
}

func (l *Leaderboard) GetTopPlayers(gameID string, count int64) ([]redis.Z, error) {
    return redis.ZRevRangeWithScores(ctx, "leaderboard:"+gameID, 0, count-1)
}
```

---

## 9. 캐싱 전략

### 면접 질문
> "In-Memory Cache의 위험성은?"
> "Look-aside vs Write-through 전략을 설명해주세요"

### In-Memory Cache 위험성

| 위험 | 설명 | 해결책 |
|-----|------|-------|
| OOM | TTL 없이 무한 증가 | TTL, 최대 크기 제한 |
| 데이터 불일치 | 다중 인스턴스 동기화 불가 | Redis 등 분산 캐시 |
| Cold Start | 재시작 시 캐시 유실 | 캐시 워밍업 |
| GC 부담 | Old Gen 객체 증가 | sync.Pool 활용 |

### Cache-Aside (Look-Aside)

```
읽기: 캐시 조회 → (미스 시) DB 조회 → 캐시 저장 → 반환
쓰기: DB 쓰기 → 캐시 삭제 (무효화)
```

```go
func (ca *CacheAside) Get(ctx context.Context, key string) (interface{}, error) {
    // 1. 캐시에서 조회
    value, err := ca.cache.Get(ctx, key)
    if err == nil {
        return value, nil  // 캐시 히트
    }

    // 2. DB 조회
    value, err = ca.db.Get(ctx, key)
    if err != nil {
        return nil, err
    }

    // 3. 캐시에 저장
    ca.cache.Set(ctx, key, value, ca.ttl)
    return value, nil
}

func (ca *CacheAside) Set(ctx context.Context, key string, value interface{}) error {
    // 1. DB에 쓰기
    if err := ca.db.Set(ctx, key, value); err != nil {
        return err
    }
    // 2. 캐시 무효화 (삭제)
    ca.cache.Delete(ctx, key)
    return nil
}
```

**장점:** 단순, 읽기 최적화, 캐시 실패해도 서비스 가능
**단점:** 캐시 미스 시 DB 부하

### Write-Through

```
읽기: Cache-Aside와 동일
쓰기: 캐시 쓰기 → DB 쓰기 (동기)
```

```go
func (wt *WriteThrough) Set(ctx context.Context, key string, value interface{}) error {
    // 1. 캐시에 쓰기
    if err := wt.cache.Set(ctx, key, value, wt.ttl); err != nil {
        return err
    }
    // 2. DB에 쓰기 (동기)
    if err := wt.db.Set(ctx, key, value); err != nil {
        wt.cache.Delete(ctx, key)  // 롤백
        return err
    }
    return nil
}
```

**장점:** 높은 일관성, 쓰기 후 즉시 캐시 히트
**단점:** 쓰기 지연 증가

### Write-Behind (Write-Back)

```
쓰기: 캐시에만 즉시 쓰기 → 비동기로 DB에 배치 쓰기
```

**장점:** 매우 빠른 쓰기, DB 부하 감소
**단점:** 데이터 손실 위험, 일관성 낮음

### 전략 비교

| 전략 | 읽기 성능 | 쓰기 성능 | 일관성 | 데이터 손실 |
|-----|---------|---------|-------|-----------|
| Cache-Aside | 높음 | 중간 | 낮음 | 낮음 |
| Write-Through | 높음 | 낮음 | 높음 | 낮음 |
| Write-Behind | 높음 | 높음 | 낮음 | 중간 |

---

## 10. gRPC와 Protocol Buffers

### 면접 질문
> "gRPC의 장점/단점과 Protocol Buffers의 이점을 설명해주세요"

### gRPC 장점

1. **고성능**: HTTP/2 멀티플렉싱, 바이너리 직렬화
2. **강타입**: IDL 기반 컴파일 타임 체크
3. **코드 생성**: 클라이언트/서버 스텁 자동 생성
4. **스트리밍**: 4가지 패턴 지원
5. **내장 기능**: 인터셉터, 타임아웃, 취소

### gRPC 단점

1. **브라우저 제한**: grpc-web 필요
2. **디버깅 어려움**: 바이너리 형식, curl 불가
3. **학습 곡선**: protobuf 문법, 빌드 파이프라인
4. **복잡성**: 단순 API에는 과함

### Protocol Buffers 이점

| 비교 항목 | JSON | Protobuf |
|---------|------|----------|
| 메시지 크기 | 1x | 0.3x |
| 직렬화 속도 | 1x | 0.2x |
| 역직렬화 속도 | 1x | 0.1x |
| 스키마 | 없음 | 명확함 |

### .proto 파일 예시

```protobuf
syntax = "proto3";

message User {
    string id = 1;
    string name = 2;
    string email = 3;
    repeated string roles = 4;
}

service UserService {
    // Unary RPC
    rpc GetUser(GetUserRequest) returns (User);

    // Server Streaming
    rpc ListUsers(ListUsersRequest) returns (stream User);

    // Client Streaming
    rpc CreateUsers(stream User) returns (CreateUsersResponse);

    // Bidirectional Streaming
    rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}
```

### gRPC vs REST 선택 기준

**gRPC 선택:**
- 마이크로서비스 내부 통신
- 스트리밍 필요
- 성능 중요
- 다국어 환경

**REST 선택:**
- 공개 API
- 브라우저 직접 호출
- 단순 CRUD
- HTTP 캐싱 필요

---

## 시니어 면접 요약 체크리스트

### Go 내부 구조
- [ ] Slice 헤더 (24 bytes: ptr, len, cap)
- [ ] Map 버킷 구조 (8개 key-value)
- [ ] iface vs eface 차이
- [ ] hchan 구조와 Ring Buffer

### 동시성
- [ ] G-M-P 스케줄러 설명
- [ ] 채널 vs Mutex 선택 기준
- [ ] Goroutine 누수 방지

### 에러 처리
- [ ] Exception vs Error 철학
- [ ] %w vs %v 래핑 차이
- [ ] errors.Is vs errors.As

### 시스템 설계
- [ ] Redis 데이터 타입별 활용
- [ ] Cache-Aside vs Write-Through
- [ ] gRPC vs REST 선택 기준
