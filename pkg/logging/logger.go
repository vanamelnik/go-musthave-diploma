package logging

import (
	"os"

	"github.com/rs/zerolog"
)

// LoggerOption customizes the logger.
type LoggerOption func(zerolog.Logger) zerolog.Logger

// WithLevel sets log level.
func WithLevel(levelStr string) LoggerOption {
	return func(logger zerolog.Logger) zerolog.Logger {
		var level zerolog.Level
		switch levelStr {
		case "trace":
			level = zerolog.TraceLevel
		case "debug":
			level = zerolog.DebugLevel
		case "info":
			level = zerolog.InfoLevel
		case "warn":
			level = zerolog.WarnLevel
		case "error":
			level = zerolog.ErrorLevel
		}
		return logger.Level(level)
	}
}

// WithConsoleOutput sets colorized output to the console.
func WithConsoleOutput(console bool) LoggerOption {
	return func(logger zerolog.Logger) zerolog.Logger {
		if !console {
			return logger
		}
		return logger.Output(zerolog.NewConsoleWriter())
	}
}

// NewLogger creates a new zerolog.Logger with provided options.
func NewLogger(opts ...LoggerOption) zerolog.Logger {
	logger := zerolog.New(os.Stdout).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Logger()
	for _, opt := range opts {
		logger = opt(logger)
	}

	return logger
}
