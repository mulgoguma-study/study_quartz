package main

import (
	"container/heap"
	"fmt"
	"math"
)

/*
===========================================
6-7. 다익스트라 알고리즘 (Dijkstra's Algorithm)
===========================================

다익스트라란?
- 가중치가 있는 그래프에서 단일 출발점 최단 경로를 찾는 알고리즘
- 음수 가중치가 없는 그래프에서 사용 가능
- 그리디 + 우선순위 큐 사용

시간복잡도:
- 배열 사용: O(V²)
- 우선순위 큐 사용: O((V + E) log V)
- 피보나치 힙 사용: O(E + V log V)

사용 사례:
- GPS 네비게이션 (최단 경로)
- 네트워크 라우팅 프로토콜
- 게임 AI 경로 찾기
*/

// ============================================
// 우선순위 큐 구현 (container/heap 사용)
// ============================================

// PQItem 우선순위 큐 아이템
type PQItem struct {
	node     int // 노드 번호
	distance int // 출발점으로부터의 거리
	index    int // 힙에서의 인덱스 (update용)
}

// PriorityQueue 최소 힙 기반 우선순위 큐
type PriorityQueue []*PQItem

func (pq PriorityQueue) Len() int { return len(pq) }

// Less: 거리가 작은 것이 높은 우선순위
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PQItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[:n-1]
	return item
}

// ============================================
// 가중치 그래프 정의
// ============================================

// Edge 가중치가 있는 간선
type Edge struct {
	To     int // 목적지 노드
	Weight int // 가중치
}

// WeightedGraph 가중치 그래프 (인접 리스트)
type WeightedGraph struct {
	Vertices int
	AdjList  map[int][]Edge
}

// NewWeightedGraph 새 가중치 그래프 생성
func NewWeightedGraph(vertices int) *WeightedGraph {
	return &WeightedGraph{
		Vertices: vertices,
		AdjList:  make(map[int][]Edge),
	}
}

// AddEdge 방향 간선 추가
func (g *WeightedGraph) AddEdge(from, to, weight int) {
	g.AdjList[from] = append(g.AdjList[from], Edge{to, weight})
}

// AddUndirectedEdge 무방향 간선 추가
func (g *WeightedGraph) AddUndirectedEdge(from, to, weight int) {
	g.AdjList[from] = append(g.AdjList[from], Edge{to, weight})
	g.AdjList[to] = append(g.AdjList[to], Edge{from, weight})
}

// ============================================
// 다익스트라 알고리즘 구현
// ============================================

// Dijkstra 우선순위 큐를 사용한 다익스트라
// 시간복잡도: O((V + E) log V)
func (g *WeightedGraph) Dijkstra(start int) ([]int, []int) {
	// 거리 배열: 무한대로 초기화
	dist := make([]int, g.Vertices)
	for i := range dist {
		dist[i] = math.MaxInt32
	}
	dist[start] = 0

	// 이전 노드 배열: 경로 추적용
	prev := make([]int, g.Vertices)
	for i := range prev {
		prev[i] = -1
	}

	// 우선순위 큐 초기화
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &PQItem{node: start, distance: 0})

	// 방문 여부 (선택적 최적화)
	visited := make([]bool, g.Vertices)

	for pq.Len() > 0 {
		// 가장 가까운 노드 꺼내기
		item := heap.Pop(pq).(*PQItem)
		u := item.node

		// 이미 처리된 노드는 건너뜀
		if visited[u] {
			continue
		}
		visited[u] = true

		// 인접 노드들에 대해 거리 갱신 (Relaxation)
		for _, edge := range g.AdjList[u] {
			v := edge.To
			weight := edge.Weight

			// 더 짧은 경로를 발견하면 갱신
			newDist := dist[u] + weight
			if newDist < dist[v] {
				dist[v] = newDist
				prev[v] = u
				heap.Push(pq, &PQItem{node: v, distance: newDist})
			}
		}
	}

	return dist, prev
}

// GetPath 출발점에서 목적지까지의 경로 재구성
func GetPath(prev []int, dest int) []int {
	path := []int{}
	current := dest

	// 역순으로 경로 추적
	for current != -1 {
		path = append([]int{current}, path...) // 앞에 추가
		current = prev[current]
	}

	return path
}

