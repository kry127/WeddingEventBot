package main

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Debug                 bool
	RestartTimeoutSeconds int
}

func LoadConfig() (*Config, error) {
	result := new(Config)

	_, isDebug := os.LookupEnv("DEBUG")
	result.Debug = isDebug

	restartTimeout, isRestarting := os.LookupEnv("RESTART_TIMEOUT")
	if !isRestarting || restartTimeout == "" {
		result.RestartTimeoutSeconds = 0
	} else {
		restartTimeoutInt, err := strconv.Atoi(restartTimeout)
		if err != nil {
			return nil, fmt.Errorf("invalid format of RESTART_TIMEOUT env variable: should be integer (like '5'), got '%v'", restartTimeout)
		}
		result.RestartTimeoutSeconds = restartTimeoutInt
	}

	return result, nil
}
