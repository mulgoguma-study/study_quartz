# WebRTC & SIP 오디오 처리 완벽 가이드

## 📌 핵심 개념 정리

### 1. 기본 용어 정의

#### Sample (샘플)

- **정의**: 아날로그 음성 신호를 디지털로 변환할 때의 **하나의 측정값**
- **크기**: 보통 16bit (2 bytes) 사용
- **예시**: 1개의 샘플 = 특정 시점의 음압 수치

#### Sample Rate (샘플레이트)

- **정의**: 1초당 샘플링하는 횟수 (Hz 단위)
- **일반적인 값**:
    - `8000 Hz` - 전화 품질 (Narrowband)
    - `16000 Hz` - 광대역 음성 (Wideband)
    - `48000 Hz` - 고품질 오디오 (Fullband)

#### Channel (채널)

- **Mono (1채널)**: 단일 오디오 스트림
- **Stereo (2채널)**: 좌우 분리된 오디오 스트림
- WebRTC/SIP는 주로 **Mono** 사용

#### Frame (프레임)

- **정의**: 일정 시간 동안의 샘플들을 묶은 단위
- **크기**: 보통 10ms, 20ms, 30ms 등
- **용도**: 코덱이 처리하는 기본 단위

#### Chunk (청크)

- **정의**: 실제로 전송/처리되는 데이터 덩어리
- **관계**: 1 chunk = 1 frame 또는 여러 frame의 집합

---

## 🧮 계산 공식

### 기본 공식들

```
1. 샘플 수 계산
Samples = Sample Rate × Duration (초)

2. 바이트 크기 계산
Bytes = Samples × Bit Depth / 8 × Channels

3. 시간 계산
Duration (초) = Samples / Sample Rate

4. Frame당 샘플 수
Samples per Frame = Sample Rate × Frame Duration (초)
```

### 실전 계산 예시

#### 예시 1: 8kHz, 20ms frame

```
- Sample Rate: 8000 Hz
- Frame Duration: 20ms = 0.02초
- Bit Depth: 16bit
- Channels: 1 (Mono)

→ Samples = 8000 × 0.02 = 160 samples
→ Bytes = 160 × 16/8 × 1 = 320 bytes
```

#### 예시 2: 16kHz, 20ms frame

```
- Sample Rate: 16000 Hz
- Frame Duration: 20ms = 0.02초

→ Samples = 16000 × 0.02 = 320 samples
→ Bytes = 320 × 2 × 1 = 640 bytes
```

#### 예시 3: 48kHz, 10ms frame

```
- Sample Rate: 48000 Hz
- Frame Duration: 10ms = 0.01초

→ Samples = 48000 × 0.01 = 480 samples
→ Bytes = 480 × 2 × 1 = 960 bytes
```

---

## 📊 Quick Reference Table

### Sample Rate별 Frame 크기 (16bit Mono 기준)

|Sample Rate|Frame Duration|Samples|Bytes|비고|
|---|---|---|---|---|
|8kHz|10ms|80|160|전화 품질|
|8kHz|20ms|160|320|✅ 가장 일반적|
|8kHz|30ms|240|480||
|16kHz|10ms|160|320||
|16kHz|20ms|320|640|✅ Wideband 표준|
|16kHz|30ms|480|960||
|48kHz|10ms|480|960||
|48kHz|20ms|960|1920|✅ WebRTC 기본값|
|48kHz|30ms|1440|2880||

### 초당 데이터 전송량 (Bitrate)

|Sample Rate|Bit Depth|Channels|Bitrate (kbps)|
|---|---|---|---|
|8kHz|16bit|1|128 kbps|
|16kHz|16bit|1|256 kbps|
|48kHz|16bit|1|768 kbps|
|48kHz|16bit|2|1536 kbps|

**계산식**: `Bitrate = Sample Rate × Bit Depth × Channels`

---

## 🎯 실전 시나리오별 가이드

### WebRTC 기본 설정

```javascript
// WebRTC AudioContext 기본값
const audioContext = new AudioContext({
  sampleRate: 48000  // 48kHz 고정
});

// ScriptProcessorNode (Deprecated)
const bufferSize = 4096;  // 샘플 수 (약 85ms @ 48kHz)

// AudioWorklet (권장)
// 128 samples 단위로 처리 (약 2.67ms @ 48kHz)
```

**중요**: WebRTC는 내부적으로 **48kHz**를 사용하고, 코덱에서 8kHz/16kHz로 다운샘플링

### SIP/VoIP 일반 설정

