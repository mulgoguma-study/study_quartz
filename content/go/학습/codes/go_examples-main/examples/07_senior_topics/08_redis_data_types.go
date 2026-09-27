package senior_topics

/*
================================================================================
Redis Data Types 활용 가이드 (시니어 레벨)
================================================================================

면접 질문: "Redis의 주요 데이터 타입과 각각의 활용 사례를 설명해주세요"

시니어급 답변 포인트:
1. String - 기본, 카운터, 캐시
2. Hash - 객체 저장
3. List - 큐, 스택, 최근 항목
4. Set - 고유 항목, 태그
5. Sorted Set - 랭킹, 타임라인
6. Stream - 이벤트 스트림 (Kafka-like)
7. Pub/Sub - 실시간 메시징
*/

import (
	"context"
	"fmt"
	"time"
)

// Redis 클라이언트 인터페이스 (실제 구현은 go-redis/redis 사용)
type RedisClient interface {
	// String
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Incr(ctx context.Context, key string) (int64, error)
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)

	// Hash
	HSet(ctx context.Context, key string, values ...interface{}) error
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error)

	// List
	LPush(ctx context.Context, key string, values ...interface{}) error
	RPush(ctx context.Context, key string, values ...interface{}) error
	LPop(ctx context.Context, key string) (string, error)
	RPop(ctx context.Context, key string) (string, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	LTrim(ctx context.Context, key string, start, stop int64) error

	// Set
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	SInter(ctx context.Context, keys ...string) ([]string, error)

	// Sorted Set
	ZAdd(ctx context.Context, key string, members ...ZMember) error
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]ZMember, error)
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]ZMember, error)
	ZScore(ctx context.Context, key string, member string) (float64, error)
	ZRank(ctx context.Context, key string, member string) (int64, error)
	ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error)
}

type ZMember struct {
	Score  float64
	Member string
}

// ============================================================================
// 1. String - 가장 기본적인 타입
// ============================================================================

/*
String 특징:
- 최대 512MB 저장 가능
- Binary-safe (이미지, 직렬화된 객체 등)
- Atomic 연산 지원 (INCR, DECR)

활용 사례:
1. 단순 캐시 (JSON, 직렬화된 객체)
2. 카운터 (페이지뷰, API 호출 횟수)
3. 분산 락 (SETNX + TTL)
4. Rate Limiting

명령어:
- SET, GET, DEL
- INCR, DECR, INCRBY
- SETNX, SETEX
- MSET, MGET (다중 키)
*/

// 캐시 예시
type UserCache struct {
	client RedisClient
}

func (c *UserCache) GetUser(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("user:%s", userID)

	// 캐시 조회
	data, err := c.client.Get(ctx, key)
	if err == nil {
		return data, nil // 캐시 히트
	}

	// 캐시 미스 - DB 조회 후 캐싱
	// userData := fetchFromDB(userID)
	userData := `{"id":"123","name":"John"}` // 시뮬레이션

	// 1시간 TTL로 캐시
	_ = c.client.Set(ctx, key, userData, 1*time.Hour)

	return userData, nil
}

// 분산 락 예시
type DistributedLock struct {
	client RedisClient
}

func (l *DistributedLock) TryLock(ctx context.Context, resource string, ttl time.Duration) (bool, error) {
	key := fmt.Sprintf("lock:%s", resource)
	// SETNX: 키가 없을 때만 설정 (atomic)
	acquired, err := l.client.SetNX(ctx, key, "locked", ttl)
	return acquired, err
}

// Rate Limiter 예시 (Sliding Window)
type RateLimiter struct {
	client RedisClient
}

func (r *RateLimiter) IsAllowed(ctx context.Context, userID string, limit int64) (bool, error) {
	key := fmt.Sprintf("ratelimit:%s:%d", userID, time.Now().Unix()/60) // 분당

	count, err := r.client.Incr(ctx, key)
	if err != nil {
		return false, err
	}

	// 첫 요청이면 TTL 설정
	if count == 1 {
		_ = r.client.Set(ctx, key, count, 1*time.Minute)
	}

	return count <= limit, nil
}

