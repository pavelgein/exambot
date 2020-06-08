package config

import (
	"fmt"
	"os"
)

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
