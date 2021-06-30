package dragonfly

// 알람 이벤트 핸들러
type VmMonitoringAlertEventHandlerInfoReg struct {
	Name string `json:"name"` // 이벤트 핸들러 이름
	Type string `json:"type"` // 네임스페이스 아이디
	//Options EventHandlerOption `json:"options"` //
	Url     string `json:"url"`
	Channel string `json:"channel"`
}
