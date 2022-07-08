package mcis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/modules/monitoring/push/mcis/collector"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

/* CollectManager */
// 1. CollectorAddrSlice
//  - 생성한 Go-routine 기반 collector 의 주소값을 보관하는 배열 변수입니다.
// 2. CollectorPolicy
//  - Collector Policy 로 MaxAgentHost, CSP 중 한개의 입력값을 받습니다.
//  - MaxAgentHost 방식은 동작 방식(deployType) - dev, compose, helm 에서 모두 정상 동작 및 기능테스트까지 완료하였습니다.
//  - CSP 방식은 dev, compose 환경에서만 정상 동작 및 기능 테스트까지 완료하였습니다.
// 3. K8sClientSet
//  - 동작 방식(deployType)이 helm 일 경우, k8s와 통신하기 위한 conn 객체입니다.
//  - 해당 객체는 k8s in-cluster 모드에만 동작합니다.
type CollectManager struct {
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
	// k8s conn set Start
	inClusterK8sConfig, errK8s := rest.InClusterConfig()
	if errK8s != nil {
		err = errK8s
	}
	clientSet, errK8s2 := kubernetes.NewForConfig(inClusterK8sConfig)
	if errK8s2 != nil {
		err = errK8s2
	}
	manager.K8sClientSet = clientSet
	// k8s conn set End

	// helm 으로 배포할 경우, df 는 collector 를 deployment 형태로 배포합니다.
	// df 와 collector 는 configmap 으로 topic 정보를 동기화합니다.
	// 아래 코드는 configmap 을 설정 및 배포하는 코드입니다.
	configMapsClient := manager.K8sClientSet.CoreV1().ConfigMaps(config.GetInstance().Dragonfly.HelmNamespace)
	configMap := &apiv1.ConfigMap{Data: map[string]string{}, ObjectMeta: metav1.ObjectMeta{
		Name: types.ConfigMapName,
	}}
	// Deploy ConfigMap => (1) 드래곤 플라이가 배포한 컨피그맵이 이미 생성되어 있는지 조회 (2) 컨피그 맵이 없을 경우, 배포
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

	// collector 생성 요청이 들어왔을 경우
	manager.WaitGroup.Add(1)
	collectorCreateOrder := len(manager.CollectorAddrSlice)
	// 생성 순서 idx 값을 collector 객체에 넣고, 생성합니다.
	newCollector, err := collector.NewMetricCollector(
		types.AVG,
		collectorCreateOrder,
	)
	if err != nil {
		return err
	}
	// 생성한 Collector 의 주소 값을 CollectorAddrSlice 배열에 추가합니다.
	manager.CollectorAddrSlice = append(manager.CollectorAddrSlice, &newCollector)

	switch config.GetInstance().Monitoring.DeployType {
	// helm 배포일 경우, collector 를 deployment 로 배포합니다.
	case types.Helm:
		collectorUUID := fmt.Sprintf("%p", &newCollector)
		env := []apiv1.EnvVar{
			{Name: "kafka_endpoint_url", Value: config.GetInstance().Kafka.EndpointUrl},
			{Name: "create_order", Value: strconv.Itoa(collectorCreateOrder)},
			{Name: "namespace", Value: config.GetInstance().Dragonfly.HelmNamespace},
			{Name: "df_addr", Value: fmt.Sprintf("%s:%d", config.GetInstance().Dragonfly.DragonflyIP, config.GetInstance().Dragonfly.HelmPort)},
			{Name: "collect_interval", Value: strconv.Itoa(config.GetInstance().Monitoring.MCISCollectorInterval)},
			{Name: "collect_uuid", Value: collectorUUID},
		}
		deploymentTemplate := util.DeploymentTemplate(collectorCreateOrder, collectorUUID, env)
		fmt.Println("Creating deployment...")
		result, err := manager.K8sClientSet.AppsV1().Deployments(config.GetInstance().Dragonfly.HelmNamespace).Create(context.TODO(), deploymentTemplate, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		fmt.Println("Created deployment: ", result.GetObjectMeta().GetName())
		return nil
	// dev, compose 배포일 경우, collector 를 newCollector.Collector 메소드를 통해 Go-routine 으로 동작 시킵니다.
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

	// collector 삭제 요청이 들어왔을 경우, df 는 collector 객체를 CollectorAddrSlice 에서 삭제합니다.
	lastCollectorIdx := len(manager.CollectorAddrSlice) - 1
	cAddr := manager.CollectorAddrSlice[lastCollectorIdx]
	switch config.GetInstance().Monitoring.DeployType {
	// 동작 방식이 dev, compose 배포일 경우, 해당 collector 가 동작하는 Go-routine 에게 "close" 채널 값을 보내 종료시킵니다.
	case types.Dev, types.Compose:
		(*cAddr).Ch <- []string{"close"}
	}
	manager.CollectorAddrSlice = manager.CollectorAddrSlice[:lastCollectorIdx]
	return nil
}
