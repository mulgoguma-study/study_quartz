package main

import (
	"fmt"
	"sort"
)

/*
===========================================
6-4. Two Pointer (투 포인터)
===========================================

투 포인터란?
- 두 개의 포인터를 사용하여 배열/문자열을 효율적으로 탐색
- 보통 정렬된 배열에서 사용하거나, 양쪽 끝에서 중앙으로 이동
- O(n²)를 O(n)으로 최적화 가능

종류:
1. 반대 방향 포인터: 양 끝에서 시작하여 중앙으로
2. 같은 방향 포인터: 둘 다 한쪽에서 시작
3. 빠른/느린 포인터: 서로 다른 속도로 이동

사용 사례:
- 정렬된 배열에서 합/차 찾기
- 팰린드롬 검사
- 중복 제거
- 연결 리스트 사이클 감지
*/

// ============================================
// 반대 방향 포인터
// ============================================

// TwoSumSorted 정렬된 배열에서 합이 target인 두 수 찾기
// 시간복잡도: O(n)
func TwoSumSorted(arr []int, target int) (int, int, bool) {
	left := 0
	right := len(arr) - 1

	for left < right {
		sum := arr[left] + arr[right]

		if sum == target {
			return left, right, true
		} else if sum < target {
			// 합이 작으면 왼쪽 포인터를 오른쪽으로 (더 큰 값 선택)
			left++
		} else {
			// 합이 크면 오른쪽 포인터를 왼쪽으로 (더 작은 값 선택)
			right--
		}
	}

	return 0, 0, false
}

// IsPalindrome 문자열이 팰린드롬인지 확인
// 시간복잡도: O(n)
func IsPalindrome(s string) bool {
	left := 0
	right := len(s) - 1

	for left < right {
		// 양쪽 문자가 다르면 팰린드롬 아님
		if s[left] != s[right] {
			return false
		}
		left++
		right--
	}

	return true
}

// IsPalindromeIgnoreCase 대소문자/공백 무시하고 팰린드롬 확인
func IsPalindromeIgnoreCase(s string) bool {
	left := 0
	right := len(s) - 1

	for left < right {
		// 알파벳이 아닌 문자 건너뛰기
		for left < right && !isAlphanumeric(s[left]) {
			left++
		}
		for left < right && !isAlphanumeric(s[right]) {
			right--
		}

		// 소문자로 변환 후 비교
		if toLower(s[left]) != toLower(s[right]) {
			return false
		}
		left++
		right--
	}

	return true
}

// 헬퍼 함수들
func isAlphanumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}

// ReverseArray 배열 뒤집기
func ReverseArray(arr []int) {
	left := 0
	right := len(arr) - 1

	for left < right {
		// 두 요소 교환
		arr[left], arr[right] = arr[right], arr[left]
		left++
		right--
	}
}

// ContainerWithMostWater 물을 가장 많이 담을 수 있는 컨테이너
// arr[i]는 i번째 막대의 높이
// 시간복잡도: O(n)
func ContainerWithMostWater(heights []int) int {
	left := 0
	right := len(heights) - 1
	maxArea := 0

	for left < right {
		// 물의 양 = 너비 × 높이(낮은 막대 기준)
		width := right - left
		height := min(heights[left], heights[right])
		area := width * height

		if area > maxArea {
			maxArea = area
		}

		// 낮은 쪽 포인터를 이동 (더 높은 막대 찾기)
		if heights[left] < heights[right] {
			left++
		} else {
			right--
		}
	}

	return maxArea
}

// ThreeSum 세 수의 합이 target인 모든 조합 찾기
// 시간복잡도: O(n²)
func ThreeSum(arr []int, target int) [][]int {
	result := [][]int{}
	n := len(arr)

	// 먼저 정렬
	sort.Ints(arr)

	for i := 0; i < n-2; i++ {
		// 중복 건너뛰기
		if i > 0 && arr[i] == arr[i-1] {
			continue
		}

		// 나머지 두 수는 투 포인터로 찾기
		left := i + 1
		right := n - 1

		for left < right {
			sum := arr[i] + arr[left] + arr[right]

			if sum == target {
				result = append(result, []int{arr[i], arr[left], arr[right]})

				// 중복 건너뛰기
				for left < right && arr[left] == arr[left+1] {
					left++
				}
				for left < right && arr[right] == arr[right-1] {
					right--
				}

				left++
				right--
			} else if sum < target {
				left++
			} else {
				right--
			}
		}
	}

	return result
}

// ============================================
// 같은 방향 포인터
// ============================================

