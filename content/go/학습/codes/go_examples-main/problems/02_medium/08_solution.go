package medium

/*
================================================================================
문제 8: Trie (Prefix Tree) - 솔루션
================================================================================
*/

/*
================================================================================
Trie 구조
================================================================================

【 트라이란? 】
- 문자열 검색에 특화된 트리 자료구조
- 각 노드는 하나의 문자를 나타냄
- 루트에서 리프까지의 경로 = 하나의 단어

┌─────────────────────────────────────────────────────────────────────────────┐
│                          Trie 시각화                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  단어: "apple", "app", "apt", "cat"                                        │
│                                                                             │
│                    (root)                                                   │
│                   /      \                                                  │
│                  a        c                                                 │
│                  |        |                                                 │
│                  p        a                                                 │
│                 / \       |                                                 │
│                p   t     t*                                                 │
│                |   *                                                        │
│                l                                                            │
│                |                                                            │
│                e*                                                           │
│                                                                             │
│  * = isEnd (단어의 끝)                                                      │
│                                                                             │
│  검색 "app": root → a → p → p* (true)                                      │
│  검색 "ap":  root → a → p (isEnd=false, false)                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
*/

// ============================================================================
// 방법 1: Map 기반 (유니코드 지원)
// ============================================================================

// TrieNode 트라이 노드
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
	count    int // 이 접두사를 가진 단어 수 (확장 기능용)
}

// Trie 트라이 자료구조
type Trie struct {
	root *TrieNode
}

// NewTrie 생성자
func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
		},
	}
}

// Insert 단어 삽입
/*
【 시간 복잡도 】O(m) - m은 단어 길이
【 공간 복잡도 】O(m) - 새 노드 생성
*/
func (t *Trie) Insert(word string) {
	node := t.root
	for _, ch := range word {
		if _, exists := node.children[ch]; !exists {
			node.children[ch] = &TrieNode{
				children: make(map[rune]*TrieNode),
			}
		}
		node = node.children[ch]
		node.count++ // 접두사 카운트 증가
	}
	node.isEnd = true
}

// Search 단어 검색
func (t *Trie) Search(word string) bool {
	node := t.findNode(word)
	return node != nil && node.isEnd
}

// StartsWith 접두사 검색
func (t *Trie) StartsWith(prefix string) bool {
	return t.findNode(prefix) != nil
}

// findNode 내부 헬퍼 - 접두사에 해당하는 노드 찾기
func (t *Trie) findNode(prefix string) *TrieNode {
	node := t.root
	for _, ch := range prefix {
		if _, exists := node.children[ch]; !exists {
			return nil
		}
		node = node.children[ch]
	}
	return node
}

// ============================================================================
// 확장 기능
// ============================================================================

// CountPrefix 접두사를 가진 단어 개수
func (t *Trie) CountPrefix(prefix string) int {
	node := t.findNode(prefix)
	if node == nil {
		return 0
	}
	return node.count
}

// GetWordsWithPrefix 접두사로 시작하는 모든 단어 반환
func (t *Trie) GetWordsWithPrefix(prefix string) []string {
	node := t.findNode(prefix)
	if node == nil {
		return []string{}
	}

	var results []string
	t.collectWords(node, prefix, &results)
	return results
}

// collectWords DFS로 단어 수집
func (t *Trie) collectWords(node *TrieNode, current string, results *[]string) {
	if node.isEnd {
		*results = append(*results, current)
	}
	for ch, child := range node.children {
		t.collectWords(child, current+string(ch), results)
	}
}

// Delete 단어 삭제
func (t *Trie) Delete(word string) bool {
	return t.deleteHelper(t.root, word, 0)
}

func (t *Trie) deleteHelper(node *TrieNode, word string, depth int) bool {
	if node == nil {
		return false
	}

	runes := []rune(word)

	if depth == len(runes) {
		// 단어 끝에 도달
		if !node.isEnd {
			return false // 단어가 존재하지 않음
		}
		node.isEnd = false
		return len(node.children) == 0 // 자식이 없으면 삭제 가능
	}

	ch := runes[depth]
	child, exists := node.children[ch]
	if !exists {
		return false
	}

	shouldDeleteChild := t.deleteHelper(child, word, depth+1)

	if shouldDeleteChild {
		delete(node.children, ch)
		child.count--
		return len(node.children) == 0 && !node.isEnd
	}

	child.count--
	return false
}

