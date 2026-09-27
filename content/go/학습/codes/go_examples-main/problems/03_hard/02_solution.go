package hard

/*
================================================================================
문제 2: Circuit Breaker - 솔루션
================================================================================
*/

import (
	"errors"
	"sync"
	"time"
)

/*
================================================================================
Circuit Breaker 패턴
================================================================================

【 왜 필요한가? 】

문제 상황:
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│  Service A ──────▶ Service B (장애) ──────▶ Service C                       │
│      │                  ✕                                                   │
│      │                                                                      │
│      └──────▶ 요청 대기 → 타임아웃 → 리소스 고갈 → A도 장애!                │
│                                                                             │
│  연쇄 장애 (Cascading Failure)                                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

해결:
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│  Service A ──▶ [Circuit Breaker] ──▶ Service B (장애)                       │
│                      │                    ✕                                 │
│                      │                                                      │
│                      └──▶ 빠른 실패 (Fast Fail) → A 보호!                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘


【 상태 전이 다이어그램 】

┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│                         실패 N회                                            │
│            ┌────────────────────────────────┐                               │
│            │                                │                               │
│            ▼                                │                               │
│     ┌──────────┐     타임아웃 후      ┌─────┴────┐                         │
│     │   OPEN   │ ──────────────────▶ │ HALF-OPEN │                         │
│     │ (열림)   │                      │ (반열림)  │                         │
│     └──────────┘                      └─────┬────┘                         │
│            ▲                                │                               │
│            │                           성공 M회                             │
│            │    실패                        │                               │
│            └────────────────────────────────┤                               │
│                                             ▼                               │
│                                       ┌──────────┐                         │
│                                       │  CLOSED  │◀──── 초기 상태           │
│                                       │  (닫힘)  │                          │
│                                       └──────────┘                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
*/

// State 서킷 상태
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF-OPEN"
	default:
		return "UNKNOWN"
	}
}

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

// counters 카운터 구조체
type counters struct {
	failures       int // 연속 실패 횟수
	successes      int // 연속 성공 횟수
	halfOpenCalls  int // Half-Open 상태에서의 호출 수
}

// CircuitBreaker 서킷 브레이커 구현
type CircuitBreaker struct {
	config    Config
	state     State
	counters  counters
	openedAt  time.Time    // Open 상태로 전환된 시간
	mu        sync.RWMutex // 동시성 제어
}

// NewCircuitBreaker 생성자
func NewCircuitBreaker(config Config) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
	}
}

// Execute 함수 실행
/*
【 동작 순서 】
1. 현재 상태 확인
2. 상태에 따라:
   - Closed: 실행 허용
   - Open: 타임아웃 확인, 즉시 실패 또는 Half-Open 전환
   - Half-Open: 제한된 실행 허용
3. 결과에 따라 상태 업데이트
*/
func (cb *CircuitBreaker) Execute(fn func() error) error {
	// 1. 요청 가능 여부 확인
	if err := cb.beforeRequest(); err != nil {
		return err
	}

	// 2. 실제 함수 실행
	err := fn()

	// 3. 결과 처리
	cb.afterRequest(err)

	return err
}

// beforeRequest 요청 전 상태 확인
func (cb *CircuitBreaker) beforeRequest() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return nil

	case StateOpen:
		// 타임아웃 경과 확인
		if time.Since(cb.openedAt) >= cb.config.Timeout {
			cb.toHalfOpen()
			return nil
		}
		return ErrCircuitOpen

	case StateHalfOpen:
		// 허용된 호출 수 확인
		if cb.counters.halfOpenCalls >= cb.config.HalfOpenMaxCalls {
			return ErrTooManyCalls
		}
		cb.counters.halfOpenCalls++
		return nil
	}

	return nil
}

// afterRequest 요청 후 상태 업데이트
func (cb *CircuitBreaker) afterRequest(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.onFailure()
	} else {
		cb.onSuccess()
	}
}

