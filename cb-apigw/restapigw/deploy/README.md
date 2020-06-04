# Deployment

> Notes
> ---
> **<font color="red">이 테스트 구성은 모두 Docker Container를 기준으로 하고 있으므로 사전에 docker 와 docker-compose 가 설치되어 있어야 합니다.</font>**

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

## REAT API G/W 환경 설정

> Notes 
> ---
> - **Configuration 설정은 `YAML` 포맷을 사용한다.**
> - **각 Host Address는 Docker Compose Network 상에서 처리되므로 container name과 container port 기준으로 작성해야 합니다.**</br>
>   ex) Fake API : http://localhost:8100 이 아닌 http://fake_api:8080

### 주요 설정은 다음과 같이 구성된다.
  - Service
    - Service 식별 설정

      | 설정 | 내용 | 필수 | 기본값 |
      |---|---|:-:|---|
      | name | 서비스 식별 명 | 필수 |  |
      | version | 설정 파일 버전 | 필수 | 1 |
      | host | 서비스에서 모든 Backend host로 공통 적용할 host 리스트 </br>(개별 Backend에 host 미 지정시 적용 됨) |  |  |
      | port | 서비스 운영 포트 </br>설정 값보다 실행 옵션 (-p) 값이 우선 적용 됨|  | 8000 |
      | timeout | 서비스 처리 제한 시간 |  | 2s |

    - TLS 설정

      | 설정 | 내용 | 필수 | 기본값 |
      |---|---|:-:|---|
      | public_key | TLS에 적용할 공개키 파일  | 필수 |  |
      | private_key | TLS에 적용할 비밀키 파일  | 필수 |  |

    - Endpoint List
        - Endpoint 설정

          | 설정 | 내용 | 필수 | 기본값 |
          |---|---|:-:|---|
          | endpoint | 클라이언트에 노출할 URL | 필수 |  |
          | method | REST 요청 메서드 (GET/PUT/POST/DELETE/...) |  | GET |
          | timeout | 엔드포인트 처리 제한 시간 </br>지정하지 않으면 서비스에 지정한 timeout 사용 |  | 2s |
          | except_querystrings | 클라이언트 요청에서 백엔드 요청으로 전달할 때 제외할 쿼리스트링 리스트 (기본은 전체 전달)|  |  |
          | except_headers | 클라이언트 요청에서 백엔드 요청으로 전달할 때 제외할 헤더 명 리스트 (기본은 전체 전달) |  |  |

        - Backend List
          - Backend 설정
          
            | 설정 | 내용 | 필수 | 기본값 |
            |---|---|:-:|---|
            | host | 백엔드 호스트 및 포트 </br>지정하지 않으면 서비스에 지정한 host 사용|  |  |
            | method | 백엔드 요청 메서드 </br>지정하지 않으면 endpoint에 지정된 메서드 사용  |  | GET |
            | url_pattern | 백엔드 요청 URL 패턴 | 필수 |  |
            | timeout | 백엔드 처리 제한 시간 |  | 2s |
            | group | 응답 데이터를 지정한 이름으로 묶어서 반환 |  |  |
            | blacklist | 응답 데이터 중에서 제외할 필드들 </br>나머지 필드들은 그대로 반환 됨 |  |  |
            | whitelist | 응답 데이터 중에서 추출할 필드들 </br>나머지 필드들은 모두 제외 됨 |  |  |
            | mapping | 응답 데이터 중에서 지정한 필드를 지정한 이름으로 변경 |  |  |
            | target | 응답 데이터 중에서 지정한 필드만을 반환함 </br>나머지 필드들은 모두 제외 됨 |  |  |
            | wrap_collection_to_json | 응답 결과가 컬랙션인 경우에 JSON 객체로 반환 여부 (true이면 collection 을 "collection" 필드로 JSON 반환, false이면 collection 상태로 반환) |  | false |
            | is_collection | 응답 결과가 JSON객체가 아닌 컬랙션인 경우 ("collection" 필드로 컬랙션을 Wrapping한 JSON 반환하며, mapping 정책에 따라서 필드명 변경 가능) |  | false |

