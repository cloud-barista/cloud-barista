# cb-log
CB-Log is the logger library for the Cloud-Barista Multi-Cloud Framework.


# 1.	install CB-Log library pkg
  A.	$ go get github.com/cloud-barista/cb-log

  - cf) if meet errors
    - error msg: "gosrc/src/go.etcd.io/etcd/vendor/google.golang.org/grpc/clientconn.go:49:2: use of internal package google.golang.org/grpc/internal/resolver/dns not allowed"    
    - sol: $ rm -rf $GOPATH/src/go.etcd.io/etcd/vendor/google.golang.org/grpc
  
  B.  export CBLOG_ROOT=$GOPATH/src/github.com/cloud-barista/cb-log

  - cf) if meet errors
    - error msg: "go run sample.go 2020/10/13 22:42:13 error: open /src/github.com/cloud-barista/cb-log/conf/log_conf.yaml: no such file or directory exit status 1"
    - sol: 
      - 1) check there is $CBLOG_ROOT  (cd $CBLOG_ROOT shoud be equal to $HOME/go/src/github.com/cloud-barista/cb-log. It means that you can see conf folder when you type cd $CBLOG_ROOT)
      - 2) check your shell type and .~rc file matched. (e.g: bash & .bashrc, zsh & .zshrc ...)
    
# 2.	example
  A.	https://github.com/cloud-barista/cb-log/blob/master/test/sample.go

# 3.	test example
  A.	$ cd $CBLOG_ROOT/test
  
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
      

  C. set Log Level: `debug` => `error`   
    i.	$ vi ../conf/log_conf.yaml
    
      …
      [CB-SPIDER].[ERROR]: 2019-08-16 23:22:57 sample.go:69, main.createUser2() - DBMS Session is closed!!

      [CB-SPIDER].[ERROR]: 2019-08-16 23:22:59 sample.go:69, main.createUser2() - DBMS Session is closed!!
      …

