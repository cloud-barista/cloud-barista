package mcir

import (
	"encoding/json"
	"fmt"
	"net/http"

	//"strings"

	"github.com/cloud-barista/cb-tumblebug/src/core/common"
	"github.com/labstack/echo/v4"
)

type subnetReq struct {
	//Id                 string `json:"id"`
	Name               string `json:"name"`
	ConnectionName     string `json:"connectionName"`
	CspSubnetId        string `json:"cspSubnetId"`
	CspSubnetName      string `json:"cspSubnetName"`
	VirtualNetworkId   string `json:"virtualNetworkId"`
	VirtualNetworkName string `json:"virtualNetworkName"`
	CidrBlock          string `json:"cidrBlock"`
	Region             string `json:"region"`
	ResourceGroupName  string `json:"resourceGroupName"`
	Description        string `json:"description"`
}

type subnetInfo struct {
	Id                 string `json:"id"`
	Name               string `json:"name"`
	ConnectionName     string `json:"connectionName"`
	CspSubnetId        string `json:"cspSubnetId"`
	CspSubnetName      string `json:"cspSubnetName"`
	VirtualNetworkId   string `json:"virtualNetworkId"`
	VirtualNetworkName string `json:"virtualNetworkName"`
	CidrBlock          string `json:"cidrBlock"`
	Region             string `json:"region"`
	ResourceGroupName  string `json:"resourceGroupName"`
	Description        string `json:"description"`
}

// MCIS API Proxy: Subnet
func RestPostSubnet(c echo.Context) error {

	nsId := c.Param("nsId")

	u := &subnetReq{}
	if err := c.Bind(u); err != nil {
		return err
	}

	/*
		action := c.QueryParam("action")
		fmt.Println("[POST Subnet requested action: " + action)
		if action == "create" {
			fmt.Println("[Creating Subnet]")
			content, _ := createSubnet(nsId, u)
			return c.JSON(http.StatusCreated, content)

		} else if action == "register" {
			fmt.Println("[Registering Subnet]")
			content, _ := registerSubnet(nsId, u)
			return c.JSON(http.StatusCreated, content)

		} else {
			mapA := map[string]string{"message": "You must specify: action=create or action=register"}
			return c.JSON(http.StatusFailedDependency, &mapA)
		}
	*/

	fmt.Println("[POST Subnet")
	fmt.Println("[Creating Subnet]")
	content, err := createSubnet(nsId, u)
	if err != nil {
		common.CBLog.Error(err)
		mapA := map[string]string{
			"message": "Failed to create a Subnet"}
		return c.JSON(http.StatusInternalServerError, &mapA)
	}
	return c.JSON(http.StatusCreated, content)
}

func RestGetSubnet(c echo.Context) error {

	nsId := c.Param("nsId")

	id := c.Param("subnetId")

	content := subnetInfo{}

	fmt.Println("[Get subnet for id]" + id)
	key := common.GenResourceKey(nsId, "subnet", id)
	fmt.Println(key)

	keyValue, _ := common.CBStore.Get(key)
	if keyValue == nil {
		mapA := map[string]string{"message": "Failed to find the subnet with given ID."}
		return c.JSON(http.StatusNotFound, &mapA)
	} else {
		fmt.Println("<" + keyValue.Key + "> \n" + keyValue.Value)
		fmt.Println("===============================================")

		json.Unmarshal([]byte(keyValue.Value), &content)
		content.Id = id // Optional. Can be omitted.

		return c.JSON(http.StatusOK, &content)
	}
}

func RestGetAllSubnet(c echo.Context) error {

	nsId := c.Param("nsId")

	var content struct {
		//Name string     `json:"name"`
		Subnet []subnetInfo `json:"subnet"`
	}

	subnetList := ListResourceId(nsId, "subnet")

	for _, v := range subnetList {

		key := common.GenResourceKey(nsId, "subnet", v)
		fmt.Println(key)
		keyValue, _ := common.CBStore.Get(key)
		fmt.Println("<" + keyValue.Key + "> \n" + keyValue.Value)
		subnetTmp := subnetInfo{}
		json.Unmarshal([]byte(keyValue.Value), &subnetTmp)
		subnetTmp.Id = v
		content.Subnet = append(content.Subnet, subnetTmp)

	}
	fmt.Printf("content %+v\n", content)

	return c.JSON(http.StatusOK, &content)

}

func RestPutSubnet(c echo.Context) error {
	//nsId := c.Param("nsId")

	return nil
}

func RestDelSubnet(c echo.Context) error {

	nsId := c.Param("nsId")
	id := c.Param("subnetId")
	forceFlag := c.QueryParam("force")

	//responseCode, _, err := delSubnet(nsId, id, forceFlag)

	responseCode, _, err := DelResource(nsId, "subnet", id, forceFlag)
	if err != nil {
		common.CBLog.Error(err)
		mapA := map[string]string{"message": "Failed to delete the subnet"}
		return c.JSON(responseCode, &mapA)
	}

	mapA := map[string]string{"message": "The subnet has been deleted"}
	return c.JSON(http.StatusOK, &mapA)
}

func RestDelAllSubnet(c echo.Context) error {

	nsId := c.Param("nsId")
	forceFlag := c.QueryParam("force")

	subnetList := ListResourceId(nsId, "subnet")

	if len(subnetList) == 0 {
		mapA := map[string]string{"message": "There is no subnet element in this namespace."}
		return c.JSON(http.StatusNotFound, &mapA)
	} else {
		for _, v := range subnetList {
			//responseCode, _, err := delSubnet(nsId, v, forceFlag)

			responseCode, _, err := DelResource(nsId, "subnet", v, forceFlag)
			if err != nil {
				common.CBLog.Error(err)
				mapA := map[string]string{"message": "Failed to delete the subnet"}
				return c.JSON(responseCode, &mapA)
			}

		}

		mapA := map[string]string{"message": "All subnets has been deleted"}
		return c.JSON(http.StatusOK, &mapA)
	}
}

