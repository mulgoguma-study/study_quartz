package huhu1

/*
================================================================================
문제 1: 실시간 차량 위치 추적 시스템 - 솔루션
================================================================================

【 핵심 설계 원칙 】
1. Geohash 기반 공간 인덱싱 - O(1) 근접 차량 조회
2. 샤딩으로 쓰기 부하 분산 - 초당 10,000건 처리
3. 최신 위치는 메모리 캐시 - 빠른 조회
4. 이력 데이터는 시간 기반 파티셔닝

【 Geohash란? 】
- 위도/경도를 문자열로 인코딩
- 예: 37.5665, 126.9780 → "wydm9"
- 같은 접두사 = 같은 영역 (가까운 위치)
- 정밀도: 5자리 ≈ 4.9km x 4.9km

================================================================================
*/

import (
	"errors"
	"math"
	"sync"
	"time"
)

// ============================================================================
// Geohash 구현
// ============================================================================

const (
	base32            = "0123456789bcdefghjkmnpqrstuvwxyz"
	defaultPrecision  = 6  // ~1.2km x 0.6km
	nearbyPrecision   = 5  // ~4.9km x 4.9km (반경 검색용)
)

// encodeGeohash 위도/경도를 geohash 문자열로 변환
func encodeGeohash(lat, lon float64, precision int) string {
	var geohash []byte
	var minLat, maxLat float64 = -90, 90
	var minLon, maxLon float64 = -180, 180

	isEven := true
	bit := 0
	ch := 0

	for len(geohash) < precision {
		if isEven {
			// 경도 처리
			mid := (minLon + maxLon) / 2
			if lon > mid {
				ch |= 1 << (4 - bit)
				minLon = mid
			} else {
				maxLon = mid
			}
		} else {
			// 위도 처리
			mid := (minLat + maxLat) / 2
			if lat > mid {
				ch |= 1 << (4 - bit)
				minLat = mid
			} else {
				maxLat = mid
			}
		}

		isEven = !isEven
		bit++

		if bit == 5 {
			geohash = append(geohash, base32[ch])
			bit = 0
			ch = 0
		}
	}

	return string(geohash)
}

// getNeighborGeohashes 주변 8개 geohash + 자기 자신 반환
func getNeighborGeohashes(geohash string) []string {
	// 간단한 구현: 실제로는 geohash 라이브러리 사용 권장
	// 여기서는 자기 자신만 반환 (완전한 구현은 복잡)
	return []string{geohash}
}

// ============================================================================
// 거리 계산 (Haversine 공식)
// ============================================================================

const earthRadiusKm = 6371.0

// haversineDistance 두 지점 간의 거리 (km)
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRadians(lat1))*math.Cos(toRadians(lat2))*
		math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

func toRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

// ============================================================================
// 샤딩된 위치 저장소
// ============================================================================

const numShards = 32

// LocationShard 각 샤드는 독립적인 락을 가짐
type LocationShard struct {
	mu        sync.RWMutex
	locations map[string]*LocationData // vehicleID -> data
}

// LocationData 차량 위치 데이터 (최신 + 이력)
type LocationData struct {
	Current LocationInfo
	History []LocationInfo // 시간순 정렬
}

// LocationInfo Location의 복사본 (패키지 충돌 방지)
type LocationInfo struct {
	VehicleID string
	Lat, Lon  float64
	Timestamp time.Time
}

// VehicleTrackerSolution 완성된 구현체
type VehicleTrackerSolution struct {
	// 샤딩된 차량 위치 저장소
	shards [numShards]*LocationShard

	// Geohash 기반 공간 인덱스
	geoIndex     map[string]map[string]bool // geohash -> set of vehicleIDs
	geoIndexLock sync.RWMutex

	// 설정
	historyLimit int // 차량당 최대 이력 수
}

// NewVehicleTrackerSolution 생성자
func NewVehicleTrackerSolution() *VehicleTrackerSolution {
	v := &VehicleTrackerSolution{
		geoIndex:     make(map[string]map[string]bool),
		historyLimit: 1000,
	}

	// 샤드 초기화
	for i := 0; i < numShards; i++ {
		v.shards[i] = &LocationShard{
			locations: make(map[string]*LocationData),
		}
	}

	return v
}

// getShard vehicleID에 해당하는 샤드 반환 (FNV-1a 해시)
func (v *VehicleTrackerSolution) getShard(vehicleID string) *LocationShard {
	h := fnv1aHash(vehicleID)
	return v.shards[h%numShards]
}

