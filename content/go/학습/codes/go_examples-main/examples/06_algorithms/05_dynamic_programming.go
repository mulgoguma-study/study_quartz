package main

import "fmt"

/*
===========================================
6-5. 동적 프로그래밍 (Dynamic Programming)
===========================================

동적 프로그래밍이란?
- 큰 문제를 작은 하위 문제로 나누어 해결하는 알고리즘 기법
- 중복되는 하위 문제의 결과를 저장(메모이제이션)하여 재사용
- 최적 부분 구조 + 중복 부분 문제 성질을 가진 문제에 적용

접근 방식:
1. Top-Down (메모이제이션): 재귀 + 캐싱
2. Bottom-Up (타뷸레이션): 반복문 + 테이블

DP 문제 해결 단계:
1. 상태(state) 정의: dp[i]가 무엇을 의미하는지 정의
2. 점화식(recurrence) 도출: dp[i]와 이전 상태들의 관계
3. 기저 조건(base case) 설정: 가장 작은 문제의 답
4. 계산 순서 결정: 어떤 순서로 채워나갈지
*/

// ============================================
// 1차원 DP
// ============================================

// Fibonacci 피보나치 수열 (가장 기본적인 DP)
// Top-Down 방식 (메모이제이션)
func FibonacciTopDown(n int) int {
	// 메모이제이션을 위한 캐시
	memo := make(map[int]int)

	var fib func(n int) int
	fib = func(n int) int {
		// 기저 조건
		if n <= 1 {
			return n
		}

		// 이미 계산된 값이면 바로 반환
		if val, exists := memo[n]; exists {
			return val
		}

		// 계산 후 저장
		memo[n] = fib(n-1) + fib(n-2)
		return memo[n]
	}

	return fib(n)
}

// FibonacciBottomUp 피보나치 Bottom-Up 방식 (타뷸레이션)
// 시간복잡도: O(n), 공간복잡도: O(n)
func FibonacciBottomUp(n int) int {
	if n <= 1 {
		return n
	}

	// dp[i] = i번째 피보나치 수
	dp := make([]int, n+1)
	dp[0] = 0
	dp[1] = 1

	// 작은 문제부터 순서대로 해결
	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}

	return dp[n]
}

// FibonacciOptimized 공간 최적화된 피보나치
// 공간복잡도: O(1)
func FibonacciOptimized(n int) int {
	if n <= 1 {
		return n
	}

	prev2, prev1 := 0, 1
	for i := 2; i <= n; i++ {
		curr := prev1 + prev2
		prev2 = prev1
		prev1 = curr
	}

	return prev1
}

// ClimbingStairs 계단 오르기
// 한 번에 1계단 또는 2계단을 오를 수 있을 때, n계단을 오르는 방법의 수
// dp[i] = i계단에 도달하는 방법의 수
func ClimbingStairs(n int) int {
	if n <= 2 {
		return n
	}

	dp := make([]int, n+1)
	dp[1] = 1 // 1계단: 1가지 방법
	dp[2] = 2 // 2계단: 2가지 방법 (1+1, 2)

	// 점화식: dp[i] = dp[i-1] + dp[i-2]
	// i번째 계단은 i-1에서 1칸 오르거나 i-2에서 2칸 오르면 됨
	for i := 3; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}

	return dp[n]
}

// HouseRobber 집 털기 문제
// 인접한 집은 털 수 없을 때 최대 금액
// dp[i] = i번째 집까지 고려했을 때 최대 금액
func HouseRobber(houses []int) int {
	n := len(houses)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return houses[0]
	}

	dp := make([]int, n)
	dp[0] = houses[0]
	dp[1] = max(houses[0], houses[1])

	// 점화식: dp[i] = max(dp[i-1], dp[i-2] + houses[i])
	// i번째 집을 털거나(dp[i-2] + houses[i]), 안 털거나(dp[i-1])
	for i := 2; i < n; i++ {
		dp[i] = max(dp[i-1], dp[i-2]+houses[i])
	}

	return dp[n-1]
}

// MaxSubarraySum 최대 부분 배열 합 (Kadane's Algorithm)
// dp[i] = i번째 요소로 끝나는 최대 부분 배열 합
func MaxSubarraySum(arr []int) int {
	if len(arr) == 0 {
		return 0
	}

	// 현재 위치에서 끝나는 최대 합
	currentMax := arr[0]
	// 전체 최대 합
	globalMax := arr[0]

	for i := 1; i < len(arr); i++ {
		// 점화식: currentMax[i] = max(arr[i], currentMax[i-1] + arr[i])
		// 이전 합에 현재 요소를 더하거나, 현재 요소에서 새로 시작
		currentMax = max(arr[i], currentMax+arr[i])
		globalMax = max(globalMax, currentMax)
	}

	return globalMax
}

