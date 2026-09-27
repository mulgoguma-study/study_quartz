package medium

/*
================================================================================
문제 2: Rate Limiter - 솔루션
================================================================================
*/

import (
	"sync"
	"time"
)

/*
================================================================================
알고리즘 비교
================================================================================

┌─────────────────────────────────────────────────────────────────────────────┐
│                        Rate Limiting 알고리즘 비교                          │
├──────────────────┬────────────────────┬─────────────────────────────────────┤
│     알고리즘     │       특징         │            사용 사례                │
├──────────────────┼────────────────────┼─────────────────────────────────────┤
│ Token Bucket     │ 버스트 허용        │ API Gateway, 네트워크 트래픽        │
│                  │ 메모리 효율적 O(1) │ AWS API Gateway 사용                │
├──────────────────┼────────────────────┼─────────────────────────────────────┤
│ Leaky Bucket     │ 일정한 출력 속도   │ 트래픽 쉐이핑                       │
│                  │ 버스트 불허        │ 네트워크 QoS                        │
├──────────────────┼────────────────────┼─────────────────────────────────────┤
│ Sliding Window   │ 정확한 제한        │ 엄격한 제한이 필요한 경우           │
│ Log              │ 메모리 O(n)        │ 금융 거래 제한                      │
├──────────────────┼────────────────────┼─────────────────────────────────────┤
│ Sliding Window   │ 메모리 효율적      │ Redis Rate Limiting                 │
│ Counter          │ 근사치 계산        │ 대규모 시스템                       │
└──────────────────┴────────────────────┴─────────────────────────────────────┘
*/

// ============================================================================
// 방법 1: Token Bucket (토큰 버킷)
// ============================================================================

/*
【 토큰 버킷 동작 원리 】

┌─────────────────────────────────────────────────────────────────────────────┐
│                          Token Bucket 시각화                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│     토큰 생성 (refillRate/sec)                                              │
│           │                                                                 │
│           ▼                                                                 │
│     ┌───────────┐                                                           │
│     │ ● ● ● ● ○ │  ← 버킷 (capacity = 5, 현재 4개)                         │
│     │   Bucket  │                                                           │
│     └─────┬─────┘                                                           │
│           │                                                                 │
│           ▼ 요청 시 토큰 1개 소비                                           │
│       ┌───────┐                                                             │
│       │Request│                                                             │
│       └───────┘                                                             │
│                                                                             │
│   ● = 토큰 있음, ○ = 빈 슬롯                                                │
│                                                                             │
│   【 특징 】                                                                │
│   - 버스트 허용: 버킷이 가득 차면 capacity만큼 한 번에 처리 가능            │
│   - 평균 처리율: refillRate로 수렴                                          │
│   - 메모리: 사용자당 O(1)                                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

【 왜 Token Bucket이 좋은가? 】
1. 버스트 허용: 순간적인 트래픽 스파이크 처리 가능
2. 메모리 효율: 사용자당 상수 메모리
3. 구현 단순: 타임스탬프와 토큰 수만 저장
4. 실시간 계산: 요청 시점에 토큰 계산 (lazy evaluation)
*/

// userBucket 사용자별 버킷 상태
type userBucket struct {
	tokens     float64   // 현재 토큰 수
	lastRefill time.Time // 마지막 리필 시간
}

// TokenBucketLimiter Token Bucket 구현
type TokenBucketLimiter struct {
	capacity   float64                // 버킷 최대 용량
	refillRate float64                // 초당 토큰 충전 개수
	buckets    map[string]*userBucket // 사용자별 버킷
	mu         sync.Mutex             // 동시성 제어
}

// NewTokenBucketLimiter 생성자
func NewTokenBucketLimiter(capacity int, refillRate float64) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		capacity:   float64(capacity),
		refillRate: refillRate,
		buckets:    make(map[string]*userBucket),
	}
}

