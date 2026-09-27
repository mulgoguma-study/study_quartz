package huhu1

/*
================================================================================
문제 8: 이벤트 소싱 시스템 - 솔루션
================================================================================

【 핵심 설계 원칙 】
1. Append-only - 이벤트는 절대 수정/삭제 불가
2. 낙관적 동시성 - 버전 충돌 감지
3. 스냅샷 최적화 - 재생 시간 단축
4. 이벤트 핸들러 - Read Model 비동기 업데이트

【 이벤트 소싱 흐름 】
┌─────────────────────────────────────────────────────────────────┐
│  Command → Validate → Append Event → Publish → Update Read     │
└─────────────────────────────────────────────────────────────────┘

【 CQRS 패턴 】
Write Side: Event Store (이벤트 저장)
Read Side: Read Model (쿼리 최적화된 뷰)

================================================================================
*/

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// 이벤트 스토어 구현
// ============================================================================

// EventStoreSolution 완성된 이벤트 스토어
type EventStoreSolution struct {
	mu sync.RWMutex

	// 이벤트 저장소 (aggregateID -> events)
	events map[string][]Event

	// 스냅샷 저장소 (aggregateID -> snapshot)
	snapshots map[string]*Snapshot

	// 현재 버전 (aggregateID -> version)
	versions map[string]int64

	// 구독자 (eventType -> handlers)
	subsMu      sync.RWMutex
	subscribers map[string][]func(Event)

	// 스냅샷 설정
	snapshotInterval int64 // N개 이벤트마다 스냅샷
}

// NewEventStoreSolution 생성자
func NewEventStoreSolution() *EventStoreSolution {
	return &EventStoreSolution{
		events:           make(map[string][]Event),
		snapshots:        make(map[string]*Snapshot),
		versions:         make(map[string]int64),
		subscribers:      make(map[string][]func(Event)),
		snapshotInterval: 100, // 100개 이벤트마다 스냅샷
	}
}

// ============================================================================
// 이벤트 추가
// ============================================================================

// Append 이벤트 추가
/*
【 동작 원리 】
1. 버전 검증 (낙관적 동시성)
2. 이벤트 ID 생성
3. 이벤트 저장 (Append-only)
4. 구독자에게 알림
5. 스냅샷 필요 여부 확인
*/
func (s *EventStoreSolution) Append(event Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 버전 검증
	currentVersion := s.versions[event.AggregateID]
	if event.Version != 0 && event.Version != currentVersion+1 {
		return errors.New("version conflict: expected " + string(rune(currentVersion+1)))
	}

	// 이벤트 설정
	event.ID = uuid.New().String()
	event.Version = currentVersion + 1
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// 저장
	s.events[event.AggregateID] = append(s.events[event.AggregateID], event)
	s.versions[event.AggregateID] = event.Version

	// 구독자 알림 (비동기)
	go s.notifySubscribers(event)

	// 스냅샷 필요 여부 확인
	if event.Version%s.snapshotInterval == 0 {
		go s.createSnapshotInternal(event.AggregateID)
	}

	return nil
}

// notifySubscribers 구독자에게 알림
func (s *EventStoreSolution) notifySubscribers(event Event) {
	s.subsMu.RLock()
	defer s.subsMu.RUnlock()

	// 특정 이벤트 타입 구독자
	if handlers, ok := s.subscribers[event.Type]; ok {
		for _, handler := range handlers {
			go handler(event)
		}
	}

	// 전체 이벤트 구독자
	if handlers, ok := s.subscribers["*"]; ok {
		for _, handler := range handlers {
			go handler(event)
		}
	}
}

// ============================================================================
// 이벤트 조회
// ============================================================================

// GetEvents 이벤트 조회
func (s *EventStoreSolution) GetEvents(aggregateID string, fromVersion int64) ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, ok := s.events[aggregateID]
	if !ok {
		return []Event{}, nil
	}

	result := make([]Event, 0)
	for _, event := range events {
		if event.Version >= fromVersion {
			result = append(result, event)
		}
	}

	return result, nil
}

// ============================================================================
// 스냅샷
// ============================================================================

// GetSnapshot 스냅샷 조회
func (s *EventStoreSolution) GetSnapshot(aggregateID string) (*Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot, ok := s.snapshots[aggregateID]
	if !ok {
		return nil, nil
	}

	return snapshot, nil
}

