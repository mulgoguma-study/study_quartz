package medium

/*
================================================================================
문제 4: Worker Pool with Context & Error Handling - 솔루션
================================================================================
*/

import (
	"context"
	"errors"
	"sync"
	"time"
)

/*
================================================================================
Worker Pool 아키텍처
================================================================================

┌─────────────────────────────────────────────────────────────────────────────┐
│                          Worker Pool 구조                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│    ┌──────────┐                                                             │
│    │  Submit  │                                                             │
│    └────┬─────┘                                                             │
│         │                                                                   │
│         ▼                                                                   │
│    ┌────────────────────────────────────┐                                   │
│    │         Task Channel (Buffer)       │                                   │
│    │    [Task1] [Task2] [Task3] ...     │                                   │
│    └──────────────┬─────────────────────┘                                   │
│                   │                                                         │
│         ┌─────────┼─────────┐                                               │
│         ▼         ▼         ▼                                               │
│    ┌─────────┐ ┌─────────┐ ┌─────────┐                                      │
│    │ Worker1 │ │ Worker2 │ │ Worker3 │   ← 고정 개수의 워커                  │
│    └────┬────┘ └────┬────┘ └────┬────┘                                      │
│         │           │           │                                           │
│         └─────────┬─┴───────────┘                                           │
│                   ▼                                                         │
│    ┌────────────────────────────────────┐                                   │
│    │         Error Collection            │                                   │
│    │    [err1] [err2] ...               │                                   │
│    └────────────────────────────────────┘                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
*/

// ErrPoolShutdown 풀이 종료되었을 때 반환되는 에러
var ErrPoolShutdown = errors.New("worker pool is shut down")

// Task 작업 함수 타입
type Task func(ctx context.Context) error

// WorkerPool 워커 풀 구현
type WorkerPool struct {
	workerCount int             // 워커 수
	tasks       chan Task       // 작업 채널
	errors      []error         // 수집된 에러
	errorsMu    sync.Mutex      // 에러 보호 뮤텍스
	wg          sync.WaitGroup  // 작업 완료 추적
	ctx         context.Context // 컨텍스트
	cancel      context.CancelFunc
	shutdown    bool            // 종료 상태
	shutdownMu  sync.RWMutex    // 종료 상태 보호
}

// NewWorkerPool 생성자
/*
【 설계 결정 】
1. workerCount: CPU 바운드 작업 → CPU 코어 수
               I/O 바운드 작업 → 코어 수 * 2~4
2. taskBuffer: 너무 작으면 Submit 블록, 너무 크면 메모리 사용
              일반적으로 workerCount * 10 정도
*/
func NewWorkerPool(workerCount, taskBuffer int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		workerCount: workerCount,
		tasks:       make(chan Task, taskBuffer),
		errors:      make([]error, 0),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start 워커들 시작
/*
【 워커 동작 】
1. 채널에서 작업 수신
2. 컨텍스트 취소 확인
3. 작업 실행 및 에러 수집
4. WaitGroup 완료 알림
*/
func (p *WorkerPool) Start(ctx context.Context) {
	// 외부 컨텍스트와 내부 컨텍스트 병합
	p.ctx, p.cancel = context.WithCancel(ctx)

	for i := 0; i < p.workerCount; i++ {
		go p.worker(i)
	}
}

// worker 개별 워커 고루틴
func (p *WorkerPool) worker(id int) {
	for {
		select {
		case <-p.ctx.Done():
			// 컨텍스트 취소됨 → 워커 종료
			return

		case task, ok := <-p.tasks:
			if !ok {
				// 채널 닫힘 → 워커 종료
				return
			}

			// 작업 실행
			if err := p.executeTask(task); err != nil {
				p.collectError(err)
			}
		}
	}
}

// executeTask 작업 실행 (패닉 복구 포함)
/*
【 패닉 복구 이유 】
- 하나의 작업 패닉이 전체 워커 풀을 중단시키면 안 됨
- 패닉을 에러로 변환하여 수집
*/
func (p *WorkerPool) executeTask(task Task) (err error) {
	// 패닉 복구
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			case string:
				err = errors.New(v)
			default:
				err = errors.New("unknown panic")
			}
		}
		p.wg.Done()
	}()

	// 작업 실행
	return task(p.ctx)
}

// collectError 에러 수집 (thread-safe)
func (p *WorkerPool) collectError(err error) {
	p.errorsMu.Lock()
	defer p.errorsMu.Unlock()
	p.errors = append(p.errors, err)
}

// Submit 작업 제출
/*
【 동작 】
1. 종료 상태 확인
2. 작업 채널에 전송 (블록될 수 있음)
3. WaitGroup 카운트 증가
*/
func (p *WorkerPool) Submit(task Task) error {
	p.shutdownMu.RLock()
	if p.shutdown {
		p.shutdownMu.RUnlock()
		return ErrPoolShutdown
	}
	p.shutdownMu.RUnlock()

	// WaitGroup 증가는 채널 전송 전에!
	p.wg.Add(1)

	select {
	case <-p.ctx.Done():
		p.wg.Done()
		return p.ctx.Err()
	case p.tasks <- task:
		return nil
	}
}

