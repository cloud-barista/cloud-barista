package tumblebug

// VM의 상태정보
// 원래는 상태정보(VmStatus), 위치정보(Latitude, Longitude) 만 있었으나
// export 기능 추가로 param 추가 됨.
type VmSimpleInfo struct {
	VmIndex  int    `json:"vmIndex"`
	VmID     string `json:"vmID"`
	VmName   string `json:"vmName"`
	VmStatus string `json:"vmStatus"`

	// Latitude  float64 `json:"latitude"`
	// Longitude float64 `json:"longitude"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`

	// export 를 위한 param들 추가
	VmConnectionName   string   `json:"vmConnectionName"`
	VmDescription      string   `json:"vmDescription"`
	VmImageId          string   `json:"vmImageId"`
	VmLabel            string   `json:"vmLabel"`
	VmSecurityGroupIds []string `json:"vmSecurityGroupIds"` //"securityGroupIIds": [		{		  "nameId": "string",		  "systemId": "string"		}	  ],
	VmSpecId           string   `json:"vmSpecId"`
	VmSshKeyId         string   `json:"vmSshKeyId"`
	VmSubnetId         string   `json:"vmSubnetId"`
	VmVnetId           string   `json:"vmVnetId"`
	VmGroupSize        int      `json:"vmGroupSize"` //? 는 없는데.. vmGroupId만 있는데...
	VmUserAccount      string   `json:"vmUserAccount"`
	VmUserPassword     string   `json:"vmUserPassword"`
}
