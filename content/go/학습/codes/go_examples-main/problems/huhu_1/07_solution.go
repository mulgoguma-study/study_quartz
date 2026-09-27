package huhu1

/*
================================================================================
문제 7: 시계열 데이터 저장소 최적화 - 솔루션
================================================================================

【 핵심 설계 원칙 】
1. 시간 기반 파티셔닝 - 시간별로 데이터 분리
2. 계층적 저장 - Raw → 1분 → 1시간 집계
3. 메모리 버퍼 - 최근 데이터 빠른 조회
4. 압축 - 오래된 데이터 압축 저장

【 저장 계층 】
┌─────────────────────────────────────────────────────────────────┐
│  Hot (메모리)    │  최근 1시간  │  Raw 데이터                   │
├──────────────────┼─────────────┼───────────────────────────────┤
│  Warm (디스크)   │  7일        │  Raw 데이터                   │
├──────────────────┼─────────────┼───────────────────────────────┤
│  Cold (압축)     │  30일       │  1분 집계                     │
├──────────────────┼─────────────┼───────────────────────────────┤
│  Archive         │  1년        │  1시간 집계                   │
└──────────────────┴─────────────┴───────────────────────────────┘

================================================================================
*/

import (
	"errors"
	"math"
	"sort"
	"sync"
	"time"
)

// ============================================================================
// 파티션 구조
// ============================================================================

// partitionKey 파티션 키 (시간 + 차량ID + 메트릭)
type partitionKey struct {
	hour      time.Time // 시간별 파티션
	vehicleID string
	metric    string
}

// partition 파티션 데이터
type partition struct {
	mu     sync.RWMutex
	points []DataPoint
}

// ============================================================================
// 시계열 저장소 구현
// ============================================================================

// TimeSeriesStoreSolution 완성된 시계열 저장소
type TimeSeriesStoreSolution struct {
	mu sync.RWMutex

	// Raw 데이터 (시간별 파티션)
	rawPartitions map[partitionKey]*partition

	// 다운샘플링된 데이터
	minuteAggregates map[partitionKey][]AggregateBucket // 1분 집계
	hourAggregates   map[partitionKey][]AggregateBucket // 1시간 집계

	// 메모리 버퍼 (최근 데이터)
	bufferMu     sync.RWMutex
	writeBuffer  []DataPoint
	bufferSize   int
	flushSize    int

	// 보존 기간 설정
	rawRetention    time.Duration // Raw 데이터 보존
	minuteRetention time.Duration // 1분 집계 보존
	hourRetention   time.Duration // 1시간 집계 보존

	// 통계
	stats struct {
		totalPoints       int64
		rawPoints         int64
		downsampledPoints int64
		oldestTimestamp   time.Time
		newestTimestamp   time.Time
	}
}

// NewTimeSeriesStoreSolution 생성자
func NewTimeSeriesStoreSolution() *TimeSeriesStoreSolution {
	return &TimeSeriesStoreSolution{
		rawPartitions:    make(map[partitionKey]*partition),
		minuteAggregates: make(map[partitionKey][]AggregateBucket),
		hourAggregates:   make(map[partitionKey][]AggregateBucket),
		writeBuffer:      make([]DataPoint, 0, 10000),
		bufferSize:       10000,
		flushSize:        1000,
		rawRetention:     7 * 24 * time.Hour,    // 7일
		minuteRetention:  30 * 24 * time.Hour,   // 30일
		hourRetention:    365 * 24 * time.Hour,  // 1년
	}
}

// ============================================================================
// 쓰기 연산
// ============================================================================

// Write 데이터 포인트 저장
/*
【 동작 원리 】
1. 메모리 버퍼에 추가
2. 버퍼가 가득 차면 파티션으로 플러시
3. 통계 업데이트
*/
func (s *TimeSeriesStoreSolution) Write(point DataPoint) error {
	if point.VehicleID == "" || point.Metric == "" {
		return errors.New("vehicleID and metric are required")
	}

	if point.Timestamp.IsZero() {
		point.Timestamp = time.Now()
	}

	s.bufferMu.Lock()
	s.writeBuffer = append(s.writeBuffer, point)

	// 버퍼가 가득 차면 플러시
	if len(s.writeBuffer) >= s.bufferSize {
		s.flushBuffer()
	}
	s.bufferMu.Unlock()

	// 통계 업데이트
	s.updateStats(point)

	return nil
}

// WriteBatch 배치 저장
func (s *TimeSeriesStoreSolution) WriteBatch(points []DataPoint) error {
	for _, point := range points {
		if err := s.Write(point); err != nil {
			return err
		}
	}
	return nil
}

