package hard

/*
================================================================================
문제 2: Circuit Breaker (서킷 브레이커) 패턴 구현
================================================================================

난이도: Hard
주제: 분산 시스템, 장애 복구, 상태 머신

【 문제 설명 】
외부 서비스 호출의 장애 전파를 방지하는 Circuit Breaker를 구현하세요.

상태:
- Closed (닫힘): 정상 동작, 요청 통과
- Open (열림): 장애 감지, 요청 즉시 실패
- Half-Open (반열림): 복구 테스트 중

【 동작 원리 】
Closed:
  - 요청 정상 통과
  - 연속 실패 N회 → Open 전환

Open:
  - 모든 요청 즉시 실패 (fallback)
  - 타임아웃 후 → Half-Open 전환

Half-Open:
  - 제한된 요청만 통과
  - 성공 → Closed 전환
  - 실패 → Open 전환

【 인터페이스 】
type CircuitBreaker interface {
    Execute(fn func() error) error
    State() State
    Reset()
}

【 설정 】
- FailureThreshold: Open 전환 실패 횟수
- SuccessThreshold: Closed 전환 성공 횟수
- Timeout: Open 상태 유지 시간
- HalfOpenMaxCalls: Half-Open 시 허용 요청 수

【 예시 】
cb := NewCircuitBreaker(Config{
    FailureThreshold:  5,
    SuccessThreshold:  2,
    Timeout:           10 * time.Second,
    HalfOpenMaxCalls:  3,
})

err := cb.Execute(func() error {
    return callExternalService()
})

if errors.Is(err, ErrCircuitOpen) {
    // 서킷 열림, fallback 처리
}

【 실무 연관성 】
- 마이크로서비스 간 통신
- 외부 API 호출
- 데이터베이스 연결
================================================================================
*/

import (
	"errors"
	"time"
)

// State 서킷 상태
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

var (
	ErrCircuitOpen    = errors.New("circuit breaker is open")
	ErrTooManyCalls   = errors.New("too many calls in half-open state")
)

// Config 서킷 브레이커 설정
type Config struct {
	FailureThreshold  int           // Open 전환 실패 횟수
	SuccessThreshold  int           // Closed 전환 성공 횟수
	Timeout           time.Duration // Open 상태 유지 시간
	HalfOpenMaxCalls  int           // Half-Open 시 허용 요청 수
}

// CircuitBreaker 서킷 브레이커 구조체
type CircuitBreaker struct {
	// 여기에 필드를 정의하세요
}

// NewCircuitBreaker 생성자
func NewCircuitBreaker(config Config) *CircuitBreaker {
	// 여기에 코드를 작성하세요
	return nil
}

// Execute 함수 실행
func (cb *CircuitBreaker) Execute(fn func() error) error {
	// 여기에 코드를 작성하세요
	return nil
}

// State 현재 상태 반환
func (cb *CircuitBreaker) GetState() State {
	// 여기에 코드를 작성하세요
	return StateClosed
}

// Reset 상태 초기화
func (cb *CircuitBreaker) Reset() {
	// 여기에 코드를 작성하세요
}

// 테스트
func main() {
	cb := NewCircuitBreaker(Config{
		FailureThreshold:  3,
		SuccessThreshold:  2,
		Timeout:           5 * time.Second,
		HalfOpenMaxCalls:  2,
	})

	// 연속 실패 시 서킷 열림
	for i := 0; i < 5; i++ {
		err := cb.Execute(func() error {
			return errors.New("service unavailable")
		})
		println("Call", i+1, "error:", err.Error())
	}

	println("Circuit state:", cb.GetState())
}
