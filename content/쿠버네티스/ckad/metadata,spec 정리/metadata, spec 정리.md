네, `metadata`는 거의 모든 k8s 리소스 공통이고, `spec`은 리소스마다 다릅니다.

![[Pasted image 20260428172502.png]]

핵심 포인트 몇 가지 보충하면:

`metadata`는 Pod/Deployment/Service/Job 가릴 것 없이 구조가 동일합니다. `labels`과 `annotations`의 차이를 헷갈리면 안 되는데, selector로 Pod를 고를 때 쓰는 건 `labels`이고, Prometheus 스크랩 설정처럼 운영 메타데이터는 `annotations`입니다.

`spec`은 리소스마다 완전히 다르지만, Pod spec이 핵심입니다. Deployment의 `spec.template`이 결국 Pod spec이라 Pod를 잘 알면 나머지는 거기에 레이어가 얹히는 구조입니다.

CKAD에서 가장 많이 나오는 패턴은 `env` (ConfigMap/Secret에서 값 주입), `volumeMounts` + `volumes` 연결, `resources.requests/limits`, probe 설정 이 네 가지입니다.