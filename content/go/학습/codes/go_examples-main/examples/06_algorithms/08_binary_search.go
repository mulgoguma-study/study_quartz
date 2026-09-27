package main

import (
	"fmt"
	"math"
	"sort"
)

/*
===========================================
6-8. 이진 탐색 (Binary Search)
===========================================

이진 탐색이란?
- 정렬된 배열에서 특정 값을 O(log n)에 찾는 알고리즘
- 매번 탐색 범위를 절반으로 줄임
- 분할 정복의 대표적인 예

전제 조건:
- 배열이 정렬되어 있어야 함
- 또는 이분 결정이 가능한 단조 함수

응용:
- 값 찾기: target이 있는 인덱스
- 경계 찾기: target 이상/이하인 첫 번째 위치
- 매개변수 탐색: 조건을 만족하는 최소/최대값
*/

// ============================================
// 기본 이진 탐색
// ============================================

// BinarySearch 기본 이진 탐색 (target의 인덱스 반환, 없으면 -1)
// 시간복잡도: O(log n)
func BinarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1

	for left <= right {
		// 오버플로우 방지: (left + right) / 2 대신 사용
		mid := left + (right-left)/2

		if arr[mid] == target {
			return mid // 찾음
		} else if arr[mid] < target {
			left = mid + 1 // 오른쪽 절반 탐색
		} else {
			right = mid - 1 // 왼쪽 절반 탐색
		}
	}

	return -1 // 못 찾음
}

// BinarySearchRecursive 재귀 이진 탐색
func BinarySearchRecursive(arr []int, target, left, right int) int {
	if left > right {
		return -1
	}

	mid := left + (right-left)/2

	if arr[mid] == target {
		return mid
	} else if arr[mid] < target {
		return BinarySearchRecursive(arr, target, mid+1, right)
	} else {
		return BinarySearchRecursive(arr, target, left, mid-1)
	}
}

// ============================================
// 경계 찾기 (Lower/Upper Bound)
// ============================================

// LowerBound target 이상인 첫 번째 인덱스 반환
// target보다 작은 모든 요소의 다음 위치
// C++ std::lower_bound와 동일
func LowerBound(arr []int, target int) int {
	left, right := 0, len(arr)

	for left < right {
		mid := left + (right-left)/2

		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid // arr[mid] >= target, mid가 답이 될 수 있음
		}
	}

	return left
}

// UpperBound target보다 큰 첫 번째 인덱스 반환
// target 이하인 모든 요소의 다음 위치
// C++ std::upper_bound와 동일
func UpperBound(arr []int, target int) int {
	left, right := 0, len(arr)

	for left < right {
		mid := left + (right-left)/2

		if arr[mid] <= target {
			left = mid + 1
		} else {
			right = mid // arr[mid] > target, mid가 답이 될 수 있음
		}
	}

	return left
}

// CountOccurrences 정렬된 배열에서 target 등장 횟수
func CountOccurrences(arr []int, target int) int {
	lower := LowerBound(arr, target)
	upper := UpperBound(arr, target)
	return upper - lower
}

// SearchRange target의 시작과 끝 인덱스
func SearchRange(arr []int, target int) []int {
	lower := LowerBound(arr, target)

	// target이 없으면 [-1, -1]
	if lower == len(arr) || arr[lower] != target {
		return []int{-1, -1}
	}

	upper := UpperBound(arr, target) - 1
	return []int{lower, upper}
}

// ============================================
// 변형: 회전된 배열에서 탐색
// ============================================

// SearchRotated 회전된 정렬 배열에서 target 찾기
// 예: [4, 5, 6, 7, 0, 1, 2] (7에서 회전)
func SearchRotated(arr []int, target int) int {
	left, right := 0, len(arr)-1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			return mid
		}

		// 왼쪽 절반이 정렬되어 있는지 확인
		if arr[left] <= arr[mid] {
			// target이 왼쪽 정렬 구간에 있는지
			if arr[left] <= target && target < arr[mid] {
				right = mid - 1
			} else {
				left = mid + 1
			}
		} else {
			// 오른쪽 절반이 정렬됨
			if arr[mid] < target && target <= arr[right] {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
	}

	return -1
}

