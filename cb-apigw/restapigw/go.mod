module github.com/cloud-barista/cb-apigw/restapigw

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.1.0
	github.com/cloud-barista/cb-log v0.0.0-20190829061936-c402c97c951a
	github.com/gin-contrib/sse v0.0.0-20170109093832-22d885f9ecc7 // indirect
	github.com/gin-gonic/gin v1.1.5-0.20170702092826-d459835d2b07
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/influxdata/influxdb v1.7.8
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.4.2
	github.com/snowzach/rotatefilehook v0.0.0-20180327172521-2f64f265f58c // indirect
	github.com/spf13/cobra v0.0.6
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/ugorji/go v1.1.7 // indirect
	github.com/unrolled/secure v1.0.4
	go.opencensus.io v0.22.1
	golang.org/x/net v0.0.0-20190921015927-1a5e07d1ff72 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20200317113312-5766fd39f98d // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v3 v3.0.0-20191120175047-4206685974f2
)

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
