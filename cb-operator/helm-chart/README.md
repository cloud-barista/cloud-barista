# Cloud-Barista Platform
---
CCloud-Barista는 멀티 클라우드 서비스 및 솔루션 개발에 필요한 공통 소프트웨어 기술입니다.


## TL;DR;
---

//TODO
```
# helm repo add cloudbarista https://xxxxx
# helm install release cloudbarista/xxxxxxx
```


## Introduction
---

이 차트는 Helm 패키지 관리자를 사용하여 Kubernetes 클러스터에 Cloud-Barista Platform을 배포하고 시작합니다.

## Prerequisites
---

* Kubernetes 1.14+
* Helm 3.0+
* PV provisioner support in the underlying infrastructure
* ReadWriteMany volumes for docker-registry

## Installing the Chart
---

릴리즈명 `release` 로 설치합니다.

```
▒ helm install release .
```

이 명령은 기본 구성으로 Kubernetes 클러스터에 Cloud-Barista를 배포합니다. 
자세한 Parameter는 아래를 참조합니다.


## Uninstalling the Chart
---
릴리즈명 `release` 를 삭제합니다.

```
▒ helm delete release
```

## Parameters
---
아래 표는 Cloud-Barista 차트의 구성 가능한 매개 변수 및 기본값을 서술합니다.

 Parameter                               | Description                                                                                | Default                    |
|----------------------------------------|--------------------------------------------------------------------------------------------|----------------------------|
| `cb-dragonfly.enabled`                 | cb-drafonfly 설치 여부                                                                     | `true`                     |
| `cb-restapigw.enabled`                 | cb-restapigw 설치 여부                                                                     | `true`                     |
| `cb-spider.enabled`                    | cb-spider 설치 여부                                                                        | `true`                     |
| `cb-tumblebug.enabled`                 | cb-tumblebug 설치 여부                                                                     | `true`                     |
| `cb-webtool.enabled`                   | cb-webtool 설치 여부                                                                       | `true`                     |
| `docker-registry.enabled`              | docker-registry 설치 여부                                                                  | `true`                     |
| `docker-registry.tlsSecretName`        | HTTPS를 위한 TLS secret 을 지정, nil 이면 HTTP                                             | `nil`                      |
| `docker-registry.persistence.enabled`  | persistence valumn 지정여부                                                                | `false`                    |

각 parameter는  helm install 시 --set key=value[,key=value] argument 로 지정할 수 있습니다. 아래예 참조

```
▒ helm install release . \
  --set docker-registry.enabled=false
```

또한 아래와 같이 Parameter 값을 지정하는 YAML 파일를 활용할 수도 있습니다.

```
▒ helm install elease -f values.yaml .
```

각 컴포넌트에서 사용하는 dependency sub-chart 들의 파라메터들은 차트 홈페이지를 참조할 수 있습니다.


### Dependency sub-charts

* cb-dragon-fly
  * bitnami/influxdb : https://github.com/bitnami/charts/tree/master/bitnami/influxdb
  * bitnami/etcd : https://github.com/bitnami/charts/tree/master/bitnami/etcd

* cb-restapigw
  * bitnami/influxdb : https://github.com/bitnami/charts/tree/master/bitnami/influxdb
  * stable/grafana : https://github.com/helm/charts/tree/master/stable/grafana
  * jaegertracing/jaeger : https://github.com/jaegertracing/helm-charts

* docker-registry
  * stable/docker-registry : https://github.com/helm/charts/tree/master/stable/docker-registry


## Configuration and installation details
---

### cb-webtool  NodePort 로 지정

```
▒ helm install release . \
  --set cb-webtool.service.type=NodePort
```

`curl`로 NodePort 동작 여부를 확인합니다.

```
▒ IP="<Cluster IP>"  # cluster ip를 지정합니다.
▒ curl -i -s http://${IP}:$(kubectl get svc/cb-webtool -o jsonpath='{.spec.ports[?(@.name=="http")].nodePort}')
```


## Private docker registry 설정


### 서비스를 NodePort 로 지정하고 Plain HTTP 프로토콜 사용

private docker-registry 외부에서 NodePort 로 접근할 수 있도록  `docker-registry.service.type=NodPort`로 값을 지정합니다
아래 예제는 `nodePort=30500` 도 임의로 지정해 주었습니다.

```
▒ helm install release . \
  --set docker-registry.service.type=NodePort \
  --set docker-registry.service.nodePort=30500 
```