### Bypass 설정하는 방법
  - 위에서 설명한 설정 중에서 Endpoint 와 Backend 설정을 조정해서 사용한다.
  - 적용 예
    ```yaml
    ...
      - endpoint: "/<prefix_url>/*bypass"
        - backend:
            - host: "http//<apiserver_host>:<apiserver_port>"
              url_pattern: "*bypass"
    ...

> Notes
> ---
> - **<font color="red">endpoint 와 url_pattern 에는 `*bypass` 라는 접미사를 사용한다.</font>**
> - 단일 Endpoint 기준으로 동작한다.
> - 각 Endpoint에 대해 단일 Backend 설정만 가능하다.
> - API G/W의 기능인 Filtering 기능 등을 사용할 수 없다. (그대로 전달하는 기능만 가능)
> - 특정 Method로 제한할 수 없기 때문에 전체 Method를 대상으로 운영된다. (실제 API Server에서 해당 Method를 검증해야 한다)

### 현재 지원되는 Middleware 들은 다음과 같다.
- Service 레벨
  - **CORS** : Cross-Origin Resource Sharing 관련 지원
    ```yaml
    middleware:
      mw-cors:
        allow_origins:  # 배열, 와일드카드 사용
          - "*"
        allow_methods:  # 배열
          - POST
          - GET
        allow_headers:  # 배열, 허용할 헤더
          - Origin
          - Authorization
          - Content-Type
        expose_headers: # 배열, 클라이언트가 연결할 수 있도록 노출할 헤더
          - Content-Length
        max_age: 12h    # 캐시 유지 시간
        allow_credentials: true # 브라우저에서 응답에 대한 자격 증명 (쿠키, 인증 헤서, TLS 등)을 자바스크립트에 노출할지 여부
    ```
  - **METRICS/INFLUXDB/Grafana** : Metric 정보 수집 및 저장 지원 및 Grafana Dashboard
    ```yaml
    middleware:
      mw-metrics:
        router_enabled: true    # 라우터 레이어 측정 여부
        proxy_enabled: true     # 프록시 레이어 측정 여부
        backend_enabled: true   # 백엔드 레이어 측정 여부
        collection_period: 10s  # 수집 주기
        expose_metrics: false   # 타 메트릭 수집기 (eg. Prometheus)에 메트릭 정보 노출 여부 (Gin 서버 구동)
        listen_address: 0.0.0.0:9000  # 노출시 사용할 주소
        influxdb:
          address: "http://localhost:8086"    # 저장소 주소
          database: cbrestapigw               # 데이터베이스 명
          reporting_period: 11s               # 수집 데이터 전송 주기
          buffer_size: 0                      # 전송에 사용할 버퍼 크기
    ```
  - **TRACE/OPENCENSUS/JAEGER** : Opencensus 기반의 Trace 수집 및 저장 지원 및 Jaeger UI
    ```yaml
    middleware:
      mw-opencensus:
        sample_rate: 100      # 샘플 비율 (0 - 측정없음 ~ 100 - 전체측정)
        reporting_period: 5s  # 데이터 전송 주기
        enabled_layers:
          router: true        # 라우터 레이어 측정 여부
          proxy: true         # 프록시 레이어 측정 여부
          backend: true       # 백엔드 레이어 측정 여부
        exporters:
          jaeger:
            endpoint: http://localhost:14268/api/traces   # Jaeger Exporter 엔드포인트 주소
            service_name: cbrestapigw                     # Jaeger 식별용 서비스 명
    ```
- Endpoint 레벨
  - AUTH (Simple HMAC)
    ```yaml
    middleware:
      mw-auth:
        secure_key: "###TEST_SECURE_KEY###"   # 해시 생성에 사용할 비밀키
        access_ids:                           # 배열, 해시 인증 후 액세스 허용 아이디 리스트
          - etri
    ```
  - **PROXY (Sequential)**
    ```yaml
    middleware:
      mw-proxy:
        sequential: true        # 지정한 여러 백엔드를 순차적으로 처리할지 여부
    ```
  - **Rate Limit (Endpoint Rate Limit)**
    - 설정이 없거나 0으로 지정된 경우는 무제한 허용
    - Rate Limit는 초당 허용 하는 호출 수를 기준으로 한다. (TokenBucket 알고리즘 적용)
    - 시간은 연속적인 흐름이므로 지정한 호출 수를 기반으로 사용 비율을 계산하여 1개씩의 호출이 가능하도록 추가한다.
    - Rate Limit는 Endpoint 단위 또는 Client 단위로 설정 가능하다.
      - Endpoint 단위 설정
        ```yaml
        middleware:
          mw-ratelimit:
            maxRate: 10   # Endpoint URL 단위로 초당 10개 호출 허용
        ```
      - Client 단위 설정
        - Client IP 단위
          ```yaml
          middleware:
            mw-ratelimit:
              clientMaxRate: 5  # 클라이언트 IP 단위로 초당 5개 허용
              strategy: "ip"
          ```
        - Request에 특정 Header 값을 지정하는 단위
          ```yaml
          middleware:
            mw-ratelimit:
              clientMaxRate: 5  # 클라이언트의 Request Header에 설정된 값을 기준으로 초당 5개 허용
              strategy: "header"
              key: "<header로 전달할 Key 명, ex. 'X-Private-Token'>"
          ```
      - Endpoint 및 Client 모두 설정
        ```yaml
        middleware:
          mw-ratelimit:
            maxRate: 10         # Endpoint URL 단위로 초당 10개 허용
            clientMaxRate: 5    # 클라이언트 Rqeuest Header에 설정된 값을 기준으로 초당 5개 허용
            strategy: "header"
            key: "X-Private-Token"
        ```
    - Endpoint 단위 호출 허용 수를 초과하는 경우는 API G/W 자체가 실패한 것이므로 <font color="red">`503 - Service Unavailable 오류`</font> 상태를 반환한다.
    - Client 단위 호출 허용 수를 초과하는 경우는 특정 사용자의 호출이 실패한 것이므로 <font color="red">`429 - Too many requests 오류`</font> 상태를 반환한다.

- Backend 레벨
  - **HTTPCACHE (Backend Reponse cache)**
    ```yaml
    middleware:
      mw-httpcache: 
        enabled: true     # 응답 캐시 활성화 여부 (In-Memory)
    ```
  - **PROXY (Flatmap filter)** : 응답 결과에 배열이 존재하는 경우에 사용
    ```yaml
    middleware:
      mw-proxy:
        flatmap_filter:         # depth는 "." 을 사용, array는 숫자 인덱스 또는 "*" 사용
          - type: "move"        # args 지정에 따라서 변경
            args:
              - "products.*.id"       # 원본, 응답결과의 products 밑의 모든 배열 중에서 id 선택
              - "products.*.id-"      # 변경, 응답결과의 products 밑의 모든 배열 중에서 id 를 id- 로 변경
              - type: "del"     # args에 지정된 결과를 모두 삭제
              - "products.*.image"
              - "products.*.body_html"
              - "products.*.created_at"
              - "products.*.handle"
              - "products.*.product_type"
              - "products.*.published_at"
              - "products.*.published_scope"
              - "products.*.tags"
              - "products.*.template_suffix"
              - "products.*.updated_at"
              - "products.*.vendor"
    ```
  - **HTTP (Error Details)**
    ```yaml
      ...
        middleware:
          mw-http:
            return_error_details: "test"  # 오류 식별을 위한 문자열
      ...
    ```
  - **Rate Limit (Endpoint Rate Limit)**
    - 설정이 없거나 0으로 지정된 경우는 무제한 허용
    - Rate Limit는 초당 허용 하는 호출 수를 기준으로 한다. (TokenBucket 알고리즘 적용)
    - 시간은 연속적인 흐름이므로 지정한 호출 수를 기반으로 사용 비율을 계산하여 1개씩의 호출이 가능하도록 추가한다.
    - Rate Limit는 Backend 단위로 설정 가능하다.
      - Backend 단위 설정
        ```yaml
        middleware:
          mw-ratelimit:
            maxRate: 10   # Backend URL 단위로 초당 10개 호출 허용
            capacity: 10  # 초당 maxRate 소비 비율로 계산된 구간마다 1개의 토큰을 추가할 수 있는 최대 값 (일반적으로 maxRate == capacity 로 설정)
        ```
    - Rate Limit 가 지정되어 호출이 제한 되는 경우에도 여러 개의 Backend가 존재할 수 있으므로 API G/W가 아닌 Backend 호출에 대한 제한이므로 성공한 Backend가 존재하는 경우라면 `200 정상` 으로 상태 코드를 처리한다.
    - 단, 단일 Backend이며 Rate Limit에 걸리는 경우는 `503, Service unavailable` 로 상태 코드를 처리한다.
    - <font color="red">`단, 제한된 Backend의 경우는 Response Header 정보 ("X-Cb-Restapigw-Completed", "X-Cb-Restapigw-Messages") 를 확인해서 오류 여부를 검증`</font>해야 한다.

### 현재 지원되는 응답 데이터 처리용 필터들은 다음과 같다.

> Notes
> ---
> **<font color="red">Bypass 처리를 한 경우는 특정 Endpoint, Backend를 대상으로 하는 것이 아니므로 응답 데이터 처리를 적용할 수 없다.</font>**

  - **whitelist** : 응답 결과중에서 추출할 필드들 지정, nested field들은 '.' 을 사용해서 설정 가능
    ```yaml
    backend:
      - url_pattern: "/hotels/1.json"
        whitelist:                    # 배열, Array가 존재하는 경우는 flatmap 사용
          - "destination_id"
    ```
    ex)
    ```json
    # 벡엔드 응답 결과
    {
      "hotel_id": 1,
      "name": "Hotel California",
      "destination_id": 1
    }
    # 필터링된 데이터
    {
      "destination_id": 1
    }
    ```
  - **blacklist** : 응답 결과중에서 제외할 필드들 지정, nested field들은 '.' 을 사용해서 설정 가능
    ```yaml
    backend:
      - url_pattern: "/hotels/1.json"
        blacklist:                    # 배열, Array가 존재하는 경우는 flatmap 사용
          - "hotel_id"
          - "name"
    ```
    ex)
    ```json
    # 벡엔드 응답 결과
    {
      "hotel_id": 1,
      "name": "Hotel California",
      "destination_id": 1
    }
    # 필터링된 데이터
    {
      "destination_id": 1
    }
    ```
  - **group** : 응답 결과를 지정한 이름의 묶음 결과로 처리
    ```yaml
    backend:
      - url_pattern: "/hotels/1.json"
        group: "hotel_info"
    ```
    ex)
    ```json
    # 벡엔드 응답 결과
    {
      "hotel_id": 1,
      "name": "Hotel California",
      "destination_id": 1
    }
    # 필터링된 데이터
    hotel_info: {
      "hotel_id": 1,
      "name": "Hotel California",
      "destination_id": 1
    }
    ```
  - **target** : 응답 결과 중에서 특정 필드만 추출할 떄 사용
    ```yaml
    backend:
      - url_pattern: "/destination/1.json"
        target: "destinations"
    ```
    ex)
    ```json
    # 벡엔드 응답 결과
    {
      "destination_id": 1,
      "description": "Top Tourist Attractions in the USA",
      "destinations": [
        "Mount Rushmore",
        "Pike Place Market in Seattle",
        "Venice Beach in LA",
        "Mesa Verde",
        "Faneuil Hall in Boston",
        "Kennedy Space Center",
        "Navy Pier in Chicago",
        "Great Smoky Mountains National Park",
        "River Walk in San Antonio",
        "Carlsbad Caverns",
        "Bryce Canyon National Park",
        "French Quarter in New Orleans",
        "Sedona Red Rock Country",
        "Walt Disney World in Orlando",
        "Yosemite National Park",
        "White House in Washington D.C.",
        "Denali National Park",
        "Las Vegas Strip",
        "Florida Keys",
        "Kilauea",
        "Niagara Falls",
        "Golden Gate Bridge in San Francisco",
        "Yellowstone",
        "Manhattan",
        "Grand Canyon"
      ]
    }
    # 필터링된 데이터
    [
      "Mount Rushmore",
      "Pike Place Market in Seattle",
      "Venice Beach in LA",
      "Mesa Verde",
      "Faneuil Hall in Boston",
      "Kennedy Space Center",
      "Navy Pier in Chicago",
      "Great Smoky Mountains National Park",
      "River Walk in San Antonio",
      "Carlsbad Caverns",
      "Bryce Canyon National Park",
      "French Quarter in New Orleans",
      "Sedona Red Rock Country",
      "Walt Disney World in Orlando",
      "Yosemite National Park",
      "White House in Washington D.C.",
      "Denali National Park",
      "Las Vegas Strip",
      "Florida Keys",
      "Kilauea",
      "Niagara Falls",
      "Golden Gate Bridge in San Francisco",
      "Yellowstone",
      "Manhattan",
      "Grand Canyon"
    ]
    ```
  - **mapping** : 응답 결과 중에서 특정 필드의 이름을 변경할 때 사용
    ```yaml
    backend:
      - url_pattern: "/hotels/1.json"
        mapping:
          "hotel_id": "id"
    ```
    ex)
    ```json
    # 벡엔드 응답 결과
    {
      "hotel_id": 1,
      "name": "Hotel California",
      "destination_id": 1
    }
    # 필터링된 데이터
    {
      "id": 1,
      "name": "Hotel California",
      "destination_id": 1
    }
    ```
  - is_collection: 응답의 결과가 객체가 아닌 컬랙션인 경우 `wrap_collection_to_json` 설정이 true인 경우는 "collection" 이라는 필드의 객체 형식으로 응답을 반환하고, 그 외의 경우는 Array인 상태로 반환한다.
    ```yaml
    backend:
      - url_pattern: "/destinations/2.json"
        is_collection: true
        wrap_collection_to_json: true
    ```
    ex)
    ```json
    # 벡엔드 응답 결과
    [
      "Mount Rushmore",
      "Pike Place Market in Seattle",
      "Venice Beach in LA",
      "Mesa Verde",
      "Faneuil Hall in Boston",
      "Kennedy Space Center",
      "Navy Pier in Chicago",
      "Great Smoky Mountains National Park",
      "River Walk in San Antonio",
      "Carlsbad Caverns",
      "Bryce Canyon National Park",
      "French Quarter in New Orleans",
      "Sedona Red Rock Country",
      "Walt Disney World in Orlando",
      "Yosemite National Park",
      "White House in Washington D.C.",
      "Denali National Park",
      "Las Vegas Strip",
      "Florida Keys",
      "Kilauea",
      "Niagara Falls",
      "Golden Gate Bridge in San Francisco",
      "Yellowstone",
      "Manhattan",
      "Grand Canyon"
    ]
    # 필터링된 데이터 (wrap_collection_to_json = true인 경우)
    {
      "collection": [
        "Mount Rushmore",
        "Pike Place Market in Seattle",
        "Venice Beach in LA",
        "Mesa Verde",
        "Faneuil Hall in Boston",
        "Kennedy Space Center",
        "Navy Pier in Chicago",
        "Great Smoky Mountains National Park",
        "River Walk in San Antonio",
        "Carlsbad Caverns",
        "Bryce Canyon National Park",
        "French Quarter in New Orleans",
        "Sedona Red Rock Country",
        "Walt Disney World in Orlando",
        "Yosemite National Park",
        "White House in Washington D.C.",
        "Denali National Park",
        "Las Vegas Strip",
        "Florida Keys",
        "Kilauea",
        "Niagara Falls",
        "Golden Gate Bridge in San Francisco",
        "Yellowstone",
        "Manhattan",
        "Grand Canyon"
      ]
    }
    ```
  - **flatmap** (결과 JSON의 중첩 구조 및 Array 처리용)
    - flatmap 은 응답 결과에 중첩 구조와 Array 가 존재할 때 사용하는 것으로 이를 사용하면 blacklist나 whitelist는 무시된다.
    - group과 target은 flatmap과 같이 사용할 수 있다.
    - flatmap에서 지원되는 기능은 아래의 2가지가 존재한다.
      - move : 2개의 arguments 를 지정해서 이름을 변경할 때 사용
      - del : arguments로 지정한 데이터를 모두 삭제
    ex)
    ```json
    # 벡엔드 응답 결과
    {
      "campaigns": [
        {
          "discounts": [
            {
              "discount": 0.15,
              "id_product": 1
            },
            {
              "discount": 0.5,
              "id_product": 2
            }
          ],
          "end_date": "2017/02/15",
          "id_campaign": 1,
          "name": "Saint Calentine",
          "start_date": "2017/02/10"
        },
        {
          "discounts": [
            {
              "discount": 0.2,
              "id_product": 1
            },
            {
              "discount": 0.1,
              "id_product": 2
            }
          ],
          "end_date": "2017/09/15",
          "id_campaign": 2,
          "name": "Summer break",
          "start_date": "2017/06/01"
        }
      ],
      "products": [
        {
          "body_html": "<p>It's the small iPod with one very big idea: Carrying files like an animal. Now the world's most popular music player, available in 8PB models, lets you enjoy TV shows, movies, video podcasts, and more. The larger, brighter display means amazing picture quality. In six eye-catching colors, iPod Maño is stunning all around. And with models arting at just $149, little speaks volumes.</p>",
          "created_at": "2017-03-16T13:03:15-04:00",
          "handle": "ipod-nano",
          "id": 1,
          "image": {
              "created_at": "2017-03-16T13:03:15-04:00",
              "id": 850703190,
              "position": 1,
              "product_id": 1,
              "src": "https://cdn.your-site.co/ipod-manyo.png",
              "updated_at": "2017-03-16T13:03:15-04:00"
          },
          "product_type": "Cult Products",
          "published_at": "2007-12-31T19:00:00-05:00",
          "published_scope": "web",
          "tags": "Emotive, Flash Memory, MP3, Music",
          "template_suffix": null,
          "title": "IPod Maño - 8PB",
          "updated_at": "2017-03-16T13:03:15-04:00",
          "vendor": "Apple"
        },
        {
          "body_html": "<p>McBook Er surpasses its previous model by removing the thunderbolt 3 port and adding the ultrafast Stormybolt 1. Conversors from USB -> Thunderbolt 3 -> Stormybolt 1 are sold separately.</p>",
          "created_at": "2017-03-16T13:03:15-04:00",
          "handle": "ipod-touch",
          "id": 2,
          "image": null,
          "product_type": "Cult Products",
          "published_at": "2008-09-25T20:00:00-04:00",
          "published_scope": "global",
          "tags": "",
          "template_suffix": null,
          "title": "McBook Er?",
          "updated_at": "2017-03-16T13:03:15-04:00",
          "vendor": "Apple"
        }
      ]
    }
    # 필터링된 데이터
    {
      "campaigns": [
        {
          "discounts": [
            {
              "discount": 0.15,
              "id_product": 1
            },
            {
              "discount": 0.5,
              "id_product": 2
            }
          ],
          "end_date": "2017/02/15",
          "id_campaign": 1,
          "name": "Saint Calentine",
          "start_date": "2017/02/10"
        },
        {
          "discounts": [
            {
              "discount": 0.2,
              "id_product": 1
            },
            {
              "discount": 0.1,
              "id_product": 2
            }
          ],
          "end_date": "2017/09/15",
          "id_campaign": 2,
          "name": "Summer break",
          "start_date": "2017/06/01"
        }
      ],
      "products": [
        {
          "id-": 1,
          "title": "IPod Maño - 8PB"
        },
        {
          "id-": 2,
          "title": "McBook Er?"
        }
      ]
    }
    ```