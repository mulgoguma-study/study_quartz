package huhu1

/*
================================================================================
문제 7: 시계열 데이터 저장소 최적화
================================================================================

난이도: Hard
주제: 시계열 데이터, 파티셔닝, 다운샘플링
회사: 42dot (백엔드/데이터 플랫폼)

【 문제 설명 】
차량 텔레메트리 데이터를 효율적으로 저장하고 조회하는 시계열 저장소를 구현하세요.
수십억 건의 데이터를 빠르게 조회하고, 오래된 데이터는 자동으로 다운샘플링합니다.

요구사항:
1. Write(point): 데이터 포인트 저장
2. Query(vehicleID, metric, from, to): 시간 범위 조회
3. Aggregate(vehicleID, metric, from, to, interval): 집계 조회 (평균, 최대, 최소)
4. Downsample(): 오래된 데이터 다운샘플링
5. Expire(): 만료된 데이터 삭제
6. GetStorageStats(): 저장소 통계

【 성능 요구사항 】
- 초당 10만 건 쓰기
- 1시간 범위 조회: 100ms 이내
- 저장 기간: Raw 7일, 1분 집계 30일, 1시간 집계 1년

【 인터페이스 】
type DataPoint struct {
    VehicleID string
    Metric    string    // "speed", "battery", "temperature"
    Value     float64
    Timestamp time.Time
    Tags      map[string]string
}

type QueryResult struct {
    Points    []DataPoint
    StartTime time.Time
    EndTime   time.Time
    Count     int
}

type AggregateResult struct {
    VehicleID string
    Metric    string
    Interval  time.Duration
    Buckets   []AggregateBucket
}

type AggregateBucket struct {
    Timestamp time.Time
    Min       float64
    Max       float64
    Avg       float64
    Sum       float64
    Count     int
}

type StorageStats struct {
    TotalPoints     int64
    RawPoints       int64
    DownsampledPoints int64
    StorageBytes    int64
    OldestTimestamp time.Time
    NewestTimestamp time.Time
}

type TimeSeriesStore interface {
    Write(point DataPoint) error
    WriteBatch(points []DataPoint) error
    Query(vehicleID, metric string, from, to time.Time) (*QueryResult, error)
    Aggregate(vehicleID, metric string, from, to time.Time, interval time.Duration) (*AggregateResult, error)
    Downsample() error
    Expire() error
    GetStorageStats() StorageStats
}

【 힌트 】
1. 시간 기반 파티셔닝 (일별/시간별)
2. 메트릭별 인덱싱
3. 메모리 버퍼 + 디스크 저장
4. 다운샘플링: Raw → 1분 평균 → 1시간 평균

【 실무 연관성 】
- TimescaleDB, InfluxDB 개념
- 모니터링 시스템 (Prometheus)
- 차량 데이터 분석 플랫폼
================================================================================
*/

import "time"

// DataPoint 시계열 데이터 포인트
type DataPoint struct {
	VehicleID string
	Metric    string
	Value     float64
	Timestamp time.Time
	Tags      map[string]string
}

// QueryResult 조회 결과
type QueryResult struct {
	Points    []DataPoint
	StartTime time.Time
	EndTime   time.Time
	Count     int
}

// AggregateResult 집계 결과
type AggregateResult struct {
	VehicleID string
	Metric    string
	Interval  time.Duration
	Buckets   []AggregateBucket
}

// AggregateBucket 집계 버킷
type AggregateBucket struct {
	Timestamp time.Time
	Min       float64
	Max       float64
	Avg       float64
	Sum       float64
	Count     int
}

// StorageStats 저장소 통계
type StorageStats struct {
	TotalPoints       int64
	RawPoints         int64
	DownsampledPoints int64
	StorageBytes      int64
	OldestTimestamp   time.Time
	NewestTimestamp   time.Time
}

// TimeSeriesStore 시계열 저장소 인터페이스
type TimeSeriesStore interface {
	Write(point DataPoint) error
	WriteBatch(points []DataPoint) error
	Query(vehicleID, metric string, from, to time.Time) (*QueryResult, error)
	Aggregate(vehicleID, metric string, from, to time.Time, interval time.Duration) (*AggregateResult, error)
	Downsample() error
	Expire() error
	GetStorageStats() StorageStats
}

// TimeSeriesStoreImpl 구현체
type TimeSeriesStoreImpl struct {
	// 여기에 필드를 정의하세요
}

// NewTimeSeriesStore 생성자
func NewTimeSeriesStore() *TimeSeriesStoreImpl {
	// 여기에 코드를 작성하세요
	return nil
}

// Write 데이터 포인트 저장
func (s *TimeSeriesStoreImpl) Write(point DataPoint) error {
	// 여기에 코드를 작성하세요
	return nil
}

// WriteBatch 배치 저장
func (s *TimeSeriesStoreImpl) WriteBatch(points []DataPoint) error {
	// 여기에 코드를 작성하세요
	return nil
}

// Query 시간 범위 조회
func (s *TimeSeriesStoreImpl) Query(vehicleID, metric string, from, to time.Time) (*QueryResult, error) {
	// 여기에 코드를 작성하세요
	return nil, nil
}

// Aggregate 집계 조회
func (s *TimeSeriesStoreImpl) Aggregate(vehicleID, metric string, from, to time.Time, interval time.Duration) (*AggregateResult, error) {
	// 여기에 코드를 작성하세요
	return nil, nil
}

// Downsample 다운샘플링
func (s *TimeSeriesStoreImpl) Downsample() error {
	// 여기에 코드를 작성하세요
	return nil
}

// Expire 만료 데이터 삭제
func (s *TimeSeriesStoreImpl) Expire() error {
	// 여기에 코드를 작성하세요
	return nil
}

// GetStorageStats 저장소 통계
func (s *TimeSeriesStoreImpl) GetStorageStats() StorageStats {
	// 여기에 코드를 작성하세요
	return StorageStats{}
}
