## cb-log
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cloud-barista/cb-log?label=go.mod)](https://github.com/cloud-barista/cb-log/blob/master/go.mod)
[![GoDoc](https://godoc.org/github.com/cloud-barista/cb-log?status.svg)](https://pkg.go.dev/github.com/cloud-barista/cb-log)&nbsp;&nbsp;&nbsp;
[![Release Version](https://img.shields.io/github/v/release/cloud-barista/cb-log)](https://github.com/cloud-barista/cb-log/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/cloud-barista/cb-log/blob/master/LICENSE)

cb-log is a logger library for the Cloud-Barista Multi-Cloud Framework.

```
[NOTE]
cb-log is currently under development. (the latest version is 0.3.0 espresso)
So, we do not recommend using the current release in production.
Please note that the functionalities of cb-log are not stable and secure yet.
If you have any difficulties in using cb-log, please let us know.
(Open an issue or Join the cloud-barista Slack)
```
***

#### [목    차]

1. [실행 환경](#실행-환경)
2. [설치 방법](#설치-방법)
3. [설정 방법](#설정-방법)
4. [활용 예시](#활용-예시)
 
***

#### [실행 환경]

- 배포환경: Ubuntu 18.04, Docker 19.03, Go 1.15
- 개발환경: Ubuntu 20.04, Ubuntu 18.04, Debian 10.6, macOS Catalina 10.15, Android 8.1
  - latest Docker, latest Go

#### [설치 방법]
- Go 설치
  ```
  $ sudo apt update
  $ sudo apt install -y make gcc
  $ sudo snap install go --classic
  ```
- cb-log 설치 방법
  - 모듈 다운로드 방법(Go Module mode, default)
    - cb-log 응용 개발
      - 응용 예시: `$ wget https://raw.githubusercontent.com/cloud-barista/cb-log/master/test/test.go`
    - 모듈 초기화: `$ go mod init test.go`
    - 모듈 다운로드: `$ go mod tidy`
    - 다운로드 위치 예시: `/home/ubuntu/go/pkg/mod/github.com/cloud-barista/cb-log@v0.3.1`
  - 소스 다운로드 방법(Old GOPATH mode):
    - 소스 다운로드: `$ go get -u -v github.com/cloud-barista/cb-log`
    - 다운로드 위치 예시: `/home/ubuntu/go/src/github.com/cloud-barista/cb-log`

#### [설정 방법]
- 설정 정보

  | Configurations | Descriptions          | Default |
  |:-------------:|:--------------|:-------------|
  | loopcheck | 설정값 변경시 자동 반영 여부 설정. <br>설정값: true, false | false |
  | loglevel | 로그 레벨 설정. <br>설정값: trace, debug, info, warn, error, fatal, panic | info |
  | logfile | 로그 파일 출력 여부 설정. <br>설정값: true, false | true |
  | logfileinfo: | ----- 이하 logfile true 일때 유효 ----- ||
  | filename | 로그를 저장할 파일 path 및 이름. <br>설정값: {path}logfilename | ./log/cblogs.log |
  | maxsize | 개별 로그 파일 크기. <br>설정값: integer #megabytes | 10 |
  | maxbackups | 로그 파일 개수. <br>설정값: integer #number  | 50 |
  | maxage | 로그 파일 유지 기간. <br>설정값: integer #days  | 31 |


- 설정 파일 예시
  - `$ wget -r -nH --cut-dirs=3 https://raw.githubusercontent.com/cloud-barista/cb-log/master/conf/log_conf.yaml`
  - `$ vi ./conf/log_conf.yaml`
    ```yaml
    #### Config for CB-Log Lib. ####

    cblog:
      ## true | false
      loopcheck: false # This temp method for development is busy wait. cf) cblogger.go:levelSetupLoop().

      ## trace | debug | info | warn/warning | error | fatal | panic
      loglevel: error # If loopcheck is true, You can set this online.

      ## true | false
      logfile: true

    ## Config for File Output ##
    logfileinfo:
      filename: ./log/cblogs.log
      maxsize: 10 # megabytes
      maxbackups: 50
      maxage: 31 # days
    ```

- 설정 적용 방법
  - 서버 재가동: loopcheck=false 설정시
  - 자동 반영: loopcheck=true 설정시

- 설정파일 위치 지정 방법
  - 환경변수 사용 방법: 
    - 환경변수 CBLOG_ROOT 설정: `(ex) export CBLOG_ROOT=$HOME/go/src/github.com/cloud-barista/cb-log`
    - 설정파일 위치: $CBLOG_ROOT/conf/log_conf.yaml
  - 설정파일 지정 방법: 
    - 설정파일 생성: `(ex) /etc/my_conf.yaml`
    - 코드 내 설정파일 위치 설정: GetLoggerWithConfigPath() 이용

      ```go
      import (
        "github.com/cloud-barista/cb-log"
        "github.com/sirupsen/logrus"
      )
      
      var cblogger *logrus.Logger

      func init() {
        // cblog is a global variable.
        cblogger = cblog.GetLoggerWithConfigPath("MY_PROJ", "/etc/my_conf.yaml")
      }
      ```

#### [활용 예시]
- 기본 사용 예시: GetLogger(), SetLevel(), GetLevel(), WithFields()
  - 대상 소스: https://github.com/cloud-barista/cb-log/blob/master/test/test.go
  - 실행 및 결과:
      ```
      $ export CBLOG_ROOT=$HOME/go/src/github.com/cloud-barista/cb-log
      $ cd $CBLOG_ROOT/test
      $ go run test.go   # setup cb-log with $CBLOG_ROOT/conf/log_conf.yaml
    
      ####LogLevel: info
      [CB-SPIDER].[INFO]: 2021-02-14 11:16:49 test.go:24, main.main() - Log Info message
      [CB-SPIDER].[WARNING]: 2021-02-14 11:16:49 test.go:25, main.main() - Log Waring message
      [CB-SPIDER].[ERROR]: 2021-02-14 11:16:49 test.go:26, main.main() - Log Error message
      [CB-SPIDER].[ERROR]: 2021-02-14 11:16:49 test.go:27, main.main() - Log Error message:internal error message

      ####LogLevel: warning
      [CB-SPIDER].[WARNING]: 2021-02-14 11:16:49 test.go:32, main.main() - Log Waring message
      [CB-SPIDER].[ERROR]: 2021-02-14 11:16:49 test.go:33, main.main() - Log Error message
      [CB-SPIDER].[ERROR]: 2021-02-14 11:16:49 test.go:34, main.main() - Log Error message:internal error message

      ####LogLevel: error
      [CB-SPIDER].[ERROR]: 2021-02-14 11:16:49 test.go:40, main.main() - Log Error message
      [CB-SPIDER].[ERROR]: 2021-02-14 11:16:49 test.go:41, main.main() - Log Error message:internal error message

      ####LogLevel: debug
      [CB-SPIDER].[DEBUG]: 2021-02-14 11:16:49 test.go:46, main.main() - WithField 테스트 	[TestField=test]
      [CB-SPIDER].[DEBUG]: 2021-02-14 11:16:49 test.go:48, main.main() - WithFields 테스트 	[Field2=value2,Field3=value3,Field1=value1]
      [CB-SPIDER].[DEBUG]: 2021-02-14 11:16:49 test.go:50, main.main() - WithError 테스트 	[error=테스트 오류]

- 활용 예시: DBMS 응용 로그 예시
  - 설정파일 환경변수 지정 방법
    - 환경변수 설정: `export CBLOG_ROOT=$HOME/go/src/github.com/cloud-barista/cb-log`
    - 대상 소스: https://github.com/cloud-barista/cb-log/blob/master/test/sample.go
    - 실행 및 결과
      ```
      $ export CBLOG_ROOT=$HOME/go/src/github.com/cloud-barista/cb-log   
      $ cd $CBLOG_ROOT/test  
      $ go run sample.go   # setup cb-log with $CBLOG_ROOT/conf/log_conf.yaml
  
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:51 sample.go:25, main.main() - start.........
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:51 sample.go:45, main.createUser1() - start creating user.
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:51 sample.go:59, main.createUser1() - finish creating user.
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:51 sample.go:64, main.createUser2() - start creating user.
      [CB-SPIDER].[ERROR]: 2019-08-16 23:22:51 sample.go:69, main.createUser2() - DBMS Session is closed!!
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:51 sample.go:78, main.createUser2() - finish creating user.
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:51 sample.go:37, main.main() - end.........

      [CB-SPIDER].[INFO]: 2019-08-16 23:22:53 sample.go:25, main.main() - start.........
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:53 sample.go:45, main.createUser1() - start creating user.
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:53 sample.go:59, main.createUser1() - finish creating user.
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:53 sample.go:64, main.createUser2() - start creating user.
      [CB-SPIDER].[ERROR]: 2019-08-16 23:22:53 sample.go:69, main.createUser2() - DBMS Session is closed!!
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:53 sample.go:78, main.createUser2() - finish creating user.
      [CB-SPIDER].[INFO]: 2019-08-16 23:22:53 sample.go:37, main.main() - end.........
      ```

    - Log Level 변경 및 실행 결과: `debug` => `error`   
      ```
      $ cd $CBLOG_ROOT/test
      $ vi $CBLOG_ROOT/conf/log_conf.yaml  ## debug => error
      $ go run sample.go   # setup cb-log with a user defined configuration file in code
      [CB-SPIDER].[ERROR]: 2019-08-16 23:22:57 sample.go:69, main.createUser2() - DBMS Session is closed!!
      [CB-SPIDER].[ERROR]: 2019-08-16 23:22:59 sample.go:69, main.createUser2() - DBMS Session is closed!!
      ```

  - 설정파일 지정 방법: GetLoggerWithConfigPath()
    - 대상 소스: https://github.com/cloud-barista/cb-log/blob/master/test/sample-with-config-path.go
    - 실행 및 결과
      ```
      $ cd ./test
      $ go run sample-with-config-path.go   # setup cb-log with a user defined configuration file in code

      [CB-SPIDER ..\conf\log_conf.yaml]
      [CB-SPIDER].[INFO]: 2020-12-23 17:46:09 sample-with-config-path.go:27, main.main() - start.........
      [CB-SPIDER].[INFO]: 2020-12-23 17:46:09 sample-with-config-path.go:48, main.createUser3() - start creating user.
      [CB-SPIDER].[DEBUG]: 2020-12-23 17:46:09 sample-with-config-path.go:58, main.createUser3() - msg for debugging msg!!
      [CB-SPIDER].[INFO]: 2020-12-23 17:46:09 sample-with-config-path.go:63, main.createUser3() - finish creating user.
      [CB-SPIDER].[DEBUG]: 2020-12-23 17:46:09 sample-with-config-path.go:30, main.main() - msg for debugging msg!!
      [CB-SPIDER].[INFO]: 2020-12-23 17:46:09 sample-with-config-path.go:68, main.createUser4() - start creating user.
      [CB-SPIDER].[ERROR]: 2020-12-23 17:46:09 sample-with-config-path.go:73, main.createUser4() - DBMS Session is closed!!
      [CB-SPIDER].[INFO]: 2020-12-23 17:46:09 sample-with-config-path.go:82, main.createUser4() - finish creating user.
      [CB-SPIDER].[INFO]: 2020-12-23 17:46:09 sample-with-config-path.go:40, main.main() - end.........
      ```
      
    - Log Level 변경 및 실행 결과: `debug` => `error`   
      ```
      $ cd ./test
      $ vi ../conf/log_conf.yaml  ## debug => error
      $ go run sample-with-config-path.go   # setup cb-log with a user defined configuration file in code
      [CB-SPIDER].[ERROR]: 2020-12-23 18:08:12 sample-with-config-path.go:73, main.createUser4() - DBMS Session is closed!!
      [CB-SPIDER].[ERROR]: 2020-12-23 18:08:14 sample-with-config-path.go:73, main.createUser4() - DBMS Session is closed!!
      ```
      