// ============================================================================
// 2. Hash - 객체 저장에 최적
// ============================================================================

/*
Hash 특징:
- 필드-값 쌍의 집합
- 부분 업데이트 가능 (특정 필드만)
- 메모리 효율적 (작은 해시는 ziplist로 저장)

활용 사례:
1. 사용자 프로필 (개별 필드 업데이트)
2. 세션 데이터
3. 쇼핑 카트
4. 설정 정보

명령어:
- HSET, HGET, HDEL
- HMSET, HMGET (다중 필드)
- HGETALL
- HINCRBY (필드 증가)
- HKEYS, HVALS

vs String(JSON):
- Hash: 부분 업데이트 가능, 필드별 접근
- String: 전체 조회/업데이트, 직렬화 필요
*/

// 사용자 프로필 저장
type UserProfile struct {
	client RedisClient
}

func (u *UserProfile) SetProfile(ctx context.Context, userID string, profile map[string]interface{}) error {
	key := fmt.Sprintf("user:profile:%s", userID)

	// HSET으로 여러 필드 설정
	args := make([]interface{}, 0, len(profile)*2)
	for k, v := range profile {
		args = append(args, k, v)
	}
	return u.client.HSet(ctx, key, args...)
}

func (u *UserProfile) GetProfile(ctx context.Context, userID string) (map[string]string, error) {
	key := fmt.Sprintf("user:profile:%s", userID)
	return u.client.HGetAll(ctx, key)
}

func (u *UserProfile) UpdateField(ctx context.Context, userID, field, value string) error {
	key := fmt.Sprintf("user:profile:%s", userID)
	return u.client.HSet(ctx, key, field, value)
}

// 쇼핑 카트
type ShoppingCart struct {
	client RedisClient
}

func (c *ShoppingCart) AddItem(ctx context.Context, userID, productID string, quantity int64) error {
	key := fmt.Sprintf("cart:%s", userID)
	_, err := c.client.HIncrBy(ctx, key, productID, quantity)
	return err
}

func (c *ShoppingCart) GetCart(ctx context.Context, userID string) (map[string]string, error) {
	key := fmt.Sprintf("cart:%s", userID)
	return c.client.HGetAll(ctx, key)
}

// ============================================================================
// 3. List - 순서가 있는 컬렉션
// ============================================================================

/*
List 특징:
- 양방향 연결 리스트
- 양쪽 끝에서 O(1) 삽입/삭제
- 인덱스 접근 O(N)

활용 사례:
1. 메시지 큐 (LPUSH + RPOP)
2. 스택 (LPUSH + LPOP)
3. 최근 활동 로그
4. 타임라인 (최근 N개)

명령어:
- LPUSH, RPUSH (삽입)
- LPOP, RPOP (삭제)
- LRANGE (범위 조회)
- LTRIM (범위 유지)
- BLPOP, BRPOP (블로킹 팝)
*/

// 메시지 큐
type MessageQueue struct {
	client RedisClient
}

func (q *MessageQueue) Enqueue(ctx context.Context, queueName, message string) error {
	return q.client.LPush(ctx, queueName, message)
}

func (q *MessageQueue) Dequeue(ctx context.Context, queueName string) (string, error) {
	// RPOP: 오른쪽에서 꺼냄 (FIFO)
	return q.client.RPop(ctx, queueName)
}

// 최근 활동 로그 (최근 N개만 유지)
type ActivityLog struct {
	client RedisClient
}

func (a *ActivityLog) AddActivity(ctx context.Context, userID, activity string) error {
	key := fmt.Sprintf("activity:%s", userID)

	// 왼쪽에 추가
	if err := a.client.LPush(ctx, key, activity); err != nil {
		return err
	}

	// 최근 100개만 유지
	return a.client.LTrim(ctx, key, 0, 99)
}

