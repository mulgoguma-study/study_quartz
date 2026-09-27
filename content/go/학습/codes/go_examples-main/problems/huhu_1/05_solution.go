package huhu1

/*
================================================================================
문제 5: 실시간 배차 알고리즘 - 솔루션
================================================================================

【 핵심 설계 원칙 】
1. Greedy 매칭 + 전역 최적화 - 빠른 응답 + 품질
2. 스코어링 함수 - 다중 요소 균형
3. 배치 처리 - 1초 단위로 모아서 최적화
4. 공간 인덱싱 - 후보 차량 빠르게 필터링

【 스코어링 공식 】
Score = α × (1/distance) + β × waitingPenalty + γ × batteryScore

- distance: 차량-승객 거리 (km)
- waitingPenalty: 대기 시간에 따른 가중치 (오래 기다릴수록 높음)
- batteryScore: 배터리 잔량 (충분할수록 높음)

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
// 스코어링 가중치
// ============================================================================

const (
	weightDistance = 0.5 // 거리 가중치
	weightWaiting  = 0.3 // 대기 시간 가중치
	weightBattery  = 0.2 // 배터리 가중치

	avgSpeedKmH        = 30.0  // 평균 속도 (km/h)
	minBatteryPercent  = 20    // 최소 필요 배터리
	maxSearchRadiusKm  = 10.0  // 최대 검색 반경
	waitingBonusPerMin = 0.1   // 대기 1분당 보너스
)

// ============================================================================
// Dispatcher 구현
// ============================================================================

// DispatcherSolution 완성된 배차 시스템
type DispatcherSolution struct {
	mu sync.RWMutex

	// 요청 관리
	pendingRequests map[string]*RideRequest // requestID -> request
	assignedRides   map[string]string       // rideID -> vehicleID

	// 차량 관리
	vehicles map[string]*VehicleStatus // vehicleID -> status

	// 공간 인덱스 (Geohash)
	vehicleIndex map[string]map[string]bool // geohash -> vehicleIDs
}

// NewDispatcherSolution 생성자
func NewDispatcherSolution() *DispatcherSolution {
	return &DispatcherSolution{
		pendingRequests: make(map[string]*RideRequest),
		assignedRides:   make(map[string]string),
		vehicles:        make(map[string]*VehicleStatus),
		vehicleIndex:    make(map[string]map[string]bool),
	}
}

// ============================================================================
// 기본 연산
// ============================================================================

// RequestRide 탑승 요청 등록
func (d *DispatcherSolution) RequestRide(request RideRequest) error {
	if request.ID == "" {
		return errors.New("request ID is required")
	}
	if request.Passengers <= 0 {
		return errors.New("passengers must be positive")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.pendingRequests[request.ID]; exists {
		return errors.New("request already exists")
	}

	// 요청 시간이 없으면 현재 시간으로 설정
	if request.RequestTime.IsZero() {
		request.RequestTime = time.Now()
	}

	d.pendingRequests[request.ID] = &request
	return nil
}

// UpdateVehicleStatus 차량 상태 업데이트
func (d *DispatcherSolution) UpdateVehicleStatus(status VehicleStatus) error {
	if status.VehicleID == "" {
		return errors.New("vehicle ID is required")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	// 이전 위치에서 인덱스 제거
	if oldStatus, exists := d.vehicles[status.VehicleID]; exists {
		oldGeohash := encodeGeohashDispatch(oldStatus.Lat, oldStatus.Lon, 5)
		if d.vehicleIndex[oldGeohash] != nil {
			delete(d.vehicleIndex[oldGeohash], status.VehicleID)
		}
	}

	// 새 위치로 인덱스 추가
	newGeohash := encodeGeohashDispatch(status.Lat, status.Lon, 5)
	if d.vehicleIndex[newGeohash] == nil {
		d.vehicleIndex[newGeohash] = make(map[string]bool)
	}
	d.vehicleIndex[newGeohash][status.VehicleID] = true

	d.vehicles[status.VehicleID] = &status
	return nil
}

// CancelRequest 요청 취소
func (d *DispatcherSolution) CancelRequest(requestID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.pendingRequests[requestID]; !exists {
		return errors.New("request not found")
	}

	delete(d.pendingRequests, requestID)
	return nil
}

// GetWaitingRequests 대기 중인 요청 목록
func (d *DispatcherSolution) GetWaitingRequests() []RideRequest {
	d.mu.RLock()
	defer d.mu.RUnlock()

	result := make([]RideRequest, 0, len(d.pendingRequests))
	for _, req := range d.pendingRequests {
		result = append(result, *req)
	}

	// 요청 시간순 정렬
	sort.Slice(result, func(i, j int) bool {
		return result[i].RequestTime.Before(result[j].RequestTime)
	})

	return result
}

// ============================================================================
// 최적 배차 알고리즘
// ============================================================================

/*
【 알고리즘 개요 】
1. 모든 대기 요청에 대해
2. 가용 차량 중 후보 필터링 (거리, 배터리, 용량)
3. 각 후보에 대해 스코어 계산
4. Greedy로 최고 스코어 매칭

【 시간 복잡도 】O(R × V) where R=요청 수, V=차량 수
【 개선 】헝가리안 알고리즘으로 전역 최적화 가능 (O(R³))
*/

