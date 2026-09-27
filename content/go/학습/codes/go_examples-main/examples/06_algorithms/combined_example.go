package main

import (
	"container/heap"
	"fmt"
)

/*
===========================================
알고리즘 통합 예제
===========================================

이 파일은 여러 알고리즘을 조합하여 실제 문제를 해결하는 예시입니다.
알고리즘 선택 기준과 조합 방법을 배울 수 있습니다.
*/

// ============================================
// 예제 1: 미로 최단 경로 (BFS + 2D 탐색)
// ============================================

type Point struct {
	row, col, dist int
}

// ShortestPathMaze BFS로 미로의 최단 경로 찾기
func ShortestPathMaze(maze [][]int, start, end [2]int) (int, [][]int) {
	rows, cols := len(maze), len(maze[0])
	directions := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	// 방문 여부 및 이전 위치 기록
	visited := make([][]bool, rows)
	prev := make([][][2]int, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
		prev[i] = make([][2]int, cols)
		for j := range prev[i] {
			prev[i][j] = [2]int{-1, -1}
		}
	}

	// BFS
	queue := []Point{{start[0], start[1], 0}}
	visited[start[0]][start[1]] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		// 도착점 도달
		if curr.row == end[0] && curr.col == end[1] {
			// 경로 재구성
			path := [][]int{}
			r, c := end[0], end[1]
			for r != -1 && c != -1 {
				path = append([][]int{{r, c}}, path...)
				pr, pc := prev[r][c][0], prev[r][c][1]
				r, c = pr, pc
			}
			return curr.dist, path
		}

		for _, dir := range directions {
			nr, nc := curr.row+dir[0], curr.col+dir[1]
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols &&
				!visited[nr][nc] && maze[nr][nc] == 0 {
				visited[nr][nc] = true
				prev[nr][nc] = [2]int{curr.row, curr.col}
				queue = append(queue, Point{nr, nc, curr.dist + 1})
			}
		}
	}

	return -1, nil // 도달 불가
}

// ============================================
// 예제 2: 단어 사다리 (BFS + HashMap)
// ============================================

// WordLadder 시작 단어에서 끝 단어까지의 최단 변환 횟수
func WordLadder(beginWord, endWord string, wordList []string) int {
	// 단어 목록을 Set으로 변환
	wordSet := make(map[string]bool)
	for _, word := range wordList {
		wordSet[word] = true
	}

	if !wordSet[endWord] {
		return 0
	}

	// BFS
	queue := []struct {
		word  string
		steps int
	}{{beginWord, 1}}

	visited := make(map[string]bool)
	visited[beginWord] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.word == endWord {
			return curr.steps
		}

		// 한 글자씩 변경
		wordBytes := []byte(curr.word)
		for i := 0; i < len(wordBytes); i++ {
			original := wordBytes[i]
			for c := byte('a'); c <= 'z'; c++ {
				if c == original {
					continue
				}
				wordBytes[i] = c
				newWord := string(wordBytes)

				if wordSet[newWord] && !visited[newWord] {
					visited[newWord] = true
					queue = append(queue, struct {
						word  string
						steps int
					}{newWord, curr.steps + 1})
				}
			}
			wordBytes[i] = original
		}
	}

	return 0
}

// ============================================
// 예제 3: 주식 최대 이익 (DP + 슬라이딩 윈도우 개념)
// ============================================

// MaxProfitWithCooldown 쿨다운이 있는 주식 매매 최대 이익
// 매도 후 하루 쉬어야 함
func MaxProfitWithCooldown(prices []int) int {
	n := len(prices)
	if n < 2 {
		return 0
	}

	// 상태:
	// hold[i]: i일에 주식을 보유한 상태의 최대 이익
	// sold[i]: i일에 주식을 판 상태의 최대 이익
	// rest[i]: i일에 쉬는 상태의 최대 이익

	hold := make([]int, n)
	sold := make([]int, n)
	rest := make([]int, n)

	hold[0] = -prices[0]
	sold[0] = 0
	rest[0] = 0

	for i := 1; i < n; i++ {
		hold[i] = max(hold[i-1], rest[i-1]-prices[i])
		sold[i] = hold[i-1] + prices[i]
		rest[i] = max(rest[i-1], sold[i-1])
	}

	return max(sold[n-1], rest[n-1])
}

