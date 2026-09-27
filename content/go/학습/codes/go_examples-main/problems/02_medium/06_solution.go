package medium

/*
================================================================================
문제 6: Merge Intervals - 솔루션
================================================================================
*/

import "sort"

/*
================================================================================
구간 병합 알고리즘
================================================================================

【 핵심 아이디어 】
1. 시작점 기준 정렬
2. 순차적으로 병합 가능 여부 확인

┌─────────────────────────────────────────────────────────────────────────────┐
│                          병합 조건                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Case 1: 병합 필요 (겹침)                                                   │
│  ─────────────────────────                                                  │
│      ├────────┤                  현재 구간                                  │
│           ├────────┤             다음 구간                                  │
│      ├─────────────┤             병합 결과                                  │
│                                                                             │
│      조건: 다음.Start <= 현재.End                                           │
│      결과: End = max(현재.End, 다음.End)                                    │
│                                                                             │
│  Case 2: 병합 불필요 (분리)                                                 │
│  ────────────────────────                                                   │
│      ├────────┤                  현재 구간                                  │
│                   ├────────┤     다음 구간                                  │
│                                                                             │
│      조건: 다음.Start > 현재.End                                            │
│      결과: 현재 구간 확정, 다음 구간으로 이동                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
*/

// Interval 구간 타입
type Interval struct {
	Start, End int
}

// MergeIntervals 구간 병합
/*
【 알고리즘 】
1. 빈 배열 체크
2. 시작점 기준 정렬
3. 첫 번째 구간을 결과에 추가
4. 나머지 구간들을 순회하며:
   - 겹치면: 마지막 결과의 End 확장
   - 안 겹치면: 새 구간 추가

【 시간 복잡도 】O(n log n) - 정렬 때문
【 공간 복잡도 】O(n) - 결과 저장
*/
func MergeIntervals(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return []Interval{}
	}

	// 1. 시작점 기준 정렬
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].Start < intervals[j].Start
	})

	// 2. 결과 배열 초기화 (첫 번째 구간)
	result := []Interval{intervals[0]}

	// 3. 순차적으로 병합
	for i := 1; i < len(intervals); i++ {
		last := &result[len(result)-1] // 마지막 병합 구간
		curr := intervals[i]           // 현재 검사 구간

		if curr.Start <= last.End {
			// 겹침: End 확장 (더 큰 값으로)
			if curr.End > last.End {
				last.End = curr.End
			}
		} else {
			// 안 겹침: 새 구간 추가
			result = append(result, curr)
		}
	}

	return result
}

/*
================================================================================
관련 문제 패턴
================================================================================

【 1. Insert Interval (구간 삽입) 】
기존 정렬된 구간에 새 구간 삽입 후 병합

func InsertInterval(intervals []Interval, newInterval Interval) []Interval {
    result := []Interval{}
    i := 0
    n := len(intervals)

    // 1. 새 구간 이전 구간들 추가
    for i < n && intervals[i].End < newInterval.Start {
        result = append(result, intervals[i])
        i++
    }

    // 2. 겹치는 구간들 병합
    for i < n && intervals[i].Start <= newInterval.End {
        newInterval.Start = min(newInterval.Start, intervals[i].Start)
        newInterval.End = max(newInterval.End, intervals[i].End)
        i++
    }
    result = append(result, newInterval)

    // 3. 나머지 구간들 추가
    for i < n {
        result = append(result, intervals[i])
        i++
    }

    return result
}

【 2. Meeting Rooms (회의실) 】
모든 회의에 참석 가능한지 확인

func CanAttendMeetings(intervals []Interval) bool {
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i].Start < intervals[j].Start
    })

    for i := 1; i < len(intervals); i++ {
        if intervals[i].Start < intervals[i-1].End {
            return false  // 겹침 발견
        }
    }
    return true
}

【 3. Meeting Rooms II (최소 회의실 수) 】
동시에 필요한 최대 회의실 수

func MinMeetingRooms(intervals []Interval) int {
    starts := make([]int, len(intervals))
    ends := make([]int, len(intervals))

    for i, iv := range intervals {
        starts[i] = iv.Start
        ends[i] = iv.End
    }

    sort.Ints(starts)
    sort.Ints(ends)

    rooms := 0
    maxRooms := 0
    s, e := 0, 0

    for s < len(starts) {
        if starts[s] < ends[e] {
            rooms++
            s++
        } else {
            rooms--
            e++
        }
        if rooms > maxRooms {
            maxRooms = rooms
        }
    }

    return maxRooms
}

【 4. Non-overlapping Intervals (제거할 최소 구간 수) 】
겹치지 않도록 하기 위해 제거해야 할 최소 구간 수

func EraseOverlapIntervals(intervals []Interval) int {
    if len(intervals) == 0 {
        return 0
    }

    // 끝점 기준 정렬 (Greedy)
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i].End < intervals[j].End
    })

    count := 0
    prevEnd := intervals[0].End

    for i := 1; i < len(intervals); i++ {
        if intervals[i].Start < prevEnd {
            count++  // 제거
        } else {
            prevEnd = intervals[i].End
        }
    }

    return count
}
================================================================================
*/

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 테스트
func main() {
	testCases := []struct {
		intervals []Interval
		expect    []Interval
	}{
		{
			[]Interval{{1, 3}, {2, 6}, {8, 10}, {15, 18}},
			[]Interval{{1, 6}, {8, 10}, {15, 18}},
		},
		{
			[]Interval{{1, 4}, {4, 5}},
			[]Interval{{1, 5}},
		},
		{
			[]Interval{{1, 4}, {0, 4}},
			[]Interval{{0, 4}},
		},
	}

	for i, tc := range testCases {
		result := MergeIntervals(tc.intervals)
		print("Test ", i+1, ": ")
		for _, iv := range result {
			print("[", iv.Start, ",", iv.End, "] ")
		}
		println()
	}
}
