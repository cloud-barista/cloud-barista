package alert

import (
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	"time"

	kclient "github.com/shaodan/kapacitor-client"

	"github.com/cloud-barista/cb-dragonfly/pkg/config"
)

const (
	kapacitorTimeout = 5 * time.Minute
)

//var once sync.Once
//var client *kclient.Client

func newClient() (*kclient.Client, error) {
	var kapacitorPort int
	if config.GetInstance().GetMonConfig().DeployType == types.Dev {
		kapacitorPort = config.GetInstance().Kapacitor.HelmPort // 29092
	} else {
		kapacitorPort = types.KapacitorDefaultPort // 9092
	}
	kapacitorConfig := kclient.Config{
		URL:                fmt.Sprintf("http://%s:%d", config.GetDefaultConfig().Kapacitor.EndpointUrl, kapacitorPort),
		Timeout:            time.Duration(kapacitorTimeout),
		InsecureSkipVerify: true,
	}
	client, err := kclient.New(kapacitorConfig)
	if client != nil {
		if _, _, err := client.Ping(); err != nil {
			util.GetLogger().Error(fmt.Sprintf("failed to ping kapacitor, error=%s", err))
		}
	}
	return client, err
}

func GetClient() *kclient.Client {
	c, err := newClient()
	if err != nil {
		fmt.Println(err)
	}
	return c
}
