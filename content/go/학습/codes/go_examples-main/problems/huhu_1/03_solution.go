package huhu1

/*
================================================================================
문제 3: V2X 메시지 처리 시스템 - 솔루션
================================================================================

【 핵심 설계 원칙 】
1. 우선순위 큐 - EMERGENCY 메시지 최우선 처리
2. 워커 풀 - 병렬 처리로 처리량 극대화
3. Pub/Sub - 효율적인 메시지 전달
4. 공간 인덱싱 - 근처 차량에만 브로드캐스트

【 우선순위 체계 】
EMERGENCY (0) > SAFETY (1) > TRAFFIC (2) > INFO (3)
낮은 숫자 = 높은 우선순위

================================================================================
*/

import (
	"container/heap"
	"context"
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// 우선순위 큐 구현
// ============================================================================

func getMessagePriority(msgType MessageType) int {
	switch msgType {
	case EMERGENCY:
		return 0
	case SAFETY:
		return 1
	case TRAFFIC:
		return 2
	case INFO:
		return 3
	default:
		return 4
	}
}

type messageItem struct {
	msg      V2XMessage
	priority int
	index    int
}

type messagePriorityQueue []*messageItem

func (pq messagePriorityQueue) Len() int { return len(pq) }

func (pq messagePriorityQueue) Less(i, j int) bool {
	// 우선순위가 같으면 먼저 도착한 메시지 우선 (타임스탬프)
	if pq[i].priority == pq[j].priority {
		return pq[i].msg.Timestamp.Before(pq[j].msg.Timestamp)
	}
	return pq[i].priority < pq[j].priority
}

func (pq messagePriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *messagePriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*messageItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *messagePriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// ============================================================================
// 구독자 관리
// ============================================================================

type subscription struct {
	vehicleID string
	msgType   MessageType
	handler   func(V2XMessage)
	lat, lon  float64 // 차량 위치 (Publish 시 필터링용)
}

// ============================================================================
// V2X 프로세서 구현
// ============================================================================

type V2XProcessorSolution struct {
	// 메시지 큐
	queueMu sync.Mutex
	queue   messagePriorityQueue
	cond    *sync.Cond

	// 구독자
	subsMu       sync.RWMutex
	subscribers  map[MessageType]map[string]*subscription // msgType -> vehicleID -> subscription
	vehicleLocs  map[string]struct{ lat, lon float64 }   // 차량 위치

	// 워커 풀
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// 통계
	stats struct {
		totalProcessed int64
		byType         sync.Map // MessageType -> *int64
		totalLatency   int64    // 나노초 합계
		latencyCount   int64
		errorCount     int64
	}
}

// NewV2XProcessorSolution 생성자
func NewV2XProcessorSolution(workerCount int) *V2XProcessorSolution {
	ctx, cancel := context.WithCancel(context.Background())

	p := &V2XProcessorSolution{
		subscribers: make(map[MessageType]map[string]*subscription),
		vehicleLocs: make(map[string]struct{ lat, lon float64 }),
		ctx:         ctx,
		cancel:      cancel,
	}

	p.queue = make(messagePriorityQueue, 0)
	heap.Init(&p.queue)
	p.cond = sync.NewCond(&p.queueMu)

	// 메시지 타입별 카운터 초기화
	for _, mt := range []MessageType{EMERGENCY, SAFETY, TRAFFIC, INFO} {
		var count int64
		p.stats.byType.Store(mt, &count)
	}

	// 워커 시작
	for i := 0; i < workerCount; i++ {
		p.wg.Add(1)
		go p.worker()
	}

	return p
}

// worker 메시지 처리 워커
func (p *V2XProcessorSolution) worker() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		default:
		}

		p.queueMu.Lock()

		// 큐가 비어있으면 대기
		for p.queue.Len() == 0 {
			p.cond.Wait()

			// 컨텍스트 취소 확인
			select {
			case <-p.ctx.Done():
				p.queueMu.Unlock()
				return
			default:
			}
		}

		// 메시지 꺼내기
		item := heap.Pop(&p.queue).(*messageItem)
		p.queueMu.Unlock()

		// 메시지 처리
		startTime := time.Now()
		p.handleMessage(item.msg)
		latency := time.Since(startTime)

		// 통계 업데이트
		atomic.AddInt64(&p.stats.totalProcessed, 1)
		if counter, ok := p.stats.byType.Load(item.msg.Type); ok {
			atomic.AddInt64(counter.(*int64), 1)
		}
		atomic.AddInt64(&p.stats.totalLatency, int64(latency))
		atomic.AddInt64(&p.stats.latencyCount, 1)
	}
}