docker image pull/push 는 기본적으로 `https` 프로토콜을 사용하도록 되어 있습니다.
위 예제에서는 NodePort 30500은 TLS 지정을 하지 않았으므로 http 프로토콜을 사용하게 됩니다.
그러므로 외부에서 docker-registry 에 pull/push를 하기 위해서는 클라이언트 docker daemon 설정 파일 `daemon.json` 에  `insecure-registries` 속성에   `IP:PORT`  값을 추가해 주어야 합니다.
또한 K8s 클러스터 내부에서 private image 를  pull 하기 위해서는 클러스터의 Node들에도 `insecure-registries` 설정을 해주어야 합니다.
[공식문서 - Deploy a plain HTTP registry](https://docs.docker.com/registry/insecure/#deploy-a-plain-http-registry) 참조합니다.

```
{
  "debug": true,
  "experimental": false,
  "insecure-registries": [
    "101.55.xxx.xxx:30500"
  ]
}
```


`docker push` 를 통해 클러스터 외부에서 정상 동작을 확인합니다.

```
# variables
▒ IP="<Cluster IP>"                                                                                 # cluster endpoint ip
▒ HOST=${IP}:$(kubectl get svc/cb-webtool -o jsonpath='{.spec.ports[?(@.name=="http")].nodePort}')  # service endpoint - IP:PORT

# image push
▒ docker pull honester/httpbin
▒ docker image tag honester/httpbin ${HOST}/httpbin
▒ docker push ${HOST}/httpbin

# image push verify
▒ curl -s http://${HOST}/v2/_catalog

# clean-up image
▒ docker rmi ${HOST}/httpbin
```


### 서비스를 NodePort 로 지정하고 Self-Sign 인증서로 HTTPS 프로토콜 사용

self-sign 인증서를 생성합니다. `server.key`, `server.crt` 가 생성됩니다.

```
# variables
▒ KEY_FILE=server.key                   # key 파일명
▒ CERT_FILE=server.crt                  # certificate 파일명
▒ CERT_NAME=docker-registry-tls-cert    # k8s TLS secret 이름
▒ IP="<Cluster IP>"                     # cluster endpoint ip 
▒ HOST=${IP}:30500                      # service endpoint -  IP:PORT

# 인증서 생성
▒ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ${KEY_FILE} -out ${CERT_FILE} -subj "/CN=${IP}/O=${IP}" -reqexts SAN -extensions SAN -config <(cat /etc/ssl/openssl.cnf <(printf "[SAN]\nsubjectAltName=IP:${IP}"))
```

생성한 인증서로 k8s 클러스터에 TLS secret 생성합니다.

```
▒ kubectl create secret tls ${CERT_NAME} --key ${KEY_FILE} --cert ${CERT_FILE}
```


`docker-registry.tlsSecretName` 파라메터에 생성된 TLS secret 이름을 지정하여 Cloud-Barista 를 설치합니다.

```
▒ helm install release . \
  --set docker-registry.service.type=NodePort \
  --set docker-registry.service.nodePort=30500 \
  --set docker-registry.tlsSecretName="${CERT_NAME}"
```


외부에서 docker-registry 로 NodePort로 HTTPS 프로토콜을 통해 접근하기 위해서  Self-Sign 인증서를 클라이언트에 등록합니다.
각 OS 별로 Self-Sign 인증서를 등록하는 방법이 다릅니다.
아래는 Mac에서 키체인에 시스템 인증서로 등록 하는 방법입니다. (https://docs.docker.com/docker-for-mac/#add-custom-ca-certificates-server-side)

```
▒ sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ${CERT_FILE}
```

`docker push` 를 통해 클러스터 외부에서 정상 동작을 확인합니다.

```
# docker 로그인
▒ docker login $HOST

# 이미지 push
▒ docker pull honester/httpbin
▒ docker image tag honester/httpbin $HOST/httpbin
▒ docker push $HOST/httpbin

# 확인
▒ curl -s https://${HOST}/v2/_catalog

# clean-up image
▒ docker rmi ${HOST}/httpbin
```


### Basic Authentication (htpasswd) 적용

`htpasswd` 를 사용하여 `user.htpasswd` 파일을 생성합니다.

```
# variables
▒ USERNAME=admin	# htpasswd 사용자명
▒ PASSWD=1234		# htpasswd 비밀번호

# htpasswd  파일 생성
▒ docker run --entrypoint htpasswd registry -Bbn ${USERNAME} ${PASSWD} > user.htpasswd
```

`docker-registry.secrets.htpasswd` 파라메터에 htpasswd 값를 지정하여 Cloud-Barista 를 설치합니다.
아래 예제는 HTTP 프로토콜을 활용한 예제입니다. HTTPs 프로토콜에도 동일한 방식으로  적용할 수 있습니다.

```
# Case - NodePort-HTTP 
▒ helm install release . \
  --set docker-registry.service.type=NodePort \
  --set docker-registry.service.nodePort=30500 \
  --set docker-registry.secrets.htpasswd="$(cat user.htpasswd)" 
```

`docker push` 를 통해 클러스터 외부에서 정상 동작을 확인합니다.
```
# variables
▒ IP="<Cluster IP>"                         # cluster endpoint ip
▒ HOST=${IP}:30500                          # service endpoint
▒ PULL_SECRET=docker-registry-pull-secret   # image pull secret 명

# 이미지 push
▒ docker pull honester/httpbin
▒ docker image tag honester/httpbin $HOST/httpbin
▒ docker push $HOST/httpbin

# 확인
▒ curl -s --user ${USERNAME}:${PASSWD} http://${HOST}/v2/_catalog

# clean-up image
▒ docker rmi ${HOST}/httpbin
```

k8s 클러스터 내부에서 private docker registry 의  Basic Authentication을 위해서는 docker-registry secret을 생성하고 해당 secret을 `imagePullSecrets` 스펙에 지정합니다.

```
# docker-registry secret 을 생성
▒  kubectl create secret docker-registry ${PULL_SECRET} --docker-server=${HOST} --docker-username=${USERNAME} --docker-password=${PASSWD} 

# Pod 배포 테스트
▒  kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: test
spec:
  containers:
  - name: httpbin
    image: honester/httpbin
    imagePullPolicy: Always
  imagePullSecrets:
  - name: ${PULL_SECRET}
EOF

# Image pull 확인
▒  kubectl get pod/test

# clean-up 
▒  kubectl delete pod/test
```

### hostPath PersistentVolume 지정

docker-registry Pod가 배포될 노드를 선택하기 위해 특정 노드에 label을 지정합니다.
아래는 NodeList 중 첫번째 노드를 임의로 선택하여 label을 지정하는 예제입니다.

```
# variables
LABEL_NAME="docker-registry"                                         # label 이름
LABEL_VALUE="on"                                                     # label 값
NODE_NAME=$(kubectl get node -o jsonpath={.items[0].metadata.name})  # label 지정할 노드

# set label
kubectl label nodes ${NODE_NAME} ${LABEL_NAME}=${LABEL_VALUE}

# verify
kubectl get nodes -l ${LABEL_NAME}=${LABEL_VALUE}
```

`docker-registry.persistence.enabled` 파라메터를 `true`로 지정하고  `docker-registry.persistence.size` 에 Valume 크기를  지정한 후 Cloud-Barista 를 설치합니다.

```
▒ helm install release . \
  --set docker-registry.service.type=NodePort \
  --set docker-registry.service.nodePort=30500 \
  --set docker-registry.nodeSelector.${LABEL_NAME}=${LABEL_VALUE} \
  --set docker-registry.persistence.enabled=true \
  --set docker-registry.persistence.size=8Gi 
```

### Private Docker Registry 설정 참조
* [Docker Registry Helm Chart](https://github.com/helm/charts/tree/master/stable/docker-registry)
* [Pull an Image from a Private Registry](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/)


## Verify installation
---

### 로그 확인

```
▒ kubectl logs $(kubectl get pod -l app.kubernetes.io/name=cb-dragonfly -o jsonpath={.items..metadata.name})
▒ kubectl logs $(kubectl get pod -l app.kubernetes.io/name=cb-restapigw -o jsonpath={.items..metadata.name})
▒ kubectl logs $(kubectl get pod -l app.kubernetes.io/name=cb-spider -o jsonpath={.items..metadata.name})
▒ kubectl logs $(kubectl get pod -l app.kubernetes.io/name=cb-tumblebug -o jsonpath={.items..metadata.name})
▒ kubectl logs $(kubectl get pod -l app.kubernetes.io/name=cb-webtool -o jsonpath={.items..metadata.name})
```

### 서비스 포트 확인

아래와 같이 검증을 위한 httpbin 파드를 준비(배포) 합니다.

```
▒ kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: httpbin
spec:
  containers:
  - image: docker.io/honester/httpbin:latest
    name: httpbin
EOF
```

httpbin Pod 를 통해 컴포너트 서비스 포트를 검증합니다.

```
▒ kubectl exec httpbin -- curl -i -s http://cb-dragonfly:9090/
▒ kubectl exec httpbin -- curl -i -s http://cb-restapigw:8000/
▒ kubectl exec httpbin -- curl -i -s http://cb-spider:1024/
▒ kubectl exec httpbin -- curl -i -s http://cb-tumblebug:1323/
▒ kubectl exec httpbin -- curl -i -s http://cb-webtool:1234/login
▒ kubectl exec httpbin -- curl -i -s http://docker-registry:5000/v2/
```
