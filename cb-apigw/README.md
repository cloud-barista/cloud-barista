# cb-apigw
cb-apigw is the API Gateway for Cloud-Barista. There are two API Gateway types: REST and gRPC API Gateways.

- API Gateway for REST API is cb-restapigw (released an initial version)
- API Gateway for gRPC API is cb-grpcapigw (have a plan to develop)


# cb-restapigw : REST API Gateway in `cb-apigw`

cb-restapigw는 PoC (Proof of Concepts) 수준의 RESTful API Gateway 기능을 제공한다.

```
[NOTE]
cb-restapigw is currently under development. (the latest version is 0.2 cappuccino)
So, we do not recommend using the current release in production.
Please note that the functionalities of cb-restapigw are not stable and secure yet.
If you have any difficulties in using cb-restapigw, please let us know.
(Open an issue or Join the cloud-barista Slack)
```


# [ 목차 ]

- [컨테이너 기반 실행](#컨테이너-기반-실행)
- Cloud-Barista 시스템 통합 실행 참고 (Docker-Compose 기반)
- [소스 기반 설치 및 실행](#소스-기반-설치-및-실행)
  - [설치](#설치)
  - [설정](#설정)
  - [실행](#실행)

# [컨테이너 기반 실행]
- cb-restapigw 이미지 확인(https://hub.docker.com/r/cloudbaristaorg/cb-restapigw/tags)
- cb-restapigw 컨테이너 실행

```
docker run -p 8000:8000 --name cb-restapigw \
-v /root/go/src/github.com/cloud-barista/cb-apigw/restapigw/conf:/app/conf \
cloudbaristaorg/cb-restapigw:v0.1-yyyymmdd
```

# [Cloud-Barista 시스템 통합 실행 참고 (Docker-Compose 기반)]

```
# git clone https://github.com/jihoon-seo/cb-deployer.git
# cd cb-deployer
# docker-compose up
```

# [소스 기반 설치 및 실행]

## [설치]

설치는 Ubuntu Latest 버전을 기준으로 한다.

- **Git 설치**
  ```shell
  # apt update
  # apt install git
  ```

- **Go 설치 (v1.12 이상)**
  - https://golang.org/dl 에서 최신 버전 확인 (현재 1.14.1)
  - 다운로드 및 설치
    ```shell
    $ wget https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz
    $ tar -C /usr/local -xzf go1.13.4.linux-amd64.tar.gz
    $ export PATH=$PATH:/usr/local/go/bin
    $ which go
    /usr/local/go/bin/go
    $ go version
    go version go1.14.1 linux/amd64
    $ go env
    GO111MODULE="on"
    GOARCH="amd64"
    GOBIN=""
    GOCACHE="/root/.cache/go-build"
    GOEXE=""
    GOFLAGS=""
    GOHOSTARCH="amd64"
    GOHOSTOS="linux"
    GOOS="linux"
    GOPATH="/root/go"
    GOPROXY=""
    GORACE=""
    GOROOT="/usr/local/go"
    GOTMPDIR=""
    GOTOOLDIR="/usr/local/go/pkg/tool/linux_amd64"
    GCCGO="gccgo"
    CC="gcc"
    CXX="g++"
    CGO_ENABLED="1"
    GOMOD=""
    CGO_CFLAGS="-g -O2"
    CGO_CPPFLAGS=""
    CGO_CXXFLAGS="-g -O2"
    CGO_FFLAGS="-g -O2"
    CGO_LDFLAGS="-g -O2"
    PKG_CONFIG="pkg-config"
    GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build341274058=/tmp/go-build -gno-record-gcc-switches"
    ```
  - 환경 파일을 통해서 설정할 경우는 다음과 같이 처리한다.
    - .bashrc 파일 하단에 경로 관련 추가
      ```txt
      ...
      export PATH=$PATH:/usr/local/go/bin
      ```
    - 적용을 위한 bash 재 실행
      ```shell
      $ source ~/.bashrc
      $ . ~/.bashrc
      ```

- **소스 다운로드**
  ```shell
  # git clone https://github.com/cloud-barista/cb-apigw.git
  ```

- **빌드**
  - Mac 환경
    ```shell
    $ cd cb-apigw/restapigw
    $ make build
    ```
  - Linux 환경
    ```shell
    $ cd cb-apigw/restapigw
    $ go build -tags cb-restapigw -o cb-restapigw -v
    ```

## [설정]

Configuration 설정은 `YAML` 포맷을 사용한다.

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
          | endpoint | 클라이언트에 노출할 URL (".../*bypass" 로 지정하면 API G/W의 기능을 사용하지 않는 Bypass 처리로 동작) | 필수 |  |
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
            | url_pattern | 백엔드 요청 URL 패턴 ("bypass" 로 지정하면 API G/W의 기능을 사용하지 않는 Bypass 처리로 동작)| 필수 |  |
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
  - is_collection: 응답의 결과가 객체가 아닌 컬랙션인 경우 `wrap_collection_to_json` 설정이 true인 경우는 "collection" 이라는 필드의 객체 형식으로 응답을 반환하고, 그 외의 경우는 Array인 상태로 반환한다. `단 "collection" 필드명은 Mapping 처리를 통해서 다른 이름으로 변경 가능하다`
    ```yaml
    backend:
      - url_pattern: "/destinations/2.json"
        is_collection: true
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

### 구성 샘플 (소스 상의 ./conf/*.yaml 경로의 샘플들 참조)
  ```yaml
  version: 1
  name: cb-restapigw
  port: 8000
  cache_ttl: 3600s
  timeout: 3s
  debug: true
  host:
    - "http://localhost:8100"
  middleware:
    mw-metrics:
      router_enabled: true
      proxy_enabled: true
      backend_enabled: true
      collection_period: 10s
      expose_metrics: false
      listen_address: 0.0.0.0:9000
      influxdb:
        address: "http://localhost:8086"
        database: cbrestapigw
        reporting_period: 11s
        buffer_size: 0 
    mw-opencensus:
      sample_rate: 100
      reporting_period: 1s
      enabled_layers:
        router: true
        proxy: true
        backend: true
      exporters:
        jaeger:
          endpoint: http://localhost:14268/api/traces
          service_name: cbrestapigw
    mw-cors:
      allow_origins:
        - "*"
      allow_methods:
        - POST
        - GET
      allow_headers:
        - Origin
        - Authorization
        - Content-Type
      expose_headers:
        - Content-Length
      max_age: 12h
      allow_credentials: true
  endpoints:
    - endpoint: "/splash"
      backend:
        - url_pattern: "/shop/campaigns.json"
          whitelist:
            - "campaigns"
        - url_pattern: "/shop/products.json"
          middleware:
            mw-proxy:
              flatmap_filter:
                - type: "move"
                  args:
                    - "products.*.id"
                    - "products.*.id-"
                - type: "del"
                  args:
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
    - endpoint: "/sequential"
      backend:
        - url_pattern: "/hotels/1.json"
          whitelist:
            - "destination_id"
        - url_pattern: "/destinations/{resp0_destination_id}.json"
          middleware:
            mw-http:
              return_error_details: "sequential"
      middleware:
        mw-proxy:
          sequential: true
    - endpoint: "/fail"
      backend:
        - url_pattern: "/user/1.json"
          group: "user"
          target: "address"
        - host:
            - "http://fake_url_that_should_not_resolve.tld"
          url_pattern: "/"
          group: "none"
          middleware:
            mw-http:
              return_error_details: "fail"
    - endpoint: "/public"
      method: GET
      backend:
        - host: 
            - "https://api.github.com"
          url_pattern: "/users/ccambo"
          whitelist:
            - "avatar_url"
            - "name"
            - "company"
            - "blog"
            - "location"
            - "mail"
            - "hireable"
            - "followers"
            - "public_repos"
            - "public_gists"
          mapping:
            "blog": "website"
          group: "github"
          middleware:
            mw-httpcache: 
              enabled: true
        - host: 
            - "https://api.bitbucket.org"
          url_pattern: "/2.0/users/kpacha"
          whitelist:
            - "links.avatar"
            - "display_name"
            - "website"
            - "location"
          mapping: 
            "display_name": "name"
          group: "bitbucket"
    - endpoint: "/github/{user}"
      method: GET
      backend:
        - host:
            - https://api.github.com
          url_pattern: "/users/{user}"
          disable_host_sanitize: true
    - endpoint: "/collection"
      method: GET
      backend:
        - url_pattern: "/destinations/2.json"
          wrap_collection_to_json: true
          is_collection: true
          mapping:
            "collection": "data"
    - endpoint: "/private/custom"
      backend:
        - url_pattern: "/user/1.json"
      middleware:
        mw-auth:
          secure_key: "###TEST_SECURE_KEY###"
          access_ids:
            - etri
            - acorn
  ```

## [실행]

### Background 서비스들 실행

- 내부 API 서버 (Fake API) 는 **jaxgeller/lwan** Docker image를 사용해서 Fake API로 사용.
- API G/W Metrics는 **InfluxDB + Grafana** 를 사용.
- API G/W Trace 정보는 **Opencensus + Jaeger** 를 사용.

API G/W 실행 테스트를 위한 백그라운드 서비스들은 `Deploy` 폴더에 구성되어 있으므로 이를 활용한다.

- <b>테스트를 위한 설정은 /deploy/docker-compose.yaml을 기준으로 Fake-API 부분을 용도에 맞도록 변경하고 설정을 맞춰서 사용.</b>
- <b>HMAC 관련된 Server 기능은 내부 테스트를 위한 것으로 공식적으로는 지원하지 않음.</b>

실행 방법은 deploy/READ.md 참조

### 소스 빌드 및 실행

- **실행 명령**
  - 바이너리 빌드
    - Mac 환경
      ```shell
      # make build
      ```
    - Linux 환경
      ```shell
      # go build -tags cb-restapigw -o cb-restapigw -v
      ```
  - 설정 검사
    - 지정한 설정 파일의 문법 검사
      - Mac 환경
      ```shell
      # make build-check
      ```
    - Linux 환경
      ```shell
      # ./cb-restapigw -c [configuration file] -d check
      ```
  - 실행
    - 지정한 설정파일의 문법 검사 및 서비스 실행
      - Mac 환경
        ```shell
        # make build-run
        ```
      - Linux 환경
        ```shell
        # ./cb-restapigw -c [configuration file] -d -p 8000 run
        ```

- **실행 옵션**
  - `-c` : 특정 경로에 존재하는 설정 파일을 지정. (eg. ./conf/cb-restapigw.yaml)
  - `-d` : Debug 모드, 설정에 지정된 것을 무시하고 적용한다.
  - `-p` : 서비스를 실행할 HTTP Server Port 지정. (eg. 8000), 설정에 지정된 것을 무시하고 적용한다.

- **클라이언트의 테스트**<br/>
  - [Postman으로 작성된 문서](https://documenter.getpostman.com/view/1735092/SW15wbJf?version=latest#c720d518-1830-4283-b512-5153ef879747)를 참고
  - _**HMAC 적용 부분은 내부 검증용으로 공식적으로는 지원하지 않음**_
  - API 호출의 결과는 CB-RESTAPIGW 수행 기준으로 판단한다.

- **클라이언트의 결과 확인**

  API Gateway 운영에 대한 부분과 Backend 호출 부분이 존재하므로 다음과 같은 규칙을 적용한다.

  - API Gateway 운영 문제
    - 정상 (200, 201), Backend 호출 결과와는 무관
    - 오류 (500)
    - 권한 (401) 
  - Backend 운영 문제
    - Response Header의 처리 완료 여부 **`(X-Cb-Restapigw-Completed)`**
      - true : 모든 Backend 처리 성공
      - false : 일부 또는 모든 Backend 처리 실패
    - Response Header의 오류 메시지 여부 **`(X-Cb-Restapigw-Messages)`**
      - 각 Backend 별 발생한 오류 메시지를 **"\n"** 구분자로 연결한 문자열 처리
      - 모두 정상이라면 **`(X-Cb-Restapigw-Messages)`** Header 정보가 존재하지 않는다.

### Docker Container 실행

백그라운드 서비스들을 구동한 후에 API G/W를 Docker 기반으로 생성하여 실행한다.

1. Docker Image 생성
   ```shell
   docker build -t cb-restapigw -f ./Dockerfile .
   ```

2. Docker Contaienr 실행
   ```shell
   docker run --network deploy_default -p 8000:8000 cb-restapigw
   ```
   * 상기 명령어의 `--network deploy_default` 는 Background 서비스가 docker-compose로 동작하면서 구성된 Docker Bridge Network의 이름이다. 별도 옵션을 주지 않았기 때문에 folder 명을 기준으로 생성된 이름을 가진다.
