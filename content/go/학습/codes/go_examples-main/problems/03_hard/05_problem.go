package hard

/*
================================================================================
문제 5: Median from Data Stream (스트림에서 중앙값)
================================================================================

난이도: Hard
주제: 힙, 설계

【 문제 설명 】
정수 스트림에서 지금까지의 중앙값을 효율적으로 계산하세요.

【 인터페이스 】
type MedianFinder struct{}
func (mf *MedianFinder) AddNum(num int)
func (mf *MedianFinder) FindMedian() float64

【 예시 】
mf := Constructor()
mf.AddNum(1)    // [1]
mf.AddNum(2)    // [1, 2]
mf.FindMedian() // 1.5
mf.AddNum(3)    // [1, 2, 3]
mf.FindMedian() // 2.0

【 힌트 】
- Two Heaps: Max-Heap(작은 절반) + Min-Heap(큰 절반)
- 두 힙의 크기 균형 유지
================================================================================
*/

// MedianFinder 중앙값 계산기
type MedianFinder struct {
	// 여기에 필드를 정의하세요
}

func Constructor() MedianFinder {
	return MedianFinder{}
}

func (mf *MedianFinder) AddNum(num int) {
	// 여기에 코드를 작성하세요
}

func (mf *MedianFinder) FindMedian() float64 {
	// 여기에 코드를 작성하세요
	return 0
}

func main() {
	mf := Constructor()
	mf.AddNum(1)
	mf.AddNum(2)
	println("Median:", mf.FindMedian()) // 1.5
	mf.AddNum(3)
	println("Median:", mf.FindMedian()) // 2.0
}
