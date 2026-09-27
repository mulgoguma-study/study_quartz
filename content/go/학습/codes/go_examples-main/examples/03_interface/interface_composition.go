// =============================================================================
// 인터페이스 컴포지션(조합) 및 상속 예제
// =============================================================================
// Go는 클래스 기반 상속이 없지만, 인터페이스 조합을 통해 유사한 효과를 얻습니다.
//
// 인터페이스 조합(Embedding):
// - 작은 인터페이스들을 조합하여 큰 인터페이스 생성
// - 코드 재사용성과 유연성 증가
// - Go의 표준 라이브러리에서 널리 사용되는 패턴
//
// 구조체 임베딩:
// - 구조체 안에 다른 구조체/인터페이스를 임베딩
// - 임베딩된 타입의 메서드를 외부에서 직접 호출 가능
// - 상속이 아닌 위임(delegation)을 통한 코드 재사용
// =============================================================================

package main

import (
	"fmt"
	"io"
	"strings"
)

// =============================================================================
// 1. 기본 인터페이스 조합
// =============================================================================

// Reader 읽기 기능 인터페이스
type Reader interface {
	Read() string
}

// Writer 쓰기 기능 인터페이스
type Writer interface {
	Write(data string)
}

// Closer 닫기 기능 인터페이스
type Closer interface {
	Close() error
}

// ReadWriter Reader와 Writer를 조합한 인터페이스
// 두 인터페이스의 모든 메서드를 포함
type ReadWriter interface {
	Reader
	Writer
}

// ReadWriteCloser 세 개의 인터페이스를 조합
type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}

// File 파일 구조체 - ReadWriteCloser 인터페이스 구현
type File struct {
	name    string
	content string
	isOpen  bool
}

func NewFile(name string) *File {
	return &File{name: name, isOpen: true}
}

func (f *File) Read() string {
	if !f.isOpen {
		return "에러: 파일이 닫혀있습니다"
	}
	return f.content
}

func (f *File) Write(data string) {
	if !f.isOpen {
		fmt.Println("에러: 파일이 닫혀있습니다")
		return
	}
	f.content += data
}

func (f *File) Close() error {
	if !f.isOpen {
		return fmt.Errorf("파일이 이미 닫혀있습니다")
	}
	f.isOpen = false
	fmt.Printf("파일 '%s' 닫힘\n", f.name)
	return nil
}

// =============================================================================
// 2. 구조체 임베딩
// =============================================================================

// Animal 기본 동물 구조체
type Animal struct {
	Name string
	Age  int
}

func (a Animal) Eat() {
	fmt.Printf("%s이(가) 먹이를 먹습니다\n", a.Name)
}

func (a Animal) Sleep() {
	fmt.Printf("%s이(가) 잠을 잡니다\n", a.Name)
}

// Dog2 Animal을 임베딩한 개 구조체
type Dog2 struct {
	Animal // 임베딩 - 익명 필드
	Breed  string
}

// Bark Dog2만의 메서드
func (d Dog2) Bark() {
	fmt.Printf("%s: 멍멍!\n", d.Name)
}

// Eat 오버라이딩 - Animal의 Eat 메서드를 재정의
func (d Dog2) Eat() {
	fmt.Printf("%s이(가) 개 사료를 맛있게 먹습니다\n", d.Name)
}

// Bird Animal을 임베딩한 새 구조체
type Bird struct {
	Animal
	CanFly bool
}

func (b Bird) Fly() {
	if b.CanFly {
		fmt.Printf("%s이(가) 하늘을 날아갑니다\n", b.Name)
	} else {
		fmt.Printf("%s은(는) 날 수 없습니다\n", b.Name)
	}
}

// =============================================================================
// 3. 다중 인터페이스 구현
// =============================================================================

// Swimmer 수영 인터페이스
type Swimmer interface {
	Swim() string
}

// Flyer 비행 인터페이스
type Flyer interface {
	Fly() string
}

// Walker 걷기 인터페이스
type Walker interface {
	Walk() string
}

// Duck2 오리는 세 가지 인터페이스 모두 구현
type Duck2 struct {
	Name string
}

func (d Duck2) Swim() string {
	return fmt.Sprintf("%s이(가) 수영합니다", d.Name)
}

func (d Duck2) Fly() string {
	return fmt.Sprintf("%s이(가) 날아갑니다", d.Name)
}

func (d Duck2) Walk() string {
	return fmt.Sprintf("%s이(가) 걸어갑니다", d.Name)
}

// Penguin 펭귄은 수영과 걷기만 가능
type Penguin struct {
	Name string
}

func (p Penguin) Swim() string {
	return fmt.Sprintf("%s이(가) 헤엄칩니다", p.Name)
}

func (p Penguin) Walk() string {
	return fmt.Sprintf("%s이(가) 뒤뚱뒤뚱 걷습니다", p.Name)
}

