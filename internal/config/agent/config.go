package agent

type AgentConfig struct {
	SendMetricsURL        string `env:"ADDRESS"`
	SendMetricsEndPoint   string
	UpdateMetricsInterval int `env:"POLL_INTERVAL"`
	SendMetricsInterval   int `env:"REPORT_INTERVAL"`
	LoggerLvl             string
}
