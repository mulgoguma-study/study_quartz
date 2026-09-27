package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
===========================================
6-3. HashMap 빈도 계산 (Frequency Counting)
===========================================

HashMap 빈도 계산이란?
- 해시맵(Go에서는 map)을 사용하여 요소의 출현 빈도를 계산
- O(1) 시간에 삽입, 조회, 갱신 가능
- 다양한 문제에서 핵심 패턴으로 사용

사용 사례:
- 배열/문자열에서 중복 요소 찾기
- 특정 조건을 만족하는 요소 쌍 찾기
- 아나그램 판별
- 가장 빈번한/드문 요소 찾기
*/

// ============================================
// 기본 빈도 계산
// ============================================

// CountFrequency 배열에서 각 요소의 빈도 계산
// 시간복잡도: O(n), 공간복잡도: O(n)
func CountFrequency(arr []int) map[int]int {
	freq := make(map[int]int)

	for _, num := range arr {
		// 키가 없으면 0으로 초기화되므로 바로 증가 가능
		freq[num]++
	}

	return freq
}

// CountCharFrequency 문자열에서 문자 빈도 계산
func CountCharFrequency(s string) map[rune]int {
	freq := make(map[rune]int)

	// range로 문자열 순회 시 rune(유니코드 코드포인트) 반환
	for _, char := range s {
		freq[char]++
	}

	return freq
}

// ============================================
// 중복 찾기
// ============================================

// FindDuplicates 중복된 요소들 반환
func FindDuplicates(arr []int) []int {
	freq := make(map[int]int)
	duplicates := []int{}

	for _, num := range arr {
		freq[num]++
		// 정확히 2번 등장할 때만 결과에 추가 (중복 방지)
		if freq[num] == 2 {
			duplicates = append(duplicates, num)
		}
	}

	return duplicates
}

// FindFirstDuplicate 첫 번째 중복 요소 찾기
func FindFirstDuplicate(arr []int) (int, bool) {
	seen := make(map[int]bool)

	for _, num := range arr {
		if seen[num] {
			return num, true
		}
		seen[num] = true
	}

	return 0, false
}

// FindUnique 고유한 요소 찾기 (한 번만 등장)
func FindUnique(arr []int) []int {
	freq := make(map[int]int)

	for _, num := range arr {
		freq[num]++
	}

	unique := []int{}
	for _, num := range arr {
		if freq[num] == 1 {
			unique = append(unique, num)
		}
	}

	return unique
}

// ============================================
// 가장 빈번한/드문 요소
// ============================================

// MostFrequent 가장 빈번한 요소 찾기
func MostFrequent(arr []int) (int, int) {
	freq := make(map[int]int)

	for _, num := range arr {
		freq[num]++
	}

	maxCount := 0
	mostFreq := 0

	for num, count := range freq {
		if count > maxCount {
			maxCount = count
			mostFreq = num
		}
	}

	return mostFreq, maxCount
}

// TopKFrequent 가장 빈번한 k개 요소 (버킷 정렬 활용)
// 시간복잡도: O(n)
func TopKFrequent(arr []int, k int) []int {
	// 1. 빈도 계산
	freq := make(map[int]int)
	for _, num := range arr {
		freq[num]++
	}

	// 2. 버킷 생성 (인덱스 = 빈도, 값 = 해당 빈도를 가진 숫자들)
	// 최대 빈도는 배열 길이를 넘지 않음
	buckets := make([][]int, len(arr)+1)
	for i := range buckets {
		buckets[i] = []int{}
	}

	for num, count := range freq {
		buckets[count] = append(buckets[count], num)
	}

	// 3. 높은 빈도부터 k개 수집
	result := []int{}
	for i := len(buckets) - 1; i >= 0 && len(result) < k; i-- {
		result = append(result, buckets[i]...)
	}

	// k개만 반환
	if len(result) > k {
		result = result[:k]
	}

	return result
}

// LeastFrequent 가장 드문 요소 찾기
func LeastFrequent(arr []int) (int, int) {
	freq := make(map[int]int)

	for _, num := range arr {
		freq[num]++
	}

	minCount := len(arr) + 1
	leastFreq := 0

	for num, count := range freq {
		if count < minCount {
			minCount = count
			leastFreq = num
		}
	}

	return leastFreq, minCount
}

// ============================================
// 아나그램 관련
// ============================================