// DijkstraSimple 단순 배열 기반 다익스트라
// 시간복잡도: O(V²), 간선이 적은 희소 그래프에서는 비효율적
func (g *WeightedGraph) DijkstraSimple(start int) []int {
	dist := make([]int, g.Vertices)
	for i := range dist {
		dist[i] = math.MaxInt32
	}
	dist[start] = 0

	visited := make([]bool, g.Vertices)

	for count := 0; count < g.Vertices; count++ {
		// 방문하지 않은 노드 중 최소 거리 노드 찾기
		u := -1
		minDist := math.MaxInt32
		for v := 0; v < g.Vertices; v++ {
			if !visited[v] && dist[v] < minDist {
				minDist = dist[v]
				u = v
			}
		}

		if u == -1 {
			break // 더 이상 도달 가능한 노드 없음
		}

		visited[u] = true

		// 거리 갱신
		for _, edge := range g.AdjList[u] {
			v := edge.To
			if !visited[v] && dist[u]+edge.Weight < dist[v] {
				dist[v] = dist[u] + edge.Weight
			}
		}
	}

	return dist
}

// ============================================
// 실용적인 예제: 네트워크 지연 시간
// ============================================

// NetworkDelayTime n개 노드 네트워크에서 모든 노드에 신호가 도달하는 최소 시간
// times[i] = [출발, 도착, 시간]
func NetworkDelayTime(times [][]int, n, k int) int {
	// 그래프 구성
	graph := NewWeightedGraph(n + 1) // 1-indexed
	for _, time := range times {
		graph.AddEdge(time[0], time[1], time[2])
	}

	// 다익스트라 실행
	dist, _ := graph.Dijkstra(k)

	// 모든 노드까지의 최대 거리 찾기
	maxTime := 0
	for i := 1; i <= n; i++ {
		if dist[i] == math.MaxInt32 {
			return -1 // 도달 불가능한 노드 있음
		}
		if dist[i] > maxTime {
			maxTime = dist[i]
		}
	}

	return maxTime
}

// ============================================
// 2D 그리드에서의 다익스트라
// ============================================

// GridPQItem 그리드용 우선순위 큐 아이템
type GridPQItem struct {
	row, col int
	distance int
}

// GridPQ 그리드용 우선순위 큐
type GridPQ []GridPQItem

func (pq GridPQ) Len() int            { return len(pq) }
func (pq GridPQ) Less(i, j int) bool  { return pq[i].distance < pq[j].distance }
func (pq GridPQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *GridPQ) Push(x interface{}) { *pq = append(*pq, x.(GridPQItem)) }
func (pq *GridPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

// MinCostPath 가중치가 있는 2D 그리드에서 최소 비용 경로
// grid[i][j]는 해당 셀을 방문하는 비용
func MinCostPath(grid [][]int) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}

	rows, cols := len(grid), len(grid[0])
	directions := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	// 거리 배열
	dist := make([][]int, rows)
	for i := range dist {
		dist[i] = make([]int, cols)
		for j := range dist[i] {
			dist[i][j] = math.MaxInt32
		}
	}
	dist[0][0] = grid[0][0]

	// 우선순위 큐
	pq := &GridPQ{}
	heap.Init(pq)
	heap.Push(pq, GridPQItem{0, 0, grid[0][0]})

	for pq.Len() > 0 {
		item := heap.Pop(pq).(GridPQItem)
		row, col, d := item.row, item.col, item.distance

		// 이미 더 짧은 경로로 방문했으면 건너뜀
		if d > dist[row][col] {
			continue
		}

		// 도착점 도달
		if row == rows-1 && col == cols-1 {
			return d
		}

		// 4방향 탐색
		for _, dir := range directions {
			newRow, newCol := row+dir[0], col+dir[1]

			if newRow >= 0 && newRow < rows && newCol >= 0 && newCol < cols {
				newDist := d + grid[newRow][newCol]
				if newDist < dist[newRow][newCol] {
					dist[newRow][newCol] = newDist
					heap.Push(pq, GridPQItem{newRow, newCol, newDist})
				}
			}
		}
	}

	return dist[rows-1][cols-1]
}

