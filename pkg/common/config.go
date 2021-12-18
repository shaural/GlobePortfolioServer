package common

import (
	"os"
	"sync"
)

var config *ConfigEnv
var configInit sync.Once

// Config returns single global instance of the configuration object
func Config() *ConfigEnv {
	configInit.Do(func() {
		config = newConfig()
	})
	return config
}

// ConfigEnv environment configuration
type ConfigEnv struct {
	Port string
	DatabaseURL string
}

func newConfig() *ConfigEnv {
	return &ConfigEnv{
		Port: os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}