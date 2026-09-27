package hard

/*
================================================================================
문제 4: Word Ladder (단어 사다리)
================================================================================

난이도: Hard
주제: BFS, 그래프

【 문제 설명 】
시작 단어에서 끝 단어까지 변환하는 최단 경로를 찾으세요.
한 번에 한 글자만 바꿀 수 있고, 중간 단어들은 모두 사전에 있어야 합니다.

【 예시 】
beginWord = "hit"
endWord = "cog"
wordList = ["hot","dot","dog","lot","log","cog"]

반환: 5  ("hit" → "hot" → "dot" → "dog" → "cog")

【 제약 조건 】
- 모든 단어 길이 동일
- 소문자만 사용
- endWord가 wordList에 없으면 0 반환

【 힌트 】
1. BFS로 최단 경로 탐색
2. 각 위치에 a-z 대입해서 다음 단어 후보 생성
================================================================================
*/

// LadderLength 최단 변환 경로 길이
func LadderLength(beginWord string, endWord string, wordList []string) int {
	// 여기에 코드를 작성하세요
	return 0
}

// FindLadders 모든 최단 경로 찾기 (Hard+)
func FindLadders(beginWord string, endWord string, wordList []string) [][]string {
	// 여기에 코드를 작성하세요
	return nil
}

func main() {
	wordList := []string{"hot", "dot", "dog", "lot", "log", "cog"}
	result := LadderLength("hit", "cog", wordList)
	println("Shortest path length:", result) // 5
}