// handleMessage 실제 메시지 처리
func (p *V2XProcessorSolution) handleMessage(msg V2XMessage) {
	// TTL 확인
	if time.Since(msg.Timestamp) > msg.TTL {
		return // 만료된 메시지
	}

	// 해당 타입 구독자들에게 전달
	p.subsMu.RLock()
	subs := p.subscribers[msg.Type]
	if subs != nil {
		for _, sub := range subs {
			// 고루틴으로 비동기 전달 (핸들러 블로킹 방지)
			go func(s *subscription) {
				defer func() {
					if r := recover(); r != nil {
						atomic.AddInt64(&p.stats.errorCount, 1)
					}
				}()
				s.handler(msg)
			}(sub)
		}
	}
	p.subsMu.RUnlock()
}

// ProcessMessage 메시지 처리 요청
/*
【 시간 복잡도 】O(log n) - 힙 삽입
【 동작 원리 】
1. 우선순위 계산
2. 우선순위 큐에 삽입
3. 워커에게 신호
*/
func (p *V2XProcessorSolution) ProcessMessage(msg V2XMessage) error {
	if msg.ID == "" {
		return errors.New("message ID is required")
	}

	priority := getMessagePriority(msg.Type)

	p.queueMu.Lock()
	heap.Push(&p.queue, &messageItem{
		msg:      msg,
		priority: priority,
	})
	p.queueMu.Unlock()

	// 워커 깨우기
	p.cond.Signal()

	return nil
}

// Subscribe 메시지 타입 구독
func (p *V2XProcessorSolution) Subscribe(vehicleID string, msgType MessageType, handler func(V2XMessage)) error {
	if vehicleID == "" {
		return errors.New("vehicle ID is required")
	}
	if handler == nil {
		return errors.New("handler is required")
	}

	p.subsMu.Lock()
	defer p.subsMu.Unlock()

	if p.subscribers[msgType] == nil {
		p.subscribers[msgType] = make(map[string]*subscription)
	}

	p.subscribers[msgType][vehicleID] = &subscription{
		vehicleID: vehicleID,
		msgType:   msgType,
		handler:   handler,
	}

	return nil
}

// Unsubscribe 구독 해제
func (p *V2XProcessorSolution) Unsubscribe(vehicleID string, msgType MessageType) error {
	p.subsMu.Lock()
	defer p.subsMu.Unlock()

	if p.subscribers[msgType] != nil {
		delete(p.subscribers[msgType], vehicleID)
	}

	return nil
}

// UpdateVehicleLocation 차량 위치 업데이트 (Publish 필터링용)
func (p *V2XProcessorSolution) UpdateVehicleLocation(vehicleID string, lat, lon float64) {
	p.subsMu.Lock()
	p.vehicleLocs[vehicleID] = struct{ lat, lon float64 }{lat, lon}
	p.subsMu.Unlock()
}

// Publish 근처 차량들에게 메시지 브로드캐스트
/*
【 동작 원리 】
1. 메시지 송신자 위치 기준
2. 반경 내 구독자만 필터링
3. 해당 타입 구독자들에게 전달
*/
func (p *V2XProcessorSolution) Publish(msg V2XMessage, radiusKm float64) error {
	if msg.ID == "" {
		return errors.New("message ID is required")
	}

	p.subsMu.RLock()
	defer p.subsMu.RUnlock()

	subs := p.subscribers[msg.Type]
	if subs == nil {
		return nil
	}

	for vehicleID, sub := range subs {
		// 자기 자신 제외
		if vehicleID == msg.SenderID {
			continue
		}

		// 거리 확인
		loc, ok := p.vehicleLocs[vehicleID]
		if !ok {
			continue
		}

		dist := haversineDistanceV2X(msg.Lat, msg.Lon, loc.lat, loc.lon)
		if dist <= radiusKm {
			// 비동기 전달
			go func(s *subscription, m V2XMessage) {
				defer func() {
					if r := recover(); r != nil {
						atomic.AddInt64(&p.stats.errorCount, 1)
					}
				}()
				s.handler(m)
			}(sub, msg)
		}
	}

	return nil
}

