package main

import (
	"fmt"
	"sort"
)

/*
===========================================
6-6. 조건 있는 그리디 (Greedy Algorithm)
===========================================

그리디 알고리즘이란?
- 매 순간 가장 좋아 보이는 선택을 하는 알고리즘
- 지역적 최적해(local optimum)를 선택하여 전역 최적해(global optimum)를 찾음
- 모든 문제에 적용되지 않음 - 그리디 선택 속성과 최적 부분 구조 필요

그리디가 적용 가능한 조건:
1. 그리디 선택 속성: 현재의 최선 선택이 최종 해결책에 포함
2. 최적 부분 구조: 부분 문제의 최적해가 전체 문제의 최적해를 구성

DP vs Greedy:
- DP: 모든 가능성 고려, 항상 최적해 보장
- Greedy: 현재 최선만 선택, 빠르지만 최적해 보장 X (문제에 따라 다름)
*/

// ============================================
// 기본 그리디 문제
// ============================================

// CoinChangeGreedy 동전 거스름돈 (그리디 - 특정 화폐 시스템에서만 최적)
// 화폐 단위가 배수 관계일 때만 최적해 보장
func CoinChangeGreedy(coins []int, amount int) int {
	// 큰 동전부터 사용하기 위해 내림차순 정렬
	sortedCoins := make([]int, len(coins))
	copy(sortedCoins, coins)
	sort.Sort(sort.Reverse(sort.IntSlice(sortedCoins)))

	count := 0
	remaining := amount

	for _, coin := range sortedCoins {
		if remaining >= coin {
			// 현재 동전을 최대한 많이 사용
			numCoins := remaining / coin
			count += numCoins
			remaining -= numCoins * coin
		}
	}

	if remaining > 0 {
		return -1 // 거스름돈을 만들 수 없음
	}
	return count
}

// ============================================
// 활동 선택 문제 (Activity Selection)
// ============================================

// Activity 활동 구조체
type Activity struct {
	Start int // 시작 시간
	End   int // 종료 시간
	Name  string
}

// ActivitySelection 최대 활동 수 선택
// 가장 빨리 끝나는 활동부터 선택하면 최적해
// 시간복잡도: O(n log n)
func ActivitySelection(activities []Activity) []Activity {
	if len(activities) == 0 {
		return nil
	}

	// 종료 시간 기준 오름차순 정렬
	sorted := make([]Activity, len(activities))
	copy(sorted, activities)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].End < sorted[j].End
	})

	// 첫 번째 활동 선택
	selected := []Activity{sorted[0]}
	lastEnd := sorted[0].End

	// 그리디: 이전 활동 종료 후 시작하는 가장 빨리 끝나는 활동 선택
	for i := 1; i < len(sorted); i++ {
		if sorted[i].Start >= lastEnd {
			selected = append(selected, sorted[i])
			lastEnd = sorted[i].End
		}
	}

	return selected
}

// ============================================
// 회의실 배정 (Meeting Rooms)
// ============================================

// MinMeetingRooms 필요한 최소 회의실 수
// 시간복잡도: O(n log n)
func MinMeetingRooms(intervals [][]int) int {
	if len(intervals) == 0 {
		return 0
	}

	n := len(intervals)

	// 시작 시간과 종료 시간을 분리하여 정렬
	starts := make([]int, n)
	ends := make([]int, n)

	for i, interval := range intervals {
		starts[i] = interval[0]
		ends[i] = interval[1]
	}

	sort.Ints(starts)
	sort.Ints(ends)

	rooms := 0
	endPtr := 0

	// 시작 시간순으로 순회
	for _, start := range starts {
		// 회의가 끝난 방이 있으면 재사용
		if start >= ends[endPtr] {
			endPtr++
		} else {
			// 새 방 필요
			rooms++
		}
	}

	return rooms
}

// ============================================
// 구간 관련 문제
// ============================================

// Interval 구간 구조체
type Interval struct {
	Start, End int
}