// GetOptimalAssignment 최적 배차 계산
func (d *DispatcherSolution) GetOptimalAssignment() ([]Assignment, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if len(d.pendingRequests) == 0 {
		return []Assignment{}, nil
	}

	// 요청을 대기 시간순으로 정렬 (오래 기다린 요청 우선)
	requests := make([]*RideRequest, 0, len(d.pendingRequests))
	for _, req := range d.pendingRequests {
		requests = append(requests, req)
	}
	sort.Slice(requests, func(i, j int) bool {
		return requests[i].RequestTime.Before(requests[j].RequestTime)
	})

	// 사용된 차량 추적
	usedVehicles := make(map[string]bool)
	assignments := make([]Assignment, 0)

	for _, request := range requests {
		// 후보 차량 찾기
		candidates := d.findCandidateVehicles(request, usedVehicles)

		if len(candidates) == 0 {
			continue // 가용 차량 없음
		}

		// 최고 점수 후보 선택
		bestAssignment := d.findBestAssignment(request, candidates)
		if bestAssignment != nil {
			assignments = append(assignments, *bestAssignment)
			usedVehicles[bestAssignment.VehicleID] = true
		}
	}

	return assignments, nil
}

// findCandidateVehicles 후보 차량 찾기
func (d *DispatcherSolution) findCandidateVehicles(request *RideRequest, usedVehicles map[string]bool) []*VehicleStatus {
	candidates := make([]*VehicleStatus, 0)

	for _, vehicle := range d.vehicles {
		// 이미 사용된 차량 제외
		if usedVehicles[vehicle.VehicleID] {
			continue
		}

		// 가용하지 않은 차량 제외
		if !vehicle.Available || vehicle.CurrentRide != nil {
			continue
		}

		// 용량 확인
		if vehicle.Capacity < request.Passengers {
			continue
		}

		// 배터리 확인
		if vehicle.Battery < minBatteryPercent {
			continue
		}

		// 거리 확인
		dist := haversineDistanceDispatch(
			request.PickupLat, request.PickupLon,
			vehicle.Lat, vehicle.Lon,
		)
		if dist > maxSearchRadiusKm {
			continue
		}

		candidates = append(candidates, vehicle)
	}

	return candidates
}

// findBestAssignment 최적 매칭 찾기
func (d *DispatcherSolution) findBestAssignment(request *RideRequest, candidates []*VehicleStatus) *Assignment {
	var bestAssignment *Assignment
	bestScore := math.Inf(-1)

	for _, vehicle := range candidates {
		assignment := d.calculateAssignment(request, vehicle)

		if assignment.Score > bestScore {
			bestScore = assignment.Score
			bestAssignment = &assignment
		}
	}

	return bestAssignment
}

// calculateAssignment 배차 점수 계산
/*
【 스코어 계산 】
- 거리 점수: 가까울수록 높음 (1/distance 정규화)
- 대기 점수: 오래 기다릴수록 높음
- 배터리 점수: 여유 있을수록 높음
*/
func (d *DispatcherSolution) calculateAssignment(request *RideRequest, vehicle *VehicleStatus) Assignment {
	// 거리 계산
	distance := haversineDistanceDispatch(
		request.PickupLat, request.PickupLon,
		vehicle.Lat, vehicle.Lon,
	)

	// 예상 대기 시간
	waitTime := time.Duration(distance/avgSpeedKmH*60) * time.Minute

	// 스코어 계산
	// 1. 거리 점수 (가까울수록 높음, 최대 1.0)
	distanceScore := 1.0 / (1.0 + distance)

	// 2. 대기 시간 보너스 (오래 기다릴수록 높음)
	waitingMinutes := time.Since(request.RequestTime).Minutes()
	waitingScore := waitingMinutes * waitingBonusPerMin
	if waitingScore > 1.0 {
		waitingScore = 1.0 // 최대 1.0
	}

	// 3. 배터리 점수 (100%일 때 1.0)
	batteryScore := float64(vehicle.Battery) / 100.0

	// 가중 합계
	score := weightDistance*distanceScore +
		weightWaiting*waitingScore +
		weightBattery*batteryScore

	return Assignment{
		RideID:      request.ID,
		VehicleID:   vehicle.VehicleID,
		EstWaitTime: waitTime,
		EstDistance: distance,
		Score:       score,
	}
}

