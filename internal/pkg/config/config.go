package config

// ServerConfig holds server application's configuration.
type ServerConfig struct {
	Address        string `yaml:"address"`
	Port           int    `yaml:"port"`
	DataSourceName string `yaml:"dsn"`
	IsDebug        bool   `yaml:"is_debug"`
}

// ClientConfig holds client application's configuration.
type ClientConfig struct {
	Server struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
	} `yaml:"server"`
	IsDebug bool `yaml:"is_debug"`
}
