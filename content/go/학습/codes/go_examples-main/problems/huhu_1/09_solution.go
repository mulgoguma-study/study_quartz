package huhu1

/*
================================================================================
문제 9: API 게이트웨이 (Rate Limiting + Caching) - 솔루션
================================================================================

【 핵심 설계 원칙 】
1. Sliding Window Counter - 정확한 Rate Limiting
2. LRU Cache + TTL - 효율적인 캐시 관리
3. Singleflight - Cache Stampede 방지
4. Lock Striping - 동시성 최적화

【 Rate Limiting 알고리즘 비교 】
┌────────────────────┬─────────────┬─────────────┬─────────────┐
│ 알고리즘           │ 정확도      │ 메모리      │ 구현 복잡도  │
├────────────────────┼─────────────┼─────────────┼─────────────┤
│ Fixed Window       │ 낮음        │ O(1)        │ 낮음        │
│ Sliding Window Log │ 높음        │ O(n)        │ 중간        │
│ Sliding Window Cnt │ 중간        │ O(1)        │ 중간        │
│ Token Bucket       │ 높음        │ O(1)        │ 중간        │
└────────────────────┴─────────────┴─────────────┴─────────────┘

================================================================================
*/

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// Sliding Window Counter (Rate Limiting)
// ============================================================================

type slidingWindowCounter struct {
	mu           sync.Mutex
	prevCount    int
	currCount    int
	windowStart  time.Time
	windowSize   time.Duration
}

func (s *slidingWindowCounter) allow(limit int, now time.Time) (bool, int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 현재 윈도우 확인
	elapsed := now.Sub(s.windowStart)

	if elapsed >= s.windowSize*2 {
		// 2 윈도우 이상 지남 → 리셋
		s.prevCount = 0
		s.currCount = 0
		s.windowStart = now.Truncate(s.windowSize)
	} else if elapsed >= s.windowSize {
		// 1 윈도우 지남 → 슬라이딩
		s.prevCount = s.currCount
		s.currCount = 0
		s.windowStart = s.windowStart.Add(s.windowSize)
	}

	// 가중 평균으로 현재 요청 수 계산
	windowProgress := float64(now.Sub(s.windowStart)) / float64(s.windowSize)
	weightedCount := float64(s.prevCount)*(1-windowProgress) + float64(s.currCount)

	if int(weightedCount) >= limit {
		remaining := limit - int(weightedCount)
		if remaining < 0 {
			remaining = 0
		}
		return false, remaining
	}

	s.currCount++
	return true, limit - int(weightedCount) - 1
}

func (s *slidingWindowCounter) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prevCount = 0
	s.currCount = 0
	s.windowStart = time.Now()
}

// ============================================================================
// LRU Cache with TTL
// ============================================================================

type lruCache struct {
	mu       sync.RWMutex
	capacity int
	items    map[string]*list.Element
	order    *list.List
}

type lruItem struct {
	key       string
	entry     CacheEntry
}

func newLRUCache(capacity int) *lruCache {
	return &lruCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		order:    list.New(),
	}
}

func (c *lruCache) get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		return nil, false
	}

	item := elem.Value.(*lruItem)

	// TTL 확인
	if time.Now().After(item.entry.ExpiresAt) {
		c.removeElement(elem)
		return nil, false
	}

	// LRU: 맨 앞으로 이동
	c.order.MoveToFront(elem)
	return item.entry.Value, true
}

func (c *lruCache) set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	entry := CacheEntry{
		Value:     value,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}

	// 기존 키 있으면 업데이트
	if elem, ok := c.items[key]; ok {
		elem.Value.(*lruItem).entry = entry
		c.order.MoveToFront(elem)
		return
	}

	// 용량 초과 시 가장 오래된 항목 제거
	if c.order.Len() >= c.capacity {
		oldest := c.order.Back()
		if oldest != nil {
			c.removeElement(oldest)
		}
	}

	// 새 항목 추가
	item := &lruItem{key: key, entry: entry}
	elem := c.order.PushFront(item)
	c.items[key] = elem
}

func (c *lruCache) delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.removeElement(elem)
	}
}

