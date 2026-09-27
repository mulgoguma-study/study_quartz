// =============================================================================
// Fan-In / Fan-Out 패턴 예제
// =============================================================================
// 이 패턴은 동시성 프로그래밍에서 작업을 분배하고 결과를 수집하는 핵심 패턴입니다.
//
// Fan-Out (팬아웃): 하나의 입력을 여러 고루틴으로 분배
//   - 하나의 채널에서 여러 워커가 데이터를 가져가서 병렬 처리
//   - CPU 집약적 작업이나 I/O 작업을 병렬화할 때 유용
//
// Fan-In (팬인): 여러 고루틴의 결과를 하나로 합침
//   - 여러 채널의 출력을 하나의 채널로 합침
//   - 분산된 결과를 수집할 때 사용
//
// 일반적인 흐름:
// Producer → Fan-Out → Workers → Fan-In → Consumer
// =============================================================================

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job 처리할 작업 구조체
type Job struct {
	ID    int
	Value int
}

// Result 작업 결과 구조체
type Result struct {
	Job       Job
	Output    int
	WorkerID  int
	ProcessedAt time.Time
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== Fan-Out 예제 ===")
	fanOutExample()

	fmt.Println("\n=== Fan-In 예제 ===")
	fanInExample()

	fmt.Println("\n=== 통합 Fan-Out/Fan-In 파이프라인 예제 ===")
	fullPipelineExample()
}

// fanOutExample 팬아웃 패턴 - 하나의 작업 스트림을 여러 워커가 처리
func fanOutExample() {
	// 작업 채널 생성 (입력)
	jobs := make(chan Job, 10)

	var wg sync.WaitGroup

	// ==========================================================================
	// Fan-Out: 3개의 워커 고루틴이 같은 채널에서 작업을 가져감
	// ==========================================================================
	// 각 워커가 채널에서 작업을 경쟁적으로 가져가서 처리
	// Go의 채널은 기본적으로 thread-safe하므로 별도의 락 없이 안전하게 공유 가능
	numWorkers := 3

	for workerID := 1; workerID <= numWorkers; workerID++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			// 채널이 닫힐 때까지 작업을 계속 가져옴
			for job := range jobs {
				// 작업 처리 시뮬레이션
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				fmt.Printf("워커 %d: Job %d 처리 완료 (값: %d)\n", id, job.ID, job.Value)
			}

			fmt.Printf("워커 %d: 종료\n", id)
		}(workerID)
	}

	// ==========================================================================
	// Producer: 작업 생성 및 채널로 전송
	// ==========================================================================
	for i := 1; i <= 9; i++ {
		jobs <- Job{ID: i, Value: i * 10}
		fmt.Printf("Producer: Job %d 전송\n", i)
	}

	// 모든 작업 전송 후 채널 닫기
	// 채널을 닫으면 워커들의 range 루프가 종료됨
	close(jobs)

	// 모든 워커가 완료될 때까지 대기
	wg.Wait()
	fmt.Println("모든 작업 완료!")
}

// fanInExample 팬인 패턴 - 여러 채널을 하나로 합침
func fanInExample() {
	// ==========================================================================
	// 여러 개의 채널 생성 (각각 다른 소스에서 데이터 생성)
	// ==========================================================================
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	// 각 채널에 데이터를 보내는 고루틴
	go func() {
		for i := 1; i <= 3; i++ {
			ch1 <- i * 100
			time.Sleep(50 * time.Millisecond)
		}
		close(ch1)
	}()

	go func() {
		for i := 1; i <= 3; i++ {
			ch2 <- i * 200
			time.Sleep(70 * time.Millisecond)
		}
		close(ch2)
	}()

	go func() {
		for i := 1; i <= 3; i++ {
			ch3 <- i * 300
			time.Sleep(90 * time.Millisecond)
		}
		close(ch3)
	}()

	// ==========================================================================
	// Fan-In: 여러 채널을 하나로 합침
	// ==========================================================================
	merged := fanIn(ch1, ch2, ch3)

	// 합쳐진 채널에서 모든 값 읽기
	for value := range merged {
		fmt.Printf("수신: %d\n", value)
	}

	fmt.Println("모든 채널 병합 및 수신 완료!")
}

// fanIn 여러 채널을 하나의 채널로 합치는 함수
// 가변 인자를 사용하여 임의 개수의 채널을 받을 수 있음
func fanIn(channels ...<-chan int) <-chan int {
	// 출력 채널 생성
	merged := make(chan int)

	var wg sync.WaitGroup

	// 각 입력 채널에 대해 고루틴 생성
	for _, ch := range channels {
		wg.Add(1)

		go func(c <-chan int) {
			defer wg.Done()

			// 채널의 모든 값을 merged 채널로 전달
			for value := range c {
				merged <- value
			}
		}(ch)
	}

	// 모든 입력 채널이 닫히면 merged 채널도 닫음
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

// fullPipelineExample 완전한 Fan-Out/Fan-In 파이프라인
func fullPipelineExample() {
	// ==========================================================================
	// 1단계: Producer - 작업 생성
	// ==========================================================================
	jobs := make(chan Job, 10)

	go func() {
		for i := 1; i <= 12; i++ {
			jobs <- Job{ID: i, Value: rand.Intn(100)}
		}
		close(jobs)
	}()

	// ==========================================================================
	// 2단계: Fan-Out - 4개의 워커가 작업 처리
	// ==========================================================================
	numWorkers := 4

	// 각 워커의 결과 채널을 저장
	workerResults := make([]<-chan Result, numWorkers)

	for i := 0; i < numWorkers; i++ {
		// 각 워커가 자신만의 결과 채널을 가짐
		workerResults[i] = processJobs(jobs, i+1)
	}

	// ==========================================================================
	// 3단계: Fan-In - 모든 워커의 결과를 하나로 합침
	// ==========================================================================
	mergedResults := fanInResults(workerResults...)

	// ==========================================================================
	// 4단계: Consumer - 결과 수집 및 처리
	// ==========================================================================
	var totalProcessed int
	var totalOutput int

	for result := range mergedResults {
		totalProcessed++
		totalOutput += result.Output
		fmt.Printf("결과 수신: Job %d → %d (워커 %d가 처리, 시간: %v)\n",
			result.Job.ID, result.Output, result.WorkerID,
			result.ProcessedAt.Format("15:04:05.000"))
	}

	fmt.Printf("\n총 처리된 작업: %d개, 출력 합계: %d\n", totalProcessed, totalOutput)
}

// processJobs 작업을 처리하고 결과 채널을 반환하는 워커
// 각 워커가 독립적인 결과 채널을 가지므로 동시성 문제가 없음
func processJobs(jobs <-chan Job, workerID int) <-chan Result {
	results := make(chan Result)

	go func() {
		defer close(results)

		for job := range jobs {
			// 작업 처리 시뮬레이션
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

			// 결과 생성 (여기서는 값을 제곱)
			results <- Result{
				Job:       job,
				Output:    job.Value * job.Value,
				WorkerID:  workerID,
				ProcessedAt: time.Now(),
			}
		}
	}()

	return results
}

// fanInResults 여러 Result 채널을 하나로 합침
func fanInResults(channels ...<-chan Result) <-chan Result {
	merged := make(chan Result)

	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)

		go func(c <-chan Result) {
			defer wg.Done()

			for result := range c {
				merged <- result
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}
