# Cloud-Barista

[![Release Version](https://img.shields.io/github/v/release/cloud-barista/cb-tumblebug?color=blue)](https://github.com/cloud-barista/cloud-barista/releases/latest)
[![Pre Release Version](https://img.shields.io/github/v/release/cloud-barista/cb-tumblebug?color=brightgreen&include_prereleases&label=release%28dev%29)](https://github.com/cloud-barista/cloud-barista/releases)
[![License](https://img.shields.io/github/license/cloud-barista/cb-tumblebug?color=blue)](https://github.com/cloud-barista/cb-tumblebug/blob/main/LICENSE)
[![Slack](https://img.shields.io/badge/Slack-Cloud--Barista-brightgreen)](https://join.slack.com/t/cloud-barista/shared_invite/zt-bda8zhkg-tlOCr7_TdQGE_oUSz4mlkA)

*The Cloud-Barista is a Multi-Cloud Service Platform SW.* 

Cloud-Barista consists of multiple frameworks (sub-systems) to accommodate microservice-like architecture. 

Please take a look [Cloud-Barista Website](https://cloud-barista.github.io/technology/) for a detail decription.

<details>
<summary>Note for developing and using Cloud-Barista</summary>

#### Development stage of Cloud-Barista
```
Cloud-Barista is currently under development. (not v1.0 yet)
We welcome any new suggestions, issues, opinions, and controbutors !
Please note that the functionalities of Cloud-Barista are not stable and secure yet.
Becareful if you plan to use the current release in production.
If you have any difficulties in using Cloud-Barista, please let us know.
(Open an issue or Join the Cloud-Barista Slack)
```

#### Localization and Globalization of CB-Tumblebug (CB-Tumblebug의 현지화 및 세계화)
```
[English] As an opensource project initiated by Korean members, 
we would like to promote participation of Korean contributors during initial stage of this project. 
So, CB-Tumblebug Repo will accept use of Korean language in its early stages.
On the other hand, we hope this project flourishes regardless of contributor's country eventually.
So, the maintainers recommend using English at least for the title of Issues, Pull Requests, and Commits, 
while CB-Tumblebug Repo accommodates local languages in the contents of them.
```

```
[한국어] CB-Tumblebug은 한국에서 시작된 오픈 소스 프로젝트로서 
프로젝트의 초기 단계에는 한국 기여자들의 참여를 촉진하고자 합니다. 
따라서 초기 단계의 CB-Tumblebug는 한국어 사용을 받아 들일 것입니다.
다른 한편으로, 이 프로젝트가 국가에 관계없이 번성하기를 희망합니다.
따라서 개발 히스토리 관리를 위해 이슈, 풀 요청, 커밋 등의 
제목에 대해서는 영어 사용을 권장하며, 내용에 대한 한국어 사용은 수용할 것입니다.
```

</details>


---
*This repository is an integrated archive for repository of major frameworks.* These repositories are included and listed in the root directory. This repo reflects the latest release only.

Main frameworks or tools are as follow,

- **CB-Spider** (connects all clouds in a single interface)
  - Upstream repo: https://github.com/cloud-barista/cb-spider
- **CB-Tumblebug** (manages multi-cloud infrastructures)
  - Upstream repo: https://github.com/cloud-barista/cb-tumblebug
- **CB-MCKS** (manages multi-cloud Kubernetes clusters)
  - Upstream repo: https://github.com/cloud-barista/cb-mcks
- **CB-Ladybug** (manages multi-cloud applications)
  - Upstream repo: https://github.com/cloud-barista/cb-ladybug
- **CB-Dragonfly** (monitors multi-cloud services)
  - Upstream repo: https://github.com/cloud-barista/cb-dragonfly
- CB-Waterstrider/**cb-webtool** (provides Web GUI to Cloud-Barista users)
  - Upstream repo: https://github.com/cloud-barista/cb-webtool
- CB-Bridge/**cb-operator** (operation tool for Cloud-Barista system runtime)
  - Upstream repo: https://github.com/cloud-barista/cb-operator
- CB-Bridge/**cb-store** (provides an unified DB interface for meta info of Cloud-Barista)
  - Upstream repo: https://github.com/cloud-barista/cb-store
- CB-Bridge/**cb-log** (provides log library to Cloud-Barista system)
  - Upstream repo: https://github.com/cloud-barista/cb-log

**[Note]** CB-Larva is a special repository that incubates (research and develop) new Multi-Cloud technologies. 
CB-Larva explores interesting ideas and shows the possibility of those (i.e., Proof of Concept (POC)). 
That's why we encourage you to take a look and contribute to the special repository.
Please note that the source code of CB-Larva would not be released and archived in this repository for the time being. 

- CB-Larva/cb-cladnet (POC for the cloud adaptive network)
  - Upstream repo: https://github.com/cloud-barista/cb-larva/tree/main/poc-cb-net


***

## [목    차]

1. [실행 환경](#실행-환경)
1. [설치 및 실행](#설치-및-실행)
1. [사용 방법 및 예시](#사용-방법-및-예시)
1. [API 및 문서](#api-및-문서)
1. [특이 사항](#특이-사항)

***

## [실행 환경]

- Linux (추천: Ubuntu v18.04)

***

## [설치 및 실행]

- Cloud-Barista 플랫폼 통합 실행 (Docker 이미지 기반)
  - cb-operator 를 통해 Cloud-Barista 전체 FW를 통합 실행할 수 있음
    - 참고: [cloud-barista/cb-operator](/cb-operator/)

- Cloud-Barista 플랫폼 개별 FW 소스 다운로드 및 설치

  - CB-Spider 설치 및 실행
    - [cloud-barista/cb-spider README를 참고하여 설정, 설치](/cb-spider/)
    - cb-spider 실행 (cb-spider API 서버 실행)

  - CB-Tumblebug 설치 및 실행
    - [cloud-barista/cb-tumblebug README를 참고하여 설정, 설치](/cb-tumblebug/)
      - CB-Spider API 서버 주소를 conf/setup.env 에 설정
      - CB-Dragonfly API 서버 주소를 conf/setup.env 에 설정
    - cb-tumblebug 실행 (cb-tumblebug API 서버 실행)

  - CB-Ladybug 설치 및 실행
    - [cloud-barista/cb-ladybug README를 참고하여 설정, 설치](/cb-ladybug/)
      - CB-Spider API 서버 주소를 conf/setup.env 에 설정
      - CB-Tumblebug API 서버 주소를 conf/setup.env 에 설정
    - cb-ladybug 실행 (cb-ladybug API 서버 실행)

  - CB-Dragonfly 설치 및 실행
    - [cloud-barista/cb-dragonfly README를 참고하여 설정, 설치](/cb-dragonfly/)
    - cb-dragonfly 실행 (cb-dragonfly API 서버 실행)

  - cb-webtool 설치 및 실행
    - [cloud-barista/cb-webtool README를 참고하여 설정, 설치](/cb-webtool/)
      - CB-Spider API 서버 주소를 conf/setup.env 에 설정
      - CB-Tumblebug API 서버 주소를 conf/setup.env 에 설정
      - CB-Dragonfly API 서버 주소를 conf/setup.env 에 설정
    - cb-webtool 실행 (cb-webtool GUI 서버 실행)

***

## [사용 방법 및 예시]

### 주요 서비스 1) 멀티 클라우드 인프라 서비스 (MCIS)
- 멀티 클라우드 인프라 서비스 환경 구성 
  - [CB-Tumblebug 설정 및 실행](https://github.com/cloud-barista/cb-tumblebug#cb-tumblebug-%EC%86%8C%EC%8A%A4-%EB%B9%8C%EB%93%9C-%EB%B0%8F-%EC%8B%A4%ED%96%89-%EB%B0%A9%EB%B2%95-%EC%83%81%EC%84%B8)
    - CB-Spider (필수)
    - CB-Dragonfly (MCIS 모니터링, CB-Tumblebug 자동 제어 기능에 필요)
    - cb-webtool (Web기반 GUI)
- [멀티 클라우드 인프라 서비스 사용 방법](https://github.com/cloud-barista/cb-tumblebug#cb-tumblebug-%EA%B8%B0%EB%8A%A5-%EC%82%AC%EC%9A%A9-%EB%B0%A9%EB%B2%95)
- [멀티 클라우드 인프라 유스케이스](https://github.com/cloud-barista/cb-tumblebug/blob/main/README.md#3-%EB%A9%80%ED%8B%B0-%ED%81%B4%EB%9D%BC%EC%9A%B0%EB%93%9C-%EC%9D%B8%ED%94%84%EB%9D%BC-%EC%9C%A0%EC%8A%A4%EC%BC%80%EC%9D%B4%EC%8A%A4)

### 주요 서비스 2) 멀티 클라우드 쿠버네티스 서비스 (MCKS)
- 멀티 클라우드 쿠버네티스 서비스 환경 구성 
  - [CB-Ladybug 설정 및 실행](https://github.com/cloud-barista/cb-ladybug#getting-started)
    - CB-Tumblebug (필수)
    - CB-Spider (필수)
    - cb-webtool (Web기반 GUI)
- [멀티 클라우드 쿠버네티스 서비스 사용 방법](https://github.com/cloud-barista/cb-ladybug/tree/master/docs/test#test)

***

## [API 및 문서]
- [API 규격](https://github.com/cloud-barista/docs/blob/master/technical_docs/API/CB-User_REST-API.md)
- [문서 통합 Repository](https://github.com/cloud-barista/docs)

***

## [특이 사항]
- 개발 단계: 기능 개발 우선 단계 (상용 활용시 안정화 및 보완 필요)
- CSP 연동 검증 상태
  - CB-Tumblebug 기준 테스트 완료된 CSP: AWS, GCP, Azure, Alibaba, Cloudit
  - CB-Ladybug 기준 테스트 완료된 CSP: AWS, GCP, Azure
  - 현재는 개발 단계이므로 기능 안정성은 낮을 수 있음 (버그 리포트 기여 환영합니다..!)

***
