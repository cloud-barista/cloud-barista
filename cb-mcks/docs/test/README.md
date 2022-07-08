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
$ docker run -d -p 1024:1024 --name cb-spider cloudbaristaorg/cb-spider:0.x.y
$ docker run -d -p 1323:1323 --name cb-tumblebug --link cb-spider:cb-spider cloudbaristaorg/cb-tumblebug:0.x.y
```
* 각 컨테이너 이미지의 최신 tag는 다음을 참조
  * https://hub.docker.com/r/cloudbaristaorg/cb-spider/tags
  * https://hub.docker.com/r/cloudbaristaorg/cb-tumblebug/tags

* 예
```
$ docker run -d -p 1024:1024 --name cb-spider cloudbaristaorg/cb-spider:0.5.0
$ docker run -d -p 1323:1323 --name cb-tumblebug --link cb-spider:cb-spider cloudbaristaorg/cb-tumblebug:0.5.0
```

### CB-Dragonfly 실행 (모니터링 에이전트 설치를 원할 경우)

- [CB-Dragonfly 실행 방법](https://github.com/cloud-barista/cb-dragonfly#2-실행-방법) 참조


### Cloud Connection Info. 등록

#### `batch-register-cloud-info.sh` 활용
다음의 단계를 수행합니다.
- `cb-mcks/docs/test/` 에 있는 `batch-register-cloud-info.sh.example` 파일을 `batch-register-cloud-info.sh` 로 복사
- `batch-register-cloud-info.sh` 파일을 텍스트 에디터로 오픈
- 자신이 발급받은 클라우드 별 연결정보를 각 환경변수에 입력
- `batch-register-cloud-info.sh` 파일을 실행

아래에 소개되는 각 CSP별 가이드는
- `batch-register-cloud-info.sh` 의 각 환경변수에 대한 설명이며, 또한
- `batch-register-cloud-info.sh` 를 사용하지 않고 수동으로 등록하는 경우를 위한 매뉴얼로도 활용될 수 있습니다.

####  GCP

* 환경변수 : 클라우드별 연결정보

```
$ export GCP_PROJECT="<project ID>"
$ export GCP_PKEY="<private key>"
$ export GCP_SA="<service account email>"
```

* 환경변수 : GCP_REGION, GCP_ZONE

```
$ export GCP_REGION="<region name>" 
$ export GCP_ZONE="<zone name>"

# 예 : asia-northeast3 (서울리전)
$ export GCP_REGION="asia-northeast3" 
$ export GCP_ZONE="asia-northeast3-a"
```

* Cloud Connection Info. 등록

```
$ ./connectioninfo-create.sh GCP
```

#### AWS

* 환경변수 : 클라우드별 연결정보

```
$ export AWS_KEY="<aws_access_key_id>"
$ export AWS_SECRET="<aws_secret_access_key>"
```

* 환경변수 : AWS_REGION, AWS_ZONE

```
$ export AWS_REGION="<region name>" 
$ export AWS_ZONE="<zone name>"

# 예: ap-northeast-2 (서울리전)
$ export AWS_REGION="ap-northeast-2"
$ export AWS_ZONE="ap-northeast-2a"

# 예: ap-northeast-1 (일본리전)
$ export AWS_REGION="ap-northeast-1"
$ export AWS_ZONE="ap-northeast-1a"
```

* Cloud Connection Info. 등록

```
$ ./connectioninfo-create.sh AWS
```

#### AZURE

* 환경변수 : 클라우드별 연결정보

```
$ export AZURE_CLIENT_ID="<azure_client_id>"
$ export AZURE_CLIENT_SECRET="<azure_client_secret>"
$ export AZURE_TENANT_ID="<azure_tenant_id>"
$ export AZURE_SUBSCRIPTION_ID="<azure_subscription_id>"
```

* 환경변수 : AZURE_REGION, AZURE_RESOURCE_GROUP

```
$ export AZURE_REGION="<region name>" 
$ export AZURE_RESOURCE_GROUP="<resource group>"

