// =============================================================================
// Go 테스트 예제 - 다양한 테스팅 기법
// =============================================================================
// 실행 방법:
// cd examples/05_test
// go test -v                           # 전체 테스트 실행
// go test -v -run TestAdd              # 특정 테스트만 실행
// go test -v -run TestDivide/          # 서브테스트 실행
// go test -cover                       # 커버리지 확인
// go test -coverprofile=coverage.out   # 커버리지 파일 생성
// go tool cover -html=coverage.out     # 커버리지 HTML 리포트
// go test -bench=.                     # 벤치마크 실행
// go test -bench=. -benchmem           # 벤치마크 + 메모리 할당 정보
// =============================================================================

package testexample

import (
	"errors"
	"fmt"
	"testing"
)

// =============================================================================
// 1. 기본 테스트 (Basic Test)
// =============================================================================

// TestAdd 기본적인 테스트 함수
func TestAdd(t *testing.T) {
	// 준비 (Arrange)
	a, b := 2, 3
	expected := 5

	// 실행 (Act)
	result := Add(a, b)

	// 검증 (Assert)
	if result != expected {
		// t.Errorf: 테스트 실패 기록, 계속 진행
		t.Errorf("Add(%d, %d) = %d; expected %d", a, b, result, expected)
	}
}

// TestSubtract 빼기 테스트
func TestSubtract(t *testing.T) {
	result := Subtract(10, 3)
	if result != 7 {
		t.Errorf("Subtract(10, 3) = %d; expected 7", result)
	}
}

// =============================================================================
// 2. 테이블 드리븐 테스트 (Table-Driven Tests)
// =============================================================================
// Go에서 가장 권장되는 테스트 패턴
// - 여러 테스트 케이스를 구조화
// - 중복 코드 최소화
// - 새 케이스 추가 용이

// TestMultiply_TableDriven 테이블 드리븐 테스트 예제
func TestMultiply_TableDriven(t *testing.T) {
	// 테스트 케이스 정의
	testCases := []struct {
		name     string // 테스트 케이스 이름
		a, b     int    // 입력값
		expected int    // 기대값
	}{
		{"positive numbers", 2, 3, 6},
		{"with zero", 5, 0, 0},
		{"negative numbers", -2, 3, -6},
		{"both negative", -2, -3, 6},
		{"large numbers", 100, 200, 20000},
	}

	// 각 테스트 케이스 실행
	for _, tc := range testCases {
		// t.Run으로 서브테스트 생성
		t.Run(tc.name, func(t *testing.T) {
			result := Multiply(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("Multiply(%d, %d) = %d; expected %d",
					tc.a, tc.b, result, tc.expected)
			}
		})
	}
}

// TestDivide_TableDriven 에러를 포함한 테이블 드리븐 테스트
func TestDivide_TableDriven(t *testing.T) {
	testCases := []struct {
		name        string
		a, b        int
		expected    int
		expectError bool
	}{
		{"normal division", 10, 2, 5, false},
		{"division by zero", 10, 0, 0, true},
		{"negative dividend", -10, 2, -5, false},
		{"integer division", 7, 2, 3, false}, // 소수점 버림
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Divide(tc.a, tc.b)

			// 에러 예상 여부 확인
			if tc.expectError {
				if err == nil {
					t.Errorf("Divide(%d, %d) expected error, got nil",
						tc.a, tc.b)
				}
				return // 에러 케이스는 여기서 종료
			}

			// 에러가 없어야 하는 케이스
			if err != nil {
				t.Errorf("Divide(%d, %d) unexpected error: %v",
					tc.a, tc.b, err)
				return
			}

			if result != tc.expected {
				t.Errorf("Divide(%d, %d) = %d; expected %d",
					tc.a, tc.b, result, tc.expected)
			}
		})
	}
}

// =============================================================================
// 3. errors.Is 테스트
// =============================================================================

// TestFindUser_ErrorsIs errors.Is를 사용한 에러 검증
func TestFindUser_ErrorsIs(t *testing.T) {
	testCases := []struct {
		name        string
		userID      int
		expectError error // 예상되는 센티널 에러
	}{
		{"valid user", 1, nil},
		{"not found", 999, ErrNotFound},
		{"invalid input", -1, ErrInvalidInput},
		{"zero id", 0, ErrInvalidInput},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := FindUser(tc.userID)

			if tc.expectError == nil {
				// 에러가 없어야 하는 케이스
				if err != nil {
					t.Errorf("FindUser(%d) unexpected error: %v",
						tc.userID, err)
				}
				return
			}

			// errors.Is로 에러 확인 (래핑된 에러도 감지)
			if !errors.Is(err, tc.expectError) {
				t.Errorf("FindUser(%d) error = %v; expected %v (using errors.Is)",
					tc.userID, err, tc.expectError)
			}
		})
	}
}