// ============================================
// 다익스트라 변형: K번째 최단 경로
// ============================================

// KthShortestPath K번째 최단 경로 길이 찾기
func (g *WeightedGraph) KthShortestPath(start, end, k int) int {
	// 각 노드에 도달한 횟수
	count := make([]int, g.Vertices)

	// 우선순위 큐: (거리, 노드)
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &PQItem{node: start, distance: 0})

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*PQItem)
		u := item.node
		d := item.distance

		count[u]++

		// end 노드에 k번째 도달하면 반환
		if u == end && count[u] == k {
			return d
		}

		// k번 이상 방문한 노드는 건너뜀
		if count[u] > k {
			continue
		}

		for _, edge := range g.AdjList[u] {
			heap.Push(pq, &PQItem{
				node:     edge.To,
				distance: d + edge.Weight,
			})
		}
	}

	return -1 // k번째 경로 없음
}

// ============================================
// 벨만-포드 알고리즘 (음수 가중치 처리)
// ============================================

// BellmanFord 음수 가중치가 있는 그래프에서 최단 경로
// 시간복잡도: O(V * E)
// 음수 사이클 존재 여부도 감지
func (g *WeightedGraph) BellmanFord(start int) ([]int, bool) {
	dist := make([]int, g.Vertices)
	for i := range dist {
		dist[i] = math.MaxInt32
	}
	dist[start] = 0

	// V-1번 반복 (최장 경로 길이)
	for i := 0; i < g.Vertices-1; i++ {
		updated := false
		// 모든 간선에 대해 relaxation
		for u := 0; u < g.Vertices; u++ {
			if dist[u] == math.MaxInt32 {
				continue
			}
			for _, edge := range g.AdjList[u] {
				if dist[u]+edge.Weight < dist[edge.To] {
					dist[edge.To] = dist[u] + edge.Weight
					updated = true
				}
			}
		}
		// 갱신이 없으면 조기 종료
		if !updated {
			break
		}
	}

	// 음수 사이클 감지: V번째에도 갱신되면 음수 사이클 존재
	for u := 0; u < g.Vertices; u++ {
		if dist[u] == math.MaxInt32 {
			continue
		}
		for _, edge := range g.AdjList[u] {
			if dist[u]+edge.Weight < dist[edge.To] {
				return nil, true // 음수 사이클 존재
			}
		}
	}

	return dist, false
}

