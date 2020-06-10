package main

import (
	"os"
	"time"

	"github.com/pavelgein/exambot/internal/config"
)

type Config struct {
	DB     config.DBConfig
	Server config.HttpServerConfig
	Salt   string
}

func CreateConfigFromEnv() Config {
	return Config{
		Salt: os.Getenv("EXAMBOT_SALT"),
		DB: config.DBConfig{
			Dialect:          config.GetEnvWithDefault("EXAMBOT_DB_DIALECT", "sqlite3"),
			ConnectionParams: config.GetEnvWithDefault("EXAMBOT_CONN_PARAMS", "test.db"),
		},
		Server: config.HttpServerConfig{
			Address:      config.GetEnvWithDefault("EXAMBOT_HTTP_ADDRESS", "0.0.0.0:12345"),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
	}
}