// ConfirmAssignment 배차 확정
func (d *DispatcherSolution) ConfirmAssignment(rideID, vehicleID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 요청 확인
	if _, exists := d.pendingRequests[rideID]; !exists {
		return errors.New("ride request not found")
	}

	// 차량 확인
	vehicle, exists := d.vehicles[vehicleID]
	if !exists {
		return errors.New("vehicle not found")
	}
	if !vehicle.Available || vehicle.CurrentRide != nil {
		return errors.New("vehicle not available")
	}

	// 배차 확정
	delete(d.pendingRequests, rideID)
	d.assignedRides[rideID] = vehicleID
	vehicle.CurrentRide = &rideID
	vehicle.Available = false

	return nil
}

// ============================================================================
// 유틸리티 함수
// ============================================================================

func haversineDistanceDispatch(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

func encodeGeohashDispatch(lat, lon float64, precision int) string {
	const base32 = "0123456789bcdefghjkmnpqrstuvwxyz"

	var geohash []byte
	var minLat, maxLat float64 = -90, 90
	var minLon, maxLon float64 = -180, 180

	isEven := true
	bit := 0
	ch := 0

	for len(geohash) < precision {
		if isEven {
			mid := (minLon + maxLon) / 2
			if lon > mid {
				ch |= 1 << (4 - bit)
				minLon = mid
			} else {
				maxLon = mid
			}
		} else {
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

// ============================================================================
// 고급 최적화: 헝가리안 알고리즘
// ============================================================================

/*
【 헝가리안 알고리즘 (Hungarian Algorithm) 】
전역 최적 매칭을 찾는 알고리즘

현재 구현: Greedy - 각 요청에 대해 최선의 차량 선택
문제점: 지역 최적해 (전체적으로는 더 나은 매칭 존재 가능)

헝가리안 알고리즘:
1. 비용 행렬 생성 (요청 × 차량)
2. 행/열 감소
3. 최소 선으로 0 덮기
4. 최적 할당 찾기

시간 복잡도: O(n³)
장점: 전역 최적해 보장
단점: 계산 비용 높음 (대규모에서 사용 어려움)

【 실무 절충안 】
- 소규모 (< 100): 헝가리안 알고리즘
- 중규모 (100-1000): 휴리스틱 + 지역 최적화
- 대규모 (> 1000): Greedy + 주기적 재최적화
*/

// ============================================================================
// 성능 최적화 포인트
// ============================================================================

/*
【 배치 처리 】
- 1초 단위로 요청 모아서 처리
- 더 좋은 전역 최적화 가능

type BatchDispatcher struct {
    batchInterval time.Duration
    requestBuffer chan RideRequest
}

func (d *BatchDispatcher) processBatch() {
    ticker := time.NewTicker(d.batchInterval)
    for range ticker.C {
        requests := d.drainBuffer()
        assignments := d.optimizeAssignments(requests)
        d.notifyAssignments(assignments)
    }
}

【 예측 기반 배차 】
- 과거 데이터로 수요 예측
- 차량 미리 이동 (rebalancing)
- 대기 시간 감소

【 동적 가격 책정 】
- 수요/공급 비율에 따라 가격 조정
- 피크 시간대 차량 공급 유도

┌─────────────────────────────────────────────────────────────────────────────┐
│                       프로덕션 아키텍처                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  [승객 앱] → [API Gateway] → [Request Service] → [Kafka]                   │
│                                                      ↓                      │
│                                            [Dispatch Engine]                │
│                                                      │                      │
│                              ┌───────────────────────┼───────────────────┐  │
│                              ↓                       ↓                   ↓  │
│                    [Demand Predictor]    [Vehicle Tracker]    [Route Engine]│
│                              └───────────────────────┼───────────────────┘  │
│                                                      ↓                      │
│                                            [Assignment Service]             │
│                                                      ↓                      │
│  [차량] ← [Push Notification] ← [Notification Service]                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
*/