// TestFindUser_WrappedError 래핑된 에러 테스트
func TestFindUser_WrappedError(t *testing.T) {
	_, err := FindUser(999) // 존재하지 않는 사용자

	// 에러가 래핑되어 있어도 errors.Is가 내부 에러를 찾아냄
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected ErrNotFound in wrapped error, got: %v", err)
	}

	// 에러 메시지에 컨텍스트 정보 포함 확인
	expectedMsg := "user id 999"
	if err == nil || !containsSubstring(err.Error(), expectedMsg) {
		t.Errorf("Error message should contain '%s', got: %v",
			expectedMsg, err)
	}
}

// =============================================================================
// 4. errors.As 테스트
// =============================================================================

// TestValidateEmail_ErrorsAs errors.As를 사용한 커스텀 에러 타입 검증
func TestValidateEmail_ErrorsAs(t *testing.T) {
	testCases := []struct {
		name         string
		email        string
		expectError  bool
		expectField  string // ValidationError의 Field 값
	}{
		{"valid email", "test@example.com", false, ""},
		{"empty email", "", true, "email"},
		{"no at sign", "invalid-email", true, "email"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEmail(tc.email)

			if !tc.expectError {
				if err != nil {
					t.Errorf("ValidateEmail(%q) unexpected error: %v",
						tc.email, err)
				}
				return
			}

			// errors.As로 ValidationError 타입 추출
			var validationErr *ValidationError
			if !errors.As(err, &validationErr) {
				t.Errorf("ValidateEmail(%q) expected ValidationError, got: %T",
					tc.email, err)
				return
			}

			// 추출한 에러의 필드 확인
			if validationErr.Field != tc.expectField {
				t.Errorf("ValidationError.Field = %q; expected %q",
					validationErr.Field, tc.expectField)
			}
		})
	}
}

// TestSaveUser_ErrorsAs 중첩된 커스텀 에러 테스트
func TestSaveUser_ErrorsAs(t *testing.T) {
	// DB 에러 시뮬레이션
	err := SaveUser("test", true)

	// DatabaseError 타입 추출
	var dbErr *DatabaseError
	if !errors.As(err, &dbErr) {
		t.Fatalf("Expected DatabaseError, got: %T", err)
	}

	// DatabaseError의 필드 확인
	if dbErr.Operation != "INSERT" {
		t.Errorf("DatabaseError.Operation = %q; expected 'INSERT'",
			dbErr.Operation)
	}

	// 내부 에러도 errors.Is로 확인 가능
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("Expected ErrInvalidInput in DatabaseError.Err")
	}
}

// =============================================================================
// 5. 서브테스트와 t.Parallel()
// =============================================================================

// TestContains_Parallel 병렬 서브테스트
func TestContains_Parallel(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{"found at start", []string{"a", "b", "c"}, "a", true},
		{"found at end", []string{"a", "b", "c"}, "c", true},
		{"not found", []string{"a", "b", "c"}, "d", false},
		{"empty slice", []string{}, "a", false},
	}

	for _, tc := range testCases {
		tc := tc // 루프 변수 캡처 (Go 1.22 이전 버전 필요)

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // 이 서브테스트를 병렬로 실행

			result := Contains(tc.slice, tc.item)
			if result != tc.expected {
				t.Errorf("Contains(%v, %q) = %v; expected %v",
					tc.slice, tc.item, result, tc.expected)
			}
		})
	}
}

// =============================================================================
// 6. 목(Mock)을 사용한 테스트
// =============================================================================

// MockUserRepository UserRepository의 목 구현
type MockUserRepository struct {
	// 반환값 설정
	GetByIDFunc func(id int) (*User, error)
	SaveFunc    func(user *User) error
	DeleteFunc  func(id int) error

	// 호출 기록
	GetByIDCalls []int
	SaveCalls    []*User
	DeleteCalls  []int
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
	m.GetByIDCalls = append(m.GetByIDCalls, id)
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockUserRepository) Save(user *User) error {
	m.SaveCalls = append(m.SaveCalls, user)
	if m.SaveFunc != nil {
		return m.SaveFunc(user)
	}
	return nil
}

func (m *MockUserRepository) Delete(id int) error {
	m.DeleteCalls = append(m.DeleteCalls, id)
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// TestUserService_GetUser 목을 사용한 서비스 테스트
func TestUserService_GetUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// 목 설정
		mockRepo := &MockUserRepository{
			GetByIDFunc: func(id int) (*User, error) {
				return &User{ID: id, Name: "Test User"}, nil
			},
		}

		// 서비스 생성 (의존성 주입)
		service := NewUserService(mockRepo)

		// 테스트 실행
		user, err := service.GetUser(1)

		// 검증
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.Name != "Test User" {
			t.Errorf("user.Name = %q; expected 'Test User'", user.Name)
		}

		// 목 호출 확인
		if len(mockRepo.GetByIDCalls) != 1 {
			t.Errorf("GetByID called %d times; expected 1",
				len(mockRepo.GetByIDCalls))
		}
		if mockRepo.GetByIDCalls[0] != 1 {
			t.Errorf("GetByID called with %d; expected 1",
				mockRepo.GetByIDCalls[0])
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		service := NewUserService(mockRepo)

		_, err := service.GetUser(0)

		// ErrInvalidInput 확인
		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got: %v", err)
		}

		// 잘못된 입력이면 Repository 호출 안 함
		if len(mockRepo.GetByIDCalls) != 0 {
			t.Errorf("GetByID should not be called for invalid id")
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			GetByIDFunc: func(id int) (*User, error) {
				return nil, ErrNotFound
			},
		}
		service := NewUserService(mockRepo)

		_, err := service.GetUser(999)

		if !errors.Is(err, ErrNotFound) {
			t.Errorf("expected ErrNotFound, got: %v", err)
		}
	})
}