// FindMin 회전된 정렬 배열의 최소값 찾기
func FindMin(arr []int) int {
	left, right := 0, len(arr)-1

	for left < right {
		mid := left + (right-left)/2

		// mid가 right보다 크면 최소값은 오른쪽에
		if arr[mid] > arr[right] {
			left = mid + 1
		} else {
			right = mid
		}
	}

	return arr[left]
}

// FindRotationCount 회전 횟수 = 최소값의 인덱스
func FindRotationCount(arr []int) int {
	left, right := 0, len(arr)-1

	// 이미 정렬됨
	if arr[left] <= arr[right] {
		return 0
	}

	for left < right {
		mid := left + (right-left)/2

		if arr[mid] > arr[right] {
			left = mid + 1
		} else {
			right = mid
		}
	}

	return left
}

// ============================================
// 매개변수 탐색 (Parametric Search)
// ============================================

// MinDaysToMakeBouquets m개의 부케를 만들 수 있는 최소 일수
// bloomDay[i]: i번째 꽃이 피는 날, k: 부케당 필요한 인접 꽃 수
func MinDaysToMakeBouquets(bloomDay []int, m, k int) int {
	n := len(bloomDay)

	// 불가능한 경우
	if m*k > n {
		return -1
	}

	// 탐색 범위: 최소 bloomDay ~ 최대 bloomDay
	left, right := bloomDay[0], bloomDay[0]
	for _, day := range bloomDay {
		if day < left {
			left = day
		}
		if day > right {
			right = day
		}
	}

	// day일에 부케를 m개 이상 만들 수 있는지 확인하는 함수
	canMake := func(day int) bool {
		bouquets := 0
		consecutive := 0

		for i := 0; i < n; i++ {
			if bloomDay[i] <= day {
				consecutive++
				if consecutive == k {
					bouquets++
					consecutive = 0
				}
			} else {
				consecutive = 0
			}
		}

		return bouquets >= m
	}

	// 조건을 만족하는 최소 day 찾기
	for left < right {
		mid := left + (right-left)/2

		if canMake(mid) {
			right = mid // 가능하면 더 작은 값 시도
		} else {
			left = mid + 1 // 불가능하면 더 큰 값 시도
		}
	}

	return left
}

// KokoEatingBananas 코코가 바나나를 다 먹는 최소 속도
// piles[i]: i번째 더미의 바나나 수, h: 주어진 시간
func KokoEatingBananas(piles []int, h int) int {
	// 탐색 범위: 1 ~ max(piles)
	left, right := 1, piles[0]
	for _, p := range piles {
		if p > right {
			right = p
		}
	}

	// 속도 k로 h시간 안에 다 먹을 수 있는지
	canFinish := func(k int) bool {
		hours := 0
		for _, pile := range piles {
			// 올림 나눗셈
			hours += (pile + k - 1) / k
		}
		return hours <= h
	}

	for left < right {
		mid := left + (right-left)/2

		if canFinish(mid) {
			right = mid
		} else {
			left = mid + 1
		}
	}

	return left
}

// ShipPackages capacity로 d일 안에 모든 패키지 배송 가능한 최소 capacity
func ShipPackages(weights []int, d int) int {
	// 최소 capacity: 가장 무거운 짐 무게
	// 최대 capacity: 모든 짐의 합
	left, right := 0, 0
	for _, w := range weights {
		if w > left {
			left = w
		}
		right += w
	}

	// capacity로 d일 안에 배송 가능한지
	canShip := func(capacity int) bool {
		days := 1
		currentLoad := 0

		for _, w := range weights {
			if currentLoad+w > capacity {
				days++
				currentLoad = w
			} else {
				currentLoad += w
			}
		}

		return days <= d
	}

	for left < right {
		mid := left + (right-left)/2

		if canShip(mid) {
			right = mid
		} else {
			left = mid + 1
		}
	}

	return left
}

// ============================================
// 실수 이진 탐색
// ============================================

