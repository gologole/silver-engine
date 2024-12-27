package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerKey is the key for storing and retrieving the logger from context.
const LoggerKey = "logger"

// NewLogger creates a new zap.Logger with custom configuration.
func NewLogger() *zap.Logger {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

// WithLoggerContext adds a zap.Logger to the provided context.
func WithLoggerContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

// FromContext retrieves the zap.Logger from the context.
// If no logger is present, it returns a default logger.
func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(LoggerKey).(*zap.Logger)
	if !ok || logger == nil {
		return zap.L().With(zap.String("fallback", "default_logger"))
	}
	return logger
}

// Usage Guide:
// 1. Create a logger instance:
//    logger := logger.NewLogger()
//
// 2. Add the logger to the context:
//    ctx := logger.WithLoggerContext(context.Background(), logger)
//
// 3. Retrieve the logger from context in any layer:
//    log := logger.FromContext(ctx)
//
// 4. Log messages:
//    log.Info("This is an info message")
//    log.Debug("This is a debug message")
//    log.Error("This is an error message")
//
// 5. The logger includes caller information to trace where the logging method was invoked.
