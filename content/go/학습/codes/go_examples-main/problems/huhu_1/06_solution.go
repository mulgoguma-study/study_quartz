package huhu1

/*
================================================================================
문제 6: 대용량 텔레메트리 데이터 파이프라인 - 솔루션
================================================================================

【 핵심 설계 원칙 】
1. Ring Buffer - 고정 크기 버퍼로 메모리 효율화
2. Worker Pool - 병렬 처리로 처리량 극대화
3. Batch Flush - I/O 최적화
4. Back-pressure - 과부하 시 흐름 제어

【 아키텍처 】
┌─────────┐    ┌─────────────┐    ┌─────────────┐    ┌──────────┐
│ Ingest  │ →  │ Ring Buffer │ →  │ Worker Pool │ →  │  Flush   │
└─────────┘    └─────────────┘    └─────────────┘    └──────────┘
                                         │
                                         ↓
                                  ┌─────────────┐
                                  │ Subscribers │
                                  └─────────────┘

================================================================================
*/

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// Ring Buffer 구현
// ============================================================================

type telemetryRingBuffer struct {
	data     []TelemetryData
	size     int
	head     int64 // 쓰기 위치
	tail     int64 // 읽기 위치
	mu       sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
}

func newTelemetryRingBuffer(size int) *telemetryRingBuffer {
	rb := &telemetryRingBuffer{
		data: make([]TelemetryData, size),
		size: size,
	}
	rb.notEmpty = sync.NewCond(&rb.mu)
	rb.notFull = sync.NewCond(&rb.mu)
	return rb
}

// Push 데이터 추가 (블로킹)
func (rb *telemetryRingBuffer) Push(data TelemetryData) error {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	// 버퍼가 가득 찼으면 대기
	for rb.isFull() {
		rb.notFull.Wait()
	}

	idx := rb.head % int64(rb.size)
	rb.data[idx] = data
	rb.head++

	rb.notEmpty.Signal()
	return nil
}

// TryPush 데이터 추가 (논블로킹)
func (rb *telemetryRingBuffer) TryPush(data TelemetryData) bool {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.isFull() {
		return false
	}

	idx := rb.head % int64(rb.size)
	rb.data[idx] = data
	rb.head++

	rb.notEmpty.Signal()
	return true
}

// Pop 데이터 꺼내기 (블로킹)
func (rb *telemetryRingBuffer) Pop() (TelemetryData, bool) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	for rb.isEmpty() {
		rb.notEmpty.Wait()
	}

	idx := rb.tail % int64(rb.size)
	data := rb.data[idx]
	rb.tail++

	rb.notFull.Signal()
	return data, true
}

// PopBatch 배치로 꺼내기
func (rb *telemetryRingBuffer) PopBatch(maxSize int) []TelemetryData {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	count := int(rb.head - rb.tail)
	if count == 0 {
		return nil
	}
	if count > maxSize {
		count = maxSize
	}

	result := make([]TelemetryData, count)
	for i := 0; i < count; i++ {
		idx := rb.tail % int64(rb.size)
		result[i] = rb.data[idx]
		rb.tail++
	}

	rb.notFull.Broadcast()
	return result
}

func (rb *telemetryRingBuffer) isFull() bool {
	return rb.head-rb.tail >= int64(rb.size)
}

func (rb *telemetryRingBuffer) isEmpty() bool {
	return rb.head == rb.tail
}

func (rb *telemetryRingBuffer) Len() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return int(rb.head - rb.tail)
}

func (rb *telemetryRingBuffer) Broadcast() {
	rb.mu.Lock()
	rb.notEmpty.Broadcast()
	rb.mu.Unlock()
}

// ============================================================================
// Data Pipeline 구현
// ============================================================================

// DataPipelineSolution 완성된 데이터 파이프라인
type DataPipelineSolution struct {
	// 버퍼
	ingestBuffer  *telemetryRingBuffer
	processBuffer *telemetryRingBuffer

	// 설정
	bufferSize    int
	flushInterval time.Duration
	batchSize     int
	workerCount   int

	// 구독자
	subsMu      sync.RWMutex
	subscribers map[string][]func(TelemetryData)

	// 저장소 (실제로는 DB, Kafka 등)
	storageMu sync.Mutex
	storage   []TelemetryData

	// 제어
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// 메트릭
	metrics struct {
		ingestedCount  int64
		processedCount int64
		flushedCount   int64
		errorCount     int64
		totalLatency   int64
		latencyCount   int64
	}
}

