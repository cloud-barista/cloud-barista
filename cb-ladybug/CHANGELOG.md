# Cloud-Barista Multi-Cloud Application Management Framework (CB-Ladybug) ChangeLog

## v0.4.0-CafeMocha (2021.06.30.)

### API Change
- Add some parametes in Cluster and Node Management APIs to support multi-cloud regions and high availability (#36)

### Feature
- Support Kubernetes(v1.18.9) Cluster Creation on Multi-Cloud Regions
- Support Control Plane High Availability (HAProxy v1.7) (#36)
- Support Multi-Cloud Network Overlay (Canal, Kilo) (#29)
- Support Microsoft Azure (#51)

### Bug Fix
- Update AWS's user account (#50)
- Use private ip for haproxy (#50)
- Fix connection error over 10 VMs (#50)
- Update MCIR naming rule/validation (#50)
- Remove return code from getHostName for adjusting mcis' vm id updates (#52)
- Add terminate, refine to delete cluster (#54)

### Note
- Default development environment: Go v1.16

***

## v0.3.0-espresso (2020.12.11.)

### API Change
- 클러스터 생성 기능 API 추가
- 클러스터 삭제 기능 API 추가
- 클러스터 정보 조회 기능 API 추가
- 노드 추가 기능 API 추가
- 노드 삭제 기능 API 추가
- 노드 정보 조회 기능 API 추가

### Feature
- 단일 클라우드 리전 대상 쿠버네티스 클러스터(v1.17.8) 생성 지원
- 아마존 AWS, 구글 GCP 지원
