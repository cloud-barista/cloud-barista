package push

import (
	"context"
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/modules/procedure/push/collector"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strconv"
	"strings"
	"sync"
)

type CollectManager struct {
	// Usage of CollectorAddrSlice => DeployType ==Dev, Compose
	CollectorAddrSlice []*collector.MetricCollector
	CollectorPolicy    string
	K8sClientSet       *kubernetes.Clientset
	WaitGroup          *sync.WaitGroup
}

func NewCollectorManager() (*CollectManager, error) {

	manager := CollectManager{}
	if config.GetInstance().Monitoring.DeployType == types.Helm {
		if err := manager.InitDFK8sEnv(); err != nil {
			return &manager, err
		}
	}
	// collectorPolicy: 1. MaxAgentHost 2. CSP
	manager.CollectorPolicy = strings.ToUpper(config.GetInstance().Monitoring.MonitoringPolicy)
	manager.CollectorAddrSlice = []*collector.MetricCollector{}

	return &manager, nil
}

func (manager *CollectManager) InitDFK8sEnv() (err error) {
	inClusterK8sConfig, errK8s := rest.InClusterConfig()
	if errK8s != nil {
		err = errK8s
	}
	clientSet, errK8s2 := kubernetes.NewForConfig(inClusterK8sConfig)
	if errK8s2 != nil {
		err = errK8s2
	}
	manager.K8sClientSet = clientSet

	// Deploy ConfigMap
	configMapsClient := manager.K8sClientSet.CoreV1().ConfigMaps(types.Namespace)
	configMap := &apiv1.ConfigMap{Data: map[string]string{}, ObjectMeta: metav1.ObjectMeta{
		Name: types.ConfigMapName,
	}}
	//If There is No ConfigMap resource, than Deploy
	_, err = configMapsClient.Get(context.TODO(), types.ConfigMapName, metav1.GetOptions{})
	if err != nil {
		result, err := configMapsClient.Create(context.TODO(), configMap, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		fmt.Println("Created ConfigMap: ", result.GetObjectMeta().GetName())
	}
	return
}

func (manager *CollectManager) CreateCollector() error {

	manager.WaitGroup.Add(1)
	collectorCreateOrder := len(manager.CollectorAddrSlice)
	newCollector, err := collector.NewMetricCollector(
		types.AVG,
		collectorCreateOrder,
	)
	if err != nil {
		return err
	}
	manager.CollectorAddrSlice = append(manager.CollectorAddrSlice, &newCollector)

	switch config.GetInstance().Monitoring.DeployType {
	case types.Helm:
		collectorUUID := fmt.Sprintf("%p", &newCollector)
		env := []apiv1.EnvVar{
			{Name: "kafka_endpoint_url", Value: config.GetInstance().Kafka.EndpointUrl},
			{Name: "create_order", Value: strconv.Itoa(collectorCreateOrder)},
			{Name: "namespace", Value: types.Namespace},
			{Name: "df_addr", Value: fmt.Sprintf("%s:%d", config.GetInstance().Dragonfly.DragonflyIP, config.GetInstance().Dragonfly.HelmPort)},
			{Name: "collect_interval", Value: strconv.Itoa(config.GetInstance().Monitoring.CollectorInterval)},
			{Name: "collect_uuid", Value: collectorUUID},
		}
		deploymentTemplate := util.DeploymentTemplate(collectorCreateOrder, collectorUUID, env)
		fmt.Println("Creating deployment...")
		result, err := manager.K8sClientSet.AppsV1().Deployments(types.Namespace).Create(context.TODO(), deploymentTemplate, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		fmt.Println("Created deployment: ", result.GetObjectMeta().GetName())
		return nil
	case types.Dev, types.Compose:
		go func() {
			err := newCollector.Collector(manager.WaitGroup)
			if err != nil {
				util.GetLogger().Error("failed to  create Collector")
			}
		}()
	}
	return nil
}

func (manager *CollectManager) DeleteCollector() error {
	lastCollectorIdx := len(manager.CollectorAddrSlice) - 1
	cAddr := manager.CollectorAddrSlice[lastCollectorIdx]
	switch config.GetInstance().Monitoring.DeployType {
	case types.Dev, types.Compose:
		(*cAddr).Ch <- []string{"close"}
	}
	manager.CollectorAddrSlice = manager.CollectorAddrSlice[:lastCollectorIdx]
	return nil
}
