// =============================================================================
// Worker Pool (워커 풀) 패턴 예제
// =============================================================================
// 워커 풀은 고정된 수의 고루틴(워커)이 작업을 분담 처리하는 패턴입니다.
//
// 장점:
// - 고루틴 생성/소멸 오버헤드 감소
// - 동시성 수준 제어 (리소스 사용량 제한)
// - 작업 큐잉 지원
// - 효율적인 리소스 활용
//
// 구성요소:
// - Job: 처리할 작업
// - Worker: 작업을 처리하는 고루틴
// - Job Queue: 작업을 대기시키는 채널
// - Result Queue: 결과를 수집하는 채널
// =============================================================================

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// =============================================================================
// 1. 기본 워커 풀
// =============================================================================

// Task 작업 구조체
type Task struct {
	ID      int
	Payload string
}

// TaskResult 작업 결과 구조체
type TaskResult struct {
	TaskID   int
	Result   string
	WorkerID int
	Duration time.Duration
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== 기본 워커 풀 예제 ===")
	basicWorkerPoolExample()

	fmt.Println("\n=== 동적 워커 풀 예제 ===")
	dynamicWorkerPoolExample()

	fmt.Println("\n=== 제네릭 워커 풀 예제 ===")
	genericWorkerPoolExample()
}

// basicWorkerPoolExample 기본적인 워커 풀 패턴
func basicWorkerPoolExample() {
	const numWorkers = 3
	const numJobs = 10

	// 작업 채널과 결과 채널 생성
	jobs := make(chan Task, numJobs)
	results := make(chan TaskResult, numJobs)

	var wg sync.WaitGroup

	// ==========================================================================
	// 워커 생성 (고정된 수의 고루틴)
	// ==========================================================================
	for workerID := 1; workerID <= numWorkers; workerID++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			// 워커는 작업 채널이 닫힐 때까지 계속 작업을 가져옴
			for task := range jobs {
				startTime := time.Now()

				// 작업 처리 시뮬레이션
				processingTime := time.Duration(rand.Intn(500)) * time.Millisecond
				time.Sleep(processingTime)

				// 결과 생성
				result := TaskResult{
					TaskID:   task.ID,
					Result:   fmt.Sprintf("Task %d 처리완료: %s", task.ID, task.Payload),
					WorkerID: id,
					Duration: time.Since(startTime),
				}

				results <- result
			}

			fmt.Printf("Worker %d: 종료\n", id)
		}(workerID)
	}

	// ==========================================================================
	// 작업 생성 및 전송
	// ==========================================================================
	for i := 1; i <= numJobs; i++ {
		jobs <- Task{
			ID:      i,
			Payload: fmt.Sprintf("데이터_%d", i),
		}
	}

	// 모든 작업 전송 후 채널 닫기
	close(jobs)

	// 별도의 고루틴에서 결과 채널 닫기
	go func() {
		wg.Wait()
		close(results)
	}()

	// ==========================================================================
	// 결과 수집
	// ==========================================================================
	for result := range results {
		fmt.Printf("결과: %s (Worker %d, 소요시간: %v)\n",
			result.Result, result.WorkerID, result.Duration)
	}
}

// =============================================================================
// 2. 동적 워커 풀 (작업량에 따라 워커 수 조절 가능)
// =============================================================================

// WorkerPool 워커 풀 구조체
type WorkerPool struct {
	numWorkers int
	jobs       chan func()
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewWorkerPool 워커 풀 생성
func NewWorkerPool(numWorkers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan func(), 100), // 버퍼드 작업 큐
		ctx:        ctx,
		cancel:     cancel,
	}

	// 워커 시작
	pool.start()

	return pool
}

// start 워커 시작
func (p *WorkerPool) start() {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)

		go func(workerID int) {
			defer p.wg.Done()

			for {
				select {
				case <-p.ctx.Done():
					// 종료 신호 수신
					fmt.Printf("Worker %d: 종료\n", workerID)
					return

				case job, ok := <-p.jobs:
					if !ok {
						// 채널이 닫힘
						return
					}
					// 작업 실행
					job()
				}
			}
		}(i + 1)
	}
}

// Submit 작업 제출
func (p *WorkerPool) Submit(job func()) {
	select {
	case p.jobs <- job:
		// 작업 제출 성공
	case <-p.ctx.Done():
		// 풀이 종료됨
		fmt.Println("워커 풀이 종료되어 작업을 제출할 수 없습니다")
	}
}

