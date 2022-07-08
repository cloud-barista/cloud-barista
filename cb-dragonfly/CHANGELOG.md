
# v0.6.0 (CafeLatte, 2022.07.08.)
## Works with
- CB-Tumblebug (https://github.com/cloud-barista/cb-tumblebug/releases/tag/v0.6.0)

### API
 - Postman API URL:  https://documenter.getpostman.com/view/11342380/UzJFvJUi

## What's Changed
* Add MCK8S monitoring features by @inno-cloudbarista in https://github.com/cloud-barista/cb-dragonfly/pull/134
* Add new MCK8s monitoring features in monitoring agents
* Add new MCK8s monitoring API

**Full Changelog**: https://github.com/cloud-barista/cb-dragonfly/compare/v0.5.1...v0.6.0

---

# v0.5.0 (Affogato, 2021.12.16.)

### Tested with
 - CB-Tumblebug (https://github.com/cloud-barista/cb-spider/releases/tag/v0.4.15)

### API
 - Postman API URL:  https://documenter.getpostman.com/view/10735617/UVCCfPcj

### Note
* Modify install agent parameter validation check logic by @hyokyungk in [#125](https://github.com/cloud-barista/cb-dragonfly/pull/125)
* Collector delete bugfix by @hyokyungk in [#127](https://github.com/cloud-barista/cb-dragonfly/pull/127)
* Upgrade monitoring logic for efficient operation(Push Only)
* Add agent health check function(Push)
* Support collector auto-scaling on k8s cluster envrionment
* Add New MCIS metrics(network packets, VM process usage)
* 버그 픽스 by @inno-cloudbarista in [#116](https://github.com/cloud-barista/cb-dragonfly/pull/116)
* DF 배포 환경 설정 변수 추가(Config) & k8s 배포를 위한 Helm 추가 by @inno-cloudbarista in [#118](https://github.com/cloud-barista/cb-dragonfly/pull/118)
* Update Config Load Logic by @inno-cloudbarista in [#121](https://github.com/cloud-barista/cb-dragonfly/pull/121)
* Update kafka log retention policy by @hyokyungk in [#119](https://github.com/cloud-barista/cb-dragonfly/pull/119)
* swagger docs 적용 by @inno-cloudbarista in [#108](https://github.com/cloud-barista/cb-dragonfly/pull/108)
* [bugfix] query 조회 버그 수정 및 메트릭 누락 항목 추가 by @inno-cloudbarista in [#109](https://github.com/cloud-barista/cb-dragonfly/pull/109)
* [gRPC] gRPC 서버, 클라이언트 기능 추가 by @inno-cloudbarista in [#110](https://github.com/cloud-barista/cb-dragonfly/pull/110)
* [공통] API FORM -> JSON 변경 by @inno-cloudbarista in [#113](https://github.com/cloud-barista/cb-dragonfly/pull/113)
* [공통] cb-store 키 규칙 통일 by @inno-cloudbarista in [#114](https://github.com/cloud-barista/cb-dragonfly/pull/114)
* [기능 개선] 모니터링 메트릭 저장/조회 이름 통일 및 에이전트 미활용 설정 정보 삭제 by @pjini in [#99](https://github.com/cloud-barista/cb-dragonfly/pull/99)
* [공통] docker-compose 환경 InfluxDB 포트 설정 관련 버그 수정 by @inno-cloudbarista in [#105](https://github.com/cloud-barista/cb-dragonfly/pull/105)
* [공통] config 파일 주석 및 기본 동작 방식 수정 (기본 push 방식, avg 함수 기반 aggregate) by @inno-cloudbarista in [#106](https://github.com/cloud-barista/cb-dragonfly/pull/106)

### What's Changed

**Full Changelog**: https://github.com/cloud-barista/cb-dragonfly/compare/v0.4.0...v0.5.0

---

# v0.4.0 (Cafe Mocha, 2021.06.30.)

### API Change 
- 에이전트 메타데이터 조회 API 추가

### Feature
- PULL 방식 모니터링 모니터링 개발
- PUSH/PULL 메커니즘 기반 CB-Dragonfly FW 구동 모듈 고도화
- PUSH/PULL 메커니즘 기반 에이전트 구동 모듈 고도화

### Bug Fix
- Kafka 토픽 관련 오류 개선
- 모니터링 데이터 조회 로직 개선

# v0.3.0-cappuccino (2020.12.03.)

### API Change
- 온디멘드 API 추가
- MCIS 모니터링 API 추가
- 알람 이벤트 핸들러, 태스크, 알람 API 추가
- CB-Dragonfly FW 헬스체크 API 추가

### Feature
- Kafka 기반 부하분산 모듈 고도화
- MCIS 성능 모니터링 메트릭 제공
- 온디멘드 기반 모니터링 제공
- CB-Store 기반 모니터링 데이터 저장
- 알람 기능 개발
- Go API(gRPC 통신) 및 CLI 도구 지원
- Cloud-Twin 환경 에이전트 구동

### Bug Fix
- 모니터링 조회 API 오류 개선

# v0.2.0-cappuccino (2020.06.02.)

### API Change
- CB-Dragonfly agent_TTL 환경설정 변수 추가
- 모니터링 에이전트 설치 API 내부 로직 개선

### Feature
- 리눅스 환경 모니터링 추가 메트릭 지원
- 윈도우즈 환경 모니터링 지원
- Load-Balancer 모듈 기반 대규모 모니터링 콜렉터 안정성 개선
- CB-Dragonfly 최초 구동 시 자동 Configuration 설정 기능 추가

### Bug Fix
- CB-Dragonfly 에이전트 설치 로직 및 오류 개선
- 대규모 모니터링 안정성 테스트 기반 모니터링 콜렉터 로직 개선

# v0.1.0-americano (2019.12.23.)

### Feature
- 리눅스 환경 모니터링 지원
- CB-Dragonfly API 기반 에이전트 설치 및 모니터링 조회
- 1분 미만의 최신 모니터링 조회
- 에이전트 수 기반 콜렉터 오토스케일링
