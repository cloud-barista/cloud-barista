module github.com/cloud-barista/cb-dragonfly

go 1.15

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)

require (
	github.com/Scalingo/go-utils v5.5.14+incompatible
	github.com/bramvdbogaerde/go-scp v0.0.0-20200119201711-987556b8bdd7
	github.com/cloud-barista/cb-log v0.2.0-cappuccino.0.20201008023843-31002c0a088d // indirect
	github.com/cloud-barista/cb-spider v0.3.0-espresso
	github.com/cloud-barista/cb-store v0.3.0-espresso
	github.com/confluentinc/confluent-kafka-go v1.4.2 // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/google/go-cmp v0.4.0
	github.com/google/uuid v1.1.2
	github.com/influxdata/influxdb v1.7.8 // indirect
	github.com/influxdata/influxdb-client-go v0.0.1
	github.com/influxdata/influxdb1-client v0.0.0-20190809212627-fc22c7df067e
	github.com/labstack/echo/v4 v4.1.10
	github.com/mitchellh/mapstructure v1.3.3
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/shaodan/kapacitor-client v0.0.0-20181228024026-84c816949946
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1
	golang.org/x/crypto v0.0.0-20201012173705-84dcc777aaee
	golang.org/x/sys v0.0.0-20210608053332-aa57babbf139 // indirect
	google.golang.org/grpc v1.33.0
	gopkg.in/confluentinc/confluent-kafka-go.v1 v1.4.2
	gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