# 예: koreacentral (한국 중부)
$ export AZURE_REGION="koreacentral"
$ export AZURE_RESOURCE_GROUP="cb-mcksRG"
```

* Cloud Connection Info. 등록

```
$ ./connectioninfo-create.sh AZURE
```

#### ALIBABA

* 환경변수 : 클라우드별 연결정보

```
$ export ALIBABA_KEY="<alibaba_access_key_id>"
$ export ALIBABA_SECRET="<alibaba_access_key_secret>"
```

* 환경변수 : ALIBABA_REGION, ALIBABA_ZONE

```
$ export ALIBABA_REGION="<region name>" 
$ export ALIBABA_ZONE="<zone name>"

# 예: ap-northeast-1 (도쿄리전)
$ export ALIBABA_REGION="ap-northeast-1"
$ export ALIBABA_ZONE="ap-northeast-1a"
```

* Cloud Connection Info. 등록

```
$ ./connectioninfo-create.sh ALIBABA
```
#### TENCENT

* 환경변수 : 클라우드별 연결정보

```
$ export TENCENT_KEY="<tencent_access_key_id>"
$ export TENCENT_SECRET="<tencent_access_key_secret>"
```

* 환경변수 : TENCENT_REGION, TENCENT_ZONE

```
$ export TENCENT_REGION="<region name>" 
$ export TENCENT_ZONE="<zone name>"

# 예: ap-seoul (서울리전)
$ export TENCENT_REGION="ap-seoul"
$ export TENCENT_ZONE="ap-seoul-1"
```

* Cloud Connection Info. 등록

```
$ ./connectioninfo-create.sh TENCENT
```

#### OPENSTACK

* 환경변수 : 클라우드별 연결정보

```
$ export OS_AUTH_URL="<openstack_auth_url>"
$ export OS_USERNAME="<openstack_username>"
$ export OS_PASSWORD="<openstack_password>"
$ export OS_USER_DOMAIN_NAME="<openstack_domainname>"
$ export OS_PROJECT_ID="<openstack_project_id>"
```

* 환경변수 : OS_REGION, OS_ZONE

```
$ export OS_REGION="RegionOne"
$ export OS_ZONE="RegionOne"
```

* Cloud Connection Info. 등록

```
$ ./connectioninfo-create.sh OPENSTACK
```

#### IBM

* 환경변수 : 클라우드별 연결정보

```
$ export IBM_API_KEY="<ibm_api_key>"
```

* 환경변수 : IBM_REGION, IBM_ZONE

```
$ export IBM_REGION="<region name>" 
$ export IBM_ZONE="<zone name>"

# 예: jp-tok (도쿄리전)
$ export IBM_REGION="jp-tok"
$ export IBM_ZONE="jp-tok-1"
```

* Cloud Connection Info. 등록

```
$ ./connectioninfo-create.sh IBM
```

#### CLOUDIT

* 환경변수 : 클라우드별 연결정보

```
$ export CI_IDENTITY_ENDPOINT="<cloudit_identity_endpoint>"
$ export CI_USERNAME="<cloudit_username>"
$ export CI_PASSWORD="<cloudit_password>"
$ export CI_AUTH_TOKEN="<cloudit_auth_token>"
$ export CI_TENANT_ID="<cloudit_tenant_id>"
```

* 환경변수 : CI_REGION, CI_ZONE

```
$ export OS_REGION="Region"
$ export OS_ZONE="default"
```

* Cloud Connection Info. 등록

```
$ ./connectioninfo-create.sh CLOUDIT
```

#### Cloud Connection Info 추가

```
# AWS/GCP/ALIBABA/TENCENT/OPENSTACK/IBM/CLOUDIT
$ export [AWS/GCP/ALIBABA/TENCENT/OS/IBM/CLOUDIT]_REGION="<region name>"
$ export [AWS/GCP/ALIBABA/TENCENT/OS/IBM/CLOUDIT]_ZONE="<zone name>"

