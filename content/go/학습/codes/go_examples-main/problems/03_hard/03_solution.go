package hard

/*
================================================================================
문제 3: Serialize and Deserialize Binary Tree - 솔루션
================================================================================
*/

import (
	"strconv"
	"strings"
)

/*
================================================================================
직렬화 전략 비교
================================================================================

【 방법 1: Pre-order DFS 】

장점:
- 재귀로 구현 간단
- 스트림 처리 가능

직렬화 순서: 루트 → 왼쪽 서브트리 → 오른쪽 서브트리

      1
     / \
    2   3
       / \
      4   5

→ "1,2,#,#,3,4,#,#,5,#,#"  (# = null)

【 방법 2: Level-order BFS 】

장점:
- 직관적 (레벨별로 읽음)
- 완전 이진 트리에 적합

→ "1,2,3,#,#,4,5,#,#,#,#"

================================================================================
*/

const NULL = "#"

// TreeNode 이진 트리 노드
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ============================================================================
// 방법 1: Pre-order DFS (권장)
// ============================================================================

// Codec 직렬화/역직렬화 코덱
type Codec struct{}

// Serialize Pre-order로 직렬화
/*
【 동작 】
1. 현재 노드 값 출력
2. 왼쪽 서브트리 재귀
3. 오른쪽 서브트리 재귀
4. null은 "#"으로 표현
*/
func (c *Codec) Serialize(root *TreeNode) string {
	var builder strings.Builder
	c.serializeHelper(root, &builder)
	result := builder.String()
	// 마지막 쉼표 제거
	if len(result) > 0 {
		result = result[:len(result)-1]
	}
	return result
}

func (c *Codec) serializeHelper(node *TreeNode, builder *strings.Builder) {
	if node == nil {
		builder.WriteString(NULL)
		builder.WriteString(",")
		return
	}

	builder.WriteString(strconv.Itoa(node.Val))
	builder.WriteString(",")
	c.serializeHelper(node.Left, builder)
	c.serializeHelper(node.Right, builder)
}

// Deserialize 역직렬화
/*
【 동작 】
1. 문자열을 토큰으로 분리
2. 인덱스 포인터로 순차 처리
3. Pre-order 순서로 트리 재구성
*/
func (c *Codec) Deserialize(data string) *TreeNode {
	if data == "" {
		return nil
	}

	tokens := strings.Split(data, ",")
	idx := 0
	return c.deserializeHelper(tokens, &idx)
}

func (c *Codec) deserializeHelper(tokens []string, idx *int) *TreeNode {
	if *idx >= len(tokens) || tokens[*idx] == NULL {
		*idx++
		return nil
	}

	val, _ := strconv.Atoi(tokens[*idx])
	*idx++

	node := &TreeNode{Val: val}
	node.Left = c.deserializeHelper(tokens, idx)
	node.Right = c.deserializeHelper(tokens, idx)

	return node
}

// ============================================================================
// 방법 2: Level-order BFS
// ============================================================================

type CodecBFS struct{}

// Serialize BFS로 직렬화
func (c *CodecBFS) Serialize(root *TreeNode) string {
	if root == nil {
		return ""
	}

	var result []string
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node == nil {
			result = append(result, NULL)
		} else {
			result = append(result, strconv.Itoa(node.Val))
			queue = append(queue, node.Left)
			queue = append(queue, node.Right)
		}
	}

	// 뒤쪽 null 제거 (압축)
	for len(result) > 0 && result[len(result)-1] == NULL {
		result = result[:len(result)-1]
	}

	return strings.Join(result, ",")
}

// Deserialize BFS로 역직렬화
func (c *CodecBFS) Deserialize(data string) *TreeNode {
	if data == "" {
		return nil
	}

	tokens := strings.Split(data, ",")
	if len(tokens) == 0 || tokens[0] == NULL {
		return nil
	}

	rootVal, _ := strconv.Atoi(tokens[0])
	root := &TreeNode{Val: rootVal}
	queue := []*TreeNode{root}
	i := 1

	for len(queue) > 0 && i < len(tokens) {
		node := queue[0]
		queue = queue[1:]

		// 왼쪽 자식
		if i < len(tokens) && tokens[i] != NULL {
			val, _ := strconv.Atoi(tokens[i])
			node.Left = &TreeNode{Val: val}
			queue = append(queue, node.Left)
		}
		i++

		// 오른쪽 자식
		if i < len(tokens) && tokens[i] != NULL {
			val, _ := strconv.Atoi(tokens[i])
			node.Right = &TreeNode{Val: val}
			queue = append(queue, node.Right)
		}
		i++
	}

	return root
}

/*
================================================================================
압축 최적화
================================================================================

【 비트 패킹 】

각 노드: 값(4바이트) + 존재 플래그(2비트)

【 차이 인코딩 】

연속된 값의 차이만 저장:
[1, 2, 3, 5] → [1, 1, 1, 2]

【 가변 길이 인코딩 】

작은 숫자는 적은 바이트로:
- 0-127: 1바이트
- 128-16383: 2바이트
- ...

================================================================================
*/

// printTree 트리 출력 (Pre-order)
func printTree(node *TreeNode, prefix string) {
	if node == nil {
		println(prefix + "nil")
		return
	}
	println(prefix + strconv.Itoa(node.Val))
	printTree(node.Left, prefix+"  L:")
	printTree(node.Right, prefix+"  R:")
}

// 테스트
func main() {
	// 트리 생성
	root := &TreeNode{Val: 1}
	root.Left = &TreeNode{Val: 2}
	root.Right = &TreeNode{Val: 3}
	root.Right.Left = &TreeNode{Val: 4}
	root.Right.Right = &TreeNode{Val: 5}

	println("=== Original Tree ===")
	printTree(root, "")

	// DFS 방식
	println("\n=== DFS Codec ===")
	codecDFS := &Codec{}
	serializedDFS := codecDFS.Serialize(root)
	println("Serialized:", serializedDFS)

	deserializedDFS := codecDFS.Deserialize(serializedDFS)
	println("\nDeserialized:")
	printTree(deserializedDFS, "")

	// BFS 방식
	println("\n=== BFS Codec ===")
	codecBFS := &CodecBFS{}
	serializedBFS := codecBFS.Serialize(root)
	println("Serialized:", serializedBFS)

	deserializedBFS := codecBFS.Deserialize(serializedBFS)
	println("\nDeserialized:")
	printTree(deserializedBFS, "")
}