func (a *ActivityLog) GetRecentActivities(ctx context.Context, userID string, count int64) ([]string, error) {
	key := fmt.Sprintf("activity:%s", userID)
	return a.client.LRange(ctx, key, 0, count-1)
}

// ============================================================================
// 4. Set - 고유한 값들의 집합
// ============================================================================

/*
Set 특징:
- 중복 없는 문자열 집합
- 추가/삭제/존재확인 O(1)
- 집합 연산 지원 (교집합, 합집합, 차집합)

활용 사례:
1. 태그 시스템
2. 팔로워/팔로잉
3. 온라인 사용자 추적
4. 추천 시스템 (공통 관심사)
5. IP 블랙리스트

명령어:
- SADD, SREM, SMEMBERS
- SISMEMBER (존재 확인)
- SINTER, SUNION, SDIFF (집합 연산)
- SRANDMEMBER (랜덤 요소)
*/

// 태그 시스템
type TagSystem struct {
	client RedisClient
}

func (t *TagSystem) AddTags(ctx context.Context, postID string, tags []string) error {
	key := fmt.Sprintf("post:tags:%s", postID)
	tagInterfaces := make([]interface{}, len(tags))
	for i, tag := range tags {
		tagInterfaces[i] = tag
	}
	return t.client.SAdd(ctx, key, tagInterfaces...)
}

func (t *TagSystem) GetTags(ctx context.Context, postID string) ([]string, error) {
	key := fmt.Sprintf("post:tags:%s", postID)
	return t.client.SMembers(ctx, key)
}

// 공통 관심사 찾기
func (t *TagSystem) GetCommonInterests(ctx context.Context, userID1, userID2 string) ([]string, error) {
	key1 := fmt.Sprintf("user:interests:%s", userID1)
	key2 := fmt.Sprintf("user:interests:%s", userID2)
	return t.client.SInter(ctx, key1, key2)
}

// 온라인 사용자 추적
type OnlineUsers struct {
	client RedisClient
}

func (o *OnlineUsers) SetOnline(ctx context.Context, userID string) error {
	return o.client.SAdd(ctx, "online_users", userID)
}

func (o *OnlineUsers) IsOnline(ctx context.Context, userID string) (bool, error) {
	return o.client.SIsMember(ctx, "online_users", userID)
}

// ============================================================================
// 5. Sorted Set (ZSet) - 점수 기반 정렬
// ============================================================================

/*
Sorted Set 특징:
- 각 요소에 점수(score) 부여
- 점수 순으로 자동 정렬
- 추가/삭제 O(log N)
- 범위 조회 O(log N + M)

활용 사례:
1. 리더보드 (게임 랭킹)
2. 타임라인 (시간순 정렬)
3. 우선순위 큐
4. 지연 작업 큐
5. Rate Limiting (Sliding Window)

명령어:
- ZADD, ZREM, ZSCORE
- ZRANK, ZREVRANK (순위)
- ZRANGE, ZREVRANGE (범위)
- ZINCRBY (점수 증가)
- ZRANGEBYSCORE (점수 범위)
*/

// 리더보드 (게임 랭킹)
type Leaderboard struct {
	client RedisClient
}

func (l *Leaderboard) UpdateScore(ctx context.Context, gameID, playerID string, score float64) error {
	key := fmt.Sprintf("leaderboard:%s", gameID)
	return l.client.ZAdd(ctx, key, ZMember{Score: score, Member: playerID})
}

func (l *Leaderboard) IncrementScore(ctx context.Context, gameID, playerID string, increment float64) error {
	key := fmt.Sprintf("leaderboard:%s", gameID)
	_, err := l.client.ZIncrBy(ctx, key, increment, playerID)
	return err
}

func (l *Leaderboard) GetTopPlayers(ctx context.Context, gameID string, count int64) ([]ZMember, error) {
	key := fmt.Sprintf("leaderboard:%s", gameID)
	// ZREVRANGE: 높은 점수 순
	return l.client.ZRevRangeWithScores(ctx, key, 0, count-1)
}

