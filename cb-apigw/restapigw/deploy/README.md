# Deployment

> Notes
> ---
> **<font color="red">이 테스트 구성은 모두 Docker Container를 기준으로 하고 있으므로 사전에 docker 와 docker-compose 가 설치되어 있어야 합니다.</font>**

## ETCD 관련 테스트 설정 (Mac에서 검증)

> API G/W 가 Docker Container가 아닌 경우에 ETCD (Docker compose)에 특정 IP로 접근하기 위한 로컬 테스트용

- Hosts 파일에 Loopback 처리용 IP 등록 (/etc/hosts)
  ```text
  172.16.238.11 etcd-0
  172.16.238.12 etcd-1
  172.16.238.13 etcd-2
  ```
- Network에 Loopback 별칭 적용
  ```shell
  $ sudo ifconfig lo0 alias 172.16.238.11
  $ sudo ifconfig lo0 alias 172.16.238.12
  $ sudo ifconfig lo0 alias 172.16.238.13
  ```
  > 단, 재 부팅되면 다시 설정해야 하며, lo0 는 시스템마다 다를 수 있으므로 확인 필요 함.
- 환경설정 파일 구성 (.env 파일)
  ```text
  NETWORK_NAME=${NETWORK_NAME:-etcd_net}
  NETWORK_CONFIG_SUBNET=172.16.238.0/24

  ETCD_00_NETWORKS_ETCD_NET_ADDRESS=172.16.238.11
  ETCD_01_NETWORKS_ETCD_NET_ADDRESS=172.16.238.12
  ETCD_02_NETWORKS_ETCD_NET_ADDRESS=172.16.238.13
  ENDPOINTS=http://${ETCD_00_NETWORKS_ETCD_NET_ADDRESS}:2379,http://${ETCD_01_NETWORKS_ETCD_NET_ADDRESS}:2379,http://${ETCD_02_NETWORKS_ETCD_NET_ADDRESS}:2379
  ```
- 실행방법
  ```shell
  $ docker-compose -f dc-etcd.yaml up
  ```

## 실행 방법

```shell
$ docker-compose up --build
```

상기의 명령으로 docker-compose 빌드 (--build) 진행 후에 실행 (up) 할 수 있습니다.
문제가 발생하면 바로 종료되므로 오류 메시지를 참고해서 문제를 해결하고 다시 실행하면 됩니다.

> Notes
> ---
> 변경된 내용이 없이 재 실행하는 경우는 Build 옵션을 사용하지 않아도 됩니다.
> ```shell
> $ docker-compose up
> ```
> 변경된 내용 (설정이나 소스 등)이 있는 경우는 반드시 Build 옵션을 사용해야 반영됩니다.

실행된 어플리케이션은 다음과 같습니다.

> Notes
> ---
> Background 서비스들을 실행하는 배포입니다.  
> Background 서비스들을 docker-compose로 실행한 후에 아래의 명령으로 API G/W를 별도 컨테이너로 구동해서 테스트를 진행합니다.  
> 
> ```shell
> $ cd ..
> # Docker Image Build
> $ docker build -t cb-restapigw .
> # Docker Container 실행
> $ docker run -itd --network deploy_default -p 8000:8000 cb-restapigw
> ```
>

- **<font color="red">InfluxDB : localhost:8086</font>**
  -  RESTAPIGW에서 Metrics 데이터를 저장하는 DB 서버입니다.
- **<font color="red">Grafana : localhost:3100</font>**
  - 수집된 Metrics 정보를 표시하는 UI 이므로 브라우저를 통해서 정보를 확인할 수 있습니다.
  - 초기 설정된 ID/PW 는 admin/admin 입니다.
  - 초기화면은 `Home Dashboard`입니다. 화면에 보이는 dashboard 리스트에서 `CB-RESTAPIGW`를 선택하시면 됩니다.