// IsAnagram 두 문자열이 아나그램인지 확인
// 시간복잡도: O(n)
func IsAnagram(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	freq := make(map[rune]int)

	// s1의 문자 빈도 증가
	for _, char := range s1 {
		freq[char]++
	}

	// s2의 문자 빈도 감소
	for _, char := range s2 {
		freq[char]--
		// 빈도가 음수가 되면 아나그램 아님
		if freq[char] < 0 {
			return false
		}
	}

	return true
}

// GroupAnagrams 아나그램끼리 그룹화
// 시간복잡도: O(n * k log k) - k는 문자열 최대 길이
func GroupAnagrams(strs []string) [][]string {
	// 정렬된 문자열을 키로 사용
	groups := make(map[string][]string)

	for _, s := range strs {
		// 문자열을 정렬하여 키 생성
		chars := []rune(s)
		sort.Slice(chars, func(i, j int) bool {
			return chars[i] < chars[j]
		})
		key := string(chars)

		groups[key] = append(groups[key], s)
	}

	// 맵의 값들을 슬라이스로 변환
	result := make([][]string, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}

	return result
}

// GroupAnagramsOptimized 아나그램 그룹화 (빈도 배열 키 사용)
// 시간복잡도: O(n * k) - 정렬 없이 더 효율적
func GroupAnagramsOptimized(strs []string) [][]string {
	groups := make(map[string][]string)

	for _, s := range strs {
		// 빈도 배열을 키로 사용 (26개 알파벳 가정)
		key := getFrequencyKey(s)
		groups[key] = append(groups[key], s)
	}

	result := make([][]string, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}

	return result
}

// getFrequencyKey 문자 빈도를 문자열 키로 변환
func getFrequencyKey(s string) string {
	freq := make([]int, 26)
	for _, char := range s {
		if char >= 'a' && char <= 'z' {
			freq[char-'a']++
		}
	}

	// 빈도 배열을 문자열로 변환
	var sb strings.Builder
	for i, count := range freq {
		if count > 0 {
			sb.WriteString(fmt.Sprintf("%c%d", 'a'+i, count))
		}
	}
	return sb.String()
}

// ============================================
// Two Sum 관련 (HashMap 활용)
// ============================================

// TwoSum 합이 target인 두 요소의 인덱스 찾기
// 시간복잡도: O(n)
func TwoSum(arr []int, target int) (int, int, bool) {
	// 값 -> 인덱스 매핑
	seen := make(map[int]int)

	for i, num := range arr {
		// 필요한 보완값(complement) 계산
		complement := target - num

		// 보완값이 이미 존재하면 쌍 발견
		if j, exists := seen[complement]; exists {
			return j, i, true
		}

		// 현재 값 저장
		seen[num] = i
	}

	return 0, 0, false
}

// TwoSumAllPairs 합이 target인 모든 쌍 찾기
func TwoSumAllPairs(arr []int, target int) [][2]int {
	freq := make(map[int]int)
	result := [][2]int{}
	used := make(map[int]bool)

	// 먼저 빈도 계산
	for _, num := range arr {
		freq[num]++
	}

	for num := range freq {
		complement := target - num

		// 이미 처리한 쌍은 건너뜀
		if used[num] {
			continue
		}

		if complement == num {
			// 같은 숫자 쌍: 2개 이상 있어야 함
			if freq[num] >= 2 {
				result = append(result, [2]int{num, num})
			}
		} else if _, exists := freq[complement]; exists {
			result = append(result, [2]int{num, complement})
			used[complement] = true
		}
		used[num] = true
	}

	return result
}

// ============================================
// 문자열 비교 및 검사
// ============================================

// CanConstruct ransomNote를 magazine 문자로 만들 수 있는지 확인
func CanConstruct(ransomNote, magazine string) bool {
	freq := make(map[rune]int)

	// magazine의 문자 빈도
	for _, char := range magazine {
		freq[char]++
	}

	// ransomNote에 필요한 문자 확인
	for _, char := range ransomNote {
		freq[char]--
		if freq[char] < 0 {
			return false
		}
	}

	return true
}

// FirstUniqueChar 첫 번째 고유 문자 인덱스
func FirstUniqueChar(s string) int {
	freq := make(map[rune]int)

	// 빈도 계산
	for _, char := range s {
		freq[char]++
	}

	// 순서대로 순회하여 첫 고유 문자 찾기
	for i, char := range s {
		if freq[char] == 1 {
			return i
		}
	}

	return -1
}

// ============================================
// 빈도 기반 정렬
// ============================================

