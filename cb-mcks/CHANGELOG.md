# Cloud-Barista Multi-Cloud Kubernetes Service Framework (CB-MCKS) ChangeLog

## v0.5.0 (Affogato, 2021.12.16.)

### Tested with

- CB-Spider (https://github.com/cloud-barista/cb-spider/releases/tag/v0.5.0)
- CB-Tumblebug (https://github.com/cloud-barista/cb-tumblebug/releases/tag/v0.5.0)
- CB-Dragonfly (https://github.com/cloud-barista/cb-dragonfly/releases/tag/v0.5.0)


### API Change
- Add installMonAgent parameter for monitoring (#115)
- Add entity for node (cspLabel, regionLabel, zoneLabel) (#100, #108)
- Add parameter and entity (label, description, createdTime) for cluster, node (#99, #107)
- Improve provisioning status information (#105)
- Change default port number (8080 -> 1470) (#92, #96)
- Rename CB-Ladybug to CB-MCKS (#96)

### Feature
- Add mcis.systemLabel for monitoring (#114)
- Add support for openstack (#104)
- Add deleteMCIS/deleteVMs for createCluster/addNode failure (#98)
- Add support for Tencent-cloud (#97)
- Add labels to kubernetes node (csp, region, zone) (#95)
- Use cb-tumblebug api instead of cb-spider api (#39, #94)
- Improve installing addons(network cni) (#64, #93)
- Add mcis label (#90)
- Add support for aws osaka region (#89)
- Improve swagger docs (#42, #88)
- Add support for Alibaba cloud (#84)

### Bug Fix
- Update default network cni (to canal) (#113)
- Update ns-delete/ns-create test scripts (#109, #112)
- Add drain/delete node when addNode failed (#111)
- Sync a node name on MCKS and kubernetes (#103, #106)
- Remove cluster name from node name (#98)
- Change mcis control url (#98)
- Remove uid from cluster & node struct ($49, #87)
- Improve ssh connection (#75, #84)

### Note
- Full Changelog: https://github.com/cloud-barista/cb-mcks/compare/v0.4.0...v0.5.0

***

## v0.4.0 (Cafe Mocha, 2021.06.30.)

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
