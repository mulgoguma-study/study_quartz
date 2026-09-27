headless 모드로 실행해야됨.
해당 모드로 dlv는 서버에서만 돌고, Goland가 원격으로 연결됨

[예제]
**예시: 간단한 Go 프로그램 디버깅**

- 프로젝트 디렉토리로 이동: cd /path/to/your/go/project
- 메인 파일 빌드: go build -o myapp main.go (디버깅 심볼 포함 위해 -gcflags=all=-N -l 플래그 추가 가능: go build -gcflags=all=-N -l -o myapp main.go)

``` bash
# 새 프로그램 디버그 (빌드 없이 소스에서 바로)
dlv debug main.go -- --arg1 value1  # main.go를 디버깅, 인자 추가 가능

# 또는 기존 바이너리 실행 (exec 모드)
dlv exec --headless --listen=:2345 --api-version=2 ./myapp --arg1 value1

# 전체 옵션 설명:
# --headless: GUI 없이 서버 모드
# --listen=:2345: 포트 2345로 리스닝 (변경 가능)
# --api-version=2: DLV API 버전 (GoLand 호환 위해 필수)
# exec/debug 뒤에 프로그램 경로나 파일
```

-  프로그램 시작하면 dlv 대기 상태. 
-  --log 옵션 시 /tmp/dlv.log 남김


