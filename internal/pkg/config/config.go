package config

// ServerConfig holds server application's configuration.
type ServerConfig struct {
	Address        string `yaml:"address"`
	Port           int    `yaml:"port"`
	DataSourceName string `yaml:"dsn"`
}

// ClientConfig holds client application's configuration.
type ClientConfig struct {
	Server struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
	} `yaml:"server"`
}
