package mcir

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"

	//"strings"

	"github.com/cloud-barista/cb-spider/interface/api"
	"github.com/cloud-barista/cb-tumblebug/src/core/common"
	"github.com/go-resty/resty/v2"

	//"github.com/cloud-barista/cb-tumblebug/src/core/mcis"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xwb1989/sqlparser"
)

type SpiderSpecInfo struct { // Spider
	// https://github.com/cloud-barista/cb-spider/blob/master/cloud-control-manager/cloud-driver/interfaces/resources/VMSpecHandler.go

	Region string
	Name   string
	VCpu   SpiderVCpuInfo
	Mem    string
	Gpu    []SpiderGpuInfo

	KeyValueList []common.KeyValue
}

type SpiderVCpuInfo struct { // Spider
	Count string
	Clock string // GHz
}

type SpiderGpuInfo struct { // Spider
	Count string
	Mfr   string
	Model string
	Mem   string
}

type TbSpecReq struct { // Tumblebug
	Name           string `json:"name"`
	ConnectionName string `json:"connectionName"`
	CspSpecName    string `json:"cspSpecName"`
	Description    string `json:"description"`
}

type TbSpecInfo struct { // Tumblebug
	Id                    string  `json:"id"`
	Name                  string  `json:"name"`
	ConnectionName        string  `json:"connectionName"`
	CspSpecName           string  `json:"cspSpecName"`
	Os_type               string  `json:"os_type"`
	Num_vCPU              uint16  `json:"num_vCPU"`
	Num_core              uint16  `json:"num_core"`
	Mem_GiB               uint16  `json:"mem_GiB"`
	Storage_GiB           uint32  `json:"storage_GiB"`
	Description           string  `json:"description"`
	Cost_per_hour         float32 `json:"cost_per_hour"`
	Num_storage           uint8   `json:"num_storage"`
	Max_num_storage       uint8   `json:"max_num_storage"`
	Max_total_storage_TiB uint16  `json:"max_total_storage_TiB"`
	Net_bw_Gbps           uint16  `json:"net_bw_Gbps"`
	Ebs_bw_Mbps           uint32  `json:"ebs_bw_Mbps"`
	Gpu_model             string  `json:"gpu_model"`
	Num_gpu               uint8   `json:"num_gpu"`
	Gpumem_GiB            uint16  `json:"gpumem_GiB"`
	Gpu_p2p               string  `json:"gpu_p2p"`
	OrderInFilteredResult uint16  `json:"orderInFilteredResult"`
	EvaluationStatus      string  `json:"evaluationStatus"`
	EvaluationScore_01    float32 `json:"evaluationScore_01"`
	EvaluationScore_02    float32 `json:"evaluationScore_02"`
	EvaluationScore_03    float32 `json:"evaluationScore_03"`
	EvaluationScore_04    float32 `json:"evaluationScore_04"`
	EvaluationScore_05    float32 `json:"evaluationScore_05"`
	EvaluationScore_06    float32 `json:"evaluationScore_06"`
	EvaluationScore_07    float32 `json:"evaluationScore_07"`
	EvaluationScore_08    float32 `json:"evaluationScore_08"`
	EvaluationScore_09    float32 `json:"evaluationScore_09"`
	EvaluationScore_10    float32 `json:"evaluationScore_10"`
}

func ConvertSpiderSpecToTumblebugSpec(spiderSpec SpiderSpecInfo) (TbSpecInfo, error) {
	if spiderSpec.Name == "" {
		err := fmt.Errorf("ConvertSpiderSpecToTumblebugSpec failed; spiderSpec.Name == \"\" ")
		emptyTumblebugSpec := TbSpecInfo{}
		return emptyTumblebugSpec, err
	}

	tumblebugSpec := TbSpecInfo{}

	tumblebugSpec.Name = spiderSpec.Name
	tumblebugSpec.CspSpecName = spiderSpec.Name
	tempUint64, _ := strconv.ParseUint(spiderSpec.VCpu.Count, 10, 16)
	tumblebugSpec.Num_vCPU = uint16(tempUint64)
	tempFloat64, _ := strconv.ParseFloat(spiderSpec.Mem, 32)
	tumblebugSpec.Mem_GiB = uint16(tempFloat64 / 1024) //fmt.Sprintf("%.0f", tempFloat64/1024)

	return tumblebugSpec, nil
}

type SpiderSpecList struct {
	Vmspec []SpiderSpecInfo `json:"vmspec"`
}

func LookupSpecList(connConfig string) (SpiderSpecList, error) {

	if os.Getenv("SPIDER_CALL_METHOD") == "REST" {

		url := common.SPIDER_REST_URL + "/vmspec"

		// Create Req body
		type JsonTemplate struct {
			ConnectionName string `json:"ConnectionName"`
		}
		tempReq := JsonTemplate{}
		tempReq.ConnectionName = connConfig

		client := resty.New()
		client.SetAllowGetMethodPayload(true)

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(tempReq).
			SetResult(&SpiderSpecList{}). // or SetResult(AuthSuccess{}).
			//SetError(&AuthError{}).       // or SetError(AuthError{}).
			Get(url)

		if err != nil {
			common.CBLog.Error(err)
			content := SpiderSpecList{}
			err := fmt.Errorf("an error occurred while requesting to CB-Spider")
			return content, err
		}

		fmt.Println(string(resp.Body()))

		fmt.Println("HTTP Status code " + strconv.Itoa(resp.StatusCode()))
		switch {
		case resp.StatusCode() >= 400 || resp.StatusCode() < 200:
			err := fmt.Errorf(string(resp.Body()))
			common.CBLog.Error(err)
			content := SpiderSpecList{}
			return content, err
		}

		temp := resp.Result().(*SpiderSpecList)
		return *temp, nil

	} else {

		// CCM API 설정
		ccm := api.NewCloudResourceHandler()
		err := ccm.SetConfigPath(os.Getenv("CBTUMBLEBUG_ROOT") + "/conf/grpc_conf.yaml")
		if err != nil {
			common.CBLog.Error("ccm failed to set config : ", err)
			return SpiderSpecList{}, err
		}
		err = ccm.Open()
		if err != nil {
			common.CBLog.Error("ccm api open failed : ", err)
			return SpiderSpecList{}, err
		}
		defer ccm.Close()

		result, err := ccm.ListVMSpecByParam(connConfig)
		if err != nil {
			common.CBLog.Error(err)
			return SpiderSpecList{}, err
		}

		temp := SpiderSpecList{}
		err2 := json.Unmarshal([]byte(result), &temp)
		if err2 != nil {
			fmt.Println("whoops:", err2)
		}
		return temp, nil

	}
}

