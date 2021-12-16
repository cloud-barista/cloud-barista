# cb-operator

[![Go Report Card](https://goreportcard.com/badge/github.com/cloud-barista/cb-operator)](https://goreportcard.com/report/github.com/cloud-barista/cb-operator)
[![Repo Size](https://img.shields.io/github/repo-size/cloud-barista/cb-operator)](#)
[![code size](https://img.shields.io/github/languages/code-size/cloud-barista/cb-operator)](#)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cloud-barista/cb-operator?label=go.mod)](https://github.com/cloud-barista/cb-operator/blob/master/go.mod)
[![GoDoc](https://godoc.org/github.com/cloud-barista/cb-operator?status.svg)](https://pkg.go.dev/github.com/cloud-barista/cb-operator@master)

[![issues](https://img.shields.io/github/issues/cloud-barista/cb-operator)](https://github.com/cloud-barista/cb-operator/issues)
[![issues](https://img.shields.io/github/issues-closed/cloud-barista/cb-operator)](https://github.com/cloud-barista/cb-operator/issues?q=is%3Aissue+is%3Aclosed)
[![pull requests](https://img.shields.io/github/issues-pr/cloud-barista/cb-operator)](https://github.com/cloud-barista/cb-operator/pulls)
[![pull requests](https://img.shields.io/github/issues-pr-closed/cloud-barista/cb-operator)](https://github.com/cloud-barista/cb-operator/pulls?q=is%3Apr+is%3Aclosed)

[![Release Version](https://img.shields.io/github/v/release/cloud-barista/cb-operator)](https://github.com/cloud-barista/cb-operator/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/cloud-barista/cb-operator/blob/master/LICENSE)

The Operation Tool for Cloud-Barista System Runtime

```
[NOTE]
cb-operator is currently under development.
So, we do not recommend using the current release in production.
Please note that the functionalities of cb-operator are not stable and secure yet.
If you have any difficulties in using cb-operator, please let us know.
(Open an issue or Join the cloud-barista Slack)
```

```
CB-Ladybug, which is included in cb-operator's Cloud-Barista deployment shape, is currently under development.
So, we do not recommend using the current release in production.
Please note that the functionalities of CB-Ladybug are not stable and secure yet.
```

## cb-operator 개요
- Cloud-Barista 시스템의 실행, 상태정보 제공, 종료 등을 지원하는 관리 도구 입니다.
- 2가지 모드를 제공합니다.
  - [Docker Compose 모드](docs/cb-operator-docker-compose-mode.md)
  - [Kubernetes 모드](docs/cb-operator-k8s-mode.md)
- Cloud-Barista 실행 환경에 따른 모드 선택 가이드
  - Cloud-Barista를 단일 노드에서, Kubernetes 클러스터 없이 간단히 실행하고 싶다 → [Docker Compose 모드](docs/cb-operator-docker-compose-mode.md)
  - Cloud-Barista를 단일 노드에서, Kubernetes 클러스터 상에 실행하고 싶다 → [Kubernetes 모드](docs/cb-operator-k8s-mode.md)
  - 다중 노드로 구성된 Kubernetes 클러스터에 Cloud-Barista를 실행하고 싶다 → [Kubernetes 모드](docs/cb-operator-k8s-mode.md)

## FAQ

1. Container conflict (Ref: https://github.com/cloud-barista/cb-operator/issues/75)

```
# ./operator run
CB_OPERATOR_MODE: DockerCompose

[Setup and Run Cloud-Barista]

[Config path] ../docker-compose-mode-files/docker-compose.yaml

Creating cb-dragonfly-influxdb ...
Creating cb-dragonfly-etcd ...
Creating cb-tumblebug-phpliteadmin ...
Creating cb-restapigw-jaeger ...
Creating cb-restapigw-influxdb ...
Creating cb-dragonfly-etcd
Creating cb-dragonfly-influxdb
Creating cb-tumblebug-phpliteadmin
Creating cb-restapigw-jaeger
Creating cb-dragonfly-influxdb ... error

ERROR: for cb-dragonfly-influxdb  Cannot create container for service cb-dragonfly-influxdb: Conflict. The container name "/cb-dragonfly-influxdb" is already in use by container "2d15cca0ed2f399ff3155648538Creating cb-tumblebug-phpliteadmin ... error

ERROR: for cb-tumblebug-phpliteadmin  Cannot create container for service cb-tumblebug-phpliteadmin: Conflict. The container name "/cb-tumblebug-phpliteadmin" is already in use by container "6947d55ae3318f3Creating cb-dragonfly-etcd ... error

ERROR: for cb-dragonfly-etcd  Cannot create container for service cb-dragonfly-etcd: Conflict. The container name "/cb-dragonfly-etcd" is already in use by container "c7800abcc11a333f4dcd79eb1527d20a2efc311Creating cb-restapigw-influxdb ... error

ERROR: for cb-restapigw-influxdb  Cannot create container for service influxdb: Conflict. The container name "/cb-restapigw-influxdb" is already in use by container "efc1fb276586b04527c02122dcac736a3d09e100Creating cb-restapigw-jaeger ... error

ERROR: for cb-restapigw-jaeger  Cannot create container for service jaeger: Conflict. The container name "/cb-restapigw-jaeger" is already in use by container "925c76c42780b892036e34b2183fc0c36fabb47e835eda48e4bbc9b83270205d". You have to remove (or rename) that container to be able to reuse that name.

ERROR: for influxdb  Cannot create container for service influxdb: Conflict. The container name "/cb-restapigw-influxdb" is already in use by container "efc1fb276586b04527c02122dcac736a3d09e100bc8c5e9b4e944c0a1a3012eb". You have to remove (or rename) that container to be able to reuse that name.

ERROR: for cb-tumblebug-phpliteadmin  Cannot create container for service cb-tumblebug-phpliteadmin: Conflict. The container name "/cb-tumblebug-phpliteadmin" is already in use by container "6947d55ae3318f3f07c969a1a60d928c3b77e48829e5aaff675b5b8001cc27da". You have to remove (or rename) that container to be able to reuse that name.

ERROR: for jaeger  Cannot create container for service jaeger: Conflict. The container name "/cb-restapigw-jaeger" is already in use by container "925c76c42780b892036e34b2183fc0c36fabb47e835eda48e4bbc9b83270205d". You have to remove (or rename) that container to be able to reuse that name.

ERROR: for cb-dragonfly-etcd  Cannot create container for service cb-dragonfly-etcd: Conflict. The container name "/cb-dragonfly-etcd" is already in use by container "c7800abcc11a333f4dcd79eb1527d20a2efc3116f7c36750d5a366b9a8fa07a6". You have to remove (or rename) that container to be able to reuse that name.

ERROR: for cb-dragonfly-influxdb  Cannot create container for service cb-dragonfly-influxdb: Conflict. The container name "/cb-dragonfly-influxdb" is already in use by container "2d15cca0ed2f399ff31556485385ef577f82eb41d81bc777c4bd71e0a97912da". You have to remove (or rename) that container to be able to reuse that name.
Encountered errors while bringing up the project.
```

Solution:
Remove conflicted containers
```
docker rm 925c76c42780b892036e34b2183fc0c36fabb47e835eda48e4bbc9b83270205d
```



