package huhu1

/*
================================================================================
문제 8: 이벤트 소싱 시스템
================================================================================

난이도: Hard
주제: 이벤트 소싱, CQRS, 이벤트 스토어
회사: 42dot (백엔드/데이터 플랫폼)

【 문제 설명 】
차량의 모든 상태 변화를 이벤트로 저장하고, 이벤트를 재생하여
언제든지 특정 시점의 상태를 복원할 수 있는 이벤트 소싱 시스템을 구현하세요.

이벤트 예시:
- VehicleRegistered: 차량 등록
- LocationUpdated: 위치 업데이트
- TripStarted: 운행 시작
- TripEnded: 운행 종료
- MaintenanceScheduled: 정비 예약
- BatteryLevelChanged: 배터리 잔량 변경

요구사항:
1. Append(event): 이벤트 추가 (불변)
2. GetEvents(aggregateID, fromVersion): 이벤트 조회
3. GetSnapshot(aggregateID): 현재 상태 스냅샷
4. ReplayTo(aggregateID, timestamp): 특정 시점으로 상태 복원
5. Subscribe(eventType, handler): 이벤트 구독
6. CreateSnapshot(aggregateID): 스냅샷 생성 (성능 최적화)

【 성능 요구사항 】
- 이벤트 추가: 1ms 이내
- 상태 복원: 100ms 이내 (스냅샷 활용)
- 이벤트 저장: 영속성 보장

【 인터페이스 】
type Event struct {
    ID          string
    AggregateID string    // 차량 ID
    Type        string    // 이벤트 타입
    Version     int64     // 버전 (순서 보장)
    Timestamp   time.Time
    Data        map[string]interface{}
    Metadata    map[string]string
}

type Snapshot struct {
    AggregateID string
    Version     int64
    State       map[string]interface{}
    Timestamp   time.Time
}

type VehicleState struct {
    VehicleID    string
    Status       string  // "ACTIVE", "INACTIVE", "MAINTENANCE"
    Location     *GPSData
    Battery      int
    TotalTrips   int
    LastUpdated  time.Time
}

type EventStore interface {
    Append(event Event) error
    GetEvents(aggregateID string, fromVersion int64) ([]Event, error)
    GetSnapshot(aggregateID string) (*Snapshot, error)
    ReplayTo(aggregateID string, timestamp time.Time) (*VehicleState, error)
    Subscribe(eventType string, handler func(Event)) error
    CreateSnapshot(aggregateID string) error
}

【 힌트 】
1. 이벤트는 불변 (Append-only)
2. 버전으로 순서 보장 (낙관적 동시성)
3. 스냅샷으로 재생 시간 단축
4. 이벤트 핸들러로 Read Model 업데이트

【 실무 연관성 】
- 감사 로그 (Audit Trail)
- 상태 복원 및 디버깅
- CQRS 패턴의 Write Side
================================================================================
*/

import "time"

// Event 이벤트
type Event struct {
	ID          string
	AggregateID string
	Type        string
	Version     int64
	Timestamp   time.Time
	Data        map[string]any
	Metadata    map[string]string
}

// Snapshot 스냅샷
type Snapshot struct {
	AggregateID string
	Version     int64
	State       map[string]any
	Timestamp   time.Time
}

// VehicleState 차량 상태 (Read Model)
type VehicleStateES struct {
	VehicleID   string
	Status      string
	Location    *GPSData
	Battery     int
	TotalTrips  int
	LastUpdated time.Time
}

// EventStore 이벤트 스토어 인터페이스
type EventStore interface {
	Append(event Event) error
	GetEvents(aggregateID string, fromVersion int64) ([]Event, error)
	GetSnapshot(aggregateID string) (*Snapshot, error)
	ReplayTo(aggregateID string, timestamp time.Time) (*VehicleStateES, error)
	Subscribe(eventType string, handler func(Event)) error
	CreateSnapshot(aggregateID string) error
}

// EventStoreImpl 구현체
type EventStoreImpl struct {
	// 여기에 필드를 정의하세요
}

// NewEventStore 생성자
func NewEventStore() *EventStoreImpl {
	// 여기에 코드를 작성하세요
	return nil
}

// Append 이벤트 추가
func (s *EventStoreImpl) Append(event Event) error {
	// 여기에 코드를 작성하세요
	return nil
}

// GetEvents 이벤트 조회
func (s *EventStoreImpl) GetEvents(aggregateID string, fromVersion int64) ([]Event, error) {
	// 여기에 코드를 작성하세요
	return nil, nil
}

// GetSnapshot 스냅샷 조회
func (s *EventStoreImpl) GetSnapshot(aggregateID string) (*Snapshot, error) {
	// 여기에 코드를 작성하세요
	return nil, nil
}

// ReplayTo 특정 시점으로 상태 복원
func (s *EventStoreImpl) ReplayTo(aggregateID string, timestamp time.Time) (*VehicleStateES, error) {
	// 여기에 코드를 작성하세요
	return nil, nil
}

// Subscribe 이벤트 구독
func (s *EventStoreImpl) Subscribe(eventType string, handler func(Event)) error {
	// 여기에 코드를 작성하세요
	return nil
}

// CreateSnapshot 스냅샷 생성
func (s *EventStoreImpl) CreateSnapshot(aggregateID string) error {
	// 여기에 코드를 작성하세요
	return nil
}