// Allow 요청 허용 여부 확인
/*
【 동작 순서 】
1. 사용자 버킷 가져오기 (없으면 생성)
2. 경과 시간에 따라 토큰 리필 (lazy refill)
3. 토큰 있으면 소비하고 true, 없으면 false
*/
func (l *TokenBucketLimiter) Allow(userID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	// 버킷 가져오기 또는 생성
	bucket, exists := l.buckets[userID]
	if !exists {
		// 새 사용자: 버킷을 가득 채워서 시작
		l.buckets[userID] = &userBucket{
			tokens:     l.capacity,
			lastRefill: now,
		}
		bucket = l.buckets[userID]
	}

	// Lazy Refill: 마지막 요청 이후 경과 시간에 따라 토큰 충전
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens += elapsed * l.refillRate

	// 최대 용량 초과 방지
	if bucket.tokens > l.capacity {
		bucket.tokens = l.capacity
	}
	bucket.lastRefill = now

	// 토큰 소비 시도
	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}

	return false
}

// ============================================================================
// 방법 2: Sliding Window Log (슬라이딩 윈도우 로그)
// ============================================================================

/*
【 슬라이딩 윈도우 로그 동작 원리 】

┌─────────────────────────────────────────────────────────────────────────────┐
│                      Sliding Window Log 시각화                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   시간 축 (1분 윈도우, limit=5)                                             │
│                                                                             │
│   ──●────●──●────●────●──●──┼──────────────────▶ 시간                       │
│     │    │  │    │    │  │  │                                               │
│    10:00:05 10:00:20  ...   현재(10:01:00)                                  │
│     ╰─────────────────────╯                                                 │
│       윈도우 내 요청들 (5개) → 다음 요청 거부                               │
│                                                                             │
│   윈도우가 이동하면서:                                                      │
│   - 오래된 요청은 윈도우 밖으로 나감                                        │
│   - 새 요청이 들어올 수 있음                                                │
│                                                                             │
│   【 특징 】                                                                │
│   - 정확한 제한: 어느 시점에서도 정확히 limit개                             │
│   - 메모리: 사용자당 O(요청수)                                              │
│   - 구현: 타임스탬프 리스트 저장                                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

【 왜 Sliding Window Log가 좋은가? 】
1. 정확성: Fixed Window의 경계 문제 없음
2. 공정성: 어느 시점에서 봐도 정확히 limit개 제한
3. 유연성: 다양한 윈도우 크기 적용 가능

【 단점 】
- 메모리: 요청마다 타임스탬프 저장 (많은 요청 시 비효율)
- 해결책: Sliding Window Counter (근사치 사용)
*/

// SlidingWindowLimiter Sliding Window Log 구현
type SlidingWindowLimiter struct {
	limit   int                        // 윈도우 내 최대 요청 수
	window  time.Duration              // 윈도우 크기
	logs    map[string][]time.Time     // 사용자별 요청 타임스탬프
	mu      sync.Mutex                 // 동시성 제어
}

// NewSlidingWindowLimiter 생성자
func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		limit:  limit,
		window: window,
		logs:   make(map[string][]time.Time),
	}
}

// Allow 요청 허용 여부 확인
/*
【 동작 순서 】
1. 현재 시간 기준으로 윈도우 시작점 계산
2. 윈도우 밖의 오래된 타임스탬프 제거
3. 현재 윈도우 내 요청 수 확인
4. limit 미만이면 새 타임스탬프 추가하고 true, 아니면 false
*/
func (l *SlidingWindowLimiter) Allow(userID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-l.window)

	// 사용자 로그 가져오기
	logs := l.logs[userID]

	// 오래된 타임스탬프 제거 (윈도우 밖)
	validLogs := make([]time.Time, 0)
	for _, ts := range logs {
		if ts.After(windowStart) {
			validLogs = append(validLogs, ts)
		}
	}

	// 제한 확인
	if len(validLogs) >= l.limit {
		l.logs[userID] = validLogs
		return false
	}

	// 새 요청 추가
	validLogs = append(validLogs, now)
	l.logs[userID] = validLogs

	return true
}

// ============================================================================
// 방법 3: Sliding Window Counter (보너스 - 메모리 효율적)
// ============================================================================

