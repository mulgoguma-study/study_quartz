package senior_topics

/*
================================================================================
Go 메모리 누수 방지 전략 (시니어 레벨)
================================================================================

면접 질문: "Go에서 메모리 누수가 발생하는 패턴과 방지법을 설명해주세요"

시니어급 답변 포인트:
1. Goroutine 누수
2. 채널 누수
3. 슬라이스/맵 메모리 누수
4. time.Ticker/Timer 누수
5. HTTP 응답 Body 누수
6. 파일/리소스 핸들 누수
7. pprof를 통한 메모리 분석
*/

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// ============================================================================
// 1. Goroutine 누수 (가장 흔한 패턴)
// ============================================================================

/*
Goroutine 누수 원인:
1. 채널 수신 대기 중 송신자가 사라짐
2. 채널 송신 대기 중 수신자가 사라짐
3. 무한 루프에서 탈출 조건 없음
4. context 취소 미처리
*/

// 나쁜 예: Goroutine이 영원히 대기
func BadGoroutineLeak() {
	ch := make(chan int)

	go func() {
		val := <-ch // 영원히 대기 - 누수!
		fmt.Println(val)
	}()

	// ch에 값을 보내지 않고 함수 종료
	// goroutine은 영원히 살아있음
}

// 좋은 예: context로 취소 가능하게
func GoodGoroutineWithContext(ctx context.Context) {
	ch := make(chan int)

	go func() {
		select {
		case val := <-ch:
			fmt.Println(val)
		case <-ctx.Done():
			fmt.Println("goroutine 정상 종료")
			return
		}
	}()
}

// 좋은 예: done 채널 패턴
func GoodGoroutineWithDone() {
	ch := make(chan int)
	done := make(chan struct{})

	go func() {
		select {
		case val := <-ch:
			fmt.Println(val)
		case <-done:
			return
		}
	}()

	// 작업 완료 후
	close(done)
}

// Goroutine 누수 감지 함수
func CountGoroutines() int {
	return runtime.NumGoroutine()
}

func GoroutineLeakDemo() {
	fmt.Println("=== Goroutine 누수 데모 ===")

	before := CountGoroutines()
	fmt.Printf("시작 전 goroutine 수: %d\n", before)

	// 누수 발생
	for i := 0; i < 10; i++ {
		ch := make(chan int)
		go func() {
			<-ch // 영원히 대기
		}()
	}

	time.Sleep(100 * time.Millisecond)
	after := CountGoroutines()
	fmt.Printf("누수 후 goroutine 수: %d (누수: %d개)\n", after, after-before)
}

// ============================================================================
// 2. 채널 누수 패턴
// ============================================================================

// 나쁜 예: Unbuffered 채널 송신 대기
func BadChannelLeak() {
	ch := make(chan int) // unbuffered

	go func() {
		ch <- 1 // 수신자 없으면 영원히 대기
	}()

	// ch를 수신하지 않고 종료
}

// 좋은 예: Buffered 채널 또는 select with default
func GoodChannelPattern() {
	ch := make(chan int, 1) // buffered

	go func() {
		select {
		case ch <- 1:
			// 성공
		default:
			// 버퍼 가득 차면 즉시 리턴
		}
	}()
}

// 좋은 예: 타임아웃과 함께
func GoodChannelWithTimeout() {
	ch := make(chan int)

	go func() {
		select {
		case ch <- 1:
			// 성공
		case <-time.After(1 * time.Second):
			fmt.Println("타임아웃 - goroutine 종료")
			return
		}
	}()
}

// ============================================================================
// 3. time.Ticker / time.Timer 누수
// ============================================================================

/*
time.Ticker 누수:
- Ticker는 명시적으로 Stop() 호출 필요
- Stop() 안 하면 내부 goroutine이 계속 실행
*/

// 나쁜 예
func BadTickerLeak() {
	ticker := time.NewTicker(1 * time.Second)
	// ticker.Stop() 호출 안 함 - 누수!

	for i := 0; i < 3; i++ {
		<-ticker.C
		fmt.Println("tick")
	}
	// ticker가 계속 실행 중
	_ = ticker
}

