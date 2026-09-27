package senior_topics

/*
================================================================================
Go 동시성 심화 (시니어 레벨)
================================================================================

면접 질문: "Go의 동시성 모델을 설명하고, 실제 프로덕션에서의 주의점은?"

시니어급 답변 포인트:
1. Goroutine 스케줄링 (G-M-P 모델)
2. 채널 vs Mutex 선택 기준
3. 동시성 패턴 (Fan-out/Fan-in, Worker Pool, Pipeline)
4. Race Condition 방지
5. 성능 최적화
*/

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// 1. G-M-P 스케줄링 모델
// ============================================================================

/*
Go 스케줄러의 G-M-P 모델:

G (Goroutine): 경량 스레드, 약 2KB 스택으로 시작
M (Machine): OS 스레드, 실제 실행 단위
P (Processor): 논리적 프로세서, GOMAXPROCS 개수만큼 존재

작동 원리:
1. 각 P는 Local Run Queue (LRQ) 보유
2. G는 P의 LRQ에 추가됨
3. M은 P에 연결되어 G를 실행
4. P가 없으면 M은 대기 상태
5. Work Stealing: 다른 P의 LRQ에서 G를 훔쳐옴

스케줄링 시점:
- go 키워드로 새 goroutine 생성
- 채널 송/수신 블록
- syscall 호출
- runtime.Gosched() 명시적 양보
- GC 실행
*/

func GoroutineSchedulingDemo() {
	fmt.Println("=== G-M-P 스케줄링 데모 ===")

	// GOMAXPROCS 확인 및 설정
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("NumCPU: %d\n", runtime.NumCPU())

	// Goroutine 수 확인
	fmt.Printf("현재 Goroutine 수: %d\n", runtime.NumGoroutine())

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// runtime.Gosched()로 명시적 양보
			runtime.Gosched()
			fmt.Printf("Goroutine %d 실행 중\n", id)
		}(i)
	}

	fmt.Printf("생성 후 Goroutine 수: %d\n", runtime.NumGoroutine())
	wg.Wait()
}

// ============================================================================
// 2. 채널 vs Mutex 선택 기준
// ============================================================================

/*
선택 기준:

채널을 사용해야 할 때:
- 데이터의 소유권 전달
- 작업 분배 (Fan-out/Fan-in)
- 이벤트 통지
- 동기화 포인트 (작업 완료 대기)

Mutex를 사용해야 할 때:
- 공유 상태 보호 (캐시, 맵 등)
- 간단한 카운터
- 복잡한 불변식(invariant) 유지
- 성능이 중요한 상황 (채널보다 빠름)

격언: "Don't communicate by sharing memory; share memory by communicating."
하지만 실제로는 상황에 맞게 선택해야 함
*/

// Mutex 기반 카운터 (간단한 상태)
type MutexCounter struct {
	mu    sync.Mutex
	count int
}

func (c *MutexCounter) Inc() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

func (c *MutexCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// 채널 기반 카운터 (소유권 전달)
type ChannelCounter struct {
	ch chan int
}

func NewChannelCounter() *ChannelCounter {
	c := &ChannelCounter{ch: make(chan int)}
	go c.run()
	return c
}

func (c *ChannelCounter) run() {
	count := 0
	for delta := range c.ch {
		if delta == 0 { // 값 요청
			c.ch <- count
		} else {
			count += delta
		}
	}
}

func (c *ChannelCounter) Inc() {
	c.ch <- 1
}

func (c *ChannelCounter) Value() int {
	c.ch <- 0
	return <-c.ch
}

// Atomic 카운터 (가장 빠름)
type AtomicCounter struct {
	count int64
}

func (c *AtomicCounter) Inc() {
	atomic.AddInt64(&c.count, 1)
}

func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.count)
}

