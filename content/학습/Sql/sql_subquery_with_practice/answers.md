# SQL 다중 서브쿼리 & WITH절(CTE) 정답 및 해설

---

## 문제 1: 평균 이상의 주문 찾기 (서브쿼리 활용)

### 정답 쿼리
```sql
SELECT
    oi.order_id,
    SUM(oi.quantity * p.price) AS total_amount
FROM
    order_items oi
JOIN
    products p ON oi.product_id = p.product_id
GROUP BY
    oi.order_id
HAVING
    SUM(oi.quantity * p.price) > (
        SELECT AVG(order_total)
        FROM (
            SELECT SUM(oi2.quantity * p2.price) AS order_total
            FROM order_items oi2
            JOIN products p2 ON oi2.product_id = p2.product_id
            GROUP BY oi2.order_id
        ) AS sub
    );
```

### 해설
- **Main Query**: `order_items`와 `products`를 조인하여 주문(`order_id`)별 총 금액을 구합니다 (`GROUP BY`).
- **HAVING 절**: 집계된 결과(`SUM`)를 조건으로 걸기 위해 사용합니다.
- **Subquery**: 전체 주문의 평균 금액을 구해야 합니다. 평균을 구하려면 먼저 "각 주문의 총액"이 리스트로 있어야 하므로, `FROM` 절 안에 서브쿼리(Inline View)를 하나 더 써서 주문별 총액을 먼저 만들고, 그 밖에서 `AVG`를 씌웠습니다.

---

## 문제 2: 카테고리별 매출 분석 (WITH 절 활용)

### 정답 쿼리
```sql
WITH CategoryRevenue AS (
    SELECT
        p.category,
        SUM(oi.quantity * p.price) AS total_revenue
    FROM
        order_items oi
    JOIN
        products p ON oi.product_id = p.product_id
    GROUP BY
        p.category
)
SELECT
    category,
    total_revenue
FROM
    CategoryRevenue
WHERE
    total_revenue >= 200;
```

### 해설
- `WITH CategoryRevenue AS (...)`: 카테고리별 매출을 미리 계산하여 `CategoryRevenue`라는 임시 이름을 붙였습니다. 복잡한 집계 로직을 미리 정의해두면 메인 쿼리가 매우 깔끔해집니다.
- 메인 쿼리에서는 단순히 `total_revenue >= 200` 조건만 걸면 되므로 가독성이 좋아집니다.

---

## 문제 3: VIP 고객 찾기 (다중 서브쿼리 or WITH 중첩)

### 정답 쿼리 (WITH 활용 권장)
```sql
WITH CustomerOrderCounts AS (
    -- 1. 고객별 주문 횟수 계산
    SELECT
        customer_id,
        COUNT(order_id) AS order_count
    FROM
        orders
    GROUP BY
        customer_id
),
AverageOrderCount AS (
    -- 2. 전체 고객의 평균 주문 횟수 계산
    SELECT
        AVG(order_count) AS avg_count
    FROM
        CustomerOrderCounts
)
SELECT
    c.name
FROM
    customers c
JOIN
    CustomerOrderCounts coc ON c.customer_id = coc.customer_id
WHERE
    coc.order_count > (SELECT avg_count FROM AverageOrderCount);
```

### 해설
- **다단계 논리**:
    1. `CustomerOrderCounts`: 고객 ID별로 몇 번 주문했는지 셉니다.
    2. `AverageOrderCount`: 위에서 만든 CTE를 다시 활용하여 '주문 횟수의 평균'값 하나를 구합니다.
- **Main Query**: 고객 테이블과 1번 CTE를 조인하고, 2번 CTE의 값을 조건으로 비교합니다.
- 이렇게 `WITH`를 사용하면 논리의 흐름대로 코드를 작성할 수 있어, 복잡한 통계 쿼리에서 매우 유용합니다.
