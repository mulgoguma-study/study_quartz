// =============================================================================
// Go 에러 처리 예제 - errors.Is와 errors.As
// =============================================================================
// Go 1.13부터 도입된 에러 래핑(wrapping)과 검사 기능
//
// errors.Is(err, target):
// - err가 target과 같은 에러인지 확인
// - 래핑된 에러 체인을 따라가며 확인
// - 센티널 에러(미리 정의된 에러 값) 비교에 사용
//
// errors.As(err, target):
// - err를 target 타입으로 변환 가능한지 확인
// - 래핑된 에러 체인을 따라가며 확인
// - 커스텀 에러 타입의 추가 정보 접근에 사용
//
// fmt.Errorf("... %w", err):
// - %w 동사를 사용하여 에러 래핑
// - 원본 에러를 보존하면서 추가 컨텍스트 제공
// =============================================================================

package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// =============================================================================
// 1. 센티널 에러 정의 (Sentinel Errors)
// =============================================================================
// 센티널 에러: 미리 정의된 에러 값으로, errors.Is로 비교

var (
	// 사용자 관련 에러
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidEmail = errors.New("invalid email format")

	// 인증 관련 에러
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")

	// 데이터 관련 에러
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
)

// =============================================================================
// 2. 커스텀 에러 타입 정의
// =============================================================================

// ValidationError 검증 에러 - 필드별 에러 정보 포함
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
}

// Error error 인터페이스 구현
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s (value: %v)",
		e.Field, e.Message, e.Value)
}

// DatabaseError 데이터베이스 에러 - SQL 정보 포함
type DatabaseError struct {
	Operation string
	Table     string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s on table '%s': %v",
		e.Operation, e.Table, e.Err)
}

// Unwrap 래핑된 에러 반환 (errors.Is, errors.As가 체인을 따라갈 수 있게)
func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// HTTPError HTTP 에러 - 상태 코드 포함
type HTTPError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

func main() {
	fmt.Println("=== 1. errors.Is 예제 ===")
	errorsIsExample()

	fmt.Println("\n=== 2. errors.As 예제 ===")
	errorsAsExample()

	fmt.Println("\n=== 3. 에러 래핑 예제 ===")
	errorWrappingExample()

	fmt.Println("\n=== 4. 실전 예제: 사용자 조회 시스템 ===")
	practicalExample()
}

// errorsIsExample errors.Is 사용 예제
func errorsIsExample() {
	// ==========================================================================
	// 기본 사용법
	// ==========================================================================
	err := findUser("unknown@example.com")

	// errors.Is로 특정 에러인지 확인
	if errors.Is(err, ErrUserNotFound) {
		fmt.Println("✓ 사용자를 찾을 수 없습니다 (ErrUserNotFound)")
	}

	// ==========================================================================
	// 래핑된 에러도 확인 가능
	// ==========================================================================
	wrappedErr := fmt.Errorf("조회 실패: %w", ErrUserNotFound)

	// 래핑되어도 errors.Is가 내부 에러를 찾아냄
	if errors.Is(wrappedErr, ErrUserNotFound) {
		fmt.Println("✓ 래핑된 에러에서도 ErrUserNotFound 감지")
	}

	// ==========================================================================
	// 표준 라이브러리 에러와 함께 사용
	// ==========================================================================
	_, fileErr := os.Open("nonexistent_file.txt")

	if errors.Is(fileErr, fs.ErrNotExist) {
		fmt.Println("✓ 파일이 존재하지 않습니다 (fs.ErrNotExist)")
	}

	// ==========================================================================
	// 다중 에러 체인
	// ==========================================================================
	// 여러 번 래핑된 에러
	level1 := fmt.Errorf("level 1: %w", ErrNotFound)
	level2 := fmt.Errorf("level 2: %w", level1)
	level3 := fmt.Errorf("level 3: %w", level2)

	// 가장 깊이 래핑된 에러도 찾아냄
	if errors.Is(level3, ErrNotFound) {
		fmt.Println("✓ 3단계 래핑된 에러에서도 ErrNotFound 감지")
	}
}

// errorsAsExample errors.As 사용 예제
func errorsAsExample() {
	// ==========================================================================
	// 기본 사용법
	// ==========================================================================
	err := validateInput("email", "invalid-email", "이메일")

	// errors.As로 특정 타입의 에러인지 확인하고 변환
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		fmt.Printf("✓ ValidationError 감지:\n")
		fmt.Printf("  - 필드: %s\n", validationErr.Field)
		fmt.Printf("  - 메시지: %s\n", validationErr.Message)
		fmt.Printf("  - 값: %v\n", validationErr.Value)
	}

	// ==========================================================================
	// 래핑된 에러에서도 타입 추출
	// ==========================================================================
	dbErr := &DatabaseError{
		Operation: "INSERT",
		Table:     "users",
		Err:       ErrInvalidInput,
	}

	// 추가 컨텍스트와 함께 래핑
	wrappedDbErr := fmt.Errorf("사용자 생성 실패: %w", dbErr)

	var extractedDbErr *DatabaseError
	if errors.As(wrappedDbErr, &extractedDbErr) {
		fmt.Printf("✓ DatabaseError 추출:\n")
		fmt.Printf("  - 작업: %s\n", extractedDbErr.Operation)
		fmt.Printf("  - 테이블: %s\n", extractedDbErr.Table)
	}

	// 내부 에러도 확인 가능
	if errors.Is(wrappedDbErr, ErrInvalidInput) {
		fmt.Println("✓ DatabaseError 내부에 ErrInvalidInput 있음")
	}

	// ==========================================================================
	// 표준 라이브러리 에러 타입 추출
	// ==========================================================================
	_, pathErr := os.Open("invalid/path/file.txt")

	var pathError *fs.PathError
	if errors.As(pathErr, &pathError) {
		fmt.Printf("✓ PathError 추출:\n")
		fmt.Printf("  - 작업: %s\n", pathError.Op)
		fmt.Printf("  - 경로: %s\n", pathError.Path)
	}
}

