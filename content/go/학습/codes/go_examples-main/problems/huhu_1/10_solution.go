package huhu1

/*
================================================================================
문제 10: 분산 작업 스케줄러 - 솔루션
================================================================================

【 핵심 설계 원칙 】
1. 우선순위 큐 - 중요한 작업 먼저 처리
2. 워커 풀 - 가용 워커에 작업 분배
3. 헬스체크 - 워커 장애 감지
4. 재시도 + 지수 백오프 - 일시적 장애 복구

【 작업 상태 전이 】
┌─────────┐  할당   ┌─────────┐  완료   ┌───────────┐
│ PENDING │───────→│ RUNNING │───────→│ COMPLETED │
└─────────┘        └─────────┘        └───────────┘
     ↑                  │
     │ 재시도           │ 실패
     │                  ↓
     │            ┌─────────┐
     └────────────│ FAILED  │ (재시도 초과 시 최종 실패)
                  └─────────┘

================================================================================
*/

import (
	"container/heap"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// 우선순위 큐
// ============================================================================

type jobItem struct {
	job      *Job
	priority int // 낮을수록 먼저 (음수 사용)
	index    int
}

type jobPriorityQueue []*jobItem

func (pq jobPriorityQueue) Len() int { return len(pq) }

func (pq jobPriorityQueue) Less(i, j int) bool {
	// 우선순위가 같으면 생성 시간순
	if pq[i].priority == pq[j].priority {
		return pq[i].job.CreatedAt.Before(pq[j].job.CreatedAt)
	}
	return pq[i].priority < pq[j].priority
}

func (pq jobPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *jobPriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*jobItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *jobPriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// ============================================================================
// 작업 스케줄러 구현
// ============================================================================

// JobSchedulerSolution 완성된 작업 스케줄러
type JobSchedulerSolution struct {
	mu sync.RWMutex

	// 작업 큐 (우선순위)
	pendingJobs *jobPriorityQueue
	queueCond   *sync.Cond

	// 작업 상태
	jobStatuses map[string]*JobStatus
	jobs        map[string]*Job

	// 워커 관리
	workers map[string]*Worker

	// 스케줄된 작업
	scheduledJobs map[string]*scheduledJob

	// 설정
	heartbeatTimeout time.Duration
	maxRetries       int
	baseRetryDelay   time.Duration

	// 라이프사이클
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// ID 생성
	nextJobID int64
}

type scheduledJob struct {
	job      *Job
	cron     string
	nextRun  time.Time
	cancel   context.CancelFunc
}

// NewJobSchedulerSolution 생성자
func NewJobSchedulerSolution() *JobSchedulerSolution {
	ctx, cancel := context.WithCancel(context.Background())

	pq := make(jobPriorityQueue, 0)
	heap.Init(&pq)

	s := &JobSchedulerSolution{
		pendingJobs:      &pq,
		jobStatuses:      make(map[string]*JobStatus),
		jobs:             make(map[string]*Job),
		workers:          make(map[string]*Worker),
		scheduledJobs:    make(map[string]*scheduledJob),
		heartbeatTimeout: 30 * time.Second,
		maxRetries:       3,
		baseRetryDelay:   time.Second,
		ctx:              ctx,
		cancel:           cancel,
	}

	s.queueCond = sync.NewCond(&s.mu)

	// 워커 모니터 시작
	s.wg.Add(1)
	go s.workerMonitor()

	return s
}

// ============================================================================
// 작업 제출
// ============================================================================

// SubmitJob 작업 제출
func (s *JobSchedulerSolution) SubmitJob(job Job) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// ID 생성
	s.nextJobID++
	job.ID = fmt.Sprintf("job-%d", s.nextJobID)

	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now()
	}
	if job.MaxRetries == 0 {
		job.MaxRetries = s.maxRetries
	}

	// 저장
	s.jobs[job.ID] = &job
	s.jobStatuses[job.ID] = &JobStatus{
		JobID: job.ID,
		State: "PENDING",
	}

	// 큐에 추가
	heap.Push(s.pendingJobs, &jobItem{
		job:      &job,
		priority: -job.Priority, // 높은 우선순위 = 작은 값
	})

	// 대기 중인 워커 깨우기
	s.queueCond.Signal()

	return job.ID, nil
}

