package dragonfly

// 알람 이벤트 핸들러
// type이 slack이면 option항목에 url, channel
type VmMonitoringAlertEventHandlerSmtpInfo struct {
	EventHandlerID string                 `json:"id"`      // 네임스페이스 아이디
	Name           string                 `json:"name"`    // 이벤트 핸들러 이름
	Type           string                 `json:"type"`    // 이벤트 핸들러 유형 ( "slack" | "smtp" )
	Options        EventHandlerOptionSmtp `json:"options"` //
}

type EventHandlerOptionSmtp struct {
	Host       string   `json:"host"`     // 서버 ex)localhost
	Port       string   `json:"port"`     //
	Username   string   `json:"username"` //
	Password   string   `json:"password"` //
	FromDomain string   `json:"from"`     //
	ToDomains  []string `json:"to"`       //
}
