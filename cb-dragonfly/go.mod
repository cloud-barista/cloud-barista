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
	github.com/cloud-barista/cb-log v0.2.0-cappuccino.0.20200913031717-ff545833c178 // indirect
	github.com/cloud-barista/cb-spider v0.2.0-cappuccino.0.20200925073009-73c399c7f818
	github.com/cloud-barista/cb-store v0.2.0-cappuccino.0.20200924125209-c313bd2a3987
	github.com/confluentinc/confluent-kafka-go v1.4.2 // indirect
	github.com/coreos/etcd v3.3.18+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.0
	github.com/google/go-cmp v0.4.0
	github.com/google/uuid v1.1.1
	github.com/influxdata/influxdb v1.7.8 // indirect
	github.com/influxdata/influxdb-client-go v0.0.1
	github.com/influxdata/influxdb1-client v0.0.0-20190809212627-fc22c7df067e
	github.com/labstack/echo/v4 v4.1.10
	github.com/mitchellh/mapstructure v1.3.3
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.10.0 // indirect
	github.com/shaodan/kapacitor-client v0.0.0-20181228024026-84c816949946
	github.com/sirupsen/logrus v1.6.0
	github.com/snowzach/rotatefilehook v0.0.0-20180327172521-2f64f265f58c // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1
	go.etcd.io/etcd v3.3.18+incompatible // indirect
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	google.golang.org/grpc v1.26.0
	gopkg.in/confluentinc/confluent-kafka-go.v1 v1.4.2
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
)