func (l *Leaderboard) GetPlayerRank(ctx context.Context, gameID, playerID string) (int64, error) {
	key := fmt.Sprintf("leaderboard:%s", gameID)
	// ZREVRANK: 높은 점수가 0등
	return l.client.ZRank(ctx, key, playerID)
}

// 타임라인 (시간순 정렬)
type Timeline struct {
	client RedisClient
}

func (t *Timeline) AddPost(ctx context.Context, userID, postID string) error {
	key := fmt.Sprintf("timeline:%s", userID)
	// 점수 = 타임스탬프
	return t.client.ZAdd(ctx, key, ZMember{
		Score:  float64(time.Now().UnixMilli()),
		Member: postID,
	})
}

func (t *Timeline) GetRecentPosts(ctx context.Context, userID string, count int64) ([]ZMember, error) {
	key := fmt.Sprintf("timeline:%s", userID)
	// 최신순 (높은 타임스탬프 먼저)
	return t.client.ZRevRangeWithScores(ctx, key, 0, count-1)
}

// 지연 작업 큐 (Delayed Job Queue)
type DelayedQueue struct {
	client RedisClient
}

func (d *DelayedQueue) ScheduleJob(ctx context.Context, queueName, jobID string, executeAt time.Time) error {
	// 실행 시간을 점수로 사용
	return d.client.ZAdd(ctx, queueName, ZMember{
		Score:  float64(executeAt.Unix()),
		Member: jobID,
	})
}

// ============================================================================
// 6. 데이터 타입 선택 가이드
// ============================================================================

/*
데이터 타입 선택 기준:

| 요구사항                  | 추천 타입    | 이유                          |
|-------------------------|------------|------------------------------|
| 단순 키-값 캐시           | String     | 단순함, 직렬화 가능            |
| 객체의 개별 필드 접근      | Hash       | 부분 업데이트 효율적           |
| 순서가 중요한 목록         | List       | 양쪽 끝 O(1) 연산             |
| 고유 항목 집합            | Set        | 중복 자동 제거                |
| 점수/순위 기반 정렬        | Sorted Set | 자동 정렬, 범위 조회 효율적     |
| 시간순 이벤트             | Stream     | 소비자 그룹 지원              |
| 실시간 메시징             | Pub/Sub    | 브로드캐스트                   |

성능 고려:
- String: 큰 값은 메모리 부담
- Hash: 작은 해시는 ziplist로 효율적
- List: 중간 요소 접근은 느림
- Set: 집합 연산은 O(N*M)
- Sorted Set: 범위 조회가 빈번하면 유리
*/

// ============================================================================
// 7. 시니어 면접 답변 요약
// ============================================================================

/*
Q: Redis의 주요 데이터 타입과 활용 사례를 설명해주세요.

A: Redis는 다양한 데이터 구조를 지원하며, 각각 최적의 사용 사례가 있습니다.

1. String:
   - 기본 캐시, 카운터, 분산 락
   - INCR로 atomic 카운터
   - SETNX + TTL로 분산 락

2. Hash:
   - 객체 저장 (사용자 프로필, 세션)
   - 부분 업데이트 가능
   - HGETALL로 전체 조회

3. List:
   - 메시지 큐 (LPUSH + RPOP)
   - 최근 활동 (LPUSH + LTRIM)
   - BLPOP으로 블로킹 큐

4. Set:
   - 태그, 팔로워 목록
   - 온라인 사용자 추적
   - SINTER로 공통 항목

5. Sorted Set:
   - 리더보드 (점수 = 스코어)
   - 타임라인 (점수 = 타임스탬프)
   - ZREVRANGE로 상위 N개

6. 선택 기준:
   - 순서 없이 고유: Set
   - 순서 있음: List 또는 Sorted Set
   - 점수/랭킹: Sorted Set
   - 객체 필드: Hash
   - 단순 값: String
*/
