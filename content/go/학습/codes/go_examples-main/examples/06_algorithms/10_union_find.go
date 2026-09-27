package main

import (
	"fmt"
	"sort"
)

/*
===========================================
6-10. 유니온 파인드 (Union-Find / Disjoint Set)
===========================================

유니온 파인드란?
- 서로소 집합(Disjoint Set)을 표현하고 관리하는 자료구조
- 두 가지 핵심 연산:
  1. Find: 요소가 속한 집합의 대표(루트) 찾기
  2. Union: 두 집합을 합치기

최적화 기법:
1. 경로 압축(Path Compression): Find 시 모든 노드를 루트에 직접 연결
2. 랭크/크기 기반 합치기: 작은 트리를 큰 트리 아래에 붙임

시간복잡도:
- 기본: O(n)
- 최적화 후: O(α(n)) ≈ O(1) (α는 역 아커만 함수)

사용 사례:
- 연결 요소 찾기
- 사이클 감지
- 최소 스패닝 트리 (Kruskal 알고리즘)
- 동적 연결성 문제
*/

// ============================================
// 기본 유니온 파인드
// ============================================

// UnionFind 유니온 파인드 구조체
type UnionFind struct {
	parent []int // parent[i]: i의 부모 노드
	rank   []int // rank[i]: i를 루트로 하는 트리의 랭크 (높이의 상한)
	count  int   // 집합(연결 요소)의 개수
}

// NewUnionFind 새 유니온 파인드 생성
func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	rank := make([]int, n)

	// 초기에는 각 요소가 자신만의 집합
	for i := 0; i < n; i++ {
		parent[i] = i // 자기 자신이 부모
		rank[i] = 0   // 초기 랭크 0
	}

	return &UnionFind{
		parent: parent,
		rank:   rank,
		count:  n,
	}
}

// Find 요소 x가 속한 집합의 루트 찾기
// 경로 압축(Path Compression) 적용
func (uf *UnionFind) Find(x int) int {
	// 루트가 아니면 재귀적으로 루트 찾기
	if uf.parent[x] != x {
		// 경로 압축: 부모를 루트로 직접 설정
		uf.parent[x] = uf.Find(uf.parent[x])
	}
	return uf.parent[x]
}

// FindIterative 반복적 Find (스택 오버플로우 방지)
func (uf *UnionFind) FindIterative(x int) int {
	root := x
	// 루트 찾기
	for uf.parent[root] != root {
		root = uf.parent[root]
	}

	// 경로 압축: 경로상의 모든 노드를 루트에 직접 연결
	for uf.parent[x] != root {
		next := uf.parent[x]
		uf.parent[x] = root
		x = next
	}

	return root
}

// Union 두 집합 합치기
// 랭크 기반 합치기(Union by Rank) 적용
func (uf *UnionFind) Union(x, y int) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)

	// 이미 같은 집합
	if rootX == rootY {
		return false
	}

	// 랭크가 낮은 트리를 높은 트리 아래에 붙임
	if uf.rank[rootX] < uf.rank[rootY] {
		uf.parent[rootX] = rootY
	} else if uf.rank[rootX] > uf.rank[rootY] {
		uf.parent[rootY] = rootX
	} else {
		// 랭크가 같으면 하나를 선택하고 랭크 증가
		uf.parent[rootY] = rootX
		uf.rank[rootX]++
	}

	uf.count-- // 집합 개수 감소
	return true
}

// Connected 두 요소가 같은 집합인지 확인
func (uf *UnionFind) Connected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}

// Count 현재 집합의 개수 반환
func (uf *UnionFind) Count() int {
	return uf.count
}

// ============================================
// 크기 기반 유니온 파인드
// ============================================

// UnionFindWithSize 크기 정보를 유지하는 유니온 파인드
type UnionFindWithSize struct {
	parent []int
	size   []int // size[i]: i를 루트로 하는 집합의 크기
	count  int
}

// NewUnionFindWithSize 새 유니온 파인드 생성 (크기 기반)
func NewUnionFindWithSize(n int) *UnionFindWithSize {
	parent := make([]int, n)
	size := make([]int, n)

	for i := 0; i < n; i++ {
		parent[i] = i
		size[i] = 1
	}

	return &UnionFindWithSize{
		parent: parent,
		size:   size,
		count:  n,
	}
}

func (uf *UnionFindWithSize) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x])
	}
	return uf.parent[x]
}

