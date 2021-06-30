
# v0.4.0 (Cafe Mocha) (2021.06.30.)

### Feature

- Remove CB-Dragonfly-related Docker network ([#110](https://github.com/cloud-barista/cb-operator/pull/110))
- Modify PV access mode to ReadWriteOnce ([#127](https://github.com/cloud-barista/cb-operator/pull/127))



# v0.3.0-espresso (2020.12.10.)

## ChangeLog

### Feature

- Add mode-selection feature (#39)
- Add Helm chart (#40)
- Add PVC for Cloud-Barista components (#44)
- Add CB-Ladybug to docker-compose.yaml and Helm chart (#71)
- Add Prometheus and Grafana to CB Helm chart (#81)
- Change docker network name (#87)



# v0.2.0-cappuccino (2020.06.02.)

### Changelog

- cb-operator 공개 (Docker Compose 기반)

### Features

- pull: CB 컨테이너 이미지들을 로컬 이미지 저장소로 다운로드
- run: CB 컨테이너들을 실행하여 CB 시스템을 구동
- info: CB 컨테이너들의 상태와 이미지 현황을 표시
- exec: CB 개별 컨테이너에 접속하여 명령을 실행
- stop: CB 컨테이너들을 중지하여 CB 시스템을 중지
- remove: CB 컨테이너 (+ 볼륨, 이미지) 를 제거
