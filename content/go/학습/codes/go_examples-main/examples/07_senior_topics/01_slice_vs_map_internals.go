package senior_topics

/*
================================================================================
시니어 Go 개발자를 위한 Slice vs Map 내부 구조 심화
================================================================================

면접 질문: "Slice와 Map의 내부 구조 차이점과 메모리 특성을 설명해주세요"

시니어급 답변 포인트:
1. 내부 데이터 구조 (헤더, 버킷 등)
2. 메모리 할당 및 증가 전략
3. GC와의 상호작용
4. 성능 특성 및 최적화 기법
*/

import (
	"fmt"
	"reflect"
	"unsafe"
)

// ============================================================================
// 1. Slice 내부 구조 (SliceHeader)
// ============================================================================

/*
Slice의 실제 구조 (runtime/slice.go):

type slice struct {
    array unsafe.Pointer  // 8 bytes: 실제 배열을 가리키는 포인터
    len   int             // 8 bytes: 현재 요소 개수
    cap   int             // 8 bytes: 용량 (할당된 공간)
}

총 24바이트의 헤더 구조체
- Slice는 "값"으로 전달되지만, 내부 array 포인터가 같은 배열을 참조
- 따라서 slice를 함수에 전달해도 원본 배열이 수정됨
*/

func SliceInternalsDemo() {
	fmt.Println("=== Slice 내부 구조 분석 ===")

	// 1. Slice 헤더 크기 확인
	s := make([]int, 5, 10)
	fmt.Printf("Slice 헤더 크기: %d bytes\n", unsafe.Sizeof(s))

	// 2. reflect.SliceHeader로 내부 구조 확인
	header := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Printf("Data 포인터: %v\n", header.Data)
	fmt.Printf("Len: %d\n", header.Len)
	fmt.Printf("Cap: %d\n", header.Cap)

	// 3. Slice가 같은 배열을 공유하는 예시
	original := []int{1, 2, 3, 4, 5}
	sliced := original[1:3] // [2, 3]

	originalHeader := (*reflect.SliceHeader)(unsafe.Pointer(&original))
	slicedHeader := (*reflect.SliceHeader)(unsafe.Pointer(&sliced))

	fmt.Printf("\n원본 배열 시작: %v\n", originalHeader.Data)
	fmt.Printf("슬라이스 배열 시작: %v\n", slicedHeader.Data)
	fmt.Printf("포인터 차이 (int 1개 크기): %d bytes\n",
		slicedHeader.Data-originalHeader.Data)

	// 4. 슬라이스 수정이 원본에 영향
	sliced[0] = 100
	fmt.Printf("sliced[0]=100 후 original: %v\n", original) // [1, 100, 3, 4, 5]
}

// ============================================================================
// 2. Slice 용량 증가 전략 (Growth Strategy)
// ============================================================================

/*
Go 1.18+ 슬라이스 증가 전략 (runtime/slice.go):

1. cap < 256: 용량을 2배로 증가
2. cap >= 256: 용량을 1.25배 + 192 증가 (점진적 증가)

이전 버전(Go 1.17 이하):
1. cap < 1024: 2배 증가
2. cap >= 1024: 1.25배 증가
*/

func SliceGrowthDemo() {
	fmt.Println("\n=== Slice 용량 증가 패턴 ===")

	var s []int
	prevCap := 0

	for i := 0; i < 2000; i++ {
		s = append(s, i)
		if cap(s) != prevCap {
			growthRate := float64(cap(s)) / float64(max(prevCap, 1))
			fmt.Printf("len=%4d, cap 변화: %4d → %4d (증가율: %.2fx)\n",
				len(s), prevCap, cap(s), growthRate)
			prevCap = cap(s)
		}
	}
}

// ============================================================================
// 3. Map 내부 구조 (hmap, bmap)
// ============================================================================

/*
Map의 실제 구조 (runtime/map.go):

type hmap struct {
    count     int            // 현재 요소 개수
    flags     uint8          // 상태 플래그
    B         uint8          // 버킷 개수 = 2^B
    noverflow uint16         // overflow 버킷 개수 (근사치)
    hash0     uint32         // 해시 시드 (랜덤)
    buckets   unsafe.Pointer // 버킷 배열 포인터
    oldbuckets unsafe.Pointer // 확장 시 이전 버킷
    nevacuate  uintptr       // 이전 진행 상황
    extra     *mapextra      // 추가 정보
}

버킷 구조 (bmap):
type bmap struct {
    tophash [8]uint8  // 각 키의 해시 상위 8비트
    // 이후 8개의 key, 8개의 value가 연속 배치
    // overflow *bmap (암시적)
}

핵심 포인트:
- 각 버킷은 8개의 key-value 쌍 저장
- tophash로 빠른 비교 (full hash 계산 전에 필터링)
- Load factor 6.5 (평균적으로 버킷당 6.5개 요소)
*/

