# SQL 다중 서브쿼리 & WITH절(CTE) 연습 문제

이 파일은 `Subquery`(서브쿼리)와 `WITH`(Common Table Expression) 절을 연습하기 위한 3가지 문제입니다.
복잡한 데이터를 단계별로 가공하는 능력을 키워보세요.

---

## 사전 준비 (Schema Setup)
연습을 위해 전자상거래 데이터를 가정합니다. (MySQL 문법 기준)

```sql
CREATE TABLE customers (
    customer_id INT PRIMARY KEY,
    name VARCHAR(50),
    city VARCHAR(50)
);

CREATE TABLE products (
    product_id INT PRIMARY KEY,
    product_name VARCHAR(50),
    category VARCHAR(50),
    price DECIMAL(10, 2)
);

CREATE TABLE orders (
    order_id INT PRIMARY KEY,
    customer_id INT,
    order_date DATE
);

CREATE TABLE order_items (
    order_id INT,
    product_id INT,
    quantity INT,
    FOREIGN KEY (order_id) REFERENCES orders(order_id),
    FOREIGN KEY (product_id) REFERENCES products(product_id)
);

-- 데이터 예시
INSERT INTO customers VALUES (1, 'John', 'Seoul'), (2, 'Sarah', 'Busan'), (3, 'Mike', 'Seoul');
INSERT INTO products VALUES (101, 'Laptop', 'Electronics', 1200), (102, 'Mouse', 'Electronics', 50), (103, 'T-Shirt', 'Fashion', 30);
INSERT INTO orders VALUES (1001, 1, '2023-01-01'), (1002, 2, '2023-01-02'), (1003, 1, '2023-01-05');
INSERT INTO order_items VALUES (1001, 101, 1), (1001, 102, 2), (1002, 103, 5), (1003, 102, 1);
```

---

## 문제 1: 평균 이상의 주문 찾기 (서브쿼리 활용)
`order_items` 테이블을 사용하여, **한 번의 주문(order_id 기준)에 포함된 총 금액(quantity * price)** 이 전체 주문들의 **평균 주문 금액보다 높은 주문**의 `order_id`와 `총 금액`을 출력하세요. This requires a subquery to calculate the average first. NOTE: `products` 테이블의 `price`를 참조해야 합니다.

**힌트:**
1. 각 주문별 총 금액을 구하는 쿼리를 먼저 생각해보세요.
2. 그 결과의 평균을 구하는 쿼리를 서브쿼리로 사용하세요.

**출력 예시:**
| order_id | total_amount |
|----------|--------------|
| 1001     | 1300.00      |
| ...      | ...          |

---

## 문제 2: 카테고리별 매출 분석 (WITH 절 활용)
`WITH` 절을 사용하여 다음 단계로 쿼리를 작성하세요.
1. **CTE1**: 각 `category`별 총 매출(`quantity * price`의 합)을 계산하는 임시 테이블을 만드세요.
2. **Main Query**: 위에서 만든 CTE를 사용하여 총 매출이 200 이상인 카테고리의 `category`와 `total_revenue`를 출력하세요.

**출력 예시:**
| category    | total_revenue |
|-------------|---------------|
| Electronics | 1350.00       |
| ...         | ...           |

---

## 문제 3: VIP 고객 찾기 (다중 서브쿼리 or WITH 중첩)
**"평균 구매 횟수"보다 더 많은 주문을 한 고객의 이름**을 찾고 싶습니다.
다음 단계의 로직을 포함하여 쿼리를 작성하세요.
1. 고객별 주문 횟수를 셉니다.
2. 전체 고객의 평균 주문 횟수를 계산합니다.
3. 내 주문 횟수가 전체 평균보다 큰 고객의 이름을 `customers` 테이블에서 찾아 출력합니다.

(서브쿼리를 WHERE 절에 중첩해서 쓰거나, CTE를 활용해 가독성 있게 풀어보세요.)

**출력 예시:**
| name |
|------|
| John |