// ============================================
// 예제 4: 최소 회의실 (그리디 + 힙)
// ============================================

// MeetingRoomsII 필요한 최소 회의실 수 (힙 사용)
type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func MeetingRoomsII(intervals [][]int) int {
	if len(intervals) == 0 {
		return 0
	}

	// 시작 시간 기준 정렬
	sortIntervals(intervals)

	// 최소 힙: 현재 사용 중인 회의실들의 종료 시간
	endTimes := &MinHeap{}
	heap.Init(endTimes)
	heap.Push(endTimes, intervals[0][1])

	for i := 1; i < len(intervals); i++ {
		// 가장 빨리 끝나는 회의실 확인
		if intervals[i][0] >= (*endTimes)[0] {
			heap.Pop(endTimes) // 회의실 재사용
		}
		heap.Push(endTimes, intervals[i][1])
	}

	return endTimes.Len()
}

func sortIntervals(intervals [][]int) {
	for i := 0; i < len(intervals); i++ {
		for j := i + 1; j < len(intervals); j++ {
			if intervals[j][0] < intervals[i][0] {
				intervals[i], intervals[j] = intervals[j], intervals[i]
			}
		}
	}
}

// ============================================
// 예제 5: 과제 스케줄링 (그리디 + DP)
// ============================================

// Assignment 과제 구조체
type Assignment struct {
	deadline int
	points   int
}

// MaxPoints 마감 기한 내 최대 점수 획득
func MaxPoints(assignments []Assignment) int {
	// 점수 기준 내림차순 정렬
	sortAssignments(assignments)

	// 최대 마감일 찾기
	maxDeadline := 0
	for _, a := range assignments {
		if a.deadline > maxDeadline {
			maxDeadline = a.deadline
		}
	}

	// 슬롯 사용 여부
	used := make([]bool, maxDeadline+1)
	total := 0

	for _, a := range assignments {
		// 마감일 직전부터 빈 슬롯 찾기
		for day := a.deadline; day > 0; day-- {
			if !used[day] {
				used[day] = true
				total += a.points
				break
			}
		}
	}

	return total
}

func sortAssignments(assignments []Assignment) {
	for i := 0; i < len(assignments); i++ {
		for j := i + 1; j < len(assignments); j++ {
			if assignments[j].points > assignments[i].points {
				assignments[i], assignments[j] = assignments[j], assignments[i]
			}
		}
	}
}

// ============================================
// 예제 6: 섬 최대 영역 (DFS + 카운팅)
// ============================================

// MaxAreaOfIsland DFS로 가장 큰 섬의 넓이
func MaxAreaOfIsland(grid [][]int) int {
	if len(grid) == 0 {
		return 0
	}

	rows, cols := len(grid), len(grid[0])
	maxArea := 0

	var dfs func(r, c int) int
	dfs = func(r, c int) int {
		if r < 0 || r >= rows || c < 0 || c >= cols || grid[r][c] == 0 {
			return 0
		}

		grid[r][c] = 0 // 방문 처리
		area := 1
		area += dfs(r-1, c)
		area += dfs(r+1, c)
		area += dfs(r, c-1)
		area += dfs(r, c+1)
		return area
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if grid[i][j] == 1 {
				area := dfs(i, j)
				if area > maxArea {
					maxArea = area
				}
			}
		}
	}

	return maxArea
}

// ============================================
// 예제 7: 연속 부분 배열 (슬라이딩 윈도우 + HashMap)
// ============================================