func MapInternalsDemo() {
	fmt.Println("\n=== Map 내부 구조 분석 ===")

	// 1. Map 헤더 크기
	m := make(map[string]int)
	fmt.Printf("Map 변수 크기 (포인터): %d bytes\n", unsafe.Sizeof(m))

	// 2. Map은 포인터이므로 함수 전달 시 복사 없음
	m["key1"] = 100
	modifyMap(m)
	fmt.Printf("함수 호출 후: m[\"key1\"] = %d\n", m["key1"])

	// 3. Map 순회 순서는 의도적으로 랜덤
	fmt.Println("\nMap 순회 순서 (매번 다름):")
	testMap := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	for i := 0; i < 3; i++ {
		fmt.Printf("  시도 %d: ", i+1)
		for k := range testMap {
			fmt.Printf("%s ", k)
		}
		fmt.Println()
	}
}

func modifyMap(m map[string]int) {
	m["key1"] = 200 // 원본이 수정됨
}

// ============================================================================
// 4. Map vs Slice 메모리 특성 비교
// ============================================================================

/*
메모리 특성 비교:

| 특성              | Slice                    | Map                       |
|------------------|--------------------------|---------------------------|
| 헤더 크기         | 24 bytes                 | 8 bytes (포인터)           |
| 함수 전달         | 헤더 복사, 배열 공유       | 포인터 복사                |
| nil 상태          | nil slice 가능           | nil map 가능               |
| 요소 주소 획득    | &s[i] 가능               | 불가능 (rehash로 이동 가능) |
| 동시성            | 안전하지 않음             | 안전하지 않음              |
| 메모리 해제       | nil 할당 또는 재슬라이싱   | delete() 또는 nil 할당     |
*/

func MemoryCharacteristicsDemo() {
	fmt.Println("\n=== 메모리 특성 비교 ===")

	// 1. Slice 요소 주소 획득 가능
	slice := []int{1, 2, 3}
	ptr := &slice[0]
	*ptr = 100
	fmt.Printf("Slice 요소 주소 변경: %v\n", slice)

	// 2. Map 요소 주소 획득 불가
	// m := map[string]int{"a": 1}
	// ptr := &m["a"]  // 컴파일 에러!
	fmt.Println("Map 요소 주소 획득: 불가능 (컴파일 에러)")

	// 3. nil slice vs nil map
	var nilSlice []int
	var nilMap map[string]int

	fmt.Printf("\nnil slice에 append: ")
	nilSlice = append(nilSlice, 1) // OK
	fmt.Printf("%v\n", nilSlice)

	// nilMap["key"] = 1  // panic: assignment to entry in nil map
	fmt.Println("nil map에 할당: panic 발생!")
}

// ============================================================================
// 5. Slice 메모리 누수 패턴
// ============================================================================

/*
Slice 메모리 누수의 대표적 패턴:

1. 큰 배열의 작은 슬라이스 유지
2. 슬라이스 요소에 포인터가 있을 때 clear 안 함
3. append 후 이전 배열이 GC되지 않음
*/

type LargeStruct struct {
	data [1024 * 1024]byte // 1MB
}

func SliceMemoryLeakDemo() {
	fmt.Println("\n=== Slice 메모리 누수 패턴 ===")

	// 패턴 1: 큰 배열의 작은 슬라이스
	fmt.Println("1. 큰 배열의 작은 슬라이스 문제:")
	bigData := make([]byte, 1024*1024) // 1MB 할당
	_ = bigData

	// 나쁜 예: 첫 10바이트만 필요하지만 전체 1MB가 유지됨
	// smallSlice := bigData[:10]

	// 좋은 예: 복사하여 새 슬라이스 생성
	smallSlice := make([]byte, 10)
	copy(smallSlice, bigData[:10])
	fmt.Printf("  새 슬라이스 용량: %d bytes\n", cap(smallSlice))

	// 패턴 2: 포인터를 포함한 슬라이스 축소
	fmt.Println("\n2. 포인터 슬라이스 축소 시 메모리 누수:")

	type Item struct {
		data *LargeStruct
	}

	items := make([]*Item, 100)
	for i := range items {
		items[i] = &Item{data: &LargeStruct{}}
	}

	// 나쁜 예: 앞 10개만 유지, 나머지 90개는 GC 안됨
	// items = items[:10]  // 뒤 90개 포인터가 여전히 배열에 존재

	// 좋은 예: nil로 초기화 후 축소
	for i := 10; i < len(items); i++ {
		items[i] = nil // GC가 수집 가능하게 함
	}
	items = items[:10]
	fmt.Printf("  정리 후 슬라이스 길이: %d\n", len(items))
}

