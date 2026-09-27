package medium

/*
================================================================================
문제 9: Clone Graph - 솔루션
================================================================================
*/

/*
================================================================================
그래프 복사 전략
================================================================================

【 핵심 문제 】
- 순환 참조: 같은 노드를 여러 번 방문하게 됨
- 해결: HashMap으로 "이미 복사한 노드" 추적

┌─────────────────────────────────────────────────────────────────────────────┐
│                        복사 과정 시각화                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  원본          HashMap              복사본                                  │
│                                                                             │
│   1 ──────────→ 1 → 1' ──────────→ 1'                                      │
│   │             ↑                   │                                       │
│   │   (이미 복사됨, 재사용)         │                                       │
│   ▼                                 ▼                                       │
│   2 ──────────→ 2 → 2' ──────────→ 2'                                      │
│   │                                 │                                       │
│   ▼                                 ▼                                       │
│   3 ──────────→ 3 → 3' ──────────→ 3'                                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
*/

// GraphNode 그래프 노드
type GraphNode struct {
	Val       int
	Neighbors []*GraphNode
}

// ============================================================================
// 방법 1: DFS (재귀)
// ============================================================================

func CloneGraph(node *GraphNode) *GraphNode {
	if node == nil {
		return nil
	}

	// 원본 → 복사본 매핑
	visited := make(map[*GraphNode]*GraphNode)
	return cloneDFS(node, visited)
}

func cloneDFS(node *GraphNode, visited map[*GraphNode]*GraphNode) *GraphNode {
	// 이미 복사한 노드면 복사본 반환
	if clone, exists := visited[node]; exists {
		return clone
	}

	// 새 노드 생성
	clone := &GraphNode{Val: node.Val}
	visited[node] = clone // 먼저 등록 (순환 참조 방지!)

	// 이웃 노드들 재귀적으로 복사
	for _, neighbor := range node.Neighbors {
		clone.Neighbors = append(clone.Neighbors, cloneDFS(neighbor, visited))
	}

	return clone
}

// ============================================================================
// 방법 2: BFS
// ============================================================================

func CloneGraphBFS(node *GraphNode) *GraphNode {
	if node == nil {
		return nil
	}

	visited := make(map[*GraphNode]*GraphNode)
	queue := []*GraphNode{node}

	// 시작 노드 복사
	visited[node] = &GraphNode{Val: node.Val}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		// 현재 노드의 모든 이웃 처리
		for _, neighbor := range curr.Neighbors {
			if _, exists := visited[neighbor]; !exists {
				// 새 노드 복사 및 큐에 추가
				visited[neighbor] = &GraphNode{Val: neighbor.Val}
				queue = append(queue, neighbor)
			}
			// 복사본에 이웃 연결
			visited[curr].Neighbors = append(
				visited[curr].Neighbors,
				visited[neighbor],
			)
		}
	}

	return visited[node]
}

/*
================================================================================
관련 문제들
================================================================================

【 1. Deep Copy of Linked List with Random Pointer 】

type NodeWithRandom struct {
    Val    int
    Next   *NodeWithRandom
    Random *NodeWithRandom
}

func CopyRandomList(head *NodeWithRandom) *NodeWithRandom {
    if head == nil {
        return nil
    }

    // 1단계: 각 노드 뒤에 복사본 삽입
    // A -> A' -> B -> B' -> C -> C'
    curr := head
    for curr != nil {
        clone := &NodeWithRandom{Val: curr.Val}
        clone.Next = curr.Next
        curr.Next = clone
        curr = clone.Next
    }

    // 2단계: Random 포인터 복사
    curr = head
    for curr != nil {
        if curr.Random != nil {
            curr.Next.Random = curr.Random.Next
        }
        curr = curr.Next.Next
    }

    // 3단계: 리스트 분리
    dummy := &NodeWithRandom{}
    cloneCurr := dummy
    curr = head
    for curr != nil {
        cloneCurr.Next = curr.Next
        cloneCurr = cloneCurr.Next
        curr.Next = curr.Next.Next
        curr = curr.Next
    }

    return dummy.Next
}

【 2. Course Schedule (위상 정렬) 】

func CanFinish(numCourses int, prerequisites [][]int) bool {
    // 그래프 구성
    graph := make([][]int, numCourses)
    inDegree := make([]int, numCourses)

    for _, p := range prerequisites {
        course, prereq := p[0], p[1]
        graph[prereq] = append(graph[prereq], course)
        inDegree[course]++
    }

    // inDegree가 0인 노드로 시작
    queue := []int{}
    for i := 0; i < numCourses; i++ {
        if inDegree[i] == 0 {
            queue = append(queue, i)
        }
    }

    count := 0
    for len(queue) > 0 {
        curr := queue[0]
        queue = queue[1:]
        count++

        for _, next := range graph[curr] {
            inDegree[next]--
            if inDegree[next] == 0 {
                queue = append(queue, next)
            }
        }
    }

    return count == numCourses
}
================================================================================
*/

// 테스트
func main() {
	// 그래프 생성
	n1 := &GraphNode{Val: 1}
	n2 := &GraphNode{Val: 2}
	n3 := &GraphNode{Val: 3}
	n4 := &GraphNode{Val: 4}
	n1.Neighbors = []*GraphNode{n2, n4}
	n2.Neighbors = []*GraphNode{n1, n3}
	n3.Neighbors = []*GraphNode{n2, n4}
	n4.Neighbors = []*GraphNode{n1, n3}

	// DFS로 복사
	cloneDFS := CloneGraph(n1)
	println("DFS Clone - Val:", cloneDFS.Val)
	println("DFS Clone - Is different:", cloneDFS != n1)
	println("DFS Clone - Neighbors count:", len(cloneDFS.Neighbors))

	// BFS로 복사
	cloneBFS := CloneGraphBFS(n1)
	println("\nBFS Clone - Val:", cloneBFS.Val)
	println("BFS Clone - Is different:", cloneBFS != n1)
}
