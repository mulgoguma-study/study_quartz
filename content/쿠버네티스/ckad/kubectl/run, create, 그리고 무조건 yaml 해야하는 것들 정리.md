---
aliases:
  - run/create 차이점
---
좋은 질문이에요. 이거 정리해두면 CKAD 시험에서 어떤 명령어를 쓸지 고민하는 시간이 확 줄어듭니다.

## 한 줄 요약

| 명령어              | 만들 수 있는 것                                                   |
| ---------------- | ----------------------------------------------------------- |
| `kubectl run`    | **Pod만**                                                    |
| `kubectl create` | Pod 빼고 거의 다 (Deployment, Service, ConfigMap, Secret, Job 등) |
| `무조건 yaml`       | NetworkPolicy, PersistentVolume, PersistentVolumeClaim, CRD |

`kubectl run`은 Pod 전용이고, `kubectl create`는 다양한 리소스를 만드는 범용 명령어예요.

## `kubectl run` — Pod 생성 전용

```bash
kubectl run nginx --image=nginx --dry-run=client -o yaml
```

**주요 옵션:**

```bash
kubectl run <이름> \
  --image=<이미지> \
  --env='KEY=VALUE' \
  --labels='app=web' \
  --port=80 \
  --command -- <명령어와 인자> \
  --dry-run=client -o yaml
```

**제한사항**: Pod만 만들 수 있어요. Deployment나 Service는 못 만듭니다. (과거엔 됐는데 지금은 deprecated 돼서 막혔어요.)

## `kubectl create` — 다양한 리소스 생성

하위 명령어(subcommand)로 리소스 타입을 지정합니다.

```bash
# Deployment
kubectl create deployment nginx --image=nginx --replicas=3 --dry-run=client -o yaml

# Service
kubectl create service clusterip myapp --tcp=80:8080 --dry-run=client -o yaml
kubectl expose deployment nginx --port=80 --dry-run=client -o yaml  # 이것도 create 계열

# ConfigMap
kubectl create configmap myconfig --from-literal=key1=value1 --dry-run=client -o yaml
kubectl create configmap myconfig --from-file=./config.txt --dry-run=client -o yaml

# Secret
kubectl create secret generic mysecret --from-literal=password=s3cret --dry-run=client -o yaml

# Job
kubectl create job myjob --image=busybox --dry-run=client -o yaml -- sleep 60

# CronJob
kubectl create cronjob mycron --image=busybox --schedule="*/5 * * * *" --dry-run=client -o yaml -- date

# ServiceAccount
kubectl create serviceaccount mysa --dry-run=client -o yaml

# Role / RoleBinding
kubectl create role myrole --verb=get,list --resource=pods --dry-run=client -o yaml

# Namespace
kubectl create namespace dev --dry-run=client -o yaml
```

## CKAD 시험에서 외워야 할 패턴

시험에서 가장 자주 쓰는 것만 추렸어요:

```bash
# ──── Pod ────
kubectl run <name> --image=<img> --dry-run=client -o yaml > pod.yaml

# ──── Deployment ────
kubectl create deployment <name> --image=<img> --replicas=3 --dry-run=client -o yaml > deploy.yaml

# ──── Service (Deployment 노출용) ────
kubectl expose deployment <name> --port=80 --target-port=8080 --dry-run=client -o yaml > svc.yaml

# ──── ConfigMap ────
kubectl create cm <name> --from-literal=KEY=VAL --dry-run=client -o yaml > cm.yaml

# ──── Secret ────
kubectl create secret generic <name> --from-literal=KEY=VAL --dry-run=client -o yaml > secret.yaml

# ──── Job ────
kubectl create job <name> --image=<img> --dry-run=client -o yaml -- <command> > job.yaml
```

## 실전 워크플로우

CKAD 시험에서의 정석 흐름이에요:

```bash
# 1. 기본 YAML 생성 (골격만)
kubectl create deployment web --image=nginx --replicas=3 --dry-run=client -o yaml > web.yaml

# 2. vim으로 열어서 추가 설정 넣기 (resources, volumes, env 등)
vim web.yaml

# 3. 문법 검증
kubectl apply -f web.yaml --dry-run=client

# 4. 진짜 적용
kubectl apply -f web.yaml
```

**"명령어로 뼈대 만들고 → YAML 편집으로 살 붙이기"**가 핵심입니다. 처음부터 YAML을 타이핑하면 시간도 오래 걸리고 오타도 납니다.

## Go 개발자 관점의 비유

```go
// kubectl run — Pod 전용 생성자
func NewPod(name, image string) *Pod { ... }

// kubectl create — 범용 팩토리
func Create(kind string, opts ...Option) Resource {
    switch kind {
    case "deployment":
        return NewDeployment(...)
    case "service":
        return NewService(...)
    case "configmap":
        return NewConfigMap(...)
    // ...
    }
}
```

`run`은 한 가지만 만드는 특화된 생성자, `create`는 종류별로 서브커맨드 가진 범용 팩토리라고 생각하시면 돼요.

## 헷갈리는 포인트 하나

**Q: "왜 Deployment는 `kubectl create deployment`고, Pod은 `kubectl create pod`이 아니라 `kubectl run`이야?"**

**A:** 역사적 이유예요. 원래 `kubectl run`이 Pod/Deployment/Job 등을 다 만들 수 있었는데, 혼란스러워서 "Pod만 만드는 명령"으로 단순화됐어요. 그래서 이름이 불균일합니다. 시험에서는 그냥 외우는 게 빨라요:

- **Pod** → `kubectl run`
- **그 외 전부** → `kubectl create <타입>`

## 치트시트로 저장해두기

```bash
# Pod
k run NAME --image=IMG [--command -- cmd args] --dry-run=client -o yaml

# Deployment  
k create deploy NAME --image=IMG --replicas=N --dry-run=client -o yaml

# Service
k expose [pod|deploy] NAME --port=P --target-port=TP --dry-run=client -o yaml

# Job
k create job NAME --image=IMG --dry-run=client -o yaml -- cmd args

# CronJob
k create cj NAME --image=IMG --schedule="* * * * *" --dry-run=client -o yaml

# ConfigMap / Secret
k create cm NAME --from-literal=K=V --dry-run=client -o yaml
k create secret generic NAME --from-literal=K=V --dry-run=client -o yaml
```

참고로 시험에서는 `alias k=kubectl`이 기본으로 설정돼 있어서 `k`로 쓰는 게 시간 절약됩니다.