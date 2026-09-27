// =============================================================================
// Context 패키지 예제
// =============================================================================
// context 패키지는 고루틴 간의 취소 신호, 타임아웃, 데드라인, 값 전달을 위한 도구입니다.
//
// 주요 함수:
// - context.Background(): 최상위 컨텍스트 (루트)
// - context.TODO(): 아직 어떤 컨텍스트를 사용할지 모를 때
// - context.WithCancel(parent): 취소 가능한 컨텍스트
// - context.WithTimeout(parent, duration): 타임아웃 컨텍스트
// - context.WithDeadline(parent, time): 데드라인 컨텍스트
// - context.WithValue(parent, key, value): 값을 담은 컨텍스트
//
// 주요 메서드:
// - ctx.Done(): 컨텍스트가 취소되면 닫히는 채널 반환
// - ctx.Err(): 취소 이유 (Canceled 또는 DeadlineExceeded)
// - ctx.Deadline(): 데드라인 시간과 존재 여부 반환
// - ctx.Value(key): 저장된 값 조회
//
// 사용 규칙:
// - Context는 함수의 첫 번째 파라미터로 전달
// - nil context를 전달하지 말고, 모르겠으면 context.TODO() 사용
// - Context는 구조체에 저장하지 말고, 함수 파라미터로 전달
// =============================================================================

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== WithCancel 예제 ===")
	withCancelExample()

	fmt.Println("\n=== WithTimeout 예제 ===")
	withTimeoutExample()

	fmt.Println("\n=== WithDeadline 예제 ===")
	withDeadlineExample()

	fmt.Println("\n=== WithValue 예제 ===")
	withValueExample()

	fmt.Println("\n=== 실전 예제: HTTP 요청 취소 시뮬레이션 ===")
	httpCancellationExample()
}

// withCancelExample WithCancel 예제 - 수동으로 취소하는 컨텍스트
func withCancelExample() {
	// 취소 가능한 컨텍스트 생성
	// cancel 함수를 호출하면 ctx.Done() 채널이 닫힘
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	// 워커 고루틴 시작
	wg.Add(1)
	go func() {
		defer wg.Done()
		worker(ctx, "Worker-1")
	}()

	// 1초 후 취소
	time.Sleep(1 * time.Second)
	fmt.Println("메인: 취소 신호 전송")
	cancel() // 컨텍스트 취소

	wg.Wait()
	fmt.Println("메인: 모든 워커 종료 확인")
}

// worker 컨텍스트의 취소를 감지하는 워커
func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			// 컨텍스트가 취소됨
			fmt.Printf("%s: 취소 감지 (이유: %v)\n", name, ctx.Err())
			return

		default:
			// 정상 작업 수행
			fmt.Printf("%s: 작업 중...\n", name)
			time.Sleep(200 * time.Millisecond)
		}
	}
}

// withTimeoutExample WithTimeout 예제 - 자동으로 타임아웃되는 컨텍스트
func withTimeoutExample() {
	// 500ms 후에 자동으로 취소되는 컨텍스트
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)

	// cancel은 리소스 정리를 위해 항상 호출해야 함 (defer 사용)
	defer cancel()

	// 오래 걸리는 작업 시뮬레이션
	result := make(chan string, 1)

	go func() {
		// 랜덤한 시간 (0~800ms) 소요
		processingTime := time.Duration(rand.Intn(800)) * time.Millisecond
		fmt.Printf("작업 시작 (예상 소요: %v, 타임아웃: 500ms)\n", processingTime)
		time.Sleep(processingTime)
		result <- "작업 완료!"
	}()

	select {
	case res := <-result:
		fmt.Printf("결과: %s\n", res)

	case <-ctx.Done():
		// 타임아웃 발생
		fmt.Printf("타임아웃! (에러: %v)\n", ctx.Err())
	}
}

