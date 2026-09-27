# SQL CONCAT & GROUP_CONCAT 정답 및 해설

---

## 문제 1: 이름 합치기 (CONCAT)

### 정답 쿼리
```sql
SELECT
    CONCAT(first_name, ' ', last_name) AS full_name
FROM
    employees;
```

### 해설
- `CONCAT(str1, str2, ...)` 함수는 여러 문자열을 하나로 합칩니다.
- 중간에 공백 `' '`을 인자로 넣어 이름과 성 사이를 띄워줍니다.

---

## 문제 2: 문장 만들기 (CONCAT)

### 정답 쿼리
```sql
SELECT
    CONCAT(first_name, ' works in ', department, ' department.') AS description
FROM
    employees;
```

### 해설
- 컬럼(`first_name`, `department`)과 고정된 문자열(`' works in '`, `' department.'`)을 자유롭게 섞어서 사용할 수 있습니다.

---

## 문제 3: 기술 목록 나열하기 (GROUP_CONCAT)

### 정답 쿼리
```sql
SELECT
    e.first_name,
    GROUP_CONCAT(s.skill ORDER BY s.skill ASC SEPARATOR ',') AS skills_list
FROM
    employees e
JOIN
    employee_skills s ON e.id = s.emp_id
GROUP BY
    e.id, e.first_name;
```

### 해설
- `GROUP_CONCAT(expr)`은 `GROUP BY`로 묶인 그룹 내의 여러 행의 값을 하나의 문자열로 합쳐줍니다.
- `ORDER BY s.skill ASC`: 합쳐지는 문자열 내부에서의 정렬 순서를 지정합니다.
- `SEPARATOR ','`: 각 값을 구분할 구분자를 지정합니다 (기본값이 `,`이므로 생략 가능하지만 명시하는 것이 좋습니다).
- `GROUP BY`: 집계 함수를 사용하므로 직원별로 그룹화해야 합니다.