// MergeIntervals 겹치는 구간 병합
// 시간복잡도: O(n log n)
func MergeIntervals(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return nil
	}

	// 시작점 기준 정렬
	sorted := make([]Interval, len(intervals))
	copy(sorted, intervals)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Start < sorted[j].Start
	})

	result := []Interval{sorted[0]}

	for i := 1; i < len(sorted); i++ {
		last := &result[len(result)-1]

		// 현재 구간이 이전 구간과 겹치면 병합
		if sorted[i].Start <= last.End {
			last.End = max(last.End, sorted[i].End)
		} else {
			// 겹치지 않으면 새 구간 추가
			result = append(result, sorted[i])
		}
	}

	return result
}

// MinArrowsToBurstBalloons 풍선 터뜨리는 최소 화살 수
// 각 풍선은 [start, end] 범위에 있음
func MinArrowsToBurstBalloons(balloons [][]int) int {
	if len(balloons) == 0 {
		return 0
	}

	// 끝점 기준 정렬
	sorted := make([][]int, len(balloons))
	copy(sorted, balloons)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i][1] < sorted[j][1]
	})

	arrows := 1
	arrowPos := sorted[0][1] // 첫 화살은 첫 풍선의 끝점에서 발사

	for i := 1; i < len(sorted); i++ {
		// 현재 풍선이 화살 위치 이후에 시작하면 새 화살 필요
		if sorted[i][0] > arrowPos {
			arrows++
			arrowPos = sorted[i][1]
		}
	}

	return arrows
}

// EraseOverlapIntervals 겹치는 구간 제거 최소 개수
func EraseOverlapIntervals(intervals [][]int) int {
	if len(intervals) <= 1 {
		return 0
	}

	// 끝점 기준 정렬
	sorted := make([][]int, len(intervals))
	copy(sorted, intervals)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i][1] < sorted[j][1]
	})

	count := 1      // 유지할 구간 수
	end := sorted[0][1]

	for i := 1; i < len(sorted); i++ {
		// 겹치지 않으면 유지
		if sorted[i][0] >= end {
			count++
			end = sorted[i][1]
		}
	}

	// 제거해야 할 구간 수 = 전체 - 유지할 구간
	return len(intervals) - count
}

// ============================================
// 작업 스케줄링
// ============================================

// Job 작업 구조체
type Job struct {
	Deadline int // 마감 기한
	Profit   int // 이익
}

// JobScheduling 작업 스케줄링 (마감 기한 내 최대 이익)
// 각 작업은 단위 시간이 걸림, 마감 기한 전에 완료해야 함
func JobScheduling(jobs []Job) (int, []int) {
	n := len(jobs)
	if n == 0 {
		return 0, nil
	}

	// 이익 기준 내림차순 정렬
	sorted := make([]Job, n)
	copy(sorted, jobs)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Profit > sorted[j].Profit
	})

	// 최대 마감 기한 찾기
	maxDeadline := 0
	for _, job := range jobs {
		if job.Deadline > maxDeadline {
			maxDeadline = job.Deadline
		}
	}

	// 슬롯 배열 (-1은 비어있음을 의미)
	slots := make([]int, maxDeadline)
	for i := range slots {
		slots[i] = -1
	}

	totalProfit := 0
	scheduled := []int{}

	for i, job := range sorted {
		// 마감 기한 직전부터 빈 슬롯 찾기
		for j := min(maxDeadline, job.Deadline) - 1; j >= 0; j-- {
			if slots[j] == -1 {
				slots[j] = i
				totalProfit += job.Profit
				scheduled = append(scheduled, i)
				break
			}
		}
	}

	return totalProfit, scheduled
}

// ============================================
// 분수 배낭 문제 (Fractional Knapsack)
// ============================================

// Item 물건 구조체
type Item struct {
	Weight float64
	Value  float64
}