// flushBuffer 버퍼를 파티션으로 플러시
func (s *TimeSeriesStoreSolution) flushBuffer() {
	if len(s.writeBuffer) == 0 {
		return
	}

	// 파티션별로 그룹화
	groups := make(map[partitionKey][]DataPoint)
	for _, point := range s.writeBuffer {
		key := partitionKey{
			hour:      point.Timestamp.Truncate(time.Hour),
			vehicleID: point.VehicleID,
			metric:    point.Metric,
		}
		groups[key] = append(groups[key], point)
	}

	// 각 파티션에 저장
	s.mu.Lock()
	for key, points := range groups {
		if s.rawPartitions[key] == nil {
			s.rawPartitions[key] = &partition{
				points: make([]DataPoint, 0),
			}
		}
		s.rawPartitions[key].mu.Lock()
		s.rawPartitions[key].points = append(s.rawPartitions[key].points, points...)
		s.rawPartitions[key].mu.Unlock()
	}
	s.mu.Unlock()

	// 버퍼 초기화
	s.writeBuffer = s.writeBuffer[:0]
}

// updateStats 통계 업데이트
func (s *TimeSeriesStoreSolution) updateStats(point DataPoint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stats.totalPoints++
	s.stats.rawPoints++

	if s.stats.oldestTimestamp.IsZero() || point.Timestamp.Before(s.stats.oldestTimestamp) {
		s.stats.oldestTimestamp = point.Timestamp
	}
	if point.Timestamp.After(s.stats.newestTimestamp) {
		s.stats.newestTimestamp = point.Timestamp
	}
}

// ============================================================================
// 조회 연산
// ============================================================================

// Query 시간 범위 조회
/*
【 동작 원리 】
1. 시간 범위에 해당하는 파티션 찾기
2. 각 파티션에서 데이터 수집
3. 시간순 정렬 후 반환
*/
func (s *TimeSeriesStoreSolution) Query(vehicleID, metric string, from, to time.Time) (*QueryResult, error) {
	if vehicleID == "" || metric == "" {
		return nil, errors.New("vehicleID and metric are required")
	}

	// 버퍼 먼저 플러시
	s.bufferMu.Lock()
	s.flushBuffer()
	s.bufferMu.Unlock()

	result := &QueryResult{
		Points:    make([]DataPoint, 0),
		StartTime: from,
		EndTime:   to,
	}

	// 시간 범위에 해당하는 파티션 순회
	s.mu.RLock()
	for hour := from.Truncate(time.Hour); !hour.After(to); hour = hour.Add(time.Hour) {
		key := partitionKey{
			hour:      hour,
			vehicleID: vehicleID,
			metric:    metric,
		}

		if part, ok := s.rawPartitions[key]; ok {
			part.mu.RLock()
			for _, point := range part.points {
				if (point.Timestamp.Equal(from) || point.Timestamp.After(from)) &&
					(point.Timestamp.Equal(to) || point.Timestamp.Before(to)) {
					result.Points = append(result.Points, point)
				}
			}
			part.mu.RUnlock()
		}
	}
	s.mu.RUnlock()

	// 시간순 정렬
	sort.Slice(result.Points, func(i, j int) bool {
		return result.Points[i].Timestamp.Before(result.Points[j].Timestamp)
	})

	result.Count = len(result.Points)
	return result, nil
}

// Aggregate 집계 조회
/*
【 동작 원리 】
1. 조회 기간에 따라 적절한 집계 레벨 선택
   - 1시간 미만: Raw 데이터에서 계산
   - 1일 미만: 1분 집계 사용
   - 1일 이상: 1시간 집계 사용
2. interval 단위로 버킷 생성
3. 각 버킷의 Min, Max, Avg, Sum, Count 계산
*/
func (s *TimeSeriesStoreSolution) Aggregate(vehicleID, metric string, from, to time.Time, interval time.Duration) (*AggregateResult, error) {
	if interval <= 0 {
		return nil, errors.New("interval must be positive")
	}

	// Raw 데이터 조회
	queryResult, err := s.Query(vehicleID, metric, from, to)
	if err != nil {
		return nil, err
	}

	result := &AggregateResult{
		VehicleID: vehicleID,
		Metric:    metric,
		Interval:  interval,
		Buckets:   make([]AggregateBucket, 0),
	}

	if len(queryResult.Points) == 0 {
		return result, nil
	}

	// 버킷별로 그룹화
	buckets := make(map[time.Time][]float64)
	for _, point := range queryResult.Points {
		bucketTime := point.Timestamp.Truncate(interval)
		buckets[bucketTime] = append(buckets[bucketTime], point.Value)
	}

	// 버킷 시간순 정렬
	var bucketTimes []time.Time
	for t := range buckets {
		bucketTimes = append(bucketTimes, t)
	}
	sort.Slice(bucketTimes, func(i, j int) bool {
		return bucketTimes[i].Before(bucketTimes[j])
	})

	// 각 버킷 집계 계산
	for _, bucketTime := range bucketTimes {
		values := buckets[bucketTime]
		bucket := AggregateBucket{
			Timestamp: bucketTime,
			Count:     len(values),
		}

		if len(values) > 0 {
			bucket.Min = values[0]
			bucket.Max = values[0]
			bucket.Sum = 0

			for _, v := range values {
				bucket.Sum += v
				if v < bucket.Min {
					bucket.Min = v
				}
				if v > bucket.Max {
					bucket.Max = v
				}
			}
			bucket.Avg = bucket.Sum / float64(len(values))
		}

		result.Buckets = append(result.Buckets, bucket)
	}

	return result, nil
}