//func LookupSpec(u *TbSpecInfo) (SpiderSpecInfo, error) {
func LookupSpec(connConfig string, specName string) (SpiderSpecInfo, error) {

	if os.Getenv("SPIDER_CALL_METHOD") == "REST" {

		//url := common.SPIDER_REST_URL + "/vmspec/" + u.CspSpecName
		url := common.SPIDER_REST_URL + "/vmspec/" + specName

		// Create Req body
		type JsonTemplate struct {
			ConnectionName string `json:"ConnectionName"`
		}
		tempReq := JsonTemplate{}
		tempReq.ConnectionName = connConfig

		client := resty.New()
		client.SetAllowGetMethodPayload(true)

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(tempReq).
			SetResult(&SpiderSpecInfo{}). // or SetResult(AuthSuccess{}).
			//SetError(&AuthError{}).       // or SetError(AuthError{}).
			Get(url)

		if err != nil {
			common.CBLog.Error(err)
			content := SpiderSpecInfo{}
			err := fmt.Errorf("an error occurred while requesting to CB-Spider")
			return content, err
		}

		fmt.Println(string(resp.Body()))

		fmt.Println("HTTP Status code " + strconv.Itoa(resp.StatusCode()))
		switch {
		case resp.StatusCode() >= 400 || resp.StatusCode() < 200:
			err := fmt.Errorf(string(resp.Body()))
			common.CBLog.Error(err)
			content := SpiderSpecInfo{}
			return content, err
		}

		temp := resp.Result().(*SpiderSpecInfo)
		return *temp, nil

	} else {

		// CCM API 설정
		ccm := api.NewCloudResourceHandler()
		err := ccm.SetConfigPath(os.Getenv("CBTUMBLEBUG_ROOT") + "/conf/grpc_conf.yaml")
		if err != nil {
			common.CBLog.Error("ccm failed to set config : ", err)
			return SpiderSpecInfo{}, err
		}
		err = ccm.Open()
		if err != nil {
			common.CBLog.Error("ccm api open failed : ", err)
			return SpiderSpecInfo{}, err
		}
		defer ccm.Close()

		result, err := ccm.GetVMSpecByParam(connConfig, specName)
		if err != nil {
			common.CBLog.Error(err)
			return SpiderSpecInfo{}, err
		}

		temp := SpiderSpecInfo{}
		err2 := json.Unmarshal([]byte(result), &temp)
		if err2 != nil {
			fmt.Errorf("an error occurred while unmarshaling: " + err2.Error())
		}
		return temp, nil

	}
}

func FetchSpecs(nsId string) (connConfigCount uint, specCount uint, err error) {

	nsId = common.GenId(nsId)

	connConfigs, err := common.GetConnConfigList()
	if err != nil {
		common.CBLog.Error(err)
		return 0, 0, err
	}

	for _, connConfig := range connConfigs.Connectionconfig {
		fmt.Println("connConfig " + connConfig.ConfigName)

		spiderSpecList, err := LookupSpecList(connConfig.ConfigName)
		if err != nil {
			common.CBLog.Error(err)
			return 0, 0, err
		}

		for _, spiderSpec := range spiderSpecList.Vmspec {
			tumblebugSpec, err := ConvertSpiderSpecToTumblebugSpec(spiderSpec)
			if err != nil {
				common.CBLog.Error(err)
				return 0, 0, err
			}

			tumblebugSpecId := connConfig.ConfigName + "-" + tumblebugSpec.Name
			//fmt.Println("tumblebugSpecId: " + tumblebugSpecId) // for debug

			check, _, err := LowerizeAndCheckResource(nsId, "spec", tumblebugSpecId)
			if check == true {
				common.CBLog.Infoln("The spec " + tumblebugSpecId + " already exists in TB; continue")
				continue
			} else if err != nil {
				common.CBLog.Infoln("Cannot check the existence of " + tumblebugSpecId + " in TB; continue")
				continue
			} else {
				tumblebugSpec.Id = tumblebugSpecId
				tumblebugSpec.Name = tumblebugSpecId
				tumblebugSpec.ConnectionName = connConfig.ConfigName

				_, err := RegisterSpecWithInfo(nsId, &tumblebugSpec)
				if err != nil {
					common.CBLog.Error(err)
					return 0, 0, err
				}
			}
			specCount++
		}
		connConfigCount++
	}
	return connConfigCount, specCount, nil
}

func RegisterSpecWithCspSpecName(nsId string, u *TbSpecReq) (TbSpecInfo, error) {

	nsId = common.GenId(nsId)

	_, lowerizedNsId, _ := common.LowerizeAndCheckNs(nsId)
	nsId = lowerizedNsId

	check, lowerizedName, err := LowerizeAndCheckResource(nsId, "spec", u.Name)
	u.Name = lowerizedName

	if check == true {
		temp := TbSpecInfo{}
		err := fmt.Errorf("The spec " + u.Name + " already exists.")
		return temp, err
	}

	if err != nil {
		temp := TbSpecInfo{}
		err := fmt.Errorf("Failed to check the existence of the spec " + lowerizedName + ".")
		return temp, err
	}

	res, err := LookupSpec(u.ConnectionName, u.CspSpecName)
	if err != nil {
		common.CBLog.Error(err)
		err := fmt.Errorf("an error occurred while lookup spec via CB-Spider")
		emptySpecInfoObj := TbSpecInfo{}
		return emptySpecInfoObj, err
	}

	content := TbSpecInfo{}
	//content.Id = common.GenUuid()
	content.Id = common.GenId(u.Name)
	content.Name = common.GenId(u.Name)
	content.CspSpecName = res.Name
	content.ConnectionName = u.ConnectionName

	tempUint64, _ := strconv.ParseUint(res.VCpu.Count, 10, 16)
	content.Num_vCPU = uint16(tempUint64)

	//content.Num_core = res.Num_core

	tempFloat64, _ := strconv.ParseFloat(res.Mem, 32)
	content.Mem_GiB = uint16(tempFloat64 / 1024)

	//content.Storage_GiB = res.Storage_GiB
	//content.Description = res.Description

	sql := "INSERT INTO `spec`(" +
		"`namespace`, " +
		"`id`, " +
		"`connectionName`, " +
		"`cspSpecName`, " +
		"`name`, " +
		"`os_type`, " +
		"`num_vCPU`, " +
		"`num_core`, " +
		"`mem_GiB`, " +
		"`storage_GiB`, " +
		"`description`, " +
		"`cost_per_hour`, " +
		"`num_storage`, " +
		"`max_num_storage`, " +
		"`max_total_storage_TiB`, " +
		"`net_bw_Gbps`, " +
		"`ebs_bw_Mbps`, " +
		"`gpu_model`, " +
		"`num_gpu`, " +
		"`gpumem_GiB`, " +
		"`gpu_p2p`, " +
		"`orderInFilteredResult`, " +
		"`evaluationStatus`, " +
		"`evaluationScore_01`, " +
		"`evaluationScore_02`, " +
		"`evaluationScore_03`, " +
		"`evaluationScore_04`, " +
		"`evaluationScore_05`, " +
		"`evaluationScore_06`, " +
		"`evaluationScore_07`, " +
		"`evaluationScore_08`, " +
		"`evaluationScore_09`, " +
		"`evaluationScore_10`) " +
		"VALUES ('" +
		nsId + "', '" +
		content.Id + "', '" +
		content.ConnectionName + "', '" +
		content.CspSpecName + "', '" +
		content.Name + "', '" +
		content.Os_type + "', '" +
		strconv.Itoa(int(content.Num_vCPU)) + "', '" +
		strconv.Itoa(int(content.Num_core)) + "', '" +
		strconv.Itoa(int(content.Mem_GiB)) + "', '" +
		strconv.Itoa(int(content.Storage_GiB)) + "', '" +
		content.Description + "', '" +
		fmt.Sprintf("%.6f", content.Cost_per_hour) + "', '" +
		strconv.Itoa(int(content.Num_storage)) + "', '" +
		strconv.Itoa(int(content.Max_num_storage)) + "', '" +
		strconv.Itoa(int(content.Max_total_storage_TiB)) + "', '" +
		strconv.Itoa(int(content.Net_bw_Gbps)) + "', '" +
		strconv.Itoa(int(content.Ebs_bw_Mbps)) + "', '" +
		content.Gpu_model + "', '" +
		strconv.Itoa(int(content.Num_gpu)) + "', '" +
		strconv.Itoa(int(content.Gpumem_GiB)) + "', '" +
		content.Gpu_p2p + "', '" +
		strconv.Itoa(int(content.OrderInFilteredResult)) + "', '" +
		content.EvaluationStatus + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_01) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_02) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_03) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_04) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_05) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_06) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_07) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_08) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_09) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_10) + "');"

	fmt.Println("sql: " + sql)
	// https://stackoverflow.com/questions/42486032/golang-sql-query-syntax-validator
	_, err = sqlparser.Parse(sql)
	if err != nil {
		return content, err
	}

	// cb-store
	fmt.Println("=========================== PUT registerSpec")
	Key := common.GenResourceKey(nsId, "spec", content.Id)
	Val, _ := json.Marshal(content)
	err = common.CBStore.Put(string(Key), string(Val))
	if err != nil {
		common.CBLog.Error(err)
		return content, err
	}
	keyValue, _ := common.CBStore.Get(string(Key))
	fmt.Println("<" + keyValue.Key + "> \n" + keyValue.Value)
	fmt.Println("===========================")

	// register information related with MCIS recommendation
	RegisterRecommendList(nsId, content.ConnectionName, content.Num_vCPU, content.Mem_GiB, content.Storage_GiB, content.Id, content.Cost_per_hour)

	stmt, err := common.MYDB.Prepare(sql)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Data inserted successfully..")
	}

	return content, nil
}

