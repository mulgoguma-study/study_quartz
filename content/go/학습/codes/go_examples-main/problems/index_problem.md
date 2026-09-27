# Go & SQL 연습 문제 목차

## 개요

이 폴더에는 시니어 백엔드 개발자를 위한 Go와 SQL 연습 문제가 포함되어 있습니다.
각 문제는 `_problem` 파일과 `_solution` 파일로 구성되어 있습니다.

- **problem**: 문제 설명과 빈 함수/쿼리
- **solution**: 상세한 풀이와 설명

---

## huhu_1 - 42dot (포티투닷) 맞춤 문제 (10문제)

### 자율주행/모빌리티 도메인 (1-5번)
| # | 문제 | 주제 | 핵심 개념 |
|---|------|------|----------|
| 1 | **실시간 차량 위치 추적** | Geospatial, 실시간 | Geohash, 샤딩, Haversine 거리 |
| 2 | **실시간 경로 최적화** | 그래프 알고리즘 | A* 알고리즘, K-shortest paths, 캐싱 |
| 3 | **V2X 메시지 처리** | 실시간 스트리밍 | 우선순위 큐, 워커 풀, Pub/Sub |
| 4 | **자율주행 상태 머신** | 상태 머신 | FSM, Guard/Action, Observer 패턴 |
| 5 | **실시간 배차 알고리즘** | 최적화 | Greedy 매칭, 스코어링 함수 |

### 백엔드/데이터 플랫폼 (6-10번)
| # | 문제 | 주제 | 핵심 개념 |
|---|------|------|----------|
| 6 | **텔레메트리 데이터 파이프라인** | 데이터 파이프라인 | Ring Buffer, 워커 풀, Back-pressure |
| 7 | **시계열 데이터 저장소** | 시계열 DB | 파티셔닝, 다운샘플링, 집계 |
| 8 | **이벤트 소싱 시스템** | 이벤트 소싱 | CQRS, 스냅샷, 이벤트 재생 |
| 9 | **API 게이트웨이** | Rate Limiting | Sliding Window, LRU 캐시, Singleflight |
| 10 | **분산 작업 스케줄러** | 분산 시스템 | 우선순위 큐, 헬스체크, 지수 백오프 |

---

## 02_medium - Go 중급 문제 (10문제)

| # | 문제 | 주제 | 핵심 개념 |
|---|------|------|----------|
| 1 | **LRU Cache** | 자료구조 | HashMap + Doubly Linked List, O(1) 연산 |
| 2 | **Rate Limiter** | 시스템 설계 | Token Bucket, Sliding Window, 동시성 |
| 3 | **Sharded Map** | 동시성 | 샤딩, RWMutex, 락 경합 최소화 |
| 4 | **Worker Pool** | 동시성 | 채널, Context, Graceful Shutdown |
| 5 | **Sliding Window Maximum** | 알고리즘 | Monotonic Deque, O(n) 최적화 |
| 6 | **Merge Intervals** | 알고리즘 | 정렬, 구간 병합 |
| 7 | **Top K Frequent** | 알고리즘 | 버킷 정렬, 힙 |
| 8 | **Trie** | 자료구조 | 접두사 트리, 자동완성 |
| 9 | **Clone Graph** | 알고리즘 | DFS/BFS, 순환 참조 처리 |
| 10 | **Design Twitter** | 시스템 설계 | k-way merge, 피드 시스템 |

---

## 03_hard - Go 고급 문제 (5문제)

| # | 문제 | 주제 | 핵심 개념 |
|---|------|------|----------|
| 1 | **Distributed Lock** | 분산 시스템 | Redis, SET NX, Lua Script, 락 안전성 |
| 2 | **Circuit Breaker** | 분산 시스템 | 상태 머신, 장애 복구, 연쇄 장애 방지 |
| 3 | **Serialize Binary Tree** | 알고리즘 | Pre-order/Level-order, 직렬화 |
| 4 | **Word Ladder** | 알고리즘 | BFS, 양방향 BFS, 최단 경로 |
| 5 | **Median from Stream** | 알고리즘 | Two Heaps, 온라인 알고리즘 |

---

## 04_sql - SQL 문제 (5문제)

| # | 문제 | 난이도 | 핵심 개념 |
|---|------|--------|----------|
| 1 | **Second Highest Salary** | Medium | 서브쿼리, LIMIT OFFSET, NULL 처리 |
| 2 | **Consecutive Numbers** | Medium | Self Join, LAG/LEAD, Gap and Islands |
| 3 | **Dept Top 3 Salaries** | Hard | DENSE_RANK, PARTITION BY |
| 4 | **Stadium Traffic** | Hard | Gap and Islands 패턴, 연속 조건 |
| 5 | **Index Optimization** | Hard | 인덱스 설계, EXPLAIN 분석, 파티셔닝 |

---

## 학습 순서 추천

### 0단계: 42dot 맞춤 문제
**자율주행/모빌리티 도메인:**
```
huhu_1/01_problem.go  - 실시간 차량 위치 추적 (Geohash, 동시성)
huhu_1/02_problem.go  - 실시간 경로 최적화 (A*, Yen's Algorithm)
huhu_1/03_problem.go  - V2X 메시지 처리 (우선순위 큐, 워커 풀)
huhu_1/04_problem.go  - 자율주행 상태 머신 (FSM, 상태 전이)
huhu_1/05_problem.go  - 실시간 배차 알고리즘 (매칭 최적화)
```

