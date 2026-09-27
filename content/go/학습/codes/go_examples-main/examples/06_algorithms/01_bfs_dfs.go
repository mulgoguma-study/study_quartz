package main

import "fmt"

/*
===========================================
6-1. BFS (너비 우선 탐색) / DFS (깊이 우선 탐색)
===========================================

BFS (Breadth-First Search):
- 가까운 노드부터 탐색하는 알고리즘
- 큐(Queue)를 사용하여 구현
- 최단 경로 문제에 적합
- 시간복잡도: O(V + E) (V: 정점 수, E: 간선 수)

DFS (Depth-First Search):
- 깊은 노드부터 탐색하는 알고리즘
- 스택(Stack) 또는 재귀를 사용하여 구현
- 경로 탐색, 사이클 감지에 적합
- 시간복잡도: O(V + E)

사용 사례:
- BFS: 최단 거리, 레벨별 탐색, 네트워크 브로드캐스트
- DFS: 미로 찾기, 위상 정렬, 연결 요소 찾기
*/

// Graph 인접 리스트로 표현한 그래프
type Graph struct {
	vertices int           // 정점의 개수
	adjList  map[int][]int // 인접 리스트
}

// NewGraph 새 그래프 생성
func NewGraph(vertices int) *Graph {
	return &Graph{
		vertices: vertices,
		adjList:  make(map[int][]int),
	}
}

// AddEdge 간선 추가 (무방향 그래프)
func (g *Graph) AddEdge(v, w int) {
	g.adjList[v] = append(g.adjList[v], w)
	g.adjList[w] = append(g.adjList[w], v) // 무방향이므로 양쪽에 추가
}

// AddDirectedEdge 간선 추가 (방향 그래프)
func (g *Graph) AddDirectedEdge(from, to int) {
	g.adjList[from] = append(g.adjList[from], to)
}

// ============================================
// BFS 구현
// ============================================

// BFS 너비 우선 탐색
func (g *Graph) BFS(start int) []int {
	// 방문 여부를 추적하는 맵
	visited := make(map[int]bool)

	// 탐색 결과를 저장할 슬라이스
	result := []int{}

	// BFS는 큐를 사용
	// Go에서는 슬라이스로 큐를 구현
	queue := []int{start}
	visited[start] = true

	for len(queue) > 0 {
		// 큐에서 맨 앞 요소를 꺼냄 (Dequeue)
		current := queue[0]
		queue = queue[1:]

		// 현재 노드를 결과에 추가
		result = append(result, current)

		// 인접한 모든 노드를 확인
		for _, neighbor := range g.adjList[current] {
			// 아직 방문하지 않은 노드만 큐에 추가
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}

	return result
}

// BFSWithLevel 레벨별 BFS (각 노드의 깊이도 함께 반환)
func (g *Graph) BFSWithLevel(start int) map[int]int {
	visited := make(map[int]bool)
	level := make(map[int]int) // 각 노드의 깊이(레벨)

	queue := []int{start}
	visited[start] = true
	level[start] = 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range g.adjList[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				// 인접 노드의 레벨 = 현재 노드 레벨 + 1
				level[neighbor] = level[current] + 1
				queue = append(queue, neighbor)
			}
		}
	}

	return level
}

// ============================================
// DFS 구현
// ============================================

// DFSRecursive 재귀를 이용한 DFS
func (g *Graph) DFSRecursive(start int) []int {
	visited := make(map[int]bool)
	result := []int{}

	// 내부 재귀 함수
	var dfs func(node int)
	dfs = func(node int) {
		// 현재 노드 방문 처리
		visited[node] = true
		result = append(result, node)

		// 인접한 모든 노드에 대해 재귀 호출
		for _, neighbor := range g.adjList[node] {
			if !visited[neighbor] {
				dfs(neighbor)
			}
		}
	}

	dfs(start)
	return result
}

// DFSIterative 스택을 이용한 반복적 DFS
func (g *Graph) DFSIterative(start int) []int {
	visited := make(map[int]bool)
	result := []int{}

	// DFS는 스택을 사용
	stack := []int{start}

	for len(stack) > 0 {
		// 스택에서 맨 위 요소를 꺼냄 (Pop)
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// 이미 방문한 노드는 건너뜀
		if visited[current] {
			continue
		}

		// 현재 노드 방문 처리
		visited[current] = true
		result = append(result, current)

		// 인접 노드를 스택에 추가 (역순으로 추가하면 순서 유지)
		neighbors := g.adjList[current]
		for i := len(neighbors) - 1; i >= 0; i-- {
			if !visited[neighbors[i]] {
				stack = append(stack, neighbors[i])
			}
		}
	}

	return result
}

// ============================================
// 실용적인 예제: 2D 그리드에서 BFS/DFS
// ============================================

