package main

import "fmt"

/*
===========================================
6-9. 백트래킹 (Backtracking)
===========================================

백트래킹이란?
- 모든 가능한 경우를 탐색하되, 유망하지 않은 경로는 빨리 포기
- DFS + 가지치기(Pruning)
- "되돌아가기" - 선택을 취소하고 다른 선택 시도

핵심 개념:
1. 선택(Choose): 가능한 선택지 중 하나 선택
2. 탐색(Explore): 선택을 기반으로 다음 단계 탐색
3. 철회(Unchoose): 탐색 후 선택 취소 (원상복구)

사용 사례:
- 순열(Permutation) / 조합(Combination) 생성
- N-Queens 문제
- 스도쿠 해결
- 미로 찾기
- 부분집합 생성
*/

// ============================================
// 순열 (Permutation)
// ============================================

// Permute 배열의 모든 순열 생성
// 시간복잡도: O(n! * n)
func Permute(nums []int) [][]int {
	result := [][]int{}

	var backtrack func(path []int, used []bool)
	backtrack = func(path []int, used []bool) {
		// 기저 조건: 모든 원소를 사용했으면 완성
		if len(path) == len(nums) {
			// path 복사본을 결과에 추가
			temp := make([]int, len(path))
			copy(temp, path)
			result = append(result, temp)
			return
		}

		for i := 0; i < len(nums); i++ {
			// 이미 사용한 원소는 건너뜀
			if used[i] {
				continue
			}

			// 선택 (Choose)
			path = append(path, nums[i])
			used[i] = true

			// 탐색 (Explore)
			backtrack(path, used)

			// 철회 (Unchoose)
			path = path[:len(path)-1]
			used[i] = false
		}
	}

	backtrack([]int{}, make([]bool, len(nums)))
	return result
}

// PermuteUnique 중복이 있는 배열의 고유 순열
func PermuteUnique(nums []int) [][]int {
	result := [][]int{}

	// 먼저 정렬 (중복 제거를 위해)
	sortInts(nums)

	var backtrack func(path []int, used []bool)
	backtrack = func(path []int, used []bool) {
		if len(path) == len(nums) {
			temp := make([]int, len(path))
			copy(temp, path)
			result = append(result, temp)
			return
		}

		for i := 0; i < len(nums); i++ {
			if used[i] {
				continue
			}

			// 중복 건너뛰기: 이전과 같은 값이고, 이전이 사용되지 않았다면
			if i > 0 && nums[i] == nums[i-1] && !used[i-1] {
				continue
			}

			path = append(path, nums[i])
			used[i] = true

			backtrack(path, used)

			path = path[:len(path)-1]
			used[i] = false
		}
	}

	backtrack([]int{}, make([]bool, len(nums)))
	return result
}

// 간단한 정렬 헬퍼
func sortInts(arr []int) {
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[j] < arr[i] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}

// ============================================
// 조합 (Combination)
// ============================================

// Combine n개 중 k개를 선택하는 모든 조합
// C(n, k) 개의 조합 생성
func Combine(n, k int) [][]int {
	result := [][]int{}

	var backtrack func(start int, path []int)
	backtrack = func(start int, path []int) {
		// k개 선택 완료
		if len(path) == k {
			temp := make([]int, len(path))
			copy(temp, path)
			result = append(result, temp)
			return
		}

		// 가지치기: 남은 요소가 부족하면 탐색 중단
		remaining := n - start + 1
		needed := k - len(path)
		if remaining < needed {
			return
		}

		// start부터 n까지 순회 (순서 고려 X)
		for i := start; i <= n; i++ {
			path = append(path, i)
			backtrack(i+1, path) // i+1부터 시작하여 중복 방지
			path = path[:len(path)-1]
		}
	}

	backtrack(1, []int{})
	return result
}

// CombinationSum 합이 target이 되는 조합 (중복 사용 가능)
func CombinationSum(candidates []int, target int) [][]int {
	result := [][]int{}

	var backtrack func(start int, path []int, remaining int)
	backtrack = func(start int, path []int, remaining int) {
		// 목표 달성
		if remaining == 0 {
			temp := make([]int, len(path))
			copy(temp, path)
			result = append(result, temp)
			return
		}

		// 초과하면 중단
		if remaining < 0 {
			return
		}

		for i := start; i < len(candidates); i++ {
			path = append(path, candidates[i])
			// 같은 수를 다시 사용할 수 있으므로 i 유지
			backtrack(i, path, remaining-candidates[i])
			path = path[:len(path)-1]
		}
	}

	backtrack(0, []int{}, target)
	return result
}

