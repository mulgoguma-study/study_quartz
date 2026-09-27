// =============================================================================
// ACID 기반 데이터 관리/처리 예제
// =============================================================================
// ACID는 데이터베이스 트랜잭션의 4가지 핵심 속성입니다:
//
// A - Atomicity (원자성)
//     : 트랜잭션의 모든 작업이 완전히 수행되거나, 전혀 수행되지 않아야 함
//     : "All or Nothing"
//
// C - Consistency (일관성)
//     : 트랜잭션 전후로 데이터베이스가 일관된 상태를 유지해야 함
//     : 모든 제약 조건(외래키, 유니크 등)이 만족되어야 함
//
// I - Isolation (격리성)
//     : 동시에 실행되는 트랜잭션들이 서로 영향을 주지 않아야 함
//     : 각 트랜잭션은 다른 트랜잭션의 중간 결과를 볼 수 없음
//
// D - Durability (지속성)
//     : 커밋된 트랜잭션의 결과는 영구적으로 저장되어야 함
//     : 시스템 장애가 발생해도 데이터가 유지됨
// =============================================================================

package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User 사용자 모델
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"uniqueIndex;size:200;not null"` // 유니크 제약조건
	Name      string    `gorm:"size:100;not null"`
	Balance   int64     `gorm:"not null;default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TransferLog 송금 기록 모델
