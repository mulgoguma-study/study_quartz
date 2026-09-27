package hard

/*
================================================================================
문제 1: Distributed Lock - 솔루션
================================================================================
*/

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

/*
================================================================================
분산 락 설계
================================================================================

【 왜 분산 락이 필요한가? 】

단일 서버:
  sync.Mutex로 충분

다중 서버:
  ┌─────────┐     ┌─────────┐     ┌─────────┐
  │ Server1 │     │ Server2 │     │ Server3 │
  └────┬────┘     └────┬────┘     └────┬────┘
       │               │               │
       │   (동시 접근) │               │
       ▼               ▼               ▼
  ┌─────────────────────────────────────────┐
  │              공유 리소스                 │
  │         (DB, 파일, 외부 API)             │
  └─────────────────────────────────────────┘

  → 모든 서버가 공유하는 락이 필요 → Redis


【 Redis 분산 락 원리 】

┌─────────────────────────────────────────────────────────────────────────────┐
│                        SET NX (Set if Not eXists)                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  SET lock:order:123 "server1-uuid" NX PX 30000                             │
│      ─────────────   ─────────────  ── ────────                            │
│          키             값          조건  TTL(ms)                           │
│                                                                             │
│  NX: 키가 없을 때만 설정 (원자적)                                           │
│  PX: 밀리초 단위 TTL                                                        │
│                                                                             │
│  결과:                                                                      │
│  - 성공 → 락 획득                                                          │
│  - 실패 → 이미 다른 프로세스가 락 보유                                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

【 안전한 락 해제 】

문제: 락 소유자만 해제해야 함

1. GET으로 확인
2. 값이 내 것이면 DEL

→ 2단계 사이에 다른 프로세스가 끼어들 수 있음!

해결: Lua Script로 원자적 처리

if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
*/

var (
	ErrLockNotHeld    = errors.New("lock not held")
	ErrLockTimeout    = errors.New("lock acquire timeout")
	ErrLockFailed     = errors.New("failed to acquire lock")
)

// RedisLock Redis 기반 분산 락
type RedisLock struct {
	client *redis.Client
	value  string            // 이 인스턴스의 고유 식별자
	locks  map[string]string // key → value 매핑 (보유 중인 락)
}

// NewRedisLock 생성자
func NewRedisLock(client *redis.Client) *RedisLock {
	return &RedisLock{
		client: client,
		value:  uuid.New().String(), // 인스턴스별 고유 ID
		locks:  make(map[string]string),
	}
}

// TryLock 락 획득 시도 (non-blocking)
/*
【 동작 】
1. SET NX PX 명령 실행
2. 성공하면 로컬에 기록하고 true 반환
3. 실패하면 false 반환

【 주의 】
- 값에는 인스턴스 고유 ID 사용 (나중에 해제 시 검증)
*/
func (l *RedisLock) TryLock(key string, ttl time.Duration) (bool, error) {
	ctx := context.Background()

	// 이 락 시도의 고유 값
	lockValue := l.value + ":" + uuid.New().String()

	// SET key value NX PX milliseconds
	result, err := l.client.SetNX(ctx, key, lockValue, ttl).Result()
	if err != nil {
		return false, err
	}

	if result {
		l.locks[key] = lockValue
		return true, nil
	}

	return false, nil
}

// Lock 락 획득 (timeout까지 대기)
/*
【 동작 】
1. TryLock 시도
2. 실패하면 짧은 간격으로 재시도
3. timeout 초과 시 에러 반환

【 Backoff 전략 】
- 고정 간격: 단순하지만 부하 집중 가능
- 지수 백오프: 재시도 간격 점점 증가
- 지터 추가: 여러 클라이언트가 동시 재시도 방지
*/
func (l *RedisLock) Lock(key string, ttl, timeout time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 재시도 간격 (지터 포함)
	retryInterval := 50 * time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return false, ErrLockTimeout
		default:
		}

		acquired, err := l.TryLock(key, ttl)
		if err != nil {
			return false, err
		}
		if acquired {
			return true, nil
		}

		// 재시도 대기
		time.Sleep(retryInterval)

		// 지수 백오프 (최대 1초)
		if retryInterval < time.Second {
			retryInterval = retryInterval * 2
		}
	}
}