// Shutdown 워커 풀 종료
func (p *WorkerPool) Shutdown() {
	p.cancel()  // 종료 신호 전송
	close(p.jobs)
	p.wg.Wait() // 모든 워커 종료 대기
}

// dynamicWorkerPoolExample 동적 워커 풀 사용 예제
func dynamicWorkerPoolExample() {
	// 5개의 워커를 가진 풀 생성
	pool := NewWorkerPool(5)

	var resultMu sync.Mutex
	results := make([]string, 0)

	var jobWg sync.WaitGroup

	// 작업 제출
	for i := 1; i <= 10; i++ {
		jobWg.Add(1)

		taskID := i // 클로저 캡처용

		pool.Submit(func() {
			defer jobWg.Done()

			// 작업 수행
			time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)

			// 결과 저장 (뮤텍스로 보호)
			resultMu.Lock()
			results = append(results, fmt.Sprintf("Task %d 완료", taskID))
			resultMu.Unlock()
		})
	}

	// 모든 작업 완료 대기
	jobWg.Wait()

	// 워커 풀 종료
	pool.Shutdown()

	// 결과 출력
	fmt.Println("처리된 작업들:")
	for _, r := range results {
		fmt.Printf("  - %s\n", r)
	}
}

// =============================================================================
// 3. 제네릭 워커 풀 (Go 1.18+)
// =============================================================================

// GenericWorkerPool 제네릭 워커 풀
type GenericWorkerPool[T any, R any] struct {
	numWorkers int
	processor  func(T) R
	jobs       chan T
	results    chan R
	wg         sync.WaitGroup
}

// NewGenericWorkerPool 제네릭 워커 풀 생성
func NewGenericWorkerPool[T any, R any](numWorkers int, processor func(T) R) *GenericWorkerPool[T, R] {
	pool := &GenericWorkerPool[T, R]{
		numWorkers: numWorkers,
		processor:  processor,
		jobs:       make(chan T, 100),
		results:    make(chan R, 100),
	}

	// 워커 시작
	for i := 0; i < numWorkers; i++ {
		pool.wg.Add(1)

		go func() {
			defer pool.wg.Done()

			for job := range pool.jobs {
				result := pool.processor(job)
				pool.results <- result
			}
		}()
	}

	return pool
}

// Submit 작업 제출
func (p *GenericWorkerPool[T, R]) Submit(job T) {
	p.jobs <- job
}

// Results 결과 채널 반환
func (p *GenericWorkerPool[T, R]) Results() <-chan R {
	return p.results
}

// Close 워커 풀 종료
func (p *GenericWorkerPool[T, R]) Close() {
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}

// ImageTask 이미지 처리 작업 (예시)
type ImageTask struct {
	ImageID   int
	ImagePath string
}

// ImageResult 이미지 처리 결과 (예시)
type ImageResult struct {
	ImageID    int
	Thumbnail  string
	Processing time.Duration
}

// genericWorkerPoolExample 제네릭 워커 풀 사용 예제
func genericWorkerPoolExample() {
	// 이미지 처리 함수 정의
	imageProcessor := func(task ImageTask) ImageResult {
		start := time.Now()

		// 이미지 처리 시뮬레이션
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		return ImageResult{
			ImageID:    task.ImageID,
			Thumbnail:  fmt.Sprintf("thumb_%s", task.ImagePath),
			Processing: time.Since(start),
		}
	}

	// 이미지 처리용 워커 풀 생성
	pool := NewGenericWorkerPool(3, imageProcessor)

	// 작업 제출
	numImages := 8
	for i := 1; i <= numImages; i++ {
		pool.Submit(ImageTask{
			ImageID:   i,
			ImagePath: fmt.Sprintf("image_%d.jpg", i),
		})
	}

	// 결과 수집 (별도 고루틴에서)
	go func() {
		time.Sleep(500 * time.Millisecond)
		pool.Close()
	}()

	// 결과 출력
	count := 0
	for result := range pool.Results() {
		count++
		fmt.Printf("이미지 %d 처리 완료: %s (소요시간: %v)\n",
			result.ImageID, result.Thumbnail, result.Processing)

		if count >= numImages {
			break
		}
	}
}
