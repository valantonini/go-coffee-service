package config

import (
	"log"
	"os"
)

const (
	BindAddress string = "BIND_ADDRESS"
	NatsAddress string = "NATS_ADDRESS"
)

// Config defines the service runtime configuration
type Config struct {
	BindAddress string
	NatsAddress string
	Logger      *log.Logger
}

// NewConfigFromEnv creates a Config from the environment variables with sensible fallbacks
func NewConfigFromEnv() *Config {
	logger := log.Default()

	bindAddress := os.Getenv(BindAddress)
	if bindAddress == "" {
		bindAddress = "localhost:80"
	}

	natsAddress := os.Getenv(NatsAddress)
	if natsAddress == "" {
		natsAddress = "nats://nats:4222"
	}

	return &Config{bindAddress, natsAddress, logger}
}