// CoinChange 동전 교환 문제
// 최소 동전 개수로 금액을 만들기
// dp[i] = 금액 i를 만드는 데 필요한 최소 동전 개수
func CoinChange(coins []int, amount int) int {
	// 불가능한 값으로 초기화 (amount+1은 절대 나올 수 없는 큰 값)
	dp := make([]int, amount+1)
	for i := range dp {
		dp[i] = amount + 1
	}
	dp[0] = 0 // 0원을 만드는 데 0개 필요

	for i := 1; i <= amount; i++ {
		for _, coin := range coins {
			// 현재 동전을 사용할 수 있으면
			if coin <= i {
				// 점화식: dp[i] = min(dp[i], dp[i-coin] + 1)
				dp[i] = min(dp[i], dp[i-coin]+1)
			}
		}
	}

	if dp[amount] > amount {
		return -1 // 불가능
	}
	return dp[amount]
}

// LongestIncreasingSubsequence 최장 증가 부분 수열 (LIS)
// dp[i] = i번째 요소로 끝나는 LIS 길이
// 시간복잡도: O(n²)
func LongestIncreasingSubsequence(arr []int) int {
	n := len(arr)
	if n == 0 {
		return 0
	}

	// 모든 dp[i]를 1로 초기화 (자기 자신만 포함)
	dp := make([]int, n)
	for i := range dp {
		dp[i] = 1
	}

	maxLen := 1

	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			// arr[j] < arr[i]이면 arr[j] 뒤에 arr[i]를 붙일 수 있음
			if arr[j] < arr[i] {
				dp[i] = max(dp[i], dp[j]+1)
			}
		}
		maxLen = max(maxLen, dp[i])
	}

	return maxLen
}

// ============================================
// 2차원 DP
// ============================================

// UniquePaths 고유 경로 수
// m x n 격자에서 좌상단에서 우하단까지 가는 경로 수
// 오른쪽 또는 아래로만 이동 가능
// dp[i][j] = (i, j)에 도달하는 경로 수
func UniquePaths(m, n int) int {
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
	}

	// 첫 행은 모두 1 (오른쪽으로만 이동)
	for j := 0; j < n; j++ {
		dp[0][j] = 1
	}
	// 첫 열은 모두 1 (아래로만 이동)
	for i := 0; i < m; i++ {
		dp[i][0] = 1
	}

	// 점화식: dp[i][j] = dp[i-1][j] + dp[i][j-1]
	// 위에서 오거나 왼쪽에서 오거나
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[i][j] = dp[i-1][j] + dp[i][j-1]
		}
	}

	return dp[m-1][n-1]
}

// MinPathSum 최소 경로 합
// grid의 좌상단에서 우하단까지 경로의 최소 합
// dp[i][j] = (i, j)까지의 최소 경로 합
func MinPathSum(grid [][]int) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}

	m, n := len(grid), len(grid[0])
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
	}

	dp[0][0] = grid[0][0]

	// 첫 행 초기화
	for j := 1; j < n; j++ {
		dp[0][j] = dp[0][j-1] + grid[0][j]
	}
	// 첫 열 초기화
	for i := 1; i < m; i++ {
		dp[i][0] = dp[i-1][0] + grid[i][0]
	}

	// 점화식: dp[i][j] = grid[i][j] + min(dp[i-1][j], dp[i][j-1])
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[i][j] = grid[i][j] + min(dp[i-1][j], dp[i][j-1])
		}
	}

	return dp[m-1][n-1]
}

// Knapsack01 0/1 배낭 문제
// 각 물건을 최대 1개만 선택할 수 있을 때 최대 가치
// dp[i][w] = i번째 물건까지 고려하고 용량 w일 때 최대 가치
func Knapsack01(weights, values []int, capacity int) int {
	n := len(weights)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}

	for i := 1; i <= n; i++ {
		for w := 0; w <= capacity; w++ {
			// i번째 물건을 넣지 않는 경우
			dp[i][w] = dp[i-1][w]

			// i번째 물건을 넣는 경우 (용량이 충분할 때만)
			if weights[i-1] <= w {
				dp[i][w] = max(dp[i][w], dp[i-1][w-weights[i-1]]+values[i-1])
			}
		}
	}

	return dp[n][capacity]
}

