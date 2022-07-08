package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cloud-barista/cb-mcks/src/grpc-api/logger"
	lb_api "github.com/cloud-barista/cb-mcks/src/grpc-api/request"
	"github.com/cloud-barista/cb-spider/interface/api"
	sp_api "github.com/cloud-barista/cb-spider/interface/api"
	tb_api "github.com/cloud-barista/cb-tumblebug/src/api/grpc/request"
	"github.com/cloud-barista/cb-tumblebug/src/core/common"
)

//  Data Structure to contain user configuration
type ConnectionInfoTestData struct {
	// driver
	DriverName        string
	DriverLibFileName string

	// credential
	CredentialName string
	KeyValueList   []api.KeyValue

	// region, zone
	RegionName string
	Region     string
	Zone       string

	// connection config
	ConnectionConfigName string
}

// 테스트 환경에 맞게 파라미터 수정 필요
//----------------------------------------------------------- Set Test Data
// 0: Variable Environments for AWS
// 1: Variable Environments for GCP

var targetCSP = []string{
	"AWS",
	"GCP",
}

var connConfigNameList = []string{
	"cb-aws-config", // for AWS
	"cb-gcp-config", // for GCP
}

// AWS credential info
var AWSCredentialList = []api.KeyValue{
	api.KeyValue{Key: "ClientId", Value: "xxxxxxxxxx"},
	api.KeyValue{Key: "ClientSecret", Value: "xxxxxxxxxxxxxxxxxxxxxx"},
}

// GCP credential info
var GCPCredentialList = []api.KeyValue{
	api.KeyValue{Key: "PrivateKey", Value: "-----BEGIN PRIVATE KEY-----\nxxxxxxxxxxx\n-----END PRIVATE KEY-----"},
	api.KeyValue{Key: "ProjectID", Value: "xxxxxxxxxx"},
	api.KeyValue{Key: "ClientEmail", Value: "xxxxxxxxxx@gmail.com"},
}

var connInfoTestData = []ConnectionInfoTestData{
	{ // for AWS
		"aws-driver-01", "aws-driver-v1.0.so", // driver
		"aws-credential-01", AWSCredentialList, // credential
		"aws-(ohio)us-east-2", "us-east-2", "us-east-2a", // region, zone
		connConfigNameList[0], // connection config
	},

	{ // for GCP
		"gcp-driver-01", "gcp-driver-v1.0.so", // driver
		"gcp-credential-01", GCPCredentialList, // credential
		"gcp-us-central1-us-central1-a", "us-central1", "us-central1-a", // region, zone
		connConfigNameList[1], // connection config
	},
}

var namespace = "cb-aws-namespace" // aws selection
var cspIdx = 0                     // aws selection

var cluster = "cb-cluster"
var spec = "t2.medium" // aws selection
var workerNodeCount = 1

func main() {

	// API 호출 테스트
	SimpleLBApiTest()
	ConfigLBApiTest()

	// CB-SPIDER, CB-TUMBLEBUG 초기화
	CIM_Create_Info_Test()
	fmt.Print("\n\n============= 3 seconds waiting.. =============\n")
	time.Sleep(3 * time.Second)

	CreateNSApiTest()
	fmt.Print("\n\n============= 3 seconds waiting.. =============\n")
	time.Sleep(3 * time.Second)

	// Cluster 테스트
	CreateClusterApiTest()
	fmt.Print("\n\n============= 60 seconds waiting.. =============\n")
	time.Sleep(60 * time.Second)

	ListClusterApiTest()
	fmt.Print("\n\n============= 3 seconds waiting.. =============\n")
	time.Sleep(3 * time.Second)

	GetClusterApiTest()
	fmt.Print("\n\n============= 3 seconds waiting.. =============\n")
	time.Sleep(3 * time.Second)

	// Node 테스트
	AddNodeApiTest()
	fmt.Print("\n\n============= 60 seconds waiting.. =============\n")
	time.Sleep(60 * time.Second)

	ListNodeApiTest()
	fmt.Print("\n\n============= 3 seconds waiting.. =============\n")
	time.Sleep(3 * time.Second)

	GetNodeApiTest()
	fmt.Print("\n\n============= 3 seconds waiting.. =============\n")
	time.Sleep(3 * time.Second)

	RemoveNodeApiTest()
	fmt.Print("\n\n============= 60 seconds waiting.. =============\n")
	time.Sleep(60 * time.Second)

	DeleteClusterApiTest()
	fmt.Print("\n\n============= 60 seconds waiting.. =============\n")
	time.Sleep(60 * time.Second)
}

