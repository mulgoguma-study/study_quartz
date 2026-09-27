package huhu1

/*
================================================================================
문제 3: V2X 메시지 처리 시스템
================================================================================

난이도: Hard
주제: 실시간 스트리밍, 동시성, 메시지 우선순위
회사: 42dot (자율주행/모빌리티)

【 문제 설명 】
V2X(Vehicle-to-Everything) 통신 메시지를 처리하는 시스템을 구현하세요.
차량은 다른 차량(V2V), 인프라(V2I), 보행자(V2P)와 실시간 통신합니다.

메시지 유형:
- EMERGENCY: 긴급 상황 (충돌 경고, 급정거) - 최우선 처리
- SAFETY: 안전 관련 (신호등, 속도 제한) - 높은 우선순위
- TRAFFIC: 교통 정보 (정체, 사고) - 중간 우선순위
- INFO: 일반 정보 (날씨, 주차) - 낮은 우선순위

요구사항:
1. ProcessMessage(msg): 메시지 처리 (우선순위 기반)
2. Subscribe(vehicleID, msgType): 특정 메시지 타입 구독
3. Publish(msg): 근처 차량들에게 메시지 브로드캐스트
4. GetStatistics(): 처리 통계 반환

【 성능 요구사항 】
- 초당 100,000건 메시지 처리
- EMERGENCY 메시지는 10ms 이내 전달
- 메시지 유실 없음 (at-least-once)

【 인터페이스 】
type MessageType string

const (
    EMERGENCY MessageType = "EMERGENCY"
    SAFETY    MessageType = "SAFETY"
    TRAFFIC   MessageType = "TRAFFIC"
    INFO      MessageType = "INFO"
)

type V2XMessage struct {
    ID        string
    Type      MessageType
    SenderID  string
    Lat, Lon  float64
    Payload   map[string]interface{}
    Timestamp time.Time
    TTL       time.Duration  // Time-to-Live
}

type V2XProcessor interface {
    ProcessMessage(msg V2XMessage) error
    Subscribe(vehicleID string, msgType MessageType, handler func(V2XMessage)) error
    Unsubscribe(vehicleID string, msgType MessageType) error
    Publish(msg V2XMessage, radiusKm float64) error
    GetStatistics() Statistics
}

type Statistics struct {
    TotalProcessed int64
    ByType         map[MessageType]int64
    AvgLatency     time.Duration
    ErrorCount     int64
}

【 힌트 】
1. 우선순위 큐로 메시지 처리 순서 관리
2. 워커 풀로 병렬 처리
3. Pub/Sub 패턴으로 메시지 전달
4. 공간 인덱싱으로 근처 차량 조회

【 실무 연관성 】
- 자율주행 차량 간 긴급 메시지 교환
- 스마트 교차로 신호 시스템
- 실시간 교통 정보 공유
================================================================================
*/

import "time"

// MessageType 메시지 유형
type MessageType string

const (
	EMERGENCY MessageType = "EMERGENCY"
	SAFETY    MessageType = "SAFETY"
	TRAFFIC   MessageType = "TRAFFIC"
	INFO      MessageType = "INFO"
)

// V2XMessage V2X 메시지
type V2XMessage struct {
	ID        string
	Type      MessageType
	SenderID  string
	Lat, Lon  float64
	Payload   map[string]interface{}
	Timestamp time.Time
	TTL       time.Duration
}

// Statistics 처리 통계
type Statistics struct {
	TotalProcessed int64
	ByType         map[MessageType]int64
	AvgLatency     time.Duration
	ErrorCount     int64
}

// V2XProcessor V2X 메시지 처리 인터페이스
type V2XProcessor interface {
	ProcessMessage(msg V2XMessage) error
	Subscribe(vehicleID string, msgType MessageType, handler func(V2XMessage)) error
	Unsubscribe(vehicleID string, msgType MessageType) error
	Publish(msg V2XMessage, radiusKm float64) error
	GetStatistics() Statistics
}

// V2XProcessorImpl 구현체
type V2XProcessorImpl struct {
	// 여기에 필드를 정의하세요
}

// NewV2XProcessor 생성자
func NewV2XProcessor() *V2XProcessorImpl {
	// 여기에 코드를 작성하세요
	return nil
}

// ProcessMessage 메시지 처리
func (p *V2XProcessorImpl) ProcessMessage(msg V2XMessage) error {
	// 여기에 코드를 작성하세요
	return nil
}

// Subscribe 메시지 타입 구독
func (p *V2XProcessorImpl) Subscribe(vehicleID string, msgType MessageType, handler func(V2XMessage)) error {
	// 여기에 코드를 작성하세요
	return nil
}

// Unsubscribe 구독 해제
func (p *V2XProcessorImpl) Unsubscribe(vehicleID string, msgType MessageType) error {
	// 여기에 코드를 작성하세요
	return nil
}

// Publish 메시지 브로드캐스트
func (p *V2XProcessorImpl) Publish(msg V2XMessage, radiusKm float64) error {
	// 여기에 코드를 작성하세요
	return nil
}

// GetStatistics 통계 반환
func (p *V2XProcessorImpl) GetStatistics() Statistics {
	// 여기에 코드를 작성하세요
	return Statistics{}
}
