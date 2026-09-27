package medium

/*
================================================================================
문제 8: Trie (Prefix Tree) 구현
================================================================================

난이도: Medium
주제: 트라이, 문자열, 자료구조

【 문제 설명 】
Trie(트라이) 자료구조를 구현하세요.

Trie는 다음 연산을 지원합니다:
- Insert(word): 단어 삽입
- Search(word): 단어가 존재하는지 확인
- StartsWith(prefix): 해당 접두사로 시작하는 단어가 있는지 확인

【 인터페이스 】
type Trie interface {
    Insert(word string)
    Search(word string) bool
    StartsWith(prefix string) bool
}

【 예시 】
trie := NewTrie()
trie.Insert("apple")
trie.Search("apple")   // true
trie.Search("app")     // false
trie.StartsWith("app") // true
trie.Insert("app")
trie.Search("app")     // true

【 제약 조건 】
- 1 <= word.length, prefix.length <= 2000
- word와 prefix는 소문자 영문자로만 구성

【 힌트 】
1. 노드 구조: children map[rune]*TrieNode, isEnd bool
2. 또는 children [26]*TrieNode (소문자만)

【 확장 과제 】
1. GetWordsWithPrefix(prefix): 접두사로 시작하는 모든 단어 반환
2. Delete(word): 단어 삭제
3. CountPrefix(prefix): 접두사를 가진 단어 개수

【 실무 연관성 】
- 자동완성
- 맞춤법 검사
- IP 라우팅 테이블
- 검색 엔진 인덱싱
================================================================================
*/

// Trie 구조체를 정의하고 구현하세요
type Trie struct {
	// 여기에 필드를 정의하세요
}

// NewTrie 생성자
func NewTrie() *Trie {
	// 여기에 코드를 작성하세요
	return nil
}

// Insert 단어 삽입
func (t *Trie) Insert(word string) {
	// 여기에 코드를 작성하세요
}

// Search 단어 검색
func (t *Trie) Search(word string) bool {
	// 여기에 코드를 작성하세요
	return false
}

// StartsWith 접두사 검색
func (t *Trie) StartsWith(prefix string) bool {
	// 여기에 코드를 작성하세요
	return false
}

// 테스트
func main() {
	trie := NewTrie()
	trie.Insert("apple")
	println("Search apple:", trie.Search("apple"))     // true
	println("Search app:", trie.Search("app"))         // false
	println("StartsWith app:", trie.StartsWith("app")) // true
	trie.Insert("app")
	println("Search app:", trie.Search("app")) // true
}
