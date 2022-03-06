package config

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	bindAddressKey := "BIND_ADDRESS"
	natsAddressKey := "NATS_ADDRESS"

	t.Run("gets bind address when supplied in env", func(t *testing.T) {
		os.Setenv(bindAddressKey, "http://coffee-service:44313")
		defer os.Unsetenv(bindAddressKey)

		config := NewConfigFromEnv("config-test")

		if config.BindAddress != "http://coffee-service:44313" {
			t.Errorf("bind address not set. got %v", config.BindAddress)
		}
	})

	t.Run("defaults to local when not supplied in env", func(t *testing.T) {
		os.Unsetenv(bindAddressKey)

		config := NewConfigFromEnv("config-test")

		if config.BindAddress != "localhost:80" {
			t.Errorf("bind address not defaulting set. got %v", config.BindAddress)
		}
	})

	t.Run("gets nats address when supplied in env", func(t *testing.T) {
		os.Setenv(natsAddressKey, "nats://custom:1234")
		defer os.Unsetenv(natsAddressKey)

		config := NewConfigFromEnv("config-test")

		if config.NatsAddress != "nats://custom:1234" {
			t.Errorf("nats address not set. got %v", config.NatsAddress)
		}
	})

	t.Run("defaults nats address", func(t *testing.T) {
		os.Unsetenv(natsAddressKey)

		config := NewConfigFromEnv("config-test")

		if config.NatsAddress != "nats://nats:4222" {
			t.Errorf("nats address not defaulting. got %v", config.NatsAddress)
		}
	})
}
