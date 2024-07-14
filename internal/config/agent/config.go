package agent

// AgentConfig holds the configuration for the agent.
type AgentConfig struct {
	SendMetricsURL        string `env:"ADDRESS"` // The URL to send metrics statistics to
	SendMetricsEndPoint   string // The endpoint for sending metrics statistics to
	UpdateMetricsInterval int    `env:"POLL_INTERVAL"`   // The interval for updating metrics statistics in seconds
	SendMetricsInterval   int    `env:"REPORT_INTERVAL"` // The interval for sending metrics statistics in seconds
	MetricsBufferSize     int    // The buffer size for metrics channel
	RateLimit             int    `env:"RATE_LIMIT"` // The size of sending metrics statistics worker pool
	SecretKey             string `env:"KEY"`        // The secret key for authentication
	CryptoKeyPath         string `env:"CRYPTO_KEY"` // The path to secret public key for asymmetric encryption
	ConfigFile            string `env:"CONFIG"`     //The path to json config file
	LoggerLvl             string // The logging level
}
