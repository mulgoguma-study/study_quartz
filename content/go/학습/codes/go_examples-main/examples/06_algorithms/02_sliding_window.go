package main

import "fmt"

/*
===========================================
6-2. 슬라이딩 윈도우 (Sliding Window)
===========================================

슬라이딩 윈도우란?
- 배열이나 문자열에서 연속적인 구간(윈도우)을 효율적으로 처리하는 기법
- 윈도우를 한 칸씩 이동하면서 필요한 값을 계산
- O(n²)를 O(n)으로 최적화 가능

종류:
1. 고정 크기 윈도우: 윈도우 크기가 k로 고정
2. 가변 크기 윈도우: 조건에 따라 윈도우 크기가 변함

사용 사례:
- 연속 부분 배열의 최대/최소 합
- 특정 조건을 만족하는 최소/최대 길이 부분 문자열
- 평균 계산, 이동 평균
*/

// ============================================
// 고정 크기 슬라이딩 윈도우
// ============================================

// MaxSumSubarray 크기 k인 연속 부분 배열의 최대 합
// 시간복잡도: O(n)
func MaxSumSubarray(arr []int, k int) int {
	n := len(arr)
	if n < k {
		return 0
	}

	// 첫 번째 윈도우의 합 계산
	windowSum := 0
	for i := 0; i < k; i++ {
		windowSum += arr[i]
	}
	maxSum := windowSum

	// 윈도우를 한 칸씩 이동
	// 새로운 요소 추가, 맨 앞 요소 제거
	for i := k; i < n; i++ {
		// 윈도우 이동:
		// - arr[i]: 새로 들어오는 요소
		// - arr[i-k]: 빠져나가는 요소
		windowSum = windowSum + arr[i] - arr[i-k]

		if windowSum > maxSum {
			maxSum = windowSum
		}
	}

	return maxSum
}

// MovingAverage 이동 평균 계산
// 각 위치에서 크기 k 윈도우의 평균 반환
func MovingAverage(arr []int, k int) []float64 {
	n := len(arr)
	if n < k {
		return nil
	}

	result := make([]float64, n-k+1)

	// 첫 번째 윈도우 합
	windowSum := 0
	for i := 0; i < k; i++ {
		windowSum += arr[i]
	}
	result[0] = float64(windowSum) / float64(k)

	// 윈도우 이동하며 평균 계산
	for i := k; i < n; i++ {
		windowSum = windowSum + arr[i] - arr[i-k]
		result[i-k+1] = float64(windowSum) / float64(k)
	}

	return result
}

// MaxInWindow 각 윈도우에서의 최대값 찾기 (deque 사용)
// 시간복잡도: O(n)
func MaxInWindow(arr []int, k int) []int {
	n := len(arr)
	if n < k {
		return nil
	}

	result := make([]int, 0, n-k+1)
	// deque: 인덱스를 저장 (내림차순 유지)
	deque := make([]int, 0, k)

	for i := 0; i < n; i++ {
		// 윈도우를 벗어난 요소 제거 (앞에서)
		if len(deque) > 0 && deque[0] <= i-k {
			deque = deque[1:]
		}

		// 현재 요소보다 작은 요소들 제거 (뒤에서)
		// 왜? 이들은 절대 최대값이 될 수 없음
		for len(deque) > 0 && arr[deque[len(deque)-1]] < arr[i] {
			deque = deque[:len(deque)-1]
		}

		// 현재 인덱스 추가
		deque = append(deque, i)

		// 첫 번째 윈도우가 완성된 후부터 결과 저장
		if i >= k-1 {
			result = append(result, arr[deque[0]])
		}
	}

	return result
}

// ============================================
// 가변 크기 슬라이딩 윈도우
// ============================================

// MinSubarrayLen 합이 target 이상인 최소 길이 부분 배열
// 시간복잡도: O(n)
func MinSubarrayLen(target int, arr []int) int {
	n := len(arr)
	minLen := n + 1 // 불가능한 큰 값으로 초기화

	left := 0
	currentSum := 0

	for right := 0; right < n; right++ {
		// 오른쪽 포인터가 가리키는 요소 추가
		currentSum += arr[right]

		// 조건을 만족하면 왼쪽 포인터 이동 (윈도우 축소)
		for currentSum >= target {
			// 현재 윈도우 길이와 최소 길이 비교
			windowLen := right - left + 1
			if windowLen < minLen {
				minLen = windowLen
			}

			// 왼쪽 요소 제거하고 포인터 이동
			currentSum -= arr[left]
			left++
		}
	}

	if minLen == n+1 {
		return 0 // 조건을 만족하는 부분 배열 없음
	}
	return minLen
}

