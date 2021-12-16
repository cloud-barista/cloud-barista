package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/modules/procedure/push/collector"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

type MetricCollector struct {
	CreateOrder       int
	ConsumerKafkaConn *kafka.Consumer
	Aggregator        collector.Aggregator
}

var KafkaConfig *kafka.ConfigMap

func PrintPanicError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func DeleteDeployment(clientSet *kubernetes.Clientset, createOrder int, collectorUUID string, namespace string) {
	fmt.Println("Deleting deployment...")
	deploymentName := fmt.Sprintf("%s%d-%s", types.DeploymentName, createOrder, collectorUUID)
	deploymentsClient := clientSet.AppsV1().Deployments(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), deploymentName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		fmt.Println("Fail to delete deployment.")
		fmt.Println(err)
	}
}

func main() {
	/** Get Env Val Start */
	kafkaEndpointUrl := os.Getenv("kafka_endpoint_url")
	var createOrder int
	createOrderString := os.Getenv("create_order")
	if createOrderString == "" {
		fmt.Println("Get Env Error")
		return
	}
	createOrder, _ = strconv.Atoi(createOrderString)
	aggregateType := types.AVG
	namespace := os.Getenv("namespace")
	dfAddr := os.Getenv("df_addr")
	collectInterval, _ := strconv.Atoi(os.Getenv("collect_interval"))
	collectorUUID := os.Getenv("collect_uuid")
	if kafkaEndpointUrl == "" || namespace == "" || dfAddr == "" {
		fmt.Println("Get Env Error")
		return
	}
	/** Get Env Val End */

	/** Set Kafka, ConfigMap Conn Start */
	KafkaConfig = &kafka.ConfigMap{
		//"bootstrap.servers":  fmt.Sprintf("%s", config.GetDefaultConfig().GetKafkaConfig().EndpointUrl),
		"bootstrap.servers":  kafkaEndpointUrl,
		"group.id":           fmt.Sprintf("%d", createOrder),
		"enable.auto.commit": true,
		//"session.timeout.ms": 15000,
		"auto.offset.reset": "earliest",
	}
	consumerKafkaConn, err := kafka.NewConsumer(KafkaConfig)
	PrintPanicError(err)
	config, errK8s := rest.InClusterConfig()
	PrintPanicError(errK8s)
	clientSet, errK8s2 := kubernetes.NewForConfig(config)
	PrintPanicError(errK8s2)
	/** Set Kafka, ConfigMap Conn End */

	/** Operate Collector Start */
	mc := MetricCollector{
		ConsumerKafkaConn: consumerKafkaConn,
		CreateOrder:       createOrder,
		Aggregator: collector.Aggregator{
			AggregateType: aggregateType,
		},
	}
	fmt.Println(fmt.Sprintf("#### Group_%d collector Create ####", createOrder))
	deadOrAliveCnt := map[string]int{}

	configMapFailCnt := 0
	for {
		time.Sleep(time.Duration(collectInterval) * time.Second)
		fmt.Println(fmt.Sprintf("#### Group_%d collector ####", createOrder))
		fmt.Println("Get ConfigMap")
		/** Get ConfigMap<Data: Collector UUID Map, BinaryData: Collector Topics> Start */
		configMap, err := clientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), "cb-dragonfly-collector-configmap", metav1.GetOptions{})
		if err != nil {
			if configMapFailCnt == 5 {
				DeleteDeployment(clientSet, createOrder, collectorUUID, namespace)
			}
			configMapFailCnt += 1
			fmt.Println("Fail to Get ConfigMap")
			fmt.Println(err)
			continue
		}
		/** Get ConfigMap<Data: Collector UUID Map, BinaryData: Collector Topics> End */

		/** Check My Collector UUID Start */
		_, alive := configMap.Data[collectorUUID]
		if !alive {
			DeleteDeployment(clientSet, createOrder, collectorUUID, namespace)
		}
		/** Check My Collector UUID End */

		/** Get My Allocated Topics Start */
		topicMap := map[int][]string{}
		if err := json.Unmarshal(configMap.BinaryData["topicMap"], &topicMap); err != nil {
			fmt.Println("Fail to unMarshal ConfigMap Object Data")
		}
		var DeliveredTopicList []string
		DeliveredTopicList, ok := topicMap[mc.CreateOrder]
		if !ok {
			fmt.Println("No topic on this Collector")
			continue
		}
		fmt.Println(fmt.Sprintf("Group_%d collector Delivered : %s", mc.CreateOrder, DeliveredTopicList))
		err = mc.ConsumerKafkaConn.SubscribeTopics(DeliveredTopicList, nil)
		if err != nil {
			fmt.Println(err)
		}
		/** Get My Allocated Topics End */

		/** Processing Topics to TSDB & Transmit Dead Topics To DF Start */
		start := time.Now()
		aliveTopics, _ := mc.Aggregator.AggregateMetric(mc.ConsumerKafkaConn, DeliveredTopicList)
		elapsed := time.Since(start)
		sort.Strings(aliveTopics)
		fmt.Println("Aggregate Time: ", elapsed)
		for _, aliveTopic := range aliveTopics {
			if _, ok := deadOrAliveCnt[aliveTopic]; ok {
				delete(deadOrAliveCnt, aliveTopic)
			}
		}
		if !cmp.Equal(DeliveredTopicList, aliveTopics) {
			_ = mc.ConsumerKafkaConn.Unsubscribe()
			deadTopics := util.ReturnDiffTopicList(DeliveredTopicList, aliveTopics)
			var err error
			for _, delTopic := range deadTopics {
				if _, ok := deadOrAliveCnt[delTopic]; !ok {
					deadOrAliveCnt[delTopic] = 0
				} else if ok {
					if deadOrAliveCnt[delTopic] == 2 {
						getUrl := fmt.Sprintf("http://%s/dragonfly/topic/delete/%s", dfAddr, delTopic)
						_, err = http.Get(getUrl)
						if err != nil {
							fmt.Println(err)
						}
						delete(deadOrAliveCnt, delTopic)
					}
					deadOrAliveCnt[delTopic] += 1
				}
			}
			if err != nil {
				fmt.Println("Sending Delete Topics to DF is Success")
			}
		}
		/** Processing Topics to TSDB & Transmit Dead Topics To DF End */
	}
	/** Operate Collector End */
}
