package senior_topics

/*
================================================================================
map[T]struct{} 사용 이유 (시니어 레벨)
================================================================================

면접 질문: "Go에서 Set을 구현할 때 map[T]struct{}를 사용하는 이유는?"

시니어급 답변 포인트:
1. struct{} 의 메모리 특성
2. Set 구현 패턴
3. 성능 비교 (struct{} vs bool)
4. 실제 사용 사례
*/

import (
	"fmt"
	"unsafe"
)

// ============================================================================
// 1. struct{} (Empty Struct) 특성
// ============================================================================

/*
struct{}의 특성:

1. 크기가 0 bytes:
   - 어떤 필드도 없음
   - 메모리 할당 없음
   - sizeof(struct{}) == 0

2. 모든 struct{}는 동일한 주소:
   - runtime.zerobase 가리킴
   - 최적화된 메모리 사용

3. 용도:
   - Set 구현의 값으로 사용
   - 신호 채널 (done channel)
   - 메서드만 가진 타입
*/

func EmptyStructDemo() {
	fmt.Println("=== struct{} 특성 데모 ===")

	// 1. 크기 확인
	var empty struct{}
	fmt.Printf("struct{} 크기: %d bytes\n", unsafe.Sizeof(empty))

	// 2. bool과 비교
	var b bool
	fmt.Printf("bool 크기: %d bytes\n", unsafe.Sizeof(b))

	// 3. 여러 struct{}의 주소
	var s1, s2, s3 struct{}
	fmt.Printf("s1 주소: %p\n", &s1)
	fmt.Printf("s2 주소: %p\n", &s2)
	fmt.Printf("s3 주소: %p\n", &s3)
	fmt.Println("(동일한 주소 - zerobase)")

	// 4. 슬라이스에서의 동작
	slice := make([]struct{}, 1000000)
	fmt.Printf("struct{} 100만개 슬라이스 크기: %d bytes\n", unsafe.Sizeof(slice))
	// 슬라이스 헤더(24 bytes)만 사용, 요소는 0 bytes
	_ = slice
}

// ============================================================================
// 2. Set 구현: map[T]struct{} vs map[T]bool
// ============================================================================

/*
왜 struct{}를 사용하는가?

map[string]bool:
- 각 엔트리에 bool 값 저장 (1 byte + padding = 8 bytes)
- 값의 의미가 모호함 (false면 없는 건가?)

map[string]struct{}:
- 값에 0 bytes 사용
- 존재 여부만 표현 (키가 있으면 존재)
- 의도가 명확함 (Set임을 나타냄)

메모리 차이:
- 1000개 요소 기준
- bool: 1000 * 8 = 8000 bytes (값 부분)
- struct{}: 0 bytes (값 부분)
*/

// map[T]bool Set
type BoolSet map[string]bool

func (s BoolSet) Add(item string) {
	s[item] = true
}

func (s BoolSet) Contains(item string) bool {
	return s[item] // false면 없음? 아니면 false가 저장된 건가?
}

func (s BoolSet) Remove(item string) {
	delete(s, item)
}

// map[T]struct{} Set (권장)
type Set map[string]struct{}

func NewSet() Set {
	return make(Set)
}

func (s Set) Add(item string) {
	s[item] = struct{}{}
}

func (s Set) Contains(item string) bool {
	_, exists := s[item]
	return exists
}

func (s Set) Remove(item string) {
	delete(s, item)
}

func (s Set) Size() int {
	return len(s)
}

func (s Set) Items() []string {
	items := make([]string, 0, len(s))
	for item := range s {
		items = append(items, item)
	}
	return items
}

func SetComparisonDemo() {
	fmt.Println("\n=== Set 구현 비교 ===")

	// bool Set의 문제
	boolSet := make(BoolSet)
	boolSet["apple"] = true
	boolSet["banana"] = false // 이게 뭘 의미하지?

	fmt.Printf("boolSet[\"apple\"]: %v\n", boolSet.Contains("apple"))
	fmt.Printf("boolSet[\"banana\"]: %v\n", boolSet.Contains("banana")) // false - 없는 건가?
	fmt.Printf("boolSet[\"cherry\"]: %v\n", boolSet.Contains("cherry")) // false - 없는 거

	// struct{} Set - 명확함
	structSet := NewSet()
	structSet.Add("apple")
	structSet.Add("banana")

	fmt.Printf("\nstructSet.Contains(\"apple\"): %v\n", structSet.Contains("apple"))
	fmt.Printf("structSet.Contains(\"banana\"): %v\n", structSet.Contains("banana"))
	fmt.Printf("structSet.Contains(\"cherry\"): %v\n", structSet.Contains("cherry"))
}

// ============================================================================
// 3. Generic Set (Go 1.18+)
// ============================================================================

type GenericSet[T comparable] map[T]struct{}

func NewGenericSet[T comparable]() GenericSet[T] {
	return make(GenericSet[T])
}

func (s GenericSet[T]) Add(item T) {
	s[item] = struct{}{}
}

func (s GenericSet[T]) Contains(item T) bool {
	_, exists := s[item]
	return exists
}

func (s GenericSet[T]) Remove(item T) {
	delete(s, item)
}

func (s GenericSet[T]) Union(other GenericSet[T]) GenericSet[T] {
	result := NewGenericSet[T]()
	for item := range s {
		result.Add(item)
	}
	for item := range other {
		result.Add(item)
	}
	return result
}

