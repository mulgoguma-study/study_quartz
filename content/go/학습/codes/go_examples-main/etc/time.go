package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== 1. 시간 생성 ===")
	createTime()

	fmt.Println("\n=== 2. 현재 시간 가져오기 ===")
	getCurrentTime()

	fmt.Println("\n=== 3. 시간 포맷팅 (문자열 변환) ===")
	formatTime()

	fmt.Println("\n=== 4. 문자열을 시간으로 파싱 ===")
	parseTime()

	fmt.Println("\n=== 5. 시간 계산 (더하기/빼기) ===")
	timeCalculation()

	fmt.Println("\n=== 6. 시간 비교 ===")
	compareTime()

	fmt.Println("\n=== 7. Duration 사용법 ===")
	useDuration()

	fmt.Println("\n=== 8. 타임존(Timezone) 처리 ===")
	handleTimezone()

	fmt.Println("\n=== 9. 시간 컴포넌트 추출 ===")
	extractComponents()

	fmt.Println("\n=== 10. Unix 타임스탬프 ===")
	unixTimestamp()

	fmt.Println("\n=== 11. Timer & Ticker ===")
	timerAndTicker()

	fmt.Println("\n=== 12. 실전 예제 ===")
	practicalExamples()
}

// 1. 시간 생성
func createTime() {
	// 특정 날짜/시간 생성
	t := time.Date(2024, time.March, 15, 14, 30, 0, 0, time.UTC)
	fmt.Println("생성된 시간:", t)

	// 현재 시간
	now := time.Now()
	fmt.Println("현재 시간:", now)

	// Unix 타임스탬프로부터 생성
	unix := time.Unix(1700000000, 0)
	fmt.Println("Unix에서 생성:", unix)
}8

// 2. 현재 시간 가져오기
func getCurrentTime() {
	now := time.Now()
	fmt.Println("전체:", now)
	fmt.Println("로컬:", now.Local())
	fmt.Println("UTC:", now.UTC())
}

// 3. 시간 포맷팅 (문자열 변환)
func formatTime() {
	now := time.Now()

	// Go의 특별한 레퍼런스 시간: "Mon Jan 2 15:04:05 MST 2006"
	// 이 순서를 기억하면 됨: 1월 2일 3시 4분 5초 2006년
	
	fmt.Println("기본 형식:", now.Format(time.RFC3339))
	fmt.Println("커스텀 1:", now.Format("2006-01-02 15:04:05"))
	fmt.Println("커스텀 2:", now.Format("2006년 01월 02일"))
	fmt.Println("시간만:", now.Format("15:04:05"))
	fmt.Println("날짜만:", now.Format("2006-01-02"))
	fmt.Println("12시간:", now.Format("2006-01-02 03:04:05 PM"))
	
	// 자주 사용하는 표준 포맷
	fmt.Println("RFC3339:", now.Format(time.RFC3339))
	fmt.Println("RFC822:", now.Format(time.RFC822))
	fmt.Println("Kitchen:", now.Format(time.Kitchen)) // 3:04PM
}

// 4. 문자열을 시간으로 파싱
func parseTime() {
	// Parse는 UTC 기준
	t1, _ := time.Parse("2006-01-02", "2024-03-15")
	fmt.Println("Parse:", t1)

	// ParseInLocation은 특정 타임존 기준
	loc, _ := time.LoadLocation("Asia/Seoul")
	t2, _ := time.ParseInLocation("2006-01-02 15:04:05", "2024-03-15 14:30:00", loc)
	fmt.Println("ParseInLocation:", t2)

	// 다양한 포맷 파싱 예제
	formats := []string{
		"2006-01-02 15:04:05",
		"2006/01/02",
		"02-Jan-2006",
		time.RFC3339,
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, "2024-03-15 14:30:00"); err == nil {
			fmt.Printf("포맷 %s 파싱 성공: %v\n", format, t)
			break
		}
	}
}