# AZURE
$ export AZURE_REGION="<region name>"
$ export AZURE_RESOURCE_GROUP="<resource group>"

$ ./connectioninfo-create.sh [AWS/GCP/AZURE/ALIBABA/TENCENT/OPENSTACK/IBM/CLOUDIT] add
```

#### 결과 확인

```
$ ./connectioninfo-list.sh all
```

#### namespace 등록

```
$ ./ns-create.sh <namespace>

# 예
$ ./ns-create.sh cb-mcks-ns

# 결과 확인
$ ./ns-get.sh cb-mcks-ns
```


## Test 

### CB-MCKS 실행

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
$ ./cluster-create.sh cb-mcks-ns cluster-01
```

### 클러스터 확인
```
$ ./cluster-get.sh <namespace> <cluster name>
```

* 예
```
$ ./cluster-get.sh cb-mcks-ns cluster-01
```

### 클러스터 삭제
```
$ ./cluster-delete.sh <namespace> <cluster name>
```

* 예
```
$ ./cluster-delete.sh cb-mcks-ns cluster-01
```

### 클러스터 리스트
```
$ ./cluster-list.sh <namespace>
```

* 예
```
$ ./cluster-list.sh cb-mcks-ns
```

### 노드 생성
```
$ ./node-add.sh <namespace> <cluster name>
```

* 예
```
$ ./node-add.sh cb-mcks-ns cluster-01
```

### 노드 확인
```
$ ./node-get.sh <namespace> <cluster name> <node name>
```

* 예
```
$ ./node-get.sh cb-mcks-ns cluster-01 cluster-01-w-1-asdflk
```

### 노드 삭제

```
$ ./node-remove.sh <namespace> <cluster name> <node name>
```

* 예
```
$ ./node-remove.sh cb-mcks-ns cluster-01 cluster-01-w-2-asdflk
```

### 노드 리스트
```
$ ./node-list.sh <namespace> <cluster name>
```

* 예
```
$ ./node-list.sh cb-mcks-ns cluster-01
```

## Kubernetes 클러스터 연결

### kubeconfig 파일 다운로드

* `kubeconfig.yaml` 파일이 생성됩니다.
```
$ ./cluster-get-kubeconfig.sh <namespace> <cluster name>
```

* 예
```
$ ./cluster-get-kubeconfig.sh cb-mcks-ns cluster-01
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
$ ./save-sshkey.sh cb-mcks-ns config-asia-northeast3
$ cat *.pem
```

### 파일에서 클라우드별 연결정보 얻기

> 참고 : CSP 별 credential 생성 가이드 
> 
> https://github.com/cloud-barista/cb-coffeehouse/wiki/A-step-by-step-guide-to-creating-credentials-of-each-cloud-service-provider

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

* AZURE
```
$ source ./env.sh AZURE "<credentials file path>"

# 예
$ source ./env.sh AZURE "${HOME}/.azure/credentials"
```

* ALIBABA
```
$ source ./env.sh ALIBABA "<credentials file path>"

# 예
$ source ./env.sh ALIBABA "${HOME}/.ssh/AccessKey.csv"
```

* TENCENT
```
$ source ./env.sh TENCENT "<credentials file path>"

# 예
$ source ./env.sh TENCENT "${HOME}/.tccli/default.credential"
```

* OPENSTACK
```
$ source ./env.sh OPENSTACK "<credentials file path>"

# 예
$ source ./env.sh OPENSTACK "${HOME}/openrc.sh"
```


### clean up

* MCIR (vpc, securityGroup, sshKey, spec, image) 확인

```
$ ./mcir-list.sh <namespace> [all/image/spec/ssh/sg/vpc]

# 예
$ ./mcir-list.sh cb-mcks-ns all
```

* MCIR (vpc, securityGroup, sshKey, spec, image) 삭제

```
$ ./mcir-delete.sh <namespace> [all/image/spec/ssh/sg/vpc]

# 예
$ ./mcir-delete.sh cb-mcks-ns all
```
