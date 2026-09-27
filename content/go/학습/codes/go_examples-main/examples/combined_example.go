// =============================================================================
// 통합 예제: 주문 처리 시스템
// =============================================================================
// 이 예제는 다음 개념들을 통합하여 사용합니다:
// - GORM 트랜잭션 (비관적 락, ACID)
// - Goroutine (WaitGroup, Channel, Worker Pool, Context)
// - 인터페이스 (덕 타이핑, 컴포지션)
// - 에러 처리 (errors.Is, errors.As)
//
// 시나리오:
// - 여러 주문을 동시에 처리하는 주문 처리 시스템
// - 재고 확인, 결제 처리, 배송 예약을 워커 풀로 처리
// - 트랜잭션으로 데이터 일관성 보장
// =============================================================================

package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// =============================================================================
// 에러 정의
// =============================================================================

var (
	ErrInsufficientStock = errors.New("재고 부족")
	ErrPaymentFailed     = errors.New("결제 실패")
	ErrOrderNotFound     = errors.New("주문을 찾을 수 없음")
)

// OrderError 주문 관련 커스텀 에러
type OrderError struct {
	OrderID   uint
	Step      string
	Err       error
	Timestamp time.Time
}

func (e *OrderError) Error() string {
	return fmt.Sprintf("[주문 %d] %s 단계 실패: %v (시간: %s)",
		e.OrderID, e.Step, e.Err, e.Timestamp.Format("15:04:05"))
}

func (e *OrderError) Unwrap() error {
	return e.Err
}

// =============================================================================
// 모델 정의
// =============================================================================

// Product 상품 모델
type Product struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:200;not null"`
	Price     int64  `gorm:"not null"`
	Stock     int    `gorm:"not null;default:0"`
	Version   int    `gorm:"not null;default:1"` // 낙관적 락
	UpdatedAt time.Time
}

