package config

import (
	"github.com/valantonini/go-coffee-service/internal/pkg/log"
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
	Logger      log.Logger
}

// NewConfigFromEnv creates a Config from the environment variables with sensible fallbacks
func NewConfigFromEnv(service string) *Config {
	logger := log.NewLogger(service)

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