func main() {
	fmt.Println("============================================")
	fmt.Println("6-7. 다익스트라 알고리즘")
	fmt.Println("============================================")

	// 그래프 생성
	/*
	   예제 그래프:
	       0 ---(4)--- 1
	       |         / |
	      (1)    (2)  (5)
	       |   /      |
	       2 ---(3)--- 3
	*/
	g := NewWeightedGraph(4)
	g.AddEdge(0, 1, 4)
	g.AddEdge(0, 2, 1)
	g.AddEdge(1, 3, 5)
	g.AddEdge(2, 1, 2)
	g.AddEdge(2, 3, 3)

	fmt.Println("\n[그래프 구조]")
	fmt.Println("    0 ---(4)--- 1")
	fmt.Println("    |         / |")
	fmt.Println("   (1)    (2)  (5)")
	fmt.Println("    |   /      |")
	fmt.Println("    2 ---(3)--- 3")

	// 다익스트라 실행
	fmt.Println("\n[다익스트라 결과 (시작점: 0)]")
	fmt.Println("--------------------------------------------")
	dist, prev := g.Dijkstra(0)

	for i := 0; i < g.Vertices; i++ {
		path := GetPath(prev, i)
		fmt.Printf("노드 %d까지: 거리=%d, 경로=%v\n", i, dist[i], path)
	}

	// 단순 다익스트라 비교
	fmt.Println("\n[단순 배열 다익스트라 결과]")
	simpleDist := g.DijkstraSimple(0)
	fmt.Printf("거리: %v\n", simpleDist)

	// 네트워크 지연 시간
	fmt.Println("\n============================================")
	fmt.Println("[네트워크 지연 시간]")
	fmt.Println("============================================")

	times := [][]int{{2, 1, 1}, {2, 3, 1}, {3, 4, 1}}
	n, k := 4, 2
	fmt.Printf("간선: %v\n", times)
	fmt.Printf("노드 수: %d, 시작 노드: %d\n", n, k)
	fmt.Printf("모든 노드 도달 시간: %d\n", NetworkDelayTime(times, n, k))

	// 2D 그리드 최소 비용 경로
	fmt.Println("\n============================================")
	fmt.Println("[2D 그리드 최소 비용 경로]")
	fmt.Println("============================================")

	grid := [][]int{
		{1, 3, 1},
		{1, 5, 1},
		{4, 2, 1},
	}
	fmt.Println("그리드:")
	for _, row := range grid {
		fmt.Printf("  %v\n", row)
	}
	fmt.Printf("최소 비용: %d\n", MinCostPath(grid))

	// K번째 최단 경로
	fmt.Println("\n============================================")
	fmt.Println("[K번째 최단 경로]")
	fmt.Println("============================================")

	g2 := NewWeightedGraph(4)
	g2.AddEdge(0, 1, 1)
	g2.AddEdge(0, 2, 3)
	g2.AddEdge(1, 2, 1)
	g2.AddEdge(1, 3, 4)
	g2.AddEdge(2, 3, 1)

	fmt.Println("그래프: 0->1(1), 0->2(3), 1->2(1), 1->3(4), 2->3(1)")
	for k := 1; k <= 3; k++ {
		dist := g2.KthShortestPath(0, 3, k)
		fmt.Printf("%d번째 최단 경로 (0->3): %d\n", k, dist)
	}

	// 벨만-포드 (음수 가중치)
	fmt.Println("\n============================================")
	fmt.Println("[벨만-포드 알고리즘 (음수 가중치)]")
	fmt.Println("============================================")

	g3 := NewWeightedGraph(5)
	g3.AddEdge(0, 1, 6)
	g3.AddEdge(0, 2, 7)
	g3.AddEdge(1, 2, 8)
	g3.AddEdge(1, 3, 5)
	g3.AddEdge(1, 4, -4) // 음수 가중치
	g3.AddEdge(2, 3, -3) // 음수 가중치
	g3.AddEdge(2, 4, 9)
	g3.AddEdge(3, 1, -2) // 음수 가중치
	g3.AddEdge(4, 3, 7)

	fmt.Println("음수 가중치 포함 그래프")
	bfDist, hasNegativeCycle := g3.BellmanFord(0)

	if hasNegativeCycle {
		fmt.Println("음수 사이클 존재!")
	} else {
		fmt.Println("벨만-포드 결과:")
		for i := 0; i < g3.Vertices; i++ {
			fmt.Printf("  노드 %d까지: %d\n", i, bfDist[i])
		}
	}

	// 알고리즘 비교
	fmt.Println("\n============================================")
	fmt.Println("최단 경로 알고리즘 비교")
	fmt.Println("============================================")
	fmt.Println(`
┌─────────────────────────────────────────────────────────────┐
│ 다익스트라 (Dijkstra)                                        │
├─────────────────────────────────────────────────────────────┤
│ - 음수 가중치 X                                              │
│ - 시간: O((V+E) log V) with 우선순위 큐                      │
│ - 단일 출발점 최단 경로                                      │
│ - 실제 활용도 높음 (GPS, 라우팅)                            │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ 벨만-포드 (Bellman-Ford)                                     │
├─────────────────────────────────────────────────────────────┤
│ - 음수 가중치 O                                              │
│ - 음수 사이클 감지 가능                                      │
│ - 시간: O(V * E)                                            │
│ - 다익스트라보다 느림                                        │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ 플로이드-워셜 (Floyd-Warshall)                               │
├─────────────────────────────────────────────────────────────┤
│ - 모든 쌍 최단 경로                                          │
│ - 시간: O(V³)                                               │
│ - 음수 가중치 O (음수 사이클 X)                              │
│ - 밀집 그래프에 적합                                        │
└─────────────────────────────────────────────────────────────┘

[선택 가이드]
- 단일 출발점, 양수 가중치 → 다익스트라
- 단일 출발점, 음수 가중치 → 벨만-포드
- 모든 쌍 최단 경로 → 플로이드-워셜
- 가중치 없는 그래프 → BFS
`)
}