/*
【 Sliding Window Counter 】

Fixed Window와 Sliding Window Log의 장점을 결합:
- 메모리 효율: 윈도우당 카운터 2개만 저장
- 근사치 계산: 이전 윈도우와 현재 윈도우의 가중 평균

공식:
count = (이전 윈도우 카운트 * (1 - 경과 비율)) + 현재 윈도우 카운트

예시 (윈도우 1분, 현재 윈도우의 30% 지점):
- 이전 윈도우: 80개 요청
- 현재 윈도우: 20개 요청
- 추정 count = 80 * 0.7 + 20 = 76개
*/

type windowCounter struct {
	prevCount    int       // 이전 윈도우 카운트
	currCount    int       // 현재 윈도우 카운트
	windowStart  time.Time // 현재 윈도우 시작 시간
}

type SlidingWindowCounterLimiter struct {
	limit    int
	window   time.Duration
	counters map[string]*windowCounter
	mu       sync.Mutex
}

func NewSlidingWindowCounterLimiter(limit int, window time.Duration) *SlidingWindowCounterLimiter {
	return &SlidingWindowCounterLimiter{
		limit:    limit,
		window:   window,
		counters: make(map[string]*windowCounter),
	}
}

func (l *SlidingWindowCounterLimiter) Allow(userID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	counter, exists := l.counters[userID]
	if !exists {
		counter = &windowCounter{windowStart: now}
		l.counters[userID] = counter
	}

	// 윈도우 경과 확인
	elapsed := now.Sub(counter.windowStart)

	// 새 윈도우로 전환
	if elapsed >= l.window {
		// 몇 윈도우가 지났는지 계산
		windowsPassed := int(elapsed / l.window)
		if windowsPassed == 1 {
			counter.prevCount = counter.currCount
		} else {
			counter.prevCount = 0
		}
		counter.currCount = 0
		counter.windowStart = counter.windowStart.Add(time.Duration(windowsPassed) * l.window)
		elapsed = now.Sub(counter.windowStart)
	}

	// 가중 평균 계산
	weight := float64(l.window-elapsed) / float64(l.window)
	estimatedCount := float64(counter.prevCount)*weight + float64(counter.currCount)

	if int(estimatedCount) >= l.limit {
		return false
	}

	counter.currCount++
	return true
}

/*
================================================================================
실무 적용 가이드
================================================================================

【 알고리즘 선택 기준 】

1. Token Bucket
   - 버스트 허용이 필요한 경우
   - API Gateway, CDN
   - 메모리가 제한적인 경우

2. Sliding Window Log
   - 정확한 제한이 필요한 경우
   - 금융 거래, 보안 관련
   - 요청 수가 적은 경우

3. Sliding Window Counter
   - 대규모 시스템
   - Redis 기반 분산 환경
   - 메모리와 정확성의 균형

【 분산 환경에서의 Rate Limiting 】

1. Redis + Lua Script
   - 원자적 연산 보장
   - 여러 서버에서 공유

2. Redis Cell (모듈)
   - Token Bucket 네이티브 구현
   - 고성능

3. Envoy/Nginx Rate Limiting
   - 인프라 레벨에서 처리
   - 애플리케이션 부하 감소

【 모니터링 포인트 】
- 제한 발생 횟수
- 사용자별 요청 패턴
- 버스트 트래픽 감지
================================================================================
*/

// 테스트
func main() {
	println("=== Token Bucket Test ===")
	tb := NewTokenBucketLimiter(5, 5.0)
	for i := 0; i < 7; i++ {
		result := tb.Allow("user1")
		println("Request", i+1, ":", result)
	}

	println("\n=== Sliding Window Test ===")
	sw := NewSlidingWindowLimiter(3, time.Second)
	for i := 0; i < 5; i++ {
		result := sw.Allow("user1")
		println("Request", i+1, ":", result)
	}

	println("\n=== Sliding Window Counter Test ===")
	swc := NewSlidingWindowCounterLimiter(3, time.Second)
	for i := 0; i < 5; i++ {
		result := swc.Allow("user1")
		println("Request", i+1, ":", result)
	}
}