// FractionalKnapsack 분수 배낭 (물건을 쪼갤 수 있음)
// 0/1 배낭과 달리 그리디로 최적해 보장
func FractionalKnapsack(items []Item, capacity float64) float64 {
	// 가치/무게 비율 기준 내림차순 정렬
	sorted := make([]Item, len(items))
	copy(sorted, items)
	sort.Slice(sorted, func(i, j int) bool {
		ratioI := sorted[i].Value / sorted[i].Weight
		ratioJ := sorted[j].Value / sorted[j].Weight
		return ratioI > ratioJ
	})

	totalValue := 0.0
	remainingCapacity := capacity

	for _, item := range sorted {
		if remainingCapacity == 0 {
			break
		}

		if item.Weight <= remainingCapacity {
			// 물건 전체를 담을 수 있으면 전부 담음
			totalValue += item.Value
			remainingCapacity -= item.Weight
		} else {
			// 일부만 담음
			fraction := remainingCapacity / item.Weight
			totalValue += item.Value * fraction
			remainingCapacity = 0
		}
	}

	return totalValue
}

// ============================================
// 점프 게임 (Jump Game)
// ============================================

// CanJump 배열의 끝까지 도달 가능한지 확인
// arr[i]는 i위치에서 점프할 수 있는 최대 거리
func CanJump(arr []int) bool {
	maxReach := 0 // 현재까지 도달할 수 있는 최대 위치

	for i := 0; i < len(arr); i++ {
		// 현재 위치에 도달할 수 없으면 실패
		if i > maxReach {
			return false
		}
		// 최대 도달 위치 갱신
		maxReach = max(maxReach, i+arr[i])

		// 이미 끝에 도달 가능하면 성공
		if maxReach >= len(arr)-1 {
			return true
		}
	}

	return true
}

// MinJumps 배열의 끝까지 도달하는 최소 점프 수
func MinJumps(arr []int) int {
	if len(arr) <= 1 {
		return 0
	}

	jumps := 0
	currentEnd := 0  // 현재 점프로 도달 가능한 마지막 위치
	farthest := 0    // 다음 점프로 도달 가능한 최대 위치

	for i := 0; i < len(arr)-1; i++ {
		farthest = max(farthest, i+arr[i])

		// 현재 점프 범위 끝에 도달하면
		if i == currentEnd {
			jumps++
			currentEnd = farthest

			// 이미 끝에 도달 가능하면 종료
			if currentEnd >= len(arr)-1 {
				break
			}
		}
	}

	return jumps
}

// ============================================
// 주유소 (Gas Station)
// ============================================

// CanCompleteCircuit 원형 경로 완주 가능한 시작점
// gas[i]: i번째 주유소에서 얻는 연료
// cost[i]: i에서 i+1로 이동하는데 필요한 연료
func CanCompleteCircuit(gas, cost []int) int {
	totalTank := 0   // 전체 연료 차이
	currentTank := 0 // 현재 연료
	startStation := 0

	for i := 0; i < len(gas); i++ {
		totalTank += gas[i] - cost[i]
		currentTank += gas[i] - cost[i]

		// 연료가 부족하면 다음 위치에서 시작
		if currentTank < 0 {
			startStation = i + 1
			currentTank = 0
		}
	}

	// 전체 연료가 충분하면 startStation에서 출발 가능
	if totalTank >= 0 {
		return startStation
	}
	return -1
}

// ============================================
// 문자열 관련 그리디
// ============================================

// RemoveKDigits 숫자 문자열에서 k개 숫자를 제거하여 최소값 만들기
func RemoveKDigits(num string, k int) string {
	stack := []byte{}

	for i := 0; i < len(num); i++ {
		// 스택의 마지막 숫자가 현재 숫자보다 크면 제거
		for len(stack) > 0 && k > 0 && stack[len(stack)-1] > num[i] {
			stack = stack[:len(stack)-1]
			k--
		}
		stack = append(stack, num[i])
	}

	// 아직 제거할 숫자가 남아있으면 뒤에서 제거
	stack = stack[:len(stack)-k]

	// 앞의 0 제거
	start := 0
	for start < len(stack) && stack[start] == '0' {
		start++
	}

	if start == len(stack) {
		return "0"
	}
	return string(stack[start:])
}