func RegisterSpecWithInfo(nsId string, content *TbSpecInfo) (TbSpecInfo, error) {

	nsId = common.GenId(nsId)

	//_, lowerizedNsId, _ := common.LowerizeAndCheckNs(nsId)
	//nsId = lowerizedNsId
	nsId = common.GenId(nsId)

	check, lowerizedName, err := LowerizeAndCheckResource(nsId, "spec", content.Name)
	content.Name = lowerizedName

	if check == true {
		temp := TbSpecInfo{}
		err := fmt.Errorf("The spec " + content.Name + " already exists.")
		return temp, err
	}

	content.Id = content.Name
	//content.Name = content.Name

	sql := "INSERT INTO `spec`(" +
		"`namespace`, " +
		"`id`, " +
		"`connectionName`, " +
		"`cspSpecName`, " +
		"`name`, " +
		"`os_type`, " +
		"`num_vCPU`, " +
		"`num_core`, " +
		"`mem_GiB`, " +
		"`storage_GiB`, " +
		"`description`, " +
		"`cost_per_hour`, " +
		"`num_storage`, " +
		"`max_num_storage`, " +
		"`max_total_storage_TiB`, " +
		"`net_bw_Gbps`, " +
		"`ebs_bw_Mbps`, " +
		"`gpu_model`, " +
		"`num_gpu`, " +
		"`gpumem_GiB`, " +
		"`gpu_p2p`, " +
		"`orderInFilteredResult`, " +
		"`evaluationStatus`, " +
		"`evaluationScore_01`, " +
		"`evaluationScore_02`, " +
		"`evaluationScore_03`, " +
		"`evaluationScore_04`, " +
		"`evaluationScore_05`, " +
		"`evaluationScore_06`, " +
		"`evaluationScore_07`, " +
		"`evaluationScore_08`, " +
		"`evaluationScore_09`, " +
		"`evaluationScore_10`) " +
		"VALUES ('" +
		nsId + "', '" +
		content.Id + "', '" +
		content.ConnectionName + "', '" +
		content.CspSpecName + "', '" +
		content.Name + "', '" +
		content.Os_type + "', '" +
		strconv.Itoa(int(content.Num_vCPU)) + "', '" +
		strconv.Itoa(int(content.Num_core)) + "', '" +
		strconv.Itoa(int(content.Mem_GiB)) + "', '" +
		strconv.Itoa(int(content.Storage_GiB)) + "', '" +
		content.Description + "', '" +
		fmt.Sprintf("%.6f", content.Cost_per_hour) + "', '" +
		strconv.Itoa(int(content.Num_storage)) + "', '" +
		strconv.Itoa(int(content.Max_num_storage)) + "', '" +
		strconv.Itoa(int(content.Max_total_storage_TiB)) + "', '" +
		strconv.Itoa(int(content.Net_bw_Gbps)) + "', '" +
		strconv.Itoa(int(content.Ebs_bw_Mbps)) + "', '" +
		content.Gpu_model + "', '" +
		strconv.Itoa(int(content.Num_gpu)) + "', '" +
		strconv.Itoa(int(content.Gpumem_GiB)) + "', '" +
		content.Gpu_p2p + "', '" +
		strconv.Itoa(int(content.OrderInFilteredResult)) + "', '" +
		content.EvaluationStatus + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_01) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_02) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_03) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_04) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_05) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_06) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_07) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_08) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_09) + "', '" +
		fmt.Sprintf("%.6f", content.EvaluationScore_10) + "');"

	fmt.Println("sql: " + sql)
	// https://stackoverflow.com/questions/42486032/golang-sql-query-syntax-validator
	_, err = sqlparser.Parse(sql)
	if err != nil {
		return *content, err
	}

	// cb-store
	fmt.Println("=========================== PUT registerSpec")
	Key := common.GenResourceKey(nsId, "spec", content.Id)
	Val, _ := json.Marshal(content)
	err = common.CBStore.Put(string(Key), string(Val))
	if err != nil {
		common.CBLog.Error(err)
		return *content, err
	}
	keyValue, _ := common.CBStore.Get(string(Key))
	fmt.Println("<" + keyValue.Key + "> \n" + keyValue.Value)
	fmt.Println("===========================")

	// register information related with MCIS recommendation
	RegisterRecommendList(nsId, content.ConnectionName, content.Num_vCPU, content.Mem_GiB, content.Storage_GiB, content.Id, content.Cost_per_hour)

	stmt, err := common.MYDB.Prepare(sql)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Data inserted successfully..")
	}

	return *content, nil
}

