package medium

/*
================================================================================
문제 9: Clone Graph (그래프 깊은 복사)
================================================================================

난이도: Medium
주제: 그래프, DFS/BFS, 해시맵

【 문제 설명 】
연결된 무방향 그래프의 깊은 복사본을 만드세요.

【 노드 정의 】
type GraphNode struct {
    Val       int
    Neighbors []*GraphNode
}

【 예시 】
입력: adjList = [[2,4],[1,3],[2,4],[1,3]]
     1 -- 2
     |    |
     4 -- 3

출력: 완전히 새로운 노드들로 구성된 동일한 구조의 그래프

【 제약 조건 】
- 노드 수: 0 ~ 100
- 1 <= Node.val <= 100
- 모든 노드 값은 유일
- 자기 자신으로의 간선 없음
- 중복 간선 없음
- 그래프는 연결되어 있음

【 힌트 】
1. HashMap으로 원본→복사본 매핑 유지
2. DFS 또는 BFS로 순회하며 복사

【 주의 】
- 순환 참조 처리 필수
- 이미 방문한 노드 재방문 방지
================================================================================
*/

// GraphNode 그래프 노드
type GraphNode struct {
	Val       int
	Neighbors []*GraphNode
}

// CloneGraph 그래프 깊은 복사
func CloneGraph(node *GraphNode) *GraphNode {
	// 여기에 코드를 작성하세요
	return nil
}

// 테스트
func main() {
	// 그래프 생성: 1-2-3-4-1 (사이클)
	n1 := &GraphNode{Val: 1}
	n2 := &GraphNode{Val: 2}
	n3 := &GraphNode{Val: 3}
	n4 := &GraphNode{Val: 4}
	n1.Neighbors = []*GraphNode{n2, n4}
	n2.Neighbors = []*GraphNode{n1, n3}
	n3.Neighbors = []*GraphNode{n2, n4}
	n4.Neighbors = []*GraphNode{n1, n3}

	clone := CloneGraph(n1)
	if clone != nil {
		println("Clone val:", clone.Val)
		println("Clone is different object:", clone != n1)
	}
}
