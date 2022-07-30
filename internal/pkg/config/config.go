package config

import (
	"flag"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// ServerConfig holds server application's configuration.
type ServerConfig struct {
	Database struct {
		Address string `yaml:"address"`
		Name    string `yaml:"name"`
	} `yaml:"database"`
	Listen struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
	} `yaml:"listen"`
	Salt    string `yaml:"salt"`
	IsDebug bool   `yaml:"is_debug"`
}

// ClientConfig holds client application's configuration.
type ClientConfig struct {
	Server struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
	} `yaml:"server"`
	IsDebug bool `yaml:"is_debug"`
}

var (
	once      sync.Once
	clientCfg *ClientConfig
	serverCfg *ServerConfig
)

// GetServerConfig once parses yaml configuration file
// and then always returns copy of filled server configuration object.
func GetServerConfig() ServerConfig {
	once.Do(func() {
		var configPath string
		flag.StringVar(&configPath, "c", "", "yaml config file")
		flag.Parse()

		serverCfg = &ServerConfig{}
		file, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to open config file")
		}

		if err = yaml.Unmarshal(file, serverCfg); err != nil {
			log.Fatal().Err(err).Msg("unable to unmarshal yaml in config file")
		}
	})

	return *serverCfg
}

// GetClientConfig once parses yaml configuration file
// and then always returns copy of filled client configuration object.
func GetClientConfig() ClientConfig {
	once.Do(func() {
		var configPath string
		flag.StringVar(&configPath, "c", "", "yaml config file")
		flag.Parse()

		clientCfg = &ClientConfig{}
		file, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to open config file")
		}

		if err = yaml.Unmarshal(file, clientCfg); err != nil {
			log.Fatal().Err(err).Msg("unable to unmarshal yaml in config file")
		}
	})

	return *clientCfg
}