// LongestSubarrayWithKDistinct K개 이하의 고유 값을 가진 최장 부분 배열
func LongestSubarrayWithKDistinct(nums []int, k int) int {
	if k == 0 {
		return 0
	}

	freq := make(map[int]int)
	maxLen := 0
	left := 0

	for right := 0; right < len(nums); right++ {
		freq[nums[right]]++

		for len(freq) > k {
			freq[nums[left]]--
			if freq[nums[left]] == 0 {
				delete(freq, nums[left])
			}
			left++
		}

		if right-left+1 > maxLen {
			maxLen = right - left + 1
		}
	}

	return maxLen
}

// ============================================
// 예제 8: 최단 경로 + 조건 (다익스트라 + DP)
// ============================================

// CheapestFlightsWithKStops K번 이하 경유로 가장 저렴한 항공편
func CheapestFlightsWithKStops(n int, flights [][]int, src, dst, k int) int {
	// 벨만-포드 변형 (k+1번 반복)
	prices := make([]int, n)
	for i := range prices {
		prices[i] = 1<<31 - 1
	}
	prices[src] = 0

	for i := 0; i <= k; i++ {
		// 현재 라운드의 결과를 저장할 임시 배열
		temp := make([]int, n)
		copy(temp, prices)

		for _, flight := range flights {
			from, to, price := flight[0], flight[1], flight[2]
			if prices[from] != 1<<31-1 {
				if prices[from]+price < temp[to] {
					temp[to] = prices[from] + price
				}
			}
		}
		prices = temp
	}

	if prices[dst] == 1<<31-1 {
		return -1
	}
	return prices[dst]
}

