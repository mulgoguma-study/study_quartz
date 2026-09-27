package senior_topics

/*
================================================================================
gRPC와 Protocol Buffers (시니어 레벨)
================================================================================

면접 질문: "gRPC의 장점/단점과 Protocol Buffers의 이점을 설명해주세요"

시니어급 답변 포인트:
1. gRPC 아키텍처와 특징
2. HTTP/2의 이점
3. Protocol Buffers vs JSON
4. gRPC 스트리밍
5. 언제 gRPC를 사용해야 하는가
6. 실제 구현 고려사항
*/

import (
	"context"
	"fmt"
	"io"
	"time"
)

// ============================================================================
// 1. gRPC 개요
// ============================================================================

/*
gRPC (gRPC Remote Procedure Call):
- Google이 개발한 고성능 RPC 프레임워크
- HTTP/2 기반
- Protocol Buffers를 기본 직렬화 형식으로 사용
- 다양한 언어 지원 (Go, Java, Python, C++, etc.)

핵심 특징:
1. IDL (Interface Definition Language) 기반
2. 강력한 타입 시스템
3. 양방향 스트리밍
4. 인터셉터 (미들웨어)
5. 코드 자동 생성
*/

// Protocol Buffers 정의 예시 (.proto 파일)
/*
syntax = "proto3";

package user;

option go_package = "github.com/example/user";

// 메시지 정의
message User {
    string id = 1;
    string name = 2;
    string email = 3;
    int32 age = 4;
    repeated string roles = 5;  // 배열
    Address address = 6;        // 중첩 메시지
}

message Address {
    string street = 1;
    string city = 2;
    string country = 3;
}

message GetUserRequest {
    string id = 1;
}

message GetUserResponse {
    User user = 1;
}

message ListUsersRequest {
    int32 page_size = 1;
    string page_token = 2;
}

message ListUsersResponse {
    repeated User users = 1;
    string next_page_token = 2;
}

// 서비스 정의
service UserService {
    // Unary RPC
    rpc GetUser(GetUserRequest) returns (GetUserResponse);

    // Server Streaming RPC
    rpc ListUsers(ListUsersRequest) returns (stream User);

    // Client Streaming RPC
    rpc CreateUsers(stream User) returns (CreateUsersResponse);

    // Bidirectional Streaming RPC
    rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}
*/

// ============================================================================
// 2. gRPC의 장점
// ============================================================================

/*
gRPC 장점:

1. 성능:
   - HTTP/2: 멀티플렉싱, 헤더 압축
   - Protocol Buffers: 바이너리 직렬화 (JSON보다 작고 빠름)
   - 연결 재사용

2. 강력한 타입 시스템:
   - 컴파일 타임 타입 체크
   - IDE 자동완성
   - 버전 호환성 관리 용이

3. 코드 생성:
   - 클라이언트/서버 스텁 자동 생성
   - 보일러플레이트 코드 제거
   - 여러 언어 동일 인터페이스

4. 스트리밍:
   - 서버 → 클라이언트 스트리밍
   - 클라이언트 → 서버 스트리밍
   - 양방향 스트리밍
   - 실시간 통신에 적합

5. 내장 기능:
   - 인터셉터 (인증, 로깅, 메트릭)
   - 타임아웃/데드라인
   - 취소 전파
   - 로드 밸런싱

6. 생태계:
   - grpc-gateway (REST 변환)
   - grpc-web (브라우저 지원)
   - 다양한 미들웨어
*/

// 인터셉터 예시
type UnaryServerInterceptor func(
	ctx context.Context,
	req interface{},
	info interface{}, // *grpc.UnaryServerInfo
	handler func(ctx context.Context, req interface{}) (interface{}, error),
) (interface{}, error)

// 로깅 인터셉터
func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info interface{},
	handler func(ctx context.Context, req interface{}) (interface{}, error),
) (interface{}, error) {
	start := time.Now()

	// 요청 로깅
	fmt.Printf("gRPC call started: %v\n", info)

	// 핸들러 실행
	resp, err := handler(ctx, req)

	// 응답 로깅
	fmt.Printf("gRPC call completed: duration=%v, error=%v\n",
		time.Since(start), err)

	return resp, err
}

// ============================================================================
// 3. gRPC의 단점
// ============================================================================

