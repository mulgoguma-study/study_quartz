● 모든 예제 파일을 생성 완료했습니다! 프로젝트 구조는 다음과 같습니다:

  examples/
  ├── 01_gorm/
  │   ├── pessimistic_lock.go   # 비관적 락 (SELECT FOR UPDATE)
  │   ├── optimistic_lock.go    # 낙관적 락 (버전 필드 사용)
  │   └── acid_transaction.go   # ACID 트랜잭션 (송금 예제)
  ├── 02_goroutine/
  │   ├── waitgroup.go          # sync.WaitGroup
  │   ├── fanin_fanout.go       # Fan-In/Fan-Out 패턴
  │   ├── mutex.go              # Mutex, RWMutex
  │   ├── channel.go            # 채널 제어 (select, timeout, done)
  │   ├── workerpool.go         # Worker Pool 패턴
  │   └── context.go            # Context (Cancel, Timeout, Value)
  ├── 03_interface/
  │   ├── duck_typing.go        # 덕 타이핑, 빈 인터페이스
  │   └── interface_composition.go  # 인터페이스 조합, 구조체 임베딩
  ├── 04_error/
  │   └── error_handling.go     # errors.Is, errors.As, 에러 래핑
  ├── 05_test/
  │   ├── testing_basics.go     # 테스트 대상 코드
  │   ├── testing_basics_test.go    # 테스팅 기법 예제
  │   ├── race_condition.go     # 레이스 컨디션 대상 코드
  │   └── race_condition_test.go    # 레이스 컨디션 테스트
  ├── 06_algorithms/
  │   ├── 01_bfs_dfs.go              # BFS/DFS 그래프 탐색
  │   ├── 02_sliding_window.go       # 슬라이딩 윈도우
  │   ├── 03_hashmap_frequency.go    # HashMap 빈도 계산
  │   ├── 04_two_pointer.go          # Two Pointer
  │   ├── 05_dynamic_programming.go  # 동적 프로그래밍 (1D/2D)
  │   ├── 06_greedy.go               # 조건 있는 그리디
  │   ├── 07_dijkstra.go             # 다익스트라
  │   ├── 08_binary_search.go        # 이진 탐색
  │   ├── 09_backtracking.go         # 백트래킹
  │   ├── 10_union_find.go           # 유니온 파인드
  │   └── combined_example.go        # 알고리즘 통합 예제
  ├── 07_senior_topics/
  │   ├── 01_slice_vs_map_internals.go  # 슬라이스/맵 내부 구조
  │   ├── 02_memory_leak_prevention.go  # 메모리 누수 방지
  │   ├── 03_concurrency_deep_dive.go   # 동시성 심화
  │   ├── 04_interface_internals.go     # 인터페이스 내부 구조
  │   ├── 05_channel_internals.go       # 채널 내부 구조
  │   ├── 06_error_handling_philosophy.go # 에러 처리 철학
  │   ├── 07_map_struct_set.go          # map[struct]로 Set 구현
  │   ├── 08_redis_data_types.go        # Redis 데이터 타입
  │   ├── 09_caching_strategies.go      # 캐싱 전략
  │   └── 10_grpc_protobuf.go           # gRPC와 Protocol Buffers
  ├── 08_database/
  │   └── 01_database_fundamentals.go   # DB 심화 (MariaDB vs PostgreSQL, 인덱스, 트랜잭션)
  └── combined_example.go       # 모든 개념 통합 예제 (주문 처리 시스템)

  실행 방법

  각 예제는 독립적으로 실행 가능합니다:

  # 개별 실행 (디렉토리로 이동 후)
  cd examples/01_gorm && go run pessimistic_lock.go
  cd examples/02_goroutine && go run waitgroup.go
  cd examples/03_interface && go run duck_typing.go
  cd examples/04_error && go run error_handling.go

  # 통합 예제 실행
  cd examples && go run combined_example.go

  # 알고리즘 예제 실행 (06_algorithms)
  cd examples/06_algorithms && go run 01_bfs_dfs.go
  cd examples/06_algorithms && go run 08_binary_search.go
  cd examples/06_algorithms && go run combined_example.go

  # 테스트 실행 (05_test)
  cd examples/05_test
  go test -v                           # 전체 테스트
  go test -race -v                     # 레이스 감지 모드
  go test -cover                       # 커버리지 확인
  go test -bench=.                     # 벤치마크 실행

  각 파일 핵심 내용

  | 파일                     | 핵심 개념                                          |
  |--------------------------|----------------------------------------------------|
  | pessimistic_lock.go      | clause.Locking{Strength: "UPDATE"}, 동시 출금 처리 |
  | optimistic_lock.go       | 버전 필드로 충돌 감지, 재시도 로직                 |
  | acid_transaction.go      | db.Transaction(), 원자성/일관성 보장               |
  | waitgroup.go             | Add(), Done(), Wait()                              |
  | fanin_fanout.go          | 여러 채널 병합, 작업 분배                          |
  | mutex.go                 | Race Condition 해결, RWMutex                       |
  | channel.go               | 버퍼드/언버퍼드, select, 타임아웃                  |
  | workerpool.go            | 고정 워커 수, 제네릭 워커 풀                       |
  | context.go               | WithCancel, WithTimeout, WithValue                 |
  | duck_typing.go           | 암시적 인터페이스 구현                             |
  | interface_composition.go | 인터페이스 임베딩, 구조체 임베딩                   |
  | error_handling.go        | 센티널 에러, 커스텀 에러, 에러 체인                |
  | testing_basics_test.go   | 테이블 드리븐 테스트, 서브테스트, 목킹             |
  | race_condition_test.go   | go test -race, 동기화 방식별 비교                  |
  | combined_example.go      | 주문 처리 시스템 (모든 개념 통합)                  |

  6. 알고리즘 (06_algorithms)

  | 파일                      | 핵심 개념                                          |
  |---------------------------|----------------------------------------------------|
  | 01_bfs_dfs.go             | BFS/DFS 그래프 탐색, 최단 경로, 연결 요소          |
  | 02_sliding_window.go      | 고정/가변 윈도우, 부분 배열 최적화                 |
  | 03_hashmap_frequency.go   | 빈도 계산, 중복 탐지, 애너그램                     |
  | 04_two_pointer.go         | 양끝 포인터, 같은 방향 포인터, 정렬 배열 탐색      |
  | 05_dynamic_programming.go | 1D/2D DP, 메모이제이션, 최적 부분 구조             |
  | 06_greedy.go              | 탐욕 선택, 활동 선택, 허프만 코딩                  |
  | 07_dijkstra.go            | 최단 경로, 우선순위 큐, 가중치 그래프              |
  | 08_binary_search.go       | O(log n) 탐색, Lower/Upper Bound, 매개변수 탐색    |
  | 09_backtracking.go        | DFS + 가지치기, 순열, 조합, N-Queens               |
  | 10_union_find.go          | 서로소 집합, 경로 압축, 사이클 감지, MST           |

  7. 시니어 개발자 토픽 (07_senior_topics)

  | 파일                          | 핵심 개념                                          |
  |-------------------------------|----------------------------------------------------|
  | 01_slice_vs_map_internals.go  | 슬라이스/맵 내부 구조, 메모리 레이아웃, 성능 특성  |
  | 02_memory_leak_prevention.go  | 메모리 누수 패턴, 고루틴 누수, 리소스 관리         |
  | 03_concurrency_deep_dive.go   | 동시성 심화, 채널 패턴, 동기화 프리미티브          |
  | 04_interface_internals.go     | 인터페이스 내부 구조, iface/eface, 타입 어설션     |
  | 05_channel_internals.go       | 채널 내부 구조, hchan, 스케줄링                    |
  | 06_error_handling_philosophy.go| 에러 처리 철학, 에러 체인, 센티널 에러            |
  | 07_map_struct_set.go          | map[struct]로 Set 구현, 키 비교 규칙               |
  | 08_redis_data_types.go        | Redis 데이터 타입, 사용 사례, 성능 특성            |
  | 09_caching_strategies.go      | 캐싱 전략, TTL, LRU, 캐시 무효화                   |
  | 10_grpc_protobuf.go           | gRPC 기초, Protocol Buffers, 스트리밍              |

  8. 데이터베이스 심화 (08_database)

  | 파일                          | 핵심 개념                                          |
  |-------------------------------|----------------------------------------------------|
  | 01_database_fundamentals.go   | MariaDB vs PostgreSQL, 인덱스, 트랜잭션, 쿼리 최적화 |

  ### 8번 항목 상세 내용

  #### 8.1 MariaDB vs PostgreSQL 비교
  - 스토리지 엔진: MariaDB(InnoDB, Aria 등 선택) vs PostgreSQL(단일 Heap 기반)
  - MVCC 구현: Undo Log vs Tuple Versioning
  - 인덱스: MariaDB(B+Tree, Hash, R-Tree) vs PostgreSQL(B-Tree, GiST, GIN, BRIN 등)
  - JSON: MariaDB(문자열) vs PostgreSQL(JSONB, 인덱싱 가능)
  - 선택 기준: 단순 CRUD -> MariaDB, 복잡한 쿼리/분석 -> PostgreSQL

  #### 8.2 복합 인덱스 (Composite Index)
  **장점:**
  - 커버링 인덱스로 테이블 접근 없이 쿼리 가능
  - 다중 조건 쿼리 최적화
  - 정렬(ORDER BY) 최적화
  - 공간 효율성 (개별 인덱스 여러 개보다 효율적)

  **문제점:**
  - Leftmost Prefix Rule: 첫 번째 컬럼 없이는 인덱스 사용 불가
  - Range 조건 이후 컬럼은 인덱스 무효화
  - INSERT/UPDATE/DELETE 성능 저하

  ```sql
  -- 인덱스: (user_id, status, created_at)
  -- 사용 가능: WHERE user_id = 1 AND status = 'pending'
  -- 사용 불가: WHERE status = 'pending' (user_id 없음)
  ```

  #### 8.3 DB 처리 vs 백엔드 처리

  | 상황                | DB에서 처리 | 백엔드에서 처리 |
  |---------------------|-------------|-----------------|
  | 데이터량 > 1000건   | O           | X               |
  | 인덱스 있음         | O           | X               |
  | 단순 정렬/집계      | O           | X               |
  | 데이터 < 100건      | △           | O               |
  | 캐시된 데이터       | X           | O               |
  | 복잡한 비즈니스 로직| X           | O               |
  | 다중 소스 조합      | X           | O               |

  **원칙:** 가능하면 DB에서 필터링/정렬/집계 후 최소 데이터만 전송

  #### 8.4 인덱스 쓰기 성능 해결법
  1. **불필요한 인덱스 제거** - 미사용 인덱스 모니터링
  2. **중복 인덱스 통합** - (a)와 (a,b) 있으면 (a) 제거
  3. **배치 INSERT** - 단건보다 다건 INSERT
  4. **대량 로드 시 인덱스 비활성화** - DISABLE KEYS -> 로드 -> ENABLE KEYS
  5. **파티셔닝** - 각 파티션의 인덱스 크기 축소
  6. **읽기/쓰기 분리** - Primary(쓰기, 인덱스 최소) / Replica(읽기, 인덱스 풍부)
  7. **Partial Index** - 조건에 맞는 행만 인덱싱 (PostgreSQL)

  #### 8.5 트랜잭션

  **장점:**
  - 데이터 무결성 (All or Nothing)
  - 동시성 제어
  - 장애 복구 용이

  **단점:**
  - WAL, 락, fsync 오버헤드
  - 락 경합과 데드락 가능성
  - 긴 트랜잭션 = 긴 락 보유

  **Best Practices:**
  - 트랜잭션은 짧게 (외부 API 호출은 밖에서)
  - 적절한 격리 수준 선택
  - 일관된 락 순서로 데드락 방지

  #### 8.6 컬럼 타입 선택

  | 타입      | 크기    | 사용 사례                    |
  |-----------|---------|------------------------------|
  | TINYINT   | 1 byte  | 상태값, 플래그, 나이         |
  | INT       | 4 bytes | 일반 ID, 수량                |
  | BIGINT    | 8 bytes | 대규모 시스템 ID             |
  | CHAR(n)   | 고정    | 국가코드, 통화코드           |
  | VARCHAR(n)| 가변    | 이름, 이메일, 주소           |
  | TIMESTAMP | 4 bytes | UTC 저장, 시간대 변환        |

  **중요:** 1억 건 테이블에서 BIGINT vs TINYINT = 667MB 차이

  #### 8.7 기타 실무 고민들

  1. **N+1 쿼리** - JOIN 또는 IN절 배치 조회로 해결
  2. **페이지네이션** - Cursor 기반 권장 (OFFSET은 느림)
  3. **Connection Pool** - MaxConns 20-50, MaxLifetime 5-30분
  4. **쿼리 최적화** - SELECT * 피하기, LIKE 앞에 % 피하기
  5. **데드락 방지** - 일관된 순서로 락 획득
  6. **Soft Delete** - deleted_at + Partial Index
  7. **EXPLAIN 분석** - Seq Scan(경고), Index Only Scan(좋음)
  8. **샤딩** - 높은 카디널리티, 균등 분산, 쿼리에 포함되는 키 선택
  9. **마이그레이션** - 무중단 스키마 변경 (CONCURRENTLY, pt-online-schema-change)

  5. 테스트 상세 가이드

  5-1. 테스팅 기법 (testing_basics_test.go)

  # 기본 테스트
  go test -v -run TestAdd

  # 테이블 드리븐 테스트
  go test -v -run TestMultiply_TableDriven

  # 서브테스트 개별 실행
  go test -v -run "TestDivide_TableDriven/division_by_zero"

  # errors.Is 테스트
  go test -v -run TestFindUser_ErrorsIs

  # errors.As 테스트
  go test -v -run TestValidateEmail_ErrorsAs

  # 목(Mock) 테스트
  go test -v -run TestUserService

  # 병렬 테스트
  go test -v -run TestContains_Parallel

  # 벤치마크
  go test -bench=. -benchmem

  # 커버리지
  go test -cover
  go test -coverprofile=coverage.out && go tool cover -html=coverage.out

  5-2. 레이스 컨디션 감지 (race_condition_test.go)

  # 레이스 감지 모드 (가장 중요!)
  go test -race -v

  # 레이스가 발생하는 테스트만 실행
  go test -race -v -run TestUnsafe

  # 안전한 구현 테스트
  go test -race -v -run "TestMutex|TestAtomic|TestSafe"

  # 동기화 방식별 성능 비교
  go test -bench=BenchmarkCounter -benchmem

  레이스 컨디션 감지 결과 예시:

  ==================
  WARNING: DATA RACE
  Read at 0x00c0000b4010 by goroutine 8:
    testexample.(*UnsafeCounter).Increment()
        race_condition.go:35 +0x3c

  Previous write at 0x00c0000b4010 by goroutine 7:
    testexample.(*UnsafeCounter).Increment()
        race_condition.go:35 +0x52
  ==================
