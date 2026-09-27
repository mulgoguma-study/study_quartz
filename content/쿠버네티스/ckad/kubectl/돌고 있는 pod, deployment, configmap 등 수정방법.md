
## 상황별 패턴

**이미 실행중인 pod 수정**

```bash
# 방법 1 - yaml 떨구고 수정
kubectl get pod aaa -o yaml > pod.yaml
vi pod.yaml
kubectl delete pod aaa
kubectl apply -f pod.yaml  # pod은 replace 안되는 경우 많음

# 방법 2 - edit으로 바로 수정
kubectl edit pod aaa
# 저장하면 자동 적용 (일부 필드는 수정 불가)
```

**deployment 수정**

```bash
# edit이 제일 편함
kubectl edit deployment aaa

# 이미지만 바꿀 때
kubectl set image deployment/aaa container명=nginx:1.19
```

**service 수정**

```bash
kubectl edit svc aaa
# 포트 변경 등
```

---

## pod vs deployment 차이

```
pod
  일부 필드 수정 불가 (이미지, 포트 등)
  kubectl edit 하면 에러남
  → delete 하고 다시 apply

deployment
  대부분 수정 가능
  kubectl edit 하면 자동으로 pod 재시작
  → edit이 편함
```

---

## 시험에서 실제 흐름

```bash
# 1. 현재 상태 확인
kubectl get pod aaa -o yaml > aaa.yaml

# 2. 수정
vi aaa.yaml

# 3-1. deployment면
kubectl apply -f aaa.yaml

# 3-2. pod이면
kubectl replace --force -f aaa.yaml
# delete + apply 한번에 해줌
```

---

## 시험 꿀팁

```bash
# replace --force 손에 익혀두세요
kubectl replace --force -f aaa.yaml
# 기존 pod 지우고 새로 만드는거 한방에
```

`kubectl edit`랑 `replace --force` 두 개가 시험에서 제일 많이 쓰는 패턴이에요.