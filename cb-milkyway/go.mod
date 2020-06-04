module github.com/cloud-barista/cb-milkyway

go 1.13

require (
	github.com/cloud-barista/cb-store v0.0.0-20200324092155-64d4df337031
	github.com/cloud-barista/cb-tumblebug v0.0.0-20200429091010-a0a79d3c399f // indirect
	github.com/google/uuid v1.1.1
	github.com/labstack/echo v3.3.10+incompatible
	github.com/mitchellh/go-homedir v1.1.0
	github.com/shirou/gopsutil v2.19.12+incompatible
	github.com/sirupsen/logrus v1.4.2
	github.com/sparrc/go-ping v0.0.0-20190613174326-4e5b6552494c
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.2
)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
