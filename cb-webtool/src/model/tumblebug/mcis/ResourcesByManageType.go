package mcis

type ResourcesByManageType struct {
	OnCspOnly    ResourceOnCsp       `json:"onCspOnly"`
	OnCspTotal   ResourceOnCsp       `json:"onCspTotal"`
	OnSpider     ResourceOnSpider    `json:"onSpider"`
	OnTumblebug  ResourceOnTumblebug `json:"onTumblebug"`
}