// errorWrappingExample 에러 래핑 예제
func errorWrappingExample() {
	// 원본 에러
	originalErr := ErrNotFound

	// 1단계 래핑: 저장소 계층
	repoErr := fmt.Errorf("repository: user lookup failed: %w", originalErr)

	// 2단계 래핑: 서비스 계층
	serviceErr := fmt.Errorf("service: could not get user profile: %w", repoErr)

	// 3단계 래핑: 핸들러 계층
	handlerErr := fmt.Errorf("handler: request processing failed: %w", serviceErr)

	// 에러 메시지 출력 (전체 체인)
	fmt.Printf("전체 에러: %v\n", handlerErr)

	// 에러 체인 분석
	fmt.Println("\n에러 체인 분석:")

	err := handlerErr
	for i := 0; err != nil; i++ {
		fmt.Printf("  Level %d: %v\n", i, err)
		err = errors.Unwrap(err) // 다음 레벨로
	}

	// 원본 에러 확인
	fmt.Printf("\n원본 에러 확인: errors.Is(handlerErr, ErrNotFound) = %v\n",
		errors.Is(handlerErr, ErrNotFound))
}

// =============================================================================
// 실전 예제: 사용자 조회 시스템
// =============================================================================

// UserRepository 사용자 저장소 (시뮬레이션)
type UserRepository struct{}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
	// 시뮬레이션: 특정 이메일만 존재
	if email == "admin@example.com" {
		return &User{ID: 1, Email: email, Name: "Admin"}, nil
	}

	// 데이터베이스 에러 시뮬레이션
	if email == "error@example.com" {
		return nil, &DatabaseError{
			Operation: "SELECT",
			Table:     "users",
			Err:       fmt.Errorf("connection timeout"),
		}
	}

	return nil, ErrUserNotFound
}

type User struct {
	ID    int
	Email string
	Name  string
}

// UserService 사용자 서비스
type UserService struct {
	repo *UserRepository
}

func (s *UserService) GetUserByEmail(email string) (*User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		// 에러 래핑하여 컨텍스트 추가
		return nil, fmt.Errorf("UserService.GetUserByEmail(%s): %w", email, err)
	}
	return user, nil
}

// APIHandler API 핸들러
type APIHandler struct {
	userService *UserService
}

func (h *APIHandler) HandleGetUser(email string) {
	user, err := h.userService.GetUserByEmail(email)

	if err != nil {
		h.handleError(err)
		return
	}

	fmt.Printf("✓ 사용자 조회 성공: ID=%d, Name=%s\n", user.ID, user.Name)
}

func (h *APIHandler) handleError(err error) {
	fmt.Printf("에러 발생: %v\n", err)

	// 에러 타입에 따른 처리
	switch {
	case errors.Is(err, ErrUserNotFound):
		// 404 Not Found
		fmt.Println("→ HTTP 404: 사용자를 찾을 수 없습니다")

	case errors.Is(err, ErrUnauthorized):
		// 401 Unauthorized
		fmt.Println("→ HTTP 401: 인증이 필요합니다")

	default:
		// 데이터베이스 에러 확인
		var dbErr *DatabaseError
		if errors.As(err, &dbErr) {
			// 503 Service Unavailable
			fmt.Printf("→ HTTP 503: 데이터베이스 오류 (%s on %s)\n",
				dbErr.Operation, dbErr.Table)
			return
		}

		// 500 Internal Server Error
		fmt.Println("→ HTTP 500: 내부 서버 오류")
	}
}

func practicalExample() {
	handler := &APIHandler{
		userService: &UserService{
			repo: &UserRepository{},
		},
	}

	fmt.Println("--- 존재하는 사용자 조회 ---")
	handler.HandleGetUser("admin@example.com")

	fmt.Println("\n--- 존재하지 않는 사용자 조회 ---")
	handler.HandleGetUser("unknown@example.com")

	fmt.Println("\n--- 데이터베이스 오류 발생 ---")
	handler.HandleGetUser("error@example.com")
}

// =============================================================================
// 헬퍼 함수들
// =============================================================================

func findUser(email string) error {
	// 시뮬레이션: 항상 사용자를 찾을 수 없음
	return ErrUserNotFound
}

func validateInput(field string, value interface{}, fieldName string) error {
	// 시뮬레이션: 검증 실패
	return &ValidationError{
		Field:   field,
		Message: fmt.Sprintf("%s 형식이 올바르지 않습니다", fieldName),
		Value:   value,
	}
}
