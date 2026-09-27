// =============================================================================
// Race Condition 테스트 예제
// =============================================================================
// 레이스 컨디션을 감지하고 검증하는 테스트 코드
//
// 실행 방법:
// cd examples/05_test
//
// 1. 레이스 감지 모드로 테스트 (핵심!)
// go test -race -v -run TestUnsafe    # 레이스가 감지되어 테스트 실패
// go test -race -v -run TestMutex     # 레이스 없음, 테스트 통과
// go test -race -v -run TestAtomic    # 레이스 없음, 테스트 통과
// go test -race -v                    # 전체 테스트
//
// 2. 레이스 감지 없이 테스트 (레이스가 있어도 통과할 수 있음!)
// go test -v -run TestUnsafe          # 레이스가 있지만 통과할 수도 있음
//
// 주의: -race 없이 테스트하면 레이스 컨디션이 있어도
//       우연히 통과할 수 있으므로 항상 -race와 함께 테스트하세요!
// =============================================================================

package testexample

import (
	"sync"
	"testing"
)

// =============================================================================
// 1. UnsafeCounter 테스트 (레이스 발생)
// =============================================================================
// go test -race -v -run TestUnsafeCounter
// 위 명령으로 실행하면 레이스 컨디션이 감지됨

func TestUnsafeCounter_Race(t *testing.T) {
	counter := &UnsafeCounter{}

	var wg sync.WaitGroup
	iterations := 1000
	goroutines := 10

	// 여러 고루틴에서 동시에 Increment 호출
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				counter.Increment() // 레이스 발생!
			}
		}()
	}

	wg.Wait()

	// 예상값: goroutines * iterations = 10000
	// 실제값: 10000보다 작을 수 있음 (레이스 때문)
	expected := goroutines * iterations
	actual := counter.Value()

	// 주의: -race 플래그 없이 실행하면 이 테스트가 통과할 수도 있음
	// 하지만 -race와 함께 실행하면 레이스가 감지됨
	t.Logf("Expected: %d, Actual: %d", expected, actual)

	// 레이스가 발생하면 값이 다를 수 있음을 보여주기 위한 로그
	if actual != expected {
		t.Logf("Race condition detected! Values don't match.")
	}
}

// =============================================================================
// 2. MutexCounter 테스트 (레이스 없음)
// =============================================================================
// go test -race -v -run TestMutexCounter
// 레이스 감지 모드에서도 통과

func TestMutexCounter_Safe(t *testing.T) {
	counter := &MutexCounter{}

	var wg sync.WaitGroup
	iterations := 1000
	goroutines := 10

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				counter.Increment() // Mutex로 보호됨
			}
		}()
	}

	wg.Wait()

	expected := goroutines * iterations
	actual := counter.Value()

	if actual != expected {
		t.Errorf("MutexCounter: expected %d, got %d", expected, actual)
	}
}

// =============================================================================
// 3. AtomicCounter 테스트 (레이스 없음)
// =============================================================================
// go test -race -v -run TestAtomicCounter

func TestAtomicCounter_Safe(t *testing.T) {
	counter := &AtomicCounter{}

	var wg sync.WaitGroup
	iterations := 1000
	goroutines := 10

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				counter.Increment() // atomic 연산
			}
		}()
	}

	wg.Wait()

	expected := int64(goroutines * iterations)
	actual := counter.Value()

	if actual != expected {
		t.Errorf("AtomicCounter: expected %d, got %d", expected, actual)
	}
}

// =============================================================================
// 4. UnsafeMap 테스트 (레이스 발생 - 패닉 가능)
// =============================================================================
// go test -race -v -run TestUnsafeMap
// 레이스 감지되며, 때로는 패닉 발생

func TestUnsafeMap_Race(t *testing.T) {
	m := NewUnsafeMap()

	var wg sync.WaitGroup

	// 동시에 읽기/쓰기
	for i := 0; i < 100; i++ {
		wg.Add(2)

		// 쓰기 고루틴
		go func(id int) {
			defer wg.Done()
			key := "key"
			m.Set(key, id) // 레이스 발생!
		}(i)

		// 읽기 고루틴
		go func(id int) {
			defer wg.Done()
			key := "key"
			m.Get(key) // 레이스 발생!
		}(i)
	}

	wg.Wait()
	t.Log("UnsafeMap test completed (may have race condition)")
}

// =============================================================================
// 5. SafeMapInt 테스트 (레이스 없음)
// =============================================================================
// go test -race -v -run TestSafeMapInt

func TestSafeMapInt_Safe(t *testing.T) {
	m := NewSafeMapInt()

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)

		go func(id int) {
			defer wg.Done()
			m.Set("key", id) // RWMutex로 보호됨
		}(i)

		go func(id int) {
			defer wg.Done()
			m.Get("key") // RWMutex로 보호됨
		}(i)
	}

	wg.Wait()
	t.Log("SafeMapInt test completed successfully")
}