// ============================================================================
// 방법 2: 배열 기반 (소문자만, 더 빠름)
// ============================================================================

/*
【 배열 vs Map 비교 】
- 배열: O(1) 접근, 메모리 고정 (26 * 포인터 크기)
- Map: O(1) 평균, 메모리 동적, 해시 오버헤드

소문자 알파벳만 다루는 경우 배열이 더 효율적
*/

type TrieNodeArray struct {
	children [26]*TrieNodeArray
	isEnd    bool
}

type TrieArray struct {
	root *TrieNodeArray
}

func NewTrieArray() *TrieArray {
	return &TrieArray{root: &TrieNodeArray{}}
}

func (t *TrieArray) Insert(word string) {
	node := t.root
	for _, ch := range word {
		idx := ch - 'a'
		if node.children[idx] == nil {
			node.children[idx] = &TrieNodeArray{}
		}
		node = node.children[idx]
	}
	node.isEnd = true
}

func (t *TrieArray) Search(word string) bool {
	node := t.root
	for _, ch := range word {
		idx := ch - 'a'
		if node.children[idx] == nil {
			return false
		}
		node = node.children[idx]
	}
	return node.isEnd
}

/*
================================================================================
Trie 응용
================================================================================

【 1. 자동완성 】

type Autocomplete struct {
    trie *Trie
}

func (a *Autocomplete) Suggest(prefix string, limit int) []string {
    words := a.trie.GetWordsWithPrefix(prefix)
    if len(words) > limit {
        return words[:limit]
    }
    return words
}

【 2. 와일드카드 검색 (.은 아무 문자) 】

func (t *Trie) SearchWithWildcard(pattern string) bool {
    return t.searchWildcard(t.root, []rune(pattern), 0)
}

func (t *Trie) searchWildcard(node *TrieNode, pattern []rune, idx int) bool {
    if node == nil {
        return false
    }
    if idx == len(pattern) {
        return node.isEnd
    }

    ch := pattern[idx]
    if ch == '.' {
        // 모든 자식 탐색
        for _, child := range node.children {
            if t.searchWildcard(child, pattern, idx+1) {
                return true
            }
        }
        return false
    }

    return t.searchWildcard(node.children[ch], pattern, idx+1)
}

【 3. 단어 교체/수정 거리 】

Trie + 편집 거리(Levenshtein Distance) 조합으로
맞춤법 검사 및 유사 단어 추천 구현 가능

================================================================================
시간/공간 복잡도
================================================================================

【 시간 복잡도 】
- Insert: O(m)
- Search: O(m)
- StartsWith: O(m)
- Delete: O(m)
여기서 m = 단어/접두사 길이

【 공간 복잡도 】
- 최악: O(N * M * C)
  - N: 단어 개수
  - M: 평균 단어 길이
  - C: 문자 집합 크기 (알파벳=26)
- 최선: 공통 접두사가 많으면 공간 절약
================================================================================
*/

// 테스트
func main() {
	trie := NewTrie()

	// 기본 테스트
	trie.Insert("apple")
	trie.Insert("app")
	trie.Insert("application")
	trie.Insert("apt")
	trie.Insert("cat")

	println("Search 'apple':", trie.Search("apple"))       // true
	println("Search 'app':", trie.Search("app"))           // true
	println("Search 'ap':", trie.Search("ap"))             // false
	println("StartsWith 'app':", trie.StartsWith("app"))   // true
	println("StartsWith 'cat':", trie.StartsWith("cat"))   // true
	println("StartsWith 'dog':", trie.StartsWith("dog"))   // false

	// 확장 기능 테스트
	println("\nCountPrefix 'app':", trie.CountPrefix("app")) // 3

	println("\nWords with prefix 'app':")
	for _, word := range trie.GetWordsWithPrefix("app") {
		println(" -", word)
	}

	// 삭제 테스트
	println("\nDelete 'app':", trie.Delete("app"))
	println("Search 'app' after delete:", trie.Search("app"))         // false
	println("Search 'apple' after delete:", trie.Search("apple"))     // true
}
