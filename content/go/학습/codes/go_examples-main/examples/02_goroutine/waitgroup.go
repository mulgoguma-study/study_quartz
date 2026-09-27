// =============================================================================
// sync.WaitGroup 예제
// =============================================================================
// WaitGroup은 고루틴의 완료를 기다리는 동기화 도구입니다.
//
// 주요 메서드:
// - Add(n): 대기할 고루틴 수를 n만큼 증가
// - Done(): 고루틴 완료 시 호출 (내부적으로 Add(-1))
// - Wait(): 모든 고루틴이 완료될 때까지 대기 (카운터가 0이 될 때까지)
//
// 사용 시 주의사항:
// - Add()는 반드시 고루틴 시작 전에 호출해야 합니다
// - Done()은 반드시 호출되어야 합니다 (defer 사용 권장)
// - WaitGroup은 재사용 가능하지만, Wait() 중에 Add()를 호출하면 안됩니다
// =============================================================================

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	// 랜덤 시드 설정 (Go 1.20 이전 버전 호환)
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== 기본 WaitGroup 예제 ===")
	basicWaitGroupExample()

	fmt.Println("\n=== 작업 분배 예제 ===")
	workDistributionExample()

	fmt.Println("\n=== 결과 수집 예제 ===")
	resultCollectionExample()
}

// basicWaitGroupExample 기본적인 WaitGroup 사용법
func basicWaitGroupExample() {
	// WaitGroup 생성
	var wg sync.WaitGroup

	// 3개의 고루틴을 실행할 예정이므로 카운터를 3으로 설정
	// 중요: Add()는 반드시 고루틴 시작 전에 호출!
	wg.Add(3)

	// 첫 번째 고루틴
	go func() {
		// defer를 사용하면 함수가 어떻게 종료되든 Done()이 호출됨
		defer wg.Done()
		fmt.Println("고루틴 1: 시작")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("고루틴 1: 완료")
	}()

	// 두 번째 고루틴
	go func() {
		defer wg.Done()
		fmt.Println("고루틴 2: 시작")
		time.Sleep(200 * time.Millisecond)
		fmt.Println("고루틴 2: 완료")
	}()

	// 세 번째 고루틴
	go func() {
		defer wg.Done()
		fmt.Println("고루틴 3: 시작")
		time.Sleep(150 * time.Millisecond)
		fmt.Println("고루틴 3: 완료")
	}()

	// 모든 고루틴이 완료될 때까지 대기
	// 카운터가 0이 될 때까지 블로킹됨
	wg.Wait()

	fmt.Println("모든 고루틴 완료!")
}

// workDistributionExample 작업을 여러 고루틴에 분배하는 예제
func workDistributionExample() {
	// 처리할 작업들
	tasks := []string{"이메일 전송", "로그 기록", "캐시 갱신", "알림 발송", "통계 수집"}

	var wg sync.WaitGroup

	// 각 작업마다 고루틴 생성
	for i, task := range tasks {
		// 루프 내에서 Add(1)을 호출하는 패턴도 가능
		wg.Add(1)

		// 중요: 루프 변수를 고루틴에 전달할 때는 파라미터로 전달해야 함
		// Go 1.22 이전 버전에서는 클로저가 루프 변수를 캡처하면 문제 발생 가능
		go func(taskID int, taskName string) {
			defer wg.Done()

			// 작업 수행 시뮬레이션
			duration := time.Duration(rand.Intn(500)) * time.Millisecond
			fmt.Printf("작업 %d (%s): 시작 (예상 소요시간: %v)\n", taskID, taskName, duration)

			time.Sleep(duration)

			fmt.Printf("작업 %d (%s): 완료!\n", taskID, taskName)
		}(i+1, task) // 루프 변수를 파라미터로 전달
	}

	wg.Wait()
	fmt.Println("모든 작업 완료!")
}

// resultCollectionExample 고루틴의 결과를 수집하는 예제
// WaitGroup과 채널을 함께 사용하는 패턴
func resultCollectionExample() {
	numbers := []int{2, 4, 6, 8, 10}

	var wg sync.WaitGroup

	// 결과를 수집할 채널
	// 버퍼 크기를 작업 수와 동일하게 설정하면 블로킹 없이 결과 전송 가능
	results := make(chan int, len(numbers))

	// 각 숫자를 제곱하는 작업을 고루틴으로 처리
	for _, num := range numbers {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			// 계산 수행
			result := n * n

			fmt.Printf("%d의 제곱: %d\n", n, result)

			// 결과를 채널로 전송
			results <- result
		}(num)
	}

	// 별도의 고루틴에서 모든 작업이 끝나면 채널을 닫음
	// 이렇게 해야 메인 고루틴에서 range로 결과를 읽을 수 있음
	go func() {
		wg.Wait()
		close(results) // 모든 작업 완료 후 채널 닫기
	}()

	// 결과 수집 및 합계 계산
	total := 0
	for result := range results {
		total += result
	}

	fmt.Printf("제곱의 합계: %d\n", total)
}