// Union 크기 기반 합치기 (Union by Size)
func (uf *UnionFindWithSize) Union(x, y int) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)

	if rootX == rootY {
		return false
	}

	// 작은 집합을 큰 집합에 붙임
	if uf.size[rootX] < uf.size[rootY] {
		uf.parent[rootX] = rootY
		uf.size[rootY] += uf.size[rootX]
	} else {
		uf.parent[rootY] = rootX
		uf.size[rootX] += uf.size[rootY]
	}

	uf.count--
	return true
}

func (uf *UnionFindWithSize) Connected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}

// GetSize 특정 요소가 속한 집합의 크기
func (uf *UnionFindWithSize) GetSize(x int) int {
	return uf.size[uf.Find(x)]
}

func (uf *UnionFindWithSize) Count() int {
	return uf.count
}

// ============================================
// 응용 1: 연결 요소 개수
// ============================================

// CountComponents 그래프의 연결 요소 개수
func CountComponents(n int, edges [][]int) int {
	uf := NewUnionFind(n)

	for _, edge := range edges {
		uf.Union(edge[0], edge[1])
	}

	return uf.Count()
}

// ============================================
// 응용 2: 사이클 감지
// ============================================

// HasCycleUF 무방향 그래프에서 사이클 존재 여부 (유니온 파인드 사용)
func HasCycleUF(n int, edges [][]int) bool {
	uf := NewUnionFind(n)

	for _, edge := range edges {
		// 이미 같은 집합에 있으면 사이클
		if uf.Connected(edge[0], edge[1]) {
			return true
		}
		uf.Union(edge[0], edge[1])
	}

	return false
}

// ============================================
// 응용 3: 최소 스패닝 트리 (Kruskal)
// ============================================

// Edge 가중치 간선
type Edge struct {
	From, To, Weight int
}

// KruskalMST 크러스컬 알고리즘으로 MST 찾기
// 시간복잡도: O(E log E)
func KruskalMST(n int, edges []Edge) (int, []Edge) {
	// 간선을 가중치 기준으로 정렬
	sortedEdges := make([]Edge, len(edges))
	copy(sortedEdges, edges)
	sort.Slice(sortedEdges, func(i, j int) bool {
		return sortedEdges[i].Weight < sortedEdges[j].Weight
	})

	uf := NewUnionFind(n)
	mstWeight := 0
	mstEdges := []Edge{}

	for _, edge := range sortedEdges {
		// 사이클을 형성하지 않으면 MST에 추가
		if uf.Union(edge.From, edge.To) {
			mstWeight += edge.Weight
			mstEdges = append(mstEdges, edge)

			// MST 완성 (n-1개의 간선)
			if len(mstEdges) == n-1 {
				break
			}
		}
	}

	return mstWeight, mstEdges
}

// ============================================
// 응용 4: 친구 네트워크
// ============================================

// FriendNetwork 친구 네트워크의 최대 크기
type FriendNetwork struct {
	uf      *UnionFindWithSize
	nameToId map[string]int
	nextId   int
}

func NewFriendNetwork() *FriendNetwork {
	return &FriendNetwork{
		uf:      NewUnionFindWithSize(100000),
		nameToId: make(map[string]int),
		nextId:   0,
	}
}

// AddFriendship 친구 관계 추가, 해당 네트워크 크기 반환
func (fn *FriendNetwork) AddFriendship(name1, name2 string) int {
	id1 := fn.getOrCreateId(name1)
	id2 := fn.getOrCreateId(name2)

	fn.uf.Union(id1, id2)

	return fn.uf.GetSize(id1)
}

func (fn *FriendNetwork) getOrCreateId(name string) int {
	if id, exists := fn.nameToId[name]; exists {
		return id
	}
	fn.nameToId[name] = fn.nextId
	fn.nextId++
	return fn.nameToId[name]
}

// ============================================
// 응용 5: 섬의 개수 (2D 그리드)
// ============================================

// NumIslands 섬의 개수 (유니온 파인드 사용)
func NumIslands(grid [][]byte) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}

	rows, cols := len(grid), len(grid[0])

	// 2D 좌표를 1D 인덱스로 변환
	getId := func(r, c int) int {
		return r*cols + c
	}

	// 육지 개수 세기 및 유니온 파인드 초기화
	landCount := 0
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if grid[i][j] == '1' {
				landCount++
			}
		}
	}

	uf := NewUnionFind(rows * cols)

	directions := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if grid[i][j] == '1' {
				for _, dir := range directions {
					ni, nj := i+dir[0], j+dir[1]
					if ni >= 0 && ni < rows && nj >= 0 && nj < cols && grid[ni][nj] == '1' {
						if uf.Union(getId(i, j), getId(ni, nj)) {
							landCount-- // 연결되면 섬 개수 감소
						}
					}
				}
			}
		}
	}

	return landCount
}

