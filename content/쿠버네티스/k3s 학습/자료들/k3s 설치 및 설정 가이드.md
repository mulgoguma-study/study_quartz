# K3s MSA 학습 환경 구축 가이드

## 1. K3s 클러스터 구성

### 마스터 노드 (개발 서버) 설치

```bash
# K3s 마스터 노드 설치
curl -sfL https://get.k3s.io | sh -

# 노드 토큰 확인 (워커 노드 연결용)
sudo cat /var/lib/rancher/k3s/server/node-token

# kubeconfig 권한 설정
sudo chmod 644 /etc/rancher/k3s/k3s.yaml
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

# 클러스터 상태 확인
kubectl get nodes
```

### 워커 노드 설치

```bash
# 워커 노드에서 실행 (마스터 IP와 토큰 입력 필요)
curl -sfL https://get.k3s.io | K3S_URL=https://<MASTER_IP>:6443 \
  K3S_TOKEN=<NODE_TOKEN> sh -
```

## 2. 로컬 레지스트리 구성 (Podman 활용)

### Podman으로 프라이빗 레지스트리 실행

```bash
# 레지스트리 컨테이너 실행
podman run -d \
  --name registry \
  -p 5000:5000 \
  --restart=always \
  -v registry-data:/var/lib/registry \
  registry:2

# K3s에서 insecure 레지스트리 허용 설정
sudo mkdir -p /etc/rancher/k3s
sudo tee /etc/rancher/k3s/registries.yaml > /dev/null <<EOF
mirrors:
  "localhost:5000":
    endpoint:
      - "http://localhost:5000"
configs:
  "localhost:5000":
    tls:
      insecure_skip_verify: true
EOF

# K3s 재시작
sudo systemctl restart k3s
```

## 3. 이미지 빌드 및 푸시

### 프로젝트 구조

```
msa-project/
├── gateway/
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── user-service/
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── product-service/
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
└── k8s/
    ├── manifests.yaml
    └── argocd/
```

### 빌드 스크립트

```bash
#!/bin/bash
# build-and-push.sh

REGISTRY="localhost:5000"

# Gateway 빌드
cd gateway
podman build -t ${REGISTRY}/gateway:latest .
podman push ${REGISTRY}/gateway:latest

# User Service 빌드
cd ../user-service
podman build -t ${REGISTRY}/user-service:latest .
podman push ${REGISTRY}/user-service:latest

# Product Service 빌드
cd ../product-service
podman build -t ${REGISTRY}/product-service:latest .
podman push ${REGISTRY}/product-service:latest

cd ..
echo "All images built and pushed successfully!"
```

## 4. 배포 및 테스트

### 애플리케이션 배포

```bash
# 매니페스트 적용
kubectl apply -f k8s/manifests.yaml

# 배포 상태 확인
kubectl get pods -o wide
kubectl get svc
kubectl get hpa

# Pod 분산 확인 (어느 노드에서 실행 중인지)
kubectl get pods -o wide | grep -E 'NAME|gateway|user|product'
```

### 서비스 테스트

```bash
# Gateway를 통한 접근 테스트
MASTER_IP=$(kubectl get nodes -o wide | grep master | awk '{print $6}')
curl http://${MASTER_IP}:30080/
curl http://${MASTER_IP}:30080/users
curl http://${MASTER_IP}:30080/products

# 여러 번 호출하여 로드밸런싱 확인 (hostname이 바뀌는지 확인)
for i in {1..10}; do
  curl -s http://${MASTER_IP}:30080/users | jq '.hostname'
done
```

## 5. Scale Out 테스트

### 수동 스케일링

```bash
# Replica 수 변경
kubectl scale deployment gateway --replicas=5
kubectl scale deployment user-service --replicas=3

# 스케일링 확인
kubectl get pods -l app=gateway
```

### HPA (Horizontal Pod Autoscaler) 테스트

```bash
# Metrics Server 확인 (K3s에는 기본 포함)
kubectl top nodes
kubectl top pods

# 부하 생성 (별도 터미널에서)
while true; do curl http://${MASTER_IP}:30080/users; done

# HPA 상태 모니터링
watch kubectl get hpa

# Pod 자동 증가 확인
watch kubectl get pods -l app=gateway
```