func RegisterRecommendList(nsId string, connectionName string, cpuSize uint16, memSize uint16, diskSize uint32, specId string, price float32) error {

	nsId = common.GenId(nsId)

	//fmt.Println("[Get MCISs")
	key := common.GenMcisKey(nsId, "", "") + "/cpuSize/" + strconv.Itoa(int(cpuSize)) + "/memSize/" + strconv.Itoa(int(memSize)) + "/diskSize/" + strconv.Itoa(int(diskSize)) + "/specId/" + specId
	fmt.Println(key)

	mapA := map[string]string{"id": specId, "price": fmt.Sprintf("%.6f", price), "connectionName": connectionName}
	Val, _ := json.Marshal(mapA)

	err := common.CBStore.Put(string(key), string(Val))
	if err != nil {
		common.CBLog.Error(err)
		return err
	}

	fmt.Println("===============================================")
	return nil

}

func DelRecommendSpec(nsId string, specId string, cpuSize uint16, memSize uint16, diskSize uint32) error {

	nsId = common.GenId(nsId)

	fmt.Println("DelRecommendSpec()")

	key := common.GenMcisKey(nsId, "", "") + "/cpuSize/" + strconv.Itoa(int(cpuSize)) + "/memSize/" + strconv.Itoa(int(memSize)) + "/diskSize/" + strconv.Itoa(int(diskSize)) + "/specId/" + specId

	err := common.CBStore.Delete(key)
	if err != nil {
		common.CBLog.Error(err)
		return err
	}

	return nil

}

func FilterSpecs(nsId string, filter TbSpecInfo) ([]TbSpecInfo, error) {

	nsId = common.GenId(nsId)

	tempList := []TbSpecInfo{}

	sqlQuery := "SELECT * FROM `spec` WHERE `namespace`='" + nsId + "'"

	if filter.Id != "" {
		sqlQuery += " AND `id`='" + filter.Id + "'"
	}
	if filter.Name != "" {
		sqlQuery += " AND `name` LIKE '%" + filter.Name + "%'"
	}
	if filter.ConnectionName != "" {
		sqlQuery += " AND `connectionName`='" + filter.ConnectionName + "'"
	}
	if filter.CspSpecName != "" {
		sqlQuery += " AND `cspSpecName` LIKE '%" + filter.CspSpecName + "%'"
	}
	if filter.Os_type != "" {
		sqlQuery += " AND `os_type` LIKE '%" + filter.Os_type + "%'"
	}

	if filter.Num_vCPU > 0 {
		sqlQuery += " AND `num_vCPU`=" + strconv.Itoa(int(filter.Num_vCPU))
	}
	if filter.Num_core > 0 {
		sqlQuery += " AND `num_core`=" + strconv.Itoa(int(filter.Num_core))
	}
	if filter.Mem_GiB > 0 {
		sqlQuery += " AND `mem_GiB`=" + strconv.Itoa(int(filter.Mem_GiB))
	}
	if filter.Storage_GiB > 0 {
		sqlQuery += " AND `storage_GiB`=" + strconv.Itoa(int(filter.Storage_GiB))
	}
	if filter.Description != "" {
		sqlQuery += " AND `description` LIKE '%" + filter.Description + "%'"
	}
	if filter.Cost_per_hour > 0 {
		sqlQuery += " AND `cost_per_hour`=" + fmt.Sprintf("%.6f", filter.Cost_per_hour)
	}
	if filter.Num_storage > 0 {
		sqlQuery += " AND `num_storage`=" + strconv.Itoa(int(filter.Num_storage))
	}
	if filter.Max_num_storage > 0 {
		sqlQuery += " AND `max_num_storage`=" + strconv.Itoa(int(filter.Max_num_storage))
	}
	if filter.Max_total_storage_TiB > 0 {
		sqlQuery += " AND `max_total_storage_TiB`=" + strconv.Itoa(int(filter.Max_total_storage_TiB))
	}
	if filter.Net_bw_Gbps > 0 {
		sqlQuery += " AND `net_bw_Gbps`=" + strconv.Itoa(int(filter.Net_bw_Gbps))
	}
	if filter.Ebs_bw_Mbps > 0 {
		sqlQuery += " AND `ebs_bw_Mbps`=" + strconv.Itoa(int(filter.Ebs_bw_Mbps))
	}
	if filter.Gpu_model != "" {
		sqlQuery += " AND `gpu_model` LIKE '%" + filter.Gpu_model + "%'"
	}
	if filter.Num_gpu > 0 {
		sqlQuery += " AND `num_gpu`=" + strconv.Itoa(int(filter.Num_gpu))
	}
	if filter.Gpumem_GiB > 0 {
		sqlQuery += " AND `gpumem_GiB`=" + strconv.Itoa(int(filter.Gpumem_GiB))
	}
	if filter.Gpu_p2p != "" {
		sqlQuery += " AND `gpu_p2p` LIKE '%" + filter.Gpu_p2p + "%'"
	}
	if filter.EvaluationStatus != "" {
		sqlQuery += " AND `evaluationStatus`='" + filter.EvaluationStatus + "'"
	}
	if filter.EvaluationScore_01 > 0 {
		sqlQuery += " AND `evaluationScore_01`=" + fmt.Sprintf("%.6f", filter.EvaluationScore_01)
	}
	if filter.EvaluationScore_02 > 0 {
		sqlQuery += " AND `evaluationScore_02`=" + fmt.Sprintf("%.6f", filter.EvaluationScore_02)
	}
	if filter.EvaluationScore_03 > 0 {
		sqlQuery += " AND `evaluationScore_03`=" + fmt.Sprintf("%.6f", filter.EvaluationScore_03)
	}
	if filter.EvaluationScore_04 > 0 {
		sqlQuery += " AND `evaluationScore_04`=" + fmt.Sprintf("%.6f", filter.EvaluationScore_04)
	}
	if filter.EvaluationScore_05 > 0 {
		sqlQuery += " AND `evaluationScore_05`=" + fmt.Sprintf("%.6f", filter.EvaluationScore_05)
	}
	if filter.EvaluationScore_06 > 0 {
		sqlQuery += " AND `evaluationScore_06`=" + fmt.Sprintf("%.6f", filter.EvaluationScore_06)
	}
	if filter.EvaluationScore_07 > 0 {
		sqlQuery += " AND `evaluationScore_07`=" + fmt.Sprintf("%.6f", filter.EvaluationScore_07)
	}
	if filter.EvaluationScore_08 > 0 {
		sqlQuery += " AND `evaluationScore_08`=" + fmt.Sprintf("%.6f", filter.EvaluationScore_08)
	}
	if filter.EvaluationScore_09 > 0 {
		sqlQuery += " AND `evaluationScore_09`=" + fmt.Sprintf("%.6f", filter.EvaluationScore_09)
	}
	if filter.EvaluationScore_10 > 0 {
		sqlQuery += " AND `evaluationScore_10`=" + fmt.Sprintf("%.6f", filter.EvaluationScore_10)
	}
	sqlQuery += ";"
	_, err := sqlparser.Parse(sqlQuery)
	if err != nil {
		return tempList, err
	}

	/*
		stmt, err := common.MYDB.Prepare(sqlQuery)
		if err != nil {
			fmt.Println(err.Error())
		}
		result, err := stmt.Exec()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("SELECTed successfully..")
		}

		result.RowsAffected

		temp := []TbSpecInfo{}
		return temp, nil
	*/

	rows, err := common.MYDB.Query(sqlQuery)
	if err != nil {
		common.CBLog.Error(err)
		return tempList, err
	}

	cols, err := rows.Columns()
	if err != nil {
		common.CBLog.Error(err)
		return tempList, err
	}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return tempList, err
		}
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		js, _ := json.Marshal(m)
		tempSpec := TbSpecInfo{}
		json.Unmarshal(js, &tempSpec)
		tempList = append(tempList, tempSpec)
	}
	return tempList, nil
}

