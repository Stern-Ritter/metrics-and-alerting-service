package server

type ServerConfig struct {
	URL             string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	StorageFilePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	LoggerLvl       string
}