// 2D 그리드에서의 이동 방향 (상, 하, 좌, 우)
var directions = [][]int{
	{-1, 0}, // 상
	{1, 0},  // 하
	{0, -1}, // 좌
	{0, 1},  // 우
}

// Point 2D 좌표
type Point struct {
	row, col int
}

// BFSGrid 2D 그리드에서 최단 거리 찾기
// grid: 0은 이동 가능, 1은 벽
// 시작점에서 도착점까지의 최단 거리 반환 (도달 불가시 -1)
func BFSGrid(grid [][]int, start, end Point) int {
	rows := len(grid)
	if rows == 0 {
		return -1
	}
	cols := len(grid[0])

	// 시작점이나 도착점이 벽이면 불가능
	if grid[start.row][start.col] == 1 || grid[end.row][end.col] == 1 {
		return -1
	}

	// 방문 여부 확인
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	// BFS를 위한 큐 (좌표와 거리를 함께 저장)
	type Node struct {
		point    Point
		distance int
	}

	queue := []Node{{start, 0}}
	visited[start.row][start.col] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// 도착점에 도달하면 거리 반환
		if current.point.row == end.row && current.point.col == end.col {
			return current.distance
		}

		// 4방향으로 이동
		for _, dir := range directions {
			newRow := current.point.row + dir[0]
			newCol := current.point.col + dir[1]

			// 유효한 좌표인지 확인
			if newRow >= 0 && newRow < rows &&
				newCol >= 0 && newCol < cols &&
				!visited[newRow][newCol] &&
				grid[newRow][newCol] == 0 {

				visited[newRow][newCol] = true
				queue = append(queue, Node{
					point:    Point{newRow, newCol},
					distance: current.distance + 1,
				})
			}
		}
	}

	return -1 // 도달 불가능
}

// DFSGrid 2D 그리드에서 경로 존재 여부 확인
func DFSGrid(grid [][]int, start, end Point) bool {
	rows := len(grid)
	if rows == 0 {
		return false
	}
	cols := len(grid[0])

	if grid[start.row][start.col] == 1 || grid[end.row][end.col] == 1 {
		return false
	}

	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	// 재귀 DFS
	var dfs func(row, col int) bool
	dfs = func(row, col int) bool {
		// 도착점에 도달
		if row == end.row && col == end.col {
			return true
		}

		visited[row][col] = true

		// 4방향으로 이동 시도
		for _, dir := range directions {
			newRow := row + dir[0]
			newCol := col + dir[1]

			if newRow >= 0 && newRow < rows &&
				newCol >= 0 && newCol < cols &&
				!visited[newRow][newCol] &&
				grid[newRow][newCol] == 0 {

				if dfs(newRow, newCol) {
					return true
				}
			}
		}

		return false
	}

	return dfs(start.row, start.col)
}

// ============================================
// 연결 요소 찾기 (Connected Components)
// ============================================

// CountConnectedComponents 연결 요소의 개수 찾기
func (g *Graph) CountConnectedComponents() int {
	visited := make(map[int]bool)
	count := 0

	// 모든 정점에 대해 DFS 수행
	for v := 0; v < g.vertices; v++ {
		if !visited[v] {
			// 새로운 연결 요소 발견
			count++
			g.dfsHelper(v, visited)
		}
	}

	return count
}

// dfsHelper DFS 헬퍼 함수
func (g *Graph) dfsHelper(node int, visited map[int]bool) {
	visited[node] = true
	for _, neighbor := range g.adjList[node] {
		if !visited[neighbor] {
			g.dfsHelper(neighbor, visited)
		}
	}
}

// ============================================
// 사이클 감지 (방향 그래프)
// ============================================

