package huhu1

/*
================================================================================
문제 6: 대용량 텔레메트리 데이터 파이프라인
================================================================================

난이도: Hard
주제: 데이터 파이프라인, 스트림 처리, 배치 처리
회사: 42dot (백엔드/데이터 플랫폼)

【 문제 설명 】
수만 대의 차량에서 초당 수백만 건의 텔레메트리 데이터가 들어옵니다.
이를 실시간으로 처리하고 저장하는 데이터 파이프라인을 구현하세요.

텔레메트리 데이터:
- 차량 ID, 타임스탬프
- GPS (위도, 경도, 속도, 방향)
- 센서 (배터리, 온도, 연료)
- 이벤트 (급정거, 급가속, 충돌 감지)

요구사항:
1. Ingest(data): 데이터 수집 (버퍼링)
2. Process(): 실시간 처리 (집계, 이상 탐지)
3. Flush(): 배치로 저장소에 쓰기
4. Subscribe(topic, handler): 실시간 이벤트 구독
5. GetMetrics(): 파이프라인 메트릭 반환

【 성능 요구사항 】
- 초당 100만 건 수집
- 처리 지연: 1초 이내
- 데이터 유실: 0% (at-least-once)
- 메모리 사용: 효율적인 버퍼 관리

【 인터페이스 】
type TelemetryData struct {
    VehicleID   string
    Timestamp   time.Time
    GPS         GPSData
    Sensors     SensorData
    Events      []EventData
}

type GPSData struct {
    Lat, Lon  float64
    Speed     float64  // km/h
    Heading   float64  // 0-360도
}

type SensorData struct {
    Battery     int     // 0-100%
    Temperature float64 // 섭씨
    Fuel        float64 // 리터
}

type EventData struct {
    Type      string  // "HARD_BRAKE", "HARD_ACCEL", "COLLISION"
    Severity  string  // "LOW", "MEDIUM", "HIGH"
    Value     float64
}

type PipelineMetrics struct {
    IngestedCount   int64
    ProcessedCount  int64
    FlushedCount    int64
    ErrorCount      int64
    BufferSize      int
    AvgLatency      time.Duration
}

type DataPipeline interface {
    Ingest(data TelemetryData) error
    Process() error
    Flush() error
    Subscribe(topic string, handler func(TelemetryData)) error
    GetMetrics() PipelineMetrics
    Start() error
    Stop() error
}

【 힌트 】
1. Ring Buffer로 메모리 효율적 버퍼링
2. Worker Pool로 병렬 처리
3. 배치 단위로 묶어서 I/O 최적화
4. Back-pressure로 과부하 방지

【 실무 연관성 】
- Kafka Consumer/Producer 패턴
- 실시간 데이터 처리 (Flink, Spark Streaming 개념)
- 차량 데이터 수집 플랫폼
================================================================================
*/

import "time"

// TelemetryData 차량 텔레메트리 데이터
type TelemetryData struct {
	VehicleID string
	Timestamp time.Time
	GPS       GPSData
	Sensors   SensorData
	Events    []EventData
}

// GPSData GPS 정보
type GPSData struct {
	Lat, Lon float64
	Speed    float64
	Heading  float64
}

// SensorData 센서 정보
type SensorData struct {
	Battery     int
	Temperature float64
	Fuel        float64
}

// EventData 이벤트 정보
type EventData struct {
	Type     string
	Severity string
	Value    float64
}

// PipelineMetrics 파이프라인 메트릭
type PipelineMetrics struct {
	IngestedCount  int64
	ProcessedCount int64
	FlushedCount   int64
	ErrorCount     int64
	BufferSize     int
	AvgLatency     time.Duration
}

// DataPipeline 데이터 파이프라인 인터페이스
type DataPipeline interface {
	Ingest(data TelemetryData) error
	Process() error
	Flush() error
	Subscribe(topic string, handler func(TelemetryData)) error
	GetMetrics() PipelineMetrics
	Start() error
	Stop() error
}

// DataPipelineImpl 구현체
type DataPipelineImpl struct {
	// 여기에 필드를 정의하세요
}

// NewDataPipeline 생성자
func NewDataPipeline(bufferSize int, flushInterval time.Duration) *DataPipelineImpl {
	// 여기에 코드를 작성하세요
	return nil
}

// Ingest 데이터 수집
func (p *DataPipelineImpl) Ingest(data TelemetryData) error {
	// 여기에 코드를 작성하세요
	return nil
}

// Process 실시간 처리
func (p *DataPipelineImpl) Process() error {
	// 여기에 코드를 작성하세요
	return nil
}

// Flush 배치 저장
func (p *DataPipelineImpl) Flush() error {
	// 여기에 코드를 작성하세요
	return nil
}

// Subscribe 이벤트 구독
func (p *DataPipelineImpl) Subscribe(topic string, handler func(TelemetryData)) error {
	// 여기에 코드를 작성하세요
	return nil
}

// GetMetrics 메트릭 반환
func (p *DataPipelineImpl) GetMetrics() PipelineMetrics {
	// 여기에 코드를 작성하세요
	return PipelineMetrics{}
}

// Start 파이프라인 시작
func (p *DataPipelineImpl) Start() error {
	// 여기에 코드를 작성하세요
	return nil
}

// Stop 파이프라인 정지
func (p *DataPipelineImpl) Stop() error {
	// 여기에 코드를 작성하세요
	return nil
}