## 6. ArgoCD 설치 및 CI/CD 구성

### ArgoCD 설치

```bash
# ArgoCD 네임스페이스 생성
kubectl create namespace argocd

# ArgoCD 설치
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# ArgoCD Server를 NodePort로 노출
kubectl patch svc argocd-server -n argocd -p '{"spec":{"type":"NodePort"}}'

# NodePort 확인
kubectl get svc argocd-server -n argocd

# 초기 admin 패스워드 확인
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

### Forgejo (Gitea Fork) 설치

```bash
# Forgejo Helm 차트 설치
helm repo add forgejo https://dl.gitea.io/charts/
kubectl create namespace forgejo

helm install forgejo forgejo/forgejo \
  --namespace forgejo \
  --set service.http.type=NodePort \
  --set service.http.nodePort=30300

# Forgejo 접속 정보
kubectl get svc -n forgejo
```

### ArgoCD Application 생성

```yaml
# argocd-app.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: msa-app
  namespace: argocd
spec:
  project: default
  source:
    repoURL: http://<FORGEJO_IP>:30300/user/msa-project.git
    targetRevision: main
    path: k8s
  destination:
    server: https://kubernetes.default.svc
    namespace: default
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

### CI/CD 워크플로우

```yaml
# .forgejo/workflows/build.yaml
name: Build and Deploy
on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Build and Push Gateway
        run: |
          podman build -t localhost:5000/gateway:${{ github.sha }} ./gateway
          podman push localhost:5000/gateway:${{ github.sha }}
      
      - name: Update K8s Manifests
        run: |
          sed -i "s|image: localhost:5000/gateway:.*|image: localhost:5000/gateway:${{ github.sha }}|" k8s/manifests.yaml
          git config user.name "CI Bot"
          git config user.email "ci@example.com"
          git add k8s/manifests.yaml
          git commit -m "Update image to ${{ github.sha }}"
          git push
```

## 7. 모니터링 설정

### 간단한 모니터링

```bash
# 실시간 리소스 사용량 확인
watch kubectl top pods
watch kubectl top nodes

# 로그 확인
kubectl logs -f deployment/gateway
kubectl logs -f deployment/user-service

# 이벤트 확인
kubectl get events --sort-by='.lastTimestamp'
```

## 8. 유용한 명령어 모음

```bash
# 전체 Pod 상태 확인
kubectl get pods -A -o wide

# 특정 Pod 내부 접속
kubectl exec -it <pod-name> -- sh

# Service Endpoint 확인
kubectl get endpoints

# ConfigMap 생성 예제
kubectl create configmap app-config --from-literal=LOG_LEVEL=debug

# Secret 생성 예제
kubectl create secret generic db-secret --from-literal=password=mypassword

# Pod 삭제 후 재생성 (롤링 업데이트)
kubectl rollout restart deployment/gateway

# 배포 히스토리 확인
kubectl rollout history deployment/gateway

# 이전 버전으로 롤백
kubectl rollout undo deployment/gateway
```

## 9. 트러블슈팅

### Pod가 시작되지 않을 때

```bash
# Pod 상세 정보 확인
kubectl describe pod <pod-name>

# 이벤트 확인
kubectl get events --field-selector involvedObject.name=<pod-name>

# 이미지 Pull 실패 시
kubectl describe pod <pod-name> | grep -A 10 Events
```

### 네트워크 연결 문제

```bash
# 클러스터 내부에서 DNS 테스트
kubectl run test-pod --image=busybox --rm -it -- sh
# Pod 내부에서
nslookup user-service
wget -O- http://user-service:8081/health
```

## 10. 정리

```bash
# 리소스 삭제
kubectl delete -f k8s/manifests.yaml

# ArgoCD 삭제
kubectl delete namespace argocd

# K3s 완전 제거 (마스터)
/usr/local/bin/k3s-uninstall.sh

# K3s 완전 제거 (워커)
/usr/local/bin/k3s-agent-uninstall.sh
```