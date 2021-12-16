package mcir

import (
	tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
)

type SpiderImageInfo struct {
	GuestOS      string                `json:"guestOS"`
	IID          tbcommon.TbIID        `json:"iid"`
	KeyValueList []tbcommon.TbKeyValue `json:"keyValueList"`
	Name         string                `json:"name"`
	Status       string                `json:"status"`
}

type SpiderImageInfos []SpiderImageInfo