// HasCycle 방향 그래프에서 사이클 존재 여부 확인
func (g *Graph) HasCycle() bool {
	visited := make(map[int]bool)
	recStack := make(map[int]bool) // 현재 재귀 스택에 있는 노드

	var hasCycleDFS func(node int) bool
	hasCycleDFS = func(node int) bool {
		visited[node] = true
		recStack[node] = true

		for _, neighbor := range g.adjList[node] {
			// 방문하지 않은 노드에서 사이클 발견
			if !visited[neighbor] {
				if hasCycleDFS(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				// 현재 재귀 스택에 있는 노드를 다시 방문 = 사이클
				return true
			}
		}

		// 현재 노드를 재귀 스택에서 제거
		recStack[node] = false
		return false
	}

	// 모든 노드에서 사이클 검사
	for v := 0; v < g.vertices; v++ {
		if !visited[v] {
			if hasCycleDFS(v) {
				return true
			}
		}
	}

	return false
}

func main() {
	fmt.Println("============================================")
	fmt.Println("6-1. BFS / DFS (그래프 탐색)")
	fmt.Println("============================================")

	// 예제 그래프 생성
	/*
	   그래프 구조:
	       0 --- 1 --- 2
	       |     |
	       3 --- 4 --- 5
	*/
	g := NewGraph(6)
	g.AddEdge(0, 1)
	g.AddEdge(0, 3)
	g.AddEdge(1, 2)
	g.AddEdge(1, 4)
	g.AddEdge(3, 4)
	g.AddEdge(4, 5)

	fmt.Println("\n[그래프 구조]")
	fmt.Println("    0 --- 1 --- 2")
	fmt.Println("    |     |")
	fmt.Println("    3 --- 4 --- 5")

	// BFS 테스트
	fmt.Println("\n[BFS 탐색 결과]")
	fmt.Printf("시작점 0: %v\n", g.BFS(0))

	// 레벨별 BFS
	fmt.Println("\n[BFS 레벨별 탐색]")
	levels := g.BFSWithLevel(0)
	for node, level := range levels {
		fmt.Printf("노드 %d: 레벨 %d\n", node, level)
	}

	// DFS 테스트
	fmt.Println("\n[DFS 재귀 탐색 결과]")
	fmt.Printf("시작점 0: %v\n", g.DFSRecursive(0))

	fmt.Println("\n[DFS 반복 탐색 결과]")
	fmt.Printf("시작점 0: %v\n", g.DFSIterative(0))

	// 2D 그리드 예제
	fmt.Println("\n============================================")
	fmt.Println("2D 그리드에서의 BFS/DFS")
	fmt.Println("============================================")

	grid := [][]int{
		{0, 0, 0, 0, 0},
		{0, 1, 1, 1, 0},
		{0, 0, 0, 1, 0},
		{1, 1, 0, 0, 0},
		{0, 0, 0, 1, 0},
	}

	fmt.Println("\n[그리드] (0: 이동 가능, 1: 벽)")
	for _, row := range grid {
		fmt.Println(row)
	}

	start := Point{0, 0}
	end := Point{4, 4}

	fmt.Printf("\n시작점: (%d, %d), 도착점: (%d, %d)\n",
		start.row, start.col, end.row, end.col)

	distance := BFSGrid(grid, start, end)
	if distance != -1 {
		fmt.Printf("BFS 최단 거리: %d\n", distance)
	} else {
		fmt.Println("BFS: 도달 불가능")
	}

	pathExists := DFSGrid(grid, start, end)
	fmt.Printf("DFS 경로 존재 여부: %v\n", pathExists)

	// 연결 요소 예제
	fmt.Println("\n============================================")
	fmt.Println("연결 요소 찾기")
	fmt.Println("============================================")

	/*
	   분리된 그래프:
	   0 - 1    2 - 3    4
	*/
	g2 := NewGraph(5)
	g2.AddEdge(0, 1)
	g2.AddEdge(2, 3)

	fmt.Println("\n[분리된 그래프]")
	fmt.Println("0 - 1    2 - 3    4 (독립)")
	fmt.Printf("연결 요소 개수: %d\n", g2.CountConnectedComponents())

	// 사이클 감지 예제
	fmt.Println("\n============================================")
	fmt.Println("사이클 감지 (방향 그래프)")
	fmt.Println("============================================")

	// 사이클이 있는 방향 그래프: 0 -> 1 -> 2 -> 0
	g3 := NewGraph(3)
	g3.AddDirectedEdge(0, 1)
	g3.AddDirectedEdge(1, 2)
	g3.AddDirectedEdge(2, 0)

	fmt.Println("\n[방향 그래프] 0 -> 1 -> 2 -> 0")
	fmt.Printf("사이클 존재: %v\n", g3.HasCycle())

	// 사이클이 없는 방향 그래프: 0 -> 1 -> 2
	g4 := NewGraph(3)
	g4.AddDirectedEdge(0, 1)
	g4.AddDirectedEdge(1, 2)

	fmt.Println("\n[방향 그래프] 0 -> 1 -> 2")
	fmt.Printf("사이클 존재: %v\n", g4.HasCycle())

	// BFS vs DFS 비교
	fmt.Println("\n============================================")
	fmt.Println("BFS vs DFS 비교")
	fmt.Println("============================================")
	fmt.Println(`
┌─────────────┬─────────────────────┬─────────────────────┐
│ 특성        │ BFS                 │ DFS                 │
├─────────────┼─────────────────────┼─────────────────────┤
│ 자료구조    │ 큐 (Queue)          │ 스택 (Stack)/재귀   │
│ 탐색 방식   │ 가까운 노드 먼저    │ 깊은 노드 먼저      │
│ 최단 경로   │ 보장 O              │ 보장 X              │
│ 메모리      │ 너비에 비례         │ 깊이에 비례         │
│ 사용 예     │ 최단거리, 레벨탐색  │ 미로, 백트래킹      │
└─────────────┴─────────────────────┴─────────────────────┘
`)
}