// ============================================
// 응용 6: 그래프 유효 트리 판별
// ============================================

// IsValidTree n개 노드와 edges로 구성된 그래프가 유효한 트리인지
// 트리 조건: 연결되어 있고 사이클이 없음 (간선 수 = n-1)
func IsValidTree(n int, edges [][]int) bool {
	// 트리는 정확히 n-1개의 간선을 가짐
	if len(edges) != n-1 {
		return false
	}

	uf := NewUnionFind(n)

	for _, edge := range edges {
		// 사이클이 생기면 트리 아님
		if !uf.Union(edge[0], edge[1]) {
			return false
		}
	}

	// 모든 노드가 연결되어야 함
	return uf.Count() == 1
}

// ============================================
// 응용 7: 중복 간선 찾기
// ============================================

// FindRedundantConnection 중복 간선 찾기
// 트리에 하나의 추가 간선이 있을 때, 제거하면 트리가 되는 간선 반환
func FindRedundantConnection(edges [][]int) []int {
	n := len(edges)
	uf := NewUnionFind(n + 1) // 1-indexed

	for _, edge := range edges {
		// 이미 연결되어 있으면 중복 간선
		if !uf.Union(edge[0], edge[1]) {
			return edge
		}
	}

	return nil
}

// ============================================
// 응용 8: 계정 병합
// ============================================

// AccountsMerge 같은 이메일을 가진 계정 병합
func AccountsMerge(accounts [][]string) [][]string {
	// 이메일 -> 첫 번째 계정 인덱스
	emailToIndex := make(map[string]int)
	// 이메일 -> 계정 이름
	emailToName := make(map[string]string)

	uf := NewUnionFind(len(accounts))

	for i, account := range accounts {
		name := account[0]
		for j := 1; j < len(account); j++ {
			email := account[j]
			emailToName[email] = name

			if firstIndex, exists := emailToIndex[email]; exists {
				// 같은 이메일이 있으면 계정 병합
				uf.Union(i, firstIndex)
			} else {
				emailToIndex[email] = i
			}
		}
	}

	// 루트별로 이메일 수집
	indexToEmails := make(map[int][]string)
	for email, index := range emailToIndex {
		root := uf.Find(index)
		indexToEmails[root] = append(indexToEmails[root], email)
	}

	// 결과 생성
	result := [][]string{}
	for root, emails := range indexToEmails {
		sort.Strings(emails)
		name := emailToName[emails[0]]
		merged := append([]string{name}, emails...)
		_ = root // unused
		result = append(result, merged)
	}

	return result
}