// SubmitWithTimeout 타임아웃이 있는 작업 제출
func (p *WorkerPool) SubmitWithTimeout(task Task, timeout time.Duration) error {
	p.shutdownMu.RLock()
	if p.shutdown {
		p.shutdownMu.RUnlock()
		return ErrPoolShutdown
	}
	p.shutdownMu.RUnlock()

	p.wg.Add(1)

	select {
	case <-p.ctx.Done():
		p.wg.Done()
		return p.ctx.Err()
	case p.tasks <- task:
		return nil
	case <-time.After(timeout):
		p.wg.Done()
		return errors.New("submit timeout")
	}
}

// Wait 모든 작업 완료 대기 및 에러 반환
func (p *WorkerPool) Wait() []error {
	p.wg.Wait()

	p.errorsMu.Lock()
	defer p.errorsMu.Unlock()

	// 에러 복사본 반환 (원본 보호)
	result := make([]error, len(p.errors))
	copy(result, p.errors)

	return result
}

// Shutdown 종료 (새 작업 거부, 진행 중인 작업 완료 대기)
/*
【 Graceful Shutdown 순서 】
1. shutdown 플래그 설정 (새 작업 거부)
2. 작업 채널 닫기 (워커들에게 종료 신호)
3. WaitGroup 대기 (진행 중인 작업 완료)
4. 컨텍스트 취소 (혹시 남은 워커 정리)
*/
func (p *WorkerPool) Shutdown() {
	p.shutdownMu.Lock()
	p.shutdown = true
	p.shutdownMu.Unlock()

	// 채널 닫기 (더 이상 작업 받지 않음)
	close(p.tasks)

	// 진행 중인 작업 완료 대기
	p.wg.Wait()

	// 컨텍스트 취소
	p.cancel()
}

// ShutdownNow 즉시 종료 (진행 중인 작업 취소)
func (p *WorkerPool) ShutdownNow() {
	p.shutdownMu.Lock()
	p.shutdown = true
	p.shutdownMu.Unlock()

	// 먼저 컨텍스트 취소 (진행 중인 작업에 취소 신호)
	p.cancel()

	// 채널 닫기
	close(p.tasks)
}

/*
================================================================================
고급 기능 확장
================================================================================
*/

// WorkerPoolWithMetrics 메트릭이 포함된 워커 풀
type WorkerPoolWithMetrics struct {
	*WorkerPool
	submitted  int64      // 제출된 작업 수
	completed  int64      // 완료된 작업 수
	failed     int64      // 실패한 작업 수
	metricsMu  sync.Mutex
}

// Result 작업 결과
type Result struct {
	Value interface{}
	Err   error
}

// TaskWithResult 결과를 반환하는 작업
type TaskWithResult func(ctx context.Context) (interface{}, error)

// WorkerPoolWithResults 결과를 수집하는 워커 풀
type WorkerPoolWithResults struct {
	pool    *WorkerPool
	results chan Result
}

func NewWorkerPoolWithResults(workerCount, taskBuffer int) *WorkerPoolWithResults {
	return &WorkerPoolWithResults{
		pool:    NewWorkerPool(workerCount, taskBuffer),
		results: make(chan Result, taskBuffer),
	}
}

/*
================================================================================
사용 패턴 예시
================================================================================

【 패턴 1: 배치 처리 】

func processBatch(items []Item) []error {
    pool := NewWorkerPool(runtime.NumCPU(), len(items))
    pool.Start(context.Background())
    defer pool.Shutdown()

    for _, item := range items {
        item := item  // 클로저 캡처
        pool.Submit(func(ctx context.Context) error {
            return processItem(ctx, item)
        })
    }

    return pool.Wait()
}

【 패턴 2: 타임아웃 있는 API 호출 】

func fetchAll(urls []string, timeout time.Duration) ([]Response, []error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    pool := NewWorkerPool(10, len(urls))
    pool.Start(ctx)

    responses := make([]Response, len(urls))
    var mu sync.Mutex

    for i, url := range urls {
        i, url := i, url
        pool.Submit(func(ctx context.Context) error {
            resp, err := fetchWithContext(ctx, url)
            if err != nil {
                return err
            }
            mu.Lock()
            responses[i] = resp
            mu.Unlock()
            return nil
        })
    }

    return responses, pool.Wait()
}

【 패턴 3: 파이프라인 】

func pipeline(input <-chan Data) <-chan Result {
    output := make(chan Result)
    pool := NewWorkerPool(5, 100)
    pool.Start(context.Background())

    go func() {
        defer close(output)
        defer pool.Shutdown()

        for data := range input {
            data := data
            pool.Submit(func(ctx context.Context) error {
                result := process(data)
                output <- result
                return nil
            })
        }
        pool.Wait()
    }()

    return output
}

================================================================================
*/

// 테스트
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool := NewWorkerPool(3, 10)
	pool.Start(ctx)

	// 다양한 작업 제출
	for i := 0; i < 10; i++ {
		taskID := i
		err := pool.Submit(func(ctx context.Context) error {
			// 컨텍스트 취소 확인
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			println("Task", taskID, "started")
			time.Sleep(100 * time.Millisecond)
			println("Task", taskID, "completed")

			// 일부 작업은 에러 반환
			if taskID%3 == 0 {
				return errors.New("task failed")
			}
			return nil
		})

		if err != nil {
			println("Submit error:", err.Error())
		}
	}

	// 결과 대기
	errs := pool.Wait()
	println("\nTotal errors:", len(errs))
	for _, err := range errs {
		println("- Error:", err.Error())
	}

	// 종료
	pool.Shutdown()
	println("\nPool shutdown complete")
}
