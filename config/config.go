package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	MongoDBConnString     string
	User                  string
	Password              string
	Debug                 bool
	RestartTimeoutSeconds int
}

func LoadConfig() (cfg *Config, err error) {
	defer func() {
		if err != nil {
			err = NewConfigError(err)
		}
	}()
	result := new(Config)

	if mongoDBConnString, ok := os.LookupEnv("MONGODB_CONN_STRING"); !ok {
		return nil, fmt.Errorf("environment variable MONGODB_CONN_STRING should be filled with MongoDB connection string. Example: mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.2.10&tls=false")
	} else {
		result.MongoDBConnString = mongoDBConnString
	}

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

type ConfigError struct {
	error
}

func NewConfigError(err error) *ConfigError {
	return &ConfigError{error: err}
}
