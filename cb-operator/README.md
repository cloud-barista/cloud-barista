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
cb-operator/src$ go build -o operator
```

# Commands to use the operator

## Help
```
~/go/src/github.com/cloud-barista/cb-operator/src$ ./operator 

The operator is a tool to operate Cloud-Barista system. 
  
  For example, you can setup and run, stop, and ... Cloud-Barista runtimes.
  
  - ./operator pull [-f ../docker-compose.yaml]
  - ./operator run [-f ../docker-compose.yaml]
  - ./operator info
  - ./operator exec -t cb-tumblebug -c "ls -al"
  - ./operator stop [-f ../docker-compose.yaml]
  - ./operator remove [-f ../docker-compose.yaml] -v -i

Usage:
  operator [command]

Available Commands:
  exec        Run commands in a target component of Cloud-Barista System
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

## Exec
```
~/go/src/github.com/cloud-barista/cb-operator/src$ ./operator exec -h

Run commands in your components of Cloud-Barista System. 
	For instance, you can get an interactive prompt of cb-tumblebug by
	[operator exec cb-tumblebug sh]

Usage:
  operator exec [flags]

Flags:
  -c, --command string   Command to excute
  -h, --help             help for exec
  -t, --target string    Name of CB component to command

Global Flags:
      --config string   config file (default is $HOME/.operator.yaml)

```
