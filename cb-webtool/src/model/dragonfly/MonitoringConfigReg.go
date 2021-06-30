package dragonfly

// 모니터링 정책 설정 Put 할때는 아래 3가지만 보냄
type MonitoringConfigReg struct {
	AgentInterval     int `json:"agent_interval"`     //모니터링 에이전트 수집 주기
	CollectorInterval int `json:"collector_interval"` //모니터링 콜렉터 Aggregate 주기
	MaxHostCount      int `json:"max_host_count"`     // 단일 콜렉터에서 수집하는 최대 가상머신 수
}