// PartitionLabels 문자열 파티션 (각 문자가 최대 한 파티션에만 등장)
func PartitionLabels(s string) []int {
	// 각 문자의 마지막 등장 위치
	lastIndex := make(map[byte]int)
	for i := 0; i < len(s); i++ {
		lastIndex[s[i]] = i
	}

	result := []int{}
	start := 0
	end := 0

	for i := 0; i < len(s); i++ {
		// 현재 문자의 마지막 위치로 end 갱신
		end = max(end, lastIndex[s[i]])

		// 현재 위치가 end와 같으면 파티션 완성
		if i == end {
			result = append(result, end-start+1)
			start = i + 1
		}
	}

	return result
}

func main() {
	fmt.Println("============================================")
	fmt.Println("6-6. 조건 있는 그리디 (Greedy Algorithm)")
	fmt.Println("============================================")

	// 동전 거스름돈
	fmt.Println("\n[동전 거스름돈 - 그리디]")
	fmt.Println("--------------------------------------------")
	coins := []int{500, 100, 50, 10}
	amount := 1260
	fmt.Printf("동전: %v, 금액: %d\n", coins, amount)
	fmt.Printf("최소 동전 개수: %d\n", CoinChangeGreedy(coins, amount))

	// 활동 선택
	fmt.Println("\n[활동 선택 문제]")
	fmt.Println("--------------------------------------------")
	activities := []Activity{
		{1, 4, "A"},
		{3, 5, "B"},
		{0, 6, "C"},
		{5, 7, "D"},
		{3, 9, "E"},
		{5, 9, "F"},
		{6, 10, "G"},
		{8, 11, "H"},
		{8, 12, "I"},
		{2, 14, "J"},
		{12, 16, "K"},
	}
	fmt.Println("활동들:")
	for _, a := range activities {
		fmt.Printf("  %s: [%d, %d)\n", a.Name, a.Start, a.End)
	}
	selected := ActivitySelection(activities)
	fmt.Printf("선택된 활동: ")
	for _, a := range selected {
		fmt.Printf("%s ", a.Name)
	}
	fmt.Printf("(%d개)\n", len(selected))

	// 회의실 배정
	fmt.Println("\n[회의실 배정]")
	meetings := [][]int{{0, 30}, {5, 10}, {15, 20}}
	fmt.Printf("회의: %v\n", meetings)
	fmt.Printf("필요한 최소 회의실 수: %d\n", MinMeetingRooms(meetings))

	// 구간 병합
	fmt.Println("\n[구간 병합]")
	intervals := []Interval{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	fmt.Printf("구간: ")
	for _, iv := range intervals {
		fmt.Printf("[%d,%d] ", iv.Start, iv.End)
	}
	fmt.Println()
	merged := MergeIntervals(intervals)
	fmt.Printf("병합 후: ")
	for _, iv := range merged {
		fmt.Printf("[%d,%d] ", iv.Start, iv.End)
	}
	fmt.Println()

	// 풍선 터뜨리기
	fmt.Println("\n[풍선 터뜨리기]")
	balloons := [][]int{{10, 16}, {2, 8}, {1, 6}, {7, 12}}
	fmt.Printf("풍선: %v\n", balloons)
	fmt.Printf("필요한 최소 화살 수: %d\n", MinArrowsToBurstBalloons(balloons))

	// 겹치는 구간 제거
	fmt.Println("\n[겹치는 구간 제거]")
	overlaps := [][]int{{1, 2}, {2, 3}, {3, 4}, {1, 3}}
	fmt.Printf("구간: %v\n", overlaps)
	fmt.Printf("제거해야 할 최소 구간 수: %d\n", EraseOverlapIntervals(overlaps))

	// 작업 스케줄링
	fmt.Println("\n[작업 스케줄링]")
	jobs := []Job{
		{4, 70},
		{2, 60},
		{4, 50},
		{3, 40},
		{1, 30},
		{4, 20},
		{6, 10},
	}
	fmt.Println("작업 (마감, 이익):")
	for i, j := range jobs {
		fmt.Printf("  작업%d: (%d, %d)\n", i, j.Deadline, j.Profit)
	}
	profit, scheduled := JobScheduling(jobs)
	fmt.Printf("최대 이익: %d, 스케줄: %v\n", profit, scheduled)

	// 분수 배낭
	fmt.Println("\n[분수 배낭 문제]")
	items := []Item{
		{10, 60},
		{20, 100},
		{30, 120},
	}
	capacity := 50.0
	fmt.Println("물건 (무게, 가치):")
	for _, item := range items {
		fmt.Printf("  (%.0f, %.0f)\n", item.Weight, item.Value)
	}
	fmt.Printf("용량: %.0f\n", capacity)
	fmt.Printf("최대 가치: %.2f\n", FractionalKnapsack(items, capacity))

	// 점프 게임
	fmt.Println("\n[점프 게임]")
	arr1 := []int{2, 3, 1, 1, 4}
	fmt.Printf("배열: %v\n", arr1)
	fmt.Printf("끝까지 도달 가능: %v\n", CanJump(arr1))
	fmt.Printf("최소 점프 수: %d\n", MinJumps(arr1))

	arr2 := []int{3, 2, 1, 0, 4}
	fmt.Printf("\n배열: %v\n", arr2)
	fmt.Printf("끝까지 도달 가능: %v\n", CanJump(arr2))

	// 주유소
	fmt.Println("\n[주유소 문제]")
	gas := []int{1, 2, 3, 4, 5}
	cost := []int{3, 4, 5, 1, 2}
	fmt.Printf("연료: %v\n", gas)
	fmt.Printf("비용: %v\n", cost)
	start := CanCompleteCircuit(gas, cost)
	fmt.Printf("시작 주유소: %d\n", start)

	// 문자열 그리디
	fmt.Println("\n[문자열 그리디]")
	num := "1432219"
	k := 3
	fmt.Printf("숫자열: \"%s\", 제거 개수: %d\n", num, k)
	fmt.Printf("최소값: \"%s\"\n", RemoveKDigits(num, k))

	str := "ababcbacadefegdehijhklij"
	fmt.Printf("\n문자열: \"%s\"\n", str)
	fmt.Printf("파티션 크기: %v\n", PartitionLabels(str))

	// 패턴 요약
	fmt.Println("\n============================================")
	fmt.Println("그리디 알고리즘 패턴 요약")
	fmt.Println("============================================")
	fmt.Println(`
┌─────────────────────────────────────────────────────────────┐
│ 그리디가 최적해를 보장하는 문제                              │
├─────────────────────────────────────────────────────────────┤
│ - 활동 선택: 종료 시간 기준 정렬                            │
│ - 분수 배낭: 가치/무게 비율 기준                            │
│ - 회의실 배정: 시작/종료 시간 분리 처리                     │
│ - 최소 화살/구간 제거: 종료점 기준 정렬                     │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ 그리디 사용 시 주의                                          │
├─────────────────────────────────────────────────────────────┤
│ - 동전 문제: 배수 관계일 때만 최적 (아니면 DP 필요)         │
│ - 0/1 배낭: 그리디 불가 (DP 필요)                           │
│ - 최적해 보장 여부 수학적 증명 필요                          │
└─────────────────────────────────────────────────────────────┘

[일반적인 그리디 전략]
1. 정렬 후 순서대로 처리
2. 매 단계 최선의 선택 (현재 상태에서)
3. 선택 후 되돌아가지 않음

[그리디 적용 가능 신호]
- "최대/최소 개수로 ~하라"
- "가장 ~한 것부터 선택"
- 정렬 후 순차 처리로 해결 가능해 보임
`)
}
