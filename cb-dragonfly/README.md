# cb-dragonfly
Cloud-Barista Integrated Monitoring Framework

***

## [목차]

0. [VM 접속](#vm-접속)
1. [설치 개요](#설치-개요)
2. [설치 절차](#설치-절차)
3. [설치 & 실행 상세 정보](#설치--실행-상세-정보)

***


## [VM 접속]

- CB-Dragonfly.pem 키를 사용 SSH 접속
  - `$ ssh cb-user@0.0.0.0 -i CB-Dragonfly.pem`    vm에 접속 ($Home = /home/cb-user)

## [설치 개요]
- 설치 환경: 리눅스(검증시험:Ubuntu 18.04)

## [설치 절차]

- Go 설치 & Git 설치
- etcd 설치 & influxdb 설치
- 환경 변수 설정

## [설치 & 실행 상세 정보]

- Git 설치
  - `$ sudo apt update`
  - `$ sudo apt install git`
  - `$ sudo apt-get install git-core`

- Go 설치
  - https://golang.org/doc/install 
  (2019년 11월 현재 `$ sudo apt install golang` 으로 설치하면 1.10 설치됨. 이 링크에서 1.12 이상 버전으로 설치할 것(Go-mod 호환성 문제))
  - `$ wget https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz` (설치 파일 다운로드)
  - `$ sudo tar -C /usr/local -xzf go1.13.4.linux-amd64.tar.gz` (압축해제)
  - `$ sudo vim ~/.bashrc 파일 맨 아래에 export GOROOT=$PATH:/usr/local/go/bin` (GOPATH 환경변수 추가)
  - `$ source ~/.bashrc` (수정한 bashrc 파일 반영)
  - `$ go version` (버전 확인)

- 모니터링 데이터베이스 저장소(의존 라이브러리 다운로드)
  - etcd 설치(3.3.11) 및 실행
  
        - `$ wget https://github.com/coreos/etcd/releases/download/v3.3.11/etcd-v3.3.11-linux-amd64.tar.gz` (설치 파일 다운로드)
        - `$ sudo tar -xvf etcd-v3.3.11-linux-amd64.tar.gz` (압축해제)
        - `$ sudo mv etcd-v3.3.11-linux-amd64/etcd* /usr/local/bin/` (추출된 실행파일을 로컬 저장소로 이동)
        - `$ etcd --version` (버전 확인)
    
        - `$ sudo mkdir -p /var/lib/etcd/` (Etcd 구성 파일 폴더 생성)
        - `$ sudo mkdir /etc/etcd` (데이터 폴더 생성)
    
        - `$ sudo groupadd --system etcd` (etcd 시스템 그룹 생성)
        - `$ sudo useradd -s /sbin/nologin --system -g etcd etcd` (etcd 시스템 사용자 생성)
        - `$ sudo chown -R etcd:etcd /var/lib/etcd/` (/var/lib/etcd/ 폴더 소유권을 etcd사용자로 설정)
    
        - `$ sudo vim /etc/systemd/system/etcd.service` (etcd에 대한 새로운 시스템 서비스 파일 작성)
        (바로 밑에 코드 붙여넣기 후)
        - `$ sudo systemctl  daemon-reload` (데몬 재시작)
        - `$ sudo systemctl  start etcd.service` (etcd 서비스 시작)
        
  - etcd.service 붙여넣기
          
          [Unit]
          Description=etcd key-value store
          Documentation=https://github.com/etcd-io/etcd
          After=network.target

          [Service]
          User=etcd
          Type=notify
          Environment=ETCD_DATA_DIR=/var/lib/etcd
          Environment=ETCD_NAME=%m
          ExecStart=/usr/local/bin/etcd
          Restart=always
          RestartSec=10s
          LimitNOFILE=40000

          [Install]
          WantedBy=multi-user.target

  - influxdb (1.7.8) 및 실행
  
        - `$ wget https://dl.influxdata.com/influxdb/releases/influxdb_1.7.8_amd64.deb` (다운로드)
        - `$ sudo dpkg -i influxdb_1.7.8_amd64.deb` (압축해제)
        - `$ sudo apt-get update && sudo apt-get install influxdb` (InfluxDB 서비스 설치)
        - `$ sudo systemctl start influxdb` (influxDB 서비스 시작)
        - `$ influx --version` (버전 확인)
    
        - `$ influx` (influxDB 사용하기)
            - CREATE DATABASE cbmon

- 멀티 클라우드 모니터링 설치

    - Git Project Clone

          - `$ sudo git config --global color.ui "auto"` (Git 소스에 색 구분)
          - `$ sudo git clone https://github.com/cloud-barista/cb-dragonfly.git` (Git 프로젝트 CLone)
          - `username = {{GitUserEmail}}` (Clone시 자격여부 확인 : 자신의 Git Email 입력)
          - `Password = {{GitUserPW}}`    (Clone시 자격여부 확인 : 자신의 Git PW 입력)
    
    - Go mod 의존성 라이브러리 로드
          
          - `$ cd ~/cb-mon` (clone한 프로젝트 파일로 들어가기)
          - `$ go mod download` (go mod 가 있는 폴더에서 다운로드 실행)
    
    - Go mod 의존성 라이브러리 다운로드 확인
    
          - `$ go mod verify` (다운로드 확인)
    
    - 라이브러리 실행
          
          - `$ sudo vim conf/setup.env` (실행에 필요한 PATH를 처리할 파일 생성  (현 위치: ~/cb-mon))
               setup.env에 추가
                
                export CBSTORE_ROOT=~/cb-mon
                export CBLOG_ROOT=~/cb-mon
                export CBMON_PATH=~/cb-mon
                export SPIDER_URL=http://localhost:1024
                
          - `$ source conf/setup.env` (수정한 setup.env 반영)         
          - `$ go run pkg/manager/main/main.go` (실행)
    
    - config 파일 설정
          
          - ``

