package huhu1

/*
================================================================================
문제 2: 실시간 경로 최적화 시스템 - 솔루션
================================================================================

【 핵심 설계 원칙 】
1. A* 알고리즘 - 휴리스틱으로 탐색 효율화
2. 인접 리스트 기반 그래프 - 메모리 효율적
3. 경로 캐싱 - 자주 사용되는 경로 캐싱
4. Yen's Algorithm - K-shortest paths

【 A* vs Dijkstra 】
- Dijkstra: 모든 방향 탐색, O((V+E)logV)
- A*: 목표 방향 우선 탐색, 휴리스틱으로 가지치기
- 지리적 데이터에서 A*가 2-10배 빠름

================================================================================
*/

import (
	"container/heap"
	"errors"
	"math"
	"sync"
)

// ============================================================================
// 우선순위 큐 구현 (A* 알고리즘용)
// ============================================================================

type routeItem struct {
	nodeID   string
	priority float64 // f(n) = g(n) + h(n)
	gScore   float64 // 실제 비용
	index    int
}

type routePriorityQueue []*routeItem

func (pq routePriorityQueue) Len() int { return len(pq) }

func (pq routePriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq routePriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *routePriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*routeItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *routePriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// ============================================================================
// 그래프 구조
// ============================================================================

// NodeInfo 노드 정보
type NodeInfo struct {
	ID  string
	Lat float64
	Lon float64
}

// EdgeInfo 간선 정보
type EdgeInfo struct {
	To       string
	Weight   float64
	Distance float64
	Blocked  bool
}

// RouteOptimizerSolution 완성된 구현체
type RouteOptimizerSolution struct {
	mu sync.RWMutex

	// 그래프 데이터
	nodes map[string]*NodeInfo
	edges map[string]map[string]*EdgeInfo // from -> to -> edge

	// 경로 캐시
	cacheMu    sync.RWMutex
	routeCache map[string]*Route // "from:to" -> route
}

// NewRouteOptimizerSolution 생성자
func NewRouteOptimizerSolution() *RouteOptimizerSolution {
	return &RouteOptimizerSolution{
		nodes:      make(map[string]*NodeInfo),
		edges:      make(map[string]map[string]*EdgeInfo),
		routeCache: make(map[string]*Route),
	}
}

// ============================================================================
// 그래프 구성
// ============================================================================

// AddNode 노드 추가
func (r *RouteOptimizerSolution) AddNode(node Node) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nodes[node.ID] = &NodeInfo{
		ID:  node.ID,
		Lat: node.Lat,
		Lon: node.Lon,
	}

	if r.edges[node.ID] == nil {
		r.edges[node.ID] = make(map[string]*EdgeInfo)
	}

	return nil
}

// AddEdge 간선 추가
func (r *RouteOptimizerSolution) AddEdge(edge Edge) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.edges[edge.From] == nil {
		r.edges[edge.From] = make(map[string]*EdgeInfo)
	}

	r.edges[edge.From][edge.To] = &EdgeInfo{
		To:       edge.To,
		Weight:   edge.Weight,
		Distance: edge.Distance,
		Blocked:  edge.Blocked,
	}

	// 캐시 무효화
	r.invalidateCache()

	return nil
}

// ============================================================================
// A* 알고리즘
// ============================================================================

/*
【 A* 알고리즘 】
f(n) = g(n) + h(n)
- g(n): 시작점에서 현재 노드까지의 실제 비용
- h(n): 현재 노드에서 목표까지의 휴리스틱 추정 비용
- f(n): 총 예상 비용

휴리스틱으로 Haversine 거리 / 평균 속도 사용
*/

// heuristic A* 휴리스틱 함수 (예상 시간)
func (r *RouteOptimizerSolution) heuristic(from, to string) float64 {
	fromNode := r.nodes[from]
	toNode := r.nodes[to]

	if fromNode == nil || toNode == nil {
		return 0
	}

	// Haversine 거리 (km)
	dist := haversineDistanceRoute(fromNode.Lat, fromNode.Lon, toNode.Lat, toNode.Lon)

	// 평균 속도 60km/h 가정 → 시간(분)으로 변환
	return dist / 60 * 60 // km / (km/h) * 60 = minutes
}

