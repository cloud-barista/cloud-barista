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
$ docker run -d -p 1024:1024 --name cb-spider cloudbaristaorg/cb-spider:v0.x.y
$ docker run -d -p 1323:1323 --name cb-tumblebug --link cb-spider:cb-spider cloudbaristaorg/cb-tumblebug:v0.x.y
```
* 각 컨테이너 이미지의 최신 tag는 다음을 참조
  * https://hub.docker.com/r/cloudbaristaorg/cb-spider/tags
  * https://hub.docker.com/r/cloudbaristaorg/cb-tumblebug/tags

* 예
```
$ docker run -d -p 1024:1024 --name cb-spider cloudbaristaorg/cb-spider:v0.3.0-espresso
$ docker run -d -p 1323:1323 --name cb-tumblebug --link cb-spider:cb-spider cloudbaristaorg/cb-tumblebug:v0.3.0-espresso
```

### Cloud Connection Info. 등록

####  GCP

* 환경변수 : 클라우드별 연결정보

```
$ export PROJECT="<project name>"
$ export PKEY="<private key>"
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
$ ./connectioninfo-create.sh GCP
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

# 예: ap-northeast-2 (서울리전)
$ export REGION="ap-northeast-2"
$ export ZONE="ap-northeast-2a"

# 예: ap-northeast-1 (일본리전)
$ export REGION="ap-northeast-1"
$ export ZONE="ap-northeast-1a"
```

* Cloud Connection Info. 등록

```
$ ./connectioninfo-create.sh AWS
```

#### AZURE

* 환경변수 : 클라우드별 연결정보

```
$ export CLIENT_ID="<azure_client_id>"
$ export CLIENT_SECRET="<azure_client_secret>"
$ export TENANT_ID="<azure_tenant_id>"
$ export SUBSCRIPTION_ID="<azure_subscription_id>"
```

* 환경변수 : REGION, RESOURCE_GROUP

```
$ export REGION="<region name>" 
$ export RESOURCE_GROUP="<resource group>"

# 예: koreacentral (한국 중부)
$ export REGION="koreacentral"
$ export RESOURCE_GROUP="cb-ladybugRG"
```

* Cloud Connection Info. 등록

```
$ ./connectioninfo-create.sh AZURE
```

#### Cloud Connection Info 추가

```
# AWS/GCP
$ export REGION="<region name>"
$ export ZONE="<zone name>"

# AZURE
$ export REGION="<region name>"
$ export RESOURCE_GROUP="<resource group>"

$ ./connectioninfo-create.sh [AWS/GCP/AZURE] add
```

#### 결과 확인

```
$ ./connectioninfo-list.sh all
```

#### namespace 등록

```
$ ./ns-create.sh <namespace>

# 예
$ ./ns-create.sh cb-ladybug-ns

# 결과 확인
$ ./ns-get.sh cb-ladybug-ns
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
$ ./cluster-create.sh <namespace> <cluster name>
```

* 예
```
$ ./cluster-create.sh cb-ladybug-ns cluster-01
```

### 클러스터 확인
```
$ ./cluster-get.sh <namespace> <cluster name>
```

* 예
```
$ ./cluster-get.sh cb-ladybug-ns cluster-01
```

### 클러스터 삭제
```
$ ./cluster-delete.sh <namespace> <cluster name>
```

* 예
```
$ ./cluster-delete.sh cb-ladybug-ns cluster-01
```

### 클러스터 리스트
```
$ ./cluster-list.sh <namespace>
```

* 예
```
$ ./cluster-list.sh cb-ladybug-ns
```

### 노드 생성
```
$ ./node-add.sh <namespace> <cluster name>
```

* 예
```
$ ./node-add.sh cb-ladybug-ns cluster-01
```

### 노드 확인
```
$ ./node-get.sh <namespace> <cluster name> <node name>
```

* 예
```
$ ./node-get.sh cb-ladybug-ns cluster-01 cluster-01-w-1-asdflk
```

### 노드 삭제

```
$ ./node-remove.sh <namespace> <cluster name> <node name>
```

* 예
```
$ ./node-remove.sh cb-ladybug-ns cluster-01 cluster-01-w-2-asdflk
```

### 노드 리스트
```
$ ./node-list.sh <namespace> <cluster name>
```

* 예
```
$ ./node-list.sh cb-ladybug-ns cluster-01
```

## Kubernetes 클러스터 연결

### kubeconfig 파일 다운로드

* `kubeconfig.yaml` 파일이 생성됩니다.
```
$ ./cluster-get-kubeconfig.sh <namespace> <cluster name>
```

* 예
```
$ ./cluster-get-kubeconfig.sh cb-ladybug-ns cluster-01
```

### kubectl 사용하기

```
$ export KUBECONFIG=$(pwd)/kubeconfig.yaml
$ kubectl get nodes
```

## 기타

### SSH key 파일 저장

```
$ ./save-sshkey.sh <namespace> <connection info>
```

* 예
```
$ ./save-sshkey.sh cb-ladybug-ns config-asia-northeast3
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

### clean up

* MCIR (vpc, securityGroup, sshKey, spec, image) 확인

```
$ ./mcir-list.sh <namespace> [all/image/spec/ssh/sg/vpc]

# 예
$ ./mcir-list.sh cb-ladybug-ns all
```

* MCIR (vpc, securityGroup, sshKey, spec, image) 삭제

```
$ ./mcir-delete.sh <namespace> [all/image/spec/ssh/sg/vpc]

# 예
$ ./mcir-delete.sh cb-ladybug-ns all
```
