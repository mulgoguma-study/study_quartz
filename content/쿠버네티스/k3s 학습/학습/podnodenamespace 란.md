---
aliases:
  - pod/node/namespace 란?
---
쿠버네티스의 핵심 세 가지 개념을 설명해드릴게요.

---

## Node

클러스터를 구성하는 **물리/가상 서버**입니다. 실제 컴퓨팅 자원(CPU, Memory)을 제공하는 인프라 단위예요.

- **Control Plane Node** (구 Master): API Server, Scheduler, etcd 등 클러스터 제어 컴포넌트 실행
- **Worker Node**: 실제 애플리케이션 워크로드(Pod)가 실행되는 곳
- 각 Worker Node에는 `kubelet`, `kube-proxy`, container runtime(containerd 등)이 실행됨

---

## Pod

쿠버네티스에서 **배포 가능한 가장 작은 단위**입니다. 컨테이너를 직접 다루지 않고 Pod로 감싸서 관리해요.

- 1개 이상의 컨테이너를 포함 (보통은 1개, 사이드카 패턴이면 2개+)
- Pod 내 컨테이너들은 **localhost로 통신** 가능하고 **볼륨을 공유**
- Pod는 일회성 — 죽으면 새로운 IP를 가진 새 Pod가 생성됨
- 그래서 직접 Pod를 쓰기보단 `Deployment`, `StatefulSet` 같은 상위 리소스로 관리

```
Node
└── Pod
    ├── container A (예: Go 앱)
    └── container B (예: envoy sidecar)
```

---

## Namespace

클러스터 내의 **논리적 격리 단위**입니다. 하나의 클러스터를 여러 팀/환경으로 나눌 때 사용해요.

- 리소스 이름은 같은 Namespace 안에서만 유일하면 됨 → `dev/my-app`, `prod/my-app` 공존 가능
- **RBAC**, **NetworkPolicy**, **ResourceQuota** 등을 Namespace 단위로 적용
- 기본 제공 Namespace: `default`, `kube-system`, `kube-public`, `kube-node-lease`
- Node는 Namespace에 속하지 않음 (cluster-scoped 리소스)

---

## 세 개념의 관계 요약

|개념|스코프|역할|
|---|---|---|
|Node|클러스터 전체|컴퓨팅 자원 제공|
|Namespace|논리적 파티션|리소스 격리/구분|
|Pod|Namespace 내|컨테이너 실행 단위|

Go 백엔드 서비스를 배포한다면 보통 `Deployment` → `Pod` 흐름으로 가고, Namespace를 `dev` / `staging` / `prod` 로 나눠서 운영하는 패턴을 많이 씁니다.