// TestUserService_CreateUser 생성 테스트
func TestUserService_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		service := NewUserService(mockRepo)

		user, err := service.CreateUser("Alice", "alice@example.com")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.Name != "Alice" {
			t.Errorf("user.Name = %q; expected 'Alice'", user.Name)
		}

		// Save 호출 확인
		if len(mockRepo.SaveCalls) != 1 {
			t.Errorf("Save called %d times; expected 1", len(mockRepo.SaveCalls))
		}
	})

	t.Run("empty name", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		service := NewUserService(mockRepo)

		_, err := service.CreateUser("", "test@example.com")

		var validationErr *ValidationError
		if !errors.As(err, &validationErr) {
			t.Fatalf("expected ValidationError, got: %T", err)
		}
		if validationErr.Field != "name" {
			t.Errorf("ValidationError.Field = %q; expected 'name'",
				validationErr.Field)
		}
	})

	t.Run("save error", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			SaveFunc: func(user *User) error {
				return &DatabaseError{Operation: "INSERT", Err: ErrInvalidInput}
			},
		}
		service := NewUserService(mockRepo)

		_, err := service.CreateUser("Alice", "alice@example.com")

		// 에러 체인 확인
		var dbErr *DatabaseError
		if !errors.As(err, &dbErr) {
			t.Errorf("expected DatabaseError in error chain, got: %v", err)
		}
	})
}

// =============================================================================
// 7. 테스트 헬퍼 함수
// =============================================================================

// TestHelper 테스트 헬퍼 사용 예제
func TestHelper(t *testing.T) {
	// assertEqual 헬퍼 함수 사용
	result := Add(2, 3)
	assertEqual(t, result, 5, "Add(2, 3)")
}

// assertEqual 값 비교 헬퍼
func assertEqual(t *testing.T, got, want interface{}, msg string) {
	t.Helper() // 이 함수를 헬퍼로 표시 (에러 위치가 호출자로 표시됨)

	if got != want {
		t.Errorf("%s: got %v, want %v", msg, got, want)
	}
}

// assertError 에러 확인 헬퍼
func assertError(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Errorf("expected error %v, got %v", target, err)
	}
}

// assertNoError 에러 없음 확인 헬퍼
func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// =============================================================================
// 8. 벤치마크 테스트
// =============================================================================

// BenchmarkAdd 덧셈 벤치마크
func BenchmarkAdd(b *testing.B) {
	// b.N은 벤치마크 프레임워크가 자동으로 조절
	for i := 0; i < b.N; i++ {
		Add(100, 200)
	}
}

// BenchmarkContains_Found Contains 벤치마크 (찾는 경우)
func BenchmarkContains_Found(b *testing.B) {
	slice := []string{"apple", "banana", "cherry", "date", "elderberry"}

	b.ResetTimer() // 타이머 리셋 (준비 시간 제외)

	for i := 0; i < b.N; i++ {
		Contains(slice, "cherry")
	}
}

// BenchmarkContains_NotFound Contains 벤치마크 (못 찾는 경우)
func BenchmarkContains_NotFound(b *testing.B) {
	slice := []string{"apple", "banana", "cherry", "date", "elderberry"}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Contains(slice, "fig")
	}
}

// BenchmarkUnique Unique 벤치마크
func BenchmarkUnique(b *testing.B) {
	slice := []string{"a", "b", "c", "a", "b", "d", "e", "c", "f", "g"}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Unique(slice)
	}
}

// =============================================================================
// 9. 예제 테스트 (Example Tests)
// =============================================================================
// Example 함수는 문서화 + 테스트 역할
// Output: 주석으로 예상 출력 명시

func ExampleAdd() {
	result := Add(2, 3)
	fmt.Println(result)
	// Output: 5
}

func ExampleDivide() {
	result, err := Divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(result)
	// Output: 5
}

func ExampleContains() {
	fruits := []string{"apple", "banana", "cherry"}
	fmt.Println(Contains(fruits, "banana"))
	fmt.Println(Contains(fruits, "grape"))
	// Output:
	// true
	// false
}

// =============================================================================
// 유틸리티 함수
// =============================================================================

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
