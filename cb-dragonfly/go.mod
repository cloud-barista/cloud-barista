module github.com/cloud-barista/cb-dragonfly

go 1.15

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.10
	k8s.io/client-go => k8s.io/client-go v0.19.10
)

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/Scalingo/go-utils v7.1.0+incompatible
	github.com/Workiva/go-datastructures v1.0.53
	github.com/bramvdbogaerde/go-scp v1.0.0
	github.com/cloud-barista/cb-log v0.4.0
	github.com/cloud-barista/cb-spider v0.4.5
	github.com/cloud-barista/cb-store v0.4.1
	github.com/confluentinc/confluent-kafka-go v1.7.0
	github.com/coreos/bbolt v1.3.4 // indirect
	github.com/go-openapi/spec v0.20.6 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.3.0
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/influxdata/influxdb v1.9.2 // indirect
	github.com/influxdata/influxdb1-client v0.0.0-20200827194710-b269163b24ab
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/labstack/echo/v4 v4.9.0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/pkg/errors v0.9.1
	github.com/shaodan/kapacitor-client v0.0.0-20181228024026-84c816949946
	github.com/sirupsen/logrus v1.8.1
	github.com/smartystreets/assertions v1.1.0 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/swaggo/echo-swagger v1.1.0
	github.com/swaggo/swag v1.8.2
	github.com/thoas/go-funk v0.9.2
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200427203606-3cfed13b9966 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/net v0.0.0-20220531201128-c960675eff93 // indirect
	golang.org/x/time v0.0.0-20210611083556-38a9dc6acbc6 // indirect
	google.golang.org/grpc v1.39.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.19.10
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
)
