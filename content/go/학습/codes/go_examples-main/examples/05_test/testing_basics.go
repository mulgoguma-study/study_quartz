// =============================================================================
// Go 테스트 기본 예제 - 테스트 대상 코드
// =============================================================================
// 이 파일은 테스트할 함수들을 정의합니다.
// 테스트 코드는 testing_basics_test.go에 있습니다.
//
// Go 테스트 규칙:
// - 테스트 파일은 _test.go로 끝나야 함
// - 테스트 함수는 Test로 시작하고 *testing.T를 매개변수로 받음
// - 벤치마크 함수는 Benchmark로 시작하고 *testing.B를 매개변수로 받음
// - 예제 함수는 Example로 시작
//
// 실행 방법:
// go test -v                    # 상세 출력
// go test -run TestAdd          # 특정 테스트만 실행
// go test -cover                # 커버리지 확인
// go test -bench=.              # 벤치마크 실행
// =============================================================================

package testexample

import (
	"errors"
	"fmt"
)

// =============================================================================
// 1. 기본 함수들 (단위 테스트 대상)
// =============================================================================

// Add 두 정수를 더함
func Add(a, b int) int {
	return a + b
}

// Subtract 두 정수를 뺌
func Subtract(a, b int) int {
	return a - b
}

// Multiply 두 정수를 곱함
func Multiply(a, b int) int {
	return a * b
}

// Divide 두 정수를 나눔 (0으로 나누면 에러)
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// =============================================================================
// 2. 에러 처리 함수들 (errors.Is, errors.As 테스트 대상)
// =============================================================================

// 센티널 에러 정의
var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized")
)

// ValidationError 검증 에러
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on '%s': %s", e.Field, e.Message)
}

// DatabaseError 데이터베이스 에러
type DatabaseError struct {
	Operation string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s: %v", e.Operation, e.Err)
}

// Unwrap 래핑된 에러 반환
func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// FindUser 사용자 조회 (시뮬레이션)
func FindUser(id int) (string, error) {
	users := map[int]string{
		1: "Alice",
		2: "Bob",
		3: "Charlie",
	}

	if id <= 0 {
		return "", ErrInvalidInput
	}

	name, exists := users[id]
	if !exists {
		return "", fmt.Errorf("user id %d: %w", id, ErrNotFound)
	}

	return name, nil
}

// ValidateEmail 이메일 검증
func ValidateEmail(email string) error {
	if email == "" {
		return &ValidationError{
			Field:   "email",
			Message: "email is required",
		}
	}

	// 간단한 검증 (실제로는 더 복잡한 검증 필요)
	hasAt := false
	for _, c := range email {
		if c == '@' {
			hasAt = true
			break
		}
	}

	if !hasAt {
		return &ValidationError{
			Field:   "email",
			Message: "email must contain @",
		}
	}

	return nil
}

// SaveUser 사용자 저장 (DB 에러 시뮬레이션)
func SaveUser(name string, simulateError bool) error {
	if simulateError {
		return &DatabaseError{
			Operation: "INSERT",
			Err:       ErrInvalidInput,
		}
	}
	return nil
}

// =============================================================================
// 3. 슬라이스/맵 처리 함수들 (테이블 드리븐 테스트 대상)
// =============================================================================

// Contains 슬라이스에 요소가 있는지 확인
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Unique 슬라이스에서 중복 제거
func Unique(slice []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// Reverse 슬라이스 뒤집기
func Reverse(slice []int) []int {
	result := make([]int, len(slice))
	for i, v := range slice {
		result[len(slice)-1-i] = v
	}
	return result
}

// =============================================================================
// 4. 인터페이스와 목킹 대상
// =============================================================================

// UserRepository 사용자 저장소 인터페이스
type UserRepository interface {
	GetByID(id int) (*User, error)
	Save(user *User) error
	Delete(id int) error
}

// User 사용자 구조체
type User struct {
	ID    int
	Name  string
	Email string
}

// UserService 사용자 서비스 (Repository 의존성 주입)
type UserService struct {
	repo UserRepository
}

// NewUserService UserService 생성자
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetUser 사용자 조회
func (s *UserService) GetUser(id int) (*User, error) {
	if id <= 0 {
		return nil, ErrInvalidInput
	}
	return s.repo.GetByID(id)
}

// CreateUser 사용자 생성
func (s *UserService) CreateUser(name, email string) (*User, error) {
	if name == "" {
		return nil, &ValidationError{Field: "name", Message: "name is required"}
	}
	if email == "" {
		return nil, &ValidationError{Field: "email", Message: "email is required"}
	}

	user := &User{
		Name:  name,
		Email: email,
	}

	if err := s.repo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}
