``` bash
# GOPATH 설정 확인 (기본: ~/go)
go env GOPATH  # 없으면 export GOPATH=$HOME/go

# DLV 설치 (최신 버전)
go install github.com/go-delve/delve/cmd/dlv@latest

# PATH에 GOPATH/bin 추가 (이미 안 되어 있으면)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc

# 설치 확인
dlv version  # Delve Debugger Version: 1.25.x 출력 확인



# 업데이트 시 (그냥 install 재실행)
go install github.com/go-delve/delve/cmd/dlv@latest


```