type Range struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
}

type FilterSpecsByRangeRequest struct {
	Num_vCPU              Range `json:"num_vCPU"`
	Num_core              Range `json:"num_core"`
	Mem_GiB               Range `json:"mem_GiB"`
	Storage_GiB           Range `json:"storage_GiB"`
	Cost_per_hour         Range `json:"cost_per_hour"`
	Num_storage           Range `json:"num_storage"`
	Max_num_storage       Range `json:"max_num_storage"`
	Max_total_storage_TiB Range `json:"max_total_storage_TiB"`
	Net_bw_Gbps           Range `json:"net_bw_Gbps"`
	Ebs_bw_Mbps           Range `json:"ebs_bw_Mbps"`
	Num_gpu               Range `json:"num_gpu"`
	Gpumem_GiB            Range `json:"gpumem_GiB"`
	EvaluationScore_01    Range `json:"evaluationScore_01"`
	EvaluationScore_02    Range `json:"evaluationScore_02"`
	EvaluationScore_03    Range `json:"evaluationScore_03"`
	EvaluationScore_04    Range `json:"evaluationScore_04"`
	EvaluationScore_05    Range `json:"evaluationScore_05"`
	EvaluationScore_06    Range `json:"evaluationScore_06"`
	EvaluationScore_07    Range `json:"evaluationScore_07"`
	EvaluationScore_08    Range `json:"evaluationScore_08"`
	EvaluationScore_09    Range `json:"evaluationScore_09"`
	EvaluationScore_10    Range `json:"evaluationScore_10"`
}

