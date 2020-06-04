
# v0.2.0-cappuccino (2020.06.04.)


## [cb-spider]

### API Change
- CloudO 목록 제공 API 추가
- VPC/Subnet API 추가 ([#9](https://github.com/cloud-barista/cb-spider/pull/9) [#226](https://github.com/cloud-barista/cb-spider/pull/226))
- VMSpec API 추가 ([#151](https://github.com/cloud-barista/cb-spider/pull/151) [#223](https://github.com/cloud-barista/cb-spider/pull/223))
- VNic API 삭제
- PublicIP API 삭제

### Feature
- 통합ID IID Manager 추가 ([#163](https://github.com/cloud-barista/cb-spider/pull/163) [#194](https://github.com/cloud-barista/cb-spider/pull/194))  
- VPC/Subnet 기능 추가  ([#9](https://github.com/cloud-barista/cb-spider/pull/9) [#226](https://github.com/cloud-barista/cb-spider/pull/226)) 
- VNic, PublicIP 자동 관리 기능으로 개선
- Cloud Driver 및 Region 정보 자동 등록 지원 도구 추가 utils/import-info/*
- Docker Driver 추가(Hetero Multi-IaaS 제어)
- Android 운영 환경을 위한 plugin off mode 추가 ([3938ea0](https://github.com/cloud-barista/cb-spider/commit/3938ea0c70e69664a62eb3cee6611cfbf26ea4ea))  

### Bug Fix


## [cb-tumblebug]

### API Change
- MCIS 통합 원격 커맨드 기능 API 추가
- 개별 VM 원격 커맨드 기능 API 추가
- MCIR Subnet 관리 API 제거
- MCIR VNic 관리 API 제거
- MCIR PublicIP 관리 API 제거
- 전체 Request 및 Response Body의 상세 항목 변경 (API 예시 참고)

### Feature
- MCIS 및 VM에 현재 수행 중인 제어 명령 정보를 관리
- 멀티 클라우드 동적 성능 밴치마킹 기능 일부 추가 (PoC 수준)
- MCIS VM 생성 및 제어시 Goroutine을 적용하여 속도 개선
- MCIS 및 VM 원격 커맨드 기능 추가
- MCIS 오브젝트 정보 보완 (VM의 위경도 정보 제공)

### Bug Fix
- MCIS 라이프사이클 오류 개선


## [cb-dragonfly]

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


## [cb-webtool]

### API Change
- Geolocation API 추가
- 전체 Request 및 Response Body의 상세 항목 변경
- 각 MCIS에서 PublicIP 추출 기능 API 추가
- Common URL API 추가

### Feature
- Location에 선택된 서비스의 위치 반영
- VM 모니터링 활성화 & 모니터링 차트 추가
- Dashboard 변경및 메인 내용 영문화 적용
- 환경 & 리소스 설정 기능 변경 및 보완
- 환경 변수에 로그인 계정 설정 추가
- cb-tumblebug & cb-spider & cb-dragonfly 변경된 API 반영

### Bug Fix
- 환경 & 리소스 설정 버그 수정


## [cb-operator]

### Changelog
- cb-operator 공개 (Docker Compose 기반)

### Features
- pull: CB 컨테이너 이미지들을 로컬 이미지 저장소로 다운로드
- run: CB 컨테이너들을 실행하여 CB 시스템을 구동
- info: CB 컨테이너들의 상태와 이미지 현황을 표시
- exec: CB 개별 컨테이너에 접속하여 명령을 실행
- stop: CB 컨테이너들을 중지하여 CB 시스템을 중지
- remove: CB 컨테이너 (+ 볼륨, 이미지) 를 제거


# v0.1.0-americano (2019.12.23.)

### Feature
- Cloud-Barista Initial Features
