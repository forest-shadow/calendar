package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/forest-shadow/calendar/internal/config"
)

func NewNotifierLogger(cfg *config.Config) (Logger, error) {
	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(os.Stdout),
				getWriteSyncer(cfg.Notifier.LogPath),
			),
			zap.NewAtomicLevelAt(getLevel(cfg)),
		),

		zap.AddCaller(),
	)

	return logger.Sugar().With("component", "notifier_job"), nil
}

type WriteSyncer struct {
	io.Writer
}

func (ws WriteSyncer) Sync() error {
	return nil
}

func getWriteSyncer(logName string) zapcore.WriteSyncer {
	ioWriter := &lumberjack.Logger{
		Filename:   logName,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		LocalTime:  true,
		Compress:   false,
	}
	sw := WriteSyncer{
		ioWriter,
	}
	return sw
}
