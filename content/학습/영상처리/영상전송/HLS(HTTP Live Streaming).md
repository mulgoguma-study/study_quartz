HLS(HTTP Live Streaming)는 애플(Apple)에서 오픈 표준으로 개발한 **HTTP 기반의 비디오 스트리밍 프로토콜**입니다.

**요약: 전체 영상을 조각내서 쪼개진 미디어파일을 HTTP로 그냥 보내는 구조. CDN 캐싱을 적극 활용하여 대규모 트래픽 확장에 굉장히 유리하고 구현난이도도 그냥 전송만 하는 구조라 낮음.** 

## 1. HLS의 핵심 동작 원리 및 구조

HLS의 본질은 "동영상을 잘게 쪼개서 일반 HTTP 웹 서버로 서빙하는 것"입니다. 크게 두 가지 파일 파일 포맷으로 구성됩니다.

- **`.m3u8` (인덱스/플레이리스트 파일):** 재생할 동영상 조각들의 순서, 경로, 재생 시간(Duration) 등의 메타데이터가 적힌 **텍스트 파일**입니다.
    
- **`.ts` 또는 `.m4s` (미디어 청크 파일):** 실제 영상과 오디오 데이터가 잘게 쪼개져 있는 **바이너리 조각 파일**입니다. (보통 2초~6초 단위로 쪼갭니다.)
    

### 전체 파이프라인 흐름

1. **Ingest (입력):** IP 카메라나 인코더가 실시간 영상 스트림(예: RTSP/RTMP)을 백엔드로 쏩니다.
    
2. **Transcoding & Chunking (변환 및 분할):** 백엔드에서 이 스트림을 받아서 미디어 코덱(주로 H.264 영상 + AAC 오디오)으로 변환하고, 동영상을 3초 단위의 `.ts` 파일들로 쪼갭니다. 그리고 이 목록을 갱신하는 `.m3u8` 파일을 만듭니다.
    
3. **Serving (서빙):** 생성된 `.m3u8`과 `.ts` 파일들을 디스크, 메모리, 혹은 AWS S3 같은 스토리지에 저장하고 일반 HTTP 서버를 통해 클라이언트에게 제공합니다.
    

## 2. Go(Golang)에서 HLS 구현 및 사용법

Go 언어로 HLS 인프라를 구축할 때, 시니어 레벨에서 가장 많이 쓰는 방식은 **FFmpeg 엔진을 백그라운드 워커로 제어하여 파일을 생성하고, Go의 고성능 `net/http` 인프라나 Gin/Echo 같은 프레임워크로 정적 파일을 서빙**하는 아키텍처입니다.

### 실무용 아키텍처 설계 예시 Code

아래 코드는 IP 카메라의 RTSP 스트림을 비동기로 HLS 조각 파일로 변환하고, 이를 HTTP로 고성능 서빙하는 Go 백엔드의 핵심 구조입니다.

Go

