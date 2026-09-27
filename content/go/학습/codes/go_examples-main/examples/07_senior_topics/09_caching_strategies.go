package senior_topics

/*
================================================================================
캐싱 전략: In-Memory 위험성과 Look-aside vs Write-through (시니어 레벨)
================================================================================

면접 질문:
1. "In-Memory Cache의 위험성은?"
2. "Look-aside vs Write-through 전략을 설명해주세요"

시니어급 답변 포인트:
1. In-Memory Cache의 장단점과 위험성
2. Cache-aside (Look-aside) 전략
3. Write-through 전략
4. Write-behind (Write-back) 전략
5. 캐시 일관성 문제
6. 분산 환경에서의 고려사항
*/

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// 1. In-Memory Cache 위험성
// ============================================================================

/*
In-Memory Cache 장점:
- 매우 빠른 접근 속도 (나노초)
- 네트워크 지연 없음
- 구현이 단순함

In-Memory Cache 위험성:

1. 메모리 부족 (OOM):
   - 캐시가 무한정 증가
   - TTL 없이 데이터 누적
   - 대용량 객체 캐싱

2. 데이터 불일치:
   - 다중 인스턴스 환경에서 캐시 동기화 불가
   - 각 인스턴스가 다른 캐시 상태
   - DB 업데이트 후 캐시 무효화 누락

3. Cold Start 문제:
   - 재시작 시 캐시 완전 유실
   - Thundering Herd (동시 요청 폭주)
   - 캐시 워밍업 필요

4. GC 부담:
   - 많은 객체가 Old Gen으로 이동
   - Full GC 시간 증가
   - 애플리케이션 일시 정지

5. 단일 장애점:
   - 캐시 서버 다운 시 전체 성능 저하
   - Fallback 전략 필요
*/

// 위험한 구현 예시 - 하지 말 것!
type DangerousCache struct {
	data map[string]interface{} // TTL 없음!
	mu   sync.RWMutex
}

func (c *DangerousCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value // 무한정 증가 가능!
}

// 안전한 구현 - TTL과 최대 크기 제한
type SafeCache struct {
	data     map[string]*cacheEntry
	mu       sync.RWMutex
	maxSize  int
	ttl      time.Duration
	onEvict  func(key string, value interface{})
}

type cacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

func NewSafeCache(maxSize int, ttl time.Duration) *SafeCache {
	cache := &SafeCache{
		data:    make(map[string]*cacheEntry),
		maxSize: maxSize,
		ttl:     ttl,
	}
	go cache.cleanupExpired()
	return cache
}

func (c *SafeCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 최대 크기 체크
	if len(c.data) >= c.maxSize {
		c.evictOldest()
	}

	c.data[key] = &cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *SafeCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.data[key]
	if !exists || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.value, true
}

func (c *SafeCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range c.data {
		if oldestKey == "" || entry.expiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.expiresAt
		}
	}

	if oldestKey != "" {
		if c.onEvict != nil {
			c.onEvict(oldestKey, c.data[oldestKey].value)
		}
		delete(c.data, oldestKey)
	}
}

func (c *SafeCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.data {
			if now.After(entry.expiresAt) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}

// ============================================================================
// 2. Cache-Aside (Look-Aside) 전략
// ============================================================================

/*
Cache-Aside (Look-Aside) 패턴:

읽기:
1. 캐시에서 조회
2. 캐시 히트: 값 반환
3. 캐시 미스: DB 조회 → 캐시에 저장 → 값 반환

쓰기:
1. DB에 쓰기
2. 캐시 무효화 (삭제)

장점:
- 간단하고 직관적
- 읽기 워크로드에 최적화
- 캐시 실패해도 서비스 가능

단점:
- 캐시 미스 시 DB 부하
- 데이터 불일치 가능 (짧은 윈도우)
- 애플리케이션이 캐시 로직 관리

사용 시점:
- 읽기가 많은 워크로드
- 데이터 일관성 요구가 상대적으로 낮은 경우
- 캐시 실패 허용 가능

                     ┌─────────────────┐
                     │   Application   │
                     └────────┬────────┘
                              │
           ┌──────────────────┼──────────────────┐
           │                  │                  │
           ▼                  │                  ▼
    ┌─────────────┐    2. Cache Miss     ┌─────────────┐
    │    Cache    │ ◄───────────────────│   Database  │
    └─────────────┘                      └─────────────┘
          │                                     │
          │ 1. Check Cache                      │ 3. Store in Cache
          ▼                                     ▼
    Cache Hit: Return                    Return & Cache
*/

// Cache-Aside 구현
type CacheAside struct {
	cache      CacheInterface
	db         DatabaseInterface
	ttl        time.Duration
}

type CacheInterface interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type DatabaseInterface interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}) error
}

