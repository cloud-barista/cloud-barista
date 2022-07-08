package mcis

type ResourceCountOverview struct {
	OnCspOnly    int  `json:"onCspOnly"`
	OnCspTotal   int  `json:"onCspTotal"`
	OnSpider     int  `json:"onSpider"`
	OnTumblebug  int  `json:"onTumblebug"`
}