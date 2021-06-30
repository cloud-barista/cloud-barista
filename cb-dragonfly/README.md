
# CB-Dragonfly
Cloud-Barista Integrated Monitoring Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/cloud-barista/cb-dragonfly)](https://goreportcard.com/report/github.com/cloud-barista/cb-dragonfly)
[![Build](https://img.shields.io/github/workflow/status/cloud-barista/cb-dragonfly/Build%20amd64%20container%20image)](https://github.com/cloud-barista/cb-dragonfly/actions?query=workflow%3A%22Build+amd64+container+image%22)
[![Top Language](https://img.shields.io/github/languages/top/cloud-barista/cb-dragonfly)](https://github.com/cloud-barista/cb-dragonfly/search?l=go)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cloud-barista/cb-dragonfly?label=go.mod)](https://github.com/cloud-barista/cb-dragonfly/blob/master/go.mod)
[![Repo Size](https://img.shields.io/github/repo-size/cloud-barista/cb-dragonfly)](#)
[![GoDoc](https://godoc.org/github.com/cloud-barista/cb-dragonfly?status.svg)](https://pkg.go.dev/github.com/cloud-barista/cb-dragonfly@master)
[![Release Version](https://img.shields.io/github/v/release/cloud-barista/cb-dragonfly?color=blue)](https://github.com/cloud-barista/cb-dragonfly/releases/latest)
[![License](https://img.shields.io/github/license/cloud-barista/cb-dragonfly?color=blue)](https://github.com/cloud-barista/cb-dragonfly/blob/master/LICENSE)

```
[NOTE]
CB-Dragonfly is currently under development. (the latest version is 0.3 espresso)
So, we do not recommend using the current release in production.
Please note that the functionalities of CB-Dragonfly are not stable and secure yet.
If you have any difficulties in using CB-Dragonfly, please let us know.
(Open an issue or Join the cloud-barista Slack)
```

***

## [목차]

1. [설치 개요](#설치-개요)
2. [설치 절차](#설치-절차)
3. [설치 & 실행 상세 정보](#설치--실행-상세-정보)
4. [CB-Dragonfly 기능 사용 방법](#cb-dragonfly-기능-사용-방법)
***


## [설치 개요]
- 설치 환경: 리눅스(검증시험:Ubuntu 18.04)

## [설치 절차]
- Git 설치
- Go 설치
- Go 환경 변수 설정 
- Docker/ Docker-Compose 설치 


## [설치 & 실행 상세 정보]

- Git 설치
  - `$ sudo apt update`
  - `$ sudo apt install git`
  - `$ sudo apt-get install git-core`

- Go 설치
  - https://golang.org/doc/install 
  (2020년 11월 현재 `$ sudo apt install golang` 으로 설치하면 1.10 설치됨. 이 링크에서 1.15 이상 버전으로 설치할 것(Go mod 호환성 문제))
  - `$ wget https://golang.org/dl/go1.15.4.linux-amd64.tar.gz` (설치 파일 다운로드)
  - `$ sudo tar -C /usr/local -xzf go1.15.4.linux-amd64.tar.gz` (압축해제)
  
- Go 환경 변수 설정
  - `$ sudo echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc` (GOROOT{/usr/local/go/bin}를 PATH 환경 변수에 추가하여 ~/.bashrc 맨 아래줄에 추가)
  - `$ source ~/.bashrc` (수정한 bashrc 파일 반영)
  - `$ go version` (GO 버전 확인)
  
 - Docker/ Docker-compose 설치
   - https://docs.docker.com/engine/install 참고
   - https://docs.docker.com/compose/install 참고
  
- 멀티 클라우드 모니터링 프레임워크 (cb-dragonfly) 설치

    - Git Project Clone

          - `$ sudo git config --global color.ui "auto"` (Git 소스에 색 구분)
          - `$ sudo git clone https://github.com/cloud-barista/cb-dragonfly.git` (Git 프로젝트 CLone)
          - `username = {{GitUserEmail}}` (Clone시 자격여부 확인 : 자신의 Git Email 입력)
          - `Password = {{GitUserPW}}`    (Clone시 자격여부 확인 : 자신의 Git PW 입력)
    
    - Go mod 기반 의존성 라이브러리 로드
          
          - `$ cd ~/cb-dragonfly` (clone한 프로젝트 파일로 들어가기)
          - `$ go mod download` (.mod 파일에 등록된 라이브러리 다운로드 실행)
    
    - Go mod 기반 의존성 라이브러리 다운로드 확인
    
          - `$ go mod verify` (다운로드 확인)
    
    - 환경변수 설정
          
          - `$ sudo vim conf/setup.env` (실행에 필요한 PATH를 처리할 파일 생성  (현 위치: ~/cb-dragonfly))
               setup.env에 추가
                
                export CBSTORE_ROOT=~/cb-dragonfly
                export CBLOG_ROOT=~/cb-dragonfly
                export CBMON_ROOT=~/cb-dragonfly
                
          - `$ source setup.env` (수정한 setup.env 반영)         
          - `$ go run pkg/manager/main/main.go` (실행)
    
    - config.yaml 파일 설정 (conf/config.yaml 파일에 cb-dragonfly 호스트 IP ( kafka IP 및 collector IP ) 정보, 배포 환경, 모니터링 정책을 순차적으로 입력)
          
          -  #### Config for cb-dragonfly ####
             
             # influxdb connection info
             influxdb:
               endpoint_url: http://cb-dragonfly-influxdb           # endpoint for influxDB
               internal_port: 8086
               external_port: 28086
               database: {{ database_name }}
               user_name: {{ user_name }}
               password: {{ password }}
             
             kapacitor:
               endpoint_url: http://cb-dragonfly-kapacitor:9092     # endpoint to kapacitor
             
             kafka:
               endpoint_url: cb-dragonfly-kafka
               external_ip: {{ external_ip }}
               deploy_type: {{ deploy_type }}                       # deploy environment "compose" => docker-compose or others , "helm" => helm chart on k8s
               compose_external_port: 9092
               helm_external_port: 32000
               internal_port: 9092
             
             # collect manager configuration info
             collectManager:
               collector_ip: {{ collector_ip }}                     # local access endpoint to cb-dragonfly API server
               collector_port: 8094                                 # udp port
               collector_group_count: 1                             # default collector group count
             
             # api server configuration info
             apiServer:
               port: 9090
             
             # monitoring interval configuration info
             monitoring:
               agent_interval: 2                                    # agent interval (s)
               collector_interval: 10                               # aggregate interval (s)
               max_host_count:  5                                   # maximum host count per collector
               monitoring_policy: {{ monitoring_policy }}           # "agentCount" => The number of agent, "csp" => csp group
             
             grpcServer:
               port: 9999


- 멀티 클라우드 모니터링 프레임워크(cb-dragonfly) Docker-Compose 실행

     - 프로젝트 빌드 및 실행 

              - `$ cd ~/cb-dragonfly`
              - `$ sudo make compose-up` (cb-dragonfly 프로젝트 빌드 및 실행)

- 멀티 클라우드 모니터링 프레임워크(cb-dragonfly) 정지

     - 프로젝트 중지 및 삭제 

              - `$ cd ~/cb-dragonfly`
              - `$ sudo make compose-rm` (cb-dragonfly 중지 및 삭제)

## [CB-Dragonfly 기능 사용 방법]

### (1) CB-Dragonfly의 REST API를 사용하여 테스트

- REST API 정의 문서: [링크](https://documenter.getpostman.com/view/7454078/TzJu8wwi)
- 활용 기능
  - 에이전트 설치/삭제
  - MCIS 모니터링 정보 조회
  - VM 모니터링 정보 조회
  - 온디멘드 정보 조회
  - 에이전트 메타데이터 조회
  - 모니터링 정책 설정
  - 알람 등록 및 관리
