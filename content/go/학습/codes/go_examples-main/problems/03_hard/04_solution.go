package hard

/*
================================================================================
문제 4: Word Ladder - 솔루션
================================================================================
*/

/*
【 BFS 접근 】

그래프 모델링:
- 노드: 단어
- 간선: 한 글자만 다른 단어 쌍

BFS로 시작 단어에서 끝 단어까지 최단 경로 탐색
*/

// LadderLength BFS로 최단 경로 찾기
func LadderLength(beginWord string, endWord string, wordList []string) int {
	// 1. wordList를 Set으로 변환 (빠른 조회)
	wordSet := make(map[string]bool)
	for _, word := range wordList {
		wordSet[word] = true
	}

	// endWord가 없으면 불가능
	if !wordSet[endWord] {
		return 0
	}

	// 2. BFS
	queue := []string{beginWord}
	visited := make(map[string]bool)
	visited[beginWord] = true
	level := 1

	for len(queue) > 0 {
		levelSize := len(queue)
		level++

		for i := 0; i < levelSize; i++ {
			word := queue[0]
			queue = queue[1:]

			// 각 위치에 a-z 대입
			wordBytes := []byte(word)
			for j := 0; j < len(wordBytes); j++ {
				original := wordBytes[j]

				for c := byte('a'); c <= 'z'; c++ {
					if c == original {
						continue
					}

					wordBytes[j] = c
					newWord := string(wordBytes)

					if newWord == endWord {
						return level
					}

					if wordSet[newWord] && !visited[newWord] {
						visited[newWord] = true
						queue = append(queue, newWord)
					}
				}

				wordBytes[j] = original
			}
		}
	}

	return 0
}

// 양방향 BFS (최적화)
func LadderLengthBidirectional(beginWord string, endWord string, wordList []string) int {
	wordSet := make(map[string]bool)
	for _, word := range wordList {
		wordSet[word] = true
	}

	if !wordSet[endWord] {
		return 0
	}

	// 양쪽에서 탐색
	beginSet := map[string]bool{beginWord: true}
	endSet := map[string]bool{endWord: true}
	visited := make(map[string]bool)
	level := 1

	for len(beginSet) > 0 && len(endSet) > 0 {
		// 작은 쪽에서 확장 (최적화)
		if len(beginSet) > len(endSet) {
			beginSet, endSet = endSet, beginSet
		}

		nextSet := make(map[string]bool)
		level++

		for word := range beginSet {
			wordBytes := []byte(word)
			for j := 0; j < len(wordBytes); j++ {
				original := wordBytes[j]

				for c := byte('a'); c <= 'z'; c++ {
					if c == original {
						continue
					}

					wordBytes[j] = c
					newWord := string(wordBytes)

					// 반대편에서 만남!
					if endSet[newWord] {
						return level
					}

					if wordSet[newWord] && !visited[newWord] {
						visited[newWord] = true
						nextSet[newWord] = true
					}
				}

				wordBytes[j] = original
			}
		}

		beginSet = nextSet
	}

	return 0
}

// FindLadders 모든 최단 경로 찾기
func FindLadders(beginWord string, endWord string, wordList []string) [][]string {
	wordSet := make(map[string]bool)
	for _, word := range wordList {
		wordSet[word] = true
	}

	if !wordSet[endWord] {
		return nil
	}

	// BFS로 각 단어까지의 거리 계산
	distance := make(map[string]int)
	parent := make(map[string][]string)

	queue := []string{beginWord}
	distance[beginWord] = 0
	found := false

	for len(queue) > 0 && !found {
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			word := queue[0]
			queue = queue[1:]
			currDist := distance[word]

			wordBytes := []byte(word)
			for j := 0; j < len(wordBytes); j++ {
				original := wordBytes[j]

				for c := byte('a'); c <= 'z'; c++ {
					if c == original {
						continue
					}

					wordBytes[j] = c
					newWord := string(wordBytes)

					if !wordSet[newWord] {
						continue
					}

					if _, exists := distance[newWord]; !exists {
						distance[newWord] = currDist + 1
						queue = append(queue, newWord)
					}

					if distance[newWord] == currDist+1 {
						parent[newWord] = append(parent[newWord], word)
					}

					if newWord == endWord {
						found = true
					}
				}

				wordBytes[j] = original
			}
		}
	}

	// DFS로 경로 복원
	var result [][]string
	if found {
		var path []string
		var dfs func(word string)
		dfs = func(word string) {
			path = append([]string{word}, path...)
			if word == beginWord {
				pathCopy := make([]string, len(path))
				copy(pathCopy, path)
				result = append(result, pathCopy)
			} else {
				for _, p := range parent[word] {
					dfs(p)
				}
			}
			path = path[1:]
		}
		dfs(endWord)
	}

	return result
}

func main() {
	wordList := []string{"hot", "dot", "dog", "lot", "log", "cog"}

	// 최단 경로 길이
	length := LadderLength("hit", "cog", wordList)
	println("Shortest path length:", length)

	// 양방향 BFS
	lengthBi := LadderLengthBidirectional("hit", "cog", wordList)
	println("Bidirectional BFS length:", lengthBi)

	// 모든 최단 경로
	paths := FindLadders("hit", "cog", wordList)
	println("\nAll shortest paths:")
	for _, path := range paths {
		for _, word := range path {
			print(word, " ")
		}
		println()
	}
}