// ScheduleJob 주기적 작업 스케줄
func (s *JobSchedulerSolution) ScheduleJob(job Job, cron string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// ID 생성
	s.nextJobID++
	job.ID = fmt.Sprintf("scheduled-job-%d", s.nextJobID)
	job.Type = "SCHEDULED"
	job.Cron = cron

	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now()
	}

	// 다음 실행 시간 계산 (간단한 구현: 1분 후)
	nextRun := time.Now().Add(time.Minute)

	ctx, cancel := context.WithCancel(s.ctx)

	s.scheduledJobs[job.ID] = &scheduledJob{
		job:     &job,
		cron:    cron,
		nextRun: nextRun,
		cancel:  cancel,
	}

	s.jobs[job.ID] = &job
	s.jobStatuses[job.ID] = &JobStatus{
		JobID: job.ID,
		State: "SCHEDULED",
	}

	// 스케줄 고루틴 시작
	s.wg.Add(1)
	go s.runScheduledJob(ctx, job.ID)

	return job.ID, nil
}

// runScheduledJob 스케줄된 작업 실행
func (s *JobSchedulerSolution) runScheduledJob(ctx context.Context, jobID string) {
	defer s.wg.Done()

	for {
		s.mu.RLock()
		sj, ok := s.scheduledJobs[jobID]
		if !ok {
			s.mu.RUnlock()
			return
		}
		nextRun := sj.nextRun
		s.mu.RUnlock()

		// 다음 실행까지 대기
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Until(nextRun)):
		}

		// 작업 복사하여 제출
		s.mu.RLock()
		sj, ok = s.scheduledJobs[jobID]
		if !ok {
			s.mu.RUnlock()
			return
		}
		jobCopy := *sj.job
		s.mu.RUnlock()

		jobCopy.Type = "ONCE"
		s.SubmitJob(jobCopy)

		// 다음 실행 시간 업데이트 (간단한 구현: 1분 후)
		s.mu.Lock()
		if sj, ok := s.scheduledJobs[jobID]; ok {
			sj.nextRun = time.Now().Add(time.Minute)
		}
		s.mu.Unlock()
	}
}

// CancelJob 작업 취소
func (s *JobSchedulerSolution) CancelJob(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	status, ok := s.jobStatuses[jobID]
	if !ok {
		return errors.New("job not found")
	}

	if status.State == "COMPLETED" || status.State == "CANCELLED" {
		return errors.New("job already finished")
	}

	status.State = "CANCELLED"
	now := time.Now()
	status.CompletedAt = &now

	// 스케줄된 작업이면 취소
	if sj, ok := s.scheduledJobs[jobID]; ok {
		sj.cancel()
		delete(s.scheduledJobs, jobID)
	}

	return nil
}

// GetJobStatus 작업 상태 조회
func (s *JobSchedulerSolution) GetJobStatus(jobID string) (*JobStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status, ok := s.jobStatuses[jobID]
	if !ok {
		return nil, errors.New("job not found")
	}

	// 복사본 반환
	statusCopy := *status
	return &statusCopy, nil
}

// ============================================================================
// 워커 관리
// ============================================================================

// RegisterWorker 워커 등록
func (s *JobSchedulerSolution) RegisterWorker(worker Worker) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	worker.Status = "ACTIVE"
	worker.LastHeartbeat = time.Now()

	s.workers[worker.ID] = &worker
	return nil
}

// UnregisterWorker 워커 해제
func (s *JobSchedulerSolution) UnregisterWorker(workerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.workers, workerID)
	return nil
}

