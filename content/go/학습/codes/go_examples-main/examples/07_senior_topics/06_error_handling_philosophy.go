package senior_topics

/*
================================================================================
Go 에러 처리 철학과 에러 래핑 전략 (시니어 레벨)
================================================================================

면접 질문: "왜 Go는 Exception이 없고 Error를 반환하나요?"
         "시니어로서 에러 래핑(Wrapping) 전략은?"

시니어급 답변 포인트:
1. Go의 에러 처리 철학
2. Exception vs Error 반환 비교
3. 에러 래핑과 Unwrap
4. errors.Is와 errors.As
5. Sentinel Error vs Custom Error Type
6. 에러 처리 모범 사례
*/

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
)

// ============================================================================
// 1. Go의 에러 처리 철학
// ============================================================================

/*
Go가 Exception 대신 Error를 선택한 이유:

1. 명시적 제어 흐름 (Explicit Control Flow):
   - Exception은 암묵적인 제어 흐름 생성
   - Error 반환은 코드에서 에러 경로가 명확히 보임
   - "Errors are values" - 에러도 그냥 값이다

2. 로컬 처리 강제 (Local Handling):
   - 각 함수 호출 지점에서 에러 처리 결정
   - 에러를 무시하는 것도 명시적 선택 (_)
   - 멀리 떨어진 catch 블록 없음

3. 성능:
   - Exception은 스택 unwinding 비용
   - Error 반환은 일반 함수 반환과 동일

4. 단순성:
   - try-catch-finally 복잡한 구문 없음
   - 에러 타입 계층 불필요

Rob Pike의 말:
"Errors are values. They can be programmed, and since they're values,
they can be used to control flow."
*/

// Exception 스타일 (다른 언어)
/*
try {
    file := openFile("data.txt")
    data := readFile(file)
    processData(data)
} catch (FileNotFoundException e) {
    // 파일 없음 처리
} catch (IOException e) {
    // IO 에러 처리
} finally {
    file.close()
}
*/

// Go 스타일 - 명시적
func GoStyleExample() error {
	file, err := os.Open("data.txt")
	if err != nil {
		return fmt.Errorf("파일 열기 실패: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("파일 읽기 실패: %w", err)
	}

	if err := processData(data); err != nil {
		return fmt.Errorf("데이터 처리 실패: %w", err)
	}

	return nil
}

func processData(data []byte) error {
	if len(data) == 0 {
		return errors.New("빈 데이터")
	}
	return nil
}

// ============================================================================
// 2. panic과 recover - Exception과 비슷하지만 다름
// ============================================================================

/*
panic의 올바른 사용:
- 프로그램이 계속될 수 없는 심각한 오류
- 프로그래머의 실수 (버그)
- 초기화 실패

panic을 사용하지 말아야 할 때:
- 예상 가능한 에러 (파일 없음, 네트워크 타임아웃 등)
- 일반적인 비즈니스 로직 에러
- API 경계를 넘어서

"Don't panic" - 일반적인 에러에 panic 사용 금지
*/

// 잘못된 사용
func badPanicUsage(filename string) []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err) // 잘못됨! 파일 없음은 예상 가능한 에러
	}
	return data
}

// 올바른 사용 - 프로그래머 실수
func mustCompile(pattern string) {
	// 컴파일 타임에 결정되는 정규식 패턴
	// 실패하면 프로그래머 실수이므로 panic 적절
	// regexp.MustCompile(pattern)
	_ = pattern
}

// recover 사용 - HTTP 핸들러 경계
func RecoverMiddleware(handler func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic 복구: %v\n", r)
			// 로깅, 500 에러 반환 등
		}
	}()
	handler()
}

// ============================================================================
// 3. 에러 래핑 (Error Wrapping) - Go 1.13+
// ============================================================================

/*
에러 래핑의 목적:
1. 컨텍스트 추가 (어디서, 왜 발생했는지)
2. 원본 에러 보존 (체인)
3. 에러 검사 가능성 유지

%w 동사:
- fmt.Errorf("context: %w", err)
- errors.Unwrap()으로 원본 접근 가능
- errors.Is(), errors.As()와 호환

%v 동사:
- fmt.Errorf("context: %v", err)
- 원본 에러와의 체인 끊김
- Unwrap 불가
*/

// 정의된 에러
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
)

// Repository 계층
func FindUser(id string) (string, error) {
	// DB에서 사용자 조회 시뮬레이션
	if id == "" {
		return "", ErrInvalidInput
	}
	if id == "deleted" {
		return "", ErrNotFound
	}
	return "user_" + id, nil
}

// Service 계층 - 에러 래핑
func GetUserProfile(id string) (string, error) {
	user, err := FindUser(id)
	if err != nil {
		// %w로 원본 에러 보존하면서 컨텍스트 추가
		return "", fmt.Errorf("GetUserProfile(id=%s): %w", id, err)
	}
	return user, nil
}