// NewDataPipelineSolution 생성자
func NewDataPipelineSolution(bufferSize int, flushInterval time.Duration) *DataPipelineSolution {
	ctx, cancel := context.WithCancel(context.Background())

	return &DataPipelineSolution{
		ingestBuffer:  newTelemetryRingBuffer(bufferSize),
		processBuffer: newTelemetryRingBuffer(bufferSize),
		bufferSize:    bufferSize,
		flushInterval: flushInterval,
		batchSize:     1000,
		workerCount:   4,
		subscribers:   make(map[string][]func(TelemetryData)),
		storage:       make([]TelemetryData, 0),
		ctx:           ctx,
		cancel:        cancel,
	}
}

// ============================================================================
// 데이터 수집
// ============================================================================

// Ingest 데이터 수집
/*
【 동작 원리 】
1. Ring Buffer에 데이터 추가
2. 버퍼가 가득 차면 back-pressure (블로킹 또는 에러)
3. 메트릭 업데이트
*/
func (p *DataPipelineSolution) Ingest(data TelemetryData) error {
	// 논블로킹으로 시도
	if !p.ingestBuffer.TryPush(data) {
		// Back-pressure: 버퍼가 가득 찼을 때
		atomic.AddInt64(&p.metrics.errorCount, 1)
		return errors.New("buffer full - back pressure applied")
	}

	atomic.AddInt64(&p.metrics.ingestedCount, 1)
	return nil
}

// IngestBlocking 블로킹 버전 (데이터 유실 방지)
func (p *DataPipelineSolution) IngestBlocking(data TelemetryData) error {
	if err := p.ingestBuffer.Push(data); err != nil {
		atomic.AddInt64(&p.metrics.errorCount, 1)
		return err
	}

	atomic.AddInt64(&p.metrics.ingestedCount, 1)
	return nil
}

// ============================================================================
// 데이터 처리
// ============================================================================

// Process 실시간 처리 (워커가 호출)
func (p *DataPipelineSolution) Process() error {
	// 배치로 데이터 가져오기
	batch := p.ingestBuffer.PopBatch(p.batchSize)
	if len(batch) == 0 {
		return nil
	}

	for _, data := range batch {
		startTime := time.Now()

		// 1. 이상 탐지
		p.detectAnomalies(data)

		// 2. 구독자에게 알림
		p.notifySubscribers(data)

		// 3. 처리 버퍼로 이동 (Flush 대기)
		p.processBuffer.TryPush(data)

		// 메트릭
		latency := time.Since(startTime)
		atomic.AddInt64(&p.metrics.totalLatency, int64(latency))
		atomic.AddInt64(&p.metrics.latencyCount, 1)
		atomic.AddInt64(&p.metrics.processedCount, 1)
	}

	return nil
}

// detectAnomalies 이상 탐지
func (p *DataPipelineSolution) detectAnomalies(data TelemetryData) {
	// 급정거 감지 (예: 속도 급격히 감소)
	// 실제로는 이전 데이터와 비교 필요

	// 배터리 이상
	if data.Sensors.Battery < 10 {
		p.notifySubscribers(TelemetryData{
			VehicleID: data.VehicleID,
			Timestamp: data.Timestamp,
			Events: []EventData{{
				Type:     "LOW_BATTERY",
				Severity: "HIGH",
				Value:    float64(data.Sensors.Battery),
			}},
		})
	}

	// 온도 이상
	if data.Sensors.Temperature > 80 {
		p.notifySubscribers(TelemetryData{
			VehicleID: data.VehicleID,
			Timestamp: data.Timestamp,
			Events: []EventData{{
				Type:     "HIGH_TEMPERATURE",
				Severity: "HIGH",
				Value:    data.Sensors.Temperature,
			}},
		})
	}
}

// notifySubscribers 구독자에게 알림
func (p *DataPipelineSolution) notifySubscribers(data TelemetryData) {
	p.subsMu.RLock()
	defer p.subsMu.RUnlock()

	// 차량별 구독자
	if handlers, ok := p.subscribers[data.VehicleID]; ok {
		for _, handler := range handlers {
			go handler(data)
		}
	}

	// 전체 구독자
	if handlers, ok := p.subscribers["*"]; ok {
		for _, handler := range handlers {
			go handler(data)
		}
	}

	// 이벤트별 구독자
	for _, event := range data.Events {
		if handlers, ok := p.subscribers[event.Type]; ok {
			for _, handler := range handlers {
				go handler(data)
			}
		}
	}
}

// ============================================================================
// 배치 저장
// ============================================================================

// Flush 배치로 저장소에 쓰기
/*
【 동작 원리 】
1. 처리 버퍼에서 배치로 데이터 가져오기
2. 저장소에 일괄 쓰기 (DB, Kafka 등)
3. 실패 시 재시도 또는 DLQ로 이동
*/
func (p *DataPipelineSolution) Flush() error {
	batch := p.processBuffer.PopBatch(p.batchSize)
	if len(batch) == 0 {
		return nil
	}

	// 저장소에 쓰기 (여기서는 메모리에 저장)
	p.storageMu.Lock()
	p.storage = append(p.storage, batch...)
	p.storageMu.Unlock()

	atomic.AddInt64(&p.metrics.flushedCount, int64(len(batch)))
	return nil
}