// 5. 시간 계산
func timeCalculation() {
	now := time.Now()

	// 시간 더하기
	after1Hour := now.Add(1 * time.Hour)
	after30Min := now.Add(30 * time.Minute)
	after2Days := now.Add(48 * time.Hour)
	fmt.Println("1시간 후:", after1Hour.Format("15:04:05"))
	fmt.Println("30분 후:", after30Min.Format("15:04:05"))
	fmt.Println("2일 후:", after2Days.Format("2006-01-02"))

	// 날짜 단위 더하기 (AddDate)
	nextMonth := now.AddDate(0, 1, 0)     // 1개월 후
	nextYear := now.AddDate(1, 0, 0)      // 1년 후
	next3Days := now.AddDate(0, 0, 3)     // 3일 후
	fmt.Println("1개월 후:", nextMonth.Format("2006-01-02"))
	fmt.Println("1년 후:", nextYear.Format("2006-01-02"))
	fmt.Println("3일 후:", next3Days.Format("2006-01-02"))

	// 시간 빼기 (음수 사용)
	before1Hour := now.Add(-1 * time.Hour)
	fmt.Println("1시간 전:", before1Hour.Format("15:04:05"))

	// 두 시간의 차이 계산
	future := now.Add(2 * time.Hour)
	diff := future.Sub(now)
	fmt.Println("시간 차이:", diff)
	fmt.Println("차이(시간):", diff.Hours())
	fmt.Println("차이(분):", diff.Minutes())
}

// 6. 시간 비교
func compareTime() {
	t1 := time.Now()
	time.Sleep(10 * time.Millisecond)
	t2 := time.Now()

	// Before, After, Equal
	fmt.Println("t1이 t2보다 이전?", t1.Before(t2))  // true
	fmt.Println("t1이 t2보다 이후?", t1.After(t2))   // false
	fmt.Println("t1과 t2가 같음?", t1.Equal(t2))      // false

	// Compare (Go 1.20+)
	cmp := t1.Compare(t2)
	fmt.Printf("Compare 결과: %d (음수=이전, 0=같음, 양수=이후)\n", cmp)

	// 시간 차이가 특정 범위 내인지 확인
	duration := t2.Sub(t1)
	if duration < 1*time.Second {
		fmt.Println("1초 이내 차이")
	}
}

// 7. Duration 사용법
func useDuration() {
	// Duration 생성
	d1 := 5 * time.Second
	d2 := 300 * time.Millisecond
	d3 := 2*time.Hour + 30*time.Minute

	fmt.Println("5초:", d1)
	fmt.Println("300ms:", d2)
	fmt.Println("2시간 30분:", d3)

	// Duration 변환
	fmt.Println("d3을 시간으로:", d3.Hours())
	fmt.Println("d3을 분으로:", d3.Minutes())
	fmt.Println("d3을 초로:", d3.Seconds())

	// Duration 파싱
	parsed, _ := time.ParseDuration("1h30m45s")
	fmt.Println("파싱된 Duration:", parsed)

	// 실전 예제: 반올림
	now := time.Now()
	rounded := now.Round(time.Minute)
	truncated := now.Truncate(time.Minute)
	fmt.Println("원본:", now.Format("15:04:05.000"))
	fmt.Println("반올림:", rounded.Format("15:04:05.000"))
	fmt.Println("내림:", truncated.Format("15:04:05.000"))
}

// 8. 타임존 처리
func handleTimezone() {
	now := time.Now()

	// 타임존 로드
	seoul, _ := time.LoadLocation("Asia/Seoul")
	tokyo, _ := time.LoadLocation("Asia/Tokyo")
	newYork, _ := time.LoadLocation("America/New_York")

	// 타임존 변환
	fmt.Println("로컬:", now.Format("15:04:05 MST"))
	fmt.Println("서울:", now.In(seoul).Format("15:04:05 MST"))
	fmt.Println("도쿄:", now.In(tokyo).Format("15:04:05 MST"))
	fmt.Println("뉴욕:", now.In(newYork).Format("15:04:05 MST"))

	// 특정 타임존에서 시간 생성
	t := time.Date(2024, 3, 15, 14, 0, 0, 0, seoul)
	fmt.Println("서울 14시:", t)
	fmt.Println("UTC로 변환:", t.UTC())
}

