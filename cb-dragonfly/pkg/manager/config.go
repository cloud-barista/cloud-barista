package manager

type Config struct {
	InfluxDB struct {
		EndpointUrl string `yaml:"endpoint_url"`
		Database    string `yaml:"database"`
		UserName    string `yaml:"user_name"`
		Password    string `yaml:"password"`
	} `yaml:"influxdb"`
	Etcd struct {
		EndpointUrl string `yaml:"endpoint_url"`
	} `yaml:"etcd"`
	CollectManager struct {
		CollectorIP   string `yaml:"collector_ip"`
		CollectorPort int    `yaml:"collector_port"`
		CollectorCnt  int    `yaml:"collector_count"`
	} `yaml:"collect_manager"`
	APIServer struct {
		Port int `yaml:"port"`
	} `yaml:"api_server"`
	Monitoring struct {
		AgentInterval     int `yaml:"agent_interval"`
		CollectorInterval int `yaml:"collector_interval"`
		ScheduleInterval  int `yaml:"schedule_interval"`
		MaxHostCount      int `yaml:"max_host_count"`
	} `yaml:"monitoring"`
}

type MonConfig struct {
	AgentInterval      int `json:"agent_interval"`     // 모니터링 에이전트 수집주기
	CollectorInterval  int `json:"collector_interval"` // 모니터링 콜렉터 Aggregate 주기
	SchedulingInterval int `json:"schedule_interval"`  // 모니터링 콜렉터 스케줄링 주기 (스케일 인/아웃 로직 체크 주기)
	MaxHostCount       int `json:"max_host_count"`     // 모니터링 콜렉터 수
}
