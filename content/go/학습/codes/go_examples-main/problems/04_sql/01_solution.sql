-- ================================================================================
-- 문제 1: Second Highest Salary - 솔루션
-- ================================================================================

-- ============================================================================
-- 방법 1: 서브쿼리 + LIMIT OFFSET (권장)
-- ============================================================================
/*
【 왜 이 방법이 좋은가? 】
1. 간단하고 직관적
2. 결과가 없을 때 자동으로 NULL 반환
3. 인덱스가 있으면 효율적
*/

SELECT (
    SELECT DISTINCT salary
    FROM Employee
    ORDER BY salary DESC
    LIMIT 1 OFFSET 1
) AS SecondHighestSalary;

-- ============================================================================
-- 방법 2: MAX + 서브쿼리
-- ============================================================================
/*
【 동작 원리 】
최대값보다 작은 값 중 최대값 = 두 번째로 큰 값
*/

SELECT MAX(salary) AS SecondHighestSalary
FROM Employee
WHERE salary < (SELECT MAX(salary) FROM Employee);

-- ============================================================================
-- 방법 3: DENSE_RANK() 윈도우 함수 (확장 가능)
-- ============================================================================
/*
【 장점 】
N번째 급여를 찾을 때 쉽게 확장 가능
*/

SELECT salary AS SecondHighestSalary
FROM (
    SELECT salary, DENSE_RANK() OVER (ORDER BY salary DESC) AS rnk
    FROM Employee
) ranked
WHERE rnk = 2
LIMIT 1;

-- 또는 CTE 사용
WITH RankedSalaries AS (
    SELECT salary, DENSE_RANK() OVER (ORDER BY salary DESC) AS rnk
    FROM Employee
)
SELECT (
    SELECT salary FROM RankedSalaries WHERE rnk = 2 LIMIT 1
) AS SecondHighestSalary;

-- ============================================================================
-- 관련 문제: N번째로 높은 급여
-- ============================================================================
/*
【 문제 】
N번째로 높은 급여를 반환하는 함수를 작성하세요.

【 MySQL 솔루션 】
CREATE FUNCTION getNthHighestSalary(N INT) RETURNS INT
BEGIN
    SET N = N - 1;
    RETURN (
        SELECT DISTINCT salary
        FROM Employee
        ORDER BY salary DESC
        LIMIT 1 OFFSET N
    );
END;

【 PostgreSQL 솔루션 】
CREATE OR REPLACE FUNCTION getNthHighestSalary(N INT)
RETURNS INT AS $$
BEGIN
    RETURN (
        SELECT DISTINCT salary
        FROM Employee
        ORDER BY salary DESC
        LIMIT 1 OFFSET N - 1
    );
END;
$$ LANGUAGE plpgsql;
*/

-- ============================================================================
-- RANK vs DENSE_RANK vs ROW_NUMBER
-- ============================================================================
/*
salary: 300, 300, 200, 100

ROW_NUMBER: 1, 2, 3, 4  (항상 연속)
RANK:       1, 1, 3, 4  (동률 후 건너뜀)
DENSE_RANK: 1, 1, 2, 3  (동률 후 안 건너뜀)

두 번째로 높은 급여에는 DENSE_RANK가 적합 (동률 처리)
*/
