// =============================================================================
// Race Condition 테스트 대상 코드
// =============================================================================
// 이 파일은 레이스 컨디션이 발생하는 코드와 해결된 코드를 포함합니다.
// 테스트 코드는 race_condition_test.go에 있습니다.
//
// Race Condition 감지 방법:
// go test -race -v                    # 레이스 감지 모드로 테스트
// go run -race race_condition.go      # 레이스 감지 모드로 실행
// go build -race                      # 레이스 감지 바이너리 빌드
//
// 주의: -race 플래그는 성능이 2~10배 느려지므로 테스트에서만 사용
// =============================================================================

package testexample

import (
	"sync"
	"sync/atomic"
)

// =============================================================================
// 1. Race Condition이 발생하는 코드 (잘못된 예)
// =============================================================================

// UnsafeCounter 안전하지 않은 카운터 (레이스 컨디션 발생)
type UnsafeCounter struct {
	count int
}

// Increment 안전하지 않은 증가 (동시 호출 시 레이스 발생)
func (c *UnsafeCounter) Increment() {
	// 이 연산은 원자적이지 않음!
	// 1. count 값 읽기
	// 2. 1 더하기
	// 3. 결과를 count에 쓰기
	c.count++
}

// Value 현재 값 반환
func (c *UnsafeCounter) Value() int {
	return c.count
}

// =============================================================================
// 2. Mutex를 사용한 해결
// =============================================================================

// MutexCounter Mutex로 보호된 카운터
type MutexCounter struct {
	mu    sync.Mutex
	count int
}

// Increment Mutex로 보호된 증가
func (c *MutexCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

// Value Mutex로 보호된 읽기
func (c *MutexCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// =============================================================================
// 3. atomic 패키지를 사용한 해결
// =============================================================================

// AtomicCounter atomic 연산을 사용한 카운터
type AtomicCounter struct {
	count int64 // atomic은 int64/uint64/int32/uint32 사용
}

// Increment 원자적 증가
func (c *AtomicCounter) Increment() {
	atomic.AddInt64(&c.count, 1)
}

// Value 원자적 읽기
func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.count)
}

// =============================================================================
// 4. Race Condition이 발생하는 맵 (잘못된 예)
// =============================================================================

// UnsafeMap 안전하지 않은 맵 (동시 접근 시 패닉 발생 가능)
type UnsafeMap struct {
	data map[string]int
}

// NewUnsafeMap 생성자
func NewUnsafeMap() *UnsafeMap {
	return &UnsafeMap{
		data: make(map[string]int),
	}
}

// Set 안전하지 않은 쓰기
func (m *UnsafeMap) Set(key string, value int) {
	m.data[key] = value
}

// Get 안전하지 않은 읽기
func (m *UnsafeMap) Get(key string) (int, bool) {
	v, ok := m.data[key]
	return v, ok
}

// =============================================================================
// 5. RWMutex를 사용한 안전한 맵
// =============================================================================

// SafeMap RWMutex로 보호된 맵
type SafeMapInt struct {
	mu   sync.RWMutex
	data map[string]int
}

// NewSafeMapInt 생성자
func NewSafeMapInt() *SafeMapInt {
	return &SafeMapInt{
		data: make(map[string]int),
	}
}

// Set 쓰기 락으로 보호된 쓰기
func (m *SafeMapInt) Set(key string, value int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// Get 읽기 락으로 보호된 읽기
func (m *SafeMapInt) Get(key string) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}

// =============================================================================
// 6. sync.Map 사용 (Go 내장 동시성 안전 맵)
// =============================================================================

// SyncMapWrapper sync.Map 래퍼
type SyncMapWrapper struct {
	data sync.Map
}

// Set sync.Map 쓰기
func (m *SyncMapWrapper) Set(key string, value int) {
	m.data.Store(key, value)
}

// Get sync.Map 읽기
func (m *SyncMapWrapper) Get(key string) (int, bool) {
	v, ok := m.data.Load(key)
	if !ok {
		return 0, false
	}
	return v.(int), true
}

// =============================================================================
// 7. 채널을 사용한 Race Condition 방지
// =============================================================================

// ChannelCounter 채널을 사용한 카운터
type ChannelCounter struct {
	incrementCh chan struct{}
	valueCh     chan int
	done        chan struct{}
}

// NewChannelCounter 채널 카운터 생성
func NewChannelCounter() *ChannelCounter {
	c := &ChannelCounter{
		incrementCh: make(chan struct{}),
		valueCh:     make(chan int),
		done:        make(chan struct{}),
	}

	// 백그라운드 고루틴이 모든 연산을 순차 처리
	go c.run()

	return c
}

// run 백그라운드 처리 고루틴
func (c *ChannelCounter) run() {
	count := 0
	for {
		select {
		case <-c.incrementCh:
			count++
		case c.valueCh <- count:
			// 현재 값 전송
		case <-c.done:
			return
		}
	}
}

// Increment 증가 요청
func (c *ChannelCounter) Increment() {
	c.incrementCh <- struct{}{}
}

// Value 현재 값 요청
func (c *ChannelCounter) Value() int {
	return <-c.valueCh
}

// Close 카운터 종료
func (c *ChannelCounter) Close() {
	close(c.done)
}

// =============================================================================
// 8. 실전 예제: 은행 계좌 (Race Condition 발생)
// =============================================================================

// UnsafeBankAccount 안전하지 않은 은행 계좌
type UnsafeBankAccount struct {
	balance int
}

// Deposit 입금 (레이스 발생 가능)
func (a *UnsafeBankAccount) Deposit(amount int) {
	a.balance += amount
}

// Withdraw 출금 (레이스 발생 가능)
func (a *UnsafeBankAccount) Withdraw(amount int) bool {
	if a.balance >= amount {
		a.balance -= amount
		return true
	}
	return false
}

// Balance 잔액 조회
func (a *UnsafeBankAccount) Balance() int {
	return a.balance
}

// =============================================================================
// 9. 실전 예제: 안전한 은행 계좌
// =============================================================================

// SafeBankAccount Mutex로 보호된 은행 계좌
type SafeBankAccount struct {
	mu      sync.Mutex
	balance int
}

// Deposit 안전한 입금
func (a *SafeBankAccount) Deposit(amount int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.balance += amount
}

// Withdraw 안전한 출금
func (a *SafeBankAccount) Withdraw(amount int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.balance >= amount {
		a.balance -= amount
		return true
	}
	return false
}

// Balance 안전한 잔액 조회
func (a *SafeBankAccount) Balance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

// Transfer 계좌 간 이체 (데드락 방지를 위한 순서 보장)
func Transfer(from, to *SafeBankAccount, amount int) bool {
	// 데드락 방지: 항상 같은 순서로 락 획득
	// 실제 구현에서는 계좌 ID 등으로 순서 결정
	from.mu.Lock()
	defer from.mu.Unlock()

	to.mu.Lock()
	defer to.mu.Unlock()

	if from.balance >= amount {
		from.balance -= amount
		to.balance += amount
		return true
	}
	return false
}
