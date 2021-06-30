module github.com/cloud-barista/cb-ladybug

go 1.16

replace (
	//	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/beego/beego/v2 v2.0.1
	github.com/cloud-barista/cb-spider v0.3.19
	github.com/cloud-barista/cb-store v0.3.15
	github.com/go-resty/resty/v2 v2.6.0
	github.com/google/uuid v1.2.0
	github.com/labstack/echo/v4 v4.3.0
	github.com/sirupsen/logrus v1.8.1
	github.com/swaggo/echo-swagger v1.1.0
	github.com/swaggo/swag v1.7.0
)
