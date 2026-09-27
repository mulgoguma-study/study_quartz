네 맞아요! **vim으로 직접 타이핑(또는 복붙)해서 넣어야** 합니다. `kubectl run`으로는 volume 설정을 만들어주지 못하거든요.

## 왜 직접 쳐야 하나

`kubectl run`의 플래그를 봐도 volume 관련은 없어요:

```bash
kubectl run --help | grep -i volume
# → 아무것도 안 나옴
```

그래서 순서가 이렇게 돼요:

```bash
# 1단계: 뼈대만 자동 생성
kubectl run app --image=nginx --dry-run=client -o yaml > app.yaml

# 2단계: vim으로 열어서 volume 부분 손으로 추가
vim app.yaml
```

## vim 사용법 간단 정리 (CKAD 필수)

Go 개발자시니까 이미 vim 쓰실 것 같긴 한데, 핵심만 짚어드리면:

```
i           → 입력(insert) 모드 진입
Esc         → 입력 모드 종료
dd          → 현재 줄 삭제
yy          → 현재 줄 복사
p           → 붙여넣기
:w          → 저장
:wq         → 저장 후 종료
:q!         → 저장 없이 종료
/keyword    → 검색
```

## CKAD 시험 꿀팁: 공식 docs에서 복붙

시험에서는 **kubernetes.io/docs 접속이 허용**돼요. volume 섹션을 외우는 것보다 **docs에서 예시 복사 → 수정**이 훨씬 빠릅니다.

### 시험 중 실전 흐름

```bash
# 1. 뼈대 만들기
k run app --image=nginx --dry-run=client -o yaml > app.yaml

# 2. 브라우저에서 docs 열기
#    검색: "kubernetes secret volume"
#    페이지: Secrets → "Using Secrets as files from a Pod"
#    거기에 YAML 예시 있음

# 3. vim으로 app.yaml 열기
vim app.yaml
```

vim 안에서:

```
:set paste        ← 들여쓰기 자동 변환 방지 (중요!)
o                 ← 다음 줄에 입력 모드 진입
(브라우저에서 복사한 내용 붙여넣기)
Esc
:wq
```

`:set paste`를 **꼭 먼저** 쳐야 해요. 안 하면 vim이 자동 들여쓰기를 추가해서 YAML이 완전히 깨집니다.

### 또 다른 방법: heredoc으로 통째로 작성

vim 편집이 어렵다면 처음부터 파일 전체를 만들어도 돼요:

```bash
cat > app.yaml <<'EOF'
apiVersion: v1
kind: Pod
metadata:
  name: app
spec:
  containers:
  - name: app
    image: nginx
    volumeMounts:
    - name: db
      mountPath: /etc/db
      readOnly: true
  volumes:
  - name: db
    secret:
      secretName: db-secret
EOF
```

`<<'EOF'`의 작은따옴표가 **중요해요**. `$`, `` ` `` 같은 특수문자가 있을 때 bash 확장되는 걸 막아줍니다.

## 현실적인 시험 전략

실제로 CKAD 합격자들이 쓰는 방식:

1. **반복되는 볼륨 패턴은 외워둔다** (emptyDir, secret volume, configMap volume 등 5~6개 정도)
2. **docs를 북마크해둔다** — 필요하면 바로 검색해서 복사
3. **YAML 들여쓰기는 vim `:set paste` 모드**에서 수정

CKAD는 속도 싸움이라 **"머리에서 YAML 나오도록 외우는 것"**과 **"docs를 빠르게 찾는 것"** 둘 중에 선택이 아니라 둘 다 해야 합니다. Secret/ConfigMap volume은 시험 단골이라 외워두는 쪽이 좋아요.

## 자주 쓰는 volume 패턴 5개 (암기 추천)

**1. Secret을 파일로 마운트**

```yaml
    volumeMounts:
    - name: vol
      mountPath: /etc/secret
  volumes:
  - name: vol
    secret:
      secretName: my-secret
```

**2. ConfigMap을 파일로 마운트**

```yaml
    volumeMounts:
    - name: vol
      mountPath: /etc/config
  volumes:
  - name: vol
    configMap:
      name: my-configmap
```

**3. emptyDir (Pod 내 컨테이너 간 공유용)**

```yaml
    volumeMounts:
    - name: shared
      mountPath: /shared
  volumes:
  - name: shared
    emptyDir: {}
```

**4. hostPath (노드의 파일시스템)**

```yaml
  volumes:
  - name: vol
    hostPath:
      path: /data
      type: Directory
```

**5. PVC (영구 스토리지)**

```yaml
  volumes:
  - name: vol
    persistentVolumeClaim:
      claimName: my-pvc
```

이 5개 구조만 머릿속에 있으면 CKAD volume 문제 90%는 커버됩니다.

한 번 직접 vim에서 편집해보시고, 안 되는 부분 있으면 얘기해주세요!