// SortByFrequency 빈도에 따라 정렬 (높은 빈도 먼저)
func SortByFrequency(arr []int) []int {
	freq := make(map[int]int)
	for _, num := range arr {
		freq[num]++
	}

	result := make([]int, len(arr))
	copy(result, arr)

	// 빈도 기준 내림차순, 같으면 값 기준 오름차순
	sort.Slice(result, func(i, j int) bool {
		if freq[result[i]] == freq[result[j]] {
			return result[i] < result[j]
		}
		return freq[result[i]] > freq[result[j]]
	})

	return result
}

// ============================================
// 부분 배열/문자열 문제
// ============================================

// SubarraySumEqualsK 합이 k인 부분 배열 개수
// 시간복잡도: O(n)
func SubarraySumEqualsK(arr []int, k int) int {
	// prefixSum -> 빈도 맵
	prefixCount := make(map[int]int)
	prefixCount[0] = 1 // 빈 부분 배열 (합 0)

	count := 0
	prefixSum := 0

	for _, num := range arr {
		prefixSum += num

		// prefixSum - k가 이전에 등장했다면
		// 그 지점부터 현재까지의 합이 k
		if freq, exists := prefixCount[prefixSum-k]; exists {
			count += freq
		}

		prefixCount[prefixSum]++
	}

	return count
}

// LongestSubarrayWithSumK 합이 k인 가장 긴 부분 배열 길이
func LongestSubarrayWithSumK(arr []int, k int) int {
	// prefixSum -> 첫 등장 인덱스
	prefixIndex := make(map[int]int)
	prefixIndex[0] = -1 // 인덱스 -1에서 시작하는 빈 배열

	maxLen := 0
	prefixSum := 0

	for i, num := range arr {
		prefixSum += num

		// prefixSum - k의 첫 등장 위치 찾기
		if startIdx, exists := prefixIndex[prefixSum-k]; exists {
			length := i - startIdx
			if length > maxLen {
				maxLen = length
			}
		}

		// 첫 등장 인덱스만 저장 (최대 길이를 위해)
		if _, exists := prefixIndex[prefixSum]; !exists {
			prefixIndex[prefixSum] = i
		}
	}

	return maxLen
}

