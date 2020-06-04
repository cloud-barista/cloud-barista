package manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/cloud-barista/cb-spider/cloud-control-manager/vm-ssh"
	"github.com/google/uuid"
	"github.com/influxdata/influxdb1-client/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/client"
	"golang.org/x/crypto/ssh"

	"github.com/cloud-barista/cb-dragonfly/pkg/collector"
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/realtimestore"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
)

type APIServer struct {
	echo       *echo.Echo
	config     Config
	aggregator *collector.Aggregator
	Etcd       realtimestore.Storage
	InfluxDB   metricstore.Storage
	manager    *CollectManager
}

// API 서버 초기화
func NewAPIServer(config Config, aggregator *collector.Aggregator, influxDB metricstore.Storage, etcd realtimestore.Storage, manager *CollectManager) (*APIServer, error) {
	e := echo.New()
	apiServer := APIServer{
		echo:       e,
		config:     config,
		aggregator: aggregator,
		InfluxDB:   influxDB,
		Etcd:       etcd,
		manager:    manager,
	}
	return &apiServer, nil
}

// 모니터링 API 서버 실행
func (apiServer *APIServer) StartAPIServer(wg *sync.WaitGroup) error {
	defer wg.Done()
	logrus.Info("Start Monitoring API Server")

	// 모니터링 API 라우팅 룰 설정
	apiServer.SetRoutingRule(apiServer.echo)

	// 모니터링 API 서버 실행
	return apiServer.echo.Start(fmt.Sprintf(":%d", apiServer.config.APIServer.Port))
}

func (apiServer *APIServer) SetRoutingRule(e *echo.Echo) {

	e.Use(middleware.CORS())

	// 멀티 클라우드 인프라 서비스 모니터링/실시간 모니터링 정보 조회
	e.GET("/dragonfly/mcis/:mcis_id/info", apiServer.GetMCISMonInfo)
	e.GET("/dragonfly/mcis/:mcis_id/rt-info", apiServer.GetMCISRealtimeMonInfo)

	// 멀티 클라우드 인프라 VM 모니터링/실시간 모니터링 정보 조회
	e.GET("/dragonfly/mcis/:mcis_id/vm/:vm_id/metric/:metric_name/info", apiServer.GetVMMonInfo)
	e.GET("/dragonfly/mcis/:mcis_id/vm/:vm_id/metric/:metric_name/rt-info", apiServer.GetVMRealtimeMonInfo)

	// 멀티 클라우드 모니터링 정책 설정
	e.PUT("/dragonfly/config", apiServer.SetMonConfig)
	e.GET("/dragonfly/config", apiServer.GetMonConfig)
	e.PUT("/dragonfly/config/reset", apiServer.ResetMonConfig)

	// 에이전트 설치 스크립트 다운로드
	e.GET("/dragonfly/file/agent/install", apiServer.GetTelegrafInstallScript)

	// 에이전트 config, package 파일 다운로드
	e.GET("/dragonfly/file/agent/conf", apiServer.GetTelegrafConfFile)
	e.GET("/dragonfly/file/agent/pkg", apiServer.GetTelegrafPkgFile)

	// 에이전트 설치
	e.POST("/dragonfly/agent/install", apiServer.InstallTelegraf)
}

// 멀티 클라우드 인프라 서비스 모니터링 정보 조회
func (apiServer *APIServer) GetMCISMonInfo(c echo.Context) error {
	mcisId := c.Param("mcis_id")
	fmt.Println(mcisId)

	// TODO: MCIS 서비스 모니터링 정보 조회 기능 개발

	return c.JSON(http.StatusOK, "")
}

// 멀티 클라우드 인프라 서비스(MCIS) 실시간 모니터링 정보 조회
func (apiServer *APIServer) GetMCISRealtimeMonInfo(c echo.Context) error {
	mcisId := c.Param("mcis_id")
	fmt.Println(mcisId)

	// TODO: MCIS 서비스 실시간 모니터링 정보 조회 기능 개발

	return c.JSON(http.StatusOK, "")
}

