package dragonfly

// 알람 이벤트 핸들러
type VmMonitoringAlertEventHandlerInfo struct {
	EventHandlerID string             `json:"id"`      // 네임스페이스 아이디
	Name           string             `json:"name"`    // 이벤트 핸들러 이름
	Type           string             `json:"type"`    // 네임스페이스 아이디
	Options        EventHandlerOption `json:"options"` //

}

type EventHandlerOption struct {
	Url     string `json:"url"`
	Channel string `json:"channel"`
}
