실제환경에선 yaml 많이 쓰나 시험에선 --from-literal 이 더 빠름

yaml의 data 부분에 들어감


```bash
# 이 명령어가
kubectl create configmap myconfig \
  --from-literal=DB_HOST=localhost \
  --from-literal=DB_PORT=5432
```

```yaml
# 이 YAML과 동일해요
apiVersion: v1
kind: ConfigMap
metadata:
  name: myconfig
data:
  DB_HOST: localhost   # --from-literal=DB_HOST=localhost
  DB_PORT: "5432"      # --from-literal=DB_PORT=5432
```

확인하고 싶으면:

```bash
kubectl get configmap myconfig -o yaml
```