func FilterSpecsByRange(nsId string, filter FilterSpecsByRangeRequest) ([]TbSpecInfo, error) {
	nsId = common.GenId(nsId)

	tempList := []TbSpecInfo{}

	sqlQuery := "SELECT * FROM `spec` WHERE `namespace`='" + nsId + "'"

	if filter.Num_vCPU.Min > 0 {
		sqlQuery += " AND `num_vCPU`>=" + fmt.Sprintf("%.6f", filter.Num_vCPU.Min)
	}
	if filter.Num_vCPU.Max > 0 {
		sqlQuery += " AND `num_vCPU`<=" + fmt.Sprintf("%.6f", filter.Num_vCPU.Max)
	}

	if filter.Num_core.Min > 0 {
		sqlQuery += " AND `num_core`>=" + fmt.Sprintf("%.6f", filter.Num_core.Min)
	}
	if filter.Num_core.Max > 0 {
		sqlQuery += " AND `num_core`<=" + fmt.Sprintf("%.6f", filter.Num_core.Max)
	}

	if filter.Mem_GiB.Min > 0 {
		sqlQuery += " AND `mem_GiB`>=" + fmt.Sprintf("%.6f", filter.Mem_GiB.Min)
	}
	if filter.Mem_GiB.Max > 0 {
		sqlQuery += " AND `mem_GiB`<=" + fmt.Sprintf("%.6f", filter.Mem_GiB.Max)
	}

	if filter.Storage_GiB.Min > 0 {
		sqlQuery += " AND `storage_GiB`>=" + fmt.Sprintf("%.6f", filter.Storage_GiB.Min)
	}
	if filter.Storage_GiB.Max > 0 {
		sqlQuery += " AND `storage_GiB`<=" + fmt.Sprintf("%.6f", filter.Storage_GiB.Max)
	}

	if filter.Cost_per_hour.Min > 0 {
		sqlQuery += " AND `cost_per_hour`>=" + fmt.Sprintf("%.6f", filter.Cost_per_hour.Min)
	}
	if filter.Cost_per_hour.Max > 0 {
		sqlQuery += " AND `cost_per_hour`<=" + fmt.Sprintf("%.6f", filter.Cost_per_hour.Max)
	}

	if filter.Num_storage.Min > 0 {
		sqlQuery += " AND `num_storage`>=" + fmt.Sprintf("%.6f", filter.Num_storage.Min)
	}
	if filter.Num_storage.Max > 0 {
		sqlQuery += " AND `num_storage`<=" + fmt.Sprintf("%.6f", filter.Num_storage.Max)
	}

	if filter.Max_num_storage.Min > 0 {
		sqlQuery += " AND `max_num_storage`>=" + fmt.Sprintf("%.6f", filter.Max_num_storage.Min)
	}
	if filter.Max_num_storage.Max > 0 {
		sqlQuery += " AND `max_num_storage`<=" + fmt.Sprintf("%.6f", filter.Max_num_storage.Max)
	}

	if filter.Max_total_storage_TiB.Min > 0 {
		sqlQuery += " AND `max_total_storage_TiB`>=" + fmt.Sprintf("%.6f", filter.Max_total_storage_TiB.Min)
	}
	if filter.Max_total_storage_TiB.Max > 0 {
		sqlQuery += " AND `max_total_storage_TiB`<=" + fmt.Sprintf("%.6f", filter.Max_total_storage_TiB.Max)
	}

	if filter.Net_bw_Gbps.Min > 0 {
		sqlQuery += " AND `net_bw_Gbps`>=" + fmt.Sprintf("%.6f", filter.Net_bw_Gbps.Min)
	}
	if filter.Net_bw_Gbps.Max > 0 {
		sqlQuery += " AND `net_bw_Gbps`<=" + fmt.Sprintf("%.6f", filter.Net_bw_Gbps.Max)
	}

	if filter.Ebs_bw_Mbps.Min > 0 {
		sqlQuery += " AND `ebs_bw_Mbps`>=" + fmt.Sprintf("%.6f", filter.Ebs_bw_Mbps.Min)
	}
	if filter.Ebs_bw_Mbps.Max > 0 {
		sqlQuery += " AND `ebs_bw_Mbps`<=" + fmt.Sprintf("%.6f", filter.Ebs_bw_Mbps.Max)
	}

	if filter.Num_gpu.Min > 0 {
		sqlQuery += " AND `num_gpu`>=" + fmt.Sprintf("%.6f", filter.Num_gpu.Min)
	}
	if filter.Num_gpu.Max > 0 {
		sqlQuery += " AND `num_gpu`<=" + fmt.Sprintf("%.6f", filter.Num_gpu.Max)
	}

	if filter.Gpumem_GiB.Min > 0 {
		sqlQuery += " AND `gpumem_GiB`>=" + fmt.Sprintf("%.6f", filter.Gpumem_GiB.Min)
	}
	if filter.Gpumem_GiB.Max > 0 {
		sqlQuery += " AND `gpumem_GiB`<=" + fmt.Sprintf("%.6f", filter.Gpumem_GiB.Max)
	}

	if filter.EvaluationScore_01.Min > 0 {
		sqlQuery += " AND `evaluationScore_01`>=" + fmt.Sprintf("%.6f", filter.EvaluationScore_01.Min)
	}
	if filter.EvaluationScore_01.Max > 0 {
		sqlQuery += " AND `evaluationScore_01`<=" + fmt.Sprintf("%.6f", filter.EvaluationScore_01.Max)
	}

	if filter.EvaluationScore_02.Min > 0 {
		sqlQuery += " AND `evaluationScore_02`>=" + fmt.Sprintf("%.6f", filter.EvaluationScore_02.Min)
	}
	if filter.EvaluationScore_02.Max > 0 {
		sqlQuery += " AND `evaluationScore_02`<=" + fmt.Sprintf("%.6f", filter.EvaluationScore_02.Max)
	}

	if filter.EvaluationScore_03.Min > 0 {
		sqlQuery += " AND `evaluationScore_03`>=" + fmt.Sprintf("%.6f", filter.EvaluationScore_03.Min)
	}
	if filter.EvaluationScore_03.Max > 0 {
		sqlQuery += " AND `evaluationScore_03`<=" + fmt.Sprintf("%.6f", filter.EvaluationScore_03.Max)
	}

	if filter.EvaluationScore_04.Min > 0 {
		sqlQuery += " AND `evaluationScore_04`>=" + fmt.Sprintf("%.6f", filter.EvaluationScore_04.Min)
	}
	if filter.EvaluationScore_04.Max > 0 {
		sqlQuery += " AND `evaluationScore_04`<=" + fmt.Sprintf("%.6f", filter.EvaluationScore_04.Max)
	}

	if filter.EvaluationScore_05.Min > 0 {
		sqlQuery += " AND `evaluationScore_05`>=" + fmt.Sprintf("%.6f", filter.EvaluationScore_05.Min)
	}
	if filter.EvaluationScore_05.Max > 0 {
		sqlQuery += " AND `evaluationScore_05`<=" + fmt.Sprintf("%.6f", filter.EvaluationScore_05.Max)
	}

	if filter.EvaluationScore_06.Min > 0 {
		sqlQuery += " AND `evaluationScore_06`>=" + fmt.Sprintf("%.6f", filter.EvaluationScore_06.Min)
	}
	if filter.EvaluationScore_06.Max > 0 {
		sqlQuery += " AND `evaluationScore_06`<=" + fmt.Sprintf("%.6f", filter.EvaluationScore_06.Max)
	}

	if filter.EvaluationScore_07.Min > 0 {
		sqlQuery += " AND `evaluationScore_07`>=" + fmt.Sprintf("%.6f", filter.EvaluationScore_07.Min)
	}
	if filter.EvaluationScore_07.Max > 0 {
		sqlQuery += " AND `evaluationScore_07`<=" + fmt.Sprintf("%.6f", filter.EvaluationScore_07.Max)
	}

	if filter.EvaluationScore_08.Min > 0 {
		sqlQuery += " AND `evaluationScore_08`>=" + fmt.Sprintf("%.6f", filter.EvaluationScore_08.Min)
	}
	if filter.EvaluationScore_08.Max > 0 {
		sqlQuery += " AND `evaluationScore_08`<=" + fmt.Sprintf("%.6f", filter.EvaluationScore_08.Max)
	}

	if filter.EvaluationScore_09.Min > 0 {
		sqlQuery += " AND `evaluationScore_09`>=" + fmt.Sprintf("%.6f", filter.EvaluationScore_09.Min)
	}
	if filter.EvaluationScore_09.Max > 0 {
		sqlQuery += " AND `evaluationScore_09`<=" + fmt.Sprintf("%.6f", filter.EvaluationScore_09.Max)
	}

	if filter.EvaluationScore_10.Min > 0 {
		sqlQuery += " AND `evaluationScore_10`>=" + fmt.Sprintf("%.6f", filter.EvaluationScore_10.Min)
	}
	if filter.EvaluationScore_10.Max > 0 {
		sqlQuery += " AND `evaluationScore_10`<=" + fmt.Sprintf("%.6f", filter.EvaluationScore_10.Max)
	}

	sqlQuery += ";"
	_, err := sqlparser.Parse(sqlQuery)
	if err != nil {
		return tempList, err
	}

	rows, err := common.MYDB.Query(sqlQuery)
	if err != nil {
		common.CBLog.Error(err)
		return tempList, err
	}

	cols, err := rows.Columns()
	if err != nil {
		common.CBLog.Error(err)
		return tempList, err
	}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return tempList, err
		}
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		js, _ := json.Marshal(m)
		tempSpec := TbSpecInfo{}
		json.Unmarshal(js, &tempSpec)
		tempList = append(tempList, tempSpec)
	}
	return tempList, nil
}