func main() {
	fmt.Println("============================================")
	fmt.Println("6-10. 유니온 파인드 (Union-Find)")
	fmt.Println("============================================")

	// 기본 사용
	fmt.Println("\n[기본 유니온 파인드]")
	fmt.Println("--------------------------------------------")
	uf := NewUnionFind(6)
	fmt.Println("6개 원소로 초기화 (0-5)")

	fmt.Println("\n연산 수행:")
	fmt.Println("Union(0, 1)")
	uf.Union(0, 1)

	fmt.Println("Union(2, 3)")
	uf.Union(2, 3)

	fmt.Println("Union(4, 5)")
	uf.Union(4, 5)

	fmt.Printf("\n현재 집합 개수: %d\n", uf.Count())
	fmt.Printf("Connected(0, 1): %v\n", uf.Connected(0, 1))
	fmt.Printf("Connected(0, 2): %v\n", uf.Connected(0, 2))

	fmt.Println("\nUnion(1, 3)")
	uf.Union(1, 3)

	fmt.Printf("\n현재 집합 개수: %d\n", uf.Count())
	fmt.Printf("Connected(0, 2): %v (이제 연결됨)\n", uf.Connected(0, 2))

	// 연결 요소 개수
	fmt.Println("\n[연결 요소 개수]")
	fmt.Println("--------------------------------------------")
	edges := [][]int{{0, 1}, {1, 2}, {3, 4}}
	fmt.Printf("n=5, edges=%v\n", edges)
	fmt.Printf("연결 요소 개수: %d\n", CountComponents(5, edges))

	// 사이클 감지
	fmt.Println("\n[사이클 감지]")
	fmt.Println("--------------------------------------------")
	edges1 := [][]int{{0, 1}, {1, 2}, {2, 0}}
	fmt.Printf("edges=%v: 사이클 존재=%v\n", edges1, HasCycleUF(3, edges1))

	edges2 := [][]int{{0, 1}, {1, 2}}
	fmt.Printf("edges=%v: 사이클 존재=%v\n", edges2, HasCycleUF(3, edges2))

	// 크러스컬 MST
	fmt.Println("\n[최소 스패닝 트리 (Kruskal)]")
	fmt.Println("--------------------------------------------")
	graphEdges := []Edge{
		{0, 1, 4},
		{0, 7, 8},
		{1, 2, 8},
		{1, 7, 11},
		{2, 3, 7},
		{2, 5, 4},
		{2, 8, 2},
		{3, 4, 9},
		{3, 5, 14},
		{4, 5, 10},
		{5, 6, 2},
		{6, 7, 1},
		{6, 8, 6},
		{7, 8, 7},
	}

	mstWeight, mstEdges := KruskalMST(9, graphEdges)
	fmt.Printf("MST 총 가중치: %d\n", mstWeight)
	fmt.Println("MST 간선:")
	for _, e := range mstEdges {
		fmt.Printf("  %d -- %d (가중치: %d)\n", e.From, e.To, e.Weight)
	}

	// 친구 네트워크
	fmt.Println("\n[친구 네트워크]")
	fmt.Println("--------------------------------------------")
	fn := NewFriendNetwork()

	friendships := [][]string{
		{"Alice", "Bob"},
		{"Charlie", "David"},
		{"Alice", "Charlie"},
		{"Eve", "Frank"},
	}

	for _, f := range friendships {
		size := fn.AddFriendship(f[0], f[1])
		fmt.Printf("%s와 %s 친구 맺음 → 네트워크 크기: %d\n", f[0], f[1], size)
	}

	// 그래프 유효 트리
	fmt.Println("\n[그래프 유효 트리 판별]")
	fmt.Println("--------------------------------------------")
	treeEdges := [][]int{{0, 1}, {0, 2}, {0, 3}, {1, 4}}
	fmt.Printf("n=5, edges=%v\n", treeEdges)
	fmt.Printf("유효한 트리: %v\n", IsValidTree(5, treeEdges))

	notTreeEdges := [][]int{{0, 1}, {1, 2}, {2, 3}, {1, 3}, {1, 4}}
	fmt.Printf("\nn=5, edges=%v\n", notTreeEdges)
	fmt.Printf("유효한 트리: %v\n", IsValidTree(5, notTreeEdges))

	// 중복 간선 찾기
	fmt.Println("\n[중복 간선 찾기]")
	fmt.Println("--------------------------------------------")
	redundantEdges := [][]int{{1, 2}, {1, 3}, {2, 3}}
	fmt.Printf("edges=%v\n", redundantEdges)
	fmt.Printf("중복 간선: %v\n", FindRedundantConnection(redundantEdges))

	// 섬의 개수
	fmt.Println("\n[섬의 개수]")
	fmt.Println("--------------------------------------------")
	grid := [][]byte{
		{'1', '1', '0', '0', '0'},
		{'1', '1', '0', '0', '0'},
		{'0', '0', '1', '0', '0'},
		{'0', '0', '0', '1', '1'},
	}
	fmt.Println("그리드:")
	for _, row := range grid {
		fmt.Printf("  %s\n", string(row))
	}
	fmt.Printf("섬의 개수: %d\n", NumIslands(grid))

	// 패턴 요약
	fmt.Println("\n============================================")
	fmt.Println("유니온 파인드 패턴 요약")
	fmt.Println("============================================")
	fmt.Println(`
┌─────────────────────────────────────────────────────────────┐
│ 유니온 파인드 기본 구조                                      │
├─────────────────────────────────────────────────────────────┤
│ type UnionFind struct {                                     │
│     parent []int  // parent[i] = i의 부모                   │
│     rank   []int  // rank[i] = 트리 높이 (선택적)           │
│ }                                                           │
│                                                             │
│ Find(x): 경로 압축으로 루트 찾기                            │
│ Union(x, y): 랭크/크기 기반으로 합치기                      │
│ Connected(x, y): Find(x) == Find(y)                         │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ 최적화 기법                                                  │
├─────────────────────────────────────────────────────────────┤
│ 1. 경로 압축 (Path Compression)                             │
│    - Find 시 모든 노드를 루트에 직접 연결                   │
│    - parent[x] = Find(parent[x])                            │
│                                                             │
│ 2. 랭크/크기 기반 합치기                                    │
│    - 작은 트리를 큰 트리 아래에 붙임                        │
│    - 트리 높이 증가 최소화                                  │
└─────────────────────────────────────────────────────────────┘

[사용 사례]
- 연결 요소 개수/판별
- 사이클 감지 (무방향 그래프)
- 최소 스패닝 트리 (Kruskal)
- 동적 연결성 쿼리
- 네트워크 연결 문제
- 이미지 세그멘테이션
`)
}