// CreateSnapshot 스냅샷 생성
func (s *EventStoreSolution) CreateSnapshot(aggregateID string) error {
	return s.createSnapshotInternal(aggregateID)
}

// createSnapshotInternal 스냅샷 생성 (내부)
func (s *EventStoreSolution) createSnapshotInternal(aggregateID string) error {
	// 현재 상태 계산
	state, err := s.replayEvents(aggregateID, time.Now())
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.snapshots[aggregateID] = &Snapshot{
		AggregateID: aggregateID,
		Version:     s.versions[aggregateID],
		State:       stateToMap(state),
		Timestamp:   time.Now(),
	}

	return nil
}

// stateToMap VehicleState를 map으로 변환
func stateToMap(state *VehicleStateES) map[string]any {
	if state == nil {
		return nil
	}

	m := map[string]any{
		"vehicleID":   state.VehicleID,
		"status":      state.Status,
		"battery":     state.Battery,
		"totalTrips":  state.TotalTrips,
		"lastUpdated": state.LastUpdated,
	}

	if state.Location != nil {
		m["location"] = map[string]any{
			"lat":     state.Location.Lat,
			"lon":     state.Location.Lon,
			"speed":   state.Location.Speed,
			"heading": state.Location.Heading,
		}
	}

	return m
}

// ============================================================================
// 상태 복원
// ============================================================================

// ReplayTo 특정 시점으로 상태 복원
/*
【 동작 원리 】
1. 스냅샷 확인 (있으면 스냅샷부터 시작)
2. 스냅샷 이후 이벤트만 재생
3. timestamp까지의 이벤트만 적용
*/
func (s *EventStoreSolution) ReplayTo(aggregateID string, timestamp time.Time) (*VehicleStateES, error) {
	return s.replayEvents(aggregateID, timestamp)
}

// replayEvents 이벤트 재생
func (s *EventStoreSolution) replayEvents(aggregateID string, until time.Time) (*VehicleStateES, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 초기 상태
	state := &VehicleStateES{
		VehicleID: aggregateID,
		Status:    "UNKNOWN",
	}

	var fromVersion int64 = 0

	// 스냅샷 확인
	if snapshot, ok := s.snapshots[aggregateID]; ok && snapshot.Timestamp.Before(until) {
		state = mapToState(snapshot.State)
		fromVersion = snapshot.Version
	}

	// 이벤트 재생
	events, ok := s.events[aggregateID]
	if !ok {
		return state, nil
	}

	for _, event := range events {
		if event.Version <= fromVersion {
			continue
		}
		if event.Timestamp.After(until) {
			break
		}

		state = applyEvent(state, event)
	}

	return state, nil
}

// mapToState map을 VehicleState로 변환
func mapToState(m map[string]any) *VehicleStateES {
	state := &VehicleStateES{}

	if v, ok := m["vehicleID"].(string); ok {
		state.VehicleID = v
	}
	if v, ok := m["status"].(string); ok {
		state.Status = v
	}
	if v, ok := m["battery"].(int); ok {
		state.Battery = v
	}
	if v, ok := m["totalTrips"].(int); ok {
		state.TotalTrips = v
	}
	if v, ok := m["lastUpdated"].(time.Time); ok {
		state.LastUpdated = v
	}

	if loc, ok := m["location"].(map[string]any); ok {
		state.Location = &GPSData{}
		if v, ok := loc["lat"].(float64); ok {
			state.Location.Lat = v
		}
		if v, ok := loc["lon"].(float64); ok {
			state.Location.Lon = v
		}
		if v, ok := loc["speed"].(float64); ok {
			state.Location.Speed = v
		}
		if v, ok := loc["heading"].(float64); ok {
			state.Location.Heading = v
		}
	}

	return state
}

