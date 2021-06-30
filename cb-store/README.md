# cb-store
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cloud-barista/cb-store?label=go.mod)](https://github.com/cloud-barista/cb-store/blob/master/go.mod)
[![GoDoc](https://godoc.org/github.com/cloud-barista/cb-store?status.svg)](https://pkg.go.dev/github.com/cloud-barista/cb-store)&nbsp;&nbsp;&nbsp;
[![Release Version](https://img.shields.io/github/v/release/cloud-barista/cb-store)](https://github.com/cloud-barista/cb-store/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/cloud-barista/cb-store/blob/master/LICENSE)

cb-store is a common repository for managing Meta Info of Cloud-Barista.
You can choose NUTSDB or ETCD for repository of cb-store.

-	[NUTSDB](https://github.com/xujiajun/nutsdb): Embedded Key-Value Store on the Local Filesystem.
- [ETCD(Client V3.0)](https://github.com/etcd-io/etcd): Distributed Key-Value Store

```
[NOTE]
cb-store is currently under development. (the latest version is 0.3.0 espresso)
So, we do not recommend using the current release in production.
Please note that the functionalities of cb-store are not stable and secure yet.
If you have any difficulties in using cb-store, please let us know.
(Open an issue or Join the cloud-barista Slack)
```
***

# 1.	install cb-store library pkg
-	$ go get github.com/cloud-barista/cb-store  
 
- $ export CBSTORE_ROOT=~/go/src/github.com/cloud-barista/cb-store
    
- $ vi conf/store_conf.yaml # set up storetype(NUTSDB|ETCD)
  
# 2.	example & test
- example: https://github.com/cloud-barista/cb-store/blob/master/test/test.go
  
- install ETCD (Client V3.0): When using ETCD
    
-	$ cd test  
    
-	$ go run test.go 

      ```
      =========================== Put(...)
      </> root
      </key1> value
      </key1> value1
      </key1/> value2
      </key1/%> value3%
      </key1/key2/key3> value4
      </space key> space value5
      </newline
       key> newline
       value6
      </a/b/c/123/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u> value/value/value
      ===========================
      =========================== Get("/")
      </> root
      ===========================
      =========================== Get("space key")
      </space key> space value5
      ===========================
      =========================== GetList("/", Ascending)
      </> root
      </a/b/c/123/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u> value/value/value
      </key1> value1
      </key1/> value2
      </key1/%> value3%
      </key1/key2/key3> value4
      </newline
       key> newline
       value6
      </space key> space value5
      ===========================
      =========================== GetList("/", Descending)
      </space key> space value5
      </newline
       key> newline
       value6
      </key1/key2/key3> value4
      </key1/%> value3%
      </key1/> value2
      </key1> value1
      </a/b/c/123/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u> value/value/value
      </> root

      ```