// 멀티 클라우드 인프라 VM 모니터링 정보 조회
func (apiServer *APIServer) GetVMMonInfo(c echo.Context) error {

	// Path 파라미터 가져오기
	//mcisId := c.Param("mcis_id")
	vmId := c.Param("vm_id")
	metricName := c.Param("metric_name")

	// Query 파라미터 가져오기
	period := c.QueryParam("periodType")
	aggregateType := c.QueryParam("statisticsCriteria")
	duration := c.QueryParam("duration")

	var metricKey string

	switch metricName {
	case "cpu":

		// cpu 메트릭 조회
		metricKey = "cpu"
		cpuMetric, err := apiServer.InfluxDB.ReadMetric(vmId, metricKey, period, aggregateType, duration)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		resultMetric, err := metricstore.MappingMonMetric(metricKey, &cpuMetric)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, resultMetric)

	case "cpufreq":
		metricKey := "cpufreq"
		cfMetric, err := apiServer.InfluxDB.ReadMetric(vmId, metricKey, period, aggregateType, duration)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		if cfMetric == nil {
			return c.JSON(http.StatusNotFound, err)
		}
		resultMetric, err := metricstore.MappingMonMetric(metricKey, &cfMetric)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, resultMetric)

	case "memory":

		// memory 메트릭 조회
		metricKey := "mem"
		memMetric, err := apiServer.InfluxDB.ReadMetric(vmId, metricKey, period, aggregateType, duration)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		resultMetric, err := metricstore.MappingMonMetric(metricKey, &memMetric)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, resultMetric)

	case "disk":

		// disk, diskio 메트릭 조회
		metricKey := "disk"
		diskMetric, err := apiServer.InfluxDB.ReadMetric(vmId, metricKey, period, aggregateType, duration)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		metricKey = "diskio"
		diskIoMetric, err := apiServer.InfluxDB.ReadMetric(vmId, metricKey, period, aggregateType, duration)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		if diskMetric == nil || diskIoMetric == nil {
			return c.JSON(http.StatusNotFound, err)
		}

		diskRow := diskMetric.(models.Row)
		diskIoRow := diskIoMetric.(models.Row)

		// Aggregate Column Info
		//diskRow.Columns = append(diskRow.Columns, diskIoRow.Columns[1:]...)

		// Aggregate Monitoring data
		/*initArr := make([]interface{}, 4)
		for idx := range diskRow.Values {
			if idx <= len(diskIoRow.Values) {
				diskRow.Values[idx] = append(diskRow.Values[idx], diskIoRow.Values[idx][1:]...)
				continue
			}
			diskRow.Values[idx] = append(diskRow.Values[idx], initArr...)
		}*/

		// Aggregate Metric
		var resultRow models.Row
		resultRow.Name = "disk"
		resultRow.Tags = diskRow.Tags
		resultRow.Columns = append(resultRow.Columns, diskRow.Columns[0:]...)
		resultRow.Columns = append(resultRow.Columns, diskIoRow.Columns[1:]...)

		// TimePoint 맵 생성 (disk, diskio 메트릭)
		timePointMap := make(map[string]string, len(diskRow.Values))
		for _, val := range diskRow.Values {
			timePoint := val[0].(string)
			timePointMap[timePoint] = timePoint
		}
		for _, val := range diskIoRow.Values {
			timePoint := val[0].(string)
			if tp, exist := timePointMap[timePoint]; !exist {
				timePointMap[tp] = tp
			}
		}

		// TimePoint 배열 생성
		idx := 0
		timePointArr := make([]string, len(timePointMap))
		for _, tp := range timePointMap {
			timePointArr[idx] = tp
			idx++
		}
		sort.Strings(timePointArr)

		// TimePoint 배열 기준 모니터링 메트릭 Aggregate
		for _, tp := range timePointArr {

			metricVal := make([]interface{}, 1)
			metricVal[0] = tp

			// disk 메트릭 aggregate
			diskMetricAdded := false
			for idx, val := range diskRow.Values {
				t := val[0].(string)
				if strings.EqualFold(t, tp) {
					metricVal = append(metricVal, val[1:]...)
					diskMetricAdded = true
					break
				}
				// 해당 TimePoint에 해당하는 disk 메트릭이 없을 경우 0으로 값 초기화
				if !diskMetricAdded && (idx == len(diskRow.Values)-1) {
					initVal := make([]interface{}, len(val)-1)
					for i := range initVal {
						initVal[i] = 0
					}
					metricVal = append(metricVal, initVal...)
				}
			}

			// diskio 메트릭 aggregate
			diskIoMetricAdded := false
			for idx, val := range diskIoRow.Values {
				t := val[0].(string)
				if strings.EqualFold(t, tp) {
					metricVal = append(metricVal, val[1:]...)
					diskIoMetricAdded = true
					break
				}
				// 해당 TimePoint에 해당하는 disk 메트릭이 없을 경우 0으로 값 초기화
				if !diskIoMetricAdded && (idx == len(diskIoRow.Values)-1) {
					initVal := make([]interface{}, len(val)-1)
					for i := range initVal {
						initVal[i] = 0
					}
					metricVal = append(metricVal, initVal...)
				}
			}

			resultRow.Values = append(resultRow.Values, metricVal)
		}

		return c.JSON(http.StatusOK, resultRow)

	case "network":

		// network 메트릭 조회
		metricKey := "net"
		netMetric, err := apiServer.InfluxDB.ReadMetric(vmId, metricKey, period, aggregateType, duration)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		if netMetric == nil {
			return c.JSON(http.StatusNotFound, err)
		}

		resultMetric, err := metricstore.MappingMonMetric(metricKey, &netMetric)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, resultMetric)

	default:
		return c.JSON(http.StatusNotFound, fmt.Sprintf("not found metric : %s", metricName))
	}
}