```
- G.711 (PCMU/PCMA): 8kHz, 20ms frame
  → 160 samples, 160 bytes (8bit)
  
- G.722: 16kHz, 20ms frame
  → 320 samples, 코덱 압축 후 ~160 bytes
  
- Opus: 8~48kHz, 20ms frame (기본)
  → 가변 bitrate, 최대 320 samples @ 16kHz
```

### RTP Packet 구조

```
RTP Header: 12 bytes
+ Payload: Frame data
= Total Packet Size

예) G.711, 20ms frame
→ 12 + 160 = 172 bytes per packet
→ 50 packets/sec (1000ms / 20ms)
```

---

## 💡 자주 하는 실수와 해결법

### 실수 1: Sample Rate 불일치

```javascript
// ❌ 잘못된 예
const audioContext = new AudioContext({ sampleRate: 8000 });
// → 브라우저는 48kHz로 강제 설정됨

// ✅ 올바른 예
const audioContext = new AudioContext();  // 48kHz
// → 필요시 Resampling 라이브러리 사용
```

### 실수 2: Buffer Size와 Latency 혼동

```
Buffer Size (samples) ≠ Network Latency

- Buffer Size: 로컬 처리 단위 (2.67~85ms)
- Frame Duration: 네트워크 전송 단위 (10~30ms)
- Network Latency: 실제 지연 시간 (수십~수백ms)
```

### 실수 3: Byte/Sample 계산 오류

```javascript
// ❌ 잘못된 계산
const bytes = samples;  // 16bit를 1byte로 착각

// ✅ 올바른 계산
const bytes = samples * 2;  // 16bit = 2 bytes
```

---

## 🔧 유용한 유틸리티 함수

```javascript
/**
 * 오디오 파라미터 계산기
 */
class AudioCalculator {
  /**
   * 샘플 수 계산
   * @param {number} sampleRate - Hz
   * @param {number} durationMs - milliseconds
   */
  static getSamples(sampleRate, durationMs) {
    return Math.floor(sampleRate * durationMs / 1000);
  }
  
  /**
   * 바이트 크기 계산
   * @param {number} samples - 샘플 수
   * @param {number} bitDepth - 비트 깊이 (기본 16)
   * @param {number} channels - 채널 수 (기본 1)
   */
  static getBytes(samples, bitDepth = 16, channels = 1) {
    return samples * (bitDepth / 8) * channels;
  }
  
  /**
   * 시간 계산
   * @param {number} samples - 샘플 수
   * @param {number} sampleRate - Hz
   * @returns {number} milliseconds
   */
  static getDuration(samples, sampleRate) {
    return (samples / sampleRate) * 1000;
  }
  
  /**
   * Bitrate 계산
   * @param {number} sampleRate - Hz
   * @param {number} bitDepth - 비트 깊이
   * @param {number} channels - 채널 수
   * @returns {number} bps (bits per second)
   */
  static getBitrate(sampleRate, bitDepth, channels) {
    return sampleRate * bitDepth * channels;
  }
}

// 사용 예시
const samples = AudioCalculator.getSamples(16000, 20);  // 320
const bytes = AudioCalculator.getBytes(samples);         // 640
const duration = AudioCalculator.getDuration(480, 48000); // 10ms
const bitrate = AudioCalculator.getBitrate(48000, 16, 2); // 1536000 bps
```

---

## 📝 체크리스트

개발 시작 전 확인사항:

- [ ] 목표 Sample Rate 확인 (8/16/48 kHz?)
- [ ] Frame Duration 결정 (10/20/30 ms?)
- [ ] 코덱 선택 (G.711/G.722/Opus?)
- [ ] Buffer Size 설정 (지연시간 고려)
- [ ] Resampling 필요 여부 확인
- [ ] RTP Packet Size 계산 (MTU 초과 주의)

디버깅 시 확인사항:

- [ ] 실제 Sample Rate vs 설정값 일치?
- [ ] Byte 크기 계산 맞나? (16bit = 2bytes!)
- [ ] Frame 경계 정렬 확인
- [ ] Under-run / Over-run 발생?
- [ ] Jitter Buffer 크기 적절?

---

## 🔗 추가 참고 자료

- **WebRTC**: `getUserMedia()` → 48kHz 고정
- **Opus Codec**: 8, 12, 16, 24, 48 kHz 지원
- **G.711**: 8kHz 고정, 64 kbps
- **G.722**: 16kHz, 48/56/64 kbps

**핵심 기억**:

> Sample Rate × Duration = Samples  
> Samples × 2 bytes = Raw PCM Size