// LongestSubstringKDistinct K개 이하의 고유 문자를 가진 최장 부분 문자열
// 시간복잡도: O(n)
func LongestSubstringKDistinct(s string, k int) int {
	if k == 0 || len(s) == 0 {
		return 0
	}

	// 문자 빈도를 저장하는 맵
	charCount := make(map[byte]int)
	maxLen := 0
	left := 0

	for right := 0; right < len(s); right++ {
		// 오른쪽 문자 추가
		charCount[s[right]]++

		// 고유 문자가 k개를 초과하면 왼쪽에서 제거
		for len(charCount) > k {
			charCount[s[left]]--
			if charCount[s[left]] == 0 {
				delete(charCount, s[left])
			}
			left++
		}

		// 최대 길이 갱신
		windowLen := right - left + 1
		if windowLen > maxLen {
			maxLen = windowLen
		}
	}

	return maxLen
}

// LongestSubstringWithoutRepeating 반복 없는 최장 부분 문자열
// 시간복잡도: O(n)
func LongestSubstringWithoutRepeating(s string) int {
	// 각 문자의 마지막 등장 위치
	lastIndex := make(map[byte]int)
	maxLen := 0
	left := 0

	for right := 0; right < len(s); right++ {
		char := s[right]

		// 이미 등장한 문자이고, 현재 윈도우 내에 있다면
		if idx, exists := lastIndex[char]; exists && idx >= left {
			// 왼쪽 포인터를 중복 문자 다음으로 이동
			left = idx + 1
		}

		// 현재 문자 위치 기록
		lastIndex[char] = right

		// 최대 길이 갱신
		windowLen := right - left + 1
		if windowLen > maxLen {
			maxLen = windowLen
		}
	}

	return maxLen
}

// ============================================
// 실용적인 예제: 문자열 아나그램 찾기
// ============================================

// FindAnagrams 문자열 s에서 p의 아나그램 시작 인덱스 찾기
// 시간복잡도: O(n)
func FindAnagrams(s, p string) []int {
	result := []int{}
	if len(s) < len(p) {
		return result
	}

	// p의 문자 빈도 계산
	pCount := make(map[byte]int)
	for i := 0; i < len(p); i++ {
		pCount[p[i]]++
	}

	// 윈도우의 문자 빈도
	windowCount := make(map[byte]int)

	// 첫 번째 윈도우 초기화
	for i := 0; i < len(p); i++ {
		windowCount[s[i]]++
	}

	// 빈도가 같은지 비교하는 헬퍼 함수
	isMatch := func() bool {
		if len(windowCount) != len(pCount) {
			return false
		}
		for k, v := range pCount {
			if windowCount[k] != v {
				return false
			}
		}
		return true
	}

	if isMatch() {
		result = append(result, 0)
	}

	// 윈도우 이동
	for i := len(p); i < len(s); i++ {
		// 새 문자 추가
		windowCount[s[i]]++

		// 맨 앞 문자 제거
		oldChar := s[i-len(p)]
		windowCount[oldChar]--
		if windowCount[oldChar] == 0 {
			delete(windowCount, oldChar)
		}

		// 아나그램인지 확인
		if isMatch() {
			result = append(result, i-len(p)+1)
		}
	}

	return result
}

// ============================================
// 실용적인 예제: 최소 윈도우 부분 문자열
// ============================================

// MinWindow s에서 t의 모든 문자를 포함하는 최소 부분 문자열
// 시간복잡도: O(n)
func MinWindow(s, t string) string {
	if len(s) == 0 || len(t) == 0 {
		return ""
	}

	// t의 문자 빈도
	tCount := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		tCount[t[i]]++
	}

	// 필요한 고유 문자 수
	required := len(tCount)

	// 윈도우 내 문자 빈도
	windowCount := make(map[byte]int)

	// 조건을 만족하는 고유 문자 수
	formed := 0

	// 결과: [윈도우 길이, 시작 인덱스, 끝 인덱스]
	minLen := len(s) + 1
	resultLeft, resultRight := 0, 0

	left := 0
	for right := 0; right < len(s); right++ {
		// 오른쪽 문자 추가
		char := s[right]
		windowCount[char]++

		// t에 있는 문자이고, 필요한 만큼 모았다면
		if count, exists := tCount[char]; exists && windowCount[char] == count {
			formed++
		}

		// 모든 조건 만족 시 왼쪽에서 축소 시도
		for left <= right && formed == required {
			char := s[left]

			// 최소 길이 갱신
			if right-left+1 < minLen {
				minLen = right - left + 1
				resultLeft = left
				resultRight = right + 1
			}

			// 왼쪽 문자 제거
			windowCount[char]--
			if count, exists := tCount[char]; exists && windowCount[char] < count {
				formed--
			}
			left++
		}
	}

	if minLen == len(s)+1 {
		return ""
	}
	return s[resultLeft:resultRight]
}

