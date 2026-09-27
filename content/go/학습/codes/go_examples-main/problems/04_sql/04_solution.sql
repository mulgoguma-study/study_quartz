-- ================================================================================
-- 문제 4: Human Traffic of Stadium - 솔루션
-- ================================================================================

-- ============================================================================
-- 방법 1: Gap and Islands
-- ============================================================================
/*
【 Gap and Islands 패턴 】

1. 100명 이상인 행만 필터링
2. ROW_NUMBER로 순번 부여
3. id - ROW_NUMBER = 그룹 식별자 (연속이면 같은 값)
4. 그룹 크기 >= 3인 것만 선택
*/
WITH filtered AS (
    SELECT *, ROW_NUMBER() OVER (ORDER BY id) AS rn
    FROM Stadium
    WHERE people >= 100
),
grouped AS (
    SELECT *, id - rn AS grp
    FROM filtered
),
valid_groups AS (
    SELECT grp
    FROM grouped
    GROUP BY grp
    HAVING COUNT(*) >= 3
)
SELECT g.id, g.visit_date, g.people
FROM grouped g
WHERE g.grp IN (SELECT grp FROM valid_groups)
ORDER BY g.id;

-- ============================================================================
-- 방법 2: Self Join (3개 테이블)
-- ============================================================================
/*
id=N, N+1, N+2가 모두 100 이상인 경우 찾기
*/
SELECT DISTINCT s.*
FROM Stadium s
WHERE s.people >= 100
AND (
    -- s가 연속 3일의 첫 번째
    (EXISTS (SELECT 1 FROM Stadium s2 WHERE s2.id = s.id + 1 AND s2.people >= 100)
     AND EXISTS (SELECT 1 FROM Stadium s3 WHERE s3.id = s.id + 2 AND s3.people >= 100))
    OR
    -- s가 연속 3일의 두 번째
    (EXISTS (SELECT 1 FROM Stadium s2 WHERE s2.id = s.id - 1 AND s2.people >= 100)
     AND EXISTS (SELECT 1 FROM Stadium s3 WHERE s3.id = s.id + 1 AND s3.people >= 100))
    OR
    -- s가 연속 3일의 세 번째
    (EXISTS (SELECT 1 FROM Stadium s2 WHERE s2.id = s.id - 1 AND s2.people >= 100)
     AND EXISTS (SELECT 1 FROM Stadium s3 WHERE s3.id = s.id - 2 AND s3.people >= 100))
)
ORDER BY s.id;

-- ============================================================================
-- 방법 3: LAG/LEAD 윈도우 함수
-- ============================================================================
WITH marked AS (
    SELECT
        id, visit_date, people,
        CASE WHEN people >= 100 THEN 1 ELSE 0 END AS valid,
        LAG(CASE WHEN people >= 100 THEN 1 ELSE 0 END, 1, 0) OVER (ORDER BY id) AS prev1,
        LAG(CASE WHEN people >= 100 THEN 1 ELSE 0 END, 2, 0) OVER (ORDER BY id) AS prev2,
        LEAD(CASE WHEN people >= 100 THEN 1 ELSE 0 END, 1, 0) OVER (ORDER BY id) AS next1,
        LEAD(CASE WHEN people >= 100 THEN 1 ELSE 0 END, 2, 0) OVER (ORDER BY id) AS next2
    FROM Stadium
)
SELECT id, visit_date, people
FROM marked
WHERE valid = 1
AND (
    (prev1 = 1 AND prev2 = 1)  -- 3번째
    OR (prev1 = 1 AND next1 = 1)  -- 2번째
    OR (next1 = 1 AND next2 = 1)  -- 1번째
)
ORDER BY id;