/*
gRPC 단점:

1. 브라우저 지원 제한:
   - 네이티브 HTTP/2 지원 불가
   - grpc-web 필요 (추가 프록시)
   - 직접 호출 불가

2. 디버깅 어려움:
   - 바이너리 형식으로 읽기 어려움
   - curl로 테스트 불가
   - 특별한 도구 필요 (grpcurl, Postman)

3. 학습 곡선:
   - Protocol Buffers 문법 학습
   - 빌드 파이프라인 설정
   - 코드 생성 관리

4. 브레이킹 체인지 관리:
   - 필드 번호 변경 불가
   - 타입 변경 주의
   - 호환성 유지 규칙

5. 생태계 성숙도:
   - REST만큼 널리 사용되지 않음
   - 일부 도구/라이브러리 부족
   - 문서화 도구 제한

6. 복잡성:
   - 단순한 API에는 과한 설정
   - 프록시/로드 밸런서 설정 필요
   - HTTP/2 요구
*/

// ============================================================================
// 4. Protocol Buffers vs JSON 비교
// ============================================================================

/*
Protocol Buffers 이점:

1. 크기:
   - JSON보다 3~10배 작음
   - 필드 이름 대신 번호 사용
   - 네트워크 대역폭 절약

2. 속도:
   - 파싱 속도 5~100배 빠름
   - 바이너리 형식 직접 접근
   - GC 부담 감소

3. 스키마:
   - 명확한 계약 정의
   - 버전 관리 용이
   - 타입 안전성

4. 진화:
   - 하위 호환성 유지 용이
   - 필드 추가/삭제 안전
   - 기본값 지원

크기 비교 예시:

JSON (115 bytes):
{
  "id": "user123",
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30,
  "roles": ["admin", "user"]
}

Protocol Buffers (~50 bytes):
바이너리 형식 (필드 번호 + 타입 + 값)

성능 비교:
| 측정 항목        | JSON   | Protobuf |
|-----------------|--------|----------|
| 직렬화 시간      | 1x     | 0.2x     |
| 역직렬화 시간    | 1x     | 0.1x     |
| 메시지 크기      | 1x     | 0.3x     |
| 메모리 할당      | 1x     | 0.3x     |
*/

// JSON 직렬화 예시
type UserJSON struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Age   int32    `json:"age"`
	Roles []string `json:"roles"`
}

// Protocol Buffers 생성 코드 (자동 생성됨)
type UserProto struct {
	Id    string
	Name  string
	Email string
	Age   int32
	Roles []string
}

func (u *UserProto) Marshal() ([]byte, error) {
	// protobuf 직렬화 로직 (자동 생성)
	return nil, nil
}

func (u *UserProto) Unmarshal(data []byte) error {
	// protobuf 역직렬화 로직 (자동 생성)
	return nil
}

// ============================================================================
// 5. gRPC 스트리밍 패턴
// ============================================================================

/*
gRPC 스트리밍 유형:

1. Unary RPC:
   - 클라이언트 단일 요청 → 서버 단일 응답
   - 일반적인 API 호출

2. Server Streaming RPC:
   - 클라이언트 단일 요청 → 서버 다중 응답
   - 대량 데이터 조회, 실시간 업데이트

3. Client Streaming RPC:
   - 클라이언트 다중 요청 → 서버 단일 응답
   - 파일 업로드, 배치 처리

4. Bidirectional Streaming RPC:
   - 양방향 다중 메시지
   - 채팅, 게임, 실시간 협업
*/

// 서버 스트리밍 인터페이스 예시
type ServerStreamingServer interface {
	Send(msg interface{}) error
	Context() context.Context
}

// Server Streaming 구현 예시
func StreamUsersHandler(req interface{}, stream ServerStreamingServer) error {
	// 대량 데이터를 청크로 전송
	for i := 0; i < 100; i++ {
		user := &UserProto{
			Id:   fmt.Sprintf("user%d", i),
			Name: fmt.Sprintf("User %d", i),
		}

		if err := stream.Send(user); err != nil {
			return err
		}

		// 백프레셔 처리
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		default:
		}
	}
	return nil
}

// 클라이언트 스트리밍 인터페이스 예시
type ClientStreamingServer interface {
	Recv() (interface{}, error)
	SendAndClose(interface{}) error
}