// Heartbeat 워커 헬스체크
func (s *JobSchedulerSolution) Heartbeat(workerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	worker, ok := s.workers[workerID]
	if !ok {
		return errors.New("worker not found")
	}

	worker.LastHeartbeat = time.Now()
	if worker.Status == "INACTIVE" {
		worker.Status = "ACTIVE"
	}

	return nil
}

// workerMonitor 워커 상태 모니터링
func (s *JobSchedulerSolution) workerMonitor() {
	defer s.wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.checkWorkerHealth()
		}
	}
}

// checkWorkerHealth 워커 헬스 체크
func (s *JobSchedulerSolution) checkWorkerHealth() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for _, worker := range s.workers {
		if now.Sub(worker.LastHeartbeat) > s.heartbeatTimeout {
			if worker.Status != "INACTIVE" {
				worker.Status = "INACTIVE"
				// 해당 워커의 실행 중인 작업 재큐잉
				s.requeueWorkerJobs(worker.ID)
			}
		}
	}
}

// requeueWorkerJobs 워커의 작업 재큐잉
func (s *JobSchedulerSolution) requeueWorkerJobs(workerID string) {
	for _, status := range s.jobStatuses {
		if status.WorkerID == workerID && status.State == "RUNNING" {
			// 작업 다시 큐에 추가
			if job, ok := s.jobs[status.JobID]; ok {
				status.State = "PENDING"
				status.WorkerID = ""

				heap.Push(s.pendingJobs, &jobItem{
					job:      job,
					priority: -job.Priority,
				})
			}
		}
	}
	s.queueCond.Broadcast()
}

// ============================================================================
// 작업 실행
// ============================================================================

// GetNextJob 다음 작업 가져오기
func (s *JobSchedulerSolution) GetNextJob(workerID string) (*Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 워커 확인
	worker, ok := s.workers[workerID]
	if !ok {
		return nil, errors.New("worker not found")
	}

	if worker.Status == "INACTIVE" {
		return nil, errors.New("worker is inactive")
	}

	if worker.CurrentJobs >= worker.Capacity {
		return nil, errors.New("worker at capacity")
	}

	// 작업 대기
	for s.pendingJobs.Len() == 0 {
		s.queueCond.Wait()

		// 컨텍스트 취소 확인
		select {
		case <-s.ctx.Done():
			return nil, errors.New("scheduler stopped")
		default:
		}
	}

	// 작업 꺼내기
	item := heap.Pop(s.pendingJobs).(*jobItem)
	job := item.job

	// 상태 업데이트
	status := s.jobStatuses[job.ID]
	status.State = "RUNNING"
	status.WorkerID = workerID
	now := time.Now()
	status.StartedAt = &now

	worker.CurrentJobs++
	worker.Status = "BUSY"

	return job, nil
}

// CompleteJob 작업 완료
func (s *JobSchedulerSolution) CompleteJob(jobID string, result any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	status, ok := s.jobStatuses[jobID]
	if !ok {
		return errors.New("job not found")
	}

	status.State = "COMPLETED"
	status.Result = result
	now := time.Now()
	status.CompletedAt = &now

	// 워커 상태 업데이트
	if worker, ok := s.workers[status.WorkerID]; ok {
		worker.CurrentJobs--
		if worker.CurrentJobs == 0 {
			worker.Status = "ACTIVE"
		}
	}

	return nil
}

