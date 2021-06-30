```
[NOTE]
cb-webtool is currently under development. (the latest version is 0.3 espresso)
So, we do not recommend using the current release in production.
Please note that the functionalities of cb-webtool are not stable and secure yet.
If you have any difficulties in using cb-webtool, please let us know.
(Open an issue or Join the cloud-barista Slack)
```
***

cb-webtool
==========
cb-webtool은 Multi-Cloud Project의 일환으로 다양한 클라우드를 cb-webtool에서 처리해 사용자로 하여금 간단하고 편안하게 클라우드를 접할 수 있게 해준다.
***
## [Index]
- [cb-webtool](#cb-webtool)
  - [[Index]](#index)
  - [[설치 환경]](#설치-환경)
  - [[의존성]](#의존성)
  - [[소스 설치]](#소스-설치)
  - [[환경 설정]](#환경-설정)
  - [[cb-webtool 실행]](#cb-webtool-실행)
  - [[cb-webtool 실행-reflex 방식]](#cb-webtool-실행-reflex-방식)
***
## [설치 환경]
cb-webtool은 1.16 이상의 Go 버전이 설치된 다양한 환경에서 실행 가능하지만 최종 동작을 검증한 OS는 Ubuntu 18.0.4입니다.

<br>

## [의존성]
cb-webtool은 내부적으로 cb-tumblebug & cb-spider & cb-dragonfly의 개방형 API를 이용하기 때문에 각 서버의 연동이 필요합니다.<br>
- [https://github.com/cloud-barista/cb-tumblebug](https://github.com/cloud-barista/cb-tumblebug) README 참고하여 설치 및 실행 (검증된 버전 : cb-tumblebug v0.2.9)
- [https://github.com/cloud-barista/cb-spider](https://github.com/cloud-barista/cb-spider) README 참고하여 설치 및 실행 (검증된 버전 : cb-spider v0.2.8)
- [https://github.com/cloud-barista/cb-dragonfly](https://github.com/cloud-barista/cb-dragonfly) README 참고하여 설치 및 실행 (검증된 버전 : cb-dragonfly v0.2.8)

<br>

## [소스 설치]
- Git 설치
  - `$ sudo apt update`
  - `$ sudo apt install git`

- Go 1.16 이상의 버전 설치<br>
  go mod 기반의 설치로 바뀌면서 Go 1.16 이상의 버전이 필요합니다.<br>

  2021년 6월 기준으로 apt install golang으로는 구 버전이 설치되기 때문에 https://golang.org/doc/install 사이트에서 1.16 이상의 버전을 직접 설치해야 합니다.<br>
  - `$ wget https://golang.org/dl/go1.16.4.linux-amd64.tar.gz`
  - `$ sudo tar -C /usr/local -xzf go1.16.4.linux-amd64.tar.gz`

- Go 환경 설정  
  - `$ echo "export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin" >> ~/.bashrc`
  - `$ echo "export GOPATH=$HOME/go" >> ~/.bashrc`
  - `$ source ~/.bashrc`
  - `$ go version`
  ```
      go version go1.16.4 linux/amd64
  ```

 - cb-webtool 설치
   - `$ mkdir -p ~/go/src/github.com/cloud-barista`
   - `$ cd ~/go/src/github.com/cloud-barista`
   - `$ git clone https://github.com/cloud-barista/cb-webtool.git`
   - `$ cd cb-webtool`
   - `$ go mod download`
   - `$ go mod verify`

<br>

## [환경 설정]
   - conf/setup.env 파일에서 cb-tumblebug & cb-spider & cb-dragonfly의 실제 URL 정보로 수정합니다.<br><br>
     **[주의사항]**<br> cb-webtool을 비롯하여 연동되는 모든 서버가 자신의 로컬 환경에서 개발되는 경우를 제외하고는 클라이언트의 웹브라우저에서 접근하기 때문에 localhost나 127.0.0.1 주소가 아닌 실제 IP 주소를 사용해야 합니다.

   - 로그인 Id와 Password의 변경은 conf/setup.env 파일의 LoginEmail & LoginPassword 정보를 수정하세요.<br>
     (기본 값은 admin/admin 입니다.)

   - 초기 Data 구축관련<br>
     내부적으로 [cb-spider](https://github.com/cloud-barista/cb-spider)와 [cb-tumblebug](https://github.com/cloud-barista/cb-tumblebug)의 개방형 API를 사용하므로 입력되는 Key Name및 Key Value는 cb-spider 및 cb-tumblebug의 API 문서를 참고하시기 바랍니다.<br>

     **[중요]**<br>
     Cloud Connection 기능을 사용할 수 없으므로 cb-tumblebug의 [활용 예시](https://github.com/cloud-barista/cb-spider#%ED%99%9C%EC%9A%A9-%EC%98%88%EC%8B%9C_)를 참고해서 **[1.configureSpider](https://github.com/cloud-barista/cb-tumblebug#1-%ED%81%B4%EB%9D%BC%EC%9A%B0%EB%93%9C%EC%A0%95%EB%B3%B4-namespace-mcir-mcis-%EB%93%B1-%EA%B0%9C%EB%B3%84-%EC%A0%9C%EC%96%B4-%EC%8B%9C%ED%97%98) 쉘 스크립트를 실행** 하시기 바랍니다.

<br>

## [cb-webtool 실행]
  - 일반 실행 
    - `$ cd ~/go/src/github.com/cloud-barista/cb-webtool`
    - `$ source ./conf/setup.env`
    - `$ go run main.go`
  
<br>

## [cb-webtool 실행-reflex 방식]
reflex를 이용한 static 파일의 자동 변경 감지및 Reload
  - reflex 설치
    - `$ go get github.com/cespare/reflex`
  - cb-webtool 실행
    - `$ cd ~/go/src/github.com/cloud-barista/cb-webtool`
    - `$ source ./conf/setup.env`
    - `$ reflex -r '\.(html|go|js)' -s go run main.go`
