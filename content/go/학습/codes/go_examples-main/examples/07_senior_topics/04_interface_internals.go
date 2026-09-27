package senior_topics

/*
================================================================================
Go Interface 내부 구조: iface, eface, itab (시니어 레벨)
================================================================================

면접 질문: "Go interface의 내부 구조를 설명해주세요. iface와 eface의 차이는?"

시니어급 답변 포인트:
1. iface vs eface 구조
2. itab (interface table) 역할
3. 메서드 디스패치 메커니즘
4. 성능 특성
5. interface{} (any)의 비용
*/

import (
	"fmt"
	"unsafe"
)

// ============================================================================
// 1. iface와 eface 구조
// ============================================================================

/*
Go runtime에서의 인터페이스 표현:

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

핵심 차이:
- eface: 메서드가 없으므로 _type만 필요
- iface: 메서드 디스패치를 위해 itab 필요
*/

// 시각화용 구조체 (실제 runtime 구조 모사)
type eface struct {
	_type unsafe.Pointer // *_type
	data  unsafe.Pointer
}

type iface struct {
	tab  unsafe.Pointer // *itab
	data unsafe.Pointer
}

/*
_type 구조 (간략화):
type _type struct {
    size       uintptr  // 타입 크기
    ptrdata    uintptr  // 포인터를 포함하는 바이트 수
    hash       uint32   // 타입 해시 (빠른 비교용)
    tflag      tflag    // 플래그
    align      uint8    // 정렬
    fieldAlign uint8    // 필드 정렬
    kind       uint8    // 타입 종류 (int, struct, etc.)
    ...
    str        nameOff  // 타입 이름
}

itab 구조:
type itab struct {
    inter *interfacetype  // 인터페이스 타입 정보
    _type *_type          // 구체적 타입 정보
    hash  uint32          // _type.hash 복사 (빠른 타입 assertion)
    _     [4]byte         // 패딩
    fun   [1]uintptr      // 메서드 주소 배열 (가변 크기)
}
*/

func InterfaceInternalsDemo() {
	fmt.Println("=== Interface 내부 구조 데모 ===")

	// 1. 인터페이스 크기 확인
	var empty interface{}
	var reader fmt.Stringer

	fmt.Printf("interface{} (eface) 크기: %d bytes\n", unsafe.Sizeof(empty))
	fmt.Printf("Stringer (iface) 크기: %d bytes\n", unsafe.Sizeof(reader))

	// 2. nil 인터페이스 vs nil 값을 가진 인터페이스
	var nilInterface interface{} = nil
	var nilPointer *int = nil
	var interfaceWithNil interface{} = nilPointer

	fmt.Printf("\nnil interface == nil: %v\n", nilInterface == nil)
	fmt.Printf("nil 포인터를 담은 interface == nil: %v\n", interfaceWithNil == nil)
	fmt.Println("(타입 정보가 있으므로 nil이 아님!)")
}

// ============================================================================
// 2. itab (Interface Table) 상세
// ============================================================================

/*
itab의 역할:
1. 구체적 타입과 인터페이스 타입의 연결
2. 메서드 주소 저장 (메서드 디스패치)
3. 타입 assertion 최적화

itab 캐싱:
- 동일한 (인터페이스 타입, 구체적 타입) 쌍은 한 번만 생성
- 전역 해시 테이블에 캐싱
- 타입 assertion 시 O(1) 조회
*/

// 예시: Stringer 인터페이스
type MyString string

func (s MyString) String() string {
	return string(s)
}

func ItabDemo() {
	fmt.Println("\n=== itab 데모 ===")

	// 동일한 타입의 여러 인터페이스
	var s1 fmt.Stringer = MyString("hello")
	var s2 fmt.Stringer = MyString("world")

	// 두 인터페이스는 같은 itab을 공유
	// (같은 인터페이스 타입 + 같은 구체적 타입)

	// itab 비교 (unsafe 사용)
	iface1 := (*iface)(unsafe.Pointer(&s1))
	iface2 := (*iface)(unsafe.Pointer(&s2))

	fmt.Printf("s1의 itab 포인터: %v\n", iface1.tab)
	fmt.Printf("s2의 itab 포인터: %v\n", iface2.tab)
	fmt.Printf("같은 itab 공유: %v\n", iface1.tab == iface2.tab)
}

// ============================================================================
// 3. 메서드 디스패치 메커니즘
// ============================================================================

/*
메서드 호출 과정:

1. 인터페이스 변수에서 itab 획득
2. itab.fun 배열에서 메서드 주소 조회
3. data 포인터와 함께 메서드 호출

최적화:
- itab.fun은 인터페이스 메서드 선언 순서대로 정렬
- 컴파일 시점에 메서드 인덱스 결정
- O(1) 메서드 조회
*/

type Writer interface {
	Write([]byte) (int, error)
}

type Reader interface {
	Read([]byte) (int, error)
}

type ReadWriter interface {
	Reader
	Writer
}

type MyBuffer struct {
	data []byte
}

func (b *MyBuffer) Read(p []byte) (int, error) {
	n := copy(p, b.data)
	b.data = b.data[n:]
	return n, nil
}

