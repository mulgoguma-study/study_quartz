
``` bash
alias k=kubectl
export do="--dry-run=client -o yaml"

# 사용 예시:
k run nginx-pod --image=nginx $do > pod.yaml
```
