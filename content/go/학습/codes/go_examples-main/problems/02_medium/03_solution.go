package medium

/*
================================================================================
문제 3: Concurrent Safe Map with Sharding - 솔루션
================================================================================
*/

import (
	"hash/fnv"
	"sync"
)

/*
================================================================================
왜 샤딩인가?
================================================================================

【 단일 락의 문제 】

┌─────────────────────────────────────────────────────────────────────────────┐
│                        단일 락 vs 샤딩 비교                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   단일 락:                                                                  │
│   ┌───────────────────────────────────────┐                                 │
│   │              전체 Map                  │                                 │
│   │         [하나의 Mutex 락]             │                                 │
│   └───────────────────────────────────────┘                                 │
│         │         │         │         │                                     │
│        G1        G2        G3        G4   (고루틴들이 순차 대기)             │
│   ──────────────────────────────────────▶                                   │
│                시간                                                         │
│                                                                             │
│   샤딩:                                                                     │
│   ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐                              │
│   │ Shard0 │ │ Shard1 │ │ Shard2 │ │ Shard3 │                              │
│   │ [Lock] │ │ [Lock] │ │ [Lock] │ │ [Lock] │                              │
│   └────────┘ └────────┘ └────────┘ └────────┘                              │
│       │          │          │          │                                    │
│      G1         G2         G3         G4   (병렬 처리!)                      │
│   ══════════════════════════════════════▶                                   │
│                시간 (훨씬 짧음)                                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

【 성능 비교 (벤치마크 예시) 】

goroutines | sync.Mutex | sync.RWMutex | sync.Map | ShardedMap(16)
-----------+------------+--------------+----------+---------------
    1      |   100ns    |    80ns      |   50ns   |    60ns
   10      |   500ns    |   200ns      |  100ns   |    70ns
  100      |  2000ns    |   800ns      |  150ns   |    80ns
 1000      |  5000ns    |  2000ns      |  300ns   |   100ns

→ 고루틴이 많아질수록 샤딩의 이점이 커짐
*/

// ============================================================================
// 개별 샤드 구조체
// ============================================================================

// shard 단일 샤드 (독립적인 락과 맵)
type shard struct {
	mu   sync.RWMutex           // 읽기/쓰기 분리 락
	data map[string]interface{} // 실제 데이터
}

// ============================================================================
// 샤딩 맵 구현
// ============================================================================

// ShardedMap 샤딩된 동시성 맵
type ShardedMap struct {
	shards     []*shard // 샤드 배열
	shardCount uint32   // 샤드 개수 (2의 거듭제곱)
	shardMask  uint32   // 비트마스크 (shardCount - 1)
}

// NewShardedMap 생성자
/*
【 샤드 개수 선택 가이드 】

CPU 코어 수의 2-4배 정도가 적당:
- 너무 적으면: 락 경합 여전히 발생
- 너무 많으면: 메모리 오버헤드, Len() 등 전체 순회 느려짐

일반적인 권장:
- 16 샤드: 대부분의 경우 충분
- 32 샤드: 매우 높은 동시성
- 256 샤드: 극한의 성능 필요 시
*/
func NewShardedMap(shardCount int) *ShardedMap {
	// 2의 거듭제곱으로 올림 (비트 연산 최적화)
	shardCount = roundUpToPowerOf2(shardCount)

	shards := make([]*shard, shardCount)
	for i := 0; i < shardCount; i++ {
		shards[i] = &shard{
			data: make(map[string]interface{}),
		}
	}

	return &ShardedMap{
		shards:     shards,
		shardCount: uint32(shardCount),
		shardMask:  uint32(shardCount - 1), // 비트마스크로 모듈로 연산 대체
	}
}

// roundUpToPowerOf2 2의 거듭제곱으로 올림
func roundUpToPowerOf2(n int) int {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}

// getShard 키에 해당하는 샤드 반환
/*
【 해시 함수 선택 】

FNV-1a 해시를 사용하는 이유:
1. 빠름: 단순한 비트 연산
2. 균등 분포: 키가 샤드에 고르게 분산
3. 충돌 적음: 해시 충돌이 적어 성능 안정적

비트마스크 최적화:
- index = hash % shardCount (나눗셈 - 느림)
- index = hash & shardMask  (비트 AND - 빠름)
- shardCount가 2의 거듭제곱이면 동일한 결과
*/
func (m *ShardedMap) getShard(key string) *shard {
	// FNV-1a 해시 계산
	h := fnv.New32a()
	h.Write([]byte(key))
	hash := h.Sum32()

	// 비트마스크로 샤드 인덱스 계산 (모듈로 연산보다 빠름)
	return m.shards[hash&m.shardMask]
}

