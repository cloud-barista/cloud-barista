```
[NOTE]
cb-webtool is currently under development. (the latest version is 0.2 cappuccino)
So, we do not recommend using the current release in production.
Please note that the functionalities of cb-webtool are not stable and secure yet.
If you have any difficulties in using cb-webtool, please let us know.
(Open an issue or Join the cloud-barista Slack)
```
***

cb-webtool
==========
cb-webtool은 Multi-Cloud Project의 일환으로 다양한 클라우드를 cb-webtool에서 처리해 <br>
사용자로 하여금 간단하고 편안하게 클라우드를 접할 수 있게 해준다.
***
## [Index]
1. [설치 환경](#설치-환경)
2. [의존성](#의존성)
3. [소스 설치](#소스-설치)
3. [환경 설정](#환경-설정)
4. [서버 실행](#서버-실행)
***
## [설치 환경]
cb-webtool은 1.12 이상의 Go 버전이 설치된 다양한 환경에서 실행 가능하지만 최종 동작을 검증한 OS는 Ubuntu 18.0.4입니다.<br>
<br>

## [의존성]
cb-webtool은 내부적으로 cb-tumblebug & cb-spider & cb-dragonfly 프로젝트를 이용하기 때문에, 각 프로젝트들의 문서를 참고하셔서 동일한 서버 또는 독립 서버에 미리 설치 및 실행해야 합니다.<br>
- [https://github.com/cloud-barista/cb-tumblebug](https://github.com/cloud-barista/cb-tumblebug) README 참고하여 설치 및 실행
- [https://github.com/cloud-barista/cb-spider](https://github.com/cloud-barista/cb-spider) README 참고하여 설치 및 실행
- [https://github.com/cloud-barista/cb-dragonfly](https://github.com/cloud-barista/cb-dragonfly) README 참고하여 설치 및 실행

## [소스 설치]
- Git 설치
  - `# apt update`
  - `# apt install git`

- Go 설치
  - https://golang.org/doc/install <br>
    (2020년 05월 현재 `apt install golang` 명령으로 설치하면 1.10 버전이 설치되므로 위 링크에서 1.12 이상의 버전으로 설치할 것)
  - `wget https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz`
  - `tar -C /usr/local -xzf go1.13.4.linux-amd64.tar.gz`
  - `.bashrc` 파일 하단에 다음을 추가: 
  ```
  export PATH=$PATH:/usr/local/go/bin
  export GOPATH=$HOME/go
  ```

- `.bashrc` 에 기재한 내용을 적용하기 위해, 다음 중 하나를 수행
  - bash 재기동
  - `source ~/.bashrc`
  - `. ~/.bashrc`

 - echo 설치
    ````bash
      $ go get -u -v github.com/labstack/echo
    ````
 
 - echo-session 설치
     ````bash
       $ go get -u -v github.com/go-session/echo-session
     ````

 - reflex 설치 (Windows 미지원 / Windows에 bash 설치 시 사용 가능)
     ````bash
       $ go get github.com/cespare/reflex 
     ````

 - cb-webtool 설치
     ````bash
       $ go get github.com/cloud-barista/cb-webtool
     ````

## [환경 설정]
   - conf/setup.env 파일에서 cb-tumblebug & cb-spider & cb-dragonfly의 실제 URL 정보로 수정합니다.<br>
     **[주의사항]** localhost나 127.0.0.1 주소를 사용할 수 없습니다.

   - conf/setup.env 파일에서 cb-webtool에 로그인할 사용자의 LoginEmail & LoginPassword 정보를 수정하세요.<br>

   - 초기 Data 구축<br>
     내부적으로 cb-spider와 cb-tumblebug을 이용하기 때문에 cb-spider의 [API규격](https://github.com/cloud-barista/cb-spider#api-%EA%B7%9C%EA%B2%A9)을 참고해서 JSON 방식의 REST 호출로 데이터를 구축하거나 [활용 예시](https://github.com/cloud-barista/cb-spider#%ED%99%9C%EC%9A%A9-%EC%98%88%EC%8B%9C_)를 참고해서 제공되는 쉘 스크립트 기반의 시험 도구를 이용해서 손쉽게 기초 데이터의 구축이 가능합니다.<br>
     **Network/Security Group/Image/Spec/Keypair는 cb-webtool v0.2.0-cappuccino에서 지원하지 않으므로 현재는 외부에서 생성해야 합니다.**

  - Credential 정보<br>
    Credential 정보의 경우 [cb-tumblebug](https://github.com/cloud-barista/cb-tumblebug)에서 각 CSP 드라이버마다 설정해야하는 Key 값들이 다르기 때문에 설정해야하는 키 값을 모를 경우 [cb-tumblebug](https://github.com/cloud-barista/cb-tumblebug)이나 [cb-spider](https://github.com/cloud-barista/cb-spider)의 [활용 예시](https://github.com/cloud-barista/cb-spider#%ED%99%9C%EC%9A%A9-%EC%98%88%EC%8B%9C_)에 있는 시험 도구 중 Credential 정보를 확인하시기 바랍니다.

  
## [서버 실행]
- Linux & Mac OS에서 실행
    ````bash (Linux & Mac OS)
    $ cd github.com/cloud-barista/cb-webtool
    $ run.sh
    ````

- Bash를 설치하지 않은 Windows 환경에서는 reflex를 사용할 수 없으므로 직접 구동해야 합니다.
    ````bash (Windows)
    $ cd github.com/cloud-barista/cb-webtool
    $ run-windows.sh
    ````
