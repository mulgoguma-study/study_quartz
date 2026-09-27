# SQL CONCAT & GROUP_CONCAT 연습 문제

이 파일은 `CONCAT`과 `GROUP_CONCAT` 함수를 연습하기 위한 3가지 문제입니다.
각 문제에 대해 적절한 SQL 쿼리를 작성해보세요.

---

## 사전 준비 (Schema Setup)
연습을 위해 아래의 테이블과 데이터를 가정합니다. (MySQL 문법 기준)

```sql
CREATE TABLE employees (
    id INT,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    department VARCHAR(50)
);

INSERT INTO employees VALUES
(1, 'Alice', 'Kim', 'HR'),
(2, 'Bob', 'Lee', 'Engineering'),
(3, 'Charlie', 'Park', 'Engineering');

CREATE TABLE employee_skills (
    emp_id INT,
    skill VARCHAR(50)
);

INSERT INTO employee_skills VALUES
(2, 'Java'), (2, 'Python'), (2, 'SQL'),
(3, 'C++'), (3, 'Go'),
(1, 'Communication');
```

---

## 문제 1: 이름 합치기 (CONCAT)
`employees` 테이블에서 `first_name`과 `last_name`을 합쳐서 `full_name`이라는 별칭으로 출력하세요.
이름 사이에는 공백(space)이 하나 있어야 합니다.

**출력 예시:**
| full_name |
|-----------|
| Alice Kim |
| Bob Lee   |
| ...       |

---

## 문제 2: 문장 만들기 (CONCAT)
`employees` 테이블을 사용하여 각 직원의 소속 부서를 설명하는 문장을 만드세요.
출력 컬럼명은 `description`으로 하세요.

**형식:** "[이름] works in [부서] department."

**출력 예시:**
| description |
|-------------|
| Alice works in HR department. |
| ... |

---

## 문제 3: 기술 목록 나열하기 (GROUP_CONCAT)
각 직원이 가진 기술(skill)을 한 줄로 나열하여 보고 싶습니다.
`employees` 테이블과 `employee_skills` 테이블을 조인하여, 각 직원의 `first_name`과 그 직원이 가진 스킬들을 쉼표(`,`)로 구분된 하나의 문자열로 출력하세요.
스킬 목록 컬럼명은 `skills_list`로 하고, 스킬은 알파벳 순으로 정렬되어야 합니다.

**출력 예시:**
| first_name | skills_list |
|------------|-------------|
| Alice      | Communication |
| Bob        | Java,Python,SQL |
| Charlie    | C++,Go |
