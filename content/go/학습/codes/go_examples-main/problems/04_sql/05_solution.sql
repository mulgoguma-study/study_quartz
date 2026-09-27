-- ================================================================================
-- 문제 5: 인덱스 최적화 - 솔루션
-- ================================================================================

-- ============================================================================
-- 인덱스 설계
-- ============================================================================

/*
【 쿼리 1: 사용자의 최근 주문 】
WHERE user_id = ? ORDER BY created_at DESC LIMIT 10

인덱스: (user_id, created_at DESC)
이유: user_id로 필터 후 created_at으로 정렬된 상태로 10개 반환
*/
CREATE INDEX idx_orders_user_created ON orders(user_id, created_at DESC);


/*
【 쿼리 2: 처리 대기 중인 주문 】
WHERE status = 'pending' AND created_at > ?

인덱스: (status, created_at)
이유:
- status로 필터 (카디널리티 낮음)
- created_at으로 Range 스캔
- pending 상태만 인덱싱 (Partial Index)
*/
-- MySQL
CREATE INDEX idx_orders_status_created ON orders(status, created_at);

-- PostgreSQL (Partial Index)
CREATE INDEX idx_orders_pending ON orders(created_at)
WHERE status = 'pending';


/*
【 쿼리 3: 관리자 대시보드 】
WHERE created_at BETWEEN ? AND ? GROUP BY status

인덱스: (created_at, status)
이유:
- created_at으로 Range 스캔
- status는 GROUP BY에 사용
- 커버링 인덱스로 만들면 더 효율적
*/
CREATE INDEX idx_orders_date_status ON orders(created_at, status);

-- 커버링 인덱스 (INCLUDE로 total_amount 포함)
-- PostgreSQL
CREATE INDEX idx_orders_date_status_cover ON orders(created_at, status)
INCLUDE (total_amount);


/*
【 쿼리 4: 월별 매출 리포트 】
WHERE status = 'delivered' AND created_at >= ? GROUP BY month

인덱스: (status, created_at) -- 쿼리 2와 공유 가능
또는 별도로: (status, DATE(created_at)) -- 함수 기반 인덱스
*/

-- ============================================================================
-- 최종 인덱스 목록
-- ============================================================================

/*
┌─────────────────────────────────────────────────────────────────────────────┐
│                        권장 인덱스 (최소 구성)                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. idx_orders_user_created (user_id, created_at DESC)                     │
│     → 쿼리 1                                                               │
│                                                                             │
│  2. idx_orders_status_created (status, created_at)                         │
│     → 쿼리 2, 4                                                            │
│                                                                             │
│  3. idx_orders_date_status (created_at, status)                            │
│     → 쿼리 3                                                               │
│                                                                             │
│  참고: 쿼리 패턴에 따라 일부 통합 가능                                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
*/

-- ============================================================================
-- 인덱스가 많을 때의 문제점
-- ============================================================================

/*
【 쓰기 성능 저하 】
- INSERT: 모든 인덱스 업데이트 필요
- UPDATE: 인덱스 컬럼 변경 시 재정렬
- DELETE: 인덱스에서도 삭제 필요

【 저장 공간 증가 】
- 인덱스도 디스크 공간 사용
- 테이블 크기의 20-50% 추가될 수 있음

【 옵티마이저 혼란 】
- 너무 많은 인덱스 → 최적 선택 어려움
- 가끔 잘못된 인덱스 선택

【 권장 】
- 테이블당 인덱스 5-7개 이하 권장
- 미사용 인덱스 정기적으로 제거
- 복합 인덱스로 개별 인덱스 통합 검토
*/

-- ============================================================================
-- 파티셔닝 전략
-- ============================================================================

/*
【 Range Partitioning by Date 】
1억 건 → 월별 파티션 → 각 800만 건

장점:
- 오래된 데이터 아카이빙 쉬움 (파티션 드롭)
- 날짜 조건 쿼리 시 파티션 프루닝
- 각 파티션 인덱스 크기 작음
*/

-- MySQL
CREATE TABLE orders (
    id BIGINT,
    user_id INT,
    status VARCHAR(20),
    created_at DATETIME,
    total_amount DECIMAL(10,2),
    PRIMARY KEY (id, created_at)  -- 파티션 키 포함 필수
)
PARTITION BY RANGE (YEAR(created_at) * 100 + MONTH(created_at)) (
    PARTITION p202301 VALUES LESS THAN (202302),
    PARTITION p202302 VALUES LESS THAN (202303),
    -- ...
    PARTITION p_future VALUES LESS THAN MAXVALUE
);

-- PostgreSQL
CREATE TABLE orders (
    id BIGINT,
    user_id INT,
    status VARCHAR(20),
    created_at TIMESTAMP,
    total_amount DECIMAL(10,2)
) PARTITION BY RANGE (created_at);

CREATE TABLE orders_2024_01 PARTITION OF orders
FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

-- ============================================================================
-- EXPLAIN 분석 예시
-- ============================================================================

/*
【 쿼리 1 분석 】
EXPLAIN SELECT * FROM orders
WHERE user_id = 12345
ORDER BY created_at DESC
LIMIT 10;

예상:
- Index Scan using idx_orders_user_created
- Index Cond: user_id = 12345
- Rows: 10 (Limit으로 조기 종료)
- 인덱스 순서대로 반환되므로 별도 정렬 없음

【 나쁜 실행 계획 예시 】
- Seq Scan (Full Table Scan) → 인덱스 추가 필요
- Filesort → ORDER BY 최적화 필요
- Using temporary → GROUP BY 최적화 필요
*/