func (b *MyBuffer) Write(p []byte) (int, error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func MethodDispatchDemo() {
	fmt.Println("\n=== 메서드 디스패치 데모 ===")

	buf := &MyBuffer{}

	// 다른 인터페이스로 할당 - 각각 다른 itab
	var w Writer = buf
	var r Reader = buf
	var rw ReadWriter = buf

	// 각 인터페이스마다 다른 itab
	// Writer의 itab.fun[0] = Write 메서드
	// Reader의 itab.fun[0] = Read 메서드
	// ReadWriter의 itab.fun[0] = Read, itab.fun[1] = Write

	_, _ = w.Write([]byte("hello"))
	_, _ = rw.Read(make([]byte, 5))
	_ = r

	fmt.Println("각 인터페이스는 자체 itab을 가짐")
}

// ============================================================================
// 4. Type Assertion 내부 동작
// ============================================================================

/*
Type Assertion 유형:

1. 구체적 타입으로 (v, ok := i.(ConcreteType))
   - itab._type.hash 비교로 빠른 체크
   - 실패 시 false 반환 또는 panic

2. 인터페이스 타입으로 (v, ok := i.(Interface))
   - 구체적 타입이 새 인터페이스를 구현하는지 확인
   - 새 itab 생성 또는 캐시에서 조회

Type Switch:
- 컴파일러가 최적화된 점프 테이블 생성
- 자주 사용되는 타입을 먼저 검사
*/

func TypeAssertionDemo() {
	fmt.Println("\n=== Type Assertion 내부 동작 ===")

	var i interface{} = MyString("test")

	// 1. 구체적 타입으로 assertion
	if s, ok := i.(MyString); ok {
		fmt.Printf("MyString assertion 성공: %s\n", s)
	}

	// 2. 인터페이스 타입으로 assertion
	if stringer, ok := i.(fmt.Stringer); ok {
		fmt.Printf("Stringer assertion 성공: %s\n", stringer.String())
	}

	// 3. Type switch
	switch v := i.(type) {
	case MyString:
		fmt.Printf("Type switch - MyString: %s\n", v)
	case fmt.Stringer:
		fmt.Printf("Type switch - Stringer: %s\n", v.String())
	default:
		fmt.Printf("알 수 없는 타입: %T\n", v)
	}
}

// ============================================================================
// 5. interface{} (any) 사용의 비용
// ============================================================================

/*
interface{} 사용 시 비용:

1. Boxing/Unboxing:
   - 값을 interface{}에 할당 시 힙 할당 가능
   - 작은 값(1 word 이하)은 직접 저장, 큰 값은 포인터

2. 타입 정보 유지:
   - 모든 interface{}는 _type 포인터 필요 (8 bytes)

3. 런타임 타입 체크:
   - 컴파일 타임 타입 체크 불가
   - Type assertion 비용

4. 최적화 방해:
   - 컴파일러의 인라이닝 제한
   - Escape analysis 어려움
*/

func InterfaceBoxingDemo() {
	fmt.Println("\n=== interface{} Boxing 비용 ===")

	// 작은 값: word 크기 이하는 직접 저장
	var small interface{} = 42
	_ = small

	// 큰 값: 힙 할당 후 포인터 저장
	type BigStruct struct {
		a, b, c, d int64
	}
	var big interface{} = BigStruct{1, 2, 3, 4}
	_ = big

	fmt.Println("작은 값 (int): 직접 저장")
	fmt.Println("큰 값 (struct): 힙 할당 후 포인터 저장")

	// 제네릭을 사용한 대안 (Go 1.18+)
	// interface{} 대신 타입 파라미터 사용
	fmt.Println("\n제네릭 대안:")
	fmt.Println("func Process[T any](v T) T  // 박싱 없음")
}

// ============================================================================
// 6. 인터페이스 설계 모범 사례
// ============================================================================

/*
시니어급 인터페이스 설계 원칙:

1. 작은 인터페이스 선호:
   - 1-3개 메서드
   - io.Reader, io.Writer가 대표적
   - 구현이 쉽고 조합 가능

2. 소비자 측에서 정의:
   - "Accept interfaces, return structs"
   - 필요한 기능만 인터페이스로 요구

3. nil 체크 주의:
   - nil 인터페이스 vs nil 값을 가진 인터페이스
   - 명시적으로 nil 반환 (var err error = nil)

4. 빈 인터페이스 최소화:
   - 타입 안전성 상실
   - 제네릭으로 대체 가능한지 검토
*/

// 좋은 예: 작은 인터페이스
type Closer interface {
	Close() error
}

// 조합으로 확장
type ReadCloser interface {
	Reader
	Closer
}

// 소비자 측 정의 예시
func ProcessData(r Reader) error {
	// Reader만 필요하면 Reader만 요구
	buf := make([]byte, 1024)
	_, err := r.Read(buf)
	return err
}

// ============================================================================
// 7. 시니어 면접 답변 요약
// ============================================================================

/*
Q: Go interface의 내부 구조를 설명해주세요.

A: Go의 인터페이스는 두 가지 내부 표현이 있습니다.

1. eface (empty interface = interface{}):
   - _type: 타입 정보 포인터
   - data: 실제 값 포인터
   - 16 bytes, 메서드가 없을 때 사용

2. iface (non-empty interface):
   - tab: itab 포인터 (인터페이스 테이블)
   - data: 실제 값 포인터
   - 16 bytes, 메서드가 있을 때 사용

3. itab (Interface Table):
   - 구체적 타입과 인터페이스 타입의 연결
   - fun 배열에 메서드 주소 저장
   - 동일한 (인터페이스, 타입) 쌍은 전역 캐시에 저장
   - 메서드 디스패치는 O(1)

4. 성능 고려사항:
   - interface{} 할당 시 큰 값은 힙 할당
   - Type assertion은 해시 비교로 최적화
   - 작은 인터페이스가 itab 크기도 작음

5. 주의점:
   - nil 인터페이스와 nil 값을 가진 인터페이스 구분
   - (타입 정보가 있으면 nil이 아님)
   - 제네릭으로 interface{} 사용 줄이기
*/
