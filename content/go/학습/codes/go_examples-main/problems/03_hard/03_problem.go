package hard

/*
================================================================================
문제 3: Serialize and Deserialize Binary Tree
================================================================================

난이도: Hard
주제: 이진 트리, DFS/BFS, 직렬화

【 문제 설명 】
이진 트리를 문자열로 직렬화하고, 다시 원래 트리로 역직렬화하세요.

【 인터페이스 】
type Codec struct{}
func (c *Codec) Serialize(root *TreeNode) string
func (c *Codec) Deserialize(data string) *TreeNode

【 예시 】
      1
     / \
    2   3
       / \
      4   5

Serialize → "1,2,null,null,3,4,null,null,5,null,null"
Deserialize → 원래 트리

【 제약 조건 】
- 노드 수: 0 ~ 10^4
- -1000 <= Node.val <= 1000
- 직렬화 포맷은 자유롭게 설계

【 힌트 】
1. Pre-order traversal 사용
2. null 값도 명시적으로 표현
3. 구분자 사용 (쉼표 등)

【 도전 과제 】
- 여러 직렬화 방식 구현 (Pre-order, Level-order)
- 더 컴팩트한 포맷 설계
================================================================================
*/

import "strings"

// TreeNode 이진 트리 노드
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// Codec 직렬화/역직렬화 코덱
type Codec struct {
	// 필요한 필드 추가
}

// Serialize 트리를 문자열로 변환
func (c *Codec) Serialize(root *TreeNode) string {
	// 여기에 코드를 작성하세요
	return ""
}

// Deserialize 문자열을 트리로 변환
func (c *Codec) Deserialize(data string) *TreeNode {
	// 여기에 코드를 작성하세요
	return nil
}

// 테스트
func main() {
	// 트리 생성:     1
	//              / \
	//             2   3
	//                / \
	//               4   5
	root := &TreeNode{Val: 1}
	root.Left = &TreeNode{Val: 2}
	root.Right = &TreeNode{Val: 3}
	root.Right.Left = &TreeNode{Val: 4}
	root.Right.Right = &TreeNode{Val: 5}

	codec := &Codec{}
	serialized := codec.Serialize(root)
	println("Serialized:", serialized)

	deserialized := codec.Deserialize(serialized)
	println("Deserialized root val:", deserialized.Val)
}

var _ = strings.Split // import 사용을 위한 더미
