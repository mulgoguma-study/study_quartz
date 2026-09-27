package medium

/*
================================================================================
문제 5: Sliding Window Maximum - 솔루션
================================================================================
*/

/*
================================================================================
Monotonic Deque (단조 덱) 개념
================================================================================

【 핵심 아이디어 】
- 덱을 단조 감소 상태로 유지
- 새 원소가 들어올 때, 더 작은 원소들은 제거 (어차피 최대값이 될 수 없음)
- 덱의 앞(front)은 항상 현재 윈도우의 최대값

┌─────────────────────────────────────────────────────────────────────────────┐
│                      Monotonic Deque 동작 예시                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  nums = [1, 3, -1, -3, 5, 3, 6, 7], k = 3                                  │
│                                                                             │
│  i=0: nums[0]=1                                                             │
│       deque = [0]           (인덱스 저장)                                   │
│                                                                             │
│  i=1: nums[1]=3 > nums[0]=1                                                 │
│       1은 절대 최대값이 될 수 없음 → 제거                                   │
│       deque = [1]                                                           │
│                                                                             │
│  i=2: nums[2]=-1 < nums[1]=3                                                │
│       -1은 나중에 3이 빠지면 최대값이 될 수 있음 → 유지                     │
│       deque = [1, 2]                                                        │
│       윈도우 완성! max = nums[1] = 3                                        │
│                                                                             │
│  i=3: nums[3]=-3                                                            │
│       deque = [1, 2, 3]                                                     │
│       인덱스 1은 아직 윈도우 안 (i-k+1=1)                                   │
│       max = nums[1] = 3                                                     │
│                                                                             │
│  i=4: nums[4]=5 > 모든 원소                                                 │
│       모두 제거                                                             │
│       deque = [4]                                                           │
│       max = nums[4] = 5                                                     │
│                                                                             │
│  ... 계속                                                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
*/

// ============================================================================
// 방법 1: Monotonic Deque - O(n)
// ============================================================================

/*
【 왜 O(n)인가? 】
- 각 원소는 덱에 최대 1번 들어가고 1번 나감
- 전체 push + pop 연산 = 2n
- 따라서 O(n)
*/
func MaxSlidingWindow(nums []int, k int) []int {
	if len(nums) == 0 || k == 0 {
		return []int{}
	}

	n := len(nums)
	result := make([]int, 0, n-k+1)

	// 덱: 인덱스를 저장 (값이 아님!)
	// 단조 감소 순서 유지 (front가 가장 큼)
	deque := make([]int, 0, k)

	for i := 0; i < n; i++ {
		// 1. 윈도우 범위를 벗어난 인덱스 제거 (front에서)
		// i - k + 1 이 현재 윈도우의 시작 인덱스
		if len(deque) > 0 && deque[0] < i-k+1 {
			deque = deque[1:] // front 제거
		}

		// 2. 현재 값보다 작은 값들 제거 (back에서)
		// 이 값들은 절대 최대값이 될 수 없음
		for len(deque) > 0 && nums[deque[len(deque)-1]] < nums[i] {
			deque = deque[:len(deque)-1] // back 제거
		}

		// 3. 현재 인덱스 추가
		deque = append(deque, i)

		// 4. 윈도우가 완성되면 (i >= k-1) 최대값 기록
		if i >= k-1 {
			result = append(result, nums[deque[0]])
		}
	}

	return result
}

// ============================================================================
// 방법 2: 브루트 포스 - O(n*k) (참고용)
// ============================================================================

func MaxSlidingWindowBruteForce(nums []int, k int) []int {
	if len(nums) == 0 || k == 0 {
		return []int{}
	}

	n := len(nums)
	result := make([]int, 0, n-k+1)

	for i := 0; i <= n-k; i++ {
		// 각 윈도우에서 최대값 찾기
		maxVal := nums[i]
		for j := i + 1; j < i+k; j++ {
			if nums[j] > maxVal {
				maxVal = nums[j]
			}
		}
		result = append(result, maxVal)
	}

	return result
}

// ============================================================================
// 방법 3: Heap 사용 - O(n log n) (참고용)
// ============================================================================

/*
Heap을 사용하면:
- 최대값 조회: O(1)
- 삽입: O(log n)
- 삭제: O(log n)

하지만 윈도우에서 벗어난 원소 제거가 어려움
→ Lazy deletion 필요 (최대값 조회 시 검증)
*/

/*
================================================================================
Monotonic Stack/Deque 패턴 활용
================================================================================

【 언제 사용하는가? 】

1. Next Greater Element (다음 큰 원소)
   - 각 원소의 오른쪽에서 첫 번째로 큰 원소 찾기
   - Monotonic Stack (증가)

2. Previous Greater Element (이전 큰 원소)
   - 각 원소의 왼쪽에서 첫 번째로 큰 원소 찾기
   - Monotonic Stack (증가)

3. Sliding Window Maximum/Minimum
   - 윈도우 내 최대/최소값 추적
   - Monotonic Deque

4. Largest Rectangle in Histogram
   - 히스토그램에서 가장 큰 직사각형
   - Monotonic Stack

5. Trapping Rain Water
   - 빗물 가두기
   - Monotonic Stack 또는 Two Pointer

【 코드 패턴 】

// Monotonic Stack (증가)
stack := []int{}
for i, v := range nums {
    for len(stack) > 0 && nums[stack[len(stack)-1]] < v {
        // 스택 top보다 현재 값이 크면 pop
        top := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        // top의 "next greater" = 현재 값
    }
    stack = append(stack, i)
}

// Monotonic Stack (감소)
for len(stack) > 0 && nums[stack[len(stack)-1]] > v {
    // 스택 top보다 현재 값이 작으면 pop
}

// Monotonic Deque
for len(deque) > 0 && deque[0] < windowStart {
    deque = deque[1:]  // front 제거 (범위 벗어남)
}
for len(deque) > 0 && nums[deque[len(deque)-1]] < v {
    deque = deque[:len(deque)-1]  // back 제거
}
deque = append(deque, i)

================================================================================
*/

// 테스트
func main() {
	testCases := []struct {
		nums   []int
		k      int
		expect []int
	}{
		{[]int{1, 3, -1, -3, 5, 3, 6, 7}, 3, []int{3, 3, 5, 5, 6, 7}},
		{[]int{1}, 1, []int{1}},
		{[]int{1, -1}, 1, []int{1, -1}},
		{[]int{9, 11}, 2, []int{11}},
		{[]int{4, -2}, 2, []int{4}},
	}

	for i, tc := range testCases {
		result := MaxSlidingWindow(tc.nums, tc.k)

		print("Test ", i+1, ": ")
		for _, v := range result {
			print(v, " ")
		}
		print(" | Expected: ")
		for _, v := range tc.expect {
			print(v, " ")
		}
		println()
	}
}
