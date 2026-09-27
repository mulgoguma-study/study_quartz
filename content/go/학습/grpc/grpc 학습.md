Go로 gRPC 개발할 때 개념이 잘 안 잡히는 부분을 정리해 드릴게요. gRPC는 **Protocol Buffers(Proto)**로 서비스와 메시지를 정의하고, 이를 기반으로 서버와 클라이언트 코드를 자동 생성하는 구조예요.

---

## 1. 핵심 개념 정리

```
┌─────────────────────────────────────────────────────────┐
│  .proto 파일 (서비스 계약 정의)                          │
│  ─────────────────────────────                          │
│  • 어떤 메서드가 있는지                                   │
│  • 요청/응답 메시지 구조                                  │
│  • 서비스 인터페이스 정의                                 │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼ (protoc 컴파일러로 코드 생성)
┌─────────────────┐    ┌─────────────────┐
│   gRPC 서버      │◄──►│   gRPC 클라이언트 │
│  (Go 코드)       │    │   (Go 코드)      │
└─────────────────┘    └─────────────────┘
```

---

## 2. Proto 파일 작성 (.proto)

```protobuf
// calculator.proto
syntax = "proto3";

option go_package = "github.com/yourname/calculator/proto";

package calculator;

// 서비스 정의: 클라이언트가 호출할 수 있는 메서드들
service Calculator {
  // Unary RPC: 1요청 → 1응답 (가장 기본)
  rpc Add (AddRequest) returns (AddResponse);
  
  // Server Streaming: 1요청 → N응답
  rpc CountDown (CountRequest) returns (stream Number);
  
  // Client Streaming: N요청 → 1응답
  rpc SumNumbers (stream Number) returns (SumResponse);
  
  // Bidirectional Streaming: N요청 ↔ N응답
  rpc Chat (stream Message) returns (stream Message);
}

// 메시지 정의: 요청/응답 데이터 구조
message AddRequest {
  int32 a = 1;  // 필드 번호 1
  int32 b = 2;  // 필드 번호 2
}

message AddResponse {
  int32 result = 1;
}

message CountRequest {
  int32 start = 1;
}

message Number {
  int32 value = 1;
}

message SumResponse {
  int32 total = 1;
}

message Message {
  string text = 1;
  string sender = 2;
}
```

### Proto 문법 핵심 포인트

| 요소 | 설명 | 예시 |
|------|------|------|
| `service` | 서비스 인터페이스 정의 | `service Calculator` |
| `rpc` | 메서드 선언 | `rpc Add (Request) returns (Response)` |
| `message` | 데이터 구조 | `message AddRequest` |
| `stream` | 스트리밍 표시 | `returns (stream Number)` |
| 필드 번호 | 바이너리 인코딩용 ID | `int32 a = 1` (1-536,870,911) |

---

## 3. Go 코드 생성

### 필요한 도구 설치
```bash
# 1. protoc 컴파일러 설치 (OS별로 다름)
# macOS: brew install protobuf
# Ubuntu: apt-get install protobuf-compiler

# 2. Go 플러그인 설치
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# PATH에 추가 (zsh 기준)
export PATH="$PATH:$(go env GOPATH)/bin"
```

### 코드 생성 명령어
```bash
# calculator.proto 파일이 있는 디렉토리에서 실행
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  calculator.proto
```

**생성되는 파일들:**
- `calculator.pb.go` - 메시지 구조체 (Request/Response)
- `calculator_grpc.pb.go` - gRPC 인터페이스 (서버/클라이언트)

---

## 4. 서버 구현 (Server)

```go
package main

import (
	"context"
	"io"
	"log"
	"net"

	pb "github.com/yourname/calculator/proto"
	"google.golang.org/grpc"
)

// 1. 서버 구조체 정의: 생성된 인터페이스 구현
type server struct {
	pb.UnimplementedCalculatorServer // 필수: forward compatibility
}

// 2. Unary RPC 구현 (1:1)
func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	log.Printf("Received: %d + %d", req.A, req.B)
	return &pb.AddResponse{
		Result: req.A + req.B,
	}, nil
}

// 3. Server Streaming 구현 (1:N)
func (s *server) CountDown(req *pb.CountRequest, stream pb.Calculator_CountDownServer) error {
	log.Printf("CountDown from %d", req.Start)
	
	for i := req.Start; i >= 0; i-- {
		// 클라이언트로 계속 전송
		if err := stream.Send(&pb.Number{Value: i}); err != nil {
			return err
		}
	}
	return nil // 스트림 종료
}

// 4. Client Streaming 구현 (N:1)
func (s *server) SumNumbers(stream pb.Calculator_SumNumbersServer) error {
	var total int32
	
	for {
		num, err := stream.Recv() // 클라이언트로부터 계속 수신
		if err == io.EOF {
			// 클라이언트가 스트림을 닫음 → 최종 응답 전송
			return stream.SendAndClose(&pb.SumResponse{Total: total})
		}
		if err != nil {
			return err
		}
		total += num.Value
	}
}

// 5. Bidirectional Streaming 구현 (N:N)
func (s *server) Chat(stream pb.Calculator_ChatServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		
		log.Printf("Received from %s: %s", msg.Sender, msg.Text)
		
		// 에코 응답
		response := &pb.Message{
			Text:   "Echo: " + msg.Text,
			Sender: "Server",
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}

func main() {
	// TCP 리스너 생성
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// gRPC 서버 생성
	s := grpc.NewServer()
	
	// 서비스 등록
	pb.RegisterCalculatorServer(s, &server{})

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

---

## 5. 클라이언트 구현 (Client)

```go
package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/yourname/calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. 서버 연결 (insecure는 개발용, 운영에서는 TLS 사용)
	conn, err := grpc.Dial("localhost:50051", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 2. 클라이언트 생성
	client := pb.NewCalculatorClient(conn)

	// 3. 각 RPC 타입별 호출
	callUnary(client)
	callServerStreaming(client)
	callClientStreaming(client)
	callBidirectionalStreaming(client)
}

