package senior_topics

/*
================================================================================
Go Channel 내부 구조: Buffered vs Unbuffered (시니어 레벨)
================================================================================

면접 질문: "Buffered와 Unbuffered 채널의 내부 동작 차이를 설명해주세요"

시니어급 답변 포인트:
1. hchan 구조체 상세
2. 송/수신 대기 큐 (sendq, recvq)
3. Ring Buffer 구현
4. 동기화 메커니즘
5. 성능 특성과 선택 기준
*/

import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

// ============================================================================
// 1. hchan 구조체 (채널의 실제 구조)
// ============================================================================

/*
runtime/chan.go의 hchan 구조체:

type hchan struct {
    qcount   uint           // 현재 버퍼에 있는 요소 수
    dataqsiz uint           // 버퍼 크기 (make의 두 번째 인자)
    buf      unsafe.Pointer // 버퍼 (Ring Buffer)
    elemsize uint16         // 요소 크기
    closed   uint32         // 채널 닫힘 여부
    elemtype *_type         // 요소 타입
    sendx    uint           // 송신 인덱스 (버퍼 내 위치)
    recvx    uint           // 수신 인덱스 (버퍼 내 위치)
    recvq    waitq          // 수신 대기 goroutine 큐
    sendq    waitq          // 송신 대기 goroutine 큐
    lock     mutex          // 뮤텍스 (채널 락)
}

type waitq struct {
    first *sudog           // 대기 큐의 첫 goroutine
    last  *sudog           // 대기 큐의 마지막 goroutine
}

type sudog struct {
    g     *g               // 대기 중인 goroutine
    elem  unsafe.Pointer   // 송/수신할 데이터 포인터
    ...
}

핵심 구조:
┌─────────────────────────────────────────────────────────┐
│                        hchan                            │
├─────────────────────────────────────────────────────────┤
│ qcount: 현재 요소 수   │ dataqsiz: 버퍼 크기            │
├─────────────────────────────────────────────────────────┤
│ buf → [Ring Buffer: elem0 | elem1 | elem2 | ... ]      │
├─────────────────────────────────────────────────────────┤
│ sendx: 다음 송신 위치  │ recvx: 다음 수신 위치          │
├─────────────────────────────────────────────────────────┤
│ recvq: ──→ [sudog] ──→ [sudog] ──→ nil (수신 대기열)   │
│ sendq: ──→ [sudog] ──→ [sudog] ──→ nil (송신 대기열)   │
├─────────────────────────────────────────────────────────┤
│ lock: 뮤텍스                                            │
└─────────────────────────────────────────────────────────┘
*/

func ChannelStructureDemo() {
	fmt.Println("=== Channel 구조 데모 ===")

	// 채널 크기 확인 (hchan 포인터 크기)
	ch := make(chan int, 10)
	fmt.Printf("채널 변수 크기: %d bytes (포인터)\n", unsafe.Sizeof(ch))

	// 버퍼 크기 확인
	fmt.Printf("버퍼 용량: %d\n", cap(ch))
	fmt.Printf("현재 요소 수: %d\n", len(ch))

	// 요소 추가
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Printf("3개 추가 후 len: %d, cap: %d\n", len(ch), cap(ch))
}

// ============================================================================
// 2. Unbuffered Channel 동작 원리
// ============================================================================

/*
Unbuffered Channel (dataqsiz = 0):
- 버퍼 없음 (buf = nil)
- 송신과 수신이 직접 동기화됨 (Synchronous)
- 송신자와 수신자가 만나야 데이터 전달

동작 흐름:
1. 송신 시: 수신 대기자가 있으면 직접 전달, 없으면 sendq에 대기
2. 수신 시: 송신 대기자가 있으면 직접 받음, 없으면 recvq에 대기

직접 전달 (Direct Send):
- 데이터를 송신자 스택에서 수신자 스택으로 직접 복사
- 버퍼를 거치지 않아 효율적
*/

func UnbufferedChannelDemo() {
	fmt.Println("\n=== Unbuffered Channel 동작 ===")

	ch := make(chan int) // dataqsiz = 0

	var wg sync.WaitGroup

	// 수신자 시작
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("수신자: 대기 시작...")
		start := time.Now()
		val := <-ch // recvq에 추가되어 대기
		fmt.Printf("수신자: 값 %d 수신 (대기시간: %v)\n", val, time.Since(start))
	}()

	// 잠시 대기 (수신자가 먼저 대기하도록)
	time.Sleep(100 * time.Millisecond)

	// 송신자
	fmt.Println("송신자: 값 전송...")
	ch <- 42 // 수신자가 있으므로 직접 전달
	fmt.Println("송신자: 전송 완료")

	wg.Wait()
}

// ============================================================================
// 3. Buffered Channel 동작 원리
// ============================================================================