// applyEvent 이벤트를 상태에 적용
func applyEvent(state *VehicleStateES, event Event) *VehicleStateES {
	state.LastUpdated = event.Timestamp

	switch event.Type {
	case "VehicleRegistered":
		state.Status = "ACTIVE"
		if name, ok := event.Data["vehicleID"].(string); ok {
			state.VehicleID = name
		}

	case "LocationUpdated":
		if state.Location == nil {
			state.Location = &GPSData{}
		}
		if lat, ok := event.Data["lat"].(float64); ok {
			state.Location.Lat = lat
		}
		if lon, ok := event.Data["lon"].(float64); ok {
			state.Location.Lon = lon
		}
		if speed, ok := event.Data["speed"].(float64); ok {
			state.Location.Speed = speed
		}

	case "TripStarted":
		state.Status = "IN_TRIP"

	case "TripEnded":
		state.Status = "ACTIVE"
		state.TotalTrips++

	case "MaintenanceScheduled":
		state.Status = "MAINTENANCE"

	case "BatteryLevelChanged":
		if battery, ok := event.Data["battery"].(int); ok {
			state.Battery = battery
		}
		if battery, ok := event.Data["battery"].(float64); ok {
			state.Battery = int(battery)
		}

	case "VehicleDeactivated":
		state.Status = "INACTIVE"
	}

	return state
}

// ============================================================================
// 구독
// ============================================================================

// Subscribe 이벤트 구독
func (s *EventStoreSolution) Subscribe(eventType string, handler func(Event)) error {
	if handler == nil {
		return errors.New("handler is required")
	}

	s.subsMu.Lock()
	defer s.subsMu.Unlock()

	s.subscribers[eventType] = append(s.subscribers[eventType], handler)
	return nil
}

// ============================================================================
// 성능 최적화 포인트
// ============================================================================

/*
【 프로덕션 최적화 】

1. 이벤트 저장소 선택
   - EventStoreDB: 전용 이벤트 스토어
   - PostgreSQL + Outbox: 트랜잭션 보장
   - Kafka: 대용량 스트리밍

2. 스냅샷 전략
   - 시간 기반: 1시간마다
   - 버전 기반: 100 이벤트마다
   - 크기 기반: 상태 크기가 클 때

3. 프로젝션 (Read Model)
   - 각 쿼리 패턴에 최적화된 뷰
   - 비동기로 업데이트
   - 여러 프로젝션 병렬 유지

4. 이벤트 업그레이드
   - 이벤트 스키마 버저닝
   - Upcaster로 구 버전 이벤트 변환

┌─────────────────────────────────────────────────────────────────────────────┐
│                       CQRS + Event Sourcing 아키텍처                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  [Command Side]                          [Query Side]                       │
│  ┌─────────────┐                        ┌─────────────┐                    │
│  │   Command   │                        │   Query     │                    │
│  │   Handler   │                        │   Handler   │                    │
│  └──────┬──────┘                        └──────┬──────┘                    │
│         │                                      │                            │
│         ↓                                      ↓                            │
│  ┌─────────────┐                        ┌─────────────┐                    │
│  │  Aggregate  │                        │  Read Model │                    │
│  │   (Domain)  │                        │   (View)    │                    │
│  └──────┬──────┘                        └──────┬──────┘                    │
│         │                                      ↑                            │
│         ↓                                      │                            │
│  ┌─────────────┐      Publish           ┌─────────────┐                    │
│  │   Event     │ ──────────────────────→│  Projector  │                    │
│  │   Store     │                        │             │                    │
│  └─────────────┘                        └─────────────┘                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

【 PostgreSQL 이벤트 스토어 스키마 】

CREATE TABLE events (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_id  VARCHAR(255) NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    event_type    VARCHAR(100) NOT NULL,
    version       BIGINT NOT NULL,
    data          JSONB NOT NULL,
    metadata      JSONB,
    timestamp     TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (aggregate_id, version)  -- 낙관적 동시성
);

CREATE TABLE snapshots (
    aggregate_id   VARCHAR(255) PRIMARY KEY,
    version        BIGINT NOT NULL,
    state          JSONB NOT NULL,
    timestamp      TIMESTAMPTZ DEFAULT NOW()
);

-- 이벤트 조회 인덱스
CREATE INDEX idx_events_aggregate ON events(aggregate_id, version);
CREATE INDEX idx_events_type ON events(event_type, timestamp);

【 Outbox 패턴 (트랜잭션 보장) 】

BEGIN;
  -- 비즈니스 로직 처리
  INSERT INTO vehicles (id, status) VALUES ('v1', 'ACTIVE');

  -- 이벤트 저장 (같은 트랜잭션)
  INSERT INTO outbox (aggregate_id, event_type, data)
  VALUES ('v1', 'VehicleRegistered', '{"status": "ACTIVE"}');
COMMIT;

-- 별도 프로세스가 outbox → Kafka 발행
-- 발행 후 outbox에서 삭제
*/