// Unlock 락 해제
/*
【 Lua Script 사용 이유 】
GET과 DEL을 원자적으로 실행해야 함

위험한 시나리오 (Lua 없이):
1. Client A: GET lock → "A"
2. (A의 락이 TTL로 만료됨)
3. Client B: SET NX lock → "B" (성공)
4. Client A: DEL lock (B의 락을 잘못 삭제!)

Lua Script로 해결:
- GET과 DEL이 하나의 원자적 연산
*/
func (l *RedisLock) Unlock(key string) error {
	ctx := context.Background()

	lockValue, exists := l.locks[key]
	if !exists {
		return ErrLockNotHeld
	}

	// Lua Script: 소유자만 해제
	script := redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`)

	result, err := script.Run(ctx, l.client, []string{key}, lockValue).Int()
	if err != nil {
		return err
	}

	delete(l.locks, key)

	if result == 0 {
		return ErrLockNotHeld
	}

	return nil
}

// Extend 락 연장 (워치독 패턴)
/*
【 언제 필요한가? 】
- 작업이 TTL보다 오래 걸릴 수 있음
- TTL이 짧으면 자주 연장 필요
- 워치독 고루틴으로 자동 연장 가능
*/
func (l *RedisLock) Extend(key string, ttl time.Duration) (bool, error) {
	ctx := context.Background()

	lockValue, exists := l.locks[key]
	if !exists {
		return false, ErrLockNotHeld
	}

	// Lua Script: 소유자인 경우에만 TTL 연장
	script := redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("pexpire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`)

	result, err := script.Run(ctx, l.client, []string{key}, lockValue, ttl.Milliseconds()).Int()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

/*
================================================================================
고급 패턴
================================================================================

【 1. 워치독 (Watchdog) 패턴 】

func (l *RedisLock) LockWithWatchdog(key string, ttl time.Duration) (context.CancelFunc, error) {
    acquired, err := l.Lock(key, ttl, 10*time.Second)
    if err != nil || !acquired {
        return nil, err
    }

    ctx, cancel := context.WithCancel(context.Background())

    // 워치독 고루틴: TTL의 1/3마다 연장
    go func() {
        ticker := time.NewTicker(ttl / 3)
        defer ticker.Stop()

        for {
            select {
            case <-ctx.Done():
                l.Unlock(key)
                return
            case <-ticker.C:
                l.Extend(key, ttl)
            }
        }
    }()

    return cancel, nil
}

// 사용
cancel, _ := lock.LockWithWatchdog("mykey", 30*time.Second)
defer cancel()  // 작업 완료 후 락 해제

【 2. Redlock 알고리즘 (다중 Redis) 】

여러 Redis 인스턴스에서 과반수 락 획득 필요:
1. 현재 시간 기록
2. 모든 인스턴스에 순차적으로 락 시도
3. 과반수 이상 획득 && 경과 시간 < TTL → 성공
4. 실패 시 모든 인스턴스에서 해제

【 3. 공정한 락 (Fair Lock) 】

요청 순서대로 락 획득:
- Redis List를 큐로 사용
- BRPOPLPUSH로 대기

================================================================================
주의사항
================================================================================

1. 시계 동기화
   - 서버 간 시간 차이가 TTL에 영향
   - NTP 동기화 필수

2. 네트워크 파티션
   - 락 보유 중 네트워크 분리 시 위험
   - 펜싱 토큰 사용 고려

3. TTL 설정
   - 너무 짧으면: 작업 완료 전 만료
   - 너무 길면: 장애 시 오래 대기
   - 워치독으로 동적 연장 권장

4. 재시작 후 복구
   - 프로세스 재시작 시 보유 락 정보 손실
   - 복구 전략 필요
================================================================================
*/

func main() {
	println("Distributed Lock 솔루션")
	println("실제 사용 예시:")
	println(`
client := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

lock := NewRedisLock(client)

// 락 획득
acquired, err := lock.Lock("order:123", 30*time.Second, 10*time.Second)
if acquired {
    defer lock.Unlock("order:123")
    // 임계 영역 작업
}
`)
}