// 멀티 클라우드 인프라 VM 실시간 모니터링 정보 조회
func (apiServer *APIServer) GetVMRealtimeMonInfo(c echo.Context) error {
	// Path 파라미터 가져오기
	//mcisId := c.Param("mcis_id")
	vmId := c.Param("vm_id")
	metricName := c.Param("metric_name")

	// Query 파라미터 가져오기
	aggregateType := c.QueryParam("statisticsCriteria")

	apiServer.aggregator.Etcd = apiServer.Etcd

	resultMap := map[string]interface{}{}
	resultMap["vmId"] = vmId
	resultMap["metricName"] = metricName
	resultMap["time"] = time.Now().UTC()
	resultMap["value"] = map[string]interface{}{}

	var metricKey string
	var metricMap map[string]interface{}
	var diskMetric, diskIoMetric, result map[string]interface{}
	var diskMetricMap, diskIoMetricMap map[string]interface{}
	var err error
	var val interface{}

	if metricName == "disk" || metricName == "diskio" {
		metricKey = "disk"
		diskMetric, err = apiServer.aggregator.GetAggregateDiskMetric(vmId, metricKey, aggregateType)
		if err != nil {
			// 만약 실시간 데이터가 없을 경우 empty Map 값 전달
			if err.(client.Error).Code == 100 {
				return c.JSON(http.StatusOK, resultMap)
			}
			return c.JSON(http.StatusInternalServerError, err)
		}
		// disk 메트릭 매핑
		diskMetricMap, err = realtimestore.MappingMonMetric(metricKey, diskMetric)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		// diskio 메트릭 조회
		metricKey = "diskio"
		diskIoMetric, err = apiServer.aggregator.GetAggregateDiskMetric(vmId, metricKey, aggregateType)
		if err != nil {
			// 만약 실시간 데이터가 없을 경우 empty Map 값 전달
			if err.(client.Error).Code == 100 {
				return c.JSON(http.StatusOK, resultMap)
			}
		}
		// diskio 메트릭 매핑
		diskIoMetricMap, err = realtimestore.MappingMonMetric(metricKey, diskIoMetric)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	} else {
		//메트릭 키 설정
		switch metricName {
		case "cpu":
			metricKey = "cpu"
		case "memory":
			metricKey = "mem"
		case "network":
			metricKey = "net"
		case "cpufreq":
			metricKey = "cpufreq"
		default:
			return c.JSON(http.StatusNotFound, fmt.Sprintf("not found metric : %s", metricName))
		}
		// disk, diskio 제외한 메트릭 조회
		result, err = apiServer.aggregator.GetAggregateMetric(vmId, metricKey, aggregateType)
		if err != nil {
			// 만약 실시간 데이터가 없을 경우 empty Map 값 전달
			if err.(client.Error).Code == 100 {
				return c.JSON(http.StatusOK, resultMap)
			}
		}
		// disk, diskio 제외한 메트릭 매핑
		metricMap, err = realtimestore.MappingMonMetric(metricKey, result)
		if _, ok := err.(client.Error); ok {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	for metricKey, val = range metricMap {
		metricMap[metricKey] = val
	}
	for metricKey, val = range diskMetricMap {
		metricMap[metricKey] = val
	}
	for metricKey, val = range diskIoMetricMap {
		metricMap[metricKey] = val
	}
	resultMap["value"] = metricMap
	return c.JSON(http.StatusOK, resultMap)
}

// 모니터링 정책 설정
func (apiServer *APIServer) SetMonConfig(c echo.Context) error {

	// form 파라미터 정보 가져오기
	agentInterval, err := strconv.Atoi(c.FormValue("agent_interval"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	collectorInterval, err := strconv.Atoi(c.FormValue("collector_interval"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	schedulingInterval, err := strconv.Atoi(c.FormValue("schedule_interval"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	maxHostCnt, err := strconv.Atoi(c.FormValue("max_host_count"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	agentTtl, err := strconv.Atoi(c.FormValue("agent_TTL"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	manager := apiServer.manager
	manager.Config.Monitoring.AgentInterval = agentInterval
	manager.Config.Monitoring.CollectorInterval = collectorInterval
	manager.Config.Monitoring.ScheduleInterval = schedulingInterval
	manager.Config.Monitoring.MaxHostCount = maxHostCnt
	manager.Config.Monitoring.AgentTtl = agentTtl

	monConfig := MonConfig{
		AgentTtl:           agentTtl,
		AgentInterval:      agentInterval,
		CollectorInterval:  collectorInterval,
		SchedulingInterval: schedulingInterval,
		MaxHostCount:       maxHostCnt,
	}

	// TODO: 구조체 map[string]interface{} 타입으로 Unmarshal
	// TODO: 추후에 별도의 map 변환 함수 (toMap() 개발)
	reqBodyBytes := new(bytes.Buffer)
	if err = json.NewEncoder(reqBodyBytes).Encode(monConfig); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	byteData := reqBodyBytes.Bytes()

	jsonMap := map[string]interface{}{}
	if err = json.Unmarshal(byteData, &jsonMap); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// etcd 저장소에 모니터링 정책 저장 후 결과 값 반환
	err = apiServer.Etcd.WriteMetric("/mon/config", jsonMap)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, manager.Config.Monitoring)
}

// 모니터링 정책 조회
func (apiServer *APIServer) GetMonConfig(c echo.Context) error {

	// etcd 저장소에서 모니터링 정책 정보 조회
	node, err := apiServer.Etcd.ReadMetric("/mon/config")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// 모니터링 정책 정보 구조체 매핑
	monConfig := &MonConfig{}
	if err := json.Unmarshal([]byte(node.Value), monConfig); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, monConfig)
}

// 모니터링 정책 초기화
func (apiServer *APIServer) ResetMonConfig(c echo.Context) error {

	// config 파일 기준 모니터링 정책 초기화
	monConfig := MonConfig{
		AgentTtl:           apiServer.config.Monitoring.AgentTtl,
		AgentInterval:      apiServer.config.Monitoring.AgentInterval,
		CollectorInterval:  apiServer.config.Monitoring.CollectorInterval,
		SchedulingInterval: apiServer.config.Monitoring.ScheduleInterval,
		MaxHostCount:       apiServer.config.Monitoring.MaxHostCount,
	}

	// TODO: 구조체 map[string]interface{} 타입으로 Unmarshal
	// TODO: 추후에 별도의 map 변환 함수 (toMap() 개발)
	reqBodyBytes := new(bytes.Buffer)
	if err := json.NewEncoder(reqBodyBytes).Encode(monConfig); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	byteData := reqBodyBytes.Bytes()

	jsonMap := map[string]interface{}{}
	if err := json.Unmarshal(byteData, &jsonMap); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// etcd 저장소에 모니터링 정책 저장 후 결과 값 반환
	if err := apiServer.Etcd.WriteMetric("/mon/config", jsonMap); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, monConfig)
}

// Telegraf agent 설치 스크립트 파일 다운로드
func (apiServer *APIServer) GetTelegrafInstallScript(c echo.Context) error {
	// Query 파라미터 가져오기
	mcisId := c.QueryParam("mcis_id")
	vmId := c.QueryParam("vm_id")

	// Query 파라미터 값 체크
	if mcisId == "" || vmId == "" {
		err := errors.New("failed to get package. query parameter is missing")
		return c.JSON(http.StatusInternalServerError, err)
	}

	collectorServer := fmt.Sprintf("%s:%d", apiServer.config.CollectManager.CollectorIP, apiServer.config.APIServer.Port)

	rootPath := os.Getenv("CBMON_PATH")
	filePath := rootPath + "/file/install_agent.sh"

	read, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// 파일 내의 변수 값 설정 (vmId, collectorServer)
	strConf := string(read)
	strConf = strings.ReplaceAll(strConf, "{{mcis_id}}", mcisId)
	strConf = strings.ReplaceAll(strConf, "{{vm_id}}", vmId)
	strConf = strings.ReplaceAll(strConf, "{{api_server}}", collectorServer)

	return c.Blob(http.StatusOK, "text/plain", []byte(strConf))
}

// Telegraf config 파일 다운로드
func (apiServer *APIServer) GetTelegrafConfFile(c echo.Context) error {
	// Query 파라미터 가져오기
	mcisId := c.QueryParam("mcis_id")
	vmId := c.QueryParam("vm_id")

	// Query 파라미터 값 체크
	if mcisId == "" || vmId == "" {
		err := errors.New("failed to get package. query parameter is missing")
		return c.JSON(http.StatusInternalServerError, err)
	}

	collectorServer := fmt.Sprintf("udp://%s:%d", apiServer.manager.Config.CollectManager.CollectorIP, apiServer.manager.Config.CollectManager.CollectorPort)
	influxDBServer := fmt.Sprintf("http://%s:8086", apiServer.manager.Config.CollectManager.CollectorIP)

	rootPath := os.Getenv("CBMON_PATH")
	filePath := rootPath + "/file/conf/telegraf.conf"

	read, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// 파일 내의 변수 값 설정 (hostId, collectorServer)
	strConf := string(read)
	strConf = strings.ReplaceAll(strConf, "{{mcis_id}}", mcisId)
	strConf = strings.ReplaceAll(strConf, "{{vm_id}}", vmId)
	strConf = strings.ReplaceAll(strConf, "{{collector_server}}", collectorServer)
	strConf = strings.ReplaceAll(strConf, "{{influxdb_server}}", influxDBServer)

	return c.Blob(http.StatusOK, "text/plain", []byte(strConf))
}

// Telegraf package 파일 다운로드
func (apiServer *APIServer) GetTelegrafPkgFile(c echo.Context) error {
	// Query 파라미터 가져오기
	osType := c.QueryParam("osType")
	arch := c.QueryParam("arch")

	// Query 파라미터 값 체크
	if osType == "" || arch == "" {
		err := errors.New("failed to get package. query parameter is missing")
		return c.JSON(http.StatusInternalServerError, err)
	}

	// osType, architecture 지원 여부 체크
	osType = strings.ToLower(osType)
	if osType != "ubuntu" && osType != "centos" {
		err := errors.New("failed to get package. not supported OS type")
		return c.JSON(http.StatusInternalServerError, err)
	}
	if !strings.Contains(arch, "32") && !strings.Contains(arch, "64") {
		err := errors.New("failed to get package. not supported architecture")
		return c.JSON(http.StatusInternalServerError, err)
	}

	if strings.Contains(arch, "64") {
		arch = "x64"
	} else {
		arch = "x32"
	}

	rootPath := os.Getenv("CBMON_PATH")
	var filePath string
	switch osType {
	case "ubuntu":
		filePath = rootPath + fmt.Sprintf("/file/pkg/%s/%s/telegraf_1.15.0~c78045c1-0_amd64.deb", osType, arch)
	case "centos":
		filePath = rootPath + fmt.Sprintf("/file/pkg/%s/%s/telegraf-1.12.0~f09f2b5-0.x86_64.rpm", osType, arch)
	default:
		err := errors.New(fmt.Sprintf("failed to get package. osType %s not supported", osType))
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.File(filePath)
}

func (apiServer *APIServer) createTelegrafConfigFile(mcisId string, vmId string) (string, error) {

	collectorServer := fmt.Sprintf("udp://%s:%d", apiServer.config.CollectManager.CollectorIP, apiServer.config.CollectManager.CollectorPort)
	influxDBServer := fmt.Sprintf("http://%s:8086", apiServer.manager.Config.CollectManager.CollectorIP)

	rootPath := os.Getenv("CBMON_PATH")
	filePath := rootPath + "/file/conf/telegraf.conf"

	read, err := ioutil.ReadFile(filePath)
	if err != nil {
		// ERROR 정보 출럭
		logrus.Error("failed to read telegraf.conf file.")
		return "", err
	}

	// 파일 내의 변수 값 설정 (hostId, collectorServer)
	strConf := string(read)
	strConf = strings.ReplaceAll(strConf, "{{mcis_id}}", mcisId)
	strConf = strings.ReplaceAll(strConf, "{{vm_id}}", vmId)
	strConf = strings.ReplaceAll(strConf, "{{collector_server}}", collectorServer)
	strConf = strings.ReplaceAll(strConf, "{{influxdb_server}}", influxDBServer)

	// telegraf.conf 파일 생성
	telegrafFilePath := rootPath + "/file/conf/"
	createFileName := "telegraf-" + uuid.New().String() + ".conf"
	telegrafConfFile := telegrafFilePath + createFileName

	err = ioutil.WriteFile(telegrafConfFile, []byte(strConf), os.FileMode(777))
	if err != nil {
		logrus.Error("failed to create telegraf.conf file.")
		return "", err
	}

	return telegrafConfFile, err
}

func (apiServer *APIServer) InstallTelegraf(c echo.Context) error {
	// form 파라미터 값 가져오기
	mcisId := c.FormValue("mcis_id")
	vmId := c.FormValue("vm_id")
	publicIp := c.FormValue("public_ip")
	userName := c.FormValue("user_name")
	sshKey := c.FormValue("ssh_key")

	// form 파라미터 값 체크
	if mcisId == "" || vmId == "" || publicIp == "" || userName == "" || sshKey == "" {
		errMsg := setMessage("failed to get package. query parameter is missing")
		return c.JSON(http.StatusInternalServerError, errMsg)
	}

	sshInfo := sshrun.SSHInfo{
		ServerPort: publicIp + ":22", //serverEndPoint
		UserName:   userName,         //userName
		PrivateKey: []byte(sshKey),   //[]byte(privateKey)
	}

	// {사용자계정}/cb-dragonfly 폴더 생성
	createFolderCmd := fmt.Sprintf("mkdir $HOME/cb-dragonfly")
	if _, err := sshrun.SSHRun(sshInfo, createFolderCmd); err != nil {
		errMsg := setMessage(fmt.Sprintf("failed to make directory cb-dragonfly, error=%s", err))
		return c.JSON(http.StatusInternalServerError, errMsg)
	}

	// 리눅스 OS 환경 체크
	osType, err := sshrun.SSHRun(sshInfo, "hostnamectl | grep 'Operating System' | awk '{print $3}' | tr 'a-z' 'A-Z'")
	if err != nil {
		cleanTelegrafInstall(sshInfo, osType)
		errMsg := setMessage(fmt.Sprintf("failed to check linux OS environments, error=%s", err))
		return c.JSON(http.StatusInternalServerError, errMsg)
	}

	rootPath := os.Getenv("CBMON_PATH")

	var sourceFile, targetFile, installCmd string
	if strings.Contains(osType, "CENTOS") {
		sourceFile = rootPath + "/file/pkg/centos/x64/telegraf-1.12.0~f09f2b5-0.x86_64.rpm"
		targetFile = fmt.Sprintf("$HOME/cb-dragonfly/cb-agent.rpm")
		installCmd = fmt.Sprintf("sudo rpm -ivh $HOME/cb-dragonfly/cb-agent.rpm")
	} else if strings.Contains(osType, "UBUNTU") {
		sourceFile = rootPath + "/file/pkg/ubuntu/x64/telegraf_1.12.0~f09f2b5-0_amd64.deb"
		targetFile = fmt.Sprintf("$HOME/cb-dragonfly/cb-agent.deb")
		installCmd = fmt.Sprintf("sudo dpkg -i $HOME/cb-dragonfly/cb-agent.deb")
	}

	// 에이전트 설치 패키지 다운로드
	if err := sshCopyWithTimeout(sshInfo, sourceFile, targetFile); err != nil {
		cleanTelegrafInstall(sshInfo, osType)
		errMsg := setMessage(fmt.Sprintf("failed to download agent package, error=%s", err))
		return c.JSON(http.StatusInternalServerError, errMsg)
	}

	// 패키지 설치 실행
	if _, err := sshrun.SSHRun(sshInfo, installCmd); err != nil {
		cleanTelegrafInstall(sshInfo, osType)
		errMsg := setMessage(fmt.Sprintf("failed to install agent package, error=%s", err))
		return c.JSON(http.StatusInternalServerError, errMsg)
	}

	sshrun.SSHRun(sshInfo, "sudo rm /etc/telegraf/telegraf.conf")

	// telegraf_conf 파일 복사
	telegrafConfSourceFile, err := apiServer.createTelegrafConfigFile(mcisId, vmId)
	telegrafConfTargetFile := "$HOME/cb-dragonfly/telegraf.conf"
	if err != nil {
		cleanTelegrafInstall(sshInfo, osType)
		errMsg := setMessage(fmt.Sprintf("failed to create telegraf.conf, error=%s", err))
		return c.JSON(http.StatusInternalServerError, errMsg)
	}
	if err := sshrun.SSHCopy(sshInfo, telegrafConfSourceFile, telegrafConfTargetFile); err != nil {
		cleanTelegrafInstall(sshInfo, osType)
		errMsg := setMessage(fmt.Sprintf("failed to copy telegraf.conf, error=%s", err))
		return c.JSON(http.StatusInternalServerError, errMsg)
	}

	if _, err := sshrun.SSHRun(sshInfo, "sudo mv $HOME/cb-dragonfly/telegraf.conf /etc/telegraf/"); err != nil {
		cleanTelegrafInstall(sshInfo, osType)
		errMsg := setMessage(fmt.Sprintf("failed to move telegraf.conf, error=%s", err))
		return c.JSON(http.StatusInternalServerError, errMsg)
	}

	// 공통 서비스 활성화 및 실행
	if _, err := sshrun.SSHRun(sshInfo, "sudo systemctl enable telegraf && sudo systemctl restart telegraf"); err != nil {
		cleanTelegrafInstall(sshInfo, osType)
		errMsg := setMessage(fmt.Sprintf("failed to enable and start telegraf service, error=%s", err))
		return c.JSON(http.StatusInternalServerError, errMsg)
	}

	// telegraf UUId conf 파일 삭제
	err = os.Remove(telegrafConfSourceFile)
	if err != nil {
		cleanTelegrafInstall(sshInfo, osType)
		errMsg := setMessage(fmt.Sprintf("failed to remove temporary telegraf.conf file, error=%s", err))
		return c.JSON(http.StatusInternalServerError, errMsg)
	}

	// 에이전트 설치에 사용한 파일 폴더 채로 제거
	removeRpmCmd := fmt.Sprintf("sudo rm -rf $HOME/cb-dragonfly")
	if _, err := sshrun.SSHRun(sshInfo, removeRpmCmd); err != nil {
		cleanTelegrafInstall(sshInfo, osType)
		errMsg := setMessage(fmt.Sprintf("failed to remove cb-dragonfly directory, error=%s", err))
		return c.JSON(http.StatusInternalServerError, errMsg)
	}

	// 정상 설치 확인
	checkCmd := "telegraf --version"
	if result, err := util.RunCommand(publicIp, userName, sshKey, checkCmd); err != nil {
		cleanTelegrafInstall(sshInfo, osType)
		errMsg := setMessage(fmt.Sprintf("failed to run telegraf command, error=%s", err))
		return c.JSON(http.StatusInternalServerError, errMsg)
	} else {
		if strings.Contains(*result, "command not found") {
			cleanTelegrafInstall(sshInfo, osType)
			errMsg := setMessage(fmt.Sprintf("failed to run telegraf command, error=%s", err))
			return c.JSON(http.StatusInternalServerError, errMsg)
		}

		successMsg := setMessage("agent installation is finished")
		return c.JSON(http.StatusOK, successMsg)
	}
}

func setMessage(msg string) echo.Map {
	errResp := echo.Map{}
	errResp["message"] = msg
	return errResp
}

func cleanTelegrafInstall(sshInfo sshrun.SSHInfo, osType string) {

	// Uninstall Telegraf
	var uninstallCmd string
	if strings.Contains(osType, "CENTOS") {
		uninstallCmd = fmt.Sprintf("sudo rpm -e telegraf")
	} else if strings.Contains(osType, "UBUNTU") {
		uninstallCmd = fmt.Sprintf("sudo dpkg -r telegraf")
	}
	sshrun.SSHRun(sshInfo, uninstallCmd)

	// Delete Install Files
	removeRpmCmd := fmt.Sprintf("sudo rm -rf $HOME/cb-dragonfly")
	sshrun.SSHRun(sshInfo, removeRpmCmd)
	removeDirCmd := fmt.Sprintf("sudo rm -rf /etc/telegraf/cb-dragonfly")
	sshrun.SSHRun(sshInfo, removeDirCmd)
}

func sshCopyWithTimeout(sshInfo sshrun.SSHInfo, sourceFile string, targetFile string) error {
	signer, err := ssh.ParsePrivateKey(sshInfo.PrivateKey)
	if err != nil {
		return err
	}
	clientConfig := ssh.ClientConfig{
		User: sshInfo.UserName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client := scp.NewClientWithTimeout(sshInfo.ServerPort, &clientConfig, 600*time.Second)
	client.Connect()

	file, _ := os.Open(sourceFile)

	defer client.Close()
	defer file.Close()

	return client.CopyFile(file, targetFile, "0755")
}
