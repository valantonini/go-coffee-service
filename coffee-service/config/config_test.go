package config

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	bindAddressKey := "BIND_ADDRESS"
	natsAddressKey := "NATS_ADDRESS"

	t.Run("gets bind address when supplied in env", func(t *testing.T) {
		bindAddress := "http://coffee-service:44313"
		os.Setenv(bindAddressKey, bindAddress)
		defer os.Unsetenv(bindAddressKey)

		config := NewConfigFromEnv()

		if config.BindAddress != bindAddress {
			t.Errorf("bind address not set. got %v", config.BindAddress)
		}
	})

	t.Run("defaults to local when not supplied in env", func(t *testing.T) {
		bindAddress := "localhost:80"
		os.Unsetenv(bindAddressKey)
		config := NewConfigFromEnv()

		if config.BindAddress != bindAddress {
			t.Errorf("bind address not defaulting set. got %v", config.BindAddress)
		}
	})

	t.Run("gets nats address when supplied in env", func(t *testing.T) {
		natsAddress := "nats://custom:1234"
		os.Setenv(natsAddressKey, natsAddress)
		defer os.Unsetenv(natsAddressKey)

		config := NewConfigFromEnv()

		if config.NatsAddress != natsAddress {
			t.Errorf("nats address not set. got %v", config.NatsAddress)
		}
	})

	t.Run("defaults nats address", func(t *testing.T) {
		natsAddress := "nats://nats:4222"
		os.Unsetenv(natsAddressKey)
		config := NewConfigFromEnv()

		if config.NatsAddress != natsAddress {
			t.Errorf("nats address not defaulting. got %v", config.NatsAddress)
		}
	})
}
