``` toml
[Unit]
Description=Go Bye-PBX Main Service Process
After=network.target ai-agent-proc-db.service ai-agent-proc-gpt.service ai-agent-tts-proc.service
Requires=ai-agent-proc-db.service ai-agent-proc-gpt.service ai-agent-tts-proc.service

[Service]
# 서비스가 비정상 종료 시 항상 재시작하도록 설정
Restart=always
RestartSec=3

# Go Bin location
ExecStart=/home/hanil/go/src/bye_pbx/bin/web_server

# 프로젝트 루트 경로를 작업 디렉토리로 설정
# 이 경로는 Go 앱이 설정 파일 등을 상대 경로로 읽을 때 기준
WorkingDirectory=/home/hanil/go/src/bye_pbx

User=hanil

# 프로세스 종료 시 kill 신호를 사용합니다.
KillMode=process

# 프로세스가 완전히 시작되었다고 간주될 때까지 대기하는 시간 (선택 사항)
# Type=simple은 ExecStart 명령어가 메인 프로세스임을 systemd에 알립니다.
Type=simple

[Install]
WantedBy=multi-user.target
```


- After, Require 에 등록된 프로세스들을 먼저 수행하고 해당 메인이 수행되도록 설정한 케이스
- 