func haversineDistanceRoute(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// FindShortestPath A* 알고리즘으로 최단 경로 탐색
/*
【 시간 복잡도 】O((V+E)logV) 최악, 휴리스틱 효과로 실제로는 훨씬 빠름
【 공간 복잡도 】O(V)
*/
func (r *RouteOptimizerSolution) FindShortestPath(from, to string) (*Route, error) {
	// 캐시 확인
	cacheKey := from + ":" + to
	r.cacheMu.RLock()
	if cached, ok := r.routeCache[cacheKey]; ok {
		r.cacheMu.RUnlock()
		return cached, nil
	}
	r.cacheMu.RUnlock()

	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.nodes[from]; !ok {
		return nil, errors.New("source node not found")
	}
	if _, ok := r.nodes[to]; !ok {
		return nil, errors.New("target node not found")
	}

	// A* 알고리즘
	pq := &routePriorityQueue{}
	heap.Init(pq)

	gScore := make(map[string]float64)
	cameFrom := make(map[string]string)
	distTo := make(map[string]float64) // 거리 추적

	gScore[from] = 0
	distTo[from] = 0

	heap.Push(pq, &routeItem{
		nodeID:   from,
		priority: r.heuristic(from, to),
		gScore:   0,
	})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*routeItem)

		if current.nodeID == to {
			// 경로 재구성
			route := r.reconstructPath(cameFrom, distTo, from, to, gScore[to])

			// 캐시 저장
			r.cacheMu.Lock()
			r.routeCache[cacheKey] = route
			r.cacheMu.Unlock()

			return route, nil
		}

		// 이미 더 좋은 경로를 찾은 경우 스킵
		if g, ok := gScore[current.nodeID]; ok && current.gScore > g {
			continue
		}

		// 인접 노드 탐색
		for _, edge := range r.edges[current.nodeID] {
			if edge.Blocked {
				continue
			}

			tentativeG := gScore[current.nodeID] + edge.Weight

			if g, ok := gScore[edge.To]; !ok || tentativeG < g {
				gScore[edge.To] = tentativeG
				distTo[edge.To] = distTo[current.nodeID] + edge.Distance
				cameFrom[edge.To] = current.nodeID

				heap.Push(pq, &routeItem{
					nodeID:   edge.To,
					priority: tentativeG + r.heuristic(edge.To, to),
					gScore:   tentativeG,
				})
			}
		}
	}

	return nil, errors.New("no path found")
}

// reconstructPath 경로 재구성
func (r *RouteOptimizerSolution) reconstructPath(cameFrom map[string]string, distTo map[string]float64, from, to string, totalTime float64) *Route {
	path := []string{to}
	current := to

	for current != from {
		current = cameFrom[current]
		path = append([]string{current}, path...)
	}

	return &Route{
		Path:      path,
		TotalTime: totalTime,
		TotalDist: distTo[to],
	}
}

// ============================================================================
// 실시간 업데이트
// ============================================================================

// UpdateEdgeWeight 간선 가중치 업데이트
/*
【 시간 복잡도 】O(1)
【 참고 】캐시 무효화로 다음 조회 시 재계산
*/
func (r *RouteOptimizerSolution) UpdateEdgeWeight(from, to string, weight float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.edges[from] == nil || r.edges[from][to] == nil {
		return errors.New("edge not found")
	}

	r.edges[from][to].Weight = weight
	r.invalidateCache()

	return nil
}

// BlockRoad 도로 차단
func (r *RouteOptimizerSolution) BlockRoad(from, to string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.edges[from] == nil || r.edges[from][to] == nil {
		return errors.New("edge not found")
	}

	r.edges[from][to].Blocked = true
	r.invalidateCache()

	return nil
}