/*
Buffered Channel (dataqsiz > 0):
- Ring Buffer로 구현
- 버퍼가 가득 차지 않으면 송신 즉시 완료
- 버퍼가 비어있지 않으면 수신 즉시 완료

Ring Buffer 동작:
- sendx: 다음 송신 위치 (쓰기 포인터)
- recvx: 다음 수신 위치 (읽기 포인터)
- 인덱스는 dataqsiz로 모듈로 연산

예시 (버퍼 크기 4):
초기: [_][_][_][_]  sendx=0, recvx=0, qcount=0
송신1: [A][_][_][_]  sendx=1, recvx=0, qcount=1
송신2: [A][B][_][_]  sendx=2, recvx=0, qcount=2
수신1: [_][B][_][_]  sendx=2, recvx=1, qcount=1
송신3: [_][B][C][_]  sendx=3, recvx=1, qcount=2
송신4: [_][B][C][D]  sendx=0, recvx=1, qcount=3  (wrap around)
*/

func BufferedChannelDemo() {
	fmt.Println("\n=== Buffered Channel 동작 ===")

	ch := make(chan string, 3) // dataqsiz = 3

	// 버퍼가 가득 찰 때까지 블록 없이 송신
	fmt.Println("버퍼에 3개 송신 (블록 없음):")
	ch <- "A"
	fmt.Printf("  송신 'A' - len: %d, cap: %d\n", len(ch), cap(ch))
	ch <- "B"
	fmt.Printf("  송신 'B' - len: %d, cap: %d\n", len(ch), cap(ch))
	ch <- "C"
	fmt.Printf("  송신 'C' - len: %d, cap: %d\n", len(ch), cap(ch))

	// 이제 버퍼가 가득 참 - 다음 송신은 블록됨
	fmt.Println("\n버퍼 가득 참. 수신 시작:")

	fmt.Printf("  수신: '%s' - len: %d\n", <-ch, len(ch))
	fmt.Printf("  수신: '%s' - len: %d\n", <-ch, len(ch))
	fmt.Printf("  수신: '%s' - len: %d\n", <-ch, len(ch))
}

// ============================================================================
// 4. 송/수신 대기 큐 동작
// ============================================================================

/*
sendq, recvq의 역할:

1. 송신 대기 (sendq):
   - 버퍼가 가득 찼을 때 송신자를 대기시킴
   - 수신 발생 시 첫 번째 대기자를 깨움

2. 수신 대기 (recvq):
   - 버퍼가 비었을 때 수신자를 대기시킴
   - 송신 발생 시 첫 번째 대기자를 깨움

sudog 구조:
- 대기 중인 goroutine (g 포인터)
- 송/수신할 데이터 (elem 포인터)
- 다음 대기자 링크 (next)
*/

func WaitQueueDemo() {
	fmt.Println("\n=== 송/수신 대기 큐 데모 ===")

	ch := make(chan int, 1) // 버퍼 1개

	var wg sync.WaitGroup

	// 3개의 송신자 시작 (버퍼 1개이므로 2개는 대기)
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("송신자 %d: 송신 시도...\n", id)
			ch <- id
			fmt.Printf("송신자 %d: 송신 완료\n", id)
		}(i)
	}

	// 잠시 대기 (송신자들이 대기열에 들어가도록)
	time.Sleep(100 * time.Millisecond)
	fmt.Println("\n--- 수신 시작 ---")

	// 3번 수신
	for i := 1; i <= 3; i++ {
		time.Sleep(50 * time.Millisecond)
		val := <-ch
		fmt.Printf("수신: %d\n", val)
	}

	wg.Wait()
}

// ============================================================================
// 5. Channel 닫기 동작
// ============================================================================

/*
채널 닫기 (close):
1. closed 플래그를 1로 설정
2. recvq의 모든 대기자를 깨움 (zero value와 false 반환)
3. sendq의 모든 대기자를 깨움 (panic 발생!)

닫힌 채널 규칙:
- 송신: panic
- 수신: 버퍼의 남은 값 반환 후 zero value 반환
- v, ok := <-ch 에서 ok가 false면 채널 닫힘
*/

func ChannelCloseDemo() {
	fmt.Println("\n=== Channel 닫기 동작 ===")

	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	// 버퍼의 남은 값들 수신
	fmt.Println("닫힌 채널에서 수신:")
	for {
		val, ok := <-ch
		if !ok {
			fmt.Printf("  채널 닫힘 (val=%d, ok=%v)\n", val, ok)
			break
		}
		fmt.Printf("  수신: %d (ok=%v)\n", val, ok)
	}

	// for range로 수신 (채널 닫힐 때까지)
	ch2 := make(chan int, 2)
	ch2 <- 10
	ch2 <- 20
	close(ch2)

	fmt.Println("\nfor range로 수신:")
	for v := range ch2 {
		fmt.Printf("  값: %d\n", v)
	}
	fmt.Println("  for range 종료")
}