// 좋은 예
func GoodTickerUsage(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() // 반드시 Stop!

	for {
		select {
		case <-ticker.C:
			fmt.Println("tick")
		case <-ctx.Done():
			return
		}
	}
}

// time.After 주의사항
func TimeAfterWarning() {
	fmt.Println("\n=== time.After 주의사항 ===")

	// 나쁜 예: 루프에서 time.After
	// 매번 새 Timer 생성 - 이전 Timer가 GC 대상이 되려면 만료까지 대기
	/*
		for {
			select {
			case <-ch:
				// 처리
			case <-time.After(5 * time.Second):  // 매 루프마다 새 Timer!
				// 타임아웃
			}
		}
	*/

	// 좋은 예: Timer 재사용
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	ch := make(chan int)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- 1
	}()

	select {
	case v := <-ch:
		if !timer.Stop() {
			<-timer.C // 채널 비우기
		}
		fmt.Printf("값 수신: %d\n", v)
	case <-timer.C:
		fmt.Println("타임아웃")
	}
}

// ============================================================================
// 4. HTTP Response Body 누수
// ============================================================================

/*
HTTP 응답 Body 누수:
- resp.Body를 Close하지 않으면 TCP 연결이 재사용 안됨
- 연결 풀 고갈로 이어짐
*/

// 나쁜 예
func BadHTTPUsage() error {
	resp, err := http.Get("https://example.com")
	if err != nil {
		return err
	}
	// resp.Body.Close() 안 함 - 연결 누수!

	_ = resp
	return nil
}

// 좋은 예
func GoodHTTPUsage() error {
	resp, err := http.Get("https://example.com")
	if err != nil {
		return err
	}
	defer resp.Body.Close() // 반드시 Close!

	// Body를 완전히 읽어야 연결 재사용 가능
	_, _ = io.Copy(io.Discard, resp.Body)

	return nil
}

// 더 나은 예: 에러 시에도 Body Close
func BetterHTTPUsage() error {
	resp, err := http.Get("https://example.com")
	if err != nil {
		return err
	}
	defer func() {
		// Body 내용을 drain해야 연결이 재사용됨
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	// 정상 처리
	return nil
}

// ============================================================================
// 5. sync.Pool 미사용으로 인한 메모리 압박
// ============================================================================

/*
sync.Pool 활용:
- 빈번하게 할당/해제되는 객체 재사용
- GC 부담 감소
*/

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 4096)
	},
}

func ProcessWithPool(data []byte) {
	// Pool에서 버퍼 획득
	buf := bufferPool.Get().([]byte)
	buf = buf[:0] // 길이 초기화

	defer func() {
		// Pool에 반환
		bufferPool.Put(buf)
	}()

	// 버퍼 사용
	buf = append(buf, data...)
	// 처리...
}

// ============================================================================
// 6. Slice 메모리 누수 패턴 (상세)
// ============================================================================

func SliceLeakPatterns() {
	fmt.Println("\n=== Slice 메모리 누수 패턴 ===")

	// 패턴 1: append 후 이전 배열 참조 유지
	data := make([]int, 1000000)
	// 새 배열이 할당되어도 data는 이전 배열을 참조
	// 이전 배열이 GC되지 않음
	newData := append(data, 1)
	_ = newData
	fmt.Println("1. append로 새 배열 생성 시 이전 배열 주의")

	// 패턴 2: 문자열에서 추출한 substring
	// Go 1.18 이전: 원본 문자열 메모리 유지
	longString := "very long string that takes up memory..."
	shortPart := longString[:4]

	// Go 1.18+에서는 이렇게 복사
	shortCopy := string([]byte(longString[:4]))
	_ = shortPart
	_ = shortCopy
	fmt.Println("2. 문자열 슬라이싱 시 복사 필요")

	// 패턴 3: 구조체 슬라이스에서 포인터 필드
	type Item struct {
		data *[1024]byte
	}

	items := make([]Item, 100)
	for i := range items {
		items[i].data = new([1024]byte)
	}

	// 앞 10개만 유지하면서 메모리 누수 방지
	for i := 10; i < len(items); i++ {
		items[i].data = nil // 포인터 해제
	}
	items = items[:10]
	fmt.Println("3. 포인터 필드는 nil로 설정 후 슬라이싱")
}

