// =============================================================================
// Channel (채널) 제어 예제
// =============================================================================
// 채널은 고루틴 간의 통신을 위한 Go의 핵심 도구입니다.
// "공유 메모리로 통신하지 말고, 통신으로 메모리를 공유하라" - Go 철학
//
// 채널의 종류:
// - 언버퍼드 채널 (make(chan T)): 송수신이 동기화됨
// - 버퍼드 채널 (make(chan T, n)): 버퍼가 가득 찰 때까지 블로킹 없이 전송
//
// 채널의 방향:
// - chan T: 양방향 (송수신 모두 가능)
// - chan<- T: 송신 전용
// - <-chan T: 수신 전용
//
// 주요 연산:
// - ch <- value: 채널로 값 전송
// - value := <-ch: 채널에서 값 수신
// - close(ch): 채널 닫기
// - select: 여러 채널 연산 중 하나 선택
// =============================================================================

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== 언버퍼드 채널 예제 ===")
	unbufferedChannelExample()

	fmt.Println("\n=== 버퍼드 채널 예제 ===")
	bufferedChannelExample()

	fmt.Println("\n=== 채널 방향 제한 예제 ===")
	channelDirectionExample()

	fmt.Println("\n=== Select 문 예제 ===")
	selectExample()

	fmt.Println("\n=== 타임아웃 패턴 예제 ===")
	timeoutPatternExample()

	fmt.Println("\n=== Done 채널 패턴 예제 ===")
	doneChannelExample()
}

// unbufferedChannelExample 언버퍼드 채널 예제
// 언버퍼드 채널은 송신자와 수신자가 동기화됨
func unbufferedChannelExample() {
	// 언버퍼드 채널 생성 (버퍼 크기 지정 없음)
	ch := make(chan string)

	// 고루틴에서 메시지 전송
	go func() {
		fmt.Println("송신자: 메시지 전송 시작")

		// 수신자가 있을 때까지 블로킹됨
		ch <- "Hello, Channel!"

		// 위 전송이 완료되어야 이 줄이 실행됨
		fmt.Println("송신자: 메시지 전송 완료")
	}()

	// 잠시 대기 후 수신 (송신자가 대기하는 것을 보여주기 위해)
	time.Sleep(100 * time.Millisecond)

	fmt.Println("수신자: 메시지 수신 시작")
	message := <-ch
	fmt.Printf("수신자: 받은 메시지 = %q\n", message)
}

// bufferedChannelExample 버퍼드 채널 예제
// 버퍼가 가득 차기 전까지는 블로킹 없이 전송 가능
func bufferedChannelExample() {
	// 버퍼 크기가 3인 채널 생성
	ch := make(chan int, 3)

	// 버퍼가 비어있으면 블로킹 없이 전송 가능
	ch <- 1
	fmt.Println("1 전송 완료 (버퍼: 1/3)")

	ch <- 2
	fmt.Println("2 전송 완료 (버퍼: 2/3)")

	ch <- 3
	fmt.Println("3 전송 완료 (버퍼: 3/3)")

	// 버퍼가 가득 찼으므로 다음 전송은 블로킹됨
	// ch <- 4  // 이 줄은 데드락 발생!

	// 수신하면 버퍼에 공간이 생김
	fmt.Printf("수신: %d (버퍼: 2/3)\n", <-ch)
	fmt.Printf("수신: %d (버퍼: 1/3)\n", <-ch)
	fmt.Printf("수신: %d (버퍼: 0/3)\n", <-ch)

	// 채널 길이와 용량 확인
	fmt.Printf("채널 길이: %d, 용량: %d\n", len(ch), cap(ch))
}

// channelDirectionExample 채널 방향 제한 예제
// 함수 시그니처에서 채널의 방향을 제한할 수 있음
func channelDirectionExample() {
	ch := make(chan int)

	// producer는 송신만, consumer는 수신만 가능하도록 제한
	go producer(ch)
	consumer(ch)
}

