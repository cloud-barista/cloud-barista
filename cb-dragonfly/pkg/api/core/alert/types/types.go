package types

type AlertTaskReq struct {
	Name string `json:"name"`

	Measurement string `json:"measurement"`

	TargetType string `json:"target_type"`
	TargetId   string `json:"target_id"`

	EventDuration string `json:"event_duration"`
	Metric        string `json:"metric"`

	AlertMathExpression string  `json:"alert_math_expression"`
	AlertThreshold      float64 `json:"alert_threshold"`

	WarnEventCnt   int64 `json:"warn_event_cnt"`
	CriticEventCnt int64 `json:"critic_event_cnt"`

	AlertEventType    string `json:"alert_event_type"`
	AlertEventName    string `json:"alert_event_name"`
	AlertEventMessage string `json:"alert_event_message"`

	AlertPostUrl string `json:"alert_post_url"`
}

type AlertTask struct {
	Name string `json:"name"`

	Measurement string `json:"measurement"`

	TargetType string `json:"target_type"`
	TargetId   string `json:"target_id"`

	EventDuration string `json:"event_duration"`
	Metric        string `json:"metric"`

	AlertMathExpression string  `json:"alert_math_expression"`
	AlertThreshold      float64 `json:"alert_threshold"`

	WarnEventCnt   int64 `json:"warn_event_cnt"`
	CriticEventCnt int64 `json:"critic_event_cnt"`

	AlertEventType    string `json:"alert_event_type,omitempty"`
	AlertEventName    string `json:"alert_event_name,omitempty"`
	AlertEventMessage string `json:"alert_event_message,omitempty"`

	AlertPostUrl string `json:"alert_post_url,omitempty"`
}

type AlertEventHandlerReq struct {
	Name string `json:"name"`
	Type string `json:"type"`

	// Parameters for Slack
	Url       string `json:"url,omitempty"`
	Workspace string `json:"workspace,omitempty"`
	Channel   string `json:"channel,omitempty"`

	// Parameters for SMTP
	Host     string   `json:"host,omitempty"`
	Port     int      `json:"port,omitempty"`
	From     string   `json:"from,omitempty"`
	To       []string `json:"to,omitempty"`
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
}

type AlertEventHandler struct {
	ID      string                 `json:"id"`
	Name    string                 `json:"name"`
	Type    string                 `json:"type"`
	Options map[string]interface{} `json:"options"`
}

type AlertEventLog struct {
	Id      string `json:"id"`
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"message"`
}