// =============================================================================
// 4. 실전 예제: 로깅 시스템
// =============================================================================

// Logger 기본 로거 인터페이스
type Logger interface {
	Log(message string)
}

// Formatter 포맷터 인터페이스
type Formatter interface {
	Format(message string) string
}

// FormattedLogger Logger와 Formatter 조합
type FormattedLogger interface {
	Logger
	Formatter
}

// SimpleLogger 기본 로거 구현
type SimpleLogger struct {
	prefix string
}

func (l SimpleLogger) Log(message string) {
	fmt.Printf("[%s] %s\n", l.prefix, message)
}

// PrefixFormatter 접두사 포맷터
type PrefixFormatter struct {
	Prefix string
}

func (f PrefixFormatter) Format(message string) string {
	return fmt.Sprintf("%s: %s", f.Prefix, message)
}

// EnhancedLogger 여러 기능을 임베딩한 로거
type EnhancedLogger struct {
	SimpleLogger    // 로깅 기능 임베딩
	PrefixFormatter // 포맷팅 기능 임베딩
	LogLevel        string
}

// Log 오버라이딩 - 포맷팅 적용
func (l EnhancedLogger) Log(message string) {
	formatted := l.Format(message)
	fmt.Printf("[%s][%s] %s\n", l.LogLevel, l.prefix, formatted)
}

// =============================================================================
// 5. 표준 라이브러리 인터페이스 조합 예시
// =============================================================================

// MyBuffer io.ReadWriter 구현
type MyBuffer struct {
	buffer strings.Builder
}

func (b *MyBuffer) Read(p []byte) (n int, err error) {
	s := b.buffer.String()
	n = copy(p, s)
	if n < len(s) {
		return n, nil
	}
	return n, io.EOF
}

func (b *MyBuffer) Write(p []byte) (n int, err error) {
	return b.buffer.Write(p)
}

func main() {
	fmt.Println("=== 1. 인터페이스 조합 예제 ===")

	file := NewFile("test.txt")
	file.Write("Hello, ")
	file.Write("Interface Composition!")
	fmt.Printf("파일 내용: %s\n", file.Read())
	file.Close()

	// ReadWriter로 사용
	var rw ReadWriter = NewFile("rw.txt")
	rw.Write("Read and Write only")
	fmt.Printf("ReadWriter: %s\n", rw.Read())

	fmt.Println("\n=== 2. 구조체 임베딩 예제 ===")

	dog := Dog2{
		Animal: Animal{Name: "바둑이", Age: 3},
		Breed:  "진돗개",
	}

	// 임베딩된 Animal의 메서드 직접 호출
	dog.Sleep()        // Animal.Sleep() 호출
	dog.Eat()          // Dog2.Eat() 호출 (오버라이딩됨)
	dog.Bark()         // Dog2.Bark() 호출
	dog.Animal.Eat()   // 원래 Animal.Eat() 명시적 호출

	fmt.Printf("품종: %s, 나이: %d세\n", dog.Breed, dog.Age)

	bird := Bird{
		Animal: Animal{Name: "참새", Age: 1},
		CanFly: true,
	}
	bird.Eat()
	bird.Fly()

	fmt.Println("\n=== 3. 다중 인터페이스 구현 예제 ===")

	duck := Duck2{Name: "도널드"}
	penguin := Penguin{Name: "뽀로로"}

	// 오리는 모든 인터페이스 사용 가능
	fmt.Println(duck.Swim())
	fmt.Println(duck.Fly())
	fmt.Println(duck.Walk())

	// 펭귄은 Swimmer, Walker로만 사용 가능
	fmt.Println(penguin.Swim())
	fmt.Println(penguin.Walk())

	// 인터페이스 타입으로 사용
	var swimmer Swimmer = penguin
	fmt.Printf("Swimmer 인터페이스로: %s\n", swimmer.Swim())

	fmt.Println("\n=== 4. 로깅 시스템 예제 ===")

	enhancedLogger := EnhancedLogger{
		SimpleLogger:    SimpleLogger{prefix: "APP"},
		PrefixFormatter: PrefixFormatter{Prefix: "DEBUG"},
		LogLevel:        "INFO",
	}

	enhancedLogger.Log("애플리케이션이 시작되었습니다")
	enhancedLogger.Log("사용자가 로그인했습니다")

	fmt.Println("\n=== 5. 표준 라이브러리 인터페이스 조합 ===")

	buf := &MyBuffer{}

	// io.Writer로 사용
	fmt.Fprint(buf, "Hello from io.Writer!")

	// io.Reader로 사용
	data := make([]byte, 100)
	n, _ := buf.Read(data)
	fmt.Printf("읽은 데이터: %s\n", string(data[:n]))
}