**백엔드/데이터 플랫폼:**
```
huhu_1/06_problem.go  - 텔레메트리 데이터 파이프라인 (Ring Buffer, Back-pressure)
huhu_1/07_problem.go  - 시계열 데이터 저장소 (파티셔닝, 다운샘플링)
huhu_1/08_problem.go  - 이벤트 소싱 시스템 (CQRS, 스냅샷)
huhu_1/09_problem.go  - API 게이트웨이 (Rate Limiting, 캐싱)
huhu_1/10_problem.go  - 분산 작업 스케줄러 (우선순위 큐, 재시도)
```

### 1단계: 기본 알고리즘
```
02_medium/05_problem.go  - Sliding Window Maximum
02_medium/06_problem.go  - Merge Intervals
02_medium/07_problem.go  - Top K Frequent
```

### 2단계: 자료구조 설계
```
02_medium/01_problem.go  - LRU Cache
02_medium/08_problem.go  - Trie
02_medium/09_problem.go  - Clone Graph
```

### 3단계: 동시성
```
02_medium/03_problem.go  - Sharded Map
02_medium/04_problem.go  - Worker Pool
02_medium/02_problem.go  - Rate Limiter
```

### 4단계: 분산 시스템
```
03_hard/01_problem.go    - Distributed Lock
03_hard/02_problem.go    - Circuit Breaker
02_medium/10_problem.go  - Design Twitter
```

### 5단계: SQL 심화
```
04_sql/01_problem.sql    - Second Highest Salary
04_sql/02_problem.sql    - Consecutive Numbers
04_sql/03_problem.sql    - Dept Top 3 Salaries
04_sql/04_problem.sql    - Stadium Traffic
04_sql/05_problem.sql    - Index Optimization
```

---

## 주제별 분류

### 42dot / 자율주행 (Autonomous Driving)
- 실시간 차량 위치 추적 (Geohash, 공간 인덱싱)
- 경로 최적화 (A* 알고리즘, K-shortest paths)
- V2X 통신 (메시지 우선순위, 실시간 처리)
- 차량 상태 머신 (FSM, 상태 전이)
- 배차 알고리즘 (매칭 최적화, 스코어링)

### 42dot / 백엔드 플랫폼 (Backend Platform)
- 텔레메트리 파이프라인 (Ring Buffer, Back-pressure)
- 시계열 데이터 저장소 (파티셔닝, 다운샘플링, 집계)
- 이벤트 소싱 (CQRS, 스냅샷, 이벤트 재생)
- API 게이트웨이 (Rate Limiting, LRU 캐시, Singleflight)
- 분산 작업 스케줄러 (우선순위 큐, 헬스체크, 재시도)

### 동시성 (Concurrency)
- Rate Limiter (Token Bucket, Sliding Window)
- Sharded Map (락 샤딩)
- Worker Pool (고루틴 풀)
- Distributed Lock (분산 락)

### 시스템 설계 (System Design)
- LRU Cache (캐시 정책)
- Rate Limiter (API 제한)
- Circuit Breaker (장애 복구)
- Design Twitter (피드 시스템)

### 알고리즘 (Algorithms)
- Sliding Window Maximum (Monotonic Deque)
- Merge Intervals (구간 처리)
- Top K Frequent (힙, 버킷 정렬)
- Word Ladder (BFS)
- Median from Stream (Two Heaps)

### 자료구조 (Data Structures)
- LRU Cache (HashMap + LinkedList)
- Trie (접두사 트리)
- Clone Graph (그래프 복사)
- Serialize Binary Tree (트리 직렬화)

### SQL
- 윈도우 함수 (RANK, DENSE_RANK, LAG, LEAD)
- Gap and Islands 패턴
- 인덱스 최적화
- 쿼리 실행 계획

---

## 실행 방법

### Go 문제
```bash
# 문제 확인
cat problems/02_medium/01_problem.go

# 솔루션 확인
cat problems/02_medium/01_solution.go

# 솔루션 실행 (독립 실행 가능한 경우)
cd problems/02_medium
go run 01_solution.go
```

### SQL 문제
```bash
# 문제 확인
cat problems/04_sql/01_problem.sql

# 솔루션 확인
cat problems/04_sql/01_solution.sql

# MySQL에서 실행
mysql -u root -p < problems/04_sql/01_solution.sql
```

---

## 참고 자료

### 알고리즘
- [LeetCode](https://leetcode.com/)
- [HackerRank](https://www.hackerrank.com/)

### Go
- [Gophercises](https://gophercises.com/)
- [Go by Example](https://gobyexample.com/)

### 시스템 설계
- [System Design Primer](https://github.com/donnemartin/system-design-primer)
- [Designing Data-Intensive Applications](https://dataintensive.net/)

---

## 난이도 가이드

| 난이도 | 예상 시간 | 설명 |
|--------|----------|------|
| Medium | 20-40분 | 핵심 개념 이해 필요 |
| Hard | 40-60분 | 복잡한 구현, 최적화 필요 |

---

## 업데이트 로그

- 2024-01: 초기 문제 세트 생성
  - Medium Go 문제 10개
  - Hard Go 문제 5개
  - SQL 문제 5개
