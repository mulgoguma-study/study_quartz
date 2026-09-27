package medium

/*
================================================================================
문제 3: Concurrent Safe Map with Sharding
================================================================================

난이도: Medium
주제: 동시성, 샤딩, 락 최적화

【 문제 설명 】
Go의 sync.Map보다 더 성능이 좋은 Concurrent Map을 샤딩 기법으로 구현하세요.

sync.Map의 한계:
- 읽기 위주 워크로드에 최적화
- 쓰기가 많으면 성능 저하

샤딩 맵의 장점:
- 키를 여러 샤드에 분산
- 샤드별로 독립적인 락 → 락 경합 감소

【 요구사항 】
- Get(key string) (value interface{}, ok bool)
- Set(key string, value interface{})
- Delete(key string)
- Len() int
- Thread-safe 해야 함

【 인터페이스 】
type ShardedMap interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{})
    Delete(key string)
    Len() int
}

【 예시 】
m := NewShardedMap(16)  // 16개 샤드

m.Set("user:1", User{Name: "Alice"})
m.Set("user:2", User{Name: "Bob"})

user, ok := m.Get("user:1")  // User{Name: "Alice"}, true
m.Delete("user:1")
user, ok = m.Get("user:1")   // nil, false

【 제약 조건 】
- 샤드 수는 2의 거듭제곱 (최적화를 위해)
- 해시 함수로 키를 샤드에 균등 분배
- 각 샤드는 독립적인 RWMutex 사용

【 힌트 】
1. FNV 해시 함수 사용 권장
2. 샤드 인덱스 = hash(key) % shardCount
3. RWMutex로 읽기/쓰기 분리

【 실무 연관성 】
- 고성능 캐시 구현
- 세션 저장소
- 동시성 높은 인메모리 DB
================================================================================
*/

// ShardedMap 구조체를 정의하고 구현하세요
type ShardedMap struct {
	// 여기에 필드를 정의하세요
	// shards: 샤드 배열
	// shardCount: 샤드 개수
}

// NewShardedMap 생성자
func NewShardedMap(shardCount int) *ShardedMap {
	// 여기에 코드를 작성하세요
	return nil
}

// Get 값 조회
func (m *ShardedMap) Get(key string) (interface{}, bool) {
	// 여기에 코드를 작성하세요
	return nil, false
}

// Set 값 저장
func (m *ShardedMap) Set(key string, value interface{}) {
	// 여기에 코드를 작성하세요
}

// Delete 값 삭제
func (m *ShardedMap) Delete(key string) {
	// 여기에 코드를 작성하세요
}

// Len 전체 항목 수
func (m *ShardedMap) Len() int {
	// 여기에 코드를 작성하세요
	return 0
}

// 테스트
func main() {
	m := NewShardedMap(16)

	// 기본 연산 테스트
	m.Set("key1", "value1")
	m.Set("key2", "value2")

	if v, ok := m.Get("key1"); ok {
		println("key1:", v.(string))
	}

	println("Len:", m.Len())

	m.Delete("key1")
	if _, ok := m.Get("key1"); !ok {
		println("key1 deleted")
	}
}