func SortSpecs(specList []TbSpecInfo, orderBy string, direction string) ([]TbSpecInfo, error) {
	var err error = nil

	sort.Slice(specList, func(i, j int) bool {
		if orderBy == "num_vCPU" {
			if direction == "descending" {
				return specList[i].Num_vCPU > specList[j].Num_vCPU
			} else if direction == "ascending" {
				return specList[i].Num_vCPU < specList[j].Num_vCPU
			} else {
				err = fmt.Errorf("'direction' should one of these: ascending, descending")
				return true
			}
		} else if orderBy == "mem_GiB" {
			if direction == "descending" {
				return specList[i].Mem_GiB > specList[j].Mem_GiB
			} else if direction == "ascending" {
				return specList[i].Mem_GiB < specList[j].Mem_GiB
			} else {
				err = fmt.Errorf("'direction' should one of these: ascending, descending")
				return true
			}
		} else if orderBy == "storage_GiB" {
			if direction == "descending" {
				return specList[i].Storage_GiB > specList[j].Storage_GiB
			} else if direction == "ascending" {
				return specList[i].Storage_GiB < specList[j].Storage_GiB
			} else {
				err = fmt.Errorf("'direction' should one of these: ascending, descending")
				return true
			}
		} else if orderBy == "evaluationScore_01" {
			if direction == "descending" {
				return specList[i].EvaluationScore_01 > specList[j].EvaluationScore_01
			} else if direction == "ascending" {
				return specList[i].EvaluationScore_01 < specList[j].EvaluationScore_01
			} else {
				err = fmt.Errorf("'direction' should one of these: ascending, descending")
				return true
			}
		} else if orderBy == "evaluationScore_02" {
			if direction == "descending" {
				return specList[i].EvaluationScore_02 > specList[j].EvaluationScore_02
			} else if direction == "ascending" {
				return specList[i].EvaluationScore_02 < specList[j].EvaluationScore_02
			} else {
				err = fmt.Errorf("'direction' should one of these: ascending, descending")
				return true
			}
		} else if orderBy == "evaluationScore_03" {
			if direction == "descending" {
				return specList[i].EvaluationScore_03 > specList[j].EvaluationScore_03
			} else if direction == "ascending" {
				return specList[i].EvaluationScore_03 < specList[j].EvaluationScore_03
			} else {
				err = fmt.Errorf("'direction' should one of these: ascending, descending")
				return true
			}
		} else if orderBy == "evaluationScore_04" {
			if direction == "descending" {
				return specList[i].EvaluationScore_04 > specList[j].EvaluationScore_04
			} else if direction == "ascending" {
				return specList[i].EvaluationScore_04 < specList[j].EvaluationScore_04
			} else {
				err = fmt.Errorf("'direction' should one of these: ascending, descending")
				return true
			}
		} else if orderBy == "evaluationScore_05" {
			if direction == "descending" {
				return specList[i].EvaluationScore_05 > specList[j].EvaluationScore_05
			} else if direction == "ascending" {
				return specList[i].EvaluationScore_05 < specList[j].EvaluationScore_05
			} else {
				err = fmt.Errorf("'direction' should one of these: ascending, descending")
				return true
			}
		} else if orderBy == "evaluationScore_06" {
			if direction == "descending" {
				return specList[i].EvaluationScore_06 > specList[j].EvaluationScore_06
			} else if direction == "ascending" {
				return specList[i].EvaluationScore_06 < specList[j].EvaluationScore_06
			} else {
				err = fmt.Errorf("'direction' should one of these: ascending, descending")
				return true
			}
		} else if orderBy == "evaluationScore_07" {
			if direction == "descending" {
				return specList[i].EvaluationScore_07 > specList[j].EvaluationScore_07
			} else if direction == "ascending" {
				return specList[i].EvaluationScore_07 < specList[j].EvaluationScore_07
			} else {
				err = fmt.Errorf("'direction' should one of these: ascending, descending")
				return true
			}
		} else if orderBy == "evaluationScore_08" {
			if direction == "descending" {
				return specList[i].EvaluationScore_08 > specList[j].EvaluationScore_08
			} else if direction == "ascending" {
				return specList[i].EvaluationScore_08 < specList[j].EvaluationScore_08
			} else {
				err = fmt.Errorf("'direction' should one of these: ascending, descending")
				return true
			}
		} else if orderBy == "evaluationScore_09" {
			if direction == "descending" {
				return specList[i].EvaluationScore_09 > specList[j].EvaluationScore_09
			} else if direction == "ascending" {
				return specList[i].EvaluationScore_09 < specList[j].EvaluationScore_09
			} else {
				err = fmt.Errorf("'direction' should one of these: ascending, descending")
				return true
			}
		} else if orderBy == "evaluationScore_10" {
			if direction == "descending" {
				return specList[i].EvaluationScore_10 > specList[j].EvaluationScore_10
			} else if direction == "ascending" {
				return specList[i].EvaluationScore_10 < specList[j].EvaluationScore_10
			} else {
				err = fmt.Errorf("'direction' should one of these: ascending, descending")
				return true
			}
		} else {
			err = fmt.Errorf("'orderBy' should one of these: num_vCPU, mem_GiB, storage_GiB")
			return true
		}
	})

	for i, _ := range specList {
		specList[i].OrderInFilteredResult = uint16(i + 1)
	}

	return specList, err
}

