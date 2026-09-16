package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type ServerConfig struct {
	Server struct {
		Port uint16
		Env  ServerEnv
	}

	Database struct {
		Host     string
		Port     uint16
		User     string
		Password string
		Dbname   string
	}
}

func LoadConfigFromFile(filePath string) (*ServerConfig, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var cfg ServerConfig
	if err = toml.Unmarshal(fileContent, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

type ServerEnv string

const (
	EnvLocal ServerEnv = "local"
	EnvProd  ServerEnv = "prod"
)
