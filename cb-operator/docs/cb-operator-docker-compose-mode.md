
<!-- https://docs.google.com/presentation/d/13a5tXC66jCtaX5lGn1reGlzZWmURXxYpVQpEhxYCQVk/edit?usp=sharing -->

## `cb-operator`의 `Docker Compose 모드`를 이용한 Cloud-Barista 설치 및 실행 가이드

이 가이드에서는 `cb-operator`의 두 가지 모드 중 하나인 `Docker Compose 모드`를 이용하여 Cloud-Barista를 설치하고 실행하는 방법에 대해 소개합니다.

## 순서
1. [참고] 프레임워크별 컨테이너 구성 및 API Endpoint
1. [참고] 그 외 컨테이너 구성 및 Endpoint
1. 개발환경 준비
1. 필요사항 설치
   1. Golang
   1. Docker
   1. Docker Compose
1. cb-operator 소스코드 다운로드
1. 환경설정 확인 및 변경
1. cb-operator 소스코드 빌드
1. cb-operator 이용하여 Cloud-Barista 실행
1. Cloud-Barista 실행상태 확인

## [참고] 프레임워크별 컨테이너 구성 및 API Endpoint
| Framework별 Container Name | REST-API Endpoint | REST-API via APIGW Endpoint | Go-API Endpoint |
|---|---|---|---|
| cb-spider | http://{{host}}:1024/spider | http://{{host}}:8000/spider | http://{{host}}:2048  |
| --- |   |   |   |
| cb-tumblebug | http://{{host}}:1323/tumblebug | http://{{host}}:8000/tumblebug | http://{{host}}:50252  |
| --- |   |   |   |
| cb-ladybug | http://{{host}}:8080/ladybug | http://{{host}}:8000/ladybug  |   |
| --- |   |   |   |
| cb-dragonfly | http://{{host}}:9090/dragonfly | http://{{host}}:8000/dragonfly | http://{{host}}:9999<!--8094/udp-->  |

## [참고] 그 외 컨테이너 구성 및 Endpoint
| Container Name | Endpoint | Misc. |
|---|---|---|
| cb-dragonfly-influxdb | http://{{host}}:28086 |   |
| cb-dragonfly-kafka | http://{{host}}:9092 |   |
| cb-dragonfly-kapacitor | http://{{host}}:29092  |   |
| cb-dragonfly-zookeeper | http://{{host}}:2181  |   |
| --- |   |   |
| cb-restapigw | GW: http://{{host}}:8000 <br> Admin: http://{{host}}:8001 | ID: admin <br> PW: test@admin00  | 
| cb-restapigw-influxdb | http://{{host}}:8086 |   |
| cb-restapigw-grafana | http://{{host}}:3100 | ID: admin <br> PW: admin  |
| cb-restapigw-jaeger | http://{{host}}:16686 |   |
| --- |   |   |
| cb-webtool | http://{{host}}:1234 |   |
| --- |   |   |
| cb-tumblebug-phpliteadmin | http://{{host}}:2015  |   |


## 개발환경 준비

[권장사항]
- Ubuntu 18.04
- Golang 1.15 또는 그 이상

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

### Docker Compose 설치
- APT 패키지 매니저를 이용하여 설치합니다.
```bash
sudo apt install docker-compose
```

## cb-operator 소스코드 다운로드
```bash
git clone https://github.com/cloud-barista/cb-operator.git
```

## 환경설정 확인 및 변경
- Cloud-Barista를 설치 및 실행하는 VM/물리머신의 Public IP 주소를 알아냅니다.
  - 예: `curl ifconfig.so`
- `cb-operator/docker-compose-mode-files/conf/cb-dragonfly/config.yaml` 파일에 Public IP 주소를 기재합니다.
```YAML
# kafka connection info
kafka:
  endpoint_url: cb-dragonfly-kafka
  external_ip: 127.0.0.1 # Cloud-Barista를 설치 및 실행하는 VM/물리머신의 Public IP 주소를 기재
  deploy_type: "compose"    # deploy environment "compose" => docker-compose or others , "helm" => helm chart on k8s
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

# 모드를 고르는 단계가 나오면, 1: Docker Compose 모드 선택

./operator run
```

