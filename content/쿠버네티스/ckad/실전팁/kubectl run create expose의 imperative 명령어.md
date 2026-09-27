CKAD에서 필수로 외워야 할 imperative 명령어들

**Pod**

```bash
# Pod 생성
kubectl run mypod --image=nginx

# Pod 생성 + 포트 노출
kubectl run mypod --image=nginx --port=80

# yaml만 뽑기 (실제 생성 X)
kubectl run mypod --image=nginx --dry-run=client -o yaml > pod.yaml
```

**Deployment**

```bash
# Deployment 생성
kubectl create deployment mydeploy --image=nginx

# 레플리카 지정
kubectl create deployment mydeploy --image=nginx --replicas=3

# 이미지 업데이트
kubectl set image deployment/mydeploy nginx=nginx:1.21
```

**Service (expose)**

```bash
# ClusterIP (기본)
kubectl expose deployment mydeploy --port=80 --target-port=8080

# NodePort
kubectl expose deployment mydeploy --port=80 --type=NodePort

# LoadBalancer
kubectl expose deployment mydeploy --port=80 --type=LoadBalancer
```

**ConfigMap / Secret**

```bash
kubectl create configmap myconfig --from-literal=key=value
kubectl create secret generic mysecret --from-literal=password=1234
```

**Job / CronJob**

```bash
kubectl create job myjob --image=busybox -- echo "hello"
kubectl create cronjob mycron --image=busybox --schedule="*/1 * * * *" -- echo "hello"
```

**기타 필수**

```bash
# 스케일
kubectl scale deployment mydeploy --replicas=5

# 라벨 확인
kubectl get pods --show-labels

# 네임스페이스 지정
kubectl get pods -n mynamespace

# 강제 삭제
kubectl delete pod mypod --force --grace-period=0
```

---

**시험 꿀팁 2개**

`--dry-run=client -o yaml` 조합이 핵심이에요. yaml 처음부터 안 쓰고 명령어로 뼈대 뽑아서 수정하는 게 훨씬 빠릅니다.

```bash
kubectl create deployment mydeploy --image=nginx --dry-run=client -o yaml > deploy.yaml
vi deploy.yaml  # 필요한 부분만 수정
kubectl apply -f deploy.yaml
```

그리고 시험 시작하자마자 alias 설정해두세요:

```bash
alias k=kubectl
export do="--dry-run=client -o yaml"
```

그럼 `k run mypod --image=nginx $do > pod.yaml` 이렇게 단축해서 쓸 수 있어요!