// CombinationSum2 합이 target이 되는 조합 (각 요소 한 번만 사용)
func CombinationSum2(candidates []int, target int) [][]int {
	result := [][]int{}
	sortInts(candidates) // 중복 처리를 위해 정렬

	var backtrack func(start int, path []int, remaining int)
	backtrack = func(start int, path []int, remaining int) {
		if remaining == 0 {
			temp := make([]int, len(path))
			copy(temp, path)
			result = append(result, temp)
			return
		}

		for i := start; i < len(candidates); i++ {
			// 가지치기: 현재 후보가 남은 값보다 크면 중단
			if candidates[i] > remaining {
				break
			}

			// 같은 레벨에서 중복 건너뛰기
			if i > start && candidates[i] == candidates[i-1] {
				continue
			}

			path = append(path, candidates[i])
			backtrack(i+1, path, remaining-candidates[i])
			path = path[:len(path)-1]
		}
	}

	backtrack(0, []int{}, target)
	return result
}

// ============================================
// 부분집합 (Subsets)
// ============================================

// Subsets 모든 부분집합 생성
// 2^n개의 부분집합
func Subsets(nums []int) [][]int {
	result := [][]int{}

	var backtrack func(start int, path []int)
	backtrack = func(start int, path []int) {
		// 현재 경로를 결과에 추가 (모든 경로가 유효한 부분집합)
		temp := make([]int, len(path))
		copy(temp, path)
		result = append(result, temp)

		for i := start; i < len(nums); i++ {
			path = append(path, nums[i])
			backtrack(i+1, path)
			path = path[:len(path)-1]
		}
	}

	backtrack(0, []int{})
	return result
}

// SubsetsWithDup 중복이 있는 배열의 고유 부분집합
func SubsetsWithDup(nums []int) [][]int {
	result := [][]int{}
	sortInts(nums)

	var backtrack func(start int, path []int)
	backtrack = func(start int, path []int) {
		temp := make([]int, len(path))
		copy(temp, path)
		result = append(result, temp)

		for i := start; i < len(nums); i++ {
			// 같은 레벨에서 중복 건너뛰기
			if i > start && nums[i] == nums[i-1] {
				continue
			}

			path = append(path, nums[i])
			backtrack(i+1, path)
			path = path[:len(path)-1]
		}
	}

	backtrack(0, []int{})
	return result
}

// ============================================
// N-Queens 문제
// ============================================

// SolveNQueens N-Queens 문제의 모든 해답
// N×N 체스판에 N개의 퀸을 서로 공격하지 않게 배치
func SolveNQueens(n int) [][]string {
	result := [][]string{}

	// 체스판 초기화
	board := make([][]byte, n)
	for i := range board {
		board[i] = make([]byte, n)
		for j := range board[i] {
			board[i][j] = '.'
		}
	}

	// 유효성 검사: (row, col)에 퀸을 놓을 수 있는지
	isValid := func(row, col int) bool {
		// 같은 열 검사
		for i := 0; i < row; i++ {
			if board[i][col] == 'Q' {
				return false
			}
		}

		// 왼쪽 위 대각선 검사
		for i, j := row-1, col-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
			if board[i][j] == 'Q' {
				return false
			}
		}

		// 오른쪽 위 대각선 검사
		for i, j := row-1, col+1; i >= 0 && j < n; i, j = i-1, j+1 {
			if board[i][j] == 'Q' {
				return false
			}
		}

		return true
	}

	// 보드를 문자열 슬라이스로 변환
	boardToStrings := func() []string {
		res := make([]string, n)
		for i := range board {
			res[i] = string(board[i])
		}
		return res
	}

	var backtrack func(row int)
	backtrack = func(row int) {
		// 모든 행에 퀸 배치 완료
		if row == n {
			result = append(result, boardToStrings())
			return
		}

		// 현재 행의 각 열에 퀸 배치 시도
		for col := 0; col < n; col++ {
			if !isValid(row, col) {
				continue
			}

			board[row][col] = 'Q' // 선택
			backtrack(row + 1)    // 다음 행
			board[row][col] = '.' // 철회
		}
	}

	backtrack(0)
	return result
}

// TotalNQueens N-Queens 해답의 개수만 반환
func TotalNQueens(n int) int {
	count := 0

	// 비트마스크로 열, 대각선 상태 관리 (더 빠름)
	var backtrack func(row, cols, diag1, diag2 int)
	backtrack = func(row, cols, diag1, diag2 int) {
		if row == n {
			count++
			return
		}

		// 사용 가능한 열 찾기
		availablePositions := ((1 << n) - 1) & ^(cols | diag1 | diag2)

		for availablePositions != 0 {
			// 가장 낮은 1비트 선택
			position := availablePositions & (-availablePositions)
			availablePositions ^= position

			backtrack(row+1,
				cols|position,
				(diag1|position)<<1,
				(diag2|position)>>1)
		}
	}

	backtrack(0, 0, 0, 0)
	return count
}

