package hard

/*
================================================================================
문제 5: Median from Data Stream - 솔루션
================================================================================

【 Two Heaps 전략 】

작은 절반 (Max-Heap)     큰 절반 (Min-Heap)
       [1, 2, 3]              [4, 5, 6]
            ↑                     ↑
           max=3               min=4

중앙값 = (3 + 4) / 2 = 3.5

규칙:
1. 항상 small.Len() == large.Len() 또는 small.Len() == large.Len() + 1
2. small의 모든 원소 <= large의 모든 원소
*/

import "container/heap"

// MaxHeap 최대 힙 (작은 절반)
type MaxHeap []int
func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// MinHeap 최소 힙 (큰 절반)
type MinHeap []int
func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// MedianFinder 중앙값 계산기
type MedianFinder struct {
	small *MaxHeap // 작은 절반 (최대 힙)
	large *MinHeap // 큰 절반 (최소 힙)
}

func Constructor() MedianFinder {
	small := &MaxHeap{}
	large := &MinHeap{}
	heap.Init(small)
	heap.Init(large)
	return MedianFinder{small: small, large: large}
}

func (mf *MedianFinder) AddNum(num int) {
	// 1. small에 먼저 추가
	heap.Push(mf.small, num)

	// 2. small의 최대값을 large로 이동
	heap.Push(mf.large, heap.Pop(mf.small))

	// 3. 크기 균형 (small >= large)
	if mf.small.Len() < mf.large.Len() {
		heap.Push(mf.small, heap.Pop(mf.large))
	}
}

func (mf *MedianFinder) FindMedian() float64 {
	if mf.small.Len() > mf.large.Len() {
		return float64((*mf.small)[0])
	}
	return float64((*mf.small)[0]+(*mf.large)[0]) / 2.0
}

/*
【 시간 복잡도 】
- AddNum: O(log n)
- FindMedian: O(1)

【 Follow-up 】
1. 모든 정수가 0-100 범위: 카운팅 배열 사용
2. 99% 정수가 0-100: 카운팅 + 예외 힙
*/

func main() {
	mf := Constructor()
	mf.AddNum(1)
	mf.AddNum(2)
	println("Median after [1,2]:", mf.FindMedian())
	mf.AddNum(3)
	println("Median after [1,2,3]:", mf.FindMedian())
}
