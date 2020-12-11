# cb-operator
The Operation Tool for Cloud-Barista System Runtime

```
[NOTE]
cb-operator is currently under development. (the latest version is 0.2 cappuccino)
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

# FAQ

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


# Launched containers and its endpoints

## Docker Compose Mode

| Name | Endpoint for direct access | Endpoint for access via cb-restapigw | Misc. |
|---|---|---|---|
| cb-restapigw | http://{{host}}:8000 |   | Admin: http://{{host}}:8001 <br> ID: admin / PW: test@admin00  |
| cb-restapigw-influxdb | http://{{host}}:8086 |   | 8083: Admin Panel <br> 8086: client-server comm. |
| cb-restapigw-grafana | http://{{host}}:3100 |   | ID: admin / PW: admin |
| cb-restapigw-jaeger | http://{{host}}:16686 |   |   |
| --- |   |   |   |
| cb-spider | http://{{host}}:1024/spider | http://{{host}}:8000/spider | gRPC: http://{{host}}:2048  |
| cb-tumblebug | http://{{host}}:1323/tumblebug | http://{{host}}:8000/tumblebug | gRPC: http://{{host}}:50252  |
| cb-tumblebug-phpliteadmin | http://{{host}}:2015  |   |   |
| cb-webtool | http://{{host}}:1234 |   |   |
| --- |   |   |   |
| cb-dragonfly | http://{{host}}:9090/dragonfly | http://{{host}}:8000/dragonfly | 8094/udp  |
| cb-dragonfly-influxdb | http://{{host}}:28086 |   |   |
| cb-dragonfly-etcd | http://{{host}}:2379 |   | 2379: client communication <br> 2380: server-to-server communication |
| cb-dragonfly-kapacitor | http://{{host}}:9092  |   |   |

## Kubernetes Mode

| Name | Endpoint for direct access | Endpoint for access via cb-restapigw | Misc. |
|---|---|---|---|
| cb-restapigw | http://{{host}}:30080 |   | Admin: http://{{host}}:30081 <br> ID: admin / PW: test@admin00  |
| cb-restapigw-influxdb | - |   | 8083: Admin Panel <br> 8086: client-server comm. |
| cb-restapigw-grafana | - |   | ID: admin / PW: admin |
| cb-restapigw-jaeger | - |   |   |
| --- |   |   |   |
| cb-spider | http://{{host}}:31024/spider | http://{{host}}:30080/spider | gRPC: http://{{host}}:32048  |
| cb-tumblebug | http://{{host}}:31323/tumblebug | http://{{host}}:30080/tumblebug | gRPC: http://{{host}}:30252  |
| cb-webtool | http://{{host}}:31234 |   |   |
| --- |   |   |   |
| cb-dragonfly | http://{{host}}:30090/dragonfly | http://{{host}}:30080/dragonfly | 30094/udp  |
| cb-dragonfly-influxdb | - |   |   |
| cb-dragonfly-etcd | - |   | 2379: client communication <br> 2380: server-to-server communication |
| --- |   |   |   |
| prometheus | -  |   |   |
| grafana | http://{{host}}:30300  |   | ID: admin / PW: admin  |

# Prerequisites

## Install Docker
- [Install Docker Engine on Ubuntu](https://docs.docker.com/engine/install/ubuntu/)
- Tested version: Docker-CE 19.03.6

<details>
  <summary>Docker version details</summary>
  
```
Client: Docker Engine - Community
 Version:           19.03.6
 API version:       1.40
 Go version:        go1.12.16
 Git commit:        369ce74a3c
 Built:             Thu Feb 13 01:27:49 2020
 OS/Arch:           linux/amd64
 Experimental:      true

Server: Docker Engine - Community
 Engine:
  Version:          19.03.6
  API version:      1.40 (minimum version 1.12)
  Go version:       go1.12.16
  Git commit:       369ce74a3c
  Built:            Thu Feb 13 01:26:21 2020
  OS/Arch:          linux/amd64
  Experimental:     false
 containerd:
  Version:          1.2.13
  GitCommit:        7ad184331fa3e55e52b890ea95e65ba581ae3429
 runc:
  Version:          1.0.0-rc10
  GitCommit:        dc9208a3303feef5b3839f4323d9beb36df0a9dd
 docker-init:
  Version:          0.18.0
  GitCommit:        fec3683
```
</details>

## Install Docker Compose
- On Ubuntu/Debian: `sudo apt install docker-compose`
- Tested version: 1.17.1

# Command to build the operator from souce code
```Shell
$ git clone https://github.com/cloud-barista/cb-operator.git

$ cd cb-operator/src

(Setup dependencies)
cb-operator/src$ go get -u

(Build a binary for cb-operator)
cb-operator/src$ go build -o operator
```

# Commands to use the operator

## Help
```
~/go/src/github.com/cloud-barista/cb-operator/src$ ./operator 

The operator is a tool to operate Cloud-Barista system. 
  
  For example, you can setup and run, stop, and ... Cloud-Barista runtimes.
  
  - ./operator pull [-f ../docker-compose-mode-files/docker-compose.yaml]
  - ./operator run [-f ../docker-compose-mode-files/docker-compose.yaml]
  - ./operator info
  - ./operator stop [-f ../docker-compose-mode-files/docker-compose.yaml]
  - ./operator remove [-f ../docker-compose-mode-files/docker-compose.yaml] -v -i

Usage:
  operator [command]

Available Commands:
  help        Help about any command
  info        Get information of Cloud-Barista System
  pull        Pull images of Cloud-Barista System containers
  remove      Stop and Remove Cloud-Barista System
  run         Setup and Run Cloud-Barista System
  stop        Stop Cloud-Barista System

Flags:
      --config string   config file (default is $HOME/.operator.yaml)
  -h, --help            help for operator
  -t, --toggle          Help message for toggle

Use "operator [command] --help" for more information about a command.
```

## Run
```
~/go/src/github.com/cloud-barista/cb-operator/src$ ./operator run -h

Setup and Run Cloud-Barista System

Usage:
  operator run [flags]

Flags:
  -f, --file string   Path to Cloud-Barista Docker-compose file (default "*.yaml")
  -h, --help          help for run

Global Flags:
      --config string   config file (default is $HOME/.operator.yaml)
```

## Stop
```
~/go/src/github.com/cloud-barista/cb-operator/src$ ./operator stop -h

Stop Cloud-Barista System

Usage:
  operator stop [flags]

Flags:
  -f, --file string   Path to Cloud-Barista Docker-compose file (default "*.yaml")
  -h, --help          help for stop

Global Flags:
      --config string   config file (default is $HOME/.operator.yaml)
```
