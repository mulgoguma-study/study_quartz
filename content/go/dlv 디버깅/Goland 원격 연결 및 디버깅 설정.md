# 1. Goland ssh 연결 후 프로젝트 열기
이건 머 아는거

# 2. 디버깅 구성
### 1. 상단 메뉴 : Run > Edit Configurations
### 2. Go Remote 선택
### 3. 설정
- **Name**: "Remote DLV Debug" (임의)
- **Host**: 서버 IP (e.g., 192.168.1.100)
- **Port**: 2345 (DLV listen 포트)
- **Remote path mappings**: 로컬 경로와 서버 경로 매핑 (e.g., Local: /Users/me/project → Remote: /home/user/project). 자동으로 설정됨.
- **Go tool args**: 필요 시 (e.g., -gcflags=all=-N -l for 디버그 심볼).
### 4. 디버깅 실행
1. 코드에서 브레이크 포인트 설정
2. 서버 터미널에서 dlv 실행 
3. goland에서 debug 버튼 클릭
	1. 연결 성공: "Connected to Delve" 메시지.
		- 프로그램 실행: F9 (Resume), F8 (Step Over), F7 (Step Into) 등 사용.
		- 변수/스택 확인: Debug 패널에서




#### 추가 팁 & 문제 해결

- **Docker 사용 시**: Dockerfile에 DLV 설치하고, EXPOSE 2345 추가. DLV 실행 시 컨테이너 포트 매핑 (e.g., -p 2345:2345).
- **성능 이슈**: 원격 연결 시 지연 발생 가능. GoLand 설정 > Build > Go > Use Go modules 자동으로 켜기.
- **Attach to running process**: 이미 실행 중인 프로세스에 붙이려면 서버에서 dlv attach --headless --listen=:2345 --api-version=2 <PID> (PID: ps aux | grep myapp).
- **에러 흔함**:
    - "Connection refused": 포트 열기 확인.
    - "No debug info": 빌드 시 -gcflags=all=-N -l 추가.
    - DLV 버전 mismatch: 서버/GoLand DLV 버전 맞추기 (go install @latest).

이 방법으로 하면 로컬처럼 편하게 디버깅할 수 있어요. 더 구체적인 에러가 있으면 로그 공유해주세요! 😊