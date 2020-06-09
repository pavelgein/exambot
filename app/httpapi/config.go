package main

import (
	"time"

	"github.com/pavelgein/exambot/internal/config"
)

type Config struct {
	DB     config.DBConfig
	Server config.HttpServerConfig
}

func CreateConfigFromEnv() Config {
	return Config{
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