func createSubnet(nsId string, u *subnetReq) (subnetInfo, error) {
	check, _ := CheckResource(nsId, "subnet", u.Name)

	if check {
		temp := subnetInfo{}
		err := fmt.Errorf("The subnet " + u.Name + " already exists.")
		return temp, err
	}

	/* FYI; as of 2020-04-17
	type subnetReq struct {
		//Id                 string `json:"id"`
		Name               string `json:"name"`
		ConnectionName     string `json:"connectionName"`
		CspSubnetId        string `json:"cspSubnetId"`
		CspSubnetName      string `json:"cspSubnetName"`
		VirtualNetworkId   string `json:"virtualNetworkId"`
		VirtualNetworkName string `json:"virtualNetworkName"`
		CidrBlock          string `json:"cidrBlock"`
		Region             string `json:"region"`
		ResourceGroupName  string `json:"resourceGroupName"`
		Description        string `json:"description"`
	}
	*/

	//url := common.SPIDER_REST_URL + "/subnet"

	// TODO: To be implemented

	content := subnetInfo{}
	content.Id = common.GenId(u.Name)
	content.Name = common.GenId(u.Name)
	content.ConnectionName = u.ConnectionName
	content.CspSubnetId = u.CspSubnetId
	content.CspSubnetName = u.CspSubnetName
	content.VirtualNetworkId = u.VirtualNetworkId
	content.VirtualNetworkName = u.VirtualNetworkName
	content.CidrBlock = u.CidrBlock
	content.Region = u.Region
	content.ResourceGroupName = u.ResourceGroupName
	content.Description = u.Description

	// cb-store
	fmt.Println("=========================== PUT createSubnet")
	Key := common.GenResourceKey(nsId, "subnet", content.Id)
	mapA := map[string]string{
		"connectionName":     content.ConnectionName,
		"cspSubnetId":        content.CspSubnetId,
		"cspSubnetName":      content.CspSubnetName,
		"virtualNetworkId":   content.VirtualNetworkId,
		"virtualNetworkName": content.VirtualNetworkName,
		"cidrBlock":          content.CidrBlock,
		"region":             content.Region,
		"resourceGroupName":  content.ResourceGroupName,
		"description":        content.Description}
	Val, _ := json.Marshal(mapA)
	//fmt.Println("Key: ", Key)
	//fmt.Println("Val: ", Val)
	err := common.CBStore.Put(string(Key), string(Val))
	if err != nil {
		common.CBLog.Error(err)
		return content, err
	}
	keyValue, _ := common.CBStore.Get(string(Key))
	fmt.Println("<" + keyValue.Key + "> \n" + keyValue.Value)
	fmt.Println("===========================")
	return content, nil
}

/*
func registerSubnet(nsId string, u *subnetReq) (subnetInfo, error) {

	content := subnetInfo{}
	content.Id = common.GenId(u.Name)
	content.ConnectionName = u.ConnectionName
	content.CspSubnetId = u.CspSubnetId
	content.CspSubnetName = u.CspSubnetName
	content.VirtualNetworkId = u.VirtualNetworkId
	content.VirtualNetworkName = u.VirtualNetworkName
	content.CidrBlock = u.CidrBlock
	content.Region = u.Region
	content.ResourceGroupName = u.ResourceGroupName
	content.Description = u.Description

	// TODO here: implement the logic
	// - Fetch the subnet info from CSP.

	// cb-store
	fmt.Println("=========================== PUT registerSubnet")
	Key := genResourceKey(nsId, "subnet", content.Id)
	mapA := map[string]string{
		"connectionName":     content.ConnectionName,
		"cspSubnetId":        content.CspSubnetId,
		"cspSubnetName":      content.CspSubnetName,
		"virtualNetworkId":   content.VirtualNetworkId,
		"virtualNetworkName": content.VirtualNetworkName,
		"cidrBlock":          content.CidrBlock,
		"region":             content.Region,
		"resourceGroupName":  content.ResourceGroupName,
		"description":        content.Description}
	Val, _ := json.Marshal(mapA)
	err := common.CBStore.Put(string(Key), string(Val))
	if err != nil {
		common.CBLog.Error(err)
		return content, err
	}
	keyValue, _ := common.CBStore.Get(string(Key))
	fmt.Println("<" + keyValue.Key + "> \n" + keyValue.Value)
	fmt.Println("===========================")
	return content, nil
}
*/

/*
func getSubnetList(nsId string) []string {

	fmt.Println("[Get subnets")
	key := "/ns/" + nsId + "/resources/subnet"
	fmt.Println(key)

	keyValue, _ := common.CBStore.GetList(key, true)
	keyValue = cbstore_utils.GetChildList(keyValue, key)

	var subnetList []string
	for _, v := range keyValue {
		//if !strings.Contains(v.Key, "vm") {
		subnetList = append(subnetList, strings.TrimPrefix(v.Key, "/ns/"+nsId+"/resources/subnet/"))
		//}
	}
	for _, v := range subnetList {
		fmt.Println("<" + v + "> \n")
	}
	fmt.Println("===============================================")
	return subnetList

}
*/

/*
func delSubnet(nsId string, Id string, forceFlag string) (int, []byte, error) {

	fmt.Println("[Delete subnet] " + Id)

	key := genResourceKey(nsId, "subnet", Id)
	fmt.Println(key)

	// delete mcis info
	err := common.CBStore.Delete(key)
	if err != nil {
		common.CBLog.Error(err)
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, nil
}
*/