// Handler 계층 - 에러 검사
func HandleGetUser(id string) {
	user, err := GetUserProfile(id)
	if err != nil {
		// errors.Is로 원본 에러 타입 검사
		if errors.Is(err, ErrNotFound) {
			fmt.Println("404 Not Found")
			return
		}
		if errors.Is(err, ErrInvalidInput) {
			fmt.Println("400 Bad Request")
			return
		}
		fmt.Printf("500 Internal Error: %v\n", err)
		return
	}
	fmt.Printf("200 OK: %s\n", user)
}

func ErrorWrappingDemo() {
	fmt.Println("=== 에러 래핑 데모 ===")

	HandleGetUser("123")     // 성공
	HandleGetUser("")        // Invalid Input
	HandleGetUser("deleted") // Not Found
}

// ============================================================================
// 4. errors.Is vs errors.As
// ============================================================================

/*
errors.Is(err, target):
- err 체인에서 target과 같은 에러가 있는지 확인
- Sentinel Error (값) 비교에 사용
- Is() 메서드 구현으로 커스텀 비교 가능

errors.As(err, target):
- err 체인에서 target 타입으로 변환 가능한 에러 찾기
- Custom Error Type에 대한 타입 assertion
- As() 메서드 구현으로 커스텀 변환 가능
*/

// Custom Error Type
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}

// 래핑된 ValidationError
func ValidateEmail(email string) error {
	if email == "" {
		return &ValidationError{Field: "email", Message: "required"}
	}
	return nil
}

func ProcessRegistration(email string) error {
	if err := ValidateEmail(email); err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}
	return nil
}

func IsVsAsDemo() {
	fmt.Println("\n=== errors.Is vs errors.As 데모 ===")

	// errors.Is 사용 (Sentinel Error)
	_, err := FindUser("deleted")
	if errors.Is(err, ErrNotFound) {
		fmt.Println("errors.Is: ErrNotFound 감지")
	}

	// errors.As 사용 (Custom Error Type)
	err = ProcessRegistration("")
	var validErr *ValidationError
	if errors.As(err, &validErr) {
		fmt.Printf("errors.As: ValidationError 감지 - Field: %s\n", validErr.Field)
	}

	// 래핑되어도 작동
	wrappedErr := fmt.Errorf("outer: %w", err)
	if errors.As(wrappedErr, &validErr) {
		fmt.Println("errors.As: 래핑된 에러에서도 타입 추출 성공")
	}
}

// ============================================================================
// 5. Sentinel Error vs Custom Error Type 선택
// ============================================================================

/*
Sentinel Error (var ErrXxx = errors.New("...")):
- 단순한 에러 조건
- 추가 정보 불필요
- 예: io.EOF, sql.ErrNoRows

사용 시점:
- 에러 종류만 알면 충분할 때
- 패키지 API의 일부로 공개할 때

Custom Error Type (type XxxError struct):
- 추가 컨텍스트 필요
- 에러 발생 상황의 상세 정보
- 에러 복구에 필요한 데이터

사용 시점:
- 에러에 메타데이터가 필요할 때
- 클라이언트가 에러 정보로 결정해야 할 때
*/

// Sentinel Error 예시
var (
	ErrConnectionClosed = errors.New("connection closed")
	ErrTimeout          = errors.New("operation timeout")
)

// Custom Error Type 예시
type HTTPError struct {
	StatusCode int
	Message    string
	RequestID  string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("[%d] %s (request_id: %s)", e.StatusCode, e.Message, e.RequestID)
}

// Temporary 에러 인터페이스 (재시도 가능 여부)
type TemporaryError interface {
	Temporary() bool
}

type NetworkError struct {
	Message string
	Retry   bool
}

func (e *NetworkError) Error() string {
	return e.Message
}

func (e *NetworkError) Temporary() bool {
	return e.Retry
}

func ErrorTypeSelectionDemo() {
	fmt.Println("\n=== 에러 타입 선택 데모 ===")

	// Sentinel Error 사용
	err := sql.ErrNoRows
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("Sentinel: 레코드 없음")
	}

	// Custom Error Type 사용
	httpErr := &HTTPError{
		StatusCode: 404,
		Message:    "User not found",
		RequestID:  "abc123",
	}
	fmt.Printf("Custom: %v\n", httpErr)

	// Temporary 인터페이스로 재시도 결정
	netErr := &NetworkError{Message: "connection reset", Retry: true}
	if tempErr, ok := error(netErr).(TemporaryError); ok && tempErr.Temporary() {
		fmt.Println("재시도 가능한 에러")
	}
}

// ============================================================================
// 6. 에러 래핑 전략 (시니어급)
// ============================================================================