// RemoveDuplicatesSorted 정렬된 배열에서 중복 제거 (in-place)
// 새 길이 반환
// 시간복잡도: O(n), 공간복잡도: O(1)
func RemoveDuplicatesSorted(arr []int) int {
	if len(arr) == 0 {
		return 0
	}

	// slow: 중복 제거된 배열의 마지막 위치
	// fast: 현재 검사 중인 위치
	slow := 0

	for fast := 1; fast < len(arr); fast++ {
		// 새로운 값을 발견하면
		if arr[fast] != arr[slow] {
			slow++
			arr[slow] = arr[fast]
		}
	}

	// slow + 1이 새 배열의 길이
	return slow + 1
}

// RemoveElement 특정 값 제거 (in-place)
func RemoveElement(arr []int, val int) int {
	slow := 0

	for fast := 0; fast < len(arr); fast++ {
		// val이 아닌 요소만 앞으로 옮김
		if arr[fast] != val {
			arr[slow] = arr[fast]
			slow++
		}
	}

	return slow
}

// MoveZeros 0을 배열 끝으로 이동 (순서 유지)
func MoveZeros(arr []int) {
	slow := 0 // 0이 아닌 요소를 놓을 위치

	// 1단계: 0이 아닌 요소를 앞으로 옮김
	for fast := 0; fast < len(arr); fast++ {
		if arr[fast] != 0 {
			arr[slow] = arr[fast]
			slow++
		}
	}

	// 2단계: 나머지를 0으로 채움
	for i := slow; i < len(arr); i++ {
		arr[i] = 0
	}
}

// MoveZerosSwap 0을 배열 끝으로 이동 (스왑 방식)
func MoveZerosSwap(arr []int) {
	slow := 0

	for fast := 0; fast < len(arr); fast++ {
		if arr[fast] != 0 {
			// 교환
			arr[slow], arr[fast] = arr[fast], arr[slow]
			slow++
		}
	}
}

// SquaresOfSortedArray 정렬된 배열의 제곱을 정렬된 상태로 반환
// 음수가 포함된 정렬된 배열에서 제곱 후 정렬
// 시간복잡도: O(n)
func SquaresOfSortedArray(arr []int) []int {
	n := len(arr)
	result := make([]int, n)

	left := 0
	right := n - 1
	pos := n - 1 // 결과 배열의 끝에서부터 채움

	for left <= right {
		leftSq := arr[left] * arr[left]
		rightSq := arr[right] * arr[right]

		// 절대값이 큰 쪽을 결과 배열 끝에 배치
		if leftSq > rightSq {
			result[pos] = leftSq
			left++
		} else {
			result[pos] = rightSq
			right--
		}
		pos--
	}

	return result
}

// ============================================
// 빠른/느린 포인터 (Fast & Slow Pointer)
// ============================================

// ListNode 연결 리스트 노드
type ListNode struct {
	Val  int
	Next *ListNode
}

// HasCycle 연결 리스트에 사이클이 있는지 확인 (Floyd's Algorithm)
// 시간복잡도: O(n), 공간복잡도: O(1)
func HasCycle(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}

	slow := head      // 한 칸씩 이동
	fast := head.Next // 두 칸씩 이동

	for fast != nil && fast.Next != nil {
		if slow == fast {
			return true // 만나면 사이클 존재
		}
		slow = slow.Next
		fast = fast.Next.Next
	}

	return false
}

// FindCycleStart 사이클 시작점 찾기
func FindCycleStart(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return nil
	}

	slow := head
	fast := head

	// 사이클 존재 여부 확인
	hasCycle := false
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			hasCycle = true
			break
		}
	}

	if !hasCycle {
		return nil
	}

	// slow를 head로 리셋하고 같은 속도로 이동
	// 다시 만나는 지점이 사이클 시작점
	slow = head
	for slow != fast {
		slow = slow.Next
		fast = fast.Next
	}

	return slow
}

// FindMiddle 연결 리스트의 중간 노드 찾기
func FindMiddle(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}

	slow := head
	fast := head

	// fast가 끝에 도달하면 slow는 중간에 있음
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	return slow
}

// FindNthFromEnd 끝에서 n번째 노드 찾기
func FindNthFromEnd(head *ListNode, n int) *ListNode {
	if head == nil {
		return nil
	}

	fast := head
	slow := head

	// fast를 n칸 먼저 이동
	for i := 0; i < n; i++ {
		if fast == nil {
			return nil // n이 리스트 길이보다 큼
		}
		fast = fast.Next
	}

	// 둘 다 한 칸씩 이동
	// fast가 끝에 도달하면 slow는 끝에서 n번째
	for fast != nil {
		slow = slow.Next
		fast = fast.Next
	}

	return slow
}

// ============================================
// 병합 관련
// ============================================

