cb-webtool
==========
cb-webtool은 Multi-Cloud Project의 일환으로 다양한 클라우드를 cb-webtool에서 처리해 <br>
사용자로 하여금 간단하고 편안하게 클라우드를 접할 수 있게 해준다.
***
## [Index]
1. [설치환경](#설치-환경)
2. [소스설치](#소스-설치)
3. [실행준비](#실행-준비)
4. [서버실행](#서버-실행)
***
## [설치 환경]
 - Linux(검증시험 : ubuntu 18.0.4)

## [소스 설치]
 - Git 설치
 - Go설치(1.12이상)
 - echo 설치
 
    ````bash
      $ go get -u -v github.com/labstack/echo
    ````
 
 - echo-session 설치
     ````bash
       $ go get -u -v github.com/go-session/echo-session
     ````
     
 - cloud-barista alliance 설치
    - cb-log install
         ````bash
          $ go get -u -v github.com/cloud-barista/cb-log
         ````
        - [https://github.com/cloud-barista/cb-log](https://github.com/cloud-barista/cb-log) README 참고하여 설치 및 설정
    - cb-store install
        ````bash
         $ go get -u -v github.com/cloud-barista/cb-store
         ````
        - [https://github.com/cloud-barista/cb-store](https://github.com/cloud-barista/cb-store) README 참고하여 설치 및 설정
    - cb-spider install
        ````bash
         $ go get -u -v github.com/cloud-barista/cb-spider
         ````
        - [https://github.com/cloud-barista/cb-spider](https://github.com/cloud-barista/cb-spider) README 참고하여 설치 및 설정
    - cb-tumblebug install
        ````bash
         $ go get -u -v github.com/cloud-barista/cb-tumblebug
         ````
        - [https://github.com/cloud-barista/cb-tumblebug](https://github.com/cloud-barista/cb-tumblebug) README 참고하여 설치 및 설정

## [실행 준비]

   - cb-tumblebug 실행에 필요한 환경변수 설정
       
       ````bash
        $  source setup.env
        ````
   
   - cb-spider 실행에 필요한 환경변수 설정
       
       ````bash
        $  source setup.env
        ````
        
## [서버 실행]

   - cb-tumblebug
    
      ````bash
      $ cd github.com/cloud-barista/cb-tumblebug/src/
      $ ./mcir
      ````
    
   - cb-spider
    
       ````bash
       $ cd github.com/cloud-barista/cb-spider/api-runtime/rest-runtime
       ````
    
       ````bash
       $ go run *.go
       ````
   
   - cb-webtool
   
       ````bash
       $ cd github.com/cloud-barista/cb-webtool
       ````
       
       ````bash
       $ go run main.go
       ````