/*
래핑 전략 원칙:

1. 컨텍스트는 행동을 설명 (동사로 시작):
   - 좋음: "querying user: %w"
   - 나쁨: "user error: %w"

2. 중복 피하기:
   - 호출 스택 전체에서 같은 정보 반복 금지
   - 새로운 정보가 있을 때만 래핑

3. 공개 API 경계에서 래핑:
   - 내부 구현 세부사항 숨기기
   - 패키지 고유 에러로 변환

4. 에러 체인 적절히 끊기:
   - 보안상 민감한 정보 노출 방지
   - %v로 체인 끊기

5. 로깅과 반환 분리:
   - 로깅: 상세 정보 포함
   - 반환: 클라이언트에 필요한 정보만
*/

// 나쁜 예: 과도한 래핑
func badWrapping() error {
	err := errors.New("db error")
	err = fmt.Errorf("repository error: %w", err)
	err = fmt.Errorf("service error: %w", err)
	err = fmt.Errorf("handler error: %w", err)
	// 결과: "handler error: service error: repository error: db error"
	// 정보가 중복되고 유용하지 않음
	return err
}

// 좋은 예: 의미 있는 컨텍스트
func goodWrapping(userID string) error {
	err := queryDatabase(userID)
	if err != nil {
		return fmt.Errorf("finding user %s: %w", userID, err)
	}
	return nil
}

func queryDatabase(id string) error {
	return errors.New("connection timeout")
}

// 에러 체인 끊기 (보안)
func PublicAPI(secret string) error {
	err := internalProcess(secret)
	if err != nil {
		// %v로 체인 끊음 - 내부 에러 노출 방지
		return fmt.Errorf("operation failed: %v", err)
	}
	return nil
}

func internalProcess(secret string) error {
	return fmt.Errorf("invalid secret: %s", secret) // 민감한 정보
}

// ============================================================================
// 7. 에러 처리 모범 사례
// ============================================================================

/*
모범 사례 체크리스트:

□ 에러 무시하지 않기
  - 명시적으로 _ 사용하고 주석 달기
  - 또는 반드시 처리

□ 에러는 한 번만 처리
  - 로깅하고 반환하면 중복 로깅 발생
  - 처리(로깅/반환) 중 하나만

□ 에러 타입 적절히 선택
  - Sentinel: 단순 조건
  - Custom Type: 추가 정보 필요

□ 래핑 시 컨텍스트 추가
  - 어떤 작업에서 실패했는지
  - 관련 파라미터 (ID 등)

□ 공개 API에서 에러 정리
  - 내부 구현 숨기기
  - 패키지 문서에 에러 명시
*/

// 예시: 에러 처리 흐름
func CompleteExample() {
	fmt.Println("\n=== 에러 처리 모범 사례 데모 ===")

	// 1. 에러 무시하지 않기
	data, err := os.ReadFile("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("설정 파일 없음 - 기본값 사용")
			// 기본값으로 처리
		} else {
			fmt.Printf("설정 로드 실패: %v\n", err)
			return
		}
	}
	_ = data

	// 2. 에러 한 번만 처리
	// 나쁨: log.Error(err); return err
	// 좋음: return fmt.Errorf("context: %w", err)  // 상위에서 로깅

	// 3. errors.Is로 특정 에러 처리
	_, err = FindUser("deleted")
	if errors.Is(err, ErrNotFound) {
		fmt.Println("사용자 없음 - 새로 생성")
	}
}

// ============================================================================
// 8. 시니어 면접 답변 요약
// ============================================================================

/*
Q: 왜 Go는 Exception이 없고 Error를 반환하나요?

A: Go는 명시적인 에러 처리를 통해 더 명확하고 예측 가능한 코드를 지향합니다.

1. Exception의 문제점:
   - 암묵적 제어 흐름 (어디서 catch될지 불명확)
   - 스택 unwinding 성능 비용
   - 에러 경로가 코드에서 보이지 않음

2. Error 반환의 장점:
   - 명시적 제어 흐름 (에러 경로가 코드에 보임)
   - 각 호출 지점에서 처리 결정
   - "Errors are values" - 프로그래밍 가능

3. panic의 위치:
   - 복구 불가능한 오류 (프로그래머 실수)
   - 일반 에러에는 사용 금지
   - API 경계에서 recover로 방어

Q: 시니어로서 에러 래핑 전략은?

A: 에러 래핑은 컨텍스트를 추가하면서 원본을 보존하는 것입니다.

1. %w 사용:
   - 체인 유지, errors.Is/As 호환
   - fmt.Errorf("context: %w", err)

2. 래핑 원칙:
   - 동사로 시작하는 컨텍스트
   - 새 정보가 있을 때만 래핑
   - 보안상 민감한 정보는 체인 끊기

3. 에러 타입 선택:
   - Sentinel: 단순 조건 (io.EOF)
   - Custom Type: 추가 정보 필요

4. 처리 원칙:
   - 에러는 한 번만 처리 (로깅 또는 반환)
   - errors.Is로 특정 에러 검사
   - errors.As로 타입 정보 추출
*/
