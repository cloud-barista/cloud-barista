# cb-milkyway (CB-Tumblebug Benchmark Agent)

## 개요
CB-Tumblebug의 최적 멀티 클라우드 인프라 배치 기능을 위한, 벤치마크 에이전트 (PoC)

(CB-Tumblebug: https://github.com/cloud-barista/cb-tumblebug)

```
[NOTE]
cb-milkyway is currently under development as a PoC. (the latest version is 0.2 cappuccino)
So, we do not recommend using the current release in production.
Please note that the functionalities of cb-milkyway are not stable and secure yet.
If you have any difficulties in using cb-milkyway, please let us know.
(Open an issue or Join the cloud-barista Slack)
```

- Sysbench (https://github.com/akopytov/sysbench) 를 활용하여 컴퓨팅 관련 성능 측정
- Ping 을 활용하여 네트워크 지연 성능 측정
- Ubuntu 18.04에서만 동작 테스트 완료

## 사용 방법

1. 벤치마크 실행이 필요한 컴퓨팅 머신(e.g. VM)에서 milkyway 실행 (API 서버 동작)
2. 클라이언트 또는 에이전트 관리기(CB-Tumblebug)에서 API Call을 통해 벤치마크 실행
   - 지원 기능
     - 벤치마킹 SW (sysbench) 설치
     - 벤치마킹 환경 자동 구성 (test files 생성, db table 및 records 생성)
     - 벤치마킹 항목: CPU (Prime number 계산), Memory (read, write), FileOI (read, write), DB-OLTP transactions (read, write)
     - 벤치마킹 항목: RTT to an address, RTT to multiple addresses
     - 벤치마킹 환경 정리 (test files 삭제, db table 및 records 삭제)

### 소스 코드로 실행
```Shell
# git clone https://github.com/cloud-barista/cb-milkyway.git
# cd cb-milkyway/src/
# go run milkyway.go
```

### 바이너리로 실행
```Shell
wget https://github.com/cloud-barista/cb-milkyway/raw/master/src/milkyway && sudo chmod 755 ~/milkyway && ~/milkyway
```

## 실행 예시 

- milkyway 실행
  - 1324 포트에서 API 서버가 실행됨

```Shell
# ~/go/src/github.com/cloud-barista/cb-milkyway/src$ go build -o milkyway && ./milkyway 

 ██████╗██████╗       ███╗   ███╗██╗██╗     ██╗  ██╗██╗   ██╗██╗    ██╗ █████╗ ██╗   ██╗
██╔════╝██╔══██╗      ████╗ ████║██║██║     ██║ ██╔╝╚██╗ ██╔╝██║    ██║██╔══██╗╚██╗ ██╔╝
██║     ██████╔╝█████╗██╔████╔██║██║██║     █████╔╝  ╚████╔╝ ██║ █╗ ██║███████║ ╚████╔╝ 
██║     ██╔══██╗╚════╝██║╚██╔╝██║██║██║     ██╔═██╗   ╚██╔╝  ██║███╗██║██╔══██║  ╚██╔╝  
╚██████╗██████╔╝      ██║ ╚═╝ ██║██║███████╗██║  ██╗   ██║   ╚███╔███╔╝██║  ██║   ██║   
 ╚═════╝╚═════╝       ╚═╝     ╚═╝╚═╝╚══════╝╚═╝  ╚═╝   ╚═╝    ╚══╝╚══╝ ╚═╝  ╚═╝   ╚═╝                    

 Benchmark Agent for CB-Tumblebug
 ________________________________________________
 Version: Cappuccino
 Repository: https://github.com/cloud-barista/cb-milkyway

⇨ http server started on [::]:1324
```

- 클라이언트 또는 에이전트 관리기(CB-Tumblebug)에서 API Call을 통해 벤치마크 시험.
  - cb-milkyway/test$ ./full_test.sh 를 통해 전체 시험 가능
  - ./full_test.sh {milkyway가 동작 중인 host address} {"install"을 입력하면 환경 세팅도 함께 진행}}

```Shell
# ~/go/src/github.com/cloud-barista/cb-milkyway/test$ ./full_test.sh localhost
####################################################################
{
   "result" : "The init is complete",
   "unit" : "",
   "elapsed" : "2.764198",
   "desc" : "128 files, 400Kb each, 50Mb total, 100000 records into 'sbtest1 are created"
}
#-----------------------------
{
   "result" : "0.035",
   "unit" : "ms",
   "elapsed" : "9.206625",
   "desc" : "Average RTT to localhost"
}
#-----------------------------
{
   "resultarray" : [
      {
         "unit" : "ms",
         "elapsed" : "9.205894",
         "result" : "0.025",
         "desc" : "Average RTT to localhost"
      },
      {
         "elapsed" : "18.421540",
         "result" : "0.033",
         "desc" : "Average RTT to localhost",
         "unit" : "ms"
      }
   ]
}
#-----------------------------
{
   "elapsed" : "10.017117",
   "desc" : "Verify prime numbers in 10000 (standard division of each number by all numbers between 2 and the square root of the number)",
   "result" : "9.9914",
   "unit" : "sec"
}
#-----------------------------
{
   "desc" : "Allocate 10G memory buffer and read (repeat reading a pointer)",
   "unit" : "MiB/sec",
   "elapsed" : "1.806758",
   "result" : "5710.32"
}
#-----------------------------
{
   "desc" : "Allocate 10G memory buffer and write (repeat writing a pointer)",
   "elapsed" : "2.082393",
   "unit" : "MiB/sec",
   "result" : "4937.90"
}
#-----------------------------
{
   "result" : "6899.90",
   "unit" : "MiB/sec",
   "desc" : "Check read throughput by excuting random reads for files in 50MiB for 30s",
   "elapsed" : "30.013554"
}
#-----------------------------
{
   "unit" : "MiB/sec",
   "elapsed" : "30.014258",
   "result" : "48.96",
   "desc" : "Check write throughput by excuting random writes for files in 50MiB for 30s"
}
#-----------------------------
{
   "desc" : "Read transactions by simulating transaction loads (OLTP) in DB for 100000 records",
   "unit" : "Transactions/s",
   "result" : "471.63",
   "elapsed" : "10.026478"
}
#-----------------------------
{
   "result" : "136.93",
   "unit" : "Transactions/s",
   "desc" : "Write transactions by simulating transaction loads (OLTP) in DB for 100000 records",
   "elapsed" : "10.052215"
}
#-----------------------------
{
   "unit" : "",
   "result" : "The cleaning is complete",
   "desc" : "The benchmark files and tables are removed",
   "elapsed" : "0.067219"
}
#-----------------------------
```
