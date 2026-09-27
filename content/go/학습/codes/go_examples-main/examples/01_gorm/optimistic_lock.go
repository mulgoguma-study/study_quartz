// =============================================================================
// 낙관적 락(Optimistic Lock) 예제
// =============================================================================
// 낙관적 락은 충돌이 적게 발생할 것이라고 "낙관적"으로 가정합니다.
// 데이터를 읽을 때는 락을 걸지 않고, 업데이트할 때 버전을 확인하여 충돌을 감지합니다.
//
// 주요 특징:
// - 버전(Version) 필드를 사용하여 충돌 감지
// - 읽을 때 락을 걸지 않아 성능이 좋음
// - 충돌 시 재시도 로직 필요
// - 읽기가 많고 쓰기가 적은 경우에 적합
// =============================================================================

package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Product 상품 모델 (낙관적 락 사용)
// Version 필드가 핵심입니다
type Product struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:200;not null"`
	Stock     int    `gorm:"not null;default:0"` // 재고 수량
	Version   int    `gorm:"not null;default:1"` // 버전 필드 (낙관적 락의 핵심)
	UpdatedAt time.Time
}

// ErrOptimisticLock 낙관적 락 충돌 에러
// 다른 트랜잭션이 먼저 데이터를 수정했을 때 발생
var ErrOptimisticLock = errors.New("낙관적 락 충돌: 다른 트랜잭션이 먼저 수정함")

func main() {
	// ==========================================================================
	// 1. 데이터베이스 연결
	// ==========================================================================
	db, err := gorm.Open(sqlite.Open("optimistic_lock.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("데이터베이스 연결 실패:", err)
	}

	// 테이블 마이그레이션
	if err := db.AutoMigrate(&Product{}); err != nil {
		log.Fatal("마이그레이션 실패:", err)
	}

	// ==========================================================================
	// 2. 테스트 데이터 초기화
	// ==========================================================================
	db.Exec("DELETE FROM products")

	product := Product{
		Name:    "맥북 프로",
		Stock:   10, // 초기 재고 10개
		Version: 1,  // 초기 버전
	}
	db.Create(&product)
	fmt.Printf("상품 등록: ID=%d, 이름=%s, 재고=%d, 버전=%d\n\n",
		product.ID, product.Name, product.Stock, product.Version)

	// ==========================================================================
	// 3. 낙관적 락을 사용한 동시 재고 차감 테스트
	// ==========================================================================
	var wg sync.WaitGroup

	// 5개의 고루틴이 각각 3개씩 구매 시도
	// 재고 10개이므로 일부는 실패해야 함
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			// 재시도 로직 포함
			purchaseWithOptimisticLock(db, product.ID, 3, workerID, 3)
		}(i)
	}

	wg.Wait()

	// ==========================================================================
	// 4. 최종 결과 확인
	// ==========================================================================
	var finalProduct Product
	db.First(&finalProduct, product.ID)
	fmt.Printf("\n최종 재고: %d개, 버전: %d\n", finalProduct.Stock, finalProduct.Version)
}

// purchaseWithOptimisticLock 낙관적 락을 사용한 구매 함수
//
// 매개변수:
//   - db: GORM 데이터베이스 연결
//   - productID: 상품 ID
//   - quantity: 구매 수량
//   - workerID: 작업자 식별 번호
//   - maxRetries: 최대 재시도 횟수
func purchaseWithOptimisticLock(db *gorm.DB, productID uint, quantity int, workerID int, maxRetries int) {
	fmt.Printf("[Worker %d] 구매 시도: %d개\n", workerID, quantity)

	// ==========================================================================
	// 재시도 루프
	// ==========================================================================
	// 낙관적 락에서는 충돌 시 재시도가 필수입니다
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := attemptPurchase(db, productID, quantity, workerID)

		if err == nil {
			// 성공
			return
		}

		if errors.Is(err, ErrOptimisticLock) {
			// 낙관적 락 충돌 - 재시도
			fmt.Printf("[Worker %d] 충돌 감지, 재시도 %d/%d\n", workerID, attempt, maxRetries)
			time.Sleep(50 * time.Millisecond) // 잠시 대기 후 재시도
			continue
		}

		// 다른 에러 (재고 부족 등) - 재시도하지 않음
		fmt.Printf("[Worker %d] 구매 실패: %v\n", workerID, err)
		return
	}

	fmt.Printf("[Worker %d] 최대 재시도 횟수 초과, 구매 실패\n", workerID)
}

// attemptPurchase 구매 시도 (단일 시도)
func attemptPurchase(db *gorm.DB, productID uint, quantity int, workerID int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// ======================================================================
		// 1단계: 현재 데이터 읽기 (락 없음)
		// ======================================================================
		// 일반 SELECT - 락을 걸지 않습니다
		var product Product
		if err := tx.First(&product, productID).Error; err != nil {
			return err
		}

		// 현재 버전 저장 (나중에 비교용)
		currentVersion := product.Version

		fmt.Printf("[Worker %d] 현재 재고: %d, 버전: %d\n",
			workerID, product.Stock, currentVersion)

		// ======================================================================
		// 2단계: 비즈니스 로직 검증
		// ======================================================================
		if product.Stock < quantity {
			return fmt.Errorf("재고 부족 (현재: %d, 요청: %d)", product.Stock, quantity)
		}

		// 처리 시간 시뮬레이션 (동시성 테스트를 위해)
		time.Sleep(100 * time.Millisecond)

		// ======================================================================
		// 3단계: 낙관적 락을 사용한 업데이트
		// ======================================================================
		// WHERE 조건에 버전을 포함하여 업데이트
		// 다른 트랜잭션이 먼저 버전을 변경했다면 이 업데이트는 영향받는 행이 0개
		//
		// SQL: UPDATE products SET stock = ?, version = version + 1
		//      WHERE id = ? AND version = ?
		result := tx.Model(&Product{}).
			Where("id = ? AND version = ?", productID, currentVersion).
			Updates(map[string]interface{}{
				"stock":      product.Stock - quantity,
				"version":    currentVersion + 1, // 버전 증가
				"updated_at": time.Now(),
			})

		if result.Error != nil {
			return result.Error
		}

		// ======================================================================
		// 4단계: 업데이트 결과 확인
		// ======================================================================
		// RowsAffected가 0이면 버전이 변경되어 업데이트가 실패한 것
		// 이것이 낙관적 락의 충돌 감지 방식입니다
		if result.RowsAffected == 0 {
			return ErrOptimisticLock
		}

		fmt.Printf("[Worker %d] 구매 성공! %d개 구매, 남은 재고: %d\n",
			workerID, quantity, product.Stock-quantity)

		return nil
	})
}