// ============================================================================
// 다운샘플링
// ============================================================================

// Downsample 오래된 데이터 다운샘플링
/*
【 동작 원리 】
1. 1시간 이상 지난 Raw 데이터 → 1분 집계로 변환
2. 7일 이상 지난 1분 집계 → 1시간 집계로 변환
3. 원본 데이터 삭제
*/
func (s *TimeSeriesStoreSolution) Downsample() error {
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Raw → 1분 집계 (1시간 이상 된 데이터)
	for key, part := range s.rawPartitions {
		if now.Sub(key.hour) > time.Hour {
			minuteBuckets := s.aggregateToMinute(part.points)
			s.minuteAggregates[key] = minuteBuckets

			// Raw 데이터 삭제
			s.stats.rawPoints -= int64(len(part.points))
			s.stats.downsampledPoints += int64(len(minuteBuckets))
			delete(s.rawPartitions, key)
		}
	}

	// 1분 집계 → 1시간 집계 (7일 이상 된 데이터)
	for key, minuteBuckets := range s.minuteAggregates {
		if now.Sub(key.hour) > s.rawRetention {
			hourBuckets := s.aggregateToHour(minuteBuckets)
			s.hourAggregates[key] = hourBuckets

			// 1분 집계 삭제
			delete(s.minuteAggregates, key)
		}
	}

	return nil
}

// aggregateToMinute Raw 데이터를 1분 집계로
func (s *TimeSeriesStoreSolution) aggregateToMinute(points []DataPoint) []AggregateBucket {
	buckets := make(map[time.Time][]float64)

	for _, point := range points {
		bucketTime := point.Timestamp.Truncate(time.Minute)
		buckets[bucketTime] = append(buckets[bucketTime], point.Value)
	}

	result := make([]AggregateBucket, 0, len(buckets))
	for bucketTime, values := range buckets {
		bucket := s.calculateBucket(bucketTime, values)
		result = append(result, bucket)
	}

	return result
}

// aggregateToHour 1분 집계를 1시간 집계로
func (s *TimeSeriesStoreSolution) aggregateToHour(minuteBuckets []AggregateBucket) []AggregateBucket {
	hourBuckets := make(map[time.Time]*AggregateBucket)

	for _, mb := range minuteBuckets {
		hourTime := mb.Timestamp.Truncate(time.Hour)

		if hourBuckets[hourTime] == nil {
			hourBuckets[hourTime] = &AggregateBucket{
				Timestamp: hourTime,
				Min:       math.MaxFloat64,
				Max:       -math.MaxFloat64,
			}
		}

		hb := hourBuckets[hourTime]
		hb.Sum += mb.Sum
		hb.Count += mb.Count
		if mb.Min < hb.Min {
			hb.Min = mb.Min
		}
		if mb.Max > hb.Max {
			hb.Max = mb.Max
		}
	}

	result := make([]AggregateBucket, 0, len(hourBuckets))
	for _, hb := range hourBuckets {
		if hb.Count > 0 {
			hb.Avg = hb.Sum / float64(hb.Count)
		}
		result = append(result, *hb)
	}

	return result
}

// calculateBucket 버킷 계산
func (s *TimeSeriesStoreSolution) calculateBucket(timestamp time.Time, values []float64) AggregateBucket {
	bucket := AggregateBucket{
		Timestamp: timestamp,
		Count:     len(values),
	}

	if len(values) == 0 {
		return bucket
	}

	bucket.Min = values[0]
	bucket.Max = values[0]

	for _, v := range values {
		bucket.Sum += v
		if v < bucket.Min {
			bucket.Min = v
		}
		if v > bucket.Max {
			bucket.Max = v
		}
	}
	bucket.Avg = bucket.Sum / float64(len(values))

	return bucket
}