var ErrCacheMiss = errors.New("cache miss")

func NewCacheAside(cache CacheInterface, db DatabaseInterface, ttl time.Duration) *CacheAside {
	return &CacheAside{
		cache: cache,
		db:    db,
		ttl:   ttl,
	}
}

// 읽기: 캐시 먼저 → DB → 캐시에 저장
func (ca *CacheAside) Get(ctx context.Context, key string) (interface{}, error) {
	// 1. 캐시에서 조회
	value, err := ca.cache.Get(ctx, key)
	if err == nil {
		return value, nil // 캐시 히트
	}

	// 2. 캐시 미스 → DB 조회
	value, err = ca.db.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	// 3. 캐시에 저장 (비동기로 해도 됨)
	_ = ca.cache.Set(ctx, key, value, ca.ttl)

	return value, nil
}

// 쓰기: DB 먼저 → 캐시 무효화
func (ca *CacheAside) Set(ctx context.Context, key string, value interface{}) error {
	// 1. DB에 쓰기
	if err := ca.db.Set(ctx, key, value); err != nil {
		return err
	}

	// 2. 캐시 무효화 (업데이트가 아닌 삭제!)
	// 이유: 캐시 업데이트와 DB 업데이트 순서 문제 방지
	_ = ca.cache.Delete(ctx, key)

	return nil
}

// ============================================================================
// 3. Write-Through 전략
// ============================================================================

/*
Write-Through 패턴:

읽기:
1. 캐시에서 조회
2. 캐시 히트: 값 반환
3. 캐시 미스: DB 조회 → 캐시에 저장 → 값 반환

쓰기:
1. 캐시에 쓰기
2. DB에 쓰기 (동기)

장점:
- 캐시와 DB 항상 동기화
- 읽기 시 항상 캐시 히트 (쓰기 후)
- 데이터 일관성 높음

단점:
- 쓰기 지연 증가 (캐시 + DB 모두 쓰기)
- 불필요한 데이터도 캐싱될 수 있음
- 캐시 장애 시 쓰기 실패 가능

사용 시점:
- 데이터 일관성이 중요한 경우
- 쓰기 후 즉시 읽기가 필요한 경우

                     ┌─────────────────┐
                     │   Application   │
                     └────────┬────────┘
                              │
                              │ Write
                              ▼
                     ┌─────────────────┐
                     │     Cache       │──────┐
                     └────────┬────────┘      │
                              │               │ Sync Write
                              │               ▼
                     ┌────────▼────────┐
                     │    Database     │
                     └─────────────────┘
*/

type WriteThrough struct {
	cache CacheInterface
	db    DatabaseInterface
	ttl   time.Duration
}

func NewWriteThrough(cache CacheInterface, db DatabaseInterface, ttl time.Duration) *WriteThrough {
	return &WriteThrough{
		cache: cache,
		db:    db,
		ttl:   ttl,
	}
}

func (wt *WriteThrough) Get(ctx context.Context, key string) (interface{}, error) {
	// Cache-Aside와 동일한 읽기 로직
	value, err := wt.cache.Get(ctx, key)
	if err == nil {
		return value, nil
	}

	value, err = wt.db.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	_ = wt.cache.Set(ctx, key, value, wt.ttl)
	return value, nil
}