func (c *lruCache) removeElement(elem *list.Element) {
	item := elem.Value.(*lruItem)
	delete(c.items, item.key)
	c.order.Remove(elem)
}

// ============================================================================
// Singleflight (Cache Stampede 방지)
// ============================================================================

type singleflight struct {
	mu    sync.Mutex
	calls map[string]*call
}

type call struct {
	wg    sync.WaitGroup
	val   any
	err   error
}

func newSingleflight() *singleflight {
	return &singleflight{
		calls: make(map[string]*call),
	}
}

func (s *singleflight) do(key string, fn func() (any, error)) (any, error) {
	s.mu.Lock()

	// 이미 실행 중인 요청이 있으면 대기
	if c, ok := s.calls[key]; ok {
		s.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	// 새 요청 등록
	c := &call{}
	c.wg.Add(1)
	s.calls[key] = c
	s.mu.Unlock()

	// 함수 실행
	c.val, c.err = fn()
	c.wg.Done()

	// 완료 후 제거
	s.mu.Lock()
	delete(s.calls, key)
	s.mu.Unlock()

	return c.val, c.err
}

// ============================================================================
// API Gateway 구현
// ============================================================================

// APIGatewaySolution 완성된 API 게이트웨이
type APIGatewaySolution struct {
	// Rate Limiting
	rateLimiters   map[string]*slidingWindowCounter
	rateLimitersMu sync.RWMutex

	// Caching
	cache *lruCache

	// Singleflight
	sf *singleflight

	// 통계
	stats struct {
		totalRequests int64
		rateLimited   int64
		cacheHits     int64
		cacheMisses   int64
	}
}

// NewAPIGatewaySolution 생성자
func NewAPIGatewaySolution(cacheSize int) *APIGatewaySolution {
	return &APIGatewaySolution{
		rateLimiters: make(map[string]*slidingWindowCounter),
		cache:        newLRUCache(cacheSize),
		sf:           newSingleflight(),
	}
}

// ============================================================================
// Rate Limiting
// ============================================================================

// RateLimit Rate Limiting 확인
/*
【 Sliding Window Counter 】
- 이전 윈도우와 현재 윈도우의 가중 평균
- Fixed Window의 경계 문제 해결
- 메모리 효율적 (O(1))
*/
func (g *APIGatewaySolution) RateLimit(key string, limit int, window time.Duration) RateLimitResult {
	atomic.AddInt64(&g.stats.totalRequests, 1)

	// Rate Limiter 가져오기 또는 생성
	g.rateLimitersMu.Lock()
	limiter, ok := g.rateLimiters[key]
	if !ok {
		limiter = &slidingWindowCounter{
			windowStart: time.Now().Truncate(window),
			windowSize:  window,
		}
		g.rateLimiters[key] = limiter
	}
	g.rateLimitersMu.Unlock()

	now := time.Now()
	allowed, remaining := limiter.allow(limit, now)

	result := RateLimitResult{
		Allowed:   allowed,
		Remaining: remaining,
		ResetAt:   limiter.windowStart.Add(window),
	}

	if !allowed {
		atomic.AddInt64(&g.stats.rateLimited, 1)
		result.RetryAfter = result.ResetAt.Sub(now)
	}

	return result
}

// ResetRateLimit Rate Limit 리셋
func (g *APIGatewaySolution) ResetRateLimit(key string) error {
	g.rateLimitersMu.Lock()
	defer g.rateLimitersMu.Unlock()

	if limiter, ok := g.rateLimiters[key]; ok {
		limiter.reset()
	}

	return nil
}

// ============================================================================
// Caching
// ============================================================================

// Get 캐시 조회
func (g *APIGatewaySolution) Get(key string) (any, bool) {
	value, ok := g.cache.get(key)

	if ok {
		atomic.AddInt64(&g.stats.cacheHits, 1)
	} else {
		atomic.AddInt64(&g.stats.cacheMisses, 1)
	}

	return value, ok
}

// Set 캐시 저장
func (g *APIGatewaySolution) Set(key string, value any, ttl time.Duration) error {
	g.cache.set(key, value, ttl)
	return nil
}

// Delete 캐시 삭제
func (g *APIGatewaySolution) Delete(key string) error {
	g.cache.delete(key)
	return nil
}

// GetOrSet 캐시 조회 또는 로더 실행
/*
【 Cache Stampede 방지 】
Singleflight로 동일 키에 대해 하나의 요청만 실행
- 캐시 미스 시 다수 요청이 동시에 DB 조회하는 것 방지
- 첫 번째 요청 결과를 나머지 요청이 공유
*/
func (g *APIGatewaySolution) GetOrSet(key string, loader func() (any, error), ttl time.Duration) (any, error) {
	// 캐시 확인
	if value, ok := g.Get(key); ok {
		return value, nil
	}

	// Singleflight로 로더 실행
	value, err := g.sf.do(key, func() (any, error) {
		// 다시 캐시 확인 (다른 요청이 이미 저장했을 수 있음)
		if v, ok := g.cache.get(key); ok {
			return v, nil
		}

		// 로더 실행
		v, err := loader()
		if err != nil {
			return nil, err
		}

		// 캐시 저장
		g.cache.set(key, v, ttl)
		return v, nil
	})

	return value, err
}

// ============================================================================
// Stats
// ============================================================================

// GetStats 통계 반환
func (g *APIGatewaySolution) GetStats() GatewayStats {
	hits := atomic.LoadInt64(&g.stats.cacheHits)
	misses := atomic.LoadInt64(&g.stats.cacheMisses)

	var hitRate float64
	total := hits + misses
	if total > 0 {
		hitRate = float64(hits) / float64(total) * 100
	}

	return GatewayStats{
		TotalRequests: atomic.LoadInt64(&g.stats.totalRequests),
		RateLimited:   atomic.LoadInt64(&g.stats.rateLimited),
		CacheHits:     hits,
		CacheMisses:   misses,
		HitRate:       hitRate,
	}
}

// ============================================================================
// 성능 최적화 포인트
// ============================================================================

/*
【 프로덕션 최적화 】

1. 분산 Rate Limiting (Redis)
   - Lua Script로 원자적 연산
   - 여러 서버 간 제한 공유

-- Redis Lua Script (Sliding Window)
local key = KEYS[1]
local window = tonumber(ARGV[1])
local limit = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

-- 오래된 요청 제거
redis.call('ZREMRANGEBYSCORE', key, 0, now - window)

-- 현재 요청 수
local count = redis.call('ZCARD', key)

if count < limit then
    -- 요청 허용
    redis.call('ZADD', key, now, now .. '-' .. math.random())
    redis.call('EXPIRE', key, window / 1000)
    return {1, limit - count - 1}
else
    -- 요청 거부
    return {0, 0}
end

2. 분산 캐시 (Redis)
   - 로컬 캐시 + Redis 2단계
   - Cache-Aside 패턴

3. 캐시 무효화 전략
   - TTL 기반 (시간 만료)
   - Event 기반 (데이터 변경 시)
   - Write-through (쓰기 시 갱신)

4. 핫 키 처리
   - 로컬 캐시로 핫 키 처리
   - 복제로 부하 분산

┌─────────────────────────────────────────────────────────────────────────────┐
│                       프로덕션 아키텍처                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  [Client] → [Load Balancer] → [API Gateway (Rate Limit)]                   │
│                                        │                                    │
│                                        ↓                                    │
│                    ┌───────────────────┼───────────────────┐               │
│                    │                   │                   │               │
│                    ↓                   ↓                   ↓               │
│              [Local Cache]      [Redis Cache]        [API Server]          │
│              (핫 데이터)         (분산 캐시)          (비즈니스)            │
│                    │                   │                   │               │
│                    └───────────────────┼───────────────────┘               │
│                                        ↓                                    │
│                                   [Database]                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

【 HTTP 헤더 예시 】

Rate Limit 응답 헤더:
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1640000000
Retry-After: 30  (429 응답 시)

캐시 응답 헤더:
Cache-Control: max-age=300
ETag: "abc123"
X-Cache: HIT  (또는 MISS)
*/
