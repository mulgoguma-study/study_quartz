package medium

/*
================================================================================
문제 1: LRU Cache - 솔루션
================================================================================
*/

// ============================================================================
// 핵심 자료구조: HashMap + Doubly Linked List
// ============================================================================
/*
【 왜 이 조합인가? 】

┌─────────────────────────────────────────────────────────────────────────────┐
│                        자료구조별 시간 복잡도                                │
├─────────────────────┬─────────────┬─────────────┬───────────────────────────┤
│      연산           │  HashMap    │ LinkedList  │  HashMap + LinkedList     │
├─────────────────────┼─────────────┼─────────────┼───────────────────────────┤
│ 검색 (by key)       │    O(1)     │    O(n)     │         O(1)              │
│ 삽입                │    O(1)     │    O(1)     │         O(1)              │
│ 삭제 (by key)       │    O(1)     │    O(n)     │         O(1) ★            │
│ 순서 유지           │     X       │     O       │          O                │
│ 앞/뒤 이동          │     X       │    O(1)     │         O(1)              │
└─────────────────────┴─────────────┴─────────────┴───────────────────────────┘

★ HashMap이 노드 포인터를 저장하므로 O(n) 검색 없이 즉시 노드 접근 가능

【 구조 】

HashMap: key → Node*
         ┌───────────────────────────────────────┐
         │  1 → Node{1, 100}                     │
         │  2 → Node{2, 200}                     │
         │  3 → Node{3, 300}                     │
         └───────────────────────────────────────┘
                    │
                    ▼
Doubly Linked List (최근 → 오래된 순서):
         ┌──────┐    ┌──────┐    ┌──────┐    ┌──────┐    ┌──────┐
         │ HEAD │ ←→ │  3   │ ←→ │  1   │ ←→ │  2   │ ←→ │ TAIL │
         │ dummy│    │ 300  │    │ 100  │    │ 200  │    │ dummy│
         └──────┘    └──────┘    └──────┘    └──────┘    └──────┘
           ↑                                               ↑
        최근 사용                                        LRU (제거 대상)
*/

// Node 이중 연결 리스트의 노드
type Node struct {
	key   int
	value int
	prev  *Node
	next  *Node
}

// LRUCache LRU 캐시 구현체
type LRUCache struct {
	capacity int
	cache    map[int]*Node // key → Node 매핑
	head     *Node         // 더미 헤드 (최근 사용)
	tail     *Node         // 더미 테일 (LRU)
}

// NewLRUCache 생성자
/*
【 Dummy Head/Tail 사용 이유 】
- 삽입/삭제 시 edge case 처리 단순화
- 빈 리스트, 단일 노드 등 특수 케이스 별도 처리 불필요
- 코드가 더 깔끔하고 버그 가능성 감소
*/
func NewLRUCache(capacity int) *LRUCache {
	// 더미 노드 생성
	head := &Node{}
	tail := &Node{}

	// 더미 노드 연결
	head.next = tail
	tail.prev = head

	return &LRUCache{
		capacity: capacity,
		cache:    make(map[int]*Node),
		head:     head,
		tail:     tail,
	}
}

// ============================================================================
// 내부 헬퍼 메서드
// ============================================================================

// addToFront 노드를 head 바로 다음에 추가 (최근 사용 표시)
/*
【 동작 】
Before: HEAD ←→ A ←→ B ←→ TAIL
After:  HEAD ←→ NEW ←→ A ←→ B ←→ TAIL
*/
func (c *LRUCache) addToFront(node *Node) {
	node.prev = c.head
	node.next = c.head.next

	c.head.next.prev = node
	c.head.next = node
}

