```sh
.PHONY: help build push deploy clean test scale-up scale-down logs status

# 변수 설정
REGISTRY ?= localhost:5000
VERSION ?= v1.0
MASTER_IP ?= $(shell hostname -I | awk '{print $$1}')
SERVICES = gateway user-service product-service

help: ## 도움말 표시
	@echo "K3s MSA 실습 환경 관리 명령어"
	@echo "=============================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## 모든 서비스 이미지 빌드
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd $$service && podman build -t $(REGISTRY)/$$service:$(VERSION) . && cd ..; \
	done
	@echo "✓ All images built successfully"

push: ## 이미지를 레지스트리에 푸시
	@echo "Pushing images to registry..."
	@for service in $(SERVICES); do \
		echo "Pushing $$service..."; \
		podman push $(REGISTRY)/$$service:$(VERSION); \
	done
	@echo "✓ All images pushed successfully"

build-push: build push ## 빌드 및 푸시를 한번에 실행

deploy: ## K3s에 애플리케이션 배포
	@echo "Deploying to K3s..."
	@sed 's|image: localhost:5000/\([^:]*\):.*|image: $(REGISTRY)/\1:$(VERSION)|g' k8s/manifests.yaml | kubectl apply -f -
	@echo "✓ Deployment completed"
	@echo "Waiting for pods to be ready..."
	@kubectl wait --for=condition=ready pod -l app=gateway --timeout=60s
	@kubectl wait --for=condition=ready pod -l app=user-service --timeout=60s
	@kubectl wait --for=condition=ready pod -l app=product-service --timeout=60s

update: build-push deploy ## 전체 업데이트 (빌드 + 푸시 + 배포)

clean: ## 모든 리소스 삭제
	@echo "Cleaning up K3s resources..."
	@kubectl delete -f k8s/manifests.yaml --ignore-not-found=true
	@echo "✓ Resources cleaned"

test: ## 서비스 테스트
	@echo "Testing services..."
	@echo "Gateway root:"
	@curl -s http://$(MASTER_IP):30080/ | jq .
	@echo "\nUser Service:"
	@curl -s http://$(MASTER_IP):30080/users | jq '.service, .hostname'
	@echo "\nProduct Service:"
	@curl -s http://$(MASTER_IP):30080/products | jq '.service, .hostname'

test-lb: ## 로드밸런싱 테스트 (10회 요청)
	@echo "Testing load balancing (10 requests)..."
	@for i in {1..10}; do \
		echo -n "Request $$i: "; \
		curl -s http://$(MASTER_IP):30080/users | jq -r '.hostname'; \
	done

scale-up: ## Gateway를 5개로 스케일 업
	@echo "Scaling up gateway to 5 replicas..."
	@kubectl scale deployment gateway --replicas=5
	@kubectl get pods -l app=gateway

scale-down: ## Gateway를 2개로 스케일 다운
	@echo "Scaling down gateway to 2 replicas..."
	@kubectl scale deployment gateway --replicas=2
	@kubectl get pods -l app=gateway

logs: ## Gateway 로그 확인
	@kubectl logs -f deployment/gateway --tail=50

logs-user: ## User Service 로그 확인
	@kubectl logs -f deployment/user-service --tail=50

logs-product: ## Product Service 로그 확인
	@kubectl logs -f deployment/product-service --tail=50

status: ## 전체 상태 확인
	@echo "=== Nodes ==="
	@kubectl get nodes
	@echo "\n=== Pods ==="
	@kubectl get pods -o wide
	@echo "\n=== Services ==="
	@kubectl get svc
	@echo "\n=== HPA ==="
	@kubectl get hpa
	@echo "\n=== Resource Usage ==="
	@kubectl top nodes 2>/dev/null || echo "Metrics not available"
	@kubectl top pods 2>/dev/null || echo "Metrics not available"

watch-pods: ## Pod 상태 실시간 모니터링
	@watch -n 2 'kubectl get pods -o wide'

watch-hpa: ## HPA 상태 실시간 모니터링
	@watch -n 2 'kubectl get hpa'

restart: ## 모든 Deployment 재시작
	@echo "Restarting all deployments..."
	@kubectl rollout restart deployment/gateway
	@kubectl rollout restart deployment/user-service
	@kubectl rollout restart deployment/product-service
	@echo "✓ All deployments restarted"

rollback: ## Gateway를 이전 버전으로 롤백
	@echo "Rolling back gateway..."
	@kubectl rollout undo deployment/gateway
	@kubectl rollout status deployment/gateway

describe-gateway: ## Gateway Deployment 상세 정보
	@kubectl describe deployment gateway

describe-hpa: ## HPA 상세 정보
	@kubectl describe hpa

events: ## 최근 이벤트 확인
	@kubectl get events --sort-by='.lastTimestamp' | tail -20

shell-gateway: ## Gateway Pod 내부 접속
	@kubectl exec -it deployment/gateway -- sh

registry-check: ## 레지스트리 이미지 목록 확인
	@echo "Images in registry:"
	@curl -s http://$(REGISTRY)/v2/_catalog | jq .
	@for service in $(SERVICES); do \
		echo "\nTags for $$service:"; \
		curl -s http://$(REGISTRY)/v2/$$service/tags/list | jq .; \
	done

install-k3s: ## K3s 마스터 노드 설치
	@echo "Installing K3s master node..."
	@curl -sfL https://get.k3s.io | sh -
	@sudo systemctl enable k3s
	@echo "✓ K3s installed"
	@echo "Node token:"
	@sudo cat /var/lib/rancher/k3s/server/node-token

setup-registry: ## 로컬 레지스트리 구성
	@echo "Setting up local registry..."
	@podman run -d --name registry -p 5000:5000 --restart=always docker.io/library/registry:2
	@sudo mkdir -p /etc/rancher/k3s
	@echo 'mirrors:\n  "localhost:5000":\n    endpoint:\n      - "http://$(MASTER_IP):5000"\nconfigs:\n  "localhost:5000":\n    tls:\n      insecure_skip_verify: true' | sudo tee /etc/rancher/k3s/registries.yaml
	@sudo systemctl restart k3s
	@echo "✓ Registry configured"

init: install-k3s setup-registry ## 초기 환경 구성 (K3s + Registry)

full-cycle: clean build-push deploy test ## 전체 사이클 실행 (정리 + 빌드 + 배포 + 테스트)
```