// onSuccess 성공 처리
func (cb *CircuitBreaker) onSuccess() {
	switch cb.state {
	case StateClosed:
		cb.counters.failures = 0

	case StateHalfOpen:
		cb.counters.successes++
		if cb.counters.successes >= cb.config.SuccessThreshold {
			cb.toClosed()
		}
	}
}

// onFailure 실패 처리
func (cb *CircuitBreaker) onFailure() {
	switch cb.state {
	case StateClosed:
		cb.counters.failures++
		if cb.counters.failures >= cb.config.FailureThreshold {
			cb.toOpen()
		}

	case StateHalfOpen:
		cb.toOpen()
	}
}

// 상태 전환 메서드
func (cb *CircuitBreaker) toClosed() {
	cb.state = StateClosed
	cb.counters = counters{}
}

func (cb *CircuitBreaker) toOpen() {
	cb.state = StateOpen
	cb.openedAt = time.Now()
	cb.counters = counters{}
}

func (cb *CircuitBreaker) toHalfOpen() {
	cb.state = StateHalfOpen
	cb.counters = counters{}
}

// GetState 현재 상태 반환
func (cb *CircuitBreaker) GetState() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Reset 상태 초기화
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.toClosed()
}

/*
================================================================================
고급 기능
================================================================================

【 1. Fallback 지원 】

func (cb *CircuitBreaker) ExecuteWithFallback(
    fn func() error,
    fallback func() error,
) error {
    err := cb.Execute(fn)
    if errors.Is(err, ErrCircuitOpen) {
        return fallback()
    }
    return err
}

【 2. 메트릭 수집 】

type Metrics struct {
    TotalCalls     int64
    SuccessCalls   int64
    FailureCalls   int64
    RejectedCalls  int64
    LastError      error
    LastErrorTime  time.Time
}

【 3. 슬라이딩 윈도우 실패율 】

고정 횟수 대신 시간 윈도우 내 실패율로 판단:
- 최근 1분 내 실패율 > 50% → Open

【 4. 여러 에러 타입 구분 】

type Settings struct {
    // 서킷을 트리거하는 에러만 카운트
    IsSuccessful func(error) bool
}

cb := NewCircuitBreaker(Settings{
    IsSuccessful: func(err error) bool {
        // 4xx 에러는 서버 문제가 아니므로 성공 취급
        var httpErr *HTTPError
        if errors.As(err, &httpErr) {
            return httpErr.StatusCode < 500
        }
        return err == nil
    },
})

================================================================================
실무 라이브러리
================================================================================

1. sony/gobreaker
   - 가장 인기 있는 Go 서킷 브레이커

2. afex/hystrix-go
   - Netflix Hystrix의 Go 포트

3. resilience4j (Java)
   - 참고용, 매우 상세한 설정 가능

================================================================================
*/

// 테스트
func main() {
	cb := NewCircuitBreaker(Config{
		FailureThreshold:  3,
		SuccessThreshold:  2,
		Timeout:           2 * time.Second,
		HalfOpenMaxCalls:  2,
	})

	println("=== Circuit Breaker Test ===\n")

	// 1. Closed 상태에서 연속 실패
	println("1. Failing 3 times to open circuit:")
	for i := 0; i < 3; i++ {
		err := cb.Execute(func() error {
			return errors.New("service error")
		})
		println("  Call", i+1, "- State:", cb.GetState().String(), "Error:", err.Error())
	}

	// 2. Open 상태 확인
	println("\n2. Circuit is now OPEN:")
	err := cb.Execute(func() error { return nil })
	println("  Call rejected:", err.Error())

	// 3. 타임아웃 대기 후 Half-Open
	println("\n3. Waiting for timeout...")
	time.Sleep(2 * time.Second)

	// 4. Half-Open에서 성공
	println("\n4. Recovery in Half-Open:")
	for i := 0; i < 2; i++ {
		err := cb.Execute(func() error {
			return nil // 성공
		})
		errStr := "nil"
		if err != nil {
			errStr = err.Error()
		}
		println("  Call", i+1, "- State:", cb.GetState().String(), "Error:", errStr)
	}

	// 5. Closed로 복구
	println("\n5. Final state:", cb.GetState().String())
}
