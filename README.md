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

</details>

---
*This repository is an integrated archive for repository of major frameworks.* These repositories are included and listed in the root directory. This repo reflects the latest release only.

The main frameworks or tools are as follows (the release version of each repository may vary),

- **CB-Spider** (connects all clouds in a single interface)
  - Upstream repo: <https://github.com/cloud-barista/cb-spider>
- **CB-Tumblebug** (manages multi-cloud infrastructures)
  - Upstream repo: <https://github.com/cloud-barista/cb-tumblebug>
- **CB-Dragonfly** (monitors multi-cloud services)
  - Upstream repo: <https://github.com/cloud-barista/cb-dragonfly>
  - Note: not updated since v0.8.0, possible to be deprecated
- CB-Bridge/**cb-operator** (operation tool for Cloud-Barista system runtime)
  - Upstream repo: <https://github.com/cloud-barista/cb-operator>
  - Note: As the components of Cloud-Barista are currently simplified, using cb-operator might require additional effort for users. It is recommended to use cb-tumblebug directly.
- CB-Bridge/**cb-store** (provides an unified DB interface for meta info of Cloud-Barista)
  - Upstream repo: <https://github.com/cloud-barista/cb-store>
- CB-Bridge/**cb-log** (provides log library to Cloud-Barista system)
  - Upstream repo: <https://github.com/cloud-barista/cb-log>

**[Note]** CB-Larva is a special repository that incubates (research and develop) new Multi-Cloud technologies.
CB-Larva explores interesting ideas and shows the possibility of those (i.e., Proof of Concept (POC)).
That's why we encourage you to take a look and contribute to the special repository.
Please note that the source code of CB-Larva would not be released and archived in this repository for the time being.

- CB-Larva/cb-cladnet (POC for the cloud adaptive network)
  - Upstream repo: <https://github.com/cloud-barista/cb-larva/tree/main/poc-cb-net>

***

## [목    차]

1. [실행 환경](#실행-환경)
1. [설치 및 실행](#설치-및-실행)
1. [사용 방법 및 예시](#사용-방법-및-예시)
1. [API 및 문서](#api-및-문서)
1. [특이 사항](#특이-사항)

***

## [실행 환경]

- Linux (추천: Ubuntu v22.04)

***

## [설치 및 실행]

- Cloud-Barista 플랫폼 통합 실행 (Docker 이미지 기반)
  - cb-operator 를 통해 Cloud-Barista 전체 FW를 통합 실행할 수 있음
    - 참고: [cloud-barista/cb-operator](/cb-operator/)
    - As the components of Cloud-Barista are currently simplified, using cb-operator might require additional effort for users. It is recommended to use cb-tumblebug directly.

- Cloud-Barista 플랫폼 개별 FW 소스 다운로드 및 설치

  - CB-Spider 설치 및 실행
    - [cloud-barista/cb-spider README를 참고하여 설정, 설치](/cb-spider/)
    - CB-Spider 실행

  - CB-Tumblebug 설치 및 실행
    - [cloud-barista/cb-tumblebug README를 참고하여 설정, 설치](/cb-tumblebug/)
      - CB-Spider API 서버 주소를 conf/setup.env 에 설정
      - CB-Dragonfly API 서버 주소를 conf/setup.env 에 설정
    - CB-Tumblebug 실행

  - CB-Dragonfly 설치 및 실행
    - [cloud-barista/cb-dragonfly README를 참고하여 설정, 설치](/cb-dragonfly/)
    - CB-Dragonfly 실행


***

## [사용 방법 및 예시]

### 주요 서비스: 멀티 클라우드 인프라 서비스 (MCIS)

- 멀티 클라우드 인프라 서비스 환경 구성
  - [CB-Tumblebug 설정 및 실행](https://github.com/cloud-barista/cb-tumblebug#cb-tumblebug-%EC%86%8C%EC%8A%A4-%EB%B9%8C%EB%93%9C-%EB%B0%8F-%EC%8B%A4%ED%96%89-%EB%B0%A9%EB%B2%95-%EC%83%81%EC%84%B8)
    - CB-Spider (필수)
    - CB-Dragonfly (MCIS 모니터링, CB-Tumblebug 자동 제어 기능에 필요)
- [멀티 클라우드 인프라 서비스 사용 방법](https://github.com/cloud-barista/cb-tumblebug#cb-tumblebug-%EA%B8%B0%EB%8A%A5-%EC%82%AC%EC%9A%A9-%EB%B0%A9%EB%B2%95)
- [멀티 클라우드 인프라 유스케이스](https://github.com/cloud-barista/cb-tumblebug/blob/main/README.md#3-%EB%A9%80%ED%8B%B0-%ED%81%B4%EB%9D%BC%EC%9A%B0%EB%93%9C-%EC%9D%B8%ED%94%84%EB%9D%BC-%EC%9C%A0%EC%8A%A4%EC%BC%80%EC%9D%B4%EC%8A%A4)

***

## [API 및 문서]

- [API 규격](https://github.com/cloud-barista/docs/blob/master/technical_docs/cloud-barista/API/CB-User_REST-API.md)
- [문서 통합 Repository](https://github.com/cloud-barista/docs)

***

## [특이 사항]

- 개발 단계: 기능 개발 우선 단계 (상용 활용시 안정화 및 보완 필요)
- CSP 연동 검증 상태
  - CB-Spider 기준 테스트 완료된 CSP: [링크 1](https://github.com/cloud-barista/cb-spider#3-제공-자원) 및 [링크 2](https://github.com/cloud-barista/cb-spider/wiki/Supported-CloudOS) 참고
  - CB-Tumblebug 기준 테스트 완료된 CSP: [링크](https://github.com/cloud-barista/cb-tumblebug/wiki/Supported-CSPs) 참고
  - 현재는 개발 단계이므로 기능 안정성은 낮을 수 있음 (버그 리포트 기여 환영합니다!)

***