func main() {
	fmt.Println("============================================")
	fmt.Println("알고리즘 통합 예제")
	fmt.Println("============================================")

	// 예제 1: 미로 최단 경로
	fmt.Println("\n[예제 1: 미로 최단 경로 (BFS)]")
	fmt.Println("--------------------------------------------")
	maze := [][]int{
		{0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 1, 1, 0},
		{1, 0, 1, 0, 0},
		{0, 0, 0, 0, 0},
	}
	fmt.Println("미로 (0: 통과 가능, 1: 벽):")
	for _, row := range maze {
		fmt.Printf("  %v\n", row)
	}
	dist, path := ShortestPathMaze(maze, [2]int{0, 0}, [2]int{4, 4})
	fmt.Printf("최단 거리: %d\n", dist)
	fmt.Printf("경로: %v\n", path)

	// 예제 2: 단어 사다리
	fmt.Println("\n[예제 2: 단어 사다리 (BFS + HashMap)]")
	fmt.Println("--------------------------------------------")
	beginWord := "hit"
	endWord := "cog"
	wordList := []string{"hot", "dot", "dog", "lot", "log", "cog"}
	fmt.Printf("시작: %s, 도착: %s\n", beginWord, endWord)
	fmt.Printf("단어 목록: %v\n", wordList)
	fmt.Printf("최단 변환 횟수: %d\n", WordLadder(beginWord, endWord, wordList))

	// 예제 3: 주식 매매 (쿨다운)
	fmt.Println("\n[예제 3: 주식 매매 with 쿨다운 (DP)]")
	fmt.Println("--------------------------------------------")
	prices := []int{1, 2, 3, 0, 2}
	fmt.Printf("주가: %v\n", prices)
	fmt.Printf("최대 이익: %d\n", MaxProfitWithCooldown(prices))

	// 예제 4: 최소 회의실
	fmt.Println("\n[예제 4: 최소 회의실 (그리디 + 힙)]")
	fmt.Println("--------------------------------------------")
	meetings := [][]int{{0, 30}, {5, 10}, {15, 20}}
	fmt.Printf("회의 시간: %v\n", meetings)
	fmt.Printf("필요한 최소 회의실: %d\n", MeetingRoomsII(meetings))

	// 예제 5: 과제 스케줄링
	fmt.Println("\n[예제 5: 과제 스케줄링 (그리디)]")
	fmt.Println("--------------------------------------------")
	assignments := []Assignment{
		{4, 70},
		{2, 60},
		{4, 50},
		{3, 40},
		{1, 30},
		{4, 20},
	}
	fmt.Println("과제 (마감, 점수):")
	for _, a := range assignments {
		fmt.Printf("  (%d, %d)\n", a.deadline, a.points)
	}
	fmt.Printf("최대 획득 점수: %d\n", MaxPoints(assignments))

	// 예제 6: 섬 최대 영역
	fmt.Println("\n[예제 6: 섬 최대 영역 (DFS)]")
	fmt.Println("--------------------------------------------")
	islandGrid := [][]int{
		{0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0},
		{0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
	}
	fmt.Printf("섬 최대 넓이: %d\n", MaxAreaOfIsland(islandGrid))

	// 예제 7: 연속 부분 배열
	fmt.Println("\n[예제 7: K개 이하 고유 값 최장 부분 배열 (슬라이딩 윈도우)]")
	fmt.Println("--------------------------------------------")
	nums := []int{1, 2, 1, 2, 3}
	k := 2
	fmt.Printf("배열: %v, k=%d\n", nums, k)
	fmt.Printf("최장 길이: %d\n", LongestSubarrayWithKDistinct(nums, k))

	// 예제 8: 항공편 최저가
	fmt.Println("\n[예제 8: K번 경유 최저가 항공편 (벨만-포드 변형)]")
	fmt.Println("--------------------------------------------")
	flights := [][]int{{0, 1, 100}, {1, 2, 100}, {0, 2, 500}}
	fmt.Printf("항공편 (출발, 도착, 가격): %v\n", flights)
	fmt.Printf("도시 수: 3, 출발: 0, 도착: 2, 최대 경유: 1\n")
	fmt.Printf("최저 가격: %d\n", CheapestFlightsWithKStops(3, flights, 0, 2, 1))

	// 알고리즘 선택 가이드
	fmt.Println("\n============================================")
	fmt.Println("알고리즘 선택 가이드")
	fmt.Println("============================================")
	fmt.Println(`
┌─────────────────────────────────────────────────────────────┐
│ 문제 유형별 알고리즘 선택                                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ "최단 경로/거리" →                                          │
│   ├─ 가중치 없음 → BFS                                      │
│   ├─ 양수 가중치 → 다익스트라                               │
│   └─ 음수 가중치 → 벨만-포드                                │
│                                                             │
│ "최대/최소/방법의 수" →                                     │
│   ├─ 중복 부분 문제 → DP                                    │
│   └─ 그리디 가능 → 그리디                                   │
│                                                             │
│ "모든 경우 탐색" →                                          │
│   ├─ 순열/조합 → 백트래킹                                   │
│   └─ 그래프 탐색 → DFS/BFS                                  │
│                                                             │
│ "연속 구간 처리" →                                          │
│   ├─ 고정 크기 → 슬라이딩 윈도우                            │
│   └─ 가변 크기 → 투 포인터/슬라이딩 윈도우                  │
│                                                             │
│ "정렬된 배열 탐색" → 이진 탐색                              │
│                                                             │
│ "집합 연결/분리" → 유니온 파인드                            │
│                                                             │
│ "빈도/중복 처리" → HashMap                                  │
│                                                             │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ 알고리즘 조합 패턴                                           │
├─────────────────────────────────────────────────────────────┤
│ • BFS + HashMap: 단어 변환, 최단 경로 + 상태 관리           │
│ • DP + 그리디: 최적화 문제에서 부분 최적 활용               │
│ • DFS + 메모이제이션: Top-Down DP                           │
│ • 이진 탐색 + 그리디: 매개변수 탐색                         │
│ • 슬라이딩 윈도우 + HashMap: 조건 있는 부분 배열            │
│ • 정렬 + 투 포인터: 합/차 찾기 문제                         │
│ • 그래프 + 유니온 파인드: 연결성 + 최적화                   │
└─────────────────────────────────────────────────────────────┘
`)
}