// ============================================================================
// 6. 성능 비교: Buffered vs Unbuffered
// ============================================================================

/*
성능 특성:

Unbuffered:
- 동기화 보장 (Rendezvous)
- 직접 전달로 버퍼 복사 없음
- goroutine 전환 비용 발생

Buffered:
- 비동기 통신 가능
- 버퍼 복사 비용
- 버퍼 크기만큼 송신 블록 없음

선택 기준:
- 동기화가 필요하면: Unbuffered
- 처리량이 중요하면: Buffered
- 백프레셔가 필요하면: 작은 Buffered
*/

func PerformanceComparisonDemo() {
	fmt.Println("\n=== 성능 비교 데모 ===")

	const count = 100000

	// Unbuffered
	start := time.Now()
	ch1 := make(chan int)
	go func() {
		for i := 0; i < count; i++ {
			ch1 <- i
		}
		close(ch1)
	}()
	for range ch1 {
	}
	unbufferedTime := time.Since(start)

	// Buffered (크기 100)
	start = time.Now()
	ch2 := make(chan int, 100)
	go func() {
		for i := 0; i < count; i++ {
			ch2 <- i
		}
		close(ch2)
	}()
	for range ch2 {
	}
	bufferedTime := time.Since(start)

	fmt.Printf("Unbuffered: %v\n", unbufferedTime)
	fmt.Printf("Buffered(100): %v\n", bufferedTime)
	fmt.Printf("성능 향상: %.2fx\n", float64(unbufferedTime)/float64(bufferedTime))
}

// ============================================================================
// 7. 채널 선택 가이드라인
// ============================================================================

/*
Unbuffered 사용 시점:
1. 동기화 지점이 필요할 때 (두 goroutine이 만나야 할 때)
2. 작업 완료 신호 (done 채널)
3. 리소스 접근 제어 (세마포어 대용)
4. 간단한 신호 전달

Buffered 사용 시점:
1. 생산자-소비자 속도 차이 완화
2. 버스트 트래픽 처리
3. 작업 큐 (Worker Pool)
4. 처리량 최적화

버퍼 크기 결정:
- 너무 작으면: 자주 블록, 성능 저하
- 너무 크면: 메모리 낭비, 지연 증가
- 경험적으로 결정하거나 모니터링으로 조정
*/

// 패턴 1: 동기화 지점 (Unbuffered)
func syncPointPattern() {
	done := make(chan struct{}) // 크기 0

	go func() {
		// 작업 수행
		close(done) // 완료 신호
	}()

	<-done // 완료 대기
}

// 패턴 2: 세마포어 (Buffered)
func semaphorePattern() {
	maxConcurrent := 3
	sem := make(chan struct{}, maxConcurrent)

	for i := 0; i < 10; i++ {
		sem <- struct{}{} // 슬롯 획득
		go func(id int) {
			defer func() { <-sem }() // 슬롯 반환
			// 작업 수행
		}(i)
	}
}

// 패턴 3: Worker Pool (Buffered)
func workerPoolPattern() {
	jobs := make(chan int, 100)    // 작업 큐
	results := make(chan int, 100) // 결과 큐

	// Workers
	for w := 0; w < 3; w++ {
		go func() {
			for job := range jobs {
				results <- job * 2
			}
		}()
	}

	// 작업 분배
	for i := 0; i < 10; i++ {
		jobs <- i
	}
	close(jobs)

	// 결과 수집
	for i := 0; i < 10; i++ {
		<-results
	}
}

// ============================================================================
// 8. 시니어 면접 답변 요약
// ============================================================================

/*
Q: Buffered와 Unbuffered 채널의 내부 동작 차이를 설명해주세요.

A: 두 채널 모두 hchan 구조체로 구현되지만 동작 방식이 다릅니다.

1. hchan 구조체:
   - buf: Ring Buffer (Buffered만 사용)
   - sendx/recvx: 버퍼 인덱스
   - sendq/recvq: 대기 goroutine 큐
   - lock: 동기화 뮤텍스

2. Unbuffered (dataqsiz = 0):
   - 버퍼 없음, 직접 데이터 전달
   - 송신자와 수신자가 만나야 통신 완료
   - 동기화 보장 (Rendezvous point)
   - 직접 전달로 복사 비용 최소화

3. Buffered (dataqsiz > 0):
   - Ring Buffer에 데이터 저장
   - 버퍼 여유 있으면 송신 즉시 완료
   - 버퍼에 데이터 있으면 수신 즉시 완료
   - 비동기 통신 가능

4. 성능 특성:
   - Unbuffered: goroutine 전환 비용, 직접 전달
   - Buffered: 버퍼 복사 비용, 전환 감소

5. 선택 기준:
   - 동기화 필요: Unbuffered
   - 처리량 중요: Buffered
   - 백프레셔: 작은 Buffered
*/
