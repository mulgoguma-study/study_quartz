package medium

/*
================================================================================
문제 10: Design Twitter - 솔루션
================================================================================
*/

import (
	"container/heap"
)

/*
================================================================================
시스템 설계 분석
================================================================================

【 핵심 설계 결정 】

1. 트윗 저장
   - 사용자별 트윗 리스트
   - 타임스탬프 필요 (최신순 정렬)

2. 팔로우 관계
   - 팔로워 → 팔로이 Set
   - O(1) 추가/삭제/확인

3. 뉴스피드
   - 자신 + 팔로이들의 트윗
   - 최신 10개 → k-way merge (Max-Heap)

┌─────────────────────────────────────────────────────────────────────────────┐
│                        데이터 구조                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  tweets: map[userId] → [{tweetId, timestamp}, ...]                         │
│                                                                             │
│  follows: map[userId] → set{followeeId, ...}                               │
│                                                                             │
│  GetNewsFeed 알고리즘:                                                      │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │                                                                       │ │
│  │  User1 tweets: [t1, t3, t5]  ─┐                                       │ │
│  │  User2 tweets: [t2, t4]      ─┼─→ Max-Heap ─→ Top 10 최신 트윗        │ │
│  │  User3 tweets: [t6, t7, t8]  ─┘                                       │ │
│  │                                                                       │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
*/

// Tweet 트윗 구조체
type Tweet struct {
	id        int
	timestamp int
}

// Twitter 트위터 구현
type Twitter struct {
	timestamp int                    // 전역 타임스탬프
	tweets    map[int][]Tweet        // userId → 트윗 리스트
	follows   map[int]map[int]bool   // followerId → set of followeeIds
}

// NewTwitter 생성자
func NewTwitter() *Twitter {
	return &Twitter{
		timestamp: 0,
		tweets:    make(map[int][]Tweet),
		follows:   make(map[int]map[int]bool),
	}
}

// PostTweet 트윗 작성
func (t *Twitter) PostTweet(userId int, tweetId int) {
	// 트윗 추가 (최신이 뒤에)
	t.tweets[userId] = append(t.tweets[userId], Tweet{
		id:        tweetId,
		timestamp: t.timestamp,
	})
	t.timestamp++
}

// Follow 팔로우
func (t *Twitter) Follow(followerId int, followeeId int) {
	if followerId == followeeId {
		return // 자기 자신 팔로우 방지
	}
	if t.follows[followerId] == nil {
		t.follows[followerId] = make(map[int]bool)
	}
	t.follows[followerId][followeeId] = true
}

// Unfollow 언팔로우
func (t *Twitter) Unfollow(followerId int, followeeId int) {
	if t.follows[followerId] != nil {
		delete(t.follows[followerId], followeeId)
	}
}

// ============================================================================
// GetNewsFeed - k-way Merge with Max-Heap
// ============================================================================

// feedItem 힙에 저장할 아이템
type feedItem struct {
	tweetId   int
	timestamp int
	userId    int
	index     int // 해당 사용자의 트윗 리스트에서의 인덱스
}

// maxHeap Max-Heap 구현
type maxHeap []feedItem

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i].timestamp > h[j].timestamp } // Max-Heap
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(feedItem)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// GetNewsFeed 뉴스피드 조회
/*
【 알고리즘: k-way Merge 】
1. 자신 + 팔로이들의 최신 트윗을 힙에 추가
2. 힙에서 최신 트윗 추출
3. 해당 사용자의 이전 트윗을 힙에 추가
4. 10개가 될 때까지 반복

【 시간 복잡도 】O(k log k + 10 log k) = O(k log k)
  k = 팔로이 수 + 1
【 공간 복잡도 】O(k)
*/
func (t *Twitter) GetNewsFeed(userId int) []int {
	result := make([]int, 0, 10)
	h := &maxHeap{}
	heap.Init(h)

	// 자신과 팔로이들의 사용자 ID 수집
	users := []int{userId}
	for followeeId := range t.follows[userId] {
		users = append(users, followeeId)
	}

	// 각 사용자의 최신 트윗을 힙에 추가
	for _, uid := range users {
		tweets := t.tweets[uid]
		if len(tweets) > 0 {
			lastIdx := len(tweets) - 1
			heap.Push(h, feedItem{
				tweetId:   tweets[lastIdx].id,
				timestamp: tweets[lastIdx].timestamp,
				userId:    uid,
				index:     lastIdx,
			})
		}
	}

	// 최신 10개 추출
	for h.Len() > 0 && len(result) < 10 {
		item := heap.Pop(h).(feedItem)
		result = append(result, item.tweetId)

		// 해당 사용자의 이전 트윗 추가
		if item.index > 0 {
			prevIdx := item.index - 1
			tweets := t.tweets[item.userId]
			heap.Push(h, feedItem{
				tweetId:   tweets[prevIdx].id,
				timestamp: tweets[prevIdx].timestamp,
				userId:    item.userId,
				index:     prevIdx,
			})
		}
	}

	return result
}

/*
================================================================================
실제 Twitter 시스템과의 비교
================================================================================

【 실제 시스템의 복잡성 】

1. Fan-out on Write vs Fan-out on Read
   ───────────────────────────────────
   - Fan-out on Write: 트윗 작성 시 모든 팔로워의 타임라인에 저장
     장점: 읽기 빠름
     단점: 쓰기 부하 (팔로워가 많으면 심각)

   - Fan-out on Read: 타임라인 요청 시 팔로이들의 트윗 조합
     장점: 쓰기 간단
     단점: 읽기 시 계산 필요

   - 실제: 하이브리드 (일반 유저는 Write, 셀럽은 Read)

2. 캐싱 전략
   ──────────
   - 타임라인 캐싱 (Redis Sorted Set)
   - 핫 사용자 캐싱
   - CDN을 통한 미디어 캐싱

3. 샤딩
   ─────
   - 사용자 ID 기반 샤딩
   - 트윗 ID 기반 샤딩

4. 비동기 처리
   ───────────
   - 메시지 큐로 팔로워 피드 업데이트
   - 비동기 알림 처리

================================================================================
*/

// 테스트
func main() {
	twitter := NewTwitter()

	// 테스트 시나리오
	twitter.PostTweet(1, 5)
	println("User 1 posts tweet 5")

	feed := twitter.GetNewsFeed(1)
	print("User 1's feed: ")
	for _, id := range feed {
		print(id, " ")
	}
	println()

	twitter.Follow(1, 2)
	println("User 1 follows User 2")

	twitter.PostTweet(2, 6)
	println("User 2 posts tweet 6")

	feed = twitter.GetNewsFeed(1)
	print("User 1's feed after follow: ")
	for _, id := range feed {
		print(id, " ")
	}
	println()

	twitter.Unfollow(1, 2)
	println("User 1 unfollows User 2")

	feed = twitter.GetNewsFeed(1)
	print("User 1's feed after unfollow: ")
	for _, id := range feed {
		print(id, " ")
	}
	println()

	// 여러 트윗 테스트
	println("\n=== Multiple tweets test ===")
	for i := 7; i <= 15; i++ {
		twitter.PostTweet(1, i)
	}
	feed = twitter.GetNewsFeed(1)
	print("User 1's feed (should be 10 tweets): ")
	for _, id := range feed {
		print(id, " ")
	}
	println()
}
