-- ================================================================================
-- 문제 2: Consecutive Numbers - 솔루션
-- ================================================================================

-- ============================================================================
-- 방법 1: Self Join
-- ============================================================================
SELECT DISTINCT l1.num AS ConsecutiveNums
FROM Logs l1
JOIN Logs l2 ON l1.id = l2.id - 1
JOIN Logs l3 ON l2.id = l3.id - 1
WHERE l1.num = l2.num AND l2.num = l3.num;

-- ============================================================================
-- 방법 2: LAG/LEAD 윈도우 함수 (확장 가능)
-- ============================================================================
SELECT DISTINCT num AS ConsecutiveNums
FROM (
    SELECT
        num,
        LAG(num, 1) OVER (ORDER BY id) AS prev1,
        LAG(num, 2) OVER (ORDER BY id) AS prev2
    FROM Logs
) t
WHERE num = prev1 AND num = prev2;

-- ============================================================================
-- 방법 3: Gap and Islands (연속 구간 찾기)
-- ============================================================================
/*
【 Gap and Islands 패턴 】
연속된 같은 값의 그룹을 식별하고, 각 그룹의 크기를 계산
*/
WITH numbered AS (
    SELECT
        id,
        num,
        id - ROW_NUMBER() OVER (PARTITION BY num ORDER BY id) AS grp
    FROM Logs
),
grouped AS (
    SELECT num, grp, COUNT(*) as cnt
    FROM numbered
    GROUP BY num, grp
)
SELECT DISTINCT num AS ConsecutiveNums
FROM grouped
WHERE cnt >= 3;

-- ============================================================================
-- 일반화: N번 이상 연속
-- ============================================================================
/*
【 N=4인 경우 】
SELECT DISTINCT l1.num
FROM Logs l1
JOIN Logs l2 ON l1.id = l2.id - 1 AND l1.num = l2.num
JOIN Logs l3 ON l2.id = l3.id - 1 AND l2.num = l3.num
JOIN Logs l4 ON l3.id = l4.id - 1 AND l3.num = l4.num;

→ N이 커지면 Gap and Islands 방식이 더 효율적
*/
