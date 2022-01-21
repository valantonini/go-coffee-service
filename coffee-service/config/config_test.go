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

		config, err := NewFromEnv()

		if err != nil {
			t.Error(err)
		}

		if config.BindAddress != bindAddress {
			t.Errorf("bind address not set. got %v", config.BindAddress)
		}
	})

	t.Run("defaults to local when not supplied in env", func(t *testing.T) {
		bindAddress := "localhost:80"
		config, err := NewFromEnv()

		if err != nil {
			t.Error(err)
		}

		if config.BindAddress != bindAddress {
			t.Errorf("bind address not set. got %v", config.BindAddress)
		}
	})

	t.Run("gets a default logger", func(t *testing.T) {
		config, err := NewFromEnv()

		if err != nil {
			t.Error(err)
		}

		if config.Logger == nil {
			t.Error("no logger set")
		}
	})
}