// Client Streaming 구현 예시
func UploadUsersHandler(stream ClientStreamingServer) error {
	var count int

	for {
		user, err := stream.Recv()
		if err == io.EOF {
			// 클라이언트가 전송 완료
			return stream.SendAndClose(&struct{ Count int }{Count: count})
		}
		if err != nil {
			return err
		}

		// 사용자 처리
		_ = user
		count++
	}
}

// ============================================================================
// 6. gRPC 사용 시점
// ============================================================================

/*
gRPC를 사용해야 할 때:

1. 마이크로서비스 간 통신:
   - 내부 서비스 간 고성능 통신
   - 타입 안전한 계약
   - 다양한 언어 혼용

2. 스트리밍이 필요한 경우:
   - 실시간 데이터
   - 대량 데이터 전송
   - 양방향 통신

3. 성능이 중요한 경우:
   - 높은 처리량
   - 낮은 지연
   - 네트워크 대역폭 제한

4. 다국어 환경:
   - 여러 언어로 작성된 서비스
   - 일관된 인터페이스 필요

REST를 사용해야 할 때:

1. 공개 API:
   - 브라우저 직접 호출
   - 다양한 클라이언트
   - 간단한 통합

2. 단순한 CRUD:
   - 복잡한 설정 불필요
   - curl로 쉽게 테스트

3. 캐싱 필요:
   - HTTP 캐싱 활용
   - CDN 사용

4. 레거시 시스템:
   - HTTP/1.1만 지원
   - 프록시 제약
*/

// ============================================================================
// 7. gRPC 구현 모범 사례
// ============================================================================

/*
모범 사례:

1. 에러 처리:
   - status.Error 사용
   - 적절한 에러 코드 (codes.NotFound 등)
   - 상세 에러 정보 포함

2. 데드라인/타임아웃:
   - 항상 데드라인 설정
   - 컨텍스트로 취소 전파
   - 적절한 타임아웃 값

3. 인터셉터 활용:
   - 로깅
   - 메트릭
   - 인증/인가
   - 재시도

4. 헬스 체크:
   - grpc_health_v1 구현
   - Kubernetes readiness/liveness

5. 리플렉션:
   - 개발/테스트 환경에서 활성화
   - 프로덕션에서는 비활성화 고려

6. 버전 관리:
   - 패키지명에 버전 포함
   - 하위 호환성 유지
   - Deprecated 필드 명시
*/

// 에러 처리 예시
// import "google.golang.org/grpc/status"
// import "google.golang.org/grpc/codes"

func GetUserWithError(ctx context.Context, userID string) (*UserProto, error) {
	if userID == "" {
		// return nil, status.Error(codes.InvalidArgument, "user_id is required")
		return nil, fmt.Errorf("user_id is required")
	}

	user := findUser(userID)
	if user == nil {
		// return nil, status.Error(codes.NotFound, "user not found")
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func findUser(id string) *UserProto {
	return nil
}

// 데드라인 설정 예시
func CallWithDeadline() {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	//
	// resp, err := client.GetUser(ctx, &GetUserRequest{Id: "123"})
	// if err != nil {
	//     if status.Code(err) == codes.DeadlineExceeded {
	//         // 타임아웃 처리
	//     }
	// }
}

// ============================================================================
// 8. 시니어 면접 답변 요약
// ============================================================================

/*
Q: gRPC의 장점/단점과 Protocol Buffers의 이점은?

A:

gRPC 장점:
1. 고성능: HTTP/2 멀티플렉싱, 바이너리 직렬화
2. 강타입: IDL 기반, 컴파일 타임 체크
3. 코드 생성: 클라이언트/서버 자동 생성
4. 스트리밍: 4가지 패턴 (Unary, Server/Client/Bidirectional)
5. 내장 기능: 인터셉터, 타임아웃, 취소

gRPC 단점:
1. 브라우저 제한: grpc-web 필요
2. 디버깅 어려움: 바이너리 형식
3. 학습 곡선: protobuf, 빌드 파이프라인
4. 복잡성: 단순 API에는 과함

Protocol Buffers 이점:
1. 크기: JSON 대비 3~10배 작음
2. 속도: 파싱 5~100배 빠름
3. 스키마: 명확한 계약, 버전 관리
4. 진화: 하위 호환성 유지 용이

사용 시점:
- gRPC: 마이크로서비스 내부 통신, 스트리밍, 고성능
- REST: 공개 API, 브라우저, 단순 CRUD
*/
