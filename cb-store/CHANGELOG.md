# Changelog

## v0.6.0 (Cafe Latte, 2022.07.08.)

- Add 'Merge' function for NutsDB mode

## v0.5.0 (Affogato, 2021.12.16.)

## v0.4.0 (Cafe Mocha, 2021.06.30.)

### Bug Fix

- bug fix, error msg: not found bucket:bucketForString,key:XXX (211fb8e)
- fix limited fetch bug in GetList(). Limitation was 10,000 of temp code (67abef9)

## v0.3.0-cappuccino (2020.12.10.)

### Feature

- `go.mod` 파일 추가
- ETCD client 가 제대로 구성되지 않을 경우에 Store 처리 메서드 호출되는 부분 검증 추가

### Bug Fix

- Driver Initialize 시점을 GetStore 시점으로 변경
- NutsDB 버전 관련한 빌드 오류 부분 해결을 위한 go.mod 수정 (grpc 버전을 v1.26.0 으로 고정)
- ETCD 환경이 없는 상태에서 ETCD Client 가 hang 걸리는 문제 해결을 위해 go.mod 수정 (Dial timeout 무시되는 문제)

## v0.2.0-cappuccino (2020.06.01.)

### Bug Fix

- nutsdb 업데이트에 대한 대응 [#11](https://github.com/cloud-barista/cb-store/issues/11)

## v0.1.0-americano (2019.12.23.)