// MergeSortedArrays 두 정렬된 배열 병합
// 시간복잡도: O(n + m)
func MergeSortedArrays(arr1, arr2 []int) []int {
	result := make([]int, 0, len(arr1)+len(arr2))
	i, j := 0, 0

	// 두 배열 비교하며 병합
	for i < len(arr1) && j < len(arr2) {
		if arr1[i] <= arr2[j] {
			result = append(result, arr1[i])
			i++
		} else {
			result = append(result, arr2[j])
			j++
		}
	}

	// 남은 요소 추가
	result = append(result, arr1[i:]...)
	result = append(result, arr2[j:]...)

	return result
}

// Intersection 두 정렬된 배열의 교집합
func Intersection(arr1, arr2 []int) []int {
	result := []int{}
	i, j := 0, 0

	for i < len(arr1) && j < len(arr2) {
		if arr1[i] == arr2[j] {
			// 중복 방지
			if len(result) == 0 || result[len(result)-1] != arr1[i] {
				result = append(result, arr1[i])
			}
			i++
			j++
		} else if arr1[i] < arr2[j] {
			i++
		} else {
			j++
		}
	}

	return result
}

// ============================================
// 문자열 관련
// ============================================

// ReverseString 문자열 뒤집기 (바이트 슬라이스로)
func ReverseString(s []byte) {
	left := 0
	right := len(s) - 1

	for left < right {
		s[left], s[right] = s[right], s[left]
		left++
		right--
	}
}

// ReverseWords 문자열에서 각 단어 뒤집기
// "hello world" -> "olleh dlrow"
func ReverseWords(s string) string {
	bytes := []byte(s)

	left := 0
	for i := 0; i <= len(bytes); i++ {
		// 공백이나 문자열 끝에서 단어 뒤집기
		if i == len(bytes) || bytes[i] == ' ' {
			// left부터 i-1까지 뒤집기
			l, r := left, i-1
			for l < r {
				bytes[l], bytes[r] = bytes[r], bytes[l]
				l++
				r--
			}
			left = i + 1
		}
	}

	return string(bytes)
}

func main() {
	fmt.Println("============================================")
	fmt.Println("6-4. Two Pointer (투 포인터)")
	fmt.Println("============================================")

	// 반대 방향 포인터
	fmt.Println("\n[반대 방향 포인터]")
	fmt.Println("--------------------------------------------")

	// Two Sum (정렬된 배열)
	arr := []int{2, 7, 11, 15}
	target := 9
	fmt.Printf("배열: %v, target=%d\n", arr, target)
	if i, j, found := TwoSumSorted(arr, target); found {
		fmt.Printf("합이 %d인 인덱스: [%d, %d] (값: %d + %d)\n", target, i, j, arr[i], arr[j])
	}

	// 팰린드롬
	fmt.Println("\n[팰린드롬 검사]")
	tests := []string{"racecar", "hello", "A man a plan a canal Panama"}
	for _, s := range tests {
		fmt.Printf("\"%s\": 팰린드롬? %v\n", s, IsPalindrome(s))
	}

	fmt.Println("\n[대소문자/공백 무시 팰린드롬]")
	fmt.Printf("\"%s\": %v\n", tests[2], IsPalindromeIgnoreCase(tests[2]))

	// 배열 뒤집기
	fmt.Println("\n[배열 뒤집기]")
	arr2 := []int{1, 2, 3, 4, 5}
	fmt.Printf("원본: %v\n", arr2)
	ReverseArray(arr2)
	fmt.Printf("뒤집기: %v\n", arr2)

	// 물 컨테이너
	fmt.Println("\n[물 컨테이너]")
	heights := []int{1, 8, 6, 2, 5, 4, 8, 3, 7}
	fmt.Printf("높이: %v\n", heights)
	fmt.Printf("최대 물의 양: %d\n", ContainerWithMostWater(heights))

	// Three Sum
	fmt.Println("\n[Three Sum]")
	arr3 := []int{-1, 0, 1, 2, -1, -4}
	target3 := 0
	fmt.Printf("배열: %v, target=%d\n", arr3, target3)
	fmt.Printf("합이 %d인 세 수: %v\n", target3, ThreeSum(arr3, target3))

	// 같은 방향 포인터
	fmt.Println("\n============================================")
	fmt.Println("[같은 방향 포인터]")
	fmt.Println("============================================")

	// 중복 제거
	arr4 := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	fmt.Printf("원본: %v\n", arr4)
	newLen := RemoveDuplicatesSorted(arr4)
	fmt.Printf("중복 제거 후: %v (길이: %d)\n", arr4[:newLen], newLen)

	// 특정 값 제거
	arr5 := []int{3, 2, 2, 3}
	fmt.Printf("\n원본: %v, 제거할 값: 3\n", arr5)
	newLen2 := RemoveElement(arr5, 3)
	fmt.Printf("제거 후: %v (길이: %d)\n", arr5[:newLen2], newLen2)

	// 0 이동
	arr6 := []int{0, 1, 0, 3, 12}
	fmt.Printf("\n원본: %v\n", arr6)
	MoveZeros(arr6)
	fmt.Printf("0을 끝으로 이동: %v\n", arr6)

	// 정렬된 배열의 제곱
	arr7 := []int{-4, -1, 0, 3, 10}
	fmt.Printf("\n원본: %v\n", arr7)
	fmt.Printf("제곱 후 정렬: %v\n", SquaresOfSortedArray(arr7))

	// 빠른/느린 포인터
	fmt.Println("\n============================================")
	fmt.Println("[빠른/느린 포인터 - 연결 리스트]")
	fmt.Println("============================================")

	// 연결 리스트 생성: 1 -> 2 -> 3 -> 4 -> 5
	head := &ListNode{Val: 1}
	head.Next = &ListNode{Val: 2}
	head.Next.Next = &ListNode{Val: 3}
	head.Next.Next.Next = &ListNode{Val: 4}
	head.Next.Next.Next.Next = &ListNode{Val: 5}

	fmt.Println("연결 리스트: 1 -> 2 -> 3 -> 4 -> 5")
	fmt.Printf("사이클 존재: %v\n", HasCycle(head))
	fmt.Printf("중간 노드 값: %d\n", FindMiddle(head).Val)

	nth := FindNthFromEnd(head, 2)
	fmt.Printf("끝에서 2번째 노드 값: %d\n", nth.Val)

	// 사이클이 있는 리스트
	cycleHead := &ListNode{Val: 1}
	cycleHead.Next = &ListNode{Val: 2}
	cycleHead.Next.Next = &ListNode{Val: 3}
	cycleHead.Next.Next.Next = cycleHead.Next // 3 -> 2 사이클

	fmt.Println("\n연결 리스트 with 사이클: 1 -> 2 -> 3 -> (2로 돌아감)")
	fmt.Printf("사이클 존재: %v\n", HasCycle(cycleHead))

	cycleStart := FindCycleStart(cycleHead)
	if cycleStart != nil {
		fmt.Printf("사이클 시작점 값: %d\n", cycleStart.Val)
	}

	// 병합
	fmt.Println("\n============================================")
	fmt.Println("[배열 병합]")
	fmt.Println("============================================")

	arr8 := []int{1, 3, 5, 7}
	arr9 := []int{2, 4, 6, 8}
	fmt.Printf("배열1: %v\n배열2: %v\n", arr8, arr9)
	fmt.Printf("병합: %v\n", MergeSortedArrays(arr8, arr9))

	arr10 := []int{1, 2, 2, 3, 4}
	arr11 := []int{2, 3, 3, 5}
	fmt.Printf("\n배열1: %v\n배열2: %v\n", arr10, arr11)
	fmt.Printf("교집합: %v\n", Intersection(arr10, arr11))

	// 문자열
	fmt.Println("\n============================================")
	fmt.Println("[문자열 관련]")
	fmt.Println("============================================")

	str := []byte("hello")
	fmt.Printf("원본: \"%s\"\n", str)
	ReverseString(str)
	fmt.Printf("뒤집기: \"%s\"\n", str)

	str2 := "hello world"
	fmt.Printf("\n원본: \"%s\"\n", str2)
	fmt.Printf("각 단어 뒤집기: \"%s\"\n", ReverseWords(str2))

	// 패턴 요약
	fmt.Println("\n============================================")
	fmt.Println("Two Pointer 패턴 요약")
	fmt.Println("============================================")
	fmt.Println(`
┌─────────────────────────────────────────────────────────────┐
│ 반대 방향 (Opposite Direction)                               │
├─────────────────────────────────────────────────────────────┤
│ - 양 끝에서 시작하여 중앙으로 이동                           │
│ - 사용: 정렬된 배열의 합 찾기, 팰린드롬, 뒤집기              │
│ - 조건에 따라 left++, right--, 또는 둘 다                    │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ 같은 방향 (Same Direction)                                   │
├─────────────────────────────────────────────────────────────┤
│ - 둘 다 한쪽 끝에서 시작                                     │
│ - slow: 결과 위치, fast: 검사 위치                           │
│ - 사용: 중복 제거, 특정 값 제거, 배열 파티셔닝               │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ 빠른/느린 (Fast/Slow)                                        │
├─────────────────────────────────────────────────────────────┤
│ - 서로 다른 속도로 이동                                      │
│ - 사용: 사이클 감지, 중간점 찾기, 끝에서 n번째 찾기          │
│ - Floyd's Algorithm (토끼와 거북이)                          │
└─────────────────────────────────────────────────────────────┘
`)
}
