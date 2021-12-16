
<!-- https://docs.google.com/presentation/d/13a5tXC66jCtaX5lGn1reGlzZWmURXxYpVQpEhxYCQVk/edit?usp=sharing -->

## `cb-operator`의 `Kubernetes 모드`를 이용한 Cloud-Barista 설치 및 실행 가이드

이 가이드에서는 `cb-operator`의 두 가지 모드 중 하나인 `Kubernetes 모드`를 이용하여 Cloud-Barista를 설치하고 실행하는 방법에 대해 소개합니다.

## 순서
1. [참고] 프레임워크별 컨테이너 구성 및 API Endpoint
1. [참고] 그 외 컨테이너 구성 및 Endpoint
1. 개발환경 준비
1. Kubernetes 버전 정하기
1. 필요사항 설치
   1. Golang
   1. Docker
   1. Helm 3
1. Kubernetes 클러스터 준비
   1. 싱글 노드로 구성되는 K8s 클러스터
   1. 멀티 노드로 구성되는 K8s 클러스터
1. cb-operator 소스코드 다운로드
1. 환경설정 확인 및 변경
1. cb-operator 소스코드 빌드
1. cb-operator 이용하여 Cloud-Barista 실행
1. Cloud-Barista 실행상태 확인

## [참고] 프레임워크별 컨테이너 구성 및 API Endpoint
| Framework별 Container Name | REST-API Endpoint | REST-API via APIGW Endpoint | Go-API Endpoint |
|---|---|---|---|
| cb-spider | http://{{host}}:31024/spider | http://{{host}}:30080/spider | http://{{host}}:32048  |
| --- |   |   |   |
| cb-tumblebug | http://{{host}}:31323/tumblebug | http://{{host}}:30080/tumblebug | http://{{host}}:30252  |
| --- |   |   |   |
| cb-mcks | http://{{host}}:31470/mcks | http://{{host}}:30080/mcks |   |
| --- |   |   |   |
| cb-ladybug | http://{{host}}:31592/ladybug | http://{{host}}:30080/ladybug |   |
| --- |   |   |   |
| cb-dragonfly | http://{{host}}:30090/dragonfly | http://{{host}}:30080/dragonfly | <!--30094/udp--> http://{{host}}:30254  |

## [참고] 그 외 컨테이너 구성 및 Endpoint
| Container Name | Endpoint | Misc. |
|---|---|---|
| cb-dragonfly-influxdb | - |   |
| cb-dragonfly-kafka | http://{{host}}:32000 |   |
| cb-dragonfly-kapacitor | -  |   |
| cb-dragonfly-zookeeper | -  |   |
| --- |   |   |
| cb-restapigw | GW: http://{{host}}:30080 <br> Admin: http://{{host}}:30081 | ID: admin <br> PW: test@admin00  |
| cb-restapigw-influxdb | - |   |
| cb-restapigw-grafana | - | ID: admin <br> PW: admin |
| cb-restapigw-jaeger | - |   |
| --- |   |   |
| cb-webtool | http://{{host}}:31234 |   |
| --- |   |   |
| prometheus | -  |   |
| grafana | http://{{host}}:30300  | ID: admin <br> PW: admin  |

## 개발환경 준비

[권장사항]
- Ubuntu 18.04
- Golang 1.15 또는 그 이상

