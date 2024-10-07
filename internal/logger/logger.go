package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

func New() (Logger, error) {
	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(zapcore.InfoLevel),
		),
		zap.AddCaller(),
	)

	return logger.Sugar(), nil
}
