package medium

/*
================================================================================
문제 2: Rate Limiter (처리율 제한기) 구현
================================================================================

난이도: Medium
주제: 동시성, 시스템 설계, 슬라이딩 윈도우

【 문제 설명 】
API 요청을 제한하는 Rate Limiter를 구현하세요.

다음 두 가지 알고리즘 중 하나 이상을 구현하세요:
1. Token Bucket (토큰 버킷)
2. Sliding Window Log (슬라이딩 윈도우 로그)

【 요구사항 】
- Allow(userID string) bool: 요청 허용 여부 반환
- 초당 N개의 요청만 허용
- Thread-safe 해야 함
- 여러 사용자를 독립적으로 처리

【 인터페이스 】
type RateLimiter interface {
    Allow(userID string) bool
}

【 예시 - Token Bucket (초당 5개 허용) 】
limiter := NewTokenBucketLimiter(5, 5)  // capacity=5, refillRate=5/sec

limiter.Allow("user1")  // true (4 tokens left)
limiter.Allow("user1")  // true (3 tokens left)
limiter.Allow("user1")  // true (2 tokens left)
limiter.Allow("user1")  // true (1 token left)
limiter.Allow("user1")  // true (0 tokens left)
limiter.Allow("user1")  // false (no tokens)

// 1초 후...
limiter.Allow("user1")  // true (tokens refilled)

【 예시 - Sliding Window (1분에 100개 허용) 】
limiter := NewSlidingWindowLimiter(100, time.Minute)

for i := 0; i < 100; i++ {
    limiter.Allow("user1")  // true
}
limiter.Allow("user1")  // false (limit reached)

【 제약 조건 】
- 동시에 여러 고루틴에서 호출될 수 있음
- 메모리 효율적으로 구현
- 사용자 수는 최대 10^6

【 힌트 】
Token Bucket:
- 토큰이 일정 속도로 채워짐
- 요청 시 토큰 1개 소비
- 버스트 허용 가능

Sliding Window Log:
- 각 요청의 타임스탬프 저장
- 윈도우 밖의 오래된 요청 제거
- 정확한 제한 가능

【 실무 연관성 】
- API Gateway의 Rate Limiting
- DDoS 방어
- 서비스 보호 및 공정한 사용
================================================================================
*/

import (
	"time"
)

// RateLimiter 인터페이스
type RateLimiter interface {
	Allow(userID string) bool
}

// ============================================================================
// Token Bucket 구현
// ============================================================================

type TokenBucketLimiter struct {
	// 여기에 필드를 정의하세요
	// capacity: 버킷 최대 용량
	// refillRate: 초당 토큰 충전 개수
}

func NewTokenBucketLimiter(capacity int, refillRate float64) *TokenBucketLimiter {
	// 여기에 코드를 작성하세요
	return nil
}

func (l *TokenBucketLimiter) Allow(userID string) bool {
	// 여기에 코드를 작성하세요
	return false
}

// ============================================================================
// Sliding Window Log 구현
// ============================================================================

type SlidingWindowLimiter struct {
	// 여기에 필드를 정의하세요
	// limit: 윈도우 내 최대 요청 수
	// window: 윈도우 크기 (예: 1분)
}

func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
	// 여기에 코드를 작성하세요
	return nil
}

func (l *SlidingWindowLimiter) Allow(userID string) bool {
	// 여기에 코드를 작성하세요
	return false
}

// 테스트
func main() {
	// Token Bucket 테스트
	tb := NewTokenBucketLimiter(5, 5.0)
	for i := 0; i < 7; i++ {
		println("Request", i+1, ":", tb.Allow("user1"))
	}

	// Sliding Window 테스트
	sw := NewSlidingWindowLimiter(3, time.Second)
	for i := 0; i < 5; i++ {
		println("Request", i+1, ":", sw.Allow("user1"))
	}
}
