# v0.4.0-CafeMocha (2021.06.30.)

### API Change 
Ref) [API ChangeLog](https://github.com/cloud-barista/cb-tumblebug/discussions/416)

- Add VMGroup parameter in create MCIS API
- Add Private IP parameter in get MCIS status API
- Add MCIS Refine option in MCIS action (get) API
- Add verifiedUserName parameter in get spec API
- Add API for ListResourceId, ListMcisId, ListVmId 
- Add TB object control API
- Add inspectResources API
- Change API style: snakeCase to camelCase


### Feature
Ref) [Supported cloud service providers](https://github.com/cloud-barista/cb-tumblebug/discussions/429)

- Add VM group feature to request multiple VMs simply [#413](https://github.com/cloud-barista/cb-tumblebug/pull/413)
- Provide SystemMessage to vm status object [#475](https://github.com/cloud-barista/cb-tumblebug/pull/475)
- Enhance and expedite mcis lifecycle handling [#625](https://github.com/cloud-barista/cb-tumblebug/pull/625)
- Add MCIS Refine feature [#572](https://github.com/cloud-barista/cb-tumblebug/pull/572)
- Add feature for general TB object retrieve [#417](https://github.com/cloud-barista/cb-tumblebug/pull/417)
- Add initial code for mcis and vm plan with location-based algo [#511](https://github.com/cloud-barista/cb-tumblebug/pull/511)
- Add inspectVMs function [#505](https://github.com/cloud-barista/cb-tumblebug/pull/505)
- Expedite auto agent installation [#448](https://github.com/cloud-barista/cb-tumblebug/pull/448)
- Enhance ssh username verification performance [#423](https://github.com/cloud-barista/cb-tumblebug/pull/423) 
- Add WeaveScope deployment script [#419](https://github.com/cloud-barista/cb-tumblebug/pull/419)
- Add jitsi video conference automation [#476](https://github.com/cloud-barista/cb-tumblebug/pull/476)
- Add script for deploying web game server [#609](https://github.com/cloud-barista/cb-tumblebug/pull/609)

### Bug Fix
- Enhance error handing for provisioning and cmd phases [#435](https://github.com/cloud-barista/cb-tumblebug/pull/435)
- Fix agent installation bug and script update [#437](https://github.com/cloud-barista/cb-tumblebug/pull/437)
- Fix initial failed status in MCIS provisioning [#467](https://github.com/cloud-barista/cb-tumblebug/pull/467)
- Fix list object key parsing bug [#607](https://github.com/cloud-barista/cb-tumblebug/pull/607)
- Patch gRPC API [#536](https://github.com/cloud-barista/cb-tumblebug/pull/536)

### Note
- Default development environment: Go v1.16 

***

# v0.3.0-espresso (2020.12.03.)

### API Change
- MCIS 자동 제어 기능 API 추가
- 동적 시스템 환경 설정 변경 기능 API 추가
- MCIS 생성 API의 모니터링 에이전트 자동 배치 옵션 제공

### Feature
- MCIS 생성시 모니터링 에이전트 자동 배치 기능 추가
- MCIS 자동 제어 기능 추가
- MCIS 시나리오 테스트 스크립트 추가
- MCIS 마스터 VM 및 VM IP 정보 제공 기능 추가
- MCIR VM 사양 패치 및 등록 기능 추가
- 동적 시스템 환경 설정 변경 기능 추가

### Bug Fix
- MCIS 종료시 런타임 오류 수정

***

# v0.2.0-cappuccino (2020.06.02.)

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

***

# v0.1.0-americano (2019.12.23.)

### Feature
- Namespace, MCIR, MCIS 관리 기본 기능 제공