```
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	hlsDir     = "./hls_output"
	rtspStream = "rtsp://wowzaec2demo.streamlock.net/vod/mp4:BigBuckBunny_115k.mp4" // 예시 RTSP 스트림 주소
)

func main() {
	// 1. HLS 조각들을 저장할 디렉토리 생성
	if err := os.MkdirAll(hlsDir, os.ModePerm); err != nil {
		log.Fatalf("디렉토리 생성 실패: %v", err)
	}

	// 2. 비동기 워커로 FFmpeg를 실행하여 RTSP -> HLS 실시간 변환 시동
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go startHlsTranscoder(ctx)

	// 3. Go 고성능 HTTP 웹 서버 가동 (Gin 프레임워크 예시)
	r := gin.Default()

	// CORS 설정 (스트리밍 웹 플레이어에서 접근 가능하도록 필수 설정)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Header("Cache-Control", "no-cache") // 인덱스 파일 캐싱 방지
		c.Next()
	})

	// HLS 파일 디렉토리를 정적 파일 서빙 폴더로 매핑
	// 클라이언트는 http://localhost:8080/live/stream.m3u8 로 접근하게 됨
	r.StaticFS("/live", http.Dir(hlsDir))

	log.Println("HLS 스트리밍 서버가 8080 포트에서 시작되었습니다.")
	r.Run(":8080")
}

// FFmpeg를 사용하여 영상을 HLS 규격으로 분할 생성하는 워커
func startHlsTranscoder(ctx context.Context) {
	outputPlaylist := filepath.Join(hlsDir, "stream.m3u8")

	// FFmpeg 옵션 설명:
	// -i: 입력 스트림 (RTSP 등)
	// -c:v libx264: H.264 비디오 코덱 변환 (카메라 원본이 H.264면 copy 가능)
	// -hls_time 3: 동영상 조각(.ts)을 3초 단위로 분할
	// -hls_list_size 5: 라이브 방송용으로 최신 5개 조각만 m3u8에 유지 (VOD용이면 0으로 설정)
	// -hls_flags delete_segments: 라이브 재생 후 지나간 old ts 파일은 자동 삭제하여 디스크 관리
	args := []string{
		"-i", rtspStream,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-f", "hls",
		"-hls_time", "3",
		"-hls_list_size", "5",
		"-hls_flags", "delete_segments",
		outputPlaylist,
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	
	// 실무에서는 대용량 처리 시 FFmpeg 프로세스의 에러 로그를 상시 모니터링해야 함
	cmd.Stderr = os.Stderr 

	log.Println("FFmpeg HLS 변환 프로세스 시작...")
	if err := cmd.Run(); err != nil {
		log.Printf("FFmpeg 프로세스 종료 또는 에러: %v", err)
	}
}
```

## 3. 16년 차 시니어가 면접에서 풀기 좋은 'HLS 대용량 처리' 아키텍처 팁

이 코드를 기반으로 면접관이 "수만 명이 동시에 HLS 영상을 보면 Go 서버가 버틸 수 있겠느냐?"라고 물었을 때 대답할 수 있는 프로페셔널한 대응 전략입니다.

### ① 정적 파일 서빙과 가용성 확보 (오토스케일링 프리)

"HLS의 가장 큰 장점은 아키텍처 유연성입니다. Go 서버 내부에서 정적 파일을 직접 `http.FileServer`로 쏴주는 것은 고루틴 덕분에 가볍지만, 대규모 동시 사용자가 몰릴 때는 Go 서버 앞단에 **AWS CloudFront나 Nginx 같은 CDN/캐시 계층**을 둡니다. `.ts` 미디어 파일은 불변(Immutable)이므로 무한 캐싱이 가능해 백엔드 I/O 부하를 0으로 수렴시킬 수 있습니다."

### ② 파일 디스크 I/O 병목 해결 (In-Memory / Tmpfs)

"실시간 카메라 스트림이 들어올 때 디스크에 `.ts` 파일을 계속 쓰고 지우면 디스크 I/O 병목이 발생하고 SSD 수명이 갉아먹힙니다. 실무에서는 Linux의 **`tmpfs` (메모리 기반 파일 시스템)** 영역에 FFmpeg 출력 디렉토리를 맵핑하여, 메모리 상에서 파일 분할 및 읽기/쓰기가 일어나도록 아키텍처를 최적화합니다."

### ③ 고가용성 라이브러리 선택지 (Pure Go 생태계)

만약 면접관이 내장 OS 커맨드로 FFmpeg를 띄우는 것(`exec.Command`) 말고 고 프로그래밍 내부에서 처리하는 방식을 꼬리 질문으로 던진다면:

"FFmpeg 바이너리 의존성을 줄이고 Pure Go로 가려면, 시니어 개발자들 사이에서 검증된 **`Bluenviron/mediamtx`** (구 rtsp-simple-server, Go로 작성된 초고성능 미디어 라우터 서버)나 **`nareix/joy4`** 같은 라이브러리를 활용해 미디어 파이프라인을 내부 고루틴 컨텍스트 안에서 제어하는 구조로 고도화할 수 있습니다."

