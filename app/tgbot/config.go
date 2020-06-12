package main

import (
	"fmt"
	"os"

	"github.com/pavelgein/exambot/internal/config"
	"github.com/pavelgein/exambot/internal/db"
)

type Config struct {
	DBConfig config.DBConfig
	Token    string
}

func GetEnv(variable string) string {
	res := os.Getenv(variable)
	if res == "" {
		panic(fmt.Sprintf("Variale %s is not set", variable))
	}

	return res
}

func GetEnvWithDefault(variable string, dflt string) string {
	res := os.Getenv(variable)
	if res == "" {
		return dflt
	}

	return res
}

func MakeConfigFromEnvironemnt() *Config {
	var config Config
	config.Token = GetEnv("EXAMBOT_TOKEN")
	config.DBConfig = db.CreateConfigFromEnvironment()
	return &config
}
