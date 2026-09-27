package medium

/*
================================================================================
문제 10: Design Twitter (시스템 설계)
================================================================================

난이도: Medium
주제: 시스템 설계, 해시맵, 힙

【 문제 설명 】
간단한 Twitter를 설계하세요.

지원 기능:
- PostTweet(userId, tweetId): 트윗 작성
- GetNewsFeed(userId): 팔로우한 사용자 + 자신의 최신 트윗 10개
- Follow(followerId, followeeId): 팔로우
- Unfollow(followerId, followeeId): 언팔로우

【 예시 】
twitter := NewTwitter()
twitter.PostTweet(1, 5)      // User 1이 tweet 5 작성
twitter.GetNewsFeed(1)       // [5]
twitter.Follow(1, 2)         // User 1이 User 2 팔로우
twitter.PostTweet(2, 6)      // User 2가 tweet 6 작성
twitter.GetNewsFeed(1)       // [6, 5] (최신순)
twitter.Unfollow(1, 2)
twitter.GetNewsFeed(1)       // [5]

【 제약 조건 】
- userId, tweetId는 양의 정수
- GetNewsFeed는 최신 10개까지만 반환
- 동일한 follow 중복 호출 처리

【 힌트 】
1. 사용자별 트윗 리스트 저장
2. 팔로우 관계는 Set으로 관리
3. GetNewsFeed는 k-way merge (힙)

【 실무 연관성 】
- 소셜 미디어 피드
- 타임라인 시스템
- 알림 시스템
================================================================================
*/

// Twitter 구조체를 정의하고 구현하세요
type Twitter struct {
	// 여기에 필드를 정의하세요
}

// NewTwitter 생성자
func NewTwitter() *Twitter {
	// 여기에 코드를 작성하세요
	return nil
}

// PostTweet 트윗 작성
func (t *Twitter) PostTweet(userId int, tweetId int) {
	// 여기에 코드를 작성하세요
}

// GetNewsFeed 뉴스피드 조회 (최신 10개)
func (t *Twitter) GetNewsFeed(userId int) []int {
	// 여기에 코드를 작성하세요
	return nil
}

// Follow 팔로우
func (t *Twitter) Follow(followerId int, followeeId int) {
	// 여기에 코드를 작성하세요
}

// Unfollow 언팔로우
func (t *Twitter) Unfollow(followerId int, followeeId int) {
	// 여기에 코드를 작성하세요
}

// 테스트
func main() {
	twitter := NewTwitter()
	twitter.PostTweet(1, 5)
	feed := twitter.GetNewsFeed(1)
	println("Feed after post:", feed[0]) // [5]

	twitter.Follow(1, 2)
	twitter.PostTweet(2, 6)
	feed = twitter.GetNewsFeed(1)
	println("Feed after follow:", feed[0], feed[1]) // [6, 5]

	twitter.Unfollow(1, 2)
	feed = twitter.GetNewsFeed(1)
	println("Feed after unfollow:", feed[0]) // [5]
}