// UnblockRoad 도로 차단 해제
func (r *RouteOptimizerSolution) UnblockRoad(from, to string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.edges[from] == nil || r.edges[from][to] == nil {
		return errors.New("edge not found")
	}

	r.edges[from][to].Blocked = false
	r.invalidateCache()

	return nil
}

// invalidateCache 캐시 무효화 (더 정교한 구현 가능)
func (r *RouteOptimizerSolution) invalidateCache() {
	r.cacheMu.Lock()
	r.routeCache = make(map[string]*Route)
	r.cacheMu.Unlock()
}

// ============================================================================
// K-Shortest Paths (Yen's Algorithm)
// ============================================================================

/*
【 Yen's Algorithm 】
1. 첫 번째 최단 경로를 찾는다
2. 이전 경로의 각 노드에서 분기하여 대안 경로를 찾는다
3. 가장 짧은 대안 경로를 선택
4. K개가 될 때까지 반복

【 시간 복잡도 】O(KN(E+NlogN)) where N=노드, E=간선
*/

// GetAlternativeRoutes K개의 대안 경로 반환
func (r *RouteOptimizerSolution) GetAlternativeRoutes(from, to string, k int) ([]Route, error) {
	if k <= 0 {
		return nil, errors.New("k must be positive")
	}

	// 첫 번째 최단 경로
	firstRoute, err := r.FindShortestPath(from, to)
	if err != nil {
		return nil, err
	}

	if k == 1 {
		return []Route{*firstRoute}, nil
	}

	result := []Route{*firstRoute}
	candidates := []Route{}

	for i := 1; i < k; i++ {
		prevRoute := result[i-1]

		// 이전 경로의 각 분기점에서 대안 탐색
		for spurIdx := 0; spurIdx < len(prevRoute.Path)-1; spurIdx++ {
			spurNode := prevRoute.Path[spurIdx]
			rootPath := prevRoute.Path[:spurIdx+1]

			// 이전 경로들에서 사용된 간선들 임시 차단
			blockedEdges := r.blockPreviousPathEdges(result, rootPath)

			// 분기점에서 목적지까지 경로 탐색
			spurRoute, err := r.findPathWithBlocked(spurNode, to, blockedEdges)
			if err == nil {
				// rootPath + spurRoute 합치기
				totalPath := make([]string, len(rootPath)-1)
				copy(totalPath, rootPath[:len(rootPath)-1])
				totalPath = append(totalPath, spurRoute.Path...)

				candidate := Route{
					Path:      totalPath,
					TotalTime: r.calculateRouteTime(totalPath),
					TotalDist: r.calculateRouteDist(totalPath),
				}

				// 중복 체크 후 추가
				if !r.isDuplicateRoute(candidates, candidate) && !r.isDuplicateRoute(result, candidate) {
					candidates = append(candidates, candidate)
				}
			}
		}

		if len(candidates) == 0 {
			break
		}

		// 가장 짧은 후보 선택
		bestIdx := 0
		for j := 1; j < len(candidates); j++ {
			if candidates[j].TotalTime < candidates[bestIdx].TotalTime {
				bestIdx = j
			}
		}

		result = append(result, candidates[bestIdx])
		candidates = append(candidates[:bestIdx], candidates[bestIdx+1:]...)
	}

	return result, nil
}

// blockPreviousPathEdges 이전 경로에서 사용된 간선 반환
func (r *RouteOptimizerSolution) blockPreviousPathEdges(routes []Route, rootPath []string) map[string]bool {
	blocked := make(map[string]bool)

	for _, route := range routes {
		// rootPath와 같은 부분을 가진 경로에서 다음 간선 차단
		if len(route.Path) > len(rootPath) {
			match := true
			for i, node := range rootPath {
				if route.Path[i] != node {
					match = false
					break
				}
			}
			if match {
				key := rootPath[len(rootPath)-1] + ":" + route.Path[len(rootPath)]
				blocked[key] = true
			}
		}
	}

	return blocked
}