func haversineDistanceV2X(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// GetStatistics 통계 반환
func (p *V2XProcessorSolution) GetStatistics() Statistics {
	stats := Statistics{
		TotalProcessed: atomic.LoadInt64(&p.stats.totalProcessed),
		ByType:         make(map[MessageType]int64),
		ErrorCount:     atomic.LoadInt64(&p.stats.errorCount),
	}

	// 타입별 통계
	for _, mt := range []MessageType{EMERGENCY, SAFETY, TRAFFIC, INFO} {
		if counter, ok := p.stats.byType.Load(mt); ok {
			stats.ByType[mt] = atomic.LoadInt64(counter.(*int64))
		}
	}

	// 평균 지연시간
	latencyCount := atomic.LoadInt64(&p.stats.latencyCount)
	if latencyCount > 0 {
		totalLatency := atomic.LoadInt64(&p.stats.totalLatency)
		stats.AvgLatency = time.Duration(totalLatency / latencyCount)
	}

	return stats
}

// Stop 프로세서 정지
func (p *V2XProcessorSolution) Stop() {
	p.cancel()
	p.cond.Broadcast() // 대기 중인 워커 깨우기
	p.wg.Wait()
}

// ============================================================================
// 성능 최적화 포인트
// ============================================================================

/*
【 EMERGENCY 메시지 최적화 】
현재: 일반 큐에서 우선순위로 처리
개선: 별도 고속 채널로 직접 처리

func (p *V2XProcessor) ProcessMessage(msg V2XMessage) error {
    if msg.Type == EMERGENCY {
        // 큐를 거치지 않고 즉시 처리
        go p.handleEmergency(msg)
        return nil
    }
    // 일반 메시지는 큐로
    ...
}

【 공간 인덱싱 최적화 】
현재: 모든 구독자와 거리 계산
개선: Geohash 기반 필터링 후 거리 계산

【 메시지 배치 처리 】
- 낮은 우선순위 메시지는 배치로 모아서 처리
- 네트워크 오버헤드 감소

【 메시지 압축 】
- Protobuf 또는 MessagePack으로 직렬화
- 대역폭 40-60% 절약

【 중복 제거 】
- 메시지 ID로 중복 감지
- Bloom Filter로 메모리 효율적 중복 체크

┌─────────────────────────────────────────────────────────────────────────────┐
│                       프로덕션 아키텍처                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  [차량 V2X 모듈] ←──→ [MQTT Broker] ←──→ [V2X Gateway]                     │
│                             ↓                                               │
│                    [Kafka (메시지 스트림)]                                   │
│                             ↓                                               │
│  ┌──────────────────────────────────────────────────────────┐              │
│  │                    V2X Processor                         │              │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │              │
│  │  │Emergency│ │ Safety  │ │ Traffic │ │  Info   │       │              │
│  │  │ Queue   │ │ Queue   │ │ Queue   │ │ Queue   │       │              │
│  │  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘       │              │
│  │       ↓           ↓           ↓           ↓            │              │
│  │  ┌────────────────────────────────────────────┐        │              │
│  │  │              Worker Pool                    │        │              │
│  │  └────────────────────────────────────────────┘        │              │
│  └──────────────────────────────────────────────────────────┘              │
│                             ↓                                               │
│                    [Pub/Sub → 차량들]                                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

【 V2X 메시지 예시 】

EMERGENCY (충돌 경고):
{
    "id": "emg-001",
    "type": "EMERGENCY",
    "sender": "vehicle-123",
    "lat": 37.5665, "lon": 126.9780,
    "payload": {
        "event": "COLLISION_WARNING",
        "severity": "HIGH",
        "description": "전방 50m 급정거 차량"
    },
    "ttl": "5s"
}

SAFETY (신호등 정보):
{
    "id": "safety-001",
    "type": "SAFETY",
    "sender": "traffic-light-42",
    "lat": 37.5660, "lon": 126.9775,
    "payload": {
        "signal": "RED",
        "remaining_seconds": 15,
        "next_signal": "GREEN"
    },
    "ttl": "30s"
}
*/