func (s GenericSet[T]) Intersection(other GenericSet[T]) GenericSet[T] {
	result := NewGenericSet[T]()
	for item := range s {
		if other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

func (s GenericSet[T]) Difference(other GenericSet[T]) GenericSet[T] {
	result := NewGenericSet[T]()
	for item := range s {
		if !other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

func GenericSetDemo() {
	fmt.Println("\n=== Generic Set 데모 ===")

	// 정수 Set
	intSet := NewGenericSet[int]()
	intSet.Add(1)
	intSet.Add(2)
	intSet.Add(3)
	fmt.Printf("intSet: %v\n", intSet)

	// 집합 연산
	setA := NewGenericSet[int]()
	setA.Add(1)
	setA.Add(2)
	setA.Add(3)

	setB := NewGenericSet[int]()
	setB.Add(2)
	setB.Add(3)
	setB.Add(4)

	fmt.Printf("A: %v\n", setA)
	fmt.Printf("B: %v\n", setB)
	fmt.Printf("A ∪ B (합집합): %v\n", setA.Union(setB))
	fmt.Printf("A ∩ B (교집합): %v\n", setA.Intersection(setB))
	fmt.Printf("A - B (차집합): %v\n", setA.Difference(setB))
}

// ============================================================================
// 4. 실제 사용 사례
// ============================================================================

/*
struct{} 사용 사례:

1. Set 구현
2. 신호 채널 (done, quit)
3. 방문 표시 (visited map)
4. 세마포어
5. 메서드만 가진 타입
*/

// 사례 1: 중복 제거
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

// 사례 2: 신호 채널
func SignalChannelExample() {
	done := make(chan struct{})

	go func() {
		// 작업 수행
		close(done) // 완료 신호 (값 전송 없이)
	}()

	<-done // 완료 대기
}

// 사례 3: 그래프 방문 표시
func BFSWithVisited(graph map[int][]int, start int) []int {
	visited := make(map[int]struct{})
	queue := []int{start}
	result := []int{}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if _, exists := visited[node]; exists {
			continue
		}
		visited[node] = struct{}{}
		result = append(result, node)

		for _, neighbor := range graph[node] {
			if _, exists := visited[neighbor]; !exists {
				queue = append(queue, neighbor)
			}
		}
	}
	return result
}

// 사례 4: 세마포어
type Semaphore chan struct{}

func NewSemaphore(max int) Semaphore {
	return make(chan struct{}, max)
}

func (s Semaphore) Acquire() {
	s <- struct{}{} // 슬롯 획득 (값은 의미 없음)
}

func (s Semaphore) Release() {
	<-s // 슬롯 반환
}

// 사례 5: 메서드만 가진 타입
type Logger struct{} // 상태 없음

func (Logger) Info(msg string) {
	fmt.Printf("[INFO] %s\n", msg)
}

func (Logger) Error(msg string) {
	fmt.Printf("[ERROR] %s\n", msg)
}

func UseCaseDemo() {
	fmt.Println("\n=== 실제 사용 사례 데모 ===")

	// 중복 제거
	items := []string{"apple", "banana", "apple", "cherry", "banana"}
	unique := RemoveDuplicates(items)
	fmt.Printf("중복 제거: %v\n", unique)

	// 그래프 BFS
	graph := map[int][]int{
		1: {2, 3},
		2: {4, 5},
		3: {5},
		4: {},
		5: {},
	}
	traversal := BFSWithVisited(graph, 1)
	fmt.Printf("BFS 순회: %v\n", traversal)

	// Logger
	logger := Logger{}
	logger.Info("Application started")
}

// ============================================================================
// 5. 메모리 비교 벤치마크
// ============================================================================

func MemoryComparisonDemo() {
	fmt.Println("\n=== 메모리 비교 ===")

	const count = 100000

	// map[string]bool
	boolMap := make(map[string]bool, count)
	for i := 0; i < count; i++ {
		boolMap[fmt.Sprintf("key%d", i)] = true
	}

	// map[string]struct{}
	structMap := make(map[string]struct{}, count)
	for i := 0; i < count; i++ {
		structMap[fmt.Sprintf("key%d", i)] = struct{}{}
	}

	fmt.Printf("10만개 요소 기준:\n")
	fmt.Printf("  map[string]bool - 값당 1 byte (+ padding)\n")
	fmt.Printf("  map[string]struct{} - 값당 0 byte\n")
	fmt.Println("  (실제 차이는 map 버킷 구조에 따라 다름)")
}

// ============================================================================
// 6. 시니어 면접 답변 요약
// ============================================================================

/*
Q: Go에서 Set을 구현할 때 map[T]struct{}를 사용하는 이유는?

A: struct{}는 크기가 0 bytes인 특별한 타입으로, Set 구현에 최적입니다.

1. struct{}의 특성:
   - sizeof(struct{}) == 0
   - 모든 struct{} 인스턴스는 같은 주소 (runtime.zerobase)
   - 메모리 할당 없음

2. map[T]struct{} vs map[T]bool:
   - bool: 각 값에 1 byte (+ padding)
   - struct{}: 값에 0 byte
   - bool은 의미가 모호함 (false가 뭘 의미?)
   - struct{}는 의도가 명확함 (존재 여부만)

3. 코드 의도 표현:
   - map[T]struct{}를 보면 Set임을 즉시 알 수 있음
   - Go 커뮤니티의 관용적 표현

4. 성능:
   - 메모리 절약 (대량의 요소 시 유의미)
   - 삽입/조회 성능은 동일

5. 다른 활용:
   - 신호 채널: chan struct{}
   - 방문 표시: map[node]struct{}
   - 세마포어: chan struct{}
   - 상태 없는 타입: type Logger struct{}

사용 예:
set := make(map[string]struct{})
set["key"] = struct{}{}  // 추가
_, exists := set["key"]  // 존재 확인
delete(set, "key")       // 삭제
*/