// producer 송신 전용 채널 파라미터
// ch<- int 는 송신만 가능한 채널을 의미
func producer(ch chan<- int) {
	for i := 1; i <= 5; i++ {
		ch <- i
		fmt.Printf("Producer: %d 전송\n", i)
	}
	close(ch) // 송신 완료 후 채널 닫기
}

// consumer 수신 전용 채널 파라미터
// <-chan int 는 수신만 가능한 채널을 의미
func consumer(ch <-chan int) {
	// range를 사용하면 채널이 닫힐 때까지 반복
	for value := range ch {
		fmt.Printf("Consumer: %d 수신\n", value)
	}
	fmt.Println("Consumer: 채널 닫힘, 종료")
}

// selectExample select 문 예제
// select는 여러 채널 연산 중 준비된 것을 실행
func selectExample() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// 두 고루틴이 서로 다른 시간에 데이터 전송
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "채널 1의 메시지"
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch2 <- "채널 2의 메시지"
	}()

	// select로 두 채널 모두에서 수신
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			// ch1에서 수신 가능하면 실행
			fmt.Printf("ch1에서 수신: %s\n", msg1)

		case msg2 := <-ch2:
			// ch2에서 수신 가능하면 실행
			fmt.Printf("ch2에서 수신: %s\n", msg2)
		}
	}

	// ==========================================================================
	// Non-blocking select (default 사용)
	// ==========================================================================
	ch3 := make(chan int, 1)

	select {
	case value := <-ch3:
		fmt.Printf("받은 값: %d\n", value)
	default:
		// 어떤 채널도 준비되지 않으면 default 실행
		fmt.Println("채널에 데이터 없음")
	}

	// 데이터를 넣고 다시 시도
	ch3 <- 42

	select {
	case value := <-ch3:
		fmt.Printf("받은 값: %d\n", value)
	default:
		fmt.Println("채널에 데이터 없음")
	}
}

// timeoutPatternExample 타임아웃 패턴 예제
// 채널 연산에 시간 제한을 두는 패턴
func timeoutPatternExample() {
	ch := make(chan string)

	// 랜덤한 시간 후에 데이터 전송
	go func() {
		delay := time.Duration(rand.Intn(200)) * time.Millisecond
		fmt.Printf("작업 시작 (소요 예상: %v)\n", delay)
		time.Sleep(delay)
		ch <- "작업 완료!"
	}()

	// 100ms 타임아웃 설정
	select {
	case result := <-ch:
		fmt.Printf("결과 수신: %s\n", result)

	case <-time.After(100 * time.Millisecond):
		// time.After()는 지정된 시간 후에 현재 시간을 보내는 채널 반환
		fmt.Println("타임아웃! 작업이 너무 오래 걸립니다.")
	}
}

// doneChannelExample Done 채널 패턴 예제
// 고루틴에게 종료 신호를 보내는 패턴
func doneChannelExample() {
	// done 채널은 종료 신호를 보내는 데 사용
	done := make(chan struct{})
	messages := make(chan int)

	// 메시지를 생성하는 워커
	go func() {
		counter := 0
		for {
			select {
			case <-done:
				// done 채널이 닫히면 종료
				fmt.Println("워커: 종료 신호 수신, 종료합니다")
				close(messages) // 메시지 채널도 닫기
				return

			case messages <- counter:
				// 메시지 전송
				counter++
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	// 5개의 메시지만 수신 후 종료
	for i := 0; i < 5; i++ {
		msg := <-messages
		fmt.Printf("메시지 수신: %d\n", msg)
	}

	// 워커에게 종료 신호 보내기
	// done 채널을 닫으면 수신 대기 중인 모든 select가 즉시 실행됨
	close(done)

	// 워커가 종료될 때까지 잠시 대기
	time.Sleep(100 * time.Millisecond)
	fmt.Println("메인: 프로그램 종료")
}
