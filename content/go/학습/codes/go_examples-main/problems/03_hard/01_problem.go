package hard

/*
================================================================================
문제 1: Distributed Lock (분산 락) 구현
================================================================================

난이도: Hard
주제: 분산 시스템, 동시성, Redis

【 문제 설명 】
Redis를 사용한 분산 락을 구현하세요.

요구사항:
- TryLock(key, ttl): 락 획득 시도 (non-blocking)
- Lock(key, ttl, timeout): 락 획득 (timeout까지 대기)
- Unlock(key): 락 해제
- 자동 만료 (TTL)
- 락 소유자만 해제 가능

【 인터페이스 】
type DistributedLock interface {
    TryLock(key string, ttl time.Duration) (bool, error)
    Lock(key string, ttl time.Duration, timeout time.Duration) (bool, error)
    Unlock(key string) error
}

【 Redis 명령어 힌트 】
SET key value NX PX milliseconds  // 키가 없을 때만 설정
DEL key                            // 키 삭제
GET key                            // 값 조회

【 Lua Script (원자적 연산) 】
// 락 해제 - 소유자만 해제 가능
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end

【 제약 조건 】
- 여러 프로세스/서버에서 동시 접근
- 네트워크 장애 고려
- 데드락 방지

【 실무 연관성 】
- 분산 시스템의 임계 영역 보호
- 중복 작업 방지
- 리더 선출
================================================================================
*/

import "time"

// DistributedLock 분산 락 인터페이스
type DistributedLock interface {
	TryLock(key string, ttl time.Duration) (bool, error)
	Lock(key string, ttl time.Duration, timeout time.Duration) (bool, error)
	Unlock(key string) error
	Extend(key string, ttl time.Duration) (bool, error) // 락 연장
}

// RedisLock Redis 기반 분산 락 구현
type RedisLock struct {
	// 여기에 필드를 정의하세요
	// client: Redis 클라이언트
	// value: 락 소유자 식별 (UUID 등)
}

// NewRedisLock 생성자
func NewRedisLock(client interface{}) *RedisLock {
	// 여기에 코드를 작성하세요
	return nil
}

// TryLock 락 획득 시도 (non-blocking)
func (l *RedisLock) TryLock(key string, ttl time.Duration) (bool, error) {
	// 여기에 코드를 작성하세요
	return false, nil
}

// Lock 락 획득 (timeout까지 대기)
func (l *RedisLock) Lock(key string, ttl, timeout time.Duration) (bool, error) {
	// 여기에 코드를 작성하세요
	return false, nil
}

// Unlock 락 해제
func (l *RedisLock) Unlock(key string) error {
	// 여기에 코드를 작성하세요
	return nil
}

// Extend 락 연장
func (l *RedisLock) Extend(key string, ttl time.Duration) (bool, error) {
	// 여기에 코드를 작성하세요
	return false, nil
}

// 테스트 (실제 Redis 필요)
func main() {
	println("Distributed Lock 구현 문제")
	println("실제 테스트는 Redis 연결 필요")
}
