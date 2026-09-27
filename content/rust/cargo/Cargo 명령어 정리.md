<% tp.file.creation_date("YYYY-MM-dd") %>



## <font color="#0070c0">기본 프로젝트 관리</font>

### cargo new

```
cargo new [프로젝트명]
```

**설명**: 새로운 Rust 프로젝트 생성 **옵션**:

- `--bin`: 바이너리 프로젝트 (기본값)
- `--lib`: 라이브러리 프로젝트
- `--vcs [git/hg/pijul/none]`: 버전 관리 시스템 선택

**사용 예시**:

```
cargo new my_app --bin
cargo new my_lib --lib
```

---

### cargo init

```
cargo init [디렉토리명]
```

**설명**: 기존 디렉토리를 Cargo 프로젝트로 초기화 **사용 예시**:

```
cargo init ./my_project
```

---

## <font color="#0070c0">빌드 및 실행</font>

### cargo build

```
cargo build [옵션]
```

**설명**: 프로젝트 빌드 (디버그 모드) **옵션**:

- `--release`: 최적화된 릴리스 빌드
- `--target [triple]`: 특정 대상 플랫폼 빌드
- `--quiet`: 출력 최소화

**사용 예시**:

```
cargo build
cargo build --release
```

---

### cargo run

```
cargo run [옵션] [-- 프로그램인자]
```

**설명**: 프로젝트 빌드 및 실행 **옵션**:

- `--release`: 최적화 빌드 후 실행
- `--quiet`: 출력 최소화

**사용 예시**:

```
cargo run
cargo run --release
cargo run -- arg1 arg2
```

---

### cargo check

```
cargo check [옵션]
```

**설명**: 컴파일 오류 확인 (빌드 아티팩트 생성 안 함) **용도**: 빠른 문법 검사 **사용 예시**:

```
cargo check
```

---

## <font color="#0070c0">패키지 관리</font>

### cargo add

```
cargo add [패키지명]
```

**설명**: 의존성 추가 **옵션**:

- `--version [버전]`: 특정 버전 지정
- `--git [URL]`: Git 저장소에서 추가
- `--dev`: 개발 의존성으로 추가
- `--optional`: 선택적 의존성으로 추가

**사용 예시**:

```
cargo add serde
cargo add serde --version 1.0
cargo add tokio --dev
```

---

### cargo remove

```
cargo remove [패키지명]
```

**설명**: 의존성 제거 **사용 예시**:

```
cargo remove serde
```

---

### cargo update

```
cargo update [옵션]
```

**설명**: 의존성 업데이트 (Cargo.lock 갱신) **옵션**:

- `-p [패키지명]`: 특정 패키지만 업데이트
- `--aggressive`: 주요 버전도 업데이트

**사용 예시**:

```
cargo update
cargo update -p serde
```

---

## 테스트 및 문서

### cargo test

```
cargo test [옵션] [테스트명]
```

**설명**: 테스트 실행 **옵션**:

- `--lib`: 라이브러리 테스트만 실행
- `--doc`: 문서 테스트만 실행
- `--release`: 최적화 빌드로 테스트
- `--quiet`: 최소 출력
- `--`: 다음 인자를 테스트에 전달

**사용 예시**:

```
cargo test
cargo test my_test_function
cargo test --doc
```

---

### cargo doc

```
cargo doc [옵션]
```

**설명**: API 문서 생성 **옵션**:

- `--open`: 생성 후 브라우저에서 열기
- `--no-deps`: 의존성 문서 제외

**사용 예시**:

```
cargo doc --open
```

---

## <font color="#0070c0">품질 및 분석</font>

### cargo clippy

```
cargo clippy [옵션]
```

**설명**: Lint 도구로 코드 품질 검사 **설치**: `rustup component add clippy` **사용 예시**:

```
cargo clippy
cargo clippy -- -D warnings
```

---

### cargo fmt

```
cargo fmt [옵션]
```

**설명**: 코드 포매팅 (Rust 스타일 가이드 준수) **설치**: `rustup component add rustfmt` **옵션**:

- `--check`: 변경 사항만 표시 (파일 수정 안 함)

**사용 예시**:

```
cargo fmt
cargo fmt -- --check
```

---

### cargo audit

```
cargo audit
```

**설명**: 보안 취약점 검사 **설치**: `cargo install cargo-audit` **사용 예시**:

```
cargo audit
```

---

## <font color="#0070c0">배포 및 공유</font>

### cargo publish

```
cargo publish [옵션]
```

**설명**: crates.io에 패키지 배포 **전제조건**: crates.io 계정 필요 **옵션**:

- `--dry-run`: 배포 시뮬레이션
- `--allow-dirty`: 미커밋 변경사항 무시

**사용 예시**:

```
cargo publish --dry-run
cargo publish
```

---

### cargo search

```
cargo search [검색어]
```

**설명**: crates.io에서 패키지 검색 **사용 예시**:

```
cargo search serde
```

---

## <font color="#0070c0">유용한 팁</font>

### 자주 사용하는 명령어 조합

|목적|명령어|
|---|---|
|빠른 개발|`cargo check`|
|완전한 빌드 테스트|`cargo test --release`|
|배포 전 검증|`cargo test && cargo clippy && cargo fmt --check`|
|실행 및 테스트|`cargo run && cargo test`|

### <font color="#0070c0">Cargo.toml 주요 섹션</font>

```toml
[package]
name = "my_project"
version = "0.1.0"
edition = "2021"

[dependencies]
serde = "1.0"

[dev-dependencies]
criterion = "0.5"

[profile.release]
opt-level = 3
lto = true
```

---

