entimport는 너무 이슈가 많아서 포기!
아틀라스가 공식이니 이거 쓰자


### 설치법
``` bash
curl -sSf https://atlasgo.sh | sh
```



### 아틀라스 CLI 활용
`atlas`는 ent의 공식 파트너 도구고, DB 스키마를 읽어서 ent 마이그레이션 파일로 변환하는 기능이 있어. 특히 `atlas schema inspect` 명령으로 HCL 또는 SQL 형식으로 추출하고, 이를 기반으로 ent 스키마 코드를 생성할 수 있어.
``` bash
#### DB 스키마 추출
atlas schema inspect \
  -u "postgres://user:pass@localhost/dbname?search_path=public" \
  --format '{{ json . }}' > schema.json

#### 이후 schema.json을 파싱해서 ent 파일 생성
```


### 순서
```
go install ariga.io/atlas/cmd/atlas@latest
```

#### 실제 추천 플로우
```
DB
 │
 ▼
atlas schema inspect  →  schema.json (또는 HCL)
 │
 ▼
커스텀 Go 스크립트   →  ent/schema/*.go (Fields만, Edges는 비워둠)
 │
 ▼
go generate ./ent
 │
 ▼
Edges/Relations 수동 추가
```