// ============================================================================
// 7. Map 메모리 누수 패턴
// ============================================================================

func MapLeakPatterns() {
	fmt.Println("\n=== Map 메모리 누수 패턴 ===")

	// 패턴 1: Map은 축소되지 않음
	m := make(map[int][]byte)

	// 큰 값들 추가
	for i := 0; i < 10000; i++ {
		m[i] = make([]byte, 1024)
	}

	// 삭제해도 버킷 메모리는 유지
	for i := 0; i < 10000; i++ {
		delete(m, i)
	}
	// m의 버킷 메모리는 그대로!

	// 해결책: 새 map으로 교체
	m = make(map[int][]byte)
	fmt.Println("1. 대량 삭제 후 새 map으로 교체 필요")

	// 패턴 2: Map에 큰 값 대신 포인터 저장
	type BigStruct struct {
		data [1024 * 1024]byte
	}

	// 나쁜 예: 값으로 저장 (복사 발생)
	// badMap := make(map[int]BigStruct)

	// 좋은 예: 포인터로 저장
	goodMap := make(map[int]*BigStruct)
	_ = goodMap
	fmt.Println("2. 큰 구조체는 포인터로 저장")
}

// ============================================================================
// 8. defer 관련 메모리 누수
// ============================================================================

func DeferLeakPatterns() {
	fmt.Println("\n=== defer 메모리 누수 패턴 ===")

	// 나쁜 예: 루프 안에서 defer
	// 함수 끝까지 모든 defer가 스택에 쌓임
	/*
		for _, file := range files {
			f, _ := os.Open(file)
			defer f.Close()  // 루프 끝까지 모든 파일이 열린 상태!
		}
	*/

	// 좋은 예: 즉시 처리 함수로 분리
	/*
		for _, file := range files {
			func() {
				f, _ := os.Open(file)
				defer f.Close()
				// 처리
			}()  // 각 반복마다 defer 실행
		}
	*/
	fmt.Println("루프 안 defer는 즉시 실행 함수로 감싸기")
}

// ============================================================================
// 9. 메모리 누수 탐지 도구 사용법
// ============================================================================

/*
pprof 사용법:

1. HTTP 서버에 pprof 추가:
   import _ "net/http/pprof"
   go func() { http.ListenAndServe(":6060", nil) }()

2. 프로파일 수집:
   go tool pprof http://localhost:6060/debug/pprof/heap
   go tool pprof http://localhost:6060/debug/pprof/goroutine

3. 분석 명령:
   - top: 메모리 사용량 상위 항목
   - list <func>: 특정 함수의 메모리 할당
   - web: 시각화 (graphviz 필요)

4. 차이 분석:
   go tool pprof -base=before.prof after.prof
*/

func PrintMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("\n=== 메모리 통계 ===\n")
	fmt.Printf("Alloc (현재 할당): %d MB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc (총 할당): %d MB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys (시스템 획득): %d MB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC (GC 횟수): %d\n", m.NumGC)
	fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
}

// ============================================================================
// 10. 체크리스트: 메모리 누수 방지
// ============================================================================

/*
메모리 누수 방지 체크리스트:

□ Goroutine
  - context 또는 done 채널로 종료 가능한가?
  - 무한 루프에 탈출 조건이 있는가?
  - 채널 송/수신에 타임아웃이 있는가?

□ Channel
  - Unbuffered 채널의 상대방이 항상 존재하는가?
  - 채널 close 후 수신측에서 처리하는가?

□ Timer/Ticker
  - Ticker.Stop()을 호출하는가?
  - 루프에서 time.After 대신 Timer를 재사용하는가?

□ HTTP
  - Response.Body를 Close하는가?
  - Body를 완전히 읽어 연결을 재사용하는가?

□ Slice/Map
  - 큰 배열의 작은 슬라이스를 유지하지 않는가?
  - 포인터 슬라이스 축소 시 nil 처리하는가?
  - Map 대량 삭제 후 재생성하는가?

□ 파일/리소스
  - 모든 리소스에 defer Close가 있는가?
  - 루프 안 defer를 함수로 분리했는가?

□ 모니터링
  - pprof endpoint가 설정되어 있는가?
  - 정기적으로 goroutine 수를 모니터링하는가?
*/
