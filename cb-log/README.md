# cb-log
CB-Log is the logger library for the Cloud-Barista Multi-Cloud Framework.


# 1.	install CB-Log library pkg
  A.	$ go get github.com/cloud-barista/cb-log

# 2.	example
  A.	https://github.com/cloud-barista/cb-log/blob/master/test/sample.go

# 3.	test example
  A.	$ cd test
  
  B.	$ go run sample.go
  
      …
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
      …
      

  C. set Log Level: info => error   
    i.	$ vi ../conf/config.yaml
    
      …
      [CB-SPIDER].[ERROR]: 2019-08-16 23:22:57 sample.go:69, main.createUser2() - DBMS Session is closed!!

      [CB-SPIDER].[ERROR]: 2019-08-16 23:22:59 sample.go:69, main.createUser2() - DBMS Session is closed!!
      …