// FailJob 작업 실패
/*
【 재시도 전략: 지수 백오프 】
1차 재시도: 1초 후
2차 재시도: 2초 후
3차 재시도: 4초 후
...
*/
func (s *JobSchedulerSolution) FailJob(jobID string, err error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	status, ok := s.jobStatuses[jobID]
	if !ok {
		return errors.New("job not found")
	}

	job, ok := s.jobs[jobID]
	if !ok {
		return errors.New("job not found")
	}

	status.Retries++
	status.Error = err.Error()

	// 워커 상태 업데이트
	if worker, ok := s.workers[status.WorkerID]; ok {
		worker.CurrentJobs--
		if worker.CurrentJobs == 0 {
			worker.Status = "ACTIVE"
		}
	}
	status.WorkerID = ""

	// 재시도 가능 여부 확인
	if status.Retries < job.MaxRetries {
		// 지수 백오프로 재시도
		delay := s.baseRetryDelay * time.Duration(1<<uint(status.Retries-1))

		go func() {
			time.Sleep(delay)

			s.mu.Lock()
			defer s.mu.Unlock()

			// 재큐잉
			status.State = "PENDING"
			heap.Push(s.pendingJobs, &jobItem{
				job:      job,
				priority: -job.Priority,
			})
			s.queueCond.Signal()
		}()
	} else {
		// 최종 실패
		status.State = "FAILED"
		now := time.Now()
		status.CompletedAt = &now
	}

	return nil
}

// Stop 스케줄러 정지
func (s *JobSchedulerSolution) Stop() {
	s.cancel()
	s.queueCond.Broadcast()
	s.wg.Wait()
}

// ============================================================================
// 성능 최적화 포인트
// ============================================================================

/*
【 프로덕션 최적화 】

1. 분산 작업 큐 (Redis/RabbitMQ)
   - 여러 스케줄러 인스턴스 간 공유
   - 작업 영속성 보장
   - 분산 락으로 중복 실행 방지

2. 작업 파티셔닝
   - 작업 타입별 큐 분리
   - 우선순위별 큐 분리
   - 워커 전문화

3. Dead Letter Queue (DLQ)
   - 최종 실패 작업 별도 저장
   - 수동 재처리 지원
   - 실패 분석

4. 작업 의존성
   - DAG (Directed Acyclic Graph)
   - 선행 작업 완료 후 실행

5. 리소스 제한
   - 워커별 메모리/CPU 제한
   - 작업별 타임아웃
   - 동시 실행 제한

┌─────────────────────────────────────────────────────────────────────────────┐
│                       프로덕션 아키텍처                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  [Client API]                                                               │
│       │                                                                     │
│       ↓                                                                     │
│  ┌─────────────┐                                                           │
│  │  Scheduler  │ ←── Cron Trigger                                          │
│  │  Service    │                                                           │
│  └──────┬──────┘                                                           │
│         │                                                                   │
│         ↓                                                                   │
│  ┌─────────────────────────────────────────────────────────────┐           │
│  │                      Redis / RabbitMQ                        │           │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐        │           │
│  │  │ High    │  │ Normal  │  │ Low     │  │ DLQ     │        │           │
│  │  │ Priority│  │ Priority│  │ Priority│  │         │        │           │
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘        │           │
│  └─────────────────────────────────────────────────────────────┘           │
│         │              │              │                                     │
│         ↓              ↓              ↓                                     │
│  ┌─────────────────────────────────────────────────────────────┐           │
│  │                      Worker Pool                             │           │
│  │  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐            │           │
│  │  │Worker 1│  │Worker 2│  │Worker 3│  │Worker N│            │           │
│  │  └────────┘  └────────┘  └────────┘  └────────┘            │           │
│  └─────────────────────────────────────────────────────────────┘           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

【 Redis 기반 구현 예시 】

-- 작업 추가 (우선순위 큐)
ZADD jobs:pending <priority> <job_json>

-- 작업 가져오기 (원자적)
BZPOPMIN jobs:pending 30

-- 분산 락 (작업 실행)
SET lock:job:<id> <worker_id> NX EX 300

-- 작업 상태 저장
HSET jobs:status:<id> state "RUNNING" worker_id "<wid>"

【 Cron 파싱 예시 】
"0 * * * *"  → 매 시 정각
"*/5 * * * *" → 5분마다
"0 9 * * 1-5" → 평일 오전 9시
*/
