package config

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("gets bind address when supplied in env", func(t *testing.T) {
		bindAddressKey := "BIND_ADDRESS"
		bindAddress := "http://coffee-service:44313"
		os.Setenv(bindAddressKey, bindAddress)
		defer os.Unsetenv(bindAddressKey)

		config, err := NewConfigFromEnv()

		if err != nil {
			t.Error(err)
		}

		if config.BindAddress != bindAddress {
			t.Errorf("bind address not set. got %v", config.BindAddress)
		}
	})

	t.Run("defaults to local when not supplied in env", func(t *testing.T) {
		bindAddress := "localhost:80"
		config, err := NewConfigFromEnv()

		if err != nil {
			t.Error(err)
		}

		if config.BindAddress != bindAddress {
			t.Errorf("bind address defaulting set. got %v", config.BindAddress)
		}
	})

	t.Run("gets nats address when supplied in env", func(t *testing.T) {
		natsAddressKey := "NATS_ADDRESS"
		natsAddress := "nats://custom:1234"
		os.Setenv(natsAddressKey, natsAddress)
		defer os.Unsetenv(natsAddressKey)

		config, err := NewConfigFromEnv()

		if err != nil {
			t.Error(err)
		}

		if config.NatsAddress != natsAddress {
			t.Errorf("nats address not set. got %v", config.NatsAddress)
		}
	})

	t.Run("defaults nats address", func(t *testing.T) {
		natsAddress := "nats://nats:4222"
		config, err := NewConfigFromEnv()

		if err != nil {
			t.Error(err)
		}

		if config.NatsAddress != natsAddress {
			t.Errorf("nats address not defaulting. got %v", config.NatsAddress)
		}
	})

	t.Run("gets a default logger", func(t *testing.T) {
		config, err := NewConfigFromEnv()

		if err != nil {
			t.Error(err)
		}

		if config.Logger == nil {
			t.Error("no logger set")
		}
	})
}
