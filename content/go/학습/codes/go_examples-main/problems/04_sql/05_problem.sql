-- ================================================================================
-- 문제 5: 인덱스 최적화 문제
-- ================================================================================
-- 난이도: Hard (실무)
-- 주제: 쿼리 최적화, 인덱스 설계
--
-- 【 시나리오 】
-- 이커머스 주문 테이블이 있습니다. 1억 건 이상의 데이터가 있습니다.
--
-- Orders:
-- - id (PK)
-- - user_id
-- - status (pending, processing, shipped, delivered, cancelled)
-- - created_at
-- - updated_at
-- - total_amount
--
-- 【 자주 실행되는 쿼리들 】
-- 1. 특정 사용자의 최근 주문 10개
-- 2. 특정 상태의 주문 (처리 필요한 것)
-- 3. 날짜 범위 + 상태로 조회
-- 4. 월별 매출 집계
--
-- 【 문제 】
-- 1. 어떤 인덱스를 만들어야 할까요?
-- 2. 각 쿼리의 실행 계획을 예측하세요.
-- 3. 인덱스가 너무 많을 때의 문제점은?
-- 4. 파티셔닝 전략은?
--
-- ================================================================================

-- 아래 쿼리들을 최적화하기 위한 인덱스를 설계하세요

-- 쿼리 1: 사용자의 최근 주문
SELECT * FROM orders
WHERE user_id = 12345
ORDER BY created_at DESC
LIMIT 10;

-- 쿼리 2: 처리 대기 중인 주문
SELECT * FROM orders
WHERE status = 'pending'
AND created_at > NOW() - INTERVAL 24 HOUR;

-- 쿼리 3: 관리자 대시보드 (상태 + 날짜 범위)
SELECT status, COUNT(*), SUM(total_amount)
FROM orders
WHERE created_at BETWEEN '2024-01-01' AND '2024-01-31'
GROUP BY status;

-- 쿼리 4: 월별 매출 리포트
SELECT
    DATE_FORMAT(created_at, '%Y-%m') as month,
    COUNT(*) as order_count,
    SUM(total_amount) as revenue
FROM orders
WHERE status = 'delivered'
AND created_at >= '2023-01-01'
GROUP BY DATE_FORMAT(created_at, '%Y-%m')
ORDER BY month;