// Unary RPC 호출
func callUnary(client pb.CalculatorClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.Add(ctx, &pb.AddRequest{A: 10, B: 20})
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}
	log.Printf("Unary Result: %d", resp.Result)
}

// Server Streaming 호출
func callServerStreaming(client pb.CalculatorClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := client.CountDown(ctx, &pb.CountRequest{Start: 5})
	if err != nil {
		log.Fatalf("could not count down: %v", err)
	}

	for {
		num, err := stream.Recv()
		if err == io.EOF {
			break // 스트림 종료
		}
		if err != nil {
			log.Fatalf("error receiving: %v", err)
		}
		log.Printf("Server Streaming: %d", num.Value)
	}
}

// Client Streaming 호출
func callClientStreaming(client pb.CalculatorClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := client.SumNumbers(ctx)
	if err != nil {
		log.Fatalf("could not sum: %v", err)
	}

	// 숫자들을 서버로 전송
	numbers := []int32{1, 2, 3, 4, 5}
	for _, n := range numbers {
		if err := stream.Send(&pb.Number{Value: n}); err != nil {
			log.Fatalf("error sending: %v", err)
		}
	}

	// 스트림 닫고 응답 받기
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error receiving response: %v", err)
	}
	log.Printf("Client Streaming Sum: %d", resp.Total)
}

// Bidirectional Streaming 호출
func callBidirectionalStreaming(client pb.CalculatorClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := client.Chat(ctx)
	if err != nil {
		log.Fatalf("could not chat: %v", err)
	}

	// 수신 고루틴
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("error receiving: %v", err)
				return
			}
			log.Printf("Bidirectional Received: %s from %s", msg.Text, msg.Sender)
		}
	}()

	// 전송
	messages := []string{"Hello", "World", "gRPC"}
	for _, text := range messages {
		if err := stream.Send(&pb.Message{
			Text:   text,
			Sender: "Client",
		}); err != nil {
			log.Fatalf("error sending: %v", err)
		}
		time.Sleep(time.Second)
	}

	stream.CloseSend()
	time.Sleep(time.Second) // 수신 고루틴이 끝날 때까지 대기
	cancel()
}
```

---

## 6. 프로젝트 구조 예시

```
my-grpc-project/
├── proto/
│   ├── calculator.proto      # 서비스 정의
│   ├── calculator.pb.go      # 자동 생성 (메시지)
│   └── calculator_grpc.pb.go # 자동 생성 (인터페이스)
├── server/
│   └── main.go               # 서버 구현
├── client/
│   └── main.go               # 클라이언트 구현
├── go.mod
└── Makefile                  # protoc 명령어 자동화
```

**Makefile 예시:**
```makefile
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/calculator.proto

server:
	go run server/main.go

client:
	go run client/main.go
```

---

## 7. 핵심 개념 체크리스트

| 개념 | 정리 |
|------|------|
| **Proto** | 서비스 계약(Contract) → 언어 중립적 IDL |
| **protoc** | Proto → Go 코드 변환 컴파일러 |
| **Server** | 인터페이스 구현 → `RegisterXXXServer`로 등록 |
| **Client** | `NewXXXClient(conn)`로 stub 생성 후 메서드 호출 |
| **Context** | 타임아웃, 메타데이터, 취소 신호 전달 |
| **Stream** | `Send()`/`Recv()`로 양방향 데이터 흐름 제어 |

---

## 8. 자주 하는 실수

1. **필드 번호 충돌**: proto 수정 시 기존 번호 재사용 금지
2. **context 무시**: 항상 timeout 설정 (무한 대기 방지)
3. **stream 닫기 안함**: `CloseAndRecv()` 또는 `CloseSend()` 호출 필요
4. **에러 처리 누락**: `io.EOF`는 정상 종료, 다른 에러는 처리
5. **import 경로**: `go_package` 옵션 정확히 설정

---

추가로 **인터셉터(미들웨어)**, **TLS 설정**, **로드밸런싱**, **health check** 같은 고급 주제도 궁금하시면 말씀해 주세요!