// Knapsack01Optimized 공간 최적화된 0/1 배낭
// 2차원 -> 1차원으로 최적화
func Knapsack01Optimized(weights, values []int, capacity int) int {
	dp := make([]int, capacity+1)

	for i := 0; i < len(weights); i++ {
		// 역순으로 순회해야 이전 결과를 덮어쓰지 않음
		for w := capacity; w >= weights[i]; w-- {
			dp[w] = max(dp[w], dp[w-weights[i]]+values[i])
		}
	}

	return dp[capacity]
}

// LongestCommonSubsequence 최장 공통 부분 수열 (LCS)
// dp[i][j] = s1[0..i-1]과 s2[0..j-1]의 LCS 길이
func LongestCommonSubsequence(s1, s2 string) int {
	m, n := len(s1), len(s2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				// 마지막 문자가 같으면 LCS에 포함
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				// 다르면 둘 중 하나를 제외
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	return dp[m][n]
}

// EditDistance 편집 거리 (Levenshtein Distance)
// s1을 s2로 변환하는 데 필요한 최소 연산 수 (삽입, 삭제, 교체)
// dp[i][j] = s1[0..i-1]을 s2[0..j-1]로 변환하는 최소 연산 수
func EditDistance(s1, s2 string) int {
	m, n := len(s1), len(s2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// 기저 조건: 빈 문자열로/에서 변환
	for i := 0; i <= m; i++ {
		dp[i][0] = i // s1의 모든 문자 삭제
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j // s2의 모든 문자 삽입
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				// 같으면 연산 불필요
				dp[i][j] = dp[i-1][j-1]
			} else {
				// 다르면 삽입, 삭제, 교체 중 최소
				dp[i][j] = 1 + min(
					dp[i-1][j],   // 삭제
					dp[i][j-1],   // 삽입
					dp[i-1][j-1], // 교체
				)
			}
		}
	}

	return dp[m][n]
}

// ============================================
// 응용 문제
// ============================================

// DecodeWays 숫자 문자열을 알파벳으로 해석하는 방법 수
// 'A' = 1, 'B' = 2, ..., 'Z' = 26
func DecodeWays(s string) int {
	if len(s) == 0 || s[0] == '0' {
		return 0
	}

	n := len(s)
	dp := make([]int, n+1)
	dp[0] = 1 // 빈 문자열
	dp[1] = 1 // 첫 문자 (0이 아닌 경우)

	for i := 2; i <= n; i++ {
		// 한 자리 숫자로 해석 (1-9)
		if s[i-1] != '0' {
			dp[i] += dp[i-1]
		}

		// 두 자리 숫자로 해석 (10-26)
		twoDigit := (s[i-2]-'0')*10 + (s[i-1] - '0')
		if twoDigit >= 10 && twoDigit <= 26 {
			dp[i] += dp[i-2]
		}
	}

	return dp[n]
}

// WordBreak 문자열을 사전의 단어들로 분리할 수 있는지
// dp[i] = s[0..i-1]을 단어들로 분리할 수 있는지
func WordBreak(s string, wordDict []string) bool {
	wordSet := make(map[string]bool)
	for _, word := range wordDict {
		wordSet[word] = true
	}

	n := len(s)
	dp := make([]bool, n+1)
	dp[0] = true // 빈 문자열

	for i := 1; i <= n; i++ {
		for j := 0; j < i; j++ {
			// s[0..j-1]이 분리 가능하고 s[j..i-1]이 사전에 있으면
			if dp[j] && wordSet[s[j:i]] {
				dp[i] = true
				break
			}
		}
	}

	return dp[n]
}