func UpdateSpec(nsId string, newSpec TbSpecInfo) (TbSpecInfo, error) {
	nsId = common.GenId(nsId)

	check, lowerizedName, err := LowerizeAndCheckResource(nsId, "spec", newSpec.Id)
	newSpec.Id = lowerizedName

	if check == false {
		temp := TbSpecInfo{}
		err := fmt.Errorf("The spec " + newSpec.Id + " does not exist.")
		return temp, err
	}

	if err != nil {
		temp := TbSpecInfo{}
		err := fmt.Errorf("Failed to check the existence of the spec " + lowerizedName + ".")
		return temp, err
	}

	tempInterface, err := GetResource(nsId, "spec", newSpec.Id)
	if err != nil {
		temp := TbSpecInfo{}
		err := fmt.Errorf("Failed to get the spec " + lowerizedName + ".")
		return temp, err
	}
	tempSpec := TbSpecInfo{}
	err = common.CopySrcToDest(&tempInterface, &tempSpec)
	if err != nil {
		temp := TbSpecInfo{}
		err := fmt.Errorf("Failed to CopySrcToDest() " + lowerizedName + ".")
		return temp, err
	}

	sqlQuery := "UPDATE `spec` SET `id`='" + newSpec.Id + "'"

	if newSpec.Name != "" {
		tempSpec.Name = newSpec.Name
		sqlQuery += ", `name`='" + newSpec.Name + "'"
	}
	if newSpec.Os_type != "" {
		tempSpec.Os_type = newSpec.Os_type
		sqlQuery += ", `os_type`='" + newSpec.Os_type + "'"
	}
	if newSpec.Num_vCPU > 0 {
		tempSpec.Num_vCPU = newSpec.Num_vCPU
		sqlQuery += ", `num_vCPU`='" + strconv.Itoa(int(newSpec.Num_vCPU)) + "'"
	}
	if newSpec.Num_core > 0 {
		tempSpec.Num_core = newSpec.Num_core
		sqlQuery += ", `num_core`='" + strconv.Itoa(int(newSpec.Num_core)) + "'"
	}
	if newSpec.Mem_GiB > 0 {
		tempSpec.Mem_GiB = newSpec.Mem_GiB
		sqlQuery += ", `mem_GiB`='" + strconv.Itoa(int(newSpec.Mem_GiB)) + "'"
	}
	if newSpec.Storage_GiB > 0 {
		tempSpec.Storage_GiB = newSpec.Storage_GiB
		sqlQuery += ", `storage_GiB`='" + strconv.Itoa(int(newSpec.Storage_GiB)) + "'"
	}
	if newSpec.Description != "" {
		tempSpec.Description = newSpec.Description
		sqlQuery += ", `description`='" + newSpec.Description + "'"
	}
	if newSpec.Cost_per_hour > 0 {
		tempSpec.Cost_per_hour = newSpec.Cost_per_hour
		sqlQuery += ", `cost_per_hour`='" + fmt.Sprintf("%.6f", newSpec.Cost_per_hour) + "'"
	}
	if newSpec.Num_storage > 0 {
		tempSpec.Num_storage = newSpec.Num_storage
		sqlQuery += ", `num_storage`='" + strconv.Itoa(int(newSpec.Num_storage)) + "'"
	}
	if newSpec.Max_num_storage > 0 {
		tempSpec.Max_num_storage = newSpec.Max_num_storage
		sqlQuery += ", `max_num_storage`='" + strconv.Itoa(int(newSpec.Max_num_storage)) + "'"
	}
	if newSpec.Max_total_storage_TiB > 0 {
		tempSpec.Max_total_storage_TiB = newSpec.Max_total_storage_TiB
		sqlQuery += ", `max_total_storage_TiB`='" + strconv.Itoa(int(newSpec.Max_total_storage_TiB)) + "'"
	}
	if newSpec.Net_bw_Gbps > 0 {
		tempSpec.Net_bw_Gbps = newSpec.Net_bw_Gbps
		sqlQuery += ", `net_bw_Gbps`='" + strconv.Itoa(int(newSpec.Net_bw_Gbps)) + "'"
	}
	if newSpec.Ebs_bw_Mbps > 0 {
		tempSpec.Ebs_bw_Mbps = newSpec.Ebs_bw_Mbps
		sqlQuery += ", `ebs_bw_Mbps`='" + strconv.Itoa(int(newSpec.Ebs_bw_Mbps)) + "'"
	}
	if newSpec.Gpu_model != "" {
		tempSpec.Gpu_model = newSpec.Gpu_model
		sqlQuery += ", `gpu_model`='" + newSpec.Gpu_model + "'"
	}
	if newSpec.Num_gpu > 0 {
		tempSpec.Num_gpu = newSpec.Num_gpu
		sqlQuery += ", `num_gpu`='" + strconv.Itoa(int(newSpec.Num_gpu)) + "'"
	}
	if newSpec.Gpumem_GiB > 0 {
		tempSpec.Gpumem_GiB = newSpec.Gpumem_GiB
		sqlQuery += ", `gpumem_GiB`='" + strconv.Itoa(int(newSpec.Gpumem_GiB)) + "'"
	}
	if newSpec.Gpu_p2p != "" {
		tempSpec.Gpu_p2p = newSpec.Gpu_p2p
		sqlQuery += ", `gpu_p2p`='" + newSpec.Gpu_p2p + "'"
	}
	if newSpec.OrderInFilteredResult > 0 {
		tempSpec.OrderInFilteredResult = newSpec.OrderInFilteredResult
		sqlQuery += ", `orderInFilteredResult`='" + strconv.Itoa(int(newSpec.OrderInFilteredResult)) + "'"
	}
	if newSpec.EvaluationStatus != "" {
		tempSpec.EvaluationStatus = newSpec.EvaluationStatus
		sqlQuery += ", `evaluationStatus`='" + newSpec.EvaluationStatus + "'"
	}
	if newSpec.EvaluationScore_01 > 0 {
		tempSpec.EvaluationScore_01 = newSpec.EvaluationScore_01
		sqlQuery += ", `evaluationScore_01`='" + fmt.Sprintf("%.6f", newSpec.EvaluationScore_01) + "'"
	}
	if newSpec.EvaluationScore_02 > 0 {
		tempSpec.EvaluationScore_02 = newSpec.EvaluationScore_02
		sqlQuery += ", `evaluationScore_02`='" + fmt.Sprintf("%.6f", newSpec.EvaluationScore_02) + "'"
	}
	if newSpec.EvaluationScore_03 > 0 {
		tempSpec.EvaluationScore_03 = newSpec.EvaluationScore_03
		sqlQuery += ", `evaluationScore_03`='" + fmt.Sprintf("%.6f", newSpec.EvaluationScore_03) + "'"
	}
	if newSpec.EvaluationScore_04 > 0 {
		tempSpec.EvaluationScore_04 = newSpec.EvaluationScore_04
		sqlQuery += ", `evaluationScore_04`='" + fmt.Sprintf("%.6f", newSpec.EvaluationScore_04) + "'"
	}
	if newSpec.EvaluationScore_05 > 0 {
		tempSpec.EvaluationScore_05 = newSpec.EvaluationScore_05
		sqlQuery += ", `evaluationScore_05`='" + fmt.Sprintf("%.6f", newSpec.EvaluationScore_05) + "'"
	}
	if newSpec.EvaluationScore_06 > 0 {
		tempSpec.EvaluationScore_06 = newSpec.EvaluationScore_06
		sqlQuery += ", `evaluationScore_06`='" + fmt.Sprintf("%.6f", newSpec.EvaluationScore_06) + "'"
	}
	if newSpec.EvaluationScore_07 > 0 {
		tempSpec.EvaluationScore_07 = newSpec.EvaluationScore_07
		sqlQuery += ", `evaluationScore_07`='" + fmt.Sprintf("%.6f", newSpec.EvaluationScore_07) + "'"
	}
	if newSpec.EvaluationScore_08 > 0 {
		tempSpec.EvaluationScore_08 = newSpec.EvaluationScore_08
		sqlQuery += ", `evaluationScore_08`='" + fmt.Sprintf("%.6f", newSpec.EvaluationScore_08) + "'"
	}
	if newSpec.EvaluationScore_09 > 0 {
		tempSpec.EvaluationScore_09 = newSpec.EvaluationScore_09
		sqlQuery += ", `evaluationScore_09`='" + fmt.Sprintf("%.6f", newSpec.EvaluationScore_09) + "'"
	}
	if newSpec.EvaluationScore_10 > 0 {
		tempSpec.EvaluationScore_10 = newSpec.EvaluationScore_10
		sqlQuery += ", `evaluationScore_10`='" + fmt.Sprintf("%.6f", newSpec.EvaluationScore_10) + "'"
	}

	sqlQuery += "WHERE `namespace`='" + nsId + "' AND `id`='" + newSpec.Id + "';"

	fmt.Println("sqlQuery: " + sqlQuery)
	// https://stackoverflow.com/questions/42486032/golang-sql-query-syntax-validator
	_, err = sqlparser.Parse(sqlQuery)
	if err != nil {
		temp := TbSpecInfo{}
		return temp, err
	}

	// cb-store
	fmt.Println("=========================== PUT registerSpec")
	Key := common.GenResourceKey(nsId, "spec", tempSpec.Id)
	Val, _ := json.Marshal(tempSpec)
	err = common.CBStore.Put(string(Key), string(Val))
	if err != nil {
		temp := TbSpecInfo{}
		common.CBLog.Error(err)
		return temp, err
	}
	keyValue, _ := common.CBStore.Get(string(Key))
	fmt.Println("<" + keyValue.Key + "> \n" + keyValue.Value)
	fmt.Println("===========================")

	// register information related with MCIS recommendation
	RegisterRecommendList(nsId, tempSpec.ConnectionName, tempSpec.Num_vCPU, tempSpec.Mem_GiB, tempSpec.Storage_GiB, tempSpec.Id, tempSpec.Cost_per_hour)

	stmt, err := common.MYDB.Prepare(sqlQuery)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Data inserted successfully..")
	}

	return tempSpec, nil
}