// removeNode 노드를 리스트에서 제거
/*
【 동작 】
Before: A ←→ NODE ←→ B
After:  A ←→ B
*/
func (c *LRUCache) removeNode(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// moveToFront 기존 노드를 맨 앞으로 이동 (최근 사용 갱신)
func (c *LRUCache) moveToFront(node *Node) {
	c.removeNode(node) // 현재 위치에서 제거
	c.addToFront(node) // 앞에 추가
}

// removeLRU 가장 오래된 노드 제거 (tail 바로 앞)
func (c *LRUCache) removeLRU() *Node {
	lru := c.tail.prev // 실제 LRU 노드 (더미 아님)
	c.removeNode(lru)
	return lru
}

// ============================================================================
// 공개 메서드
// ============================================================================

// Get key로 값 조회
/*
【 동작 순서 】
1. HashMap에서 O(1)로 노드 찾기
2. 없으면 -1 반환
3. 있으면 해당 노드를 맨 앞으로 이동 (최근 사용 표시)
4. 값 반환
*/
func (c *LRUCache) Get(key int) int {
	// HashMap에서 노드 조회 - O(1)
	node, exists := c.cache[key]
	if !exists {
		return -1
	}

	// 최근 사용되었으므로 맨 앞으로 이동 - O(1)
	c.moveToFront(node)

	return node.value
}

// Put key-value 저장
/*
【 동작 순서 】
1. key가 이미 존재하면:
   - value 업데이트
   - 노드를 맨 앞으로 이동
2. key가 새로운 경우:
   - 용량 초과 시 LRU 노드 제거
   - 새 노드 생성하여 맨 앞에 추가
   - HashMap에 등록
*/
func (c *LRUCache) Put(key int, value int) {
	// 이미 존재하는 key인 경우
	if node, exists := c.cache[key]; exists {
		node.value = value    // 값 업데이트
		c.moveToFront(node)   // 최근 사용 표시
		return
	}

	// 용량 초과 체크
	if len(c.cache) >= c.capacity {
		// LRU 노드 제거
		lru := c.removeLRU()
		delete(c.cache, lru.key) // HashMap에서도 제거
	}

	// 새 노드 생성
	newNode := &Node{
		key:   key,
		value: value,
	}

	// 맨 앞에 추가
	c.addToFront(newNode)

	// HashMap에 등록
	c.cache[key] = newNode
}

/*
================================================================================
시간/공간 복잡도 분석
================================================================================

【 시간 복잡도 】
- Get: O(1)
  - HashMap 조회: O(1)
  - 노드 이동: O(1) (포인터 조작만)

- Put: O(1)
  - HashMap 조회/삽입: O(1)
  - 노드 삭제/추가: O(1)

【 공간 복잡도 】
- O(capacity)
  - HashMap: capacity개의 엔트리
  - LinkedList: capacity개의 노드

================================================================================
실무 적용 및 확장
================================================================================

【 실무에서의 LRU 캐시 】

1. Redis의 maxmemory-policy
   - volatile-lru: 만료 설정된 키 중 LRU 제거
   - allkeys-lru: 모든 키 중 LRU 제거

2. 브라우저 캐시
   - HTTP 캐시에서 오래된 리소스 제거

3. 데이터베이스 버퍼 풀
   - MySQL InnoDB의 버퍼 풀

【 확장 가능한 기능 】

1. TTL 지원
   type Node struct {
       key       int
       value     int
       expiresAt time.Time  // 만료 시간
       // ...
   }

2. Thread-safe 버전
   type LRUCache struct {
       mu sync.RWMutex  // 읽기/쓰기 락
       // ...
   }

3. 샤딩 (Sharding)
   - 여러 LRU 캐시로 분산
   - 락 경합 감소

4. LRU-K, 2Q, ARC 등 변형 알고리즘
   - 단순 LRU의 한계 극복
================================================================================
*/

// 테스트
func main() {
	cache := NewLRUCache(2)

	cache.Put(1, 1)
	println("Put(1,1)")

	cache.Put(2, 2)
	println("Put(2,2)")

	println("Get(1):", cache.Get(1)) // 1

	cache.Put(3, 3)                  // 2 제거됨
	println("Put(3,3) - evicts 2")

	println("Get(2):", cache.Get(2)) // -1

	cache.Put(4, 4)                  // 1 제거됨
	println("Put(4,4) - evicts 1")

	println("Get(1):", cache.Get(1)) // -1
	println("Get(3):", cache.Get(3)) // 3
	println("Get(4):", cache.Get(4)) // 4
}