// Get 값 조회
/*
【 RLock 사용 이유 】
- 읽기 전용 연산이므로 RLock 사용
- 여러 고루틴이 동시에 같은 샤드 읽기 가능
- 쓰기 중에만 대기
*/
func (m *ShardedMap) Get(key string) (interface{}, bool) {
	shard := m.getShard(key)

	shard.mu.RLock()
	defer shard.mu.RUnlock()

	val, ok := shard.data[key]
	return val, ok
}

// Set 값 저장
/*
【 Lock 사용 이유 】
- 쓰기 연산이므로 배타적 Lock 필요
- 해당 샤드의 다른 모든 연산 대기
- 다른 샤드는 영향 없음!
*/
func (m *ShardedMap) Set(key string, value interface{}) {
	shard := m.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.data[key] = value
}

// Delete 값 삭제
func (m *ShardedMap) Delete(key string) {
	shard := m.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	delete(shard.data, key)
}

// Len 전체 항목 수
/*
【 주의 】
- 모든 샤드를 순회해야 함
- 순간적인 스냅샷이 아님 (순회 중 변경 가능)
- 정확한 값이 필요하면 모든 샤드 락 필요 (비효율적)
*/
func (m *ShardedMap) Len() int {
	count := 0
	for _, shard := range m.shards {
		shard.mu.RLock()
		count += len(shard.data)
		shard.mu.RUnlock()
	}
	return count
}

// ============================================================================
// 추가 유틸리티 메서드
// ============================================================================

// Keys 모든 키 반환
func (m *ShardedMap) Keys() []string {
	keys := make([]string, 0)
	for _, shard := range m.shards {
		shard.mu.RLock()
		for k := range shard.data {
			keys = append(keys, k)
		}
		shard.mu.RUnlock()
	}
	return keys
}

// ForEach 모든 항목에 함수 적용
func (m *ShardedMap) ForEach(fn func(key string, value interface{})) {
	for _, shard := range m.shards {
		shard.mu.RLock()
		for k, v := range shard.data {
			fn(k, v)
		}
		shard.mu.RUnlock()
	}
}

// GetOrSet 있으면 반환, 없으면 설정 후 반환
/*
【 원자적 연산 】
- Get과 Set을 분리하면 race condition 발생 가능
- 하나의 락 안에서 처리
*/
func (m *ShardedMap) GetOrSet(key string, value interface{}) (interface{}, bool) {
	shard := m.getShard(key)

	// 먼저 읽기 시도 (락 경합 최소화)
	shard.mu.RLock()
	if val, ok := shard.data[key]; ok {
		shard.mu.RUnlock()
		return val, true
	}
	shard.mu.RUnlock()

	// 없으면 쓰기 락으로 업그레이드
	shard.mu.Lock()
	defer shard.mu.Unlock()

	// Double-check (다른 고루틴이 먼저 설정했을 수 있음)
	if val, ok := shard.data[key]; ok {
		return val, true
	}

	shard.data[key] = value
	return value, false
}

/*
================================================================================
성능 최적화 추가 기법
================================================================================

【 1. 캐시 라인 패딩 】

CPU 캐시 라인(보통 64바이트) 공유로 인한 false sharing 방지:

type shard struct {
    mu   sync.RWMutex
    data map[string]interface{}
    _    [64 - unsafe.Sizeof(sync.RWMutex{}) - 8]byte  // 패딩
}

【 2. Lock-free 읽기 (Read-Copy-Update) 】

읽기가 압도적으로 많은 경우:
- atomic.Value로 맵 전체를 저장
- 쓰기 시 복사 후 교체

【 3. 로컬 캐시 + 샤딩 】

type ShardedMapWithCache struct {
    sharded *ShardedMap
    local   sync.Map  // 핫 데이터 캐시
}

【 4. 배치 연산 】

func (m *ShardedMap) SetBatch(items map[string]interface{}) {
    // 샤드별로 그룹핑하여 한 번에 처리
}

================================================================================
sync.Map vs ShardedMap 선택 기준
================================================================================

sync.Map 선택:
- 읽기 >> 쓰기
- 키가 한 번 쓰고 여러 번 읽는 패턴
- 고루틴별로 다른 키 세트 접근

ShardedMap 선택:
- 읽기/쓰기가 혼합
- 많은 고루틴이 동시 접근
- 예측 가능한 성능 필요
================================================================================
*/

// 테스트
func main() {
	m := NewShardedMap(16)

	// 기본 연산 테스트
	m.Set("user:1", map[string]string{"name": "Alice"})
	m.Set("user:2", map[string]string{"name": "Bob"})
	m.Set("user:3", map[string]string{"name": "Charlie"})

	if v, ok := m.Get("user:1"); ok {
		user := v.(map[string]string)
		println("user:1 name:", user["name"])
	}

	println("Total items:", m.Len())

	m.Delete("user:1")
	println("After delete:", m.Len())

	// GetOrSet 테스트
	val, existed := m.GetOrSet("user:4", map[string]string{"name": "David"})
	user := val.(map[string]string)
	println("user:4 name:", user["name"], "existed:", existed)
}
