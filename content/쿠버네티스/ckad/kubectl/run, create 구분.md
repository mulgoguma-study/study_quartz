
```bash
# run - Pod 전용
kubectl run mypod --image=nginx

# create - 나머지 거의 다
kubectl create deployment mydeploy --image=nginx
kubectl create service clusterip mysvc --tcp=80:80
kubectl create namespace mynamespace
kubectl create configmap myconfig --from-literal=KEY=VALUE
kubectl create secret generic mysecret --from-literal=PW=1234
```

**암기 팁:**

- `run` = **Pod 딱 하나** 띄울 때
- `create` = **그 외 전부**

---

근데 사실 시험에서 더 많이 쓰는 패턴이 있어요:

```bash
# dry-run으로 yaml 뽑고 수정해서 apply
kubectl run mypod --image=nginx --dry-run=client -o yaml > pod.yaml
kubectl create deployment mydeploy --image=nginx --dry-run=client -o yaml > deploy.yaml

vim pod.yaml  # 수정
kubectl apply -f pod.yaml
```

외워야 할 건 결국:

- `run` → Pod
- `create` → 나머지
- `--dry-run=client -o yaml` → yaml 뽑을 때

이 세 가지가 핵심이에요 👍