// ============================================
// 스도쿠 (Sudoku)
// ============================================

// SolveSudoku 스도쿠 풀기
// board[i][j]가 '.'이면 빈 칸
func SolveSudoku(board [][]byte) bool {
	// 다음 빈 칸 찾기
	row, col := -1, -1
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j] == '.' {
				row, col = i, j
				break
			}
		}
		if row != -1 {
			break
		}
	}

	// 빈 칸이 없으면 완성
	if row == -1 {
		return true
	}

	// 유효성 검사
	isValid := func(num byte) bool {
		// 같은 행 검사
		for j := 0; j < 9; j++ {
			if board[row][j] == num {
				return false
			}
		}

		// 같은 열 검사
		for i := 0; i < 9; i++ {
			if board[i][col] == num {
				return false
			}
		}

		// 3x3 박스 검사
		boxRow, boxCol := (row/3)*3, (col/3)*3
		for i := boxRow; i < boxRow+3; i++ {
			for j := boxCol; j < boxCol+3; j++ {
				if board[i][j] == num {
					return false
				}
			}
		}

		return true
	}

	// 1-9 시도
	for num := byte('1'); num <= '9'; num++ {
		if !isValid(num) {
			continue
		}

		board[row][col] = num // 선택

		if SolveSudoku(board) { // 탐색
			return true
		}

		board[row][col] = '.' // 철회
	}

	return false
}

// ============================================
// 문자열 관련
// ============================================

// GenerateParenthesis 올바른 괄호 조합 생성
func GenerateParenthesis(n int) []string {
	result := []string{}

	var backtrack func(path string, open, close int)
	backtrack = func(path string, open, close int) {
		// n개씩 사용 완료
		if len(path) == 2*n {
			result = append(result, path)
			return
		}

		// 여는 괄호 추가 가능
		if open < n {
			backtrack(path+"(", open+1, close)
		}

		// 닫는 괄호 추가 가능 (여는 괄호보다 적을 때만)
		if close < open {
			backtrack(path+")", open, close+1)
		}
	}

	backtrack("", 0, 0)
	return result
}

// LetterCombinations 전화 키패드 문자 조합
func LetterCombinations(digits string) []string {
	if len(digits) == 0 {
		return []string{}
	}

	// 숫자 -> 문자 매핑
	mapping := map[byte]string{
		'2': "abc",
		'3': "def",
		'4': "ghi",
		'5': "jkl",
		'6': "mno",
		'7': "pqrs",
		'8': "tuv",
		'9': "wxyz",
	}

	result := []string{}

	var backtrack func(index int, path string)
	backtrack = func(index int, path string) {
		if index == len(digits) {
			result = append(result, path)
			return
		}

		letters := mapping[digits[index]]
		for i := 0; i < len(letters); i++ {
			backtrack(index+1, path+string(letters[i]))
		}
	}

	backtrack(0, "")
	return result
}

// Partition 팰린드롬 파티셔닝
func Partition(s string) [][]string {
	result := [][]string{}

	isPalindrome := func(str string) bool {
		left, right := 0, len(str)-1
		for left < right {
			if str[left] != str[right] {
				return false
			}
			left++
			right--
		}
		return true
	}

	var backtrack func(start int, path []string)
	backtrack = func(start int, path []string) {
		// 문자열 끝에 도달
		if start == len(s) {
			temp := make([]string, len(path))
			copy(temp, path)
			result = append(result, temp)
			return
		}

		// 모든 부분 문자열 시도
		for end := start + 1; end <= len(s); end++ {
			substr := s[start:end]

			// 팰린드롬인 경우만 진행
			if isPalindrome(substr) {
				path = append(path, substr)
				backtrack(end, path)
				path = path[:len(path)-1]
			}
		}
	}

	backtrack(0, []string{})
	return result
}