// Order 주문 모델
type Order struct {
	ID         uint   `gorm:"primaryKey"`
	ProductID  uint   `gorm:"not null;index"`
	Quantity   int    `gorm:"not null"`
	TotalPrice int64  `gorm:"not null"`
	Status     string `gorm:"size:20;not null;default:'PENDING'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// =============================================================================
// 인터페이스 정의 (덕 타이핑)
// =============================================================================

// OrderProcessor 주문 처리 인터페이스
type OrderProcessor interface {
	Process(ctx context.Context, order *Order) error
}

// PaymentGateway 결제 게이트웨이 인터페이스
type PaymentGateway interface {
	Charge(amount int64) error
	Refund(amount int64) error
}

// NotificationService 알림 서비스 인터페이스
type NotificationService interface {
	Send(message string) error
}

// =============================================================================
// 인터페이스 구현 (덕 타이핑)
// =============================================================================

// StripePayment Stripe 결제 구현
type StripePayment struct{}

func (s *StripePayment) Charge(amount int64) error {
	// 20% 확률로 결제 실패 시뮬레이션
	if rand.Float32() < 0.2 {
		return ErrPaymentFailed
	}
	fmt.Printf("  💳 Stripe 결제 완료: %d원\n", amount)
	return nil
}

func (s *StripePayment) Refund(amount int64) error {
	fmt.Printf("  💳 Stripe 환불 완료: %d원\n", amount)
	return nil
}

// EmailNotification 이메일 알림 구현
type EmailNotification struct{}

func (e *EmailNotification) Send(message string) error {
	fmt.Printf("  📧 이메일 발송: %s\n", message)
	return nil
}

// =============================================================================
// 주문 처리 서비스 (인터페이스 컴포지션)
// =============================================================================

// OrderService 주문 서비스
type OrderService struct {
	db           *gorm.DB
	payment      PaymentGateway       // 인터페이스로 주입
	notification NotificationService  // 인터페이스로 주입
}

// NewOrderService 생성자
func NewOrderService(db *gorm.DB, payment PaymentGateway, notification NotificationService) *OrderService {
	return &OrderService{
		db:           db,
		payment:      payment,
		notification: notification,
	}
}

// ProcessOrder 주문 처리 (ACID 트랜잭션)
func (s *OrderService) ProcessOrder(ctx context.Context, order *Order) error {
	// 컨텍스트 취소 확인
	select {
	case <-ctx.Done():
		return &OrderError{
			OrderID:   order.ID,
			Step:      "시작",
			Err:       ctx.Err(),
			Timestamp: time.Now(),
		}
	default:
	}

	// 트랜잭션으로 ACID 보장
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1단계: 상품 조회 (비관적 락)
		var product Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&product, order.ProductID).Error; err != nil {
			return &OrderError{
				OrderID:   order.ID,
				Step:      "상품조회",
				Err:       err,
				Timestamp: time.Now(),
			}
		}

		// 2단계: 재고 확인
		if product.Stock < order.Quantity {
			return &OrderError{
				OrderID:   order.ID,
				Step:      "재고확인",
				Err:       ErrInsufficientStock,
				Timestamp: time.Now(),
			}
		}

		// 3단계: 총 가격 계산
		order.TotalPrice = product.Price * int64(order.Quantity)

		// 4단계: 결제 처리
		if err := s.payment.Charge(order.TotalPrice); err != nil {
			return &OrderError{
				OrderID:   order.ID,
				Step:      "결제",
				Err:       err,
				Timestamp: time.Now(),
			}
		}

		// 5단계: 재고 차감
		product.Stock -= order.Quantity
		if err := tx.Save(&product).Error; err != nil {
			// 결제 환불
			s.payment.Refund(order.TotalPrice)
			return &OrderError{
				OrderID:   order.ID,
				Step:      "재고차감",
				Err:       err,
				Timestamp: time.Now(),
			}
		}

		// 6단계: 주문 상태 업데이트
		order.Status = "COMPLETED"
		if err := tx.Save(order).Error; err != nil {
			s.payment.Refund(order.TotalPrice)
			return &OrderError{
				OrderID:   order.ID,
				Step:      "주문업데이트",
				Err:       err,
				Timestamp: time.Now(),
			}
		}

		return nil
	})
}

// =============================================================================
// 워커 풀 (Goroutine)
// =============================================================================

// OrderWorkerPool 주문 처리 워커 풀
type OrderWorkerPool struct {
	numWorkers   int
	orderService *OrderService
	orders       chan *Order
	results      chan OrderResult
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
}

// OrderResult 주문 처리 결과
type OrderResult struct {
	Order    *Order
	Success  bool
	Error    error
	WorkerID int
	Duration time.Duration
}

// NewOrderWorkerPool 워커 풀 생성
func NewOrderWorkerPool(numWorkers int, orderService *OrderService) *OrderWorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &OrderWorkerPool{
		numWorkers:   numWorkers,
		orderService: orderService,
		orders:       make(chan *Order, 100),
		results:      make(chan OrderResult, 100),
		ctx:          ctx,
		cancel:       cancel,
	}

	pool.start()
	return pool
}

// start 워커 시작
func (p *OrderWorkerPool) start() {
	for i := 1; i <= p.numWorkers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

// worker 워커 고루틴
func (p *OrderWorkerPool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			fmt.Printf("Worker %d: 종료\n", id)
			return

		case order, ok := <-p.orders:
			if !ok {
				return
			}

			startTime := time.Now()

			// 주문 처리 (타임아웃 컨텍스트)
			ctx, cancel := context.WithTimeout(p.ctx, 5*time.Second)
			err := p.orderService.ProcessOrder(ctx, order)
			cancel()

			p.results <- OrderResult{
				Order:    order,
				Success:  err == nil,
				Error:    err,
				WorkerID: id,
				Duration: time.Since(startTime),
			}
		}
	}
}

// Submit 주문 제출
func (p *OrderWorkerPool) Submit(order *Order) {
	p.orders <- order
}

// Results 결과 채널 반환
func (p *OrderWorkerPool) Results() <-chan OrderResult {
	return p.results
}

// Shutdown 워커 풀 종료
func (p *OrderWorkerPool) Shutdown() {
	close(p.orders)
	p.wg.Wait()
	close(p.results)
}

// =============================================================================
// 메인 함수
// =============================================================================

func main() {
	rand.Seed(time.Now().UnixNano())

	// ==========================================================================
	// 1. 데이터베이스 설정
	// ==========================================================================
	db, err := gorm.Open(sqlite.Open("combined_example.db"), &gorm.Config{})
	if err != nil {
		panic("DB 연결 실패: " + err.Error())
	}

	// 테이블 마이그레이션
	db.AutoMigrate(&Product{}, &Order{})

	// 테스트 데이터 초기화
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM orders")

	// 상품 생성
	products := []Product{
		{Name: "맥북 프로", Price: 3000000, Stock: 5},
		{Name: "아이폰 15", Price: 1500000, Stock: 10},
		{Name: "에어팟", Price: 300000, Stock: 20},
	}

	for i := range products {
		db.Create(&products[i])
	}

	fmt.Println("=== 초기 상품 목록 ===")
	var allProducts []Product
	db.Find(&allProducts)
	for _, p := range allProducts {
		fmt.Printf("  %s: %d원 (재고: %d개)\n", p.Name, p.Price, p.Stock)
	}

	// ==========================================================================
	// 2. 서비스 생성 (인터페이스 주입)
	// ==========================================================================
	orderService := NewOrderService(
		db,
		&StripePayment{},      // PaymentGateway 인터페이스
		&EmailNotification{},   // NotificationService 인터페이스
	)

	// ==========================================================================
	// 3. 워커 풀 생성 및 주문 처리
	// ==========================================================================
	fmt.Println("\n=== 주문 처리 시작 ===")

	pool := NewOrderWorkerPool(3, orderService)

	// 주문 생성 및 제출
	orders := []*Order{
		{ProductID: 1, Quantity: 2},  // 맥북 2개
		{ProductID: 2, Quantity: 3},  // 아이폰 3개
		{ProductID: 3, Quantity: 5},  // 에어팟 5개
		{ProductID: 1, Quantity: 10}, // 맥북 10개 (재고 부족 예상)
		{ProductID: 2, Quantity: 1},  // 아이폰 1개
	}

	// 주문 저장 및 제출
	for _, order := range orders {
		db.Create(order)
		pool.Submit(order)
	}

	// ==========================================================================
	// 4. 결과 수집 (Fan-In)
	// ==========================================================================
	var (
		successCount int
		failCount    int
	)

	// 결과 수집
	go func() {
		time.Sleep(3 * time.Second)
		pool.Shutdown()
	}()

	for result := range pool.Results() {
		if result.Success {
			successCount++
			fmt.Printf("✅ 주문 %d 성공 (Worker %d, 소요: %v)\n",
				result.Order.ID, result.WorkerID, result.Duration)
		} else {
			failCount++
			fmt.Printf("❌ 주문 %d 실패 (Worker %d, 소요: %v)\n",
				result.Order.ID, result.WorkerID, result.Duration)

			// errors.As를 사용한 에러 처리
			var orderErr *OrderError
			if errors.As(result.Error, &orderErr) {
				fmt.Printf("   → 단계: %s\n", orderErr.Step)

				// errors.Is를 사용한 에러 확인
				switch {
				case errors.Is(result.Error, ErrInsufficientStock):
					fmt.Println("   → 원인: 재고 부족")
				case errors.Is(result.Error, ErrPaymentFailed):
					fmt.Println("   → 원인: 결제 실패")
				default:
					fmt.Printf("   → 원인: %v\n", orderErr.Err)
				}
			}
		}
	}

	// ==========================================================================
	// 5. 최종 결과 출력
	// ==========================================================================
	fmt.Println("\n=== 처리 결과 요약 ===")
	fmt.Printf("성공: %d건, 실패: %d건\n", successCount, failCount)

	fmt.Println("\n=== 최종 상품 재고 ===")
	db.Find(&allProducts)
	for _, p := range allProducts {
		fmt.Printf("  %s: 재고 %d개\n", p.Name, p.Stock)
	}

	fmt.Println("\n=== 주문 상태 ===")
	var allOrders []Order
	db.Find(&allOrders)
	for _, o := range allOrders {
		fmt.Printf("  주문 %d: 상품ID=%d, 수량=%d, 상태=%s\n",
			o.ID, o.ProductID, o.Quantity, o.Status)
	}
}