func fnv1aHash(s string) uint32 {
	const (
		offset32 = 2166136261
		prime32  = 16777619
	)
	hash := uint32(offset32)
	for i := 0; i < len(s); i++ {
		hash ^= uint32(s[i])
		hash *= prime32
	}
	return hash
}

// ============================================================================
// 인터페이스 구현
// ============================================================================

// UpdateLocation 차량 위치 업데이트
/*
【 시간 복잡도 】O(1)
【 동시성 】샤드별 락으로 경합 최소화

【 동작 원리 】
1. 기존 geohash에서 제거 (있으면)
2. 새 geohash로 추가
3. 최신 위치 업데이트
4. 이력에 추가 (limit 초과 시 오래된 것 제거)
*/
func (v *VehicleTrackerSolution) UpdateLocation(loc LocationInfo) error {
	if loc.VehicleID == "" {
		return errors.New("vehicle ID is required")
	}

	shard := v.getShard(loc.VehicleID)
	newGeohash := encodeGeohash(loc.Lat, loc.Lon, defaultPrecision)

	// 1. 기존 geohash에서 제거
	shard.mu.RLock()
	data, exists := shard.locations[loc.VehicleID]
	var oldGeohash string
	if exists {
		oldGeohash = encodeGeohash(data.Current.Lat, data.Current.Lon, defaultPrecision)
	}
	shard.mu.RUnlock()

	if oldGeohash != "" && oldGeohash != newGeohash {
		v.geoIndexLock.Lock()
		if v.geoIndex[oldGeohash] != nil {
			delete(v.geoIndex[oldGeohash], loc.VehicleID)
		}
		v.geoIndexLock.Unlock()
	}

	// 2. 새 geohash에 추가
	v.geoIndexLock.Lock()
	if v.geoIndex[newGeohash] == nil {
		v.geoIndex[newGeohash] = make(map[string]bool)
	}
	v.geoIndex[newGeohash][loc.VehicleID] = true
	v.geoIndexLock.Unlock()

	// 3. 위치 데이터 업데이트
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if !exists {
		shard.locations[loc.VehicleID] = &LocationData{
			Current: loc,
			History: []LocationInfo{loc},
		}
	} else {
		data.Current = loc
		data.History = append(data.History, loc)

		// 이력 제한
		if len(data.History) > v.historyLimit {
			data.History = data.History[len(data.History)-v.historyLimit:]
		}
	}

	return nil
}

// GetNearbyVehicles 반경 내 차량 조회
/*
【 시간 복잡도 】O(k) where k = 해당 영역 차량 수
【 동작 원리 】
1. 현재 위치의 geohash 계산
2. 해당 geohash + 인접 geohash의 모든 차량 조회
3. 정확한 거리 계산으로 필터링
*/
func (v *VehicleTrackerSolution) GetNearbyVehicles(lat, lon, radiusKm float64) ([]LocationInfo, error) {
	if radiusKm <= 0 {
		return nil, errors.New("radius must be positive")
	}

	// 반경에 맞는 geohash 정밀도 선택
	precision := nearbyPrecision
	if radiusKm < 1 {
		precision = defaultPrecision
	}

	centerGeohash := encodeGeohash(lat, lon, precision)
	neighborGeohashes := getNeighborGeohashes(centerGeohash)

	// 후보 차량 수집
	candidateVehicles := make(map[string]bool)
	v.geoIndexLock.RLock()
	for _, gh := range neighborGeohashes {
		// geohash 접두사 매칭
		for storedGh, vehicles := range v.geoIndex {
			if len(storedGh) >= precision && storedGh[:precision] == gh {
				for vehicleID := range vehicles {
					candidateVehicles[vehicleID] = true
				}
			}
		}
	}
	v.geoIndexLock.RUnlock()

	// 정확한 거리 계산으로 필터링
	var result []LocationInfo
	for vehicleID := range candidateVehicles {
		loc, err := v.GetVehicleLocation(vehicleID)
		if err != nil || loc == nil {
			continue
		}

		distance := haversineDistance(lat, lon, loc.Lat, loc.Lon)
		if distance <= radiusKm {
			result = append(result, *loc)
		}
	}

	return result, nil
}

