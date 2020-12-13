# Test 
Test shell scripts 사용법

## Prerequisites 
> 클러스터 생성/삭제 기능을 이용하기 위해서는 Cloud Connection 정보를 등록해야 합니다.

### jq 설치
* shell 에서 json parsing 시 `jq` 유틸리티를 활용합니다.
* https://stedolan.github.io/jq/

```
$ brew install jq           # mac os
$ sudo apt-get install jq   # linux
```

### CB-Spider, CB-Tumblebug 실행

```
$  docker run -d -p 1024:1024 --name cb-spider cloudbaristaorg/cb-spider:v0.x.0-yyyymmdd
$  docker run -d -p 1323:1323 --name cb-tumblebug --link cb-spider:cb-spider cloudbaristaorg/cb-tumblebug:v0.x.0-yyyymmdd
```
* 각 컨테이너 이미지의 최신 tag는 다음을 참조
  * https://hub.docker.com/r/cloudbaristaorg/cb-spider/tags
  * https://hub.docker.com/r/cloudbaristaorg/cb-tumblebug/tags

* 예
```
$  docker run -d -p 1024:1024 --name cb-spider cloudbaristaorg/cb-spider:v0.2.0-20200715
$  docker run -d -p 1323:1323 --name cb-tumblebug --link cb-spider:cb-spider cloudbaristaorg/cb-tumblebug:v0.2.5
```

### Cloud Connection Info. 등록

####  GCP

* 환경변수 : 클라우드별 연결정보

```
$ export PROJECT="<project name>"
$ export PKEY="private key>"
$ export SA="<service account email>"
```

* 환경변수 : REGION, ZONE

```
$ export REGION="<region name>" 
$ export ZONE="<zone name>"

# 예 : asia-northeast3 (서울리전)
$ export REGION="asia-northeast3" 
$ export ZONE="asia-northeast3-a"
```

* Cloud Connection Info. 등록

```
$ ./init.sh GCP
```

* 결과 확인

```
$ ./get.sh GCP ns,config
```

#### AWS

* 환경변수 : 클라우드별 연결정보

```
$ export KEY="<aws_access_key_id>"
$ export SECRET="<aws_secret_access_key>"
```

* 환경변수 : REGION, ZONE

```
$ export REGION="<region name>" 
$ export ZONE="<zone name>"

# 예: ap-northeast-1 (일본리전)
$ export REGION="ap-northeast-1"
$ export ZONE="ap-northeast-1a"
```

* Cloud Connection Info. 등록

```
$ ./init.sh AWS
```

* 결과 확인

```
$ ./get.sh AWS ns,config
```

## Test 

### cb-ladybug 실행

```
$ export CBLOG_ROOT="$(pwd)"
$ export CBSTORE_ROOT="$(pwd)"
$ go run src/main.go
```

### 클러스터 생성
```
$ /cluster-create.sh [GCP/AWS] <cluster name> <spec:machine-type> <worker-node-count>
```

* 예
```
$ ./cluster-create.sh GCP cb-cluster n1-standard-2 1   # GCP
$ ./cluster-create.sh AWS cb-cluster t2.medium 1       # AWS
```

### 클러스터 삭제
```
$ /cluster-delete.sh [GCP/AWS] <cluster name>
```

* 예
```
$ ./cluster-delete.sh GCP cb-cluster   # GCP
$ ./cluster-delete.sh AWS cb-cluster   # AWS
```

### 노드 생성
```
$ /node-add.sh [GCP/AWS] <cluster name> <spec:machine-type> <worker-node-count>
```

* 예
```
$ ./node-add.sh GCP cb-cluster n1-standard-2 1   # GCP
$ ./node-add.sh AWS cb-cluster t2.medium 1       # AWS
```

### 노드 삭제

```
$ /node-remove.sh [GCP/AWS] <cluster name> <node name>
```

* 예
```
$ ./node-remove.sh GCP cb-cluster cb-gcp-cluster-test-1-w-q3ui2  # GCP
$ ./node-remove.sh AWS cb-cluster cb-aws-cluster-test-1-w-iqp7n  # AWS
```

## Kubernetes 클러스터 연결

### kubeconfig 파일 다운로드

* `kubeconfig.yaml` 파일이 생성됩니다.
```
$ ./cluster-kubeconfig.sh [GCP/AWS] <cluster name>
```

* 예
```
$ ./cluster-kubeconfig.sh AWS cb-cluster
```

### kubectl 사용하기

```
$ export KUBECONFIG=$(pwd)/kubeconfig.yaml
$ kubectl config set-cluster kubernetes --insecure-skip-tls-verify=true
$ kubectl get nodes
```

## 기타

### SSH key 파일 저장

```
$ ./savekey.sh [AWS/GCP] <cluster name>
```

* 예
```
$ ./savekey.sh AWS cb-cluster
$ cat *.pem
```

### 파일에서 클라우드별 연결정보 얻기

* GCP ( [jq](https://stedolan.github.io/jq/) 설치 필요)

```
$ source ./env.sh GCP "<json file path>"

# 예
$ source ./env.sh GCP "${HOME}/.ssh/google-credential-cloudbarista.json"
```

* AWS
```
$ source ./env.sh AWS "<credentials file path>"

# 예
$ source ./env.sh AWS "${HOME}/.aws/credentials"

# '${HOME}/.aws/credentials' file format which is created by awscli
$ cat ${HOME}/.aws/credentials
[default]
aws_secret_access_key = y7Ganz6A.................................
aws_access_key_id = AKIA2Z........................
```