func main() {
	fmt.Println("============================================")
	fmt.Println("6-9. 백트래킹 (Backtracking)")
	fmt.Println("============================================")

	// 순열
	fmt.Println("\n[순열 (Permutation)]")
	fmt.Println("--------------------------------------------")
	nums := []int{1, 2, 3}
	fmt.Printf("입력: %v\n", nums)
	fmt.Println("순열:")
	for _, perm := range Permute(nums) {
		fmt.Printf("  %v\n", perm)
	}

	// 중복 있는 순열
	fmt.Println("\n[중복 있는 순열]")
	nums2 := []int{1, 1, 2}
	fmt.Printf("입력: %v\n", nums2)
	fmt.Println("고유 순열:")
	for _, perm := range PermuteUnique(nums2) {
		fmt.Printf("  %v\n", perm)
	}

	// 조합
	fmt.Println("\n[조합 (Combination)]")
	fmt.Println("--------------------------------------------")
	n, k := 4, 2
	fmt.Printf("C(%d, %d):\n", n, k)
	for _, comb := range Combine(n, k) {
		fmt.Printf("  %v\n", comb)
	}

	// 조합 합
	fmt.Println("\n[조합 합 (중복 사용 가능)]")
	candidates := []int{2, 3, 6, 7}
	target := 7
	fmt.Printf("후보: %v, 목표: %d\n", candidates, target)
	for _, comb := range CombinationSum(candidates, target) {
		fmt.Printf("  %v\n", comb)
	}

	// 부분집합
	fmt.Println("\n[부분집합 (Subsets)]")
	fmt.Println("--------------------------------------------")
	nums3 := []int{1, 2, 3}
	fmt.Printf("입력: %v\n", nums3)
	fmt.Printf("부분집합 개수: %d\n", len(Subsets(nums3)))

	// N-Queens
	fmt.Println("\n============================================")
	fmt.Println("[N-Queens 문제]")
	fmt.Println("============================================")
	nq := 4
	solutions := SolveNQueens(nq)
	fmt.Printf("%d-Queens 해답 (%d개):\n", nq, len(solutions))
	for i, sol := range solutions {
		fmt.Printf("\n해답 %d:\n", i+1)
		for _, row := range sol {
			fmt.Printf("  %s\n", row)
		}
	}

	fmt.Printf("\n8-Queens 해답 개수: %d\n", TotalNQueens(8))

	// 스도쿠
	fmt.Println("\n============================================")
	fmt.Println("[스도쿠]")
	fmt.Println("============================================")
	board := [][]byte{
		{'5', '3', '.', '.', '7', '.', '.', '.', '.'},
		{'6', '.', '.', '1', '9', '5', '.', '.', '.'},
		{'.', '9', '8', '.', '.', '.', '.', '6', '.'},
		{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
		{'4', '.', '.', '8', '.', '3', '.', '.', '1'},
		{'7', '.', '.', '.', '2', '.', '.', '.', '6'},
		{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
		{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
		{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
	}

	fmt.Println("문제:")
	for _, row := range board {
		fmt.Printf("  %s\n", string(row))
	}

	SolveSudoku(board)
	fmt.Println("\n해답:")
	for _, row := range board {
		fmt.Printf("  %s\n", string(row))
	}

	// 문자열 문제
	fmt.Println("\n============================================")
	fmt.Println("[문자열 백트래킹]")
	fmt.Println("============================================")

	// 괄호 생성
	fmt.Println("\n[올바른 괄호 조합]")
	fmt.Printf("n=3: %v\n", GenerateParenthesis(3))

	// 전화 키패드
	fmt.Println("\n[전화 키패드 문자 조합]")
	digits := "23"
	fmt.Printf("digits=\"%s\": %v\n", digits, LetterCombinations(digits))

	// 팰린드롬 파티셔닝
	fmt.Println("\n[팰린드롬 파티셔닝]")
	str := "aab"
	fmt.Printf("s=\"%s\":\n", str)
	for _, p := range Partition(str) {
		fmt.Printf("  %v\n", p)
	}

	// 패턴 요약
	fmt.Println("\n============================================")
	fmt.Println("백트래킹 패턴 요약")
	fmt.Println("============================================")
	fmt.Println(`
┌─────────────────────────────────────────────────────────────┐
│ 백트래킹 템플릿                                              │
├─────────────────────────────────────────────────────────────┤
│ func backtrack(path []int, ...) {                           │
│     // 1. 종료 조건                                         │
│     if 완성 조건 {                                          │
│         결과에 path 복사본 추가                              │
│         return                                              │
│     }                                                       │
│                                                             │
│     // 2. 선택지 순회                                       │
│     for 선택 in 선택지들 {                                  │
│         if 유효하지 않은 선택 { continue }  // 가지치기     │
│                                                             │
│         path에 선택 추가    // Choose                        │
│         backtrack(path)      // Explore                      │
│         path에서 선택 제거  // Unchoose                      │
│     }                                                       │
│ }                                                           │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ 중복 처리                                                    │
├─────────────────────────────────────────────────────────────┤
│ 1. 정렬 후 인접 중복 건너뛰기                               │
│    if i > start && nums[i] == nums[i-1] { continue }        │
│                                                             │
│ 2. visited/used 배열 사용                                   │
│    if used[i] { continue }                                  │
└─────────────────────────────────────────────────────────────┘

[문제 유형별 접근]
- 순열: 모든 원소 사용, used 배열 사용
- 조합: start 인덱스로 순서 제어
- 부분집합: 각 원소 선택/비선택
- N-Queens: 행별로 열 선택, 대각선 검사
`)
}