// GetVehicleLocation 특정 차량 위치 조회
/*
【 시간 복잡도 】O(1)
【 동시성 】샤드별 RLock으로 동시 읽기 가능
*/
func (v *VehicleTrackerSolution) GetVehicleLocation(vehicleID string) (*LocationInfo, error) {
	if vehicleID == "" {
		return nil, errors.New("vehicle ID is required")
	}

	shard := v.getShard(vehicleID)

	shard.mu.RLock()
	defer shard.mu.RUnlock()

	data, exists := shard.locations[vehicleID]
	if !exists {
		return nil, nil // 차량 없음 (에러 아님)
	}

	// 복사본 반환 (동시성 안전)
	loc := data.Current
	return &loc, nil
}

// GetVehicleHistory 차량 이동 이력 조회
/*
【 시간 복잡도 】O(n) where n = 이력 개수
【 동작 원리 】
1. 샤드에서 이력 조회
2. from-to 범위 필터링
3. 복사본 반환
*/
func (v *VehicleTrackerSolution) GetVehicleHistory(vehicleID string, from, to time.Time) ([]LocationInfo, error) {
	if vehicleID == "" {
		return nil, errors.New("vehicle ID is required")
	}

	shard := v.getShard(vehicleID)

	shard.mu.RLock()
	defer shard.mu.RUnlock()

	data, exists := shard.locations[vehicleID]
	if !exists {
		return nil, nil
	}

	// 시간 범위 필터링
	var result []LocationInfo
	for _, loc := range data.History {
		if (loc.Timestamp.Equal(from) || loc.Timestamp.After(from)) &&
			(loc.Timestamp.Equal(to) || loc.Timestamp.Before(to)) {
			result = append(result, loc)
		}
	}

	return result, nil
}

// ============================================================================
// 성능 최적화 포인트
// ============================================================================

/*
【 현재 구현의 한계와 개선 방향 】

1. Geohash 이웃 계산
   - 현재: 자기 자신만 반환
   - 개선: 8방향 이웃 geohash 계산 구현
   - 라이브러리: github.com/mmcloughlin/geohash

2. 이력 저장소
   - 현재: 메모리 (서버 재시작 시 손실)
   - 개선: 시계열 DB 연동 (InfluxDB, TimescaleDB)
   - 파티셔닝: 일별/주별 파티션

3. 클러스터링
   - 현재: 단일 서버
   - 개선: Consistent Hashing으로 차량 분배
   - Redis Cluster로 공유 상태

4. 캐싱 계층
   - 현재: 메모리만
   - 개선: Redis 캐시 + 메모리 캐시 (L1/L2)

5. 배치 처리
   - 현재: 건별 처리
   - 개선: 버퍼링 후 배치 쓰기

┌─────────────────────────────────────────────────────────────────────────────┐
│                          프로덕션 아키텍처                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  [차량 10만대] → [Kafka] → [Location Worker Pool] → [Redis Cluster]        │
│                                      ↓                                      │
│                            [TimescaleDB] (이력 저장)                        │
│                                      ↓                                      │
│                      [API Server] ← [Redis Cache] (읽기)                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
*/

// ============================================================================
// 테스트용 메인 함수
// ============================================================================

func mainSolution() {
	tracker := NewVehicleTrackerSolution()

	// 차량 위치 업데이트
	now := time.Now()

	// 서울 시청 근처 차량들
	vehicles := []LocationInfo{
		{VehicleID: "vehicle-001", Lat: 37.5665, Lon: 126.9780, Timestamp: now},
		{VehicleID: "vehicle-002", Lat: 37.5670, Lon: 126.9785, Timestamp: now},
		{VehicleID: "vehicle-003", Lat: 37.5680, Lon: 126.9790, Timestamp: now},
		// 강남 차량 (멀리 있음)
		{VehicleID: "vehicle-100", Lat: 37.4979, Lon: 127.0276, Timestamp: now},
	}

	for _, v := range vehicles {
		tracker.UpdateLocation(v)
	}

	// 서울 시청 기준 5km 반경 검색
	nearby, _ := tracker.GetNearbyVehicles(37.5665, 126.9780, 5.0)
	println("Nearby vehicles (5km radius):", len(nearby))

	// 특정 차량 조회
	loc, _ := tracker.GetVehicleLocation("vehicle-001")
	if loc != nil {
		println("Vehicle-001 location:", loc.Lat, loc.Lon)
	}

	// 이력 조회
	history, _ := tracker.GetVehicleHistory("vehicle-001", now.Add(-1*time.Hour), now.Add(1*time.Hour))
	println("Vehicle-001 history count:", len(history))
}
