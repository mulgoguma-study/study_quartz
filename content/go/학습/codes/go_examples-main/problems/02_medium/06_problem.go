package medium

/*
================================================================================
문제 6: Merge Intervals (구간 병합)
================================================================================

난이도: Medium
주제: 정렬, 구간(Interval)

【 문제 설명 】
겹치는 구간들을 병합하여 겹치지 않는 구간 배열을 반환하세요.

【 입력 예시 】
intervals = [[1,3], [2,6], [8,10], [15,18]]

【 출력 예시 】
[[1,6], [8,10], [15,18]]

【 설명 】
[1,3]과 [2,6]이 겹치므로 [1,6]으로 병합

【 추가 예시 】
입력: [[1,4], [4,5]]
출력: [[1,5]]

【 제약 조건 】
- 1 <= intervals.length <= 10^4
- intervals[i].length == 2
- 0 <= start_i <= end_i <= 10^4

【 힌트 】
1. 시작점 기준으로 정렬
2. 순차적으로 순회하며 병합

【 실무 연관성 】
- 회의실 예약 시스템
- 시간대 병합
- 캘린더 이벤트 처리
================================================================================
*/

// Interval 구간 타입
type Interval struct {
	Start, End int
}

// MergeIntervals 구간 병합
func MergeIntervals(intervals []Interval) []Interval {
	// 여기에 코드를 작성하세요
	return nil
}

// 테스트
func main() {
	intervals := []Interval{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	result := MergeIntervals(intervals)
	for _, iv := range result {
		println(iv.Start, "-", iv.End)
	}
	// 예상: [1,6], [8,10], [15,18]
}