func CounterComparisonDemo() {
	fmt.Println("\n=== 카운터 구현 비교 ===")

	const iterations = 100000

	// Mutex 카운터
	start := time.Now()
	mc := &MutexCounter{}
	var wg sync.WaitGroup
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			mc.Inc()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("Mutex 카운터: %d, 소요시간: %v\n", mc.Value(), time.Since(start))

	// Atomic 카운터
	start = time.Now()
	ac := &AtomicCounter{}
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			ac.Inc()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("Atomic 카운터: %d, 소요시간: %v\n", ac.Value(), time.Since(start))
}

// ============================================================================
// 3. sync.RWMutex - 읽기/쓰기 분리
// ============================================================================

/*
RWMutex 특성:
- 다수의 Reader 동시 접근 가능
- Writer는 독점적 접근
- 읽기가 많고 쓰기가 적은 경우 유리
- Writer starvation 주의 (대기 중인 Writer가 있으면 새 Reader 블록)
*/

type ConcurrentCache struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

func NewConcurrentCache() *ConcurrentCache {
	return &ConcurrentCache{
		data: make(map[string]interface{}),
	}
}

func (c *ConcurrentCache) Get(key string) (interface{}, bool) {
	c.mu.RLock() // 읽기 락
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *ConcurrentCache) Set(key string, value interface{}) {
	c.mu.Lock() // 쓰기 락
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *ConcurrentCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// ============================================================================
// 4. sync.Once - 단 한 번만 실행
// ============================================================================

/*
sync.Once 활용:
- Singleton 패턴 구현
- 지연 초기화 (Lazy Initialization)
- 설정 로딩
- 연결 풀 초기화
*/

type Database struct {
	conn string
}

var (
	dbInstance *Database
	dbOnce     sync.Once
)

func GetDatabase() *Database {
	dbOnce.Do(func() {
		fmt.Println("데이터베이스 연결 초기화 (한 번만 실행)")
		dbInstance = &Database{conn: "connected"}
	})
	return dbInstance
}

func OnceDemo() {
	fmt.Println("\n=== sync.Once 데모 ===")

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			db := GetDatabase()
			fmt.Printf("Goroutine %d: DB=%v\n", id, db.conn)
		}(i)
	}
	wg.Wait()
}

// ============================================================================
// 5. sync.Cond - 조건 변수
// ============================================================================

/*
sync.Cond 사용 시점:
- 특정 조건이 만족될 때까지 대기
- 브로드캐스트로 여러 goroutine 깨우기
- 생산자-소비자 패턴
*/

type BoundedQueue struct {
	cond  *sync.Cond
	data  []interface{}
	limit int
}

func NewBoundedQueue(limit int) *BoundedQueue {
	return &BoundedQueue{
		cond:  sync.NewCond(&sync.Mutex{}),
		data:  make([]interface{}, 0, limit),
		limit: limit,
	}
}

func (q *BoundedQueue) Enqueue(item interface{}) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	// 큐가 가득 찰 때까지 대기
	for len(q.data) >= q.limit {
		q.cond.Wait()
	}

	q.data = append(q.data, item)
	q.cond.Signal() // 대기 중인 하나를 깨움
}

func (q *BoundedQueue) Dequeue() interface{} {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	// 큐가 빌 때까지 대기
	for len(q.data) == 0 {
		q.cond.Wait()
	}

	item := q.data[0]
	q.data = q.data[1:]
	q.cond.Signal() // 대기 중인 하나를 깨움

	return item
}

// ============================================================================
// 6. Context를 활용한 동시성 제어
// ============================================================================

/*
Context 활용:
- 취소 신호 전파
- 타임아웃 설정
- 요청 범위 값 전달
- Goroutine 생명주기 관리
*/

