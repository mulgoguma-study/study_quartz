package medium

/*
================================================================================
문제 4: Worker Pool with Context & Error Handling
================================================================================

난이도: Medium
주제: 고루틴, 채널, 컨텍스트, 에러 처리

【 문제 설명 】
실무에서 사용할 수 있는 Worker Pool을 구현하세요.

요구사항:
1. 고정된 수의 워커로 작업 처리
2. Context를 통한 취소/타임아웃 지원
3. 에러 수집 및 반환
4. Graceful Shutdown

【 인터페이스 】
type Task func(ctx context.Context) error

type WorkerPool interface {
    Submit(task Task) error     // 작업 제출
    Wait() []error              // 모든 작업 완료 대기, 에러 반환
    Shutdown()                  // 새 작업 거부, 진행 중인 작업 완료 대기
}

【 예시 】
pool := NewWorkerPool(5, 100)  // 5 workers, 100 task buffer

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 작업 제출
for i := 0; i < 20; i++ {
    taskID := i
    pool.Submit(func(ctx context.Context) error {
        // 작업 수행
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(100 * time.Millisecond):
            fmt.Printf("Task %d completed\n", taskID)
            return nil
        }
    })
}

// 결과 대기
errors := pool.Wait()
for _, err := range errors {
    fmt.Println("Error:", err)
}

pool.Shutdown()

【 제약 조건 】
- Submit은 버퍼가 가득 차면 블록되거나 에러 반환
- Context 취소 시 진행 중인 작업에 알림
- 모든 에러 수집하여 반환
- 데드락 없이 정상 종료

【 힌트 】
1. 버퍼드 채널로 작업 큐 구현
2. WaitGroup으로 작업 완료 추적
3. 별도의 에러 채널 또는 뮤텍스로 에러 수집
4. done 채널로 종료 시그널

【 실무 연관성 】
- 배치 작업 처리
- API 호출 병렬화
- 파일 처리
- 데이터베이스 마이그레이션
================================================================================
*/

import (
	"context"
)

// Task 작업 함수 타입
type Task func(ctx context.Context) error

// WorkerPool 워커 풀 구조체
type WorkerPool struct {
	// 여기에 필드를 정의하세요
	// workerCount: 워커 수
	// tasks: 작업 채널
	// errors: 에러 수집
	// ctx, cancel: 컨텍스트
	// wg: WaitGroup
}

// NewWorkerPool 생성자
func NewWorkerPool(workerCount, taskBuffer int) *WorkerPool {
	// 여기에 코드를 작성하세요
	return nil
}

// Start 워커 시작
func (p *WorkerPool) Start(ctx context.Context) {
	// 여기에 코드를 작성하세요
}

// Submit 작업 제출
func (p *WorkerPool) Submit(task Task) error {
	// 여기에 코드를 작성하세요
	return nil
}

// Wait 모든 작업 완료 대기
func (p *WorkerPool) Wait() []error {
	// 여기에 코드를 작성하세요
	return nil
}

// Shutdown 종료
func (p *WorkerPool) Shutdown() {
	// 여기에 코드를 작성하세요
}

// 테스트
func main() {
	ctx := context.Background()
	pool := NewWorkerPool(3, 10)
	pool.Start(ctx)

	// 작업 제출
	for i := 0; i < 5; i++ {
		taskID := i
		pool.Submit(func(ctx context.Context) error {
			println("Processing task", taskID)
			return nil
		})
	}

	// 결과 대기
	errors := pool.Wait()
	println("Errors:", len(errors))

	pool.Shutdown()
}