// withDeadlineExample WithDeadline 예제 - 특정 시간에 만료되는 컨텍스트
func withDeadlineExample() {
	// 현재 시간 + 300ms를 데드라인으로 설정
	deadline := time.Now().Add(300 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// 데드라인 정보 확인
	if d, ok := ctx.Deadline(); ok {
		fmt.Printf("데드라인: %v (남은 시간: %v)\n",
			d.Format("15:04:05.000"),
			time.Until(d))
	}

	// 여러 단계의 작업 수행
	for i := 1; i <= 5; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("단계 %d에서 데드라인 초과! (에러: %v)\n", i, ctx.Err())
			return

		default:
			fmt.Printf("단계 %d 처리 중...\n", i)
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Println("모든 단계 완료!")
}

// =============================================================================
// WithValue 관련
// =============================================================================

// 컨텍스트 키 타입 정의 (충돌 방지를 위해 커스텀 타입 사용)
type contextKey string

const (
	requestIDKey contextKey = "requestID"
	userIDKey    contextKey = "userID"
)

// withValueExample WithValue 예제 - 컨텍스트에 값 저장
func withValueExample() {
	// 기본 컨텍스트에 값을 추가
	ctx := context.Background()

	// 값 추가 (체이닝)
	ctx = context.WithValue(ctx, requestIDKey, "REQ-12345")
	ctx = context.WithValue(ctx, userIDKey, 42)

	// 값을 사용하는 함수 호출
	processRequest(ctx)
}

// processRequest 컨텍스트에서 값을 읽어 사용하는 함수
func processRequest(ctx context.Context) {
	// 값 조회
	// 타입 단언(type assertion) 필요
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		fmt.Printf("요청 ID: %s\n", requestID)
	}

	if userID, ok := ctx.Value(userIDKey).(int); ok {
		fmt.Printf("사용자 ID: %d\n", userID)
	}

	// 존재하지 않는 키 조회
	if value := ctx.Value("nonexistent"); value == nil {
		fmt.Println("존재하지 않는 키: nil 반환")
	}
}

// =============================================================================
// 실전 예제: HTTP 요청 취소 시뮬레이션
// =============================================================================

// fetchData 외부 API 호출 시뮬레이션
func fetchData(ctx context.Context, source string) (string, error) {
	// 결과 채널
	resultCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		// API 호출 시뮬레이션 (100~400ms 소요)
		delay := time.Duration(100+rand.Intn(300)) * time.Millisecond
		time.Sleep(delay)

		// 10% 확률로 에러 발생
		if rand.Float32() < 0.1 {
			errCh <- fmt.Errorf("%s: 연결 실패", source)
			return
		}

		resultCh <- fmt.Sprintf("%s: 데이터 (소요: %v)", source, delay)
	}()

	// 컨텍스트와 결과/에러를 함께 처리
	select {
	case <-ctx.Done():
		// 컨텍스트가 취소됨 (타임아웃 또는 수동 취소)
		return "", fmt.Errorf("%s: 취소됨 (%v)", source, ctx.Err())

	case err := <-errCh:
		return "", err

	case result := <-resultCh:
		return result, nil
	}
}

// httpCancellationExample HTTP 요청 취소 시뮬레이션
func httpCancellationExample() {
	// 200ms 타임아웃 컨텍스트
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// 여러 API를 동시에 호출
	sources := []string{"API-A", "API-B", "API-C"}

	var wg sync.WaitGroup
	results := make(chan string, len(sources))
	errors := make(chan string, len(sources))

	for _, source := range sources {
		wg.Add(1)

		go func(src string) {
			defer wg.Done()

			result, err := fetchData(ctx, src)
			if err != nil {
				errors <- err.Error()
				return
			}
			results <- result
		}(source)
	}

	// 결과 수집
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// 결과 출력
	fmt.Println("성공한 요청:")
	for result := range results {
		fmt.Printf("  ✓ %s\n", result)
	}

	fmt.Println("실패한 요청:")
	for err := range errors {
		fmt.Printf("  ✗ %s\n", err)
	}
}