func ContextCancellationDemo() {
	fmt.Println("\n=== Context 취소 데모 ===")

	// 3초 후 자동 취소
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// 여러 worker에 context 전달
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(ctx, id)
		}(i)
	}

	// 1초 후 수동 취소 (타임아웃 전에)
	time.Sleep(1 * time.Second)
	cancel()
	fmt.Println("취소 신호 전송")

	wg.Wait()
	fmt.Println("모든 worker 종료")
}

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: 종료 (%v)\n", id, ctx.Err())
			return
		case <-time.After(200 * time.Millisecond):
			fmt.Printf("Worker %d: 작업 중...\n", id)
		}
	}
}

// ============================================================================
// 7. 동시성 패턴: Pipeline
// ============================================================================

/*
Pipeline 패턴:
- 여러 단계를 채널로 연결
- 각 단계는 독립적으로 실행
- 데이터 스트림 처리에 적합
*/

func PipelineDemo() {
	fmt.Println("\n=== Pipeline 패턴 데모 ===")

	// 1단계: 숫자 생성
	gen := func(nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, n := range nums {
				out <- n
			}
		}()
		return out
	}

	// 2단계: 제곱
	sq := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n * n
			}
		}()
		return out
	}

	// 3단계: 2배
	double := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n * 2
			}
		}()
		return out
	}

	// 파이프라인 구성
	c := gen(1, 2, 3, 4, 5)
	c = sq(c)
	c = double(c)

	// 결과 출력
	for result := range c {
		fmt.Printf("%d ", result)
	}
	fmt.Println()
}

// ============================================================================
// 8. 동시성 패턴: Fan-out/Fan-in
// ============================================================================

/*
Fan-out: 하나의 입력을 여러 goroutine으로 분배
Fan-in: 여러 채널의 결과를 하나로 합침
*/

func FanOutFanInDemo() {
	fmt.Println("\n=== Fan-out/Fan-in 패턴 데모 ===")

	// 작업 채널
	jobs := make(chan int, 10)
	results := make(chan int, 10)

	// Fan-out: 3개의 worker
	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := range jobs {
				time.Sleep(50 * time.Millisecond) // 시뮬레이션
				results <- j * 2
				fmt.Printf("Worker %d processed job %d\n", workerID, j)
			}
		}(w)
	}

	// 작업 분배
	go func() {
		for i := 1; i <= 9; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Fan-in: 결과 수집
	go func() {
		wg.Wait()
		close(results)
	}()

	// 결과 출력
	var sum int
	for r := range results {
		sum += r
	}
	fmt.Printf("총합: %d\n", sum)
}

// ============================================================================
// 9. errgroup을 활용한 에러 처리
// ============================================================================

/*
errgroup 특징:
- 여러 goroutine의 에러를 하나로 수집
- 하나라도 실패하면 전체 취소
- context 기반 취소 지원

사용법:
import "golang.org/x/sync/errgroup"

g, ctx := errgroup.WithContext(context.Background())
g.Go(func() error { ... })
err := g.Wait()
*/

// ============================================================================
// 10. 시니어 면접 답변 요약
// ============================================================================

/*
Q: Go의 동시성 모델을 설명해주세요.

A: Go의 동시성은 CSP(Communicating Sequential Processes) 모델을 기반으로 합니다.

1. G-M-P 스케줄러:
   - G(Goroutine): 2KB 스택의 경량 스레드
   - M(Machine): OS 스레드
   - P(Processor): GOMAXPROCS 개수의 논리 프로세서
   - Work Stealing으로 부하 분산

2. 채널 vs Mutex 선택:
   - 채널: 소유권 전달, 작업 분배, 이벤트
   - Mutex: 공유 상태 보호, 성능 중요 시
   - Atomic: 단순 카운터 (가장 빠름)

3. 실무 주의점:
   - context로 goroutine 생명주기 관리
   - errgroup으로 에러 처리
   - race detector로 검증 (-race 플래그)
   - pprof로 goroutine 누수 모니터링

4. 패턴 활용:
   - Pipeline: 데이터 스트림 처리
   - Fan-out/Fan-in: 병렬 처리
   - Worker Pool: 부하 제어
   - Rate Limiting: 요청 제한
*/