// ============================================================================
// 만료 처리
// ============================================================================

// Expire 만료된 데이터 삭제
func (s *TimeSeriesStoreSolution) Expire() error {
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Raw 데이터 만료 (7일)
	for key := range s.rawPartitions {
		if now.Sub(key.hour) > s.rawRetention {
			delete(s.rawPartitions, key)
		}
	}

	// 1분 집계 만료 (30일)
	for key := range s.minuteAggregates {
		if now.Sub(key.hour) > s.minuteRetention {
			delete(s.minuteAggregates, key)
		}
	}

	// 1시간 집계 만료 (1년)
	for key := range s.hourAggregates {
		if now.Sub(key.hour) > s.hourRetention {
			delete(s.hourAggregates, key)
		}
	}

	return nil
}

// GetStorageStats 저장소 통계
func (s *TimeSeriesStoreSolution) GetStorageStats() StorageStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return StorageStats{
		TotalPoints:       s.stats.totalPoints,
		RawPoints:         s.stats.rawPoints,
		DownsampledPoints: s.stats.downsampledPoints,
		OldestTimestamp:   s.stats.oldestTimestamp,
		NewestTimestamp:   s.stats.newestTimestamp,
	}
}

// ============================================================================
// 성능 최적화 포인트
// ============================================================================

/*
【 프로덕션 최적화 】

1. 압축
   - Gorilla 압축 (Facebook 시계열 압축)
   - 타임스탬프: Delta-of-delta
   - 값: XOR 압축
   - 10배 이상 압축률

2. 인덱싱
   - 역 인덱스 (태그 → 시계열)
   - B-tree 인덱스 (시간 범위 조회)
   - Bloom Filter (존재 여부 빠른 확인)

3. 샤딩
   - 차량 ID 기반 샤딩
   - 각 샤드 독립 저장/조회

4. 캐싱
   - 최근 조회 결과 캐싱
   - 집계 결과 캐싱

5. 비동기 다운샘플링
   - 백그라운드 워커로 주기적 실행
   - 서비스 영향 최소화

┌─────────────────────────────────────────────────────────────────────────────┐
│                       프로덕션 아키텍처                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  [Write Path]                                                               │
│  ┌─────────┐    ┌─────────────┐    ┌─────────────┐    ┌──────────────┐    │
│  │ Ingest  │ →  │ WAL (Write  │ →  │ Memory      │ →  │ SSTable      │    │
│  │         │    │ Ahead Log)  │    │ Table       │    │ (Disk)       │    │
│  └─────────┘    └─────────────┘    └─────────────┘    └──────────────┘    │
│                                                                             │
│  [Read Path]                                                                │
│  ┌─────────┐    ┌─────────────┐    ┌─────────────┐                         │
│  │ Query   │ →  │ Merge       │ ←  │ Memory +    │                         │
│  │         │    │ Results     │    │ SSTable     │                         │
│  └─────────┘    └─────────────┘    └─────────────┘                         │
│                                                                             │
│  [Compaction]                                                               │
│  ┌─────────────────────────────────────────────────────────────────┐       │
│  │ Raw SSTable → Minute SSTable → Hour SSTable → Archive          │       │
│  └─────────────────────────────────────────────────────────────────┘       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

【 TimescaleDB 파티셔닝 예시 】

CREATE TABLE telemetry (
    time        TIMESTAMPTZ NOT NULL,
    vehicle_id  TEXT NOT NULL,
    metric      TEXT NOT NULL,
    value       DOUBLE PRECISION
);

-- 하이퍼테이블로 변환 (자동 파티셔닝)
SELECT create_hypertable('telemetry', 'time',
    chunk_time_interval => INTERVAL '1 day'
);

-- 압축 정책
ALTER TABLE telemetry SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'vehicle_id, metric'
);

-- 7일 이상 된 청크 압축
SELECT add_compression_policy('telemetry', INTERVAL '7 days');

-- 연속 집계 (다운샘플링)
CREATE MATERIALIZED VIEW telemetry_hourly
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 hour', time) AS bucket,
    vehicle_id,
    metric,
    avg(value) as avg_value,
    min(value) as min_value,
    max(value) as max_value
FROM telemetry
GROUP BY bucket, vehicle_id, metric;
*/