// 9. 시간 컴포넌트 추출
func extractComponents() {
	now := time.Now()

	// 개별 컴포넌트
	year, month, day := now.Date()
	hour, min, sec := now.Clock()

	fmt.Printf("날짜: %d년 %d월 %d일\n", year, month, day)
	fmt.Printf("시간: %d시 %d분 %d초\n", hour, min, sec)

	// 개별 접근
	fmt.Println("연도:", now.Year())
	fmt.Println("월:", now.Month(), int(now.Month()))
	fmt.Println("일:", now.Day())
	fmt.Println("시:", now.Hour())
	fmt.Println("분:", now.Minute())
	fmt.Println("초:", now.Second())
	fmt.Println("나노초:", now.Nanosecond())
	fmt.Println("요일:", now.Weekday())
	fmt.Println("연중 몇째 날:", now.YearDay())

	// 월초, 월말 구하기
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	lastDay := firstDay.AddDate(0, 1, -1)
	fmt.Println("월초:", firstDay.Format("2006-01-02"))
	fmt.Println("월말:", lastDay.Format("2006-01-02"))
}

// 10. Unix 타임스탬프
func unixTimestamp() {
	now := time.Now()

	// Unix 타임스탬프 (초)
	unix := now.Unix()
	fmt.Println("Unix 초:", unix)

	// Unix 타임스탬프 (밀리초)
	unixMilli := now.UnixMilli()
	fmt.Println("Unix 밀리초:", unixMilli)

	// Unix 타임스탬프 (마이크로초)
	unixMicro := now.UnixMicro()
	fmt.Println("Unix 마이크로초:", unixMicro)

	// Unix 타임스탬프 (나노초)
	unixNano := now.UnixNano()
	fmt.Println("Unix 나노초:", unixNano)

	// Unix 타임스탬프에서 시간 생성
	t := time.Unix(unix, 0)
	fmt.Println("Unix에서 복원:", t.Format("2006-01-02 15:04:05"))
}

// 11. Timer & Ticker
func timerAndTicker() {
	fmt.Println("Timer 시작...")
	
	// Timer: 일정 시간 후 한 번 실행
	timer := time.NewTimer(100 * time.Millisecond)
	<-timer.C
	fmt.Println("Timer 완료!")

	// Ticker: 일정 간격으로 반복 실행
	fmt.Println("Ticker 시작 (3번)...")
	ticker := time.NewTicker(50 * time.Millisecond)
	count := 0
	for range ticker.C {
		count++
		fmt.Printf("Tick %d\n", count)
		if count >= 3 {
			ticker.Stop()
			break
		}
	}

	// time.After (간단한 대기)
	fmt.Println("100ms 대기 중...")
	<-time.After(100 * time.Millisecond)
	fmt.Println("대기 완료!")
}

// 12. 실전 예제
func practicalExamples() {
	// 예제 1: 나이 계산
	birthDate := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	age := time.Now().Year() - birthDate.Year()
	fmt.Printf("나이: %d세\n", age)

	// 예제 2: D-Day 계산
	targetDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.Local)
	daysLeft := int(time.Until(targetDate).Hours() / 24)
	fmt.Printf("D-Day: %d일 남음\n", daysLeft)

	// 예제 3: 경과 시간 측정
	start := time.Now()
	time.Sleep(50 * time.Millisecond)
	elapsed := time.Since(start)
	fmt.Printf("경과 시간: %v\n", elapsed)

	// 예제 4: 특정 기간 내 체크
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	if now.After(startOfDay) && now.Before(endOfDay) {
		fmt.Println("오늘 날짜 범위 내")
	}

	// 예제 5: 시간 검증
	timeStr := "2024-13-32" // 잘못된 날짜
	if _, err := time.Parse("2006-01-02", timeStr); err != nil {
		fmt.Println("유효하지 않은 날짜:", err)
	}

	// 예제 6: 다음 정각 구하기
	next := now.Truncate(time.Hour).Add(time.Hour)
	fmt.Println("다음 정각:", next.Format("15:04:05"))

	// 예제 7: 주말 체크
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		fmt.Println("주말입니다!")
	} else {
		fmt.Println("평일입니다.")
	}
}