// 쓰기: 캐시와 DB 모두 동기적으로 쓰기
func (wt *WriteThrough) Set(ctx context.Context, key string, value interface{}) error {
	// 1. 캐시에 쓰기
	if err := wt.cache.Set(ctx, key, value, wt.ttl); err != nil {
		return fmt.Errorf("cache write failed: %w", err)
	}

	// 2. DB에 쓰기 (동기)
	if err := wt.db.Set(ctx, key, value); err != nil {
		// 롤백: 캐시에서 삭제
		_ = wt.cache.Delete(ctx, key)
		return fmt.Errorf("db write failed: %w", err)
	}

	return nil
}

// ============================================================================
// 4. Write-Behind (Write-Back) 전략
// ============================================================================

/*
Write-Behind (Write-Back) 패턴:

쓰기:
1. 캐시에만 즉시 쓰기
2. 비동기로 DB에 배치 쓰기

장점:
- 매우 빠른 쓰기 응답
- DB 부하 감소 (배치 처리)
- 쓰기 횟수 줄임 (같은 키 여러 번 쓰기 합침)

단점:
- 데이터 손실 위험 (캐시 장애 시)
- 복잡한 구현
- 데이터 일관성 보장 어려움

사용 시점:
- 쓰기 성능이 중요한 경우
- 약간의 데이터 손실 허용 가능
- 배치 처리에 적합한 경우

                     ┌─────────────────┐
                     │   Application   │
                     └────────┬────────┘
                              │
                              │ Write (Immediate)
                              ▼
                     ┌─────────────────┐
                     │     Cache       │
                     │  (Write Queue)  │
                     └────────┬────────┘
                              │
                              │ Async Batch Write
                              ▼
                     ┌─────────────────┐
                     │    Database     │
                     └─────────────────┘
*/

type WriteBehind struct {
	cache      CacheInterface
	db         DatabaseInterface
	ttl        time.Duration
	writeQueue chan writeRequest
	mu         sync.RWMutex
	pending    map[string]interface{} // 쓰기 대기 중인 데이터
}

type writeRequest struct {
	key   string
	value interface{}
}

func NewWriteBehind(cache CacheInterface, db DatabaseInterface, ttl time.Duration) *WriteBehind {
	wb := &WriteBehind{
		cache:      cache,
		db:         db,
		ttl:        ttl,
		writeQueue: make(chan writeRequest, 1000),
		pending:    make(map[string]interface{}),
	}
	go wb.processWrites()
	return wb
}

func (wb *WriteBehind) Set(ctx context.Context, key string, value interface{}) error {
	// 1. 캐시에 즉시 쓰기
	if err := wb.cache.Set(ctx, key, value, wb.ttl); err != nil {
		return err
	}

	// 2. 쓰기 큐에 추가 (비동기)
	wb.mu.Lock()
	wb.pending[key] = value
	wb.mu.Unlock()

	select {
	case wb.writeQueue <- writeRequest{key: key, value: value}:
		// 큐에 추가됨
	default:
		// 큐 가득 참 - 동기 쓰기로 폴백
		return wb.db.Set(ctx, key, value)
	}

	return nil
}

func (wb *WriteBehind) processWrites() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	batch := make(map[string]interface{})

	for {
		select {
		case req := <-wb.writeQueue:
			batch[req.key] = req.value

			// 배치 크기 도달 시 즉시 처리
			if len(batch) >= 100 {
				wb.flushBatch(batch)
				batch = make(map[string]interface{})
			}

		case <-ticker.C:
			// 주기적 플러시
			if len(batch) > 0 {
				wb.flushBatch(batch)
				batch = make(map[string]interface{})
			}
		}
	}
}

func (wb *WriteBehind) flushBatch(batch map[string]interface{}) {
	ctx := context.Background()
	for key, value := range batch {
		if err := wb.db.Set(ctx, key, value); err != nil {
			// 실패한 항목 재시도 큐에 추가
			fmt.Printf("DB write failed for key %s: %v\n", key, err)
		}

		wb.mu.Lock()
		delete(wb.pending, key)
		wb.mu.Unlock()
	}
}

// ============================================================================
// 5. 전략 비교 및 선택 가이드
// ============================================================================

/*
전략 비교:

| 전략          | 읽기 성능 | 쓰기 성능 | 일관성  | 복잡도 | 데이터 손실 |
|--------------|---------|---------|--------|-------|-----------|
| Cache-Aside  | 높음     | 중간     | 낮음    | 낮음   | 낮음       |
| Write-Through| 높음     | 낮음     | 높음    | 중간   | 낮음       |
| Write-Behind | 높음     | 높음     | 낮음    | 높음   | 중간       |

선택 기준:

Cache-Aside:
- 읽기가 많은 워크로드
- 데이터 일관성 요구가 낮음
- 구현 단순성 중요

Write-Through:
- 데이터 일관성 중요
- 쓰기 후 즉시 읽기 필요
- 쓰기 지연 허용

Write-Behind:
- 쓰기 성능 중요
- 약간의 데이터 손실 허용
- 배치 처리 가능
*/

// ============================================================================
// 6. 캐시 무효화 전략
// ============================================================================

/*
캐시 무효화 전략:

1. TTL (Time-To-Live):
   - 가장 단순
   - 일정 시간 후 자동 만료
   - 데이터 신선도 vs 캐시 효율 트레이드오프

2. 이벤트 기반 무효화:
   - DB 변경 시 캐시 무효화
   - CDC (Change Data Capture) 활용
   - 메시지 큐로 무효화 이벤트 전파

3. 수동 무효화:
   - 특정 작업 후 명시적 무효화
   - 제어가 쉬움
   - 누락 위험

분산 환경에서의 무효화:
- Redis Pub/Sub로 모든 인스턴스에 전파
- 캐시 태그로 관련 캐시 일괄 무효화
*/

type CacheInvalidator struct {
	cache    CacheInterface
	pubsub   PubSubInterface
}

type PubSubInterface interface {
	Publish(ctx context.Context, channel string, message interface{}) error
	Subscribe(ctx context.Context, channel string) (<-chan string, error)
}

func (ci *CacheInvalidator) InvalidateAndNotify(ctx context.Context, key string) error {
	// 로컬 캐시 무효화
	if err := ci.cache.Delete(ctx, key); err != nil {
		return err
	}

	// 다른 인스턴스에 알림
	return ci.pubsub.Publish(ctx, "cache:invalidate", key)
}

func (ci *CacheInvalidator) StartSubscriber(ctx context.Context) error {
	ch, err := ci.pubsub.Subscribe(ctx, "cache:invalidate")
	if err != nil {
		return err
	}

	go func() {
		for key := range ch {
			_ = ci.cache.Delete(ctx, key)
		}
	}()

	return nil
}

// ============================================================================
// 7. 시니어 면접 답변 요약
// ============================================================================

/*
Q: In-Memory Cache의 위험성은?

A:
1. 메모리 부족 (OOM):
   - TTL과 최대 크기 제한 필수
   - LRU 등 eviction 정책 적용

2. 데이터 불일치:
   - 다중 인스턴스에서 동기화 불가
   - Redis 같은 분산 캐시 고려

3. Cold Start:
   - 재시작 시 캐시 유실
   - 캐시 워밍업 전략 필요

4. 단일 장애점:
   - 폴백 전략 필수
   - 캐시 실패 시 DB 직접 조회


Q: Look-aside vs Write-through 전략?

A:
Cache-Aside (Look-Aside):
- 읽기: 캐시 → (미스 시) DB → 캐시 저장
- 쓰기: DB → 캐시 삭제
- 장점: 단순, 읽기 최적화
- 단점: 캐시 미스 시 지연

Write-Through:
- 읽기: 캐시 → (미스 시) DB → 캐시 저장
- 쓰기: 캐시 → DB (동기)
- 장점: 높은 일관성, 항상 캐시 히트
- 단점: 쓰기 지연

선택 기준:
- 읽기 중심, 일관성 낮음 → Cache-Aside
- 일관성 중요, 쓰기 후 읽기 → Write-Through
- 쓰기 성능 중요 → Write-Behind (주의: 데이터 손실 위험)
*/
