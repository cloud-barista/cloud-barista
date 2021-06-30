package dragonfly

// 알람 이벤트 핸들러
// type이 slack이면 option항목에 url, channel
type VmMonitoringAlertEventHandlerSlackInfo struct {
	EventHandlerID string                  `json:"id"`      // 네임스페이스 아이디
	Name           string                  `json:"name"`    // 이벤트 핸들러 이름
	Type           string                  `json:"type"`    // 이벤트 핸들러 유형 ( "slack" | "smtp" )
	Options        EventHandlerOptionSlack `json:"options"` //
}

type EventHandlerOptionSlack struct {
	Url     string `json:"url"`     // slack hook URL
	Channel string `json:"channel"` // slack 알람 채널 (#kapacitor-alert)
}
