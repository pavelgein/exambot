package main

import (
	"fmt"
	"os"
)

type Config struct {
	DBDialect          string
	DBConnectionParams string
	Token              string
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
	config.DBDialect = GetEnvWithDefault("EXAMBOT_DB_DIALECT", "sqlite3")
	config.DBConnectionParams = GetEnvWithDefault("EXAMBOT_CONN_PARAMS", "test.db")
	return &config
}