// ============================================================================
// 6. Map 메모리 특성 및 주의사항
// ============================================================================

/*
Map 메모리 특성:

1. Map은 축소되지 않음: delete()로 요소를 삭제해도 버킷 메모리는 유지
2. 많은 요소 삭제 후에도 메모리 사용량 동일
3. 해결책: 새 map으로 복사 또는 sync.Pool 활용
*/

func MapMemoryDemo() {
	fmt.Println("\n=== Map 메모리 특성 ===")

	// Map은 절대 축소되지 않음
	m := make(map[int]int)

	// 10만 개 추가
	for i := 0; i < 100000; i++ {
		m[i] = i
	}
	fmt.Printf("10만개 추가 후 len: %d\n", len(m))

	// 99,999개 삭제
	for i := 0; i < 99999; i++ {
		delete(m, i)
	}
	fmt.Printf("삭제 후 len: %d (버킷 메모리는 그대로!)\n", len(m))

	// 해결책: 새 map으로 복사
	newMap := make(map[int]int, len(m))
	for k, v := range m {
		newMap[k] = v
	}
	m = newMap
	fmt.Println("새 map으로 복사하여 메모리 정리 완료")
}

// ============================================================================
// 7. 시니어급 최적화 기법
// ============================================================================

func OptimizationTechniques() {
	fmt.Println("\n=== 시니어급 최적화 기법 ===")

	// 1. 용량 미리 할당
	fmt.Println("1. 용량 미리 할당:")

	// 나쁜 예
	var bad []int
	for i := 0; i < 10000; i++ {
		bad = append(bad, i) // 여러 번 재할당
	}

	// 좋은 예
	good := make([]int, 0, 10000)
	for i := 0; i < 10000; i++ {
		good = append(good, i) // 재할당 없음
	}
	fmt.Println("  용량 지정으로 재할당 방지")

	// 2. Map 용량 힌트
	fmt.Println("\n2. Map 용량 힌트:")
	m := make(map[string]int, 10000) // 초기 버킷 할당
	_ = m
	fmt.Println("  make(map[K]V, hint)로 초기 용량 지정")

	// 3. 슬라이스 재사용
	fmt.Println("\n3. 슬라이스 재사용 ([:0] 트릭):")
	buffer := make([]byte, 0, 1024)

	for i := 0; i < 3; i++ {
		buffer = buffer[:0] // 길이만 0으로, 용량 유지
		buffer = append(buffer, []byte("reused data")...)
		fmt.Printf("  반복 %d: len=%d, cap=%d\n", i+1, len(buffer), cap(buffer))
	}
}

// ============================================================================
// 면접 답변 요약
// ============================================================================

/*
시니어급 면접 답변 포인트:

Q: Slice와 Map의 차이점은?

A: (내부 구조 관점에서)

1. Slice:
   - 24바이트 헤더 (ptr, len, cap)
   - 연속된 메모리 배열
   - 함수 전달 시 헤더 복사, 배열 공유
   - 요소 주소 획득 가능 (&s[i])
   - 용량 초과 시 새 배열 할당 후 복사

2. Map:
   - 8바이트 포인터 (hmap 구조체 참조)
   - 해시 테이블 + 버킷 구조
   - 각 버킷에 8개 key-value 저장
   - 요소 주소 획득 불가 (rehash로 위치 변경 가능)
   - 삭제해도 메모리 축소 안됨

3. 메모리 관리:
   - Slice: 큰 배열의 작은 슬라이스 유지 시 누수
   - Map: delete() 후에도 버킷 메모리 유지
   - 둘 다 nil 할당 또는 새로 생성으로 메모리 해제

4. 최적화:
   - 용량 미리 할당 (make with capacity)
   - [:0] 트릭으로 슬라이스 재사용
   - Map은 새로 생성하여 메모리 정리
*/

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
