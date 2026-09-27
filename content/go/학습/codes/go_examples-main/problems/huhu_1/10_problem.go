package huhu1

/*
================================================================================
문제 10: 분산 작업 스케줄러
================================================================================

난이도: Hard
주제: 분산 시스템, 작업 스케줄링, 장애 복구
회사: 42dot (백엔드/데이터 플랫폼)

【 문제 설명 】
대규모 배치 작업을 여러 워커에 분산하여 실행하는 작업 스케줄러를 구현하세요.
작업 실패 시 재시도하고, 워커 장애 시 다른 워커가 인계받아야 합니다.

작업 예시:
- 일일 리포트 생성
- 데이터 백업
- 로그 압축/정리
- 데이터 동기화
- 알림 발송

요구사항:
1. SubmitJob(job): 작업 제출
2. ScheduleJob(job, cron): 주기적 작업 스케줄
3. CancelJob(jobID): 작업 취소
4. GetJobStatus(jobID): 작업 상태 조회
5. RegisterWorker(workerID): 워커 등록
6. Heartbeat(workerID): 워커 헬스체크

【 성능 요구사항 】
- 동시 작업: 10,000개 이상
- 작업 할당: 100ms 이내
- 장애 감지: 30초 이내
- 재시도: 최대 3회, 지수 백오프

【 인터페이스 】
type Job struct {
    ID          string
    Name        string
    Type        string    // "ONCE", "SCHEDULED"
    Priority    int       // 0-9 (높을수록 우선)
    Payload     map[string]interface{}
    Cron        string    // cron expression (scheduled jobs)
    MaxRetries  int
    Timeout     time.Duration
    CreatedAt   time.Time
}

type JobStatus struct {
    JobID       string
    State       string    // "PENDING", "RUNNING", "COMPLETED", "FAILED", "CANCELLED"
    WorkerID    string
    Progress    int       // 0-100
    Result      interface{}
    Error       string
    Retries     int
    StartedAt   *time.Time
    CompletedAt *time.Time
}

type Worker struct {
    ID          string
    Status      string    // "ACTIVE", "BUSY", "INACTIVE"
    Capacity    int       // 동시 실행 가능 작업 수
    CurrentJobs int
    LastHeartbeat time.Time
}

type JobScheduler interface {
    SubmitJob(job Job) (string, error)
    ScheduleJob(job Job, cron string) (string, error)
    CancelJob(jobID string) error
    GetJobStatus(jobID string) (*JobStatus, error)
    RegisterWorker(worker Worker) error
    UnregisterWorker(workerID string) error
    Heartbeat(workerID string) error
    GetNextJob(workerID string) (*Job, error)
    CompleteJob(jobID string, result interface{}) error
    FailJob(jobID string, err error) error
}

【 힌트 】
1. 우선순위 큐로 작업 관리
2. 헬스체크로 워커 상태 모니터링
3. 지수 백오프로 재시도
4. 분산 락으로 작업 중복 실행 방지

【 실무 연관성 】
- Celery, Sidekiq 같은 작업 큐
- Kubernetes Job/CronJob
- AWS Step Functions, Airflow
================================================================================
*/

import "time"

// Job 작업 정의
type Job struct {
	ID         string
	Name       string
	Type       string
	Priority   int
	Payload    map[string]any
	Cron       string
	MaxRetries int
	Timeout    time.Duration
	CreatedAt  time.Time
}

// JobStatus 작업 상태
type JobStatus struct {
	JobID       string
	State       string
	WorkerID    string
	Progress    int
	Result      any
	Error       string
	Retries     int
	StartedAt   *time.Time
	CompletedAt *time.Time
}

// Worker 워커
type Worker struct {
	ID            string
	Status        string
	Capacity      int
	CurrentJobs   int
	LastHeartbeat time.Time
}

// JobScheduler 작업 스케줄러 인터페이스
type JobScheduler interface {
	SubmitJob(job Job) (string, error)
	ScheduleJob(job Job, cron string) (string, error)
	CancelJob(jobID string) error
	GetJobStatus(jobID string) (*JobStatus, error)
	RegisterWorker(worker Worker) error
	UnregisterWorker(workerID string) error
	Heartbeat(workerID string) error
	GetNextJob(workerID string) (*Job, error)
	CompleteJob(jobID string, result any) error
	FailJob(jobID string, err error) error
}

// JobSchedulerImpl 구현체
type JobSchedulerImpl struct {
	// 여기에 필드를 정의하세요
}

// NewJobScheduler 생성자
func NewJobScheduler() *JobSchedulerImpl {
	// 여기에 코드를 작성하세요
	return nil
}

// SubmitJob 작업 제출
func (s *JobSchedulerImpl) SubmitJob(job Job) (string, error) {
	// 여기에 코드를 작성하세요
	return "", nil
}

// ScheduleJob 주기적 작업 스케줄
func (s *JobSchedulerImpl) ScheduleJob(job Job, cron string) (string, error) {
	// 여기에 코드를 작성하세요
	return "", nil
}

// CancelJob 작업 취소
func (s *JobSchedulerImpl) CancelJob(jobID string) error {
	// 여기에 코드를 작성하세요
	return nil
}

// GetJobStatus 작업 상태 조회
func (s *JobSchedulerImpl) GetJobStatus(jobID string) (*JobStatus, error) {
	// 여기에 코드를 작성하세요
	return nil, nil
}

// RegisterWorker 워커 등록
func (s *JobSchedulerImpl) RegisterWorker(worker Worker) error {
	// 여기에 코드를 작성하세요
	return nil
}

// UnregisterWorker 워커 해제
func (s *JobSchedulerImpl) UnregisterWorker(workerID string) error {
	// 여기에 코드를 작성하세요
	return nil
}

// Heartbeat 워커 헬스체크
func (s *JobSchedulerImpl) Heartbeat(workerID string) error {
	// 여기에 코드를 작성하세요
	return nil
}

// GetNextJob 다음 작업 가져오기
func (s *JobSchedulerImpl) GetNextJob(workerID string) (*Job, error) {
	// 여기에 코드를 작성하세요
	return nil, nil
}

// CompleteJob 작업 완료
func (s *JobSchedulerImpl) CompleteJob(jobID string, result any) error {
	// 여기에 코드를 작성하세요
	return nil
}

// FailJob 작업 실패
func (s *JobSchedulerImpl) FailJob(jobID string, err error) error {
	// 여기에 코드를 작성하세요
	return nil
}