// Sqrt 제곱근 계산 (정수)
func Sqrt(x int) int {
	if x < 2 {
		return x
	}

	left, right := 1, x/2

	for left <= right {
		mid := left + (right-left)/2
		square := mid * mid

		if square == x {
			return mid
		} else if square < x {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return right // floor(sqrt(x))
}

// SqrtFloat 제곱근 계산 (실수, 정밀도 지정)
func SqrtFloat(x float64, precision float64) float64 {
	if x < 0 {
		return math.NaN()
	}
	if x == 0 {
		return 0
	}

	left, right := 0.0, max(1.0, x)

	for right-left > precision {
		mid := (left + right) / 2

		if mid*mid < x {
			left = mid
		} else {
			right = mid
		}
	}

	return (left + right) / 2
}

// ============================================
// 2D 이진 탐색
// ============================================

// SearchMatrix 행렬에서 target 찾기
// 각 행과 열이 정렬된 m×n 행렬
func SearchMatrix(matrix [][]int, target int) bool {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return false
	}

	rows, cols := len(matrix), len(matrix[0])

	// 행렬을 1D 배열로 취급
	left, right := 0, rows*cols-1

	for left <= right {
		mid := left + (right-left)/2
		// 1D 인덱스를 2D 좌표로 변환
		row := mid / cols
		col := mid % cols
		val := matrix[row][col]

		if val == target {
			return true
		} else if val < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return false
}

// SearchMatrix2 정렬된 2D 행렬에서 탐색 (행별, 열별 정렬)
// 각 행은 왼쪽에서 오른쪽으로, 각 열은 위에서 아래로 정렬
func SearchMatrix2(matrix [][]int, target int) bool {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return false
	}

	rows, cols := len(matrix), len(matrix[0])
	row, col := 0, cols-1 // 오른쪽 상단에서 시작

	for row < rows && col >= 0 {
		if matrix[row][col] == target {
			return true
		} else if matrix[row][col] < target {
			row++ // 아래로 이동
		} else {
			col-- // 왼쪽으로 이동
		}
	}

	return false
}

func main() {
	fmt.Println("============================================")
	fmt.Println("6-8. 이진 탐색 (Binary Search)")
	fmt.Println("============================================")

	// 기본 이진 탐색
	fmt.Println("\n[기본 이진 탐색]")
	fmt.Println("--------------------------------------------")
	arr := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
	fmt.Printf("배열: %v\n", arr)

	target := 7
	idx := BinarySearch(arr, target)
	fmt.Printf("target=%d, 인덱스=%d\n", target, idx)

	target = 8
	idx = BinarySearch(arr, target)
	fmt.Printf("target=%d, 인덱스=%d (없음)\n", target, idx)

	// Lower/Upper Bound
	fmt.Println("\n[경계 찾기 (Lower/Upper Bound)]")
	fmt.Println("--------------------------------------------")
	arr2 := []int{1, 2, 2, 2, 3, 4, 5}
	fmt.Printf("배열: %v\n", arr2)

	target = 2
	fmt.Printf("LowerBound(%d) = %d (인덱스)\n", target, LowerBound(arr2, target))
	fmt.Printf("UpperBound(%d) = %d (인덱스)\n", target, UpperBound(arr2, target))
	fmt.Printf("등장 횟수: %d\n", CountOccurrences(arr2, target))
	fmt.Printf("범위: %v\n", SearchRange(arr2, target))

	// Go 표준 라이브러리 sort.Search
	fmt.Println("\n[Go sort.Search 사용]")
	idx = sort.Search(len(arr), func(i int) bool {
		return arr[i] >= 7
	})
	fmt.Printf("sort.Search로 7 이상인 첫 인덱스: %d\n", idx)

	// 회전된 배열
	fmt.Println("\n[회전된 정렬 배열]")
	fmt.Println("--------------------------------------------")
	rotated := []int{4, 5, 6, 7, 0, 1, 2}
	fmt.Printf("배열: %v\n", rotated)
	fmt.Printf("최소값: %d\n", FindMin(rotated))
	fmt.Printf("회전 횟수: %d\n", FindRotationCount(rotated))
	fmt.Printf("target=0 인덱스: %d\n", SearchRotated(rotated, 0))
	fmt.Printf("target=3 인덱스: %d (없음)\n", SearchRotated(rotated, 3))

	// 매개변수 탐색
	fmt.Println("\n============================================")
	fmt.Println("[매개변수 탐색 (Parametric Search)]")
	fmt.Println("============================================")

	// 부케 만들기
	fmt.Println("\n[부케 만들기]")
	bloomDay := []int{1, 10, 3, 10, 2}
	m, k := 3, 1
	fmt.Printf("개화일: %v, 부케 수: %d, 연속 꽃: %d\n", bloomDay, m, k)
	fmt.Printf("최소 대기 일수: %d\n", MinDaysToMakeBouquets(bloomDay, m, k))

	// 바나나 먹기
	fmt.Println("\n[코코의 바나나]")
	piles := []int{3, 6, 7, 11}
	h := 8
	fmt.Printf("바나나 더미: %v, 시간: %d\n", piles, h)
	fmt.Printf("최소 먹는 속도: %d\n", KokoEatingBananas(piles, h))

	// 패키지 배송
	fmt.Println("\n[패키지 배송]")
	weights := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	days := 5
	fmt.Printf("무게: %v, 기한: %d일\n", weights, days)
	fmt.Printf("최소 선적 용량: %d\n", ShipPackages(weights, days))

	// 제곱근
	fmt.Println("\n============================================")
	fmt.Println("[제곱근 계산]")
	fmt.Println("============================================")
	x := 8
	fmt.Printf("Sqrt(%d) = %d (정수)\n", x, Sqrt(x))

	xf := 2.0
	fmt.Printf("SqrtFloat(%.1f) = %.10f\n", xf, SqrtFloat(xf, 1e-10))

	// 2D 행렬 탐색
	fmt.Println("\n============================================")
	fmt.Println("[2D 행렬 탐색]")
	fmt.Println("============================================")

	matrix := [][]int{
		{1, 3, 5, 7},
		{10, 11, 16, 20},
		{23, 30, 34, 60},
	}
	fmt.Println("행렬:")
	for _, row := range matrix {
		fmt.Printf("  %v\n", row)
	}
	fmt.Printf("target=3 존재: %v\n", SearchMatrix(matrix, 3))
	fmt.Printf("target=13 존재: %v\n", SearchMatrix(matrix, 13))

	// 패턴 요약
	fmt.Println("\n============================================")
	fmt.Println("이진 탐색 패턴 요약")
	fmt.Println("============================================")
	fmt.Println(`
┌─────────────────────────────────────────────────────────────┐
│ 기본 이진 탐색 템플릿                                        │
├─────────────────────────────────────────────────────────────┤
│ left, right := 0, len(arr)-1                                │
│ for left <= right {                                         │
│     mid := left + (right-left)/2                            │
│     if arr[mid] == target { return mid }                    │
│     else if arr[mid] < target { left = mid + 1 }            │
│     else { right = mid - 1 }                                │
│ }                                                           │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ Lower/Upper Bound 템플릿                                     │
├─────────────────────────────────────────────────────────────┤
│ // left < right 사용                                        │
│ for left < right {                                          │
│     mid := left + (right-left)/2                            │
│     if condition(mid) { right = mid }                       │
│     else { left = mid + 1 }                                 │
│ }                                                           │
│ return left                                                 │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ 매개변수 탐색 패턴                                           │
├─────────────────────────────────────────────────────────────┤
│ 1. 정답의 범위 [left, right] 설정                           │
│ 2. mid 값으로 조건 만족 여부 판단 함수 작성                  │
│ 3. 조건 만족 시 더 좋은 값 탐색                              │
│    - 최소값 찾기: right = mid                               │
│    - 최대값 찾기: left = mid + 1                            │
└─────────────────────────────────────────────────────────────┘

[사용 시기]
- "정렬된 배열에서 찾기" → 기본 이진 탐색
- "조건을 만족하는 최소/최대값" → 매개변수 탐색
- "~이상/~이하 처음 위치" → Lower/Upper Bound
`)
}
