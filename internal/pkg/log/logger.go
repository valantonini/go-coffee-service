package log

import "github.com/hashicorp/go-hclog"

type Logger interface {
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

func NewLogger(name string) Logger {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       name,
		JSONFormat: false,
		Level:      hclog.Debug,
	})

	return logger
}
