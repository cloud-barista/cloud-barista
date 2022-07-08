# Cloud-Barista Multi-Cloud Kubernetes Service Framework (CB-MCKS) ChangeLog

## v0.6.0 (CafeLatte, 2022.07.08.)

### Tested with

- CB-Spider (https://github.com/cloud-barista/cb-spider/releases/tag/v0.6.0)
- CB-Tumblebug (https://github.com/cloud-barista/cb-tumblebug/releases/tag/v0.6.0)
- CB-Dragonfly (https://github.com/cloud-barista/cb-dragonfly/releases/tag/v0.6.0)


### API Change

- Add support for k8s 1.23 version [#136](https://github.com/cloud-barista/cb-mcks/pull/136)
- Improve a cbadm cli [#135](https://github.com/cloud-barista/cb-mcks/pull/135)

### Feature

- Add support for IBM-cloud [#119](https://github.com/cloud-barista/cb-mcks/pull/119)
- Implement 'get VM SpecList' rest API  feature [#128](https://github.com/cloud-barista/cb-mcks/pull/128)
- Implement 'get VM SpecList' gRPC API  feature [#129](https://github.com/cloud-barista/cb-mcks/pull/129)
- Add support for Cloudit [#131](https://github.com/cloud-barista/cb-mcks/pull/131)
- Add support for k8s 1.23 version [#136](https://github.com/cloud-barista/cb-mcks/pull/136)
- Add a inline parameters to cbadm commands(create a cluster & add nodes) [#137](https://github.com/cloud-barista/cb-mcks/pull/137)

### Bug Fix

- Modify the worker-join process hang when creating a cluster by applying Kilo-CNI [#121](https://github.com/cloud-barista/cb-mcks/pull/121)
- Duplicated private-ip problem on multiple VPCs [#124](https://github.com/cloud-barista/cb-mcks/pull/124)
- Delete body for http get Method [#126](https://github.com/cloud-barista/cb-mcks/pull/126)
- Bump github.com/beego/beego/v2 from 2.0.1 to 2.0.2 [#132](https://github.com/cloud-barista/cb-mcks/pull/132)
- Insert swapoff into k8s 1.23 install script [#138](https://github.com/cloud-barista/cb-mcks/pull/138)

### Refactoring

- Improve a source structure [#122](https://github.com/cloud-barista/cb-mcks/pull/122)

### Documentation

- Tidy markdown documents [#118](https://github.com/cloud-barista/cb-mcks/pull/118)
- Fix connectioninfo-create.sh 'CSP' add [#130](https://github.com/cloud-barista/cb-mcks/pull/130)


### Note
- Full Changelog: https://github.com/cloud-barista/cb-mcks/compare/v0.5.0...v0.6.0

***

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