// SimpleLBApiTest - 환경설정함수를 이용한 간단한 MCKS API 호출
func SimpleLBApiTest() {

	fmt.Print("\n\n============= SimpleLBApiTest() =============\n")

	logger := logger.NewLogger()

	mcar := lb_api.NewMCARManager()

	err := mcar.SetServerAddr("localhost:50254")
	if err != nil {
		logger.Fatal(err)
	}

	err = mcar.SetTimeout(90 * time.Second)
	if err != nil {
		logger.Fatal(err)
	}

	/* 서버가 TLS 가 설정된 경우
	err = mcar.SetTLSCA(os.Getenv("APP_ROOT") + "/certs/ca.crt")
	if err != nil {
		logger.Fatal(err)
	}
	*/

	/* 서버가 JWT 인증이 설정된 경우
	err = mcar.SetJWTToken("xxxxxxxxxxxxxxxxxxx")
	if err != nil {
		logger.Fatal(err)
	}
	*/

	err = mcar.Open()
	if err != nil {
		logger.Fatal(err)
	}

	result, err := mcar.Healthy()
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Printf("\nresult :\n%s\n", result)

	mcar.Close()
}

// ConfigLBApiTest - 환경설정파일을 이용한 MCKS API 호출
func ConfigLBApiTest() {

	fmt.Print("\n\n============= ConfigLBApiTest() =============\n")

	logger := logger.NewLogger()

	mcar := lb_api.NewMCARManager()

	err := mcar.SetConfigPath("./grpc_conf.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	err = mcar.Open()
	if err != nil {
		logger.Fatal(err)
	}

	result, err := mcar.Healthy()
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Printf("\nresult :\n%s\n", result)

	mcar.Close()
}

/******************************************
  1. Create CloudInfoManager
  2. Setup env. with a config file
  3. Open New Session
  4. Close (with defer)
  5. Call API
    (1) create driver info
    (2) create credential info
    (3) create region info
    (4) create Connection config
******************************************/
func CIM_Create_Info_Test() {

	fmt.Print("\n\n============= CloudInfoManager: Create Driver/Credential/Region/ConnConfig Info Test =============\n")

	// 1. Create CloudInfoManager
	cim := sp_api.NewCloudInfoManager()

	// 2. Setup env. with a config file
	err := cim.SetConfigPath("./grpc_conf.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// 3. Open New Session
	err = cim.Open()
	if err != nil {
		log.Fatal(err)
	}
	// 4. Close (with defer)
	defer cim.Close()

	// 5. Call API

	// (1) create driver info
	fmt.Print("\n\n\t============= CloudInfoManager: (1) create driver info Test =============\n")
	reqCloudDriver := &api.CloudDriverReq{
		DriverName:        connInfoTestData[cspIdx].DriverName,
		ProviderName:      targetCSP[cspIdx],
		DriverLibFileName: connInfoTestData[cspIdx].DriverLibFileName,
	}
	result, err := cim.CreateCloudDriverByParam(reqCloudDriver)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nresult :\n%s\n", result)

	// (2) create credential info
	fmt.Print("\n\n\t============= CloudInfoManager: (2) create credential info Test =============\n")
	reqCredential := &api.CredentialReq{
		CredentialName:   connInfoTestData[cspIdx].CredentialName,
		ProviderName:     targetCSP[cspIdx],
		KeyValueInfoList: connInfoTestData[cspIdx].KeyValueList,
	}
	result, err = cim.CreateCredentialByParam(reqCredential)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nresult :\n%s\n", result)

	// (3) create region info
	fmt.Print("\n\n\t============= CloudInfoManager: (3) create region info Test =============\n")
	reqRegion := &api.RegionReq{
		RegionName:   connInfoTestData[cspIdx].RegionName,
		ProviderName: targetCSP[cspIdx],
		KeyValueInfoList: []api.KeyValue{
			api.KeyValue{Key: "Region", Value: connInfoTestData[cspIdx].Region},
			api.KeyValue{Key: "Zone", Value: connInfoTestData[cspIdx].Zone},
		},
	}
	result, err = cim.CreateRegionByParam(reqRegion)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nresult :\n%s\n", result)

	// (4) create Connection config
	fmt.Print("\n\n\t============= CloudInfoManager: (4) create connection config info Test =============\n")
	reqConnectionConfig := &api.ConnectionConfigReq{
		ConfigName:     connInfoTestData[cspIdx].ConnectionConfigName,
		ProviderName:   targetCSP[cspIdx],
		DriverName:     connInfoTestData[cspIdx].DriverName,
		CredentialName: connInfoTestData[cspIdx].CredentialName,
		RegionName:     connInfoTestData[cspIdx].RegionName,
	}
	result, err = cim.CreateConnectionConfigByParam(reqConnectionConfig)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nresult :\n%s\n", result)
}

// CreateNSApiTest - CB-TUMBLEBUG 초기화
func CreateNSApiTest() {

	fmt.Print("\n\n============= CreateNSApiTest() =============\n")

	logger := logger.NewLogger()

	ns := tb_api.NewNSManager()

	err := ns.SetConfigPath("./grpc_conf.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	err = ns.Open()
	if err != nil {
		logger.Fatal(err)
	}

	reqNs := &common.NsReq{
		Name:        namespace,
		Description: "NameSpace for General Testing",
	}
	result, err := ns.CreateNSByParam(reqNs)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Printf("\nresult :\n%s\n", result)

	ns.Close()
}

// CreateClusterApiTest - Cluster 생성
func CreateClusterApiTest() {
	fmt.Print("\n\n============= CreateClusterApiTest() =============\n")

	logger := logger.NewLogger()

	mcar := lb_api.NewMCARManager()

	err := mcar.SetConfigPath("./grpc_conf.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	err = mcar.Open()
	if err != nil {
		logger.Fatal(err)
	}

	reqCluster := &lb_api.ClusterCreateRequest{
		Namespace: namespace,
		Item: lb_api.ClusterReq{
			Name: cluster,
			ControlPlane: []lb_api.NodeConfig{
				lb_api.NodeConfig{Connection: connConfigNameList[cspIdx], Count: 1, Spec: spec},
			},
			Worker: []lb_api.NodeConfig{
				lb_api.NodeConfig{Connection: connConfigNameList[cspIdx], Count: workerNodeCount, Spec: spec},
			},
			Config: lb_api.Config{
				Kubernetes: lb_api.Kubernetes{
					NetworkCni:       "canal",
					PodCidr:          "10.244.0.0/16",
					ServiceCidr:      "10.96.0.0/12",
					ServiceDnsDomain: "cluster.local",
				},
			},
		},
	}

	result, err := mcar.CreateClusterByParam(reqCluster)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Printf("\nresult :\n%s\n", result)

	mcar.Close()
}

// ListClusterApiTest - Cluster 목록
func ListClusterApiTest() {
	fmt.Print("\n\n============= ListClusterApiTest() =============\n")

	logger := logger.NewLogger()

	mcar := lb_api.NewMCARManager()

	err := mcar.SetConfigPath("./grpc_conf.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	err = mcar.Open()
	if err != nil {
		logger.Fatal(err)
	}

	result, err := mcar.ListClusterByParam(namespace)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Printf("\nresult :\n%s\n", result)

	mcar.Close()
}

// GetClusterApiTest - Cluster 조회
func GetClusterApiTest() {
	fmt.Print("\n\n============= GetClusterApiTest() =============\n")

	logger := logger.NewLogger()

	mcar := lb_api.NewMCARManager()

	err := mcar.SetConfigPath("./grpc_conf.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	err = mcar.Open()
	if err != nil {
		logger.Fatal(err)
	}

	result, err := mcar.GetClusterByParam(namespace, cluster)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Printf("\nresult :\n%s\n", result)

	mcar.Close()
}

// AddNodeApiTest - Node 추가
func AddNodeApiTest() {
	fmt.Print("\n\n============= AddNodeApiTest() =============\n")

	logger := logger.NewLogger()

	mcar := lb_api.NewMCARManager()

	err := mcar.SetConfigPath("./grpc_conf.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	err = mcar.Open()
	if err != nil {
		logger.Fatal(err)
	}

	reqNode := &lb_api.NodeCreateRequest{
		Namespace: namespace,
		Cluster:   cluster,
		Item: lb_api.NodeReq{
			Worker: []lb_api.NodeConfig{
				lb_api.NodeConfig{Connection: connConfigNameList[cspIdx], Count: workerNodeCount, Spec: spec},
			},
		},
	}

	result, err := mcar.AddNodeByParam(reqNode)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Printf("\nresult :\n%s\n", result)

	mcar.Close()
}

// ListNodeApiTest - Node 목록
func ListNodeApiTest() {
	fmt.Print("\n\n============= ListNodeApiTest() =============\n")

	logger := logger.NewLogger()

	mcar := lb_api.NewMCARManager()

	err := mcar.SetConfigPath("./grpc_conf.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	err = mcar.Open()
	if err != nil {
		logger.Fatal(err)
	}

	result, err := mcar.ListNodeByParam(namespace, cluster)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Printf("\nresult :\n%s\n", result)

	mcar.Close()
}

// GetNodeApiTest - Node 조회
func GetNodeApiTest() {
	fmt.Print("\n\n============= GetNodeApiTest() =============\n")

	logger := logger.NewLogger()

	mcar := lb_api.NewMCARManager()

	err := mcar.SetConfigPath("./grpc_conf.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	err = mcar.Open()
	if err != nil {
		logger.Fatal(err)
	}

	/*
		* node 이름은 동적으로 설정됨.. 확인한 후에 테스트 수행..
		*
		result, err := mcar.GetNodeByParam(namespace, cluster, "xxx")
		if err != nil {
			logger.Fatal(err)
		}

		fmt.Printf("\nresult :\n%s\n", result)
	*/

	mcar.Close()
}

// RemoveNodeApiTest - Node 삭제
func RemoveNodeApiTest() {
	fmt.Print("\n\n============= RemoveNodeApiTest() =============\n")

	logger := logger.NewLogger()

	mcar := lb_api.NewMCARManager()

	err := mcar.SetConfigPath("./grpc_conf.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	err = mcar.Open()
	if err != nil {
		logger.Fatal(err)
	}

	/*
		* node 이름은 동적으로 설정됨.. 확인한 후에 테스트 수행..
		*
		result, err := mcar.RemoveNodeByParam(namespace, cluster, "xxx"")
		if err != nil {
			logger.Fatal(err)
		}

		fmt.Printf("\nresult :\n%s\n", result)
	*/

	mcar.Close()
}

// DeleteClusterApiTest - Cluster 삭제
func DeleteClusterApiTest() {
	fmt.Print("\n\n============= DeleteClusterApiTest() =============\n")

	logger := logger.NewLogger()

	mcar := lb_api.NewMCARManager()

	err := mcar.SetConfigPath("./grpc_conf.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	err = mcar.Open()
	if err != nil {
		logger.Fatal(err)
	}

	result, err := mcar.DeleteClusterByParam(namespace, cluster)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Printf("\nresult :\n%s\n", result)

	mcar.Close()
}
