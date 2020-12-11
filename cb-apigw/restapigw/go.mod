module github.com/cloud-barista/cb-apigw/restapigw

go 1.15

// CB-STORE 관련
replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	github.com/xujiajun/nutsdb v0.5.0 => github.com/xujiajun/nutsdb v0.5.1-0.20200320023740-0cc84000d103
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	github.com/cloud-barista/cb-log v0.2.0-cappuccino.0.20201008023843-31002c0a088d
	github.com/cloud-barista/cb-store v0.2.0-cappuccino.0.20201111072717-b0bb715e2694
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/influxdata/influxdb v1.8.3
	github.com/json-iterator/go v1.1.10
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/afero v1.4.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.7.1
	github.com/ugorji/go v1.1.13 // indirect
	github.com/unrolled/secure v1.0.8
	go.opencensus.io v0.22.5
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9 // indirect
	golang.org/x/sys v0.0.0-20201020230747-6e5568b54d1a // indirect
	google.golang.org/api v0.33.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20201021134325-0d71844de594 // indirect
	google.golang.org/grpc v1.33.1 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)
