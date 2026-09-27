
## 동작 방식

```bash
kubectl edit pod aaa
# vi 에디터 열림
# 수정 후 :wq 저장
# 즉시 적용
```

---

## 근데 pod은 수정 안되는 필드가 있어요

```bash
kubectl edit pod aaa
# image, port, resource 등 수정하면
# 저장할 때 에러남

# "... field is immutable" 에러
```

이럴 때 vi가 그냥 닫히는게 아니라

```bash
# 수정 내용을 임시파일로 저장해줌
# /tmp/kubectl-edit-xxxx.yaml
# 이걸 이용해서

kubectl replace --force -f /tmp/kubectl-edit-xxxx.yaml
```

---

## 필드별 정리

```
edit으로 바로 되는 것
  replicas 수
  이미지 (deployment에서)
  환경변수 (deployment에서)
  label, annotation

edit으로 안되는 것 (pod 직접)
  이미지
  포트
  리소스 requests/limits
  → replace --force 필요
```

---

## 그래서 실전 패턴

```
deployment 수정  → kubectl edit 바로
pod 직접 수정    → replace --force
빠른 이미지 변경 → kubectl set image
```

시험에서 pod 직접 수정하라고 하면 거의 `replace --force` 쓴다고 보면 돼요.