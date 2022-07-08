# v0.6.0 (Cafe Latte, 2022.07.08.)

## What's Changed

* Bump actions/checkout from 2 to 3 by @dependabot in <https://github.com/cloud-barista/cb-operator/pull/179>
* Bump peter-evans/create-pull-request from 3 to 4 by @dependabot in <https://github.com/cloud-barista/cb-operator/pull/180>
* Update CB-Dragonfly HelmChart config file by @hyokyungk in <https://github.com/cloud-barista/cb-operator/pull/177>
* Update CB container image tags by @jihoon-seo in <https://github.com/cloud-barista/cb-operator/pull/181>
* Bump actions/setup-go from 2 to 3 by @dependabot in <https://github.com/cloud-barista/cb-operator/pull/182>
* Bump github/codeql-action from 1 to 2 by @dependabot in <https://github.com/cloud-barista/cb-operator/pull/183>

**Full Changelog**: <https://github.com/cloud-barista/cb-operator/compare/v0.5.0...v0.6.0>

# v0.5.0 (Affogato, 2021.12.16.)

## What's Changed

* [Workflow] Update CHANGELOG-generated.md by @jihoon-seo in [#132](https://github.com/cloud-barista/cb-operator/pull/132)
* Add support for GKE by @jihoon-seo in [#134](https://github.com/cloud-barista/cb-operator/pull/134)
* Remove `sudo` from Docker Compose & K8s mode by @jihoon-seo in [#141](https://github.com/cloud-barista/cb-operator/pull/141)
* Add etcd Helm chart, and Update Helm chart, `docker-compose.yaml` & `cb-operator` by @jihoon-seo in [#142](https://github.com/cloud-barista/cb-operator/pull/142)
* CB-Dragonfly FW config 파일 최신화 by @inno-cloudbarista in [#140](https://github.com/cloud-barista/cb-operator/pull/140)
* Update Helm chart and docker-compose.yaml by @jihoon-seo in [#145](https://github.com/cloud-barista/cb-operator/pull/145)
* Add 'operator update' subcommand by @jihoon-seo in [#150](https://github.com/cloud-barista/cb-operator/pull/150)
* Add initContainers to wait for etcd by @jihoon-seo in [#147](https://github.com/cloud-barista/cb-operator/pull/147)
* Update cb-webtool env vars by @jihoon-seo in [#151](https://github.com/cloud-barista/cb-operator/pull/151)
* Upload doc images by @jihoon-seo in [#153](https://github.com/cloud-barista/cb-operator/pull/153)
* Update cb-operator by @jihoon-seo in [#156](https://github.com/cloud-barista/cb-operator/pull/156)
* change some variables to constant(#152) by @computerphilosopher in [#163](https://github.com/cloud-barista/cb-operator/pull/163)
* Add build-test GitHub workflow [ci skip test] by @jihoon-seo in [#158](https://github.com/cloud-barista/cb-operator/pull/158)
* Reflect renaming: CB-Ladybug → CB-MCKS by @jihoon-seo in [#162](https://github.com/cloud-barista/cb-operator/pull/162)
* Update 'paths-to-ignore' in build test workflow by @jihoon-seo in [#165](https://github.com/cloud-barista/cb-operator/pull/165)
* Fix build test workflow by @jihoon-seo in [#167](https://github.com/cloud-barista/cb-operator/pull/167)
* Update CB-Dragonfly config files by @jihoon-seo in [#164](https://github.com/cloud-barista/cb-operator/pull/164)
* Add Chronograf to docker-compose-df-only by @seokho-son in [#168](https://github.com/cloud-barista/cb-operator/pull/168)
* Update Kong conf in docker-compose-mode by @jihoon-seo in [#170](https://github.com/cloud-barista/cb-operator/pull/170)
* Change docker compose df only image version by @seokho-son in [#173](https://github.com/cloud-barista/cb-operator/pull/173)
* Update CB-Dragonfly Docker-compose & HelmChart config file by @hyokyungk in [#172](https://github.com/cloud-barista/cb-operator/pull/172)
* Fix error in df config for helm chart by @seokho-son in [#174](https://github.com/cloud-barista/cb-operator/pull/174)
* Update config files by @jihoon-seo in [#175](https://github.com/cloud-barista/cb-operator/pull/175)

## New Contributors

* @computerphilosopher made their first contribution in [#163](https://github.com/cloud-barista/cb-operator/pull/163)
* @hyokyungk made their first contribution in [#172](https://github.com/cloud-barista/cb-operator/pull/172)

**Full Changelog**: <https://github.com/cloud-barista/cb-operator/compare/v0.4.0...v0.5.0>

# v0.4.0 (Cafe Mocha, 2021.06.30.)

### Feature

* Remove CB-Dragonfly-related Docker network ([#110](https://github.com/cloud-barista/cb-operator/pull/110))
* Modify PV access mode to ReadWriteOnce ([#127](https://github.com/cloud-barista/cb-operator/pull/127))

# v0.3.0-espresso (2020.12.10.)

## ChangeLog

### Feature

* Add mode-selection feature (#39)
* Add Helm chart (#40)
* Add PVC for Cloud-Barista components (#44)
* Add CB-Ladybug to docker-compose.yaml and Helm chart (#71)
* Add Prometheus and Grafana to CB Helm chart (#81)
* Change docker network name (#87)

# v0.2.0-cappuccino (2020.06.02.)

### Changelog

* cb-operator 공개 (Docker Compose 기반)

### Features

* pull: CB 컨테이너 이미지들을 로컬 이미지 저장소로 다운로드
* run: CB 컨테이너들을 실행하여 CB 시스템을 구동
* info: CB 컨테이너들의 상태와 이미지 현황을 표시
* exec: CB 개별 컨테이너에 접속하여 명령을 실행
* stop: CB 컨테이너들을 중지하여 CB 시스템을 중지
* remove: CB 컨테이너 (+ 볼륨, 이미지) 를 제거