// ============================================================================
// 구독
// ============================================================================

// Subscribe 이벤트 구독
func (p *DataPipelineSolution) Subscribe(topic string, handler func(TelemetryData)) error {
	if handler == nil {
		return errors.New("handler is required")
	}

	p.subsMu.Lock()
	defer p.subsMu.Unlock()

	p.subscribers[topic] = append(p.subscribers[topic], handler)
	return nil
}

// ============================================================================
// 라이프사이클
// ============================================================================

// Start 파이프라인 시작
func (p *DataPipelineSolution) Start() error {
	// 처리 워커 시작
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.processWorker()
	}

	// Flush 워커 시작
	p.wg.Add(1)
	go p.flushWorker()

	return nil
}

// processWorker 처리 워커
func (p *DataPipelineSolution) processWorker() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		default:
			p.Process()
			// 바쁜 대기 방지
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// flushWorker Flush 워커
func (p *DataPipelineSolution) flushWorker() {
	defer p.wg.Done()

	ticker := time.NewTicker(p.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			// 종료 전 마지막 Flush
			p.Flush()
			return
		case <-ticker.C:
			p.Flush()
		}
	}
}

// Stop 파이프라인 정지
func (p *DataPipelineSolution) Stop() error {
	p.cancel()

	// 버퍼 대기 중인 고루틴 깨우기
	p.ingestBuffer.Broadcast()
	p.processBuffer.Broadcast()

	p.wg.Wait()
	return nil
}

// GetMetrics 메트릭 반환
func (p *DataPipelineSolution) GetMetrics() PipelineMetrics {
	metrics := PipelineMetrics{
		IngestedCount:  atomic.LoadInt64(&p.metrics.ingestedCount),
		ProcessedCount: atomic.LoadInt64(&p.metrics.processedCount),
		FlushedCount:   atomic.LoadInt64(&p.metrics.flushedCount),
		ErrorCount:     atomic.LoadInt64(&p.metrics.errorCount),
		BufferSize:     p.ingestBuffer.Len() + p.processBuffer.Len(),
	}

	latencyCount := atomic.LoadInt64(&p.metrics.latencyCount)
	if latencyCount > 0 {
		totalLatency := atomic.LoadInt64(&p.metrics.totalLatency)
		metrics.AvgLatency = time.Duration(totalLatency / latencyCount)
	}

	return metrics
}

// ============================================================================
// 성능 최적화 포인트
// ============================================================================

/*
【 프로덕션 최적화 】

1. 메모리 풀 (sync.Pool)
   - TelemetryData 객체 재사용
   - GC 압박 감소

var telemetryPool = sync.Pool{
    New: func() interface{} {
        return &TelemetryData{}
    },
}

2. 무잠금 큐 (Lock-free Queue)
   - atomic 연산으로 락 없이 동작
   - 처리량 10배 이상 향상 가능

3. 배치 압축
   - Snappy, LZ4로 배치 압축
   - 네트워크/디스크 I/O 감소

4. 파티셔닝
   - 차량 ID로 파티션 분배
   - 각 파티션 독립 처리

5. Checkpoint
   - 처리 위치 주기적 저장
   - 장애 시 정확한 위치에서 재시작

┌─────────────────────────────────────────────────────────────────────────────┐
│                       프로덕션 아키텍처                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  [차량 10만대] → [API Gateway] → [Kafka (파티션 x 32)]                      │
│                                          ↓                                  │
│                              ┌───────────────────────┐                      │
│                              │   Consumer Group      │                      │
│                              │  ┌─────┐ ┌─────┐     │                      │
│                              │  │ W1  │ │ W2  │ ... │                      │
│                              │  └─────┘ └─────┘     │                      │
│                              └───────────────────────┘                      │
│                                          ↓                                  │
│                    ┌─────────────────────┼─────────────────────┐           │
│                    ↓                     ↓                     ↓           │
│             [TimescaleDB]         [ClickHouse]          [Redis]            │
│             (시계열 저장)         (분석 쿼리)          (실시간 캐시)        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

【 Kafka 설정 예시 】

토픽: vehicle-telemetry
- 파티션: 32 (차량 ID 해시로 분배)
- 복제 팩터: 3
- 보존 기간: 7일

Consumer 설정:
- enable.auto.commit: false (수동 커밋)
- max.poll.records: 1000 (배치 크기)
- fetch.min.bytes: 1MB (배치 수집)
*/