// findPathWithBlocked 차단된 간선 제외하고 경로 탐색
func (r *RouteOptimizerSolution) findPathWithBlocked(from, to string, blocked map[string]bool) (*Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pq := &routePriorityQueue{}
	heap.Init(pq)

	gScore := make(map[string]float64)
	cameFrom := make(map[string]string)
	distTo := make(map[string]float64)

	gScore[from] = 0
	distTo[from] = 0

	heap.Push(pq, &routeItem{
		nodeID:   from,
		priority: r.heuristic(from, to),
		gScore:   0,
	})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*routeItem)

		if current.nodeID == to {
			return r.reconstructPath(cameFrom, distTo, from, to, gScore[to]), nil
		}

		if g, ok := gScore[current.nodeID]; ok && current.gScore > g {
			continue
		}

		for _, edge := range r.edges[current.nodeID] {
			if edge.Blocked {
				continue
			}

			// 추가 차단 확인
			edgeKey := current.nodeID + ":" + edge.To
			if blocked[edgeKey] {
				continue
			}

			tentativeG := gScore[current.nodeID] + edge.Weight

			if g, ok := gScore[edge.To]; !ok || tentativeG < g {
				gScore[edge.To] = tentativeG
				distTo[edge.To] = distTo[current.nodeID] + edge.Distance
				cameFrom[edge.To] = current.nodeID

				heap.Push(pq, &routeItem{
					nodeID:   edge.To,
					priority: tentativeG + r.heuristic(edge.To, to),
					gScore:   tentativeG,
				})
			}
		}
	}

	return nil, errors.New("no path found")
}

// calculateRouteTime 경로의 총 시간 계산
func (r *RouteOptimizerSolution) calculateRouteTime(path []string) float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var total float64
	for i := 0; i < len(path)-1; i++ {
		if edge, ok := r.edges[path[i]][path[i+1]]; ok {
			total += edge.Weight
		}
	}
	return total
}

// calculateRouteDist 경로의 총 거리 계산
func (r *RouteOptimizerSolution) calculateRouteDist(path []string) float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var total float64
	for i := 0; i < len(path)-1; i++ {
		if edge, ok := r.edges[path[i]][path[i+1]]; ok {
			total += edge.Distance
		}
	}
	return total
}

// isDuplicateRoute 중복 경로 확인
func (r *RouteOptimizerSolution) isDuplicateRoute(routes []Route, route Route) bool {
	for _, existing := range routes {
		if len(existing.Path) != len(route.Path) {
			continue
		}
		same := true
		for i := range existing.Path {
			if existing.Path[i] != route.Path[i] {
				same = false
				break
			}
		}
		if same {
			return true
		}
	}
	return false
}

// ============================================================================
// 성능 최적화 포인트
// ============================================================================

/*
【 프로덕션 최적화 】

1. 계층적 경로 탐색 (Hierarchical Routing)
   - 먼 거리: 고속도로 네트워크로 탐색
   - 가까운 거리: 일반 도로 탐색
   - 탐색 공간 10-100배 감소

2. 경로 캐싱 전략
   - LRU 캐시로 자주 조회되는 경로 캐싱
   - 부분 경로 캐싱 (A→B 캐시가 있으면 A→C 계산 시 활용)

3. 증분 업데이트
   - 전체 캐시 무효화 대신 영향받는 경로만 업데이트
   - 변경된 간선을 사용하는 경로만 재계산

4. 병렬 처리
   - 여러 경로 요청 동시 처리
   - 대안 경로 탐색 병렬화

┌─────────────────────────────────────────────────────────────────────────────┐
│                        프로덕션 아키텍처                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  [교통 정보 수집] → [Kafka] → [Graph Updater] → [Graph Store (Redis)]      │
│                                                                             │
│  [경로 요청] → [API Gateway] → [Route Service] → [A* Engine]               │
│                                    ↓                                        │
│                            [Route Cache (Redis)]                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
*/
