# Cloud-Barista
Multi-Cloud Service Common Framework

Cloub-Barista consists of multiple sub-frameworks with the concept of microservice.
Main frameworks are as follows,
1. CB-Spider (connects all clouds in a single interface)
2. CB-Tumblebug (manages multi-cloud resource and integrated infra services)
3. CB-Dragonfly (monitors multi-cloud service)
4. CB-Webtool (provides a GUI to  Cloud-Barista users)
5. CB-Store (provides an unified DB interface for meta info of Cloud-Barista)
6. CB-Log (provides log system)

The frameworks are in src directory.

***

## [목    차]

1. [설치 환경](#설치-환경)
2. [설치 및 실행](#설치-및-실행)
3. [실행 준비](#실행-준비)
4. [API 규격](#API-규격)
5. [활용 예시](#활용-예시)
6. [특이 사항](#특이-사항)

***

## [설치 환경]

- 리눅스(검증시험:Ubuntu 18.04, Raspbian GNU/Linux 10)

## [설치 및 실행]

- Git 설치
- Go 설치(1.12 이상)  

- Cloud-Barista 소스 다운로드 및 설치
  - Cloud-Barista alliance 설치 (CB-Log)
    - cloud-barista/cb-log README를 참고하여 설치
  
  - Cloud-Barista alliance 설치 (CB-Store)
    - cloud-barista/cb-store README를 참고하여 설치

  - CB-Spider 설치 및 실행
    - cloud-barista/cb-spider README를 참고하여 설정, 설치
    - cb-spider 실행 (cb-spider API 서버 실행)

  - CB-Tumblebug 설치 및 실행
    - cloud-barista/cb-tumblebug README를 참고하여 설정, 설치
      - cb-spider API 서버 주소를 cb-tumblebug의 setup.env에 설정
    - cb-tumblebug 실행 (cb-tumblebug API 서버 실행)

  - CB-Dragonfly 설치 및 실행
    - cloud-barista/cb-dragonfly README를 참고하여 설정, 설치
    - cb-dragonfly 실행 (cb-dragonfly API 서버 실행)

  - CB-Webtool 설치 및 실행
    - cloud-barista/cb-webtool README를 참고하여 설정, 설치
    - cb-webtool 실행 (cb-webtool UGI 서버 실행)
  
## [API 규격]
- cloud-barista/docs/API-Specifications/User-REST-API(v0.30).md 참고
  
## [활용 예시]
- cb-spider API를 통해서 클라우드 연동 정보 입력
  - 시험 도구: `cb-spier/api-runtime/rest-runtime/test/[aws|azure|gcp|openstack|cloudit]` (AWS 경우:aws)
  - 시험 순서: 연동 정보 추가 => 자원등록 => VM 생성 및 제어 시험
  - 시험 방법
    - (연동정보관리) cb-spider/api-runtime/rest-runtime/test/aws/cim-insert-test.sh 참고(Credential 정보 수정 후 실행)

- cb-tumblebug API를 통해서 네임스페이스 등록, 멀티 클라우드 인프라 자원(MCIR) 관리, 멀티 클라우드 인프라 서비스(MCIS) 관리 수행
  - 멀티 클라우드 네임스페이스 관리 API를 통해서 Namespace 생성
  - 멀티 클라우드 인프라 자원(MCIR) 관리 API를 통해서 MCIS 및 VM 생성을 위한 자원 생성
  - 멀티 클라우드 인프라 서비스(MCIS) 관리 API를 통해서 MCIS 생성, 조회, 제어, 종료

## [특이 사항]
- 개발상태: 초기 기능 중심 개발 추진 중 / 기술 개발용 / 상용 활용시 보완필요