func main() {
	fmt.Println("============================================")
	fmt.Println("6-3. HashMap 빈도 계산")
	fmt.Println("============================================")

	// 기본 빈도 계산
	fmt.Println("\n[기본 빈도 계산]")
	fmt.Println("--------------------------------------------")
	arr := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	fmt.Printf("배열: %v\n", arr)
	freq := CountFrequency(arr)
	fmt.Println("빈도:")
	for num, count := range freq {
		fmt.Printf("  %d: %d번\n", num, count)
	}

	// 문자열 빈도
	fmt.Println("\n[문자열 문자 빈도]")
	str := "hello"
	fmt.Printf("문자열: \"%s\"\n", str)
	charFreq := CountCharFrequency(str)
	fmt.Println("빈도:")
	for char, count := range charFreq {
		fmt.Printf("  '%c': %d번\n", char, count)
	}

	// 중복 찾기
	fmt.Println("\n============================================")
	fmt.Println("[중복 요소 찾기]")
	fmt.Println("============================================")
	arr2 := []int{4, 3, 2, 7, 8, 2, 3, 1}
	fmt.Printf("배열: %v\n", arr2)
	fmt.Printf("중복된 요소: %v\n", FindDuplicates(arr2))

	if first, found := FindFirstDuplicate(arr2); found {
		fmt.Printf("첫 번째 중복: %d\n", first)
	}

	fmt.Printf("고유한 요소 (1번만 등장): %v\n", FindUnique(arr2))

	// 빈번한/드문 요소
	fmt.Println("\n============================================")
	fmt.Println("[가장 빈번한/드문 요소]")
	fmt.Println("============================================")
	arr3 := []int{1, 1, 1, 2, 2, 3}
	fmt.Printf("배열: %v\n", arr3)
	mostFreq, count := MostFrequent(arr3)
	fmt.Printf("가장 빈번한 요소: %d (등장: %d번)\n", mostFreq, count)

	leastFreq, countLeast := LeastFrequent(arr3)
	fmt.Printf("가장 드문 요소: %d (등장: %d번)\n", leastFreq, countLeast)

	// Top K 빈번한 요소
	fmt.Println("\n[Top K 빈번한 요소]")
	arr4 := []int{1, 1, 1, 2, 2, 3, 4, 4, 4, 4}
	k := 2
	fmt.Printf("배열: %v, k=%d\n", arr4, k)
	fmt.Printf("Top %d 빈번한 요소: %v\n", k, TopKFrequent(arr4, k))

	// 아나그램
	fmt.Println("\n============================================")
	fmt.Println("[아나그램]")
	fmt.Println("============================================")
	s1, s2 := "anagram", "nagaram"
	fmt.Printf("\"%s\"와 \"%s\"는 아나그램? %v\n", s1, s2, IsAnagram(s1, s2))

	s3, s4 := "rat", "car"
	fmt.Printf("\"%s\"와 \"%s\"는 아나그램? %v\n", s3, s4, IsAnagram(s3, s4))

	// 아나그램 그룹화
	fmt.Println("\n[아나그램 그룹화]")
	strs := []string{"eat", "tea", "tan", "ate", "nat", "bat"}
	fmt.Printf("입력: %v\n", strs)
	fmt.Println("그룹:")
	for _, group := range GroupAnagrams(strs) {
		fmt.Printf("  %v\n", group)
	}

	// Two Sum
	fmt.Println("\n============================================")
	fmt.Println("[Two Sum]")
	fmt.Println("============================================")
	arr5 := []int{2, 7, 11, 15}
	target := 9
	fmt.Printf("배열: %v, target=%d\n", arr5, target)
	if i, j, found := TwoSum(arr5, target); found {
		fmt.Printf("인덱스: [%d, %d] (값: %d + %d = %d)\n", i, j, arr5[i], arr5[j], target)
	}

	// 모든 쌍 찾기
	arr6 := []int{1, 2, 3, 4, 5, 6}
	target2 := 7
	fmt.Printf("\n배열: %v, target=%d\n", arr6, target2)
	fmt.Printf("합이 %d인 모든 쌍: %v\n", target2, TwoSumAllPairs(arr6, target2))

	// 문자열 문제
	fmt.Println("\n============================================")
	fmt.Println("[문자열 문제]")
	fmt.Println("============================================")

	ransomNote := "aa"
	magazine := "aab"
	fmt.Printf("RansomNote: \"%s\", Magazine: \"%s\"\n", ransomNote, magazine)
	fmt.Printf("만들 수 있나? %v\n", CanConstruct(ransomNote, magazine))

	str5 := "loveleetcode"
	fmt.Printf("\n문자열: \"%s\"\n", str5)
	fmt.Printf("첫 번째 고유 문자 인덱스: %d\n", FirstUniqueChar(str5))

	// 빈도 기반 정렬
	fmt.Println("\n============================================")
	fmt.Println("[빈도 기반 정렬]")
	fmt.Println("============================================")
	arr7 := []int{1, 1, 2, 2, 2, 3}
	fmt.Printf("원본: %v\n", arr7)
	fmt.Printf("빈도순 정렬: %v\n", SortByFrequency(arr7))

	// 부분 배열 합
	fmt.Println("\n============================================")
	fmt.Println("[부분 배열 합 문제]")
	fmt.Println("============================================")
	arr8 := []int{1, 1, 1}
	k2 := 2
	fmt.Printf("배열: %v, k=%d\n", arr8, k2)
	fmt.Printf("합이 %d인 부분 배열 개수: %d\n", k2, SubarraySumEqualsK(arr8, k2))

	arr9 := []int{1, -1, 5, -2, 3}
	k3 := 3
	fmt.Printf("\n배열: %v, k=%d\n", arr9, k3)
	fmt.Printf("합이 %d인 가장 긴 부분 배열 길이: %d\n", k3, LongestSubarrayWithSumK(arr9, k3))

	// 패턴 요약
	fmt.Println("\n============================================")
	fmt.Println("HashMap 빈도 계산 패턴 요약")
	fmt.Println("============================================")
	fmt.Println(`
┌─────────────────────────────────────────────────────────────┐
│ HashMap 활용 패턴                                            │
├─────────────────────────────────────────────────────────────┤
│ 1. 빈도 계산: map[요소]int로 등장 횟수 카운트               │
│ 2. 존재 여부: map[요소]bool로 빠른 O(1) 조회                │
│ 3. 인덱스 저장: map[값]인덱스로 위치 기억                   │
│ 4. 그룹화: map[키][]값으로 같은 특성 요소 그룹화            │
│ 5. 누적합: map[prefixSum]빈도로 부분합 문제 해결            │
└─────────────────────────────────────────────────────────────┘

[시간복잡도]
- 삽입/조회/삭제: O(1) 평균
- 전체 요소 순회: O(n)

[주의사항]
- map 순회 순서는 무작위 (순서 보장 X)
- 존재하지 않는 키 접근 시 제로값 반환
- 동시성 환경에서는 sync.Map 사용 고려
`)
}
