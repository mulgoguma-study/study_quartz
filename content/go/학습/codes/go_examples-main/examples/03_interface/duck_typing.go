// =============================================================================
// 덕 타이핑(Duck Typing) 예제
// =============================================================================
// "오리처럼 걷고, 오리처럼 꽥꽥거리면, 그것은 오리다"
//
// 덕 타이핑의 핵심 개념:
// - 객체의 타입보다 행동(메서드)이 중요
// - 명시적인 implements 선언 없이 인터페이스 구현
// - 인터페이스에 정의된 메서드만 구현하면 자동으로 해당 인터페이스를 구현한 것으로 인정
//
// Go의 인터페이스 특징:
// - 암시적 구현 (implicit implementation)
// - 구조적 타이핑 (structural typing)
// - 작은 인터페이스를 조합하여 사용 권장
// =============================================================================

package main

import (
	"fmt"
	"math"
)

// =============================================================================
// 1. 기본 덕 타이핑 예제
// =============================================================================

// Speaker 말할 수 있는 것들의 인터페이스
// 인터페이스는 메서드 시그니처의 집합
type Speaker interface {
	Speak() string
}

// Dog 개 구조체
type Dog struct {
	Name string
}

// Speak Dog의 Speak 메서드 - Speaker 인터페이스 자동 구현
// "implements Speaker"라고 선언하지 않아도 됨!
func (d Dog) Speak() string {
	return fmt.Sprintf("%s: 멍멍!", d.Name)
}

// Cat 고양이 구조체
type Cat struct {
	Name string
}

// Speak Cat의 Speak 메서드
func (c Cat) Speak() string {
	return fmt.Sprintf("%s: 야옹~", c.Name)
}

// Robot 로봇 구조체 (동물이 아니지만 말할 수 있음)
type Robot struct {
	Model string
}

// Speak Robot의 Speak 메서드
func (r Robot) Speak() string {
	return fmt.Sprintf("로봇 %s: 삐빕 삐빕", r.Model)
}

// MakeSpeak Speaker 인터페이스를 받는 함수
// Dog, Cat, Robot 모두 이 함수에 전달 가능
func MakeSpeak(s Speaker) {
	fmt.Println(s.Speak())
}

// =============================================================================
// 2. 실용적인 덕 타이핑 예제 - 도형
// =============================================================================

// Shape 도형 인터페이스
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle 직사각형
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle 원
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Triangle 삼각형
type Triangle struct {
	A, B, C float64 // 세 변의 길이
}

func (t Triangle) Area() float64 {
	// 헤론의 공식
	s := (t.A + t.B + t.C) / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

// PrintShapeInfo Shape 인터페이스를 구현한 어떤 도형이든 처리 가능
func PrintShapeInfo(s Shape) {
	fmt.Printf("면적: %.2f, 둘레: %.2f\n", s.Area(), s.Perimeter())
}

// =============================================================================
// 3. 표준 라이브러리의 덕 타이핑 예제
// =============================================================================

// MyWriter 커스텀 Writer 구현
// io.Writer 인터페이스: Write(p []byte) (n int, err error)
type MyWriter struct {
	data []byte
}

// Write io.Writer 인터페이스 구현
func (w *MyWriter) Write(p []byte) (n int, err error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

// String 저장된 데이터를 문자열로 반환
func (w *MyWriter) String() string {
	return string(w.data)
}

// =============================================================================
// 4. 빈 인터페이스 (any 타입)
// =============================================================================

// PrintAny 어떤 타입이든 출력
// interface{} 또는 any는 모든 타입을 담을 수 있음
func PrintAny(v interface{}) {
	// 타입 스위치를 사용한 타입 확인
	switch val := v.(type) {
	case int:
		fmt.Printf("정수: %d\n", val)
	case string:
		fmt.Printf("문자열: %s\n", val)
	case float64:
		fmt.Printf("실수: %.2f\n", val)
	case bool:
		fmt.Printf("불리언: %v\n", val)
	default:
		fmt.Printf("기타 타입: %T = %v\n", val, val)
	}
}

func main() {
	fmt.Println("=== 덕 타이핑 기본 예제 ===")

	// 서로 다른 타입이지만 모두 Speaker 인터페이스를 구현
	dog := Dog{Name: "바둑이"}
	cat := Cat{Name: "나비"}
	robot := Robot{Model: "T-800"}

	// 모두 MakeSpeak 함수에 전달 가능
	MakeSpeak(dog)
	MakeSpeak(cat)
	MakeSpeak(robot)

	// 슬라이스에도 다양한 타입 저장 가능
	speakers := []Speaker{dog, cat, robot}
	fmt.Println("\n슬라이스에서 반복:")
	for _, s := range speakers {
		MakeSpeak(s)
	}

	fmt.Println("\n=== 도형 예제 ===")

	rect := Rectangle{Width: 10, Height: 5}
	circle := Circle{Radius: 7}
	triangle := Triangle{A: 3, B: 4, C: 5}

	fmt.Print("직사각형 - ")
	PrintShapeInfo(rect)

	fmt.Print("원 - ")
	PrintShapeInfo(circle)

	fmt.Print("삼각형 - ")
	PrintShapeInfo(triangle)

	// 도형 슬라이스
	shapes := []Shape{rect, circle, triangle}
	fmt.Printf("\n도형 %d개의 총 면적: ", len(shapes))

	var totalArea float64
	for _, shape := range shapes {
		totalArea += shape.Area()
	}
	fmt.Printf("%.2f\n", totalArea)

	fmt.Println("\n=== 커스텀 Writer 예제 ===")

	writer := &MyWriter{}

	// io.Writer를 기대하는 함수에 사용 가능
	fmt.Fprint(writer, "Hello, ")
	fmt.Fprint(writer, "Duck Typing!")

	fmt.Printf("기록된 내용: %s\n", writer.String())

	fmt.Println("\n=== 빈 인터페이스 예제 ===")

	PrintAny(42)
	PrintAny("Hello")
	PrintAny(3.14)
	PrintAny(true)
	PrintAny(Dog{Name: "멍멍이"})
}