## Kubernetes 버전 정하기
- 참고: [쿠버네티스 버전 및 버전 차이(skew) 지원 정책](https://kubernetes.io/ko/docs/setup/release/version-skew-policy/)
- Kubernetes 버전은 `x.y.z` 로 표현되며, `x`는 메이저 버전, `y`는 마이너 버전, `z`는 패치 버전을 의미합니다.
- Kubernetes 프로젝트는 최근 3개의 마이너 릴리즈 (1.20, 1.19, 1.18) 에 대해 업데이트를 지원합니다.
- 최근 3개의 마이너 릴리즈 (1.20, 1.19, 1.18) 에 대한 최신 패치 버전은 다음과 같습니다. (as of 2021-03-19)
  - 1.20.5
  - 1.19.9
  - 1.18.17
- 먼저, 어떤 마이너 릴리즈를 사용할 지 결정합니다. (예: `1.20`)
- 아래에 등장하는 `kubectl` 은 `kube-apiserver`의 한 단계 마이너 버전(이전 또는 최신) 내에서 지원하므로, `kubectl` 은 `1.21`, `1.20`, `1.19` 중에서 선택하여 설치하면 됩니다.
- 아래에 등장하는 `minikube` 를 이용하여 Kubernetes 클러스터를 만들 때, 위에서 정한 마이너 릴리즈 (예: `1.20`) 에 맞는 패치 릴리즈  (예: `1.20.5`) 를 선택하면 됩니다.
<!-- - 아래에서 `kubectl` 은 APT 패키지 매니저로 설치할 예정인데, 위에서 소개한 최신 패치 버전의 APT 패키지가 아직 등록되지 않았을 수도 있습니다.
- 아래에서 소개되는 `minikube` 도, 위에서 소개한 최신 패치 버전의 Kubernetes 설치를 아직 지원하지 않을 수도 있습니다.
- 아래에 등장하는 `kubectl`, `minikube` 은 최신 패치 버전까 -->

## 필요사항 설치

### Golang 설치
- https://golang.org/doc/install 에서 설명하는 방법대로 설치합니다.

<details>
  <summary>[클릭하여 예시 보기]</summary>
  
```bash
# Golang 다운로드
wget https://golang.org/dl/go1.16.2.linux-amd64.tar.gz

# 기존 Golang 삭제 및 압축파일 해제
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.2.linux-amd64.tar.gz

# ~/.bashrc 또는 ~/.zshrc 등에 다음 라인을 추가
export PATH=$PATH:/usr/local/go/bin

# 셸을 재시작하고 다음을 실행하여 Go 버전 확인
go version
```
</details>

### Docker 설치
- https://docs.docker.com/engine/install/ubuntu/ 에서 설명하는 방법대로 설치합니다.

<details>
  <summary>[클릭하여 예시 보기]</summary>
  
```bash
# 기존에 Docker 가 설치되어 있었다면 삭제
sudo apt remove docker docker-engine docker.io containerd runc

# Docker 설치를 위한 APT repo 추가
sudo apt update

sudo apt install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# x86_64 / amd64
echo \
  "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update

sudo apt install docker-ce docker-ce-cli containerd.io
```
</details>

### Helm 3 설치
- https://helm.sh/docs/intro/install/#from-apt-debianubuntu 에서 설명하는 방법대로 설치합니다.

<details>
  <summary>[클릭하여 예시 보기]</summary>
  
```bash
curl https://baltocdn.com/helm/signing.asc | sudo apt-key add -
sudo apt-get install apt-transport-https --yes
echo "deb https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
sudo apt-get update
sudo apt-get install helm
```
</details>

## Kubernetes 클러스터 준비
1. 싱글 노드로 구성되는 K8s 클러스터
1. 멀티 노드로 구성되는 K8s 클러스터

### 싱글 노드로 구성되는 K8s 클러스터
- 예시로 `minikube` 를 사용합니다.
- https://minikube.sigs.k8s.io/docs/start/#debian-package 에서 설명하는 방법대로 설치합니다.

<details>
  <summary>[클릭하여 예시 보기]</summary>
  
```bash
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube_latest_amd64.deb
sudo dpkg -i minikube_latest_amd64.deb
```
</details>

- `minikube` 를 사용하여 Kubernetes 클러스터를 생성합니다.

<details>
  <summary>[클릭하여 예시 보기]</summary>
  
```bash
minikube start [--driver=none] [--kubernetes-version=v1.20.5]
```
</details>


#### kubectl 설치
- `kubectl` 이 없는 경우, 이 가이드의 내용대로 `kubectl` 을 설치하면 됩니다.
- `kubectl` 은 `kube-apiserver`의 한 단계 마이너 버전(이전 또는 최신) 내에서 지원하므로, `kubectl` 은 `1.21`, `1.20`, `1.19` 중에서 선택하여 설치하면 됩니다.
- https://kubernetes.io/ko/docs/tasks/tools/install-kubectl/#%EA%B8%B0%EB%B3%B8-%ED%8C%A8%ED%82%A4%EC%A7%80-%EA%B4%80%EB%A6%AC-%EB%8F%84%EA%B5%AC%EB%A5%BC-%EC%82%AC%EC%9A%A9%ED%95%98%EC%97%AC-%EC%84%A4%EC%B9%98  에서 설명하는 방법대로 설치합니다.


<details>
  <summary>[클릭하여 예시 보기]</summary>
  
```bash
sudo apt-get update && sudo apt-get install -y apt-transport-https gnupg2 curl
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
sudo apt-cache madison kubectl # 원하는 kubectl 버전을 확인
sudo apt-get install -y kubectl=1.20.4-00
```
</details>

### 멀티 노드로 구성되는 K8s 클러스터
- 예시로 `kubeadm` 을 사용합니다.
- https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/ 에서 설명하는 방법대로 진행합니다.

#### kubeadm, kubectl, kubelet 설치
- https://kubernetes.io/ko/docs/setup/production-environment/tools/kubeadm/install-kubeadm/ 에서 설명하는 방법대로 설치합니다.

<details>
  <summary>[클릭하여 예시 보기]</summary>
  
```bash
sudo apt-get update && sudo apt-get install -y apt-transport-https curl
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF
sudo apt-get update
sudo apt-cache madison kubeadm # 원하는 kubeadm 버전을 확인
sudo apt-get install -y kubelet=1.20.4-00 kubeadm=1.20.4-00 kubectl=1.20.4-00
sudo apt-mark hold kubelet kubeadm kubectl
```
</details>

#### 컨트롤-플레인 노드 (구. 마스터 노드) 설정하기
- CNI 로는 Calico 를 사용하는 것을 가정합니다. (`--pod-network-cidr=192.168.0.0/16`)
- https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/ 에서 설명하는 방법대로 진행합니다.

<details>
  <summary>[클릭하여 예시 보기]</summary>
  
```bash
sudo kubeadm init --pod-network-cidr=192.168.0.0/16
# 여기서 출력되는 kubeadm join ... 명령어를 메모해 둡니다.

# kubectl이 kubeconfig를 인식할 수 있도록 설정
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# Calico 설치 (https://docs.projectcalico.org/getting-started/kubernetes/quickstart)
kubectl create -f https://docs.projectcalico.org/manifests/tigera-operator.yaml
kubectl create -f https://docs.projectcalico.org/manifests/custom-resources.yaml
watch kubectl get pods -n calico-system

# 선택사항: 만약 컨트롤-플레인 노드에도 파드가 할당되도록 하려면
kubectl taint nodes --all node-role.kubernetes.io/master-
```
</details>

#### (워커) 노드 추가하기
- https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/#join-nodes 에서 설명하는 방법대로 진행합니다.

<details>
  <summary>[클릭하여 예시 보기]</summary>
  
```bash
# (워커) 노드로 추가하고자 하는 컴퓨터에서
# 위에서 메모한 kubeadm join ... 명령어를 실행합니다.
kubeadm join --token <token> <control-plane-host>:<control-plane-port> --discovery-token-ca-cert-hash sha256:<hash>
```
</details>

## cb-operator 소스코드 다운로드
```bash
git clone https://github.com/cloud-barista/cb-operator.git
```

## 환경설정 확인 및 변경
- Cloud-Barista를 설치 및 실행하는 VM/물리머신의 Public IP 주소를 알아냅니다.
  - 예: `curl ifconfig.so`
- `cb-operator/helm-chart/charts/cb-dragonfly/files/conf/config.yaml` 파일에 Public IP 주소를 기재합니다.
```YAML
# kafka connection info
kafka:
  endpoint_url: cb-dragonfly-kafka
  external_ip: 127.0.0.1 # Cloud-Barista를 설치 및 실행하는 VM/물리머신의 Public IP 주소를 기재
  deploy_type: "helm"    # deploy environment "compose" => docker-compose or others , "helm" => helm chart on k8s
  compose_external_port: 9092
  helm_external_port: 32000
  internal_port: 9092

# collect manager configuration info
collectManager:
  collector_ip: 127.0.0.1  # Cloud-Barista를 설치 및 실행하는 VM/물리머신의 Public IP 주소를 기재
  collector_port: 8094    # udp port
  collector_group_count: 1      # default collector group count
```

## cb-operator 소스코드 빌드
```bash
cd cb-operator/src
go build -o operator main.go
```

## cb-operator 이용하여 Cloud-Barista 실행
```bash
./operator

# 모드를 고르는 단계가 나오면, 2: Kubernetes 모드 선택

./operator run
```

## Cloud-Barista 실행상태 확인
```bash
./operator info
```

<details>
  <summary>[클릭하여 예시 보기]</summary>
  
```
CB_OPERATOR_MODE: Kubernetes

[Get info for Cloud-Barista runtimes]

[Config path] ../helm-chart/values.yaml

[v]Status of Cloud-Barista Helm release
NAME: cloud-barista
LAST DEPLOYED: Fri Mar 19 17:56:35 2021
NAMESPACE: cloud-barista
STATUS: deployed
REVISION: 1
TEST SUITE: None

[v]Status of Cloud-Barista pods
NAME                                                     READY   STATUS      RESTARTS   AGE
cb-dragonfly-58bf9f6f74-sjnst                            1/1     Running     3          30m
cb-dragonfly-influxdb-0                                  1/1     Running     0          30m
cb-dragonfly-kafka-0                                     1/1     Running     1          30m
cb-dragonfly-zookeeper-0                                 1/1     Running     0          30m
cb-ladybug-667468cf88-bpqtg                              1/1     Running     0          30m
cb-restapigw-6f797fc998-p2shj                            1/1     Running     0          30m
cb-restapigw-influxdb-6cdb79759c-qsssf                   1/1     Running     0          30m
cb-restapigw-jaeger-agent-ffjzr                          1/1     Running     0          30m
cb-restapigw-jaeger-cassandra-schema-cv9mx               0/1     Completed   1          30m
cb-restapigw-jaeger-collector-84f7978ff9-7f865           1/1     Running     5          30m
cb-restapigw-jaeger-query-68655bb76d-xrltf               2/2     Running     5          30m
cb-spider-57db74c88f-m6klq                               1/1     Running     0          30m
cb-tumblebug-7d8ff8ff8d-9pv52                            1/1     Running     0          30m
cb-webtool-698d445765-k7p4c                              1/1     Running     0          30m
cloud-barista-cassandra-0                                1/1     Running     0          30m
cloud-barista-cassandra-1                                1/1     Running     0          28m
cloud-barista-cassandra-2                                1/1     Running     0          27m
cloud-barista-cb-dragonfly-kapacitor-5589f8f8db-xfb5f    1/1     Running     0          30m
cloud-barista-grafana-df64c9986-dm7ws                    2/2     Running     0          30m
cloud-barista-kube-state-metrics-d845db79c-gtdlg         1/1     Running     0          30m
cloud-barista-prometheus-alertmanager-687cc55586-s5b8m   2/2     Running     0          30m
cloud-barista-prometheus-node-exporter-z974q             1/1     Running     0          30m
cloud-barista-prometheus-pushgateway-854b87889b-thc6j    1/1     Running     0          30m
cloud-barista-prometheus-server-78654ccd49-zd6fl         2/2     Running     0          30m

[v]Status of Cloud-Barista container images
bitnami/influxdb:1.8.0-debian-10-r37
cassandra:3.11.6
cloudbaristaorg/cb-dragonfly:v0.3.0-espresso
...
```
</details>

## Cloud-Barista 중지
```bash
./operator stop
```
