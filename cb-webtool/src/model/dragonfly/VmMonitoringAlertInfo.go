package dragonfly

//
type VmMonitoringAlertInfo struct {
	AlertName           string `json:"name"`                  // 알람 이름
	Measurement         string `json:"measurement"`           // 알람 메트릭 유형 ( "cpu" | "mem" | "disk" )
	TargetType          string `json:"target_type"`           // 알람 임계치 타겟 유형 ( "NS" | "MCIS" | "VM" )
	TargetID            string `json:"target_id"`             // 알람 임계치 타겟 아이디
	EventDuration       string `json:"event_duration"`        // 알람 주기
	Metric              string `json:"metric"`                // 알람 메트릭
	AlertMathExpression string `json:"alert_math_expression"` // 알람 임계치 정규식 ( "equal" | "greater" | "equalgreater" | "less" | "equalless" )
	AlertThreshold      int    `json:"alert_threshold"`       // 알람 임계치 값
	WarnEventCnt        int    `json:"warn_event_cnt"`        // 알람 임계치 횟수 (warn)
	CriticEventCnt      int    `json:"critic_event_cnt"`      // 알람 임계치 횟수 (critic)
	AlertEventType      string `json:"alert_event_type"`      // 이벤트 핸들러 유형 ( "slack" | "smtp" )
	AlertEventName      string `json:"alert_event_name"`      // 이벤트 핸들러 이름
	AlertEventMessage   string `json:"alert_event_message"`   // 이벤트 핸들러 메세지
}
