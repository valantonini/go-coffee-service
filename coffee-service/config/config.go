package config

import (
	"log"
	"os"
)

const (
	BindAddress string = "BIND_ADDRESS"
)

// Config defines the service runtime configuration
type Config struct {
	BindAddress string
	Logger      *log.Logger
}

// NewConfigFromEnv creates a Config from the environment variables
func NewConfigFromEnv() (*Config, error) {
	logger := log.Default()

	bindAddress := os.Getenv(BindAddress)
	if bindAddress == "" {
		bindAddress = "localhost:80"
	}

	logger.Printf("Binding to %s", bindAddress)

	return &Config{bindAddress, logger}, nil
}