type TransferLog struct {
	ID         uint      `gorm:"primaryKey"`
	FromUserID uint      `gorm:"not null;index"`
	ToUserID   uint      `gorm:"not null;index"`
	Amount     int64     `gorm:"not null"`
	Status     string    `gorm:"size:20;not null"` // SUCCESS, FAILED
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

// 커스텀 에러 정의
var (
	ErrInsufficientBalance = errors.New("잔액이 부족합니다")
	ErrUserNotFound        = errors.New("사용자를 찾을 수 없습니다")
	ErrSameUser            = errors.New("같은 계정으로 송금할 수 없습니다")
	ErrInvalidAmount       = errors.New("유효하지 않은 금액입니다")
)

func main() {
	// ==========================================================================
	// 데이터베이스 설정
	// ==========================================================================
	db, err := gorm.Open(sqlite.Open("acid_example.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("데이터베이스 연결 실패:", err)
	}

	// 테이블 마이그레이션
	if err := db.AutoMigrate(&User{}, &TransferLog{}); err != nil {
		log.Fatal("마이그레이션 실패:", err)
	}

	// 테스트 데이터 초기화
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM transfer_logs")

	// 테스트 사용자 생성
	users := []User{
		{Email: "alice@example.com", Name: "Alice", Balance: 100000},
		{Email: "bob@example.com", Name: "Bob", Balance: 50000},
		{Email: "charlie@example.com", Name: "Charlie", Balance: 30000},
	}

	for _, user := range users {
		db.Create(&user)
	}

	fmt.Println("=== 초기 사용자 잔액 ===")
	printAllUsers(db)

	// ==========================================================================
	// 테스트 시나리오 실행
	// ==========================================================================

	// 시나리오 1: 정상적인 송금 (ACID 속성 확인)
	fmt.Println("\n=== 시나리오 1: 정상 송금 (Alice → Bob, 30,000원) ===")
	err = transferMoney(db, 1, 2, 30000)
	if err != nil {
		fmt.Printf("송금 실패: %v\n", err)
	} else {
		fmt.Println("송금 성공!")
	}
	printAllUsers(db)

	// 시나리오 2: 잔액 부족으로 인한 롤백 (Atomicity 확인)
	fmt.Println("\n=== 시나리오 2: 잔액 부족 (Alice → Charlie, 100,000원) ===")
	err = transferMoney(db, 1, 3, 100000)
	if err != nil {
		fmt.Printf("송금 실패 (예상됨): %v\n", err)
	}
	fmt.Println("잔액 변동 없음 확인 (원자성):")
	printAllUsers(db)

	// 시나리오 3: 다중 송금 트랜잭션 (여러 작업을 하나의 트랜잭션으로)
	fmt.Println("\n=== 시나리오 3: 다중 송금 (Bob → Alice, Charlie 각각 10,000원) ===")
	err = multiTransfer(db, 2, []uint{1, 3}, 10000)
	if err != nil {
		fmt.Printf("다중 송금 실패: %v\n", err)
	} else {
		fmt.Println("다중 송금 성공!")
	}
	printAllUsers(db)

	// 시나리오 4: 다중 송금 중 실패 (중간에 실패해도 전체 롤백)
	fmt.Println("\n=== 시나리오 4: 다중 송금 중 잔액 부족 (Charlie → Alice, Bob 각각 25,000원) ===")
	err = multiTransfer(db, 3, []uint{1, 2}, 25000)
	if err != nil {
		fmt.Printf("다중 송금 실패 (예상됨): %v\n", err)
	}
	fmt.Println("모든 잔액이 원래대로 유지 (원자성):")
	printAllUsers(db)

	// 송금 기록 출력
	fmt.Println("\n=== 송금 기록 ===")
	printTransferLogs(db)
}

// transferMoney 송금 함수 - ACID 속성을 보장하는 트랜잭션
//
// 이 함수는 ACID의 모든 속성을 보여줍니다:
// - Atomicity: 출금, 입금, 로그 기록이 모두 성공하거나 모두 롤백
// - Consistency: 송금 전후 총 잔액이 동일하게 유지
// - Isolation: 다른 트랜잭션과 격리되어 실행
// - Durability: 커밋 후 결과가 영구 저장
func transferMoney(db *gorm.DB, fromUserID, toUserID uint, amount int64) error {
	// ==========================================================================
	// 입력값 검증 (Consistency 보장의 일부)
	// ==========================================================================
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if fromUserID == toUserID {
		return ErrSameUser
	}

	// ==========================================================================
	// 트랜잭션 시작
	// ==========================================================================
	// db.Transaction()을 사용하면 GORM이 자동으로 트랜잭션을 관리합니다
	// - 정상 완료 시: 자동 커밋 (Durability 보장)
	// - 에러 발생 시: 자동 롤백 (Atomicity 보장)
	return db.Transaction(func(tx *gorm.DB) error {
		var fromUser, toUser User

		// ======================================================================
		// 송금자 조회 및 검증
		// ======================================================================
		if err := tx.First(&fromUser, fromUserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return err
		}

		// 잔액 확인 (Consistency 보장)
		if fromUser.Balance < amount {
			return ErrInsufficientBalance
		}

		// ======================================================================
		// 수신자 조회
		// ======================================================================
		if err := tx.First(&toUser, toUserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return err
		}

		// ======================================================================
		// 송금 처리 (Atomicity - 이 세 작업은 모두 성공하거나 모두 실패)
		// ======================================================================

		// 1. 송금자 잔액 차감
		if err := tx.Model(&fromUser).Update("balance", fromUser.Balance-amount).Error; err != nil {
			return fmt.Errorf("출금 실패: %w", err)
		}
		fmt.Printf("  → %s 계좌에서 %d원 출금\n", fromUser.Name, amount)

		// 2. 수신자 잔액 증가
		if err := tx.Model(&toUser).Update("balance", toUser.Balance+amount).Error; err != nil {
			// 이 에러가 발생하면 위의 출금도 롤백됨 (Atomicity)
			return fmt.Errorf("입금 실패: %w", err)
		}
		fmt.Printf("  → %s 계좌로 %d원 입금\n", toUser.Name, amount)

		// 3. 송금 기록 저장
		transferLog := TransferLog{
			FromUserID: fromUserID,
			ToUserID:   toUserID,
			Amount:     amount,
			Status:     "SUCCESS",
		}
		if err := tx.Create(&transferLog).Error; err != nil {
			// 로그 저장 실패 시에도 전체 롤백 (Atomicity)
			return fmt.Errorf("로그 기록 실패: %w", err)
		}

		// ======================================================================
		// 트랜잭션 정상 완료
		// ======================================================================
		// nil을 반환하면 GORM이 자동으로 커밋합니다 (Durability)
		return nil
	})
}

// multiTransfer 다중 송금 - 한 사람이 여러 사람에게 동시에 송금
//
// 이 함수는 Atomicity의 강력한 예시입니다:
// 중간에 하나라도 실패하면 이미 완료된 송금도 모두 롤백됩니다.
func multiTransfer(db *gorm.DB, fromUserID uint, toUserIDs []uint, amountEach int64) error {
	if amountEach <= 0 {
		return ErrInvalidAmount
	}

	// 전체를 하나의 트랜잭션으로 처리
	return db.Transaction(func(tx *gorm.DB) error {
		var fromUser User

		// 송금자 조회
		if err := tx.First(&fromUser, fromUserID).Error; err != nil {
			return ErrUserNotFound
		}

		// 총 필요 금액 계산
		totalAmount := amountEach * int64(len(toUserIDs))

		// 잔액 확인
		if fromUser.Balance < totalAmount {
			return fmt.Errorf("%w (필요: %d, 보유: %d)",
				ErrInsufficientBalance, totalAmount, fromUser.Balance)
		}

		// 각 수신자에게 송금
		for i, toUserID := range toUserIDs {
			if fromUserID == toUserID {
				continue
			}

			var toUser User
			if err := tx.First(&toUser, toUserID).Error; err != nil {
				return ErrUserNotFound
			}

			// 중간 검증: 남은 잔액으로 송금 가능한지
			remainingBalance := fromUser.Balance - (amountEach * int64(i))
			if remainingBalance < amountEach {
				// 여기서 에러가 발생하면 이전의 모든 송금도 롤백!
				return ErrInsufficientBalance
			}

			// 송금자 잔액 차감
			fromUser.Balance -= amountEach
			if err := tx.Model(&fromUser).Update("balance", fromUser.Balance).Error; err != nil {
				return err
			}

			// 수신자 잔액 증가
			if err := tx.Model(&toUser).Update("balance", toUser.Balance+amountEach).Error; err != nil {
				return err
			}

			// 로그 기록
			tx.Create(&TransferLog{
				FromUserID: fromUserID,
				ToUserID:   toUserID,
				Amount:     amountEach,
				Status:     "SUCCESS",
			})

			fmt.Printf("  → %s → %s: %d원 송금 완료\n",
				fromUser.Name, toUser.Name, amountEach)
		}

		return nil
	})
}

// printAllUsers 모든 사용자 정보 출력
func printAllUsers(db *gorm.DB) {
	var users []User
	db.Order("id").Find(&users)

	var total int64
	for _, user := range users {
		fmt.Printf("  ID=%d, 이름=%-10s, 잔액=%,d원\n",
			user.ID, user.Name, user.Balance)
		total += user.Balance
	}
	fmt.Printf("  [총 잔액: %d원]\n", total)
}

// printTransferLogs 송금 기록 출력
func printTransferLogs(db *gorm.DB) {
	var logs []TransferLog
	db.Order("id").Find(&logs)

	for _, log := range logs {
		fmt.Printf("  ID=%d, From=%d → To=%d, 금액=%d원, 상태=%s\n",
			log.ID, log.FromUserID, log.ToUserID, log.Amount, log.Status)
	}
}
