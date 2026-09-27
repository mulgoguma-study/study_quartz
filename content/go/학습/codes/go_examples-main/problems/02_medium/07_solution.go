package medium

/*
================================================================================
문제 7: Top K Frequent Elements - 솔루션
================================================================================
*/

import (
	"container/heap"
)

/*
================================================================================
세 가지 접근법 비교
================================================================================

┌─────────────────────────────────────────────────────────────────────────────┐
│                        알고리즘 비교                                        │
├──────────────────┬─────────────────┬─────────────────┬──────────────────────┤
│      방법        │   시간 복잡도   │   공간 복잡도   │        특징          │
├──────────────────┼─────────────────┼─────────────────┼──────────────────────┤
│ 정렬             │   O(n log n)    │     O(n)        │ 간단, 범용적         │
│ 힙 (Min-Heap)    │   O(n log k)    │     O(k)        │ k가 작을 때 효율적   │
│ 버킷 정렬        │     O(n)        │     O(n)        │ 빈도 범위 제한 시    │
│ Quick Select     │   O(n) 평균     │     O(n)        │ 불안정, 최악 O(n²)   │
└──────────────────┴─────────────────┴─────────────────┴──────────────────────┘
*/

// ============================================================================
// 방법 1: 버킷 정렬 - O(n) ★ 최적
// ============================================================================

/*
【 핵심 아이디어 】
- 빈도가 가능한 범위: 1 ~ n (배열 길이)
- 버킷[i] = 빈도가 i인 원소들
- 뒤에서부터 k개 수집

┌─────────────────────────────────────────────────────────────────────────────┐
│  예시: nums = [1,1,1,2,2,3], k = 2                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  빈도 맵: {1: 3, 2: 2, 3: 1}                                               │
│                                                                             │
│  버킷 (인덱스 = 빈도):                                                     │
│  ┌─────┬─────┬─────┬─────┬─────┬─────┬─────┐                              │
│  │  0  │  1  │  2  │  3  │  4  │  5  │  6  │                              │
│  ├─────┼─────┼─────┼─────┼─────┼─────┼─────┤                              │
│  │ [ ] │ [3] │ [2] │ [1] │ [ ] │ [ ] │ [ ] │                              │
│  └─────┴─────┴─────┴─────┴─────┴─────┴─────┘                              │
│                        ↑     ↑                                              │
│                     2번출현  3번출현                                        │
│                                                                             │
│  뒤에서부터 k개 수집: [1, 2]                                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
*/

func TopKFrequent(nums []int, k int) []int {
	// 1. 빈도 계산
	freq := make(map[int]int)
	for _, num := range nums {
		freq[num]++
	}

	// 2. 버킷 생성 (인덱스 = 빈도)
	n := len(nums)
	buckets := make([][]int, n+1)
	for i := range buckets {
		buckets[i] = []int{}
	}

	// 3. 버킷에 분배
	for num, count := range freq {
		buckets[count] = append(buckets[count], num)
	}

	// 4. 뒤에서부터 k개 수집
	result := make([]int, 0, k)
	for i := n; i >= 0 && len(result) < k; i-- {
		result = append(result, buckets[i]...)
	}

	return result[:k]
}

// ============================================================================
// 방법 2: Min-Heap - O(n log k)
// ============================================================================

/*
【 핵심 아이디어 】
- 크기 k의 Min-Heap 유지
- 빈도가 더 큰 원소가 오면 최소값 교체
- 힙에 남는 k개가 답

【 왜 Min-Heap인가? 】
- Max-Heap이면 최대값을 빠르게 제거 불가
- Min-Heap이면 최소값(가장 빈도 낮은 것) 제거 가능
*/

// freqItem 힙에 저장할 아이템
type freqItem struct {
	num   int
	count int
}

// minHeap Min-Heap 구현
type minHeap []freqItem

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].count < h[j].count }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(freqItem)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func TopKFrequentHeap(nums []int, k int) []int {
	// 1. 빈도 계산
	freq := make(map[int]int)
	for _, num := range nums {
		freq[num]++
	}

	// 2. Min-Heap으로 상위 k개 유지
	h := &minHeap{}
	heap.Init(h)

	for num, count := range freq {
		heap.Push(h, freqItem{num, count})
		if h.Len() > k {
			heap.Pop(h) // 최소값 제거
		}
	}

	// 3. 결과 추출
	result := make([]int, k)
	for i := k - 1; i >= 0; i-- {
		result[i] = heap.Pop(h).(freqItem).num
	}

	return result
}

// ============================================================================
// 방법 3: 정렬 - O(n log n)
// ============================================================================

func TopKFrequentSort(nums []int, k int) []int {
	import "sort"

	// 1. 빈도 계산
	freq := make(map[int]int)
	for _, num := range nums {
		freq[num]++
	}

	// 2. 고유 원소 추출
	unique := make([]int, 0, len(freq))
	for num := range freq {
		unique = append(unique, num)
	}

	// 3. 빈도 기준 정렬
	sort.Slice(unique, func(i, j int) bool {
		return freq[unique[i]] > freq[unique[j]]
	})

	// 4. 상위 k개 반환
	return unique[:k]
}

/*
================================================================================
실무 적용
================================================================================

【 인기 검색어 집계 】

type SearchLogger struct {
    counts map[string]int
    mu     sync.RWMutex
}

func (s *SearchLogger) Log(query string) {
    s.mu.Lock()
    s.counts[query]++
    s.mu.Unlock()
}

func (s *SearchLogger) TopK(k int) []string {
    s.mu.RLock()
    defer s.mu.RUnlock()

    // 버킷 정렬 또는 힙 사용
    // ...
}

【 Redis ZSET 활용 】

// 실시간 순위 집계에 Redis Sorted Set 활용
ZINCRBY searches 1 "golang tutorial"
ZINCRBY searches 1 "golang concurrency"
ZREVRANGE searches 0 9 WITHSCORES  // Top 10
================================================================================
*/

// 테스트
func main() {
	testCases := []struct {
		nums []int
		k    int
	}{
		{[]int{1, 1, 1, 2, 2, 3}, 2},
		{[]int{1}, 1},
		{[]int{1, 2}, 2},
	}

	for i, tc := range testCases {
		result := TopKFrequent(tc.nums, tc.k)
		print("Test ", i+1, " (Bucket): ")
		for _, v := range result {
			print(v, " ")
		}

		resultHeap := TopKFrequentHeap(tc.nums, tc.k)
		print(" | Heap: ")
		for _, v := range resultHeap {
			print(v, " ")
		}
		println()
	}
}
