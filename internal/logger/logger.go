package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/forest-shadow/calendar/internal/config"
)

// Added an alias to the logger from the zap package
// so it can be easily replaced in the future if necessary.
type Logger = *zap.SugaredLogger

var encoderConfig = zapcore.EncoderConfig{
	TimeKey:        "ts",
	LevelKey:       "level",
	NameKey:        "logger",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	CallerKey:      "caller",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

func getLevel(cfg *config.Config) zapcore.Level {
	if cfg.Env == "production" {
		return zapcore.InfoLevel
	}

	return zapcore.DebugLevel
}

func NewLogger(cfg *config.Config) (Logger, error) {
	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(getLevel(cfg)),
		),
		zap.AddCaller(),
	)

	return logger.Sugar(), nil
}
