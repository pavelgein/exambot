package config

import (
	"time"
)

type HttpServerConfig struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func CreateHttpServerConfig() HttpServerConfig {
	return HttpServerConfig{
		Address:      GetEnvWithDefault("EXAMBOT_HTTP_ADDRESS", "0.0.0.0:12345"),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
}