- **<font color="red">Jaeger : localhost:16686</font>**
  - RESTAPIGW에서 동작한 Trace 정보를 표시하는 수집기이며 UI를 제공하므로 브라우저를 통해서 정보를 확인할 수 있습니다.
  - 왼쪽의 `Search` 탭의 `Service` 부분에 cb-restapigw를 선택하고 `Find Traces` 버튼을 누르면 Trace 정보를 확인할 수 있습니다.
  - 단, Trace 수집 주기가 있으므로 초기에는 서비스가 보이지 않을 수 있습니다. refresh를 해서 서비스가 등록되었는지를 확인이 필요하며, 수집 주기 (10s) 이후에도 서비스가 등록되지 않았다면 터미널의 로그를 통해서 문제가 있는지를 확인해야 합니다.
- **<font color="red">Fake API : localhost:8100</font>**
  - 테스트를 위한 샘플 API
- **<font color="red">HMAC Server : localhost:8010</font>**
  - 테스트를 위한 HMAC 기반 인증 발급서버

## 테스트 방법

- HMAC 테스트 방법
  - ./web/conf 폴더의 hmac.yaml 에 설정 값과 ./conf/cb-restapigw.yaml 설정의 mw-auth 부분의 secure_key 부분을 동일하게 설정
  - 사용자 인증 접근을 제한할 경우는 ./conf/cb-restapigw.yaml 설정의 mw-auth 부분의 access_ids 에 허용할 리스트 설정
  - hmac_site의 secure_key 와 access_key 가 다르거나 duration 지정을 초과한 시간은 모두 401 access denied 가 발생함. (Response Header의 message로 처리된 메시지를 확인)
- 브라우저 또는 POSTMAN을 사용해서 ./conf/cb-restapigw.yaml 설정에 맞는 Endpoint 호출
  - http://localhost:8000/splash
    - 2개 Backend API를 호출해서 결과가 Merge되는 것 확인
    - flatmap filter를 통한 결과 필드명 변경 ("id" -> "id-") 확인
    - flatmap filter를 통한 불 필요 결과값 삭제 확인
  - http://localhost:8000/sequential
    - 2개 Backend API 를 순차적으로 처리 확인
    - whitelist filter를 통한 결과값 추출 확인
    - 먼저 처리된 Backend API 의 결과 값을 다음 실행될 Backend API의 변수로 활용 확인
    - 오류 발생시 상세 메시지 출력 확인
  - http://localhost:8000/fail
    - 2개 Backend 호출 중 1개에서 오류 발생한 경우 일부 데이터만 반환 확인
    - Response Header 에 처리 완료 (X-Cb-Restapigw-Complete) 및 메시지 (X-Cb-Restapigw-Message) 확인
  - http://localhost:8000/public
    - 외부 site를 Backend로 사용 확인
    - whitelist filter를 통한 결과 추출 확인
    - mapping filter를 통한 결과 필드명 이름 변경 확인
    - group filter를 통한 결과 그룹 처리 확인
  - http://localhost:8000/github/[사용자id]
    - 외부 site를 Backend로 사용 확인
    - Path variable 사용 확인
  - http://localhost:8000/collection
    - 결과가 JSON 객체가 아닌 Collection인 경우에 core.CollectionTag ("collection") 이라는 필드로 JSON 구성 반환 확인
  - http://localhost:8000/private/custom
    - HMAC 기반 인증 동작 확인
    - HMAC 기능 검증 (Hash 인증, Access IDs, Duration)에 따른 401 발생 확인

## 실행 중지

실행 상태인 터미널에서 `Ctrl+C` 로 중지 시그널을 처리하면 종료됩니다. 아래의 명령으로 사용된 리소스를 해제해 주시면 됩니다.
```shell
$ docker-compose down
```

> Notes
> ---
> 만일 터미널을 종료한 상태라면 다음과 같이 docker-compose.yaml 파일일 존재하는 폴더에서 터미널을 열고 아래의 명령을 실행하시면 됩니다.
> ```shell
> $ docker-compose stop   # docker-compose 종료
> $ docker-compose down   # 사용한 리소스 해제
> ```