func main() {
	fmt.Println("============================================")
	fmt.Println("6-2. 슬라이딩 윈도우 (Sliding Window)")
	fmt.Println("============================================")

	// 고정 크기 슬라이딩 윈도우 예제
	fmt.Println("\n[고정 크기 슬라이딩 윈도우]")
	fmt.Println("--------------------------------------------")

	arr := []int{2, 1, 5, 1, 3, 2}
	k := 3
	fmt.Printf("배열: %v, k=%d\n", arr, k)
	fmt.Printf("크기 %d인 부분 배열의 최대 합: %d\n", k, MaxSumSubarray(arr, k))

	// 이동 평균
	fmt.Println("\n[이동 평균]")
	arr2 := []int{1, 3, 5, 7, 9, 11}
	k2 := 3
	fmt.Printf("배열: %v, k=%d\n", arr2, k2)
	fmt.Printf("이동 평균: %v\n", MovingAverage(arr2, k2))

	// 각 윈도우의 최대값
	fmt.Println("\n[각 윈도우의 최대값]")
	arr3 := []int{1, 3, -1, -3, 5, 3, 6, 7}
	k3 := 3
	fmt.Printf("배열: %v, k=%d\n", arr3, k3)
	fmt.Printf("각 윈도우 최대값: %v\n", MaxInWindow(arr3, k3))

	// 가변 크기 슬라이딩 윈도우 예제
	fmt.Println("\n============================================")
	fmt.Println("[가변 크기 슬라이딩 윈도우]")
	fmt.Println("============================================")

	// 합이 target 이상인 최소 길이 부분 배열
	arr4 := []int{2, 3, 1, 2, 4, 3}
	target := 7
	fmt.Printf("\n배열: %v, target=%d\n", arr4, target)
	fmt.Printf("합이 %d 이상인 최소 길이 부분 배열: %d\n", target, MinSubarrayLen(target, arr4))

	// K개 이하 고유 문자 최장 부분 문자열
	str1 := "eceba"
	k4 := 2
	fmt.Printf("\n문자열: \"%s\", k=%d\n", str1, k4)
	fmt.Printf("%d개 이하 고유 문자 최장 부분 문자열 길이: %d\n", k4, LongestSubstringKDistinct(str1, k4))

	// 반복 없는 최장 부분 문자열
	str2 := "abcabcbb"
	fmt.Printf("\n문자열: \"%s\"\n", str2)
	fmt.Printf("반복 없는 최장 부분 문자열 길이: %d\n", LongestSubstringWithoutRepeating(str2))

	str3 := "bbbbb"
	fmt.Printf("\n문자열: \"%s\"\n", str3)
	fmt.Printf("반복 없는 최장 부분 문자열 길이: %d\n", LongestSubstringWithoutRepeating(str3))

	// 아나그램 찾기
	fmt.Println("\n============================================")
	fmt.Println("[아나그램 찾기]")
	fmt.Println("============================================")

	s := "cbaebabacd"
	p := "abc"
	fmt.Printf("문자열 s: \"%s\", 패턴 p: \"%s\"\n", s, p)
	fmt.Printf("아나그램 시작 인덱스: %v\n", FindAnagrams(s, p))

	// 최소 윈도우 부분 문자열
	fmt.Println("\n============================================")
	fmt.Println("[최소 윈도우 부분 문자열]")
	fmt.Println("============================================")

	s2 := "ADOBECODEBANC"
	t := "ABC"
	fmt.Printf("문자열 s: \"%s\", 대상 t: \"%s\"\n", s2, t)
	fmt.Printf("모든 문자를 포함하는 최소 부분 문자열: \"%s\"\n", MinWindow(s2, t))

	// 슬라이딩 윈도우 패턴 요약
	fmt.Println("\n============================================")
	fmt.Println("슬라이딩 윈도우 패턴 요약")
	fmt.Println("============================================")
	fmt.Println(`
┌─────────────────────────────────────────────────────────────┐
│ 고정 크기 윈도우 패턴                                        │
├─────────────────────────────────────────────────────────────┤
│ 1. 첫 k개 요소로 초기 윈도우 생성                            │
│ 2. 한 칸씩 이동: 새 요소 추가 + 맨 앞 요소 제거              │
│ 3. 각 위치에서 필요한 연산 수행                              │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ 가변 크기 윈도우 패턴                                        │
├─────────────────────────────────────────────────────────────┤
│ 1. right 포인터로 윈도우 확장                                │
│ 2. 조건 만족 시 left 포인터로 윈도우 축소                    │
│ 3. 최적값(최대/최소 길이) 갱신                               │
└─────────────────────────────────────────────────────────────┘

[사용 시나리오]
- "연속된 k개의 합/평균/최대" → 고정 크기 윈도우
- "조건을 만족하는 최소/최대 길이" → 가변 크기 윈도우
- "중복 없이/K개 이하로" → 해시맵 + 가변 윈도우
`)
}
