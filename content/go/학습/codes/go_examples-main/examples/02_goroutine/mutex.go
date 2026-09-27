// =============================================================================
// Mutex (상호 배제) 예제
// =============================================================================
// Mutex는 공유 자원에 대한 동시 접근을 제어하는 동기화 도구입니다.
//
// sync.Mutex 주요 메서드:
// - Lock(): 락 획득 (다른 고루틴이 이미 락을 가지고 있으면 대기)
// - Unlock(): 락 해제 (반드시 Lock 후에 호출)
//
// sync.RWMutex (읽기/쓰기 뮤텍스):
// - Lock() / Unlock(): 쓰기 락 (배타적)
// - RLock() / RUnlock(): 읽기 락 (공유 가능)
// - 읽기가 많고 쓰기가 적은 경우 성능 향상
//
// 주의사항:
// - 락을 획득한 고루틴만 Unlock 가능
// - 데드락 주의: 락 획득 순서 통일, defer 사용 권장
// - 락을 최소한의 범위에서만 사용
// =============================================================================

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Race Condition 문제 예제 ===")
	raceConditionProblem()

	fmt.Println("\n=== Mutex로 해결 ===")
	mutexSolution()

	fmt.Println("\n=== RWMutex 예제 ===")
	rwMutexExample()

	fmt.Println("\n=== 실전 예제: 안전한 맵 ===")
	safeMapExample()
}

// raceConditionProblem Race Condition 문제 시연
// 여러 고루틴이 동시에 같은 변수를 수정하면 예상과 다른 결과가 발생
func raceConditionProblem() {
	counter := 0
	iterations := 1000

	var wg sync.WaitGroup

	// 10개의 고루틴이 각각 counter를 1000번 증가
	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < iterations; j++ {
				// 이 연산은 원자적이지 않음!
				// 1. counter 값 읽기
				// 2. 1 더하기
				// 3. 결과를 counter에 쓰기
				// 위 세 단계 사이에 다른 고루틴이 끼어들 수 있음
				counter++
			}
		}()
	}

	wg.Wait()

	// 예상값: 10 * 1000 = 10000
	// 실제값: 10000보다 작을 수 있음 (Race Condition)
	fmt.Printf("예상값: %d, 실제값: %d\n", 10*iterations, counter)
	fmt.Println("(실제값이 예상값보다 작으면 Race Condition 발생)")
}

// mutexSolution Mutex를 사용하여 Race Condition 해결
func mutexSolution() {
	counter := 0
	iterations := 1000

	// 뮤텍스 생성
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < iterations; j++ {
				// ================================================================
				// 임계 영역(Critical Section) 시작
				// ================================================================
				// Lock()을 호출하면:
				// 1. 아무도 락을 가지고 있지 않으면 → 락 획득, 계속 진행
				// 2. 다른 고루틴이 락을 가지고 있으면 → 대기
				mu.Lock()

				// 이 영역에서는 오직 하나의 고루틴만 실행됨
				counter++

				// 락 해제 - 반드시 Lock 후에 호출해야 함
				mu.Unlock()
				// ================================================================
				// 임계 영역 끝
				// ================================================================
			}
		}()
	}

	wg.Wait()

	// 이제 항상 정확한 값이 나옴
	fmt.Printf("예상값: %d, 실제값: %d\n", 10*iterations, counter)
}

// rwMutexExample RWMutex 사용 예제
// 읽기 작업이 많고 쓰기 작업이 적을 때 유용
func rwMutexExample() {
	// 공유 데이터
	data := make(map[string]int)
	data["count"] = 0

	// RWMutex 생성
	var rwmu sync.RWMutex
	var wg sync.WaitGroup

	// ==========================================================================
	// 쓰기 작업 (배타적 락)
	// ==========================================================================
	// 쓰기 작업은 5개
	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			// 쓰기 락 획득
			// Lock()은 모든 읽기/쓰기 락이 해제될 때까지 대기
			rwmu.Lock()
			defer rwmu.Unlock()

			// 임계 영역 - 쓰기
			oldValue := data["count"]
			time.Sleep(10 * time.Millisecond) // 쓰기 작업 시뮬레이션
			data["count"] = oldValue + 10

			fmt.Printf("Writer %d: count를 %d에서 %d로 변경\n",
				id, oldValue, data["count"])
		}(i + 1)
	}

	// ==========================================================================
	// 읽기 작업 (공유 락)
	// ==========================================================================
	// 읽기 작업은 20개 (쓰기보다 훨씬 많음)
	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			// 읽기 락 획득
			// RLock()은 쓰기 락이 없으면 즉시 획득
			// 여러 고루틴이 동시에 읽기 락을 획득할 수 있음
			rwmu.RLock()
			defer rwmu.RUnlock()

			// 임계 영역 - 읽기
			value := data["count"]
			time.Sleep(5 * time.Millisecond) // 읽기 작업 시뮬레이션

			fmt.Printf("Reader %d: count = %d\n", id, value)
		}(i + 1)
	}

	wg.Wait()
	fmt.Printf("최종 count 값: %d\n", data["count"])
}

// SafeMap 동시성 안전한 맵 구현
// Go의 내장 map은 동시성 안전하지 않으므로 래핑 필요
type SafeMap struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

// NewSafeMap SafeMap 생성자
func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]interface{}),
	}
}

// Set 값 설정 (쓰기 락 사용)
func (sm *SafeMap) Set(key string, value interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.data[key] = value
}

// Get 값 조회 (읽기 락 사용)
func (sm *SafeMap) Get(key string) (interface{}, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	value, exists := sm.data[key]
	return value, exists
}

// Delete 값 삭제 (쓰기 락 사용)
func (sm *SafeMap) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.data, key)
}

// Len 맵 길이 조회 (읽기 락 사용)
func (sm *SafeMap) Len() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return len(sm.data)
}

// Keys 모든 키 조회 (읽기 락 사용)
func (sm *SafeMap) Keys() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	keys := make([]string, 0, len(sm.data))
	for key := range sm.data {
		keys = append(keys, key)
	}
	return keys
}

// safeMapExample SafeMap 사용 예제
func safeMapExample() {
	sm := NewSafeMap()

	var wg sync.WaitGroup

	// 여러 고루틴에서 동시에 읽기/쓰기
	for i := 0; i < 100; i++ {
		wg.Add(2)

		// 쓰기 고루틴
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", id%10)
			sm.Set(key, id)
		}(i)

		// 읽기 고루틴
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", id%10)
			sm.Get(key)
		}(i)
	}

	wg.Wait()

	fmt.Printf("맵에 저장된 키 개수: %d\n", sm.Len())
	fmt.Printf("저장된 키들: %v\n", sm.Keys())

	// 특정 값 조회
	if value, exists := sm.Get("key_5"); exists {
		fmt.Printf("key_5의 값: %v\n", value)
	}
}