## Cloud-Barista 실행상태 확인
```bash
./operator info
```

<details>
  <summary>[클릭하여 예시 보기]</summary>
  
```
CB_OPERATOR_MODE: DockerCompose

[Get info for Cloud-Barista runtimes]

[Config path] ../docker-compose-mode-files/docker-compose.yaml


[v]Status of Cloud-Barista runtimes
          Name                         Command               State                                                    Ports
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------
cb-dragonfly                cb-dragonfly                     Exit 2
cb-dragonfly-influxdb       /entrypoint.sh influxd           Up       0.0.0.0:28083->8083/tcp, 0.0.0.0:28086->8086/tcp
cb-dragonfly-kafka          start-kafka.sh                   Up       0.0.0.0:9092->9092/tcp
cb-dragonfly-kapacitor      /entrypoint.sh kapacitord        Up       0.0.0.0:29092->9092/tcp
cb-dragonfly-zookeeper      /bin/sh -c /usr/sbin/sshd  ...   Up       0.0.0.0:2181->2181/tcp, 22/tcp, 2888/tcp, 3888/tcp
cb-ladybug                  /app/cb-ladybug                  Up       0.0.0.0:8080->8080/tcp
cb-restapigw                /app/cb-restapigw -c /app/ ...   Up       0.0.0.0:8000->8000/tcp, 0.0.0.0:8001->8001/tcp
cb-restapigw-grafana        /run.sh                          Up       0.0.0.0:3100->3000/tcp
cb-restapigw-influxdb       /entrypoint.sh influxd           Up       0.0.0.0:8083->8083/tcp, 0.0.0.0:8086->8086/tcp
cb-restapigw-jaeger         /go/bin/all-in-one-linux - ...   Up       14250/tcp, 0.0.0.0:14268->14268/tcp, 0.0.0.0:16686->16686/tcp, 5775/udp, 5778/tcp, 6831/udp, 6832/udp
cb-spider                   /root/go/src/github.com/cl ...   Up       0.0.0.0:1024->1024/tcp, 0.0.0.0:2048->2048/tcp, 4096/tcp
cb-tumblebug                /app/src/cb-tumblebug            Up       0.0.0.0:1323->1323/tcp, 0.0.0.0:50252->50252/tcp
cb-tumblebug-phpliteadmin   /usr/bin/caddy --conf /etc ...   Up       0.0.0.0:2015->2015/tcp, 443/tcp, 80/tcp

[v]Status of Cloud-Barista runtime images
        Container                    Repository                  Tag           Image Id      Size
---------------------------------------------------------------------------------------------------
cb-dragonfly                cloudbaristaorg/cb-dragonfly   v0.3.0-espresso   00badc2e5613   125 MB
cb-dragonfly-influxdb       influxdb                       1.8-alpine        97eae8355b82   175 MB
cb-dragonfly-kafka          wurstmeister/kafka             2.12-2.4.1        2dd8b556702e   412 MB
cb-dragonfly-kapacitor      kapacitor                      1.5               95490156d6f2   232 MB
cb-dragonfly-zookeeper      wurstmeister/zookeeper         latest            3f43f72cb283   486 MB
cb-ladybug                  cloudbaristaorg/cb-ladybug     v0.3.0-espresso   a8351e6ea963   29.2 MB
cb-restapigw                cloudbaristaorg/cb-restapigw   v0.3.0-espresso   119daf1d457e   96.3 MB
cb-restapigw-grafana        grafana/grafana                latest            c9e576dccd68   189 MB
cb-restapigw-influxdb       influxdb                       latest            bd69ea12fb63   270 MB
cb-restapigw-jaeger         jaegertracing/all-in-one       latest            d369432efee6   49.2 MB
cb-spider                   cloudbaristaorg/cb-spider      v0.3.0-espresso   00bf045c9748   208 MB
cb-tumblebug                cloudbaristaorg/cb-tumblebug   v0.3.0-espresso   76332875c917   113 MB
cb-tumblebug-phpliteadmin   acttaiwan/phpliteadmin         latest            f5242ee12570   78.9 MB
```
</details>

## Cloud-Barista 중지
```bash
./operator stop
```