func main() {
	fmt.Println("============================================")
	fmt.Println("6-5. 동적 프로그래밍 (Dynamic Programming)")
	fmt.Println("============================================")

	// 피보나치
	fmt.Println("\n[피보나치 수열]")
	fmt.Println("--------------------------------------------")
	n := 10
	fmt.Printf("F(%d) Top-Down: %d\n", n, FibonacciTopDown(n))
	fmt.Printf("F(%d) Bottom-Up: %d\n", n, FibonacciBottomUp(n))
	fmt.Printf("F(%d) Optimized: %d\n", n, FibonacciOptimized(n))

	// 계단 오르기
	fmt.Println("\n[계단 오르기]")
	stairs := 5
	fmt.Printf("%d계단을 오르는 방법의 수: %d\n", stairs, ClimbingStairs(stairs))

	// House Robber
	fmt.Println("\n[House Robber]")
	houses := []int{2, 7, 9, 3, 1}
	fmt.Printf("집: %v\n", houses)
	fmt.Printf("최대 금액: %d\n", HouseRobber(houses))

	// 최대 부분 배열 합
	fmt.Println("\n[최대 부분 배열 합 (Kadane)]")
	arr := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	fmt.Printf("배열: %v\n", arr)
	fmt.Printf("최대 합: %d\n", MaxSubarraySum(arr))

	// 동전 교환
	fmt.Println("\n[동전 교환]")
	coins := []int{1, 2, 5}
	amount := 11
	fmt.Printf("동전: %v, 금액: %d\n", coins, amount)
	fmt.Printf("최소 동전 개수: %d\n", CoinChange(coins, amount))

	// LIS
	fmt.Println("\n[최장 증가 부분 수열 (LIS)]")
	arr2 := []int{10, 9, 2, 5, 3, 7, 101, 18}
	fmt.Printf("배열: %v\n", arr2)
	fmt.Printf("LIS 길이: %d\n", LongestIncreasingSubsequence(arr2))

	// 2차원 DP
	fmt.Println("\n============================================")
	fmt.Println("[2차원 DP]")
	fmt.Println("============================================")

	// 고유 경로
	fmt.Println("\n[고유 경로]")
	m, n2 := 3, 7
	fmt.Printf("%dx%d 격자의 고유 경로 수: %d\n", m, n2, UniquePaths(m, n2))

	// 최소 경로 합
	fmt.Println("\n[최소 경로 합]")
	grid := [][]int{
		{1, 3, 1},
		{1, 5, 1},
		{4, 2, 1},
	}
	fmt.Println("격자:")
	for _, row := range grid {
		fmt.Printf("  %v\n", row)
	}
	fmt.Printf("최소 경로 합: %d\n", MinPathSum(grid))

	// 배낭 문제
	fmt.Println("\n[0/1 배낭 문제]")
	weights := []int{1, 2, 3}
	values := []int{6, 10, 12}
	capacity := 5
	fmt.Printf("무게: %v, 가치: %v, 용량: %d\n", weights, values, capacity)
	fmt.Printf("최대 가치 (2D DP): %d\n", Knapsack01(weights, values, capacity))
	fmt.Printf("최대 가치 (1D DP): %d\n", Knapsack01Optimized(weights, values, capacity))

	// LCS
	fmt.Println("\n[최장 공통 부분 수열 (LCS)]")
	s1, s2 := "abcde", "ace"
	fmt.Printf("문자열1: \"%s\", 문자열2: \"%s\"\n", s1, s2)
	fmt.Printf("LCS 길이: %d\n", LongestCommonSubsequence(s1, s2))

	// 편집 거리
	fmt.Println("\n[편집 거리]")
	s3, s4 := "horse", "ros"
	fmt.Printf("문자열1: \"%s\", 문자열2: \"%s\"\n", s3, s4)
	fmt.Printf("편집 거리: %d\n", EditDistance(s3, s4))

	// 응용 문제
	fmt.Println("\n============================================")
	fmt.Println("[응용 문제]")
	fmt.Println("============================================")

	// Decode Ways
	fmt.Println("\n[Decode Ways]")
	encoded := "226"
	fmt.Printf("숫자열: \"%s\"\n", encoded)
	fmt.Printf("해석 방법 수: %d\n", DecodeWays(encoded))

	// Word Break
	fmt.Println("\n[Word Break]")
	str := "leetcode"
	dict := []string{"leet", "code"}
	fmt.Printf("문자열: \"%s\", 사전: %v\n", str, dict)
	fmt.Printf("분리 가능: %v\n", WordBreak(str, dict))

	// DP 접근법 비교
	fmt.Println("\n============================================")
	fmt.Println("DP 패턴 요약")
	fmt.Println("============================================")
	fmt.Println(`
┌─────────────────────────────────────────────────────────────┐
│ Top-Down (메모이제이션)                                      │
├─────────────────────────────────────────────────────────────┤
│ - 재귀 + 캐싱                                                │
│ - 직관적, 필요한 부분만 계산                                  │
│ - 함수 호출 오버헤드, 스택 제한                              │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ Bottom-Up (타뷸레이션)                                       │
├─────────────────────────────────────────────────────────────┤
│ - 반복문 + 테이블                                            │
│ - 효율적, 공간 최적화 가능                                   │
│ - 모든 부분 문제 계산                                        │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ DP 문제 유형                                                 │
├─────────────────────────────────────────────────────────────┤
│ 1D: 피보나치, 계단, House Robber, 동전, LIS                 │
│ 2D: 경로, 배낭, LCS, 편집거리                               │
│ 구간 DP: 행렬 곱셈, 팰린드롬                                 │
│ 비트마스크 DP: TSP, 상태 압축                               │
└─────────────────────────────────────────────────────────────┘

[DP 문제 인식 신호]
- "최대/최소/방법의 수"를 구하라
- "이전 선택이 현재에 영향"
- "중복되는 하위 문제"가 보임
`)
}
