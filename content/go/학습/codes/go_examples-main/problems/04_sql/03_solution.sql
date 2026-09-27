-- ================================================================================
-- 문제 3: Department Top Three Salaries - 솔루션
-- ================================================================================

-- ============================================================================
-- 방법 1: DENSE_RANK() 윈도우 함수 (권장)
-- ============================================================================
/*
【 DENSE_RANK 사용 이유 】
- 동일 급여 = 동일 순위
- 순위 건너뛰지 않음 (1,1,2,3...)
- 상위 3개 "급여"를 찾는 것이므로 DENSE_RANK 적합
*/
SELECT d.name AS Department, e.name AS Employee, e.salary AS Salary
FROM (
    SELECT
        name,
        salary,
        departmentId,
        DENSE_RANK() OVER (
            PARTITION BY departmentId
            ORDER BY salary DESC
        ) AS rnk
    FROM Employee
) e
JOIN Department d ON e.departmentId = d.id
WHERE e.rnk <= 3;

-- ============================================================================
-- 방법 2: CTE 사용 (더 가독성 좋음)
-- ============================================================================
WITH RankedEmployees AS (
    SELECT
        e.name AS employee_name,
        e.salary,
        e.departmentId,
        d.name AS department_name,
        DENSE_RANK() OVER (
            PARTITION BY e.departmentId
            ORDER BY e.salary DESC
        ) AS salary_rank
    FROM Employee e
    JOIN Department d ON e.departmentId = d.id
)
SELECT
    department_name AS Department,
    employee_name AS Employee,
    salary AS Salary
FROM RankedEmployees
WHERE salary_rank <= 3
ORDER BY department_name, salary DESC;

-- ============================================================================
-- 방법 3: 서브쿼리 (윈도우 함수 없이)
-- ============================================================================
/*
【 동작 원리 】
같은 부서에서 자신보다 급여가 높은 "고유한" 급여 수가 3개 미만이면
상위 3위 안에 드는 것
*/
SELECT d.name AS Department, e.name AS Employee, e.salary AS Salary
FROM Employee e
JOIN Department d ON e.departmentId = d.id
WHERE (
    SELECT COUNT(DISTINCT e2.salary)
    FROM Employee e2
    WHERE e2.departmentId = e.departmentId AND e2.salary > e.salary
) < 3
ORDER BY d.name, e.salary DESC;

-- ============================================================================
-- 실무 확장: Top N 제품 by 카테고리
-- ============================================================================
/*
【 전자상거래 예시 】
각 카테고리별 판매량 상위 5개 제품
*/
-- WITH RankedProducts AS (
--     SELECT
--         p.product_name,
--         c.category_name,
--         SUM(oi.quantity) AS total_sold,
--         DENSE_RANK() OVER (
--             PARTITION BY p.category_id
--             ORDER BY SUM(oi.quantity) DESC
--         ) AS sales_rank
--     FROM Products p
--     JOIN Categories c ON p.category_id = c.id
--     JOIN OrderItems oi ON p.id = oi.product_id
--     GROUP BY p.id, p.product_name, c.category_name
-- )
-- SELECT category_name, product_name, total_sold
-- FROM RankedProducts
-- WHERE sales_rank <= 5;