// =============================================================================
// 6. sync.Map 테스트 (레이스 없음)
// =============================================================================
// go test -race -v -run TestSyncMap

func TestSyncMap_Safe(t *testing.T) {
	m := &SyncMapWrapper{}

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)

		go func(id int) {
			defer wg.Done()
			m.Set("key", id) // sync.Map은 동시성 안전
		}(i)

		go func(id int) {
			defer wg.Done()
			m.Get("key") // sync.Map은 동시성 안전
		}(i)
	}

	wg.Wait()
	t.Log("SyncMap test completed successfully")
}

// =============================================================================
// 7. ChannelCounter 테스트 (레이스 없음)
// =============================================================================
// go test -race -v -run TestChannelCounter

func TestChannelCounter_Safe(t *testing.T) {
	counter := NewChannelCounter()
	defer counter.Close()

	var wg sync.WaitGroup
	iterations := 100
	goroutines := 10

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				counter.Increment() // 채널을 통한 순차 처리
			}
		}()
	}

	wg.Wait()

	expected := goroutines * iterations
	actual := counter.Value()

	if actual != expected {
		t.Errorf("ChannelCounter: expected %d, got %d", expected, actual)
	}
}

// =============================================================================
// 8. UnsafeBankAccount 테스트 (레이스 발생)
// =============================================================================
// go test -race -v -run TestUnsafeBankAccount

func TestUnsafeBankAccount_Race(t *testing.T) {
	account := &UnsafeBankAccount{balance: 10000}

	var wg sync.WaitGroup

	// 동시에 입금/출금
	for i := 0; i < 100; i++ {
		wg.Add(2)

		go func() {
			defer wg.Done()
			account.Deposit(100) // 레이스 발생!
		}()

		go func() {
			defer wg.Done()
			account.Withdraw(100) // 레이스 발생!
		}()
	}

	wg.Wait()

	// 100번 입금(+10000), 100번 출금(-10000) = 원래 잔액(10000)이어야 함
	// 하지만 레이스로 인해 다를 수 있음
	t.Logf("Final balance: %d (expected around 10000)", account.Balance())
}

// =============================================================================
// 9. SafeBankAccount 테스트 (레이스 없음)
// =============================================================================
// go test -race -v -run TestSafeBankAccount

func TestSafeBankAccount_Safe(t *testing.T) {
	account := &SafeBankAccount{balance: 10000}

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)

		go func() {
			defer wg.Done()
			account.Deposit(100) // Mutex로 보호됨
		}()

		go func() {
			defer wg.Done()
			account.Withdraw(100) // Mutex로 보호됨
		}()
	}

	wg.Wait()

	// Mutex로 보호되므로 정확한 결과
	expected := 10000 // 100번 입금, 100번 출금 = 원래 잔액
	actual := account.Balance()

	if actual != expected {
		t.Errorf("SafeBankAccount: expected %d, got %d", expected, actual)
	}
}

// =============================================================================
// 10. 벤치마크: 동기화 방식별 성능 비교
// =============================================================================
// go test -bench=BenchmarkCounter -benchmem

func BenchmarkUnsafeCounter(b *testing.B) {
	counter := &UnsafeCounter{}
	for i := 0; i < b.N; i++ {
		counter.Increment()
	}
}

func BenchmarkMutexCounter(b *testing.B) {
	counter := &MutexCounter{}
	for i := 0; i < b.N; i++ {
		counter.Increment()
	}
}

func BenchmarkAtomicCounter(b *testing.B) {
	counter := &AtomicCounter{}
	for i := 0; i < b.N; i++ {
		counter.Increment()
	}
}

// 병렬 벤치마크 (실제 동시성 상황)
func BenchmarkMutexCounter_Parallel(b *testing.B) {
	counter := &MutexCounter{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Increment()
		}
	})
}

func BenchmarkAtomicCounter_Parallel(b *testing.B) {
	counter := &AtomicCounter{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Increment()
		}
	})
}

// =============================================================================
// 11. 테스트 요약 출력
// =============================================================================

func TestRaceCondition_Summary(t *testing.T) {
	t.Log(`
========================================
Race Condition 테스트 명령어 요약
========================================

1. 레이스 감지 모드로 전체 테스트:
   go test -race -v

2. 특정 테스트만 레이스 감지:
   go test -race -v -run TestUnsafe

3. 벤치마크 실행:
   go test -bench=BenchmarkCounter -benchmem

4. 커버리지 확인:
   go test -cover -race

========================================
결과 해석
========================================
- "WARNING: DATA RACE" 메시지가 나오면 레이스 발생
- Unsafe* 테스트는 레이스가 감지되어야 정상
- Safe*/Mutex*/Atomic* 테스트는 레이스 없이 통과해야 정상

========================================
	`)
}
