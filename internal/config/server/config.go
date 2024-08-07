package server

// ServerConfig holds the configuration for the server.
type ServerConfig struct {
	URL             string `env:"ADDRESS"`           // The address and port to run the server
	StoreInterval   int    `env:"STORE_INTERVAL"`    // The interval to store metrics statistics to file in seconds
	FileStoragePath string `env:"FILE_STORAGE_PATH"` // The file path to store metrics statistics
	Restore         bool   `env:"RESTORE"`           // Need to restore metrics statistics from the file when running the server
	DatabaseDSN     string `env:"DATABASE_DSN"`      // The database DSN
	GRPC            bool   `env:"GRPC"`              // The grpc usage flag
	SecretKey       string `env:"KEY"`               // The secret key for authentication
	CryptoKeyPath   string `env:"CRYPTO_KEY"`        // The path to secret private key for asymmetric encryption
	TLSCertPath     string `env:"TLS_CERT"`          // The path to TLS certificate
	TLSKeyPath      string `env:"TLS_KEY"`           // The path to TLS key
	TrustedSubnet   string `env:"TRUSTED_SUBNET"`    // The trusted subnet for agents
	ConfigFile      string `env:"CONFIG"`            // The path to json config file
	ShutdownTimeout int    // The server shutdown timeout in seconds
	LoggerLvl       string // The logging level
}
