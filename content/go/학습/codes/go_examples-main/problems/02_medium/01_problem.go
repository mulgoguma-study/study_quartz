package medium

/*
================================================================================
문제 1: LRU Cache 구현
================================================================================

난이도: Medium
주제: 해시맵, 이중 연결 리스트, 캐시 설계

【 문제 설명 】
LRU(Least Recently Used) 캐시를 구현하세요.

LRU 캐시는 다음 연산을 지원합니다:
- Get(key): key가 존재하면 value 반환, 없으면 -1 반환
- Put(key, value): key-value 저장. 용량 초과 시 가장 오래된 항목 제거

모든 연산은 O(1) 시간 복잡도로 동작해야 합니다.

【 인터페이스 】
type LRUCache interface {
    Get(key int) int
    Put(key int, value int)
}

【 예시 】
cache := NewLRUCache(2)  // 용량 2
cache.Put(1, 1)          // cache = {1=1}
cache.Put(2, 2)          // cache = {1=1, 2=2}
cache.Get(1)             // return 1, cache = {2=2, 1=1} (1이 최근 사용됨)
cache.Put(3, 3)          // cache = {1=1, 3=3} (2가 제거됨, LRU)
cache.Get(2)             // return -1 (2는 제거됨)
cache.Put(4, 4)          // cache = {3=3, 4=4} (1이 제거됨)
cache.Get(1)             // return -1
cache.Get(3)             // return 3
cache.Get(4)             // return 4

【 제약 조건 】
- 1 <= capacity <= 3000
- 0 <= key <= 10^4
- 0 <= value <= 10^5
- Get, Put 호출 횟수 <= 2 * 10^5

【 힌트 】
- HashMap + Doubly Linked List 조합
- 최근 사용된 항목은 리스트 앞으로 이동
- 제거 시 리스트 뒤에서 제거

【 실무 연관성 】
- Redis, Memcached의 기본 동작 원리
- 브라우저 캐시
- 데이터베이스 버퍼 풀
================================================================================
*/

// LRUCache 구조체를 정의하고 구현하세요
type LRUCache struct {
	// 여기에 필드를 정의하세요
}

// NewLRUCache 생성자
func NewLRUCache(capacity int) *LRUCache {
	// 여기에 코드를 작성하세요
	return nil
}

// Get 메서드
func (c *LRUCache) Get(key int) int {
	// 여기에 코드를 작성하세요
	return -1
}

// Put 메서드
func (c *LRUCache) Put(key int, value int) {
	// 여기에 코드를 작성하세요
}

// 테스트
func main() {
	cache := NewLRUCache(2)
	cache.Put(1, 1)
	cache.Put(2, 2)
	println("Get(1):", cache.Get(1)) // 1
	cache.Put(3, 3)
	println("Get(2):", cache.Get(2)) // -1
	cache.Put(4, 4)
	println("Get(1):", cache.Get(1)) // -1
	println("Get(3):", cache.Get(3)) // 3
	println("Get(4):", cache.Get(4)) // 4
}
