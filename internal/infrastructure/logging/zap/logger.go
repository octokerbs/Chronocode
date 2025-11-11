package zap

import (
	"github.com/octokerbs/chronocode-backend/internal/log"
	"go.uber.org/zap"
)

type ZapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

func NewZapLogger() (*ZapLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	sugaredLogger := logger.Sugar()

	return &ZapLogger{sugaredLogger: sugaredLogger}, nil
}

func (zl *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	zl.sugaredLogger.Infow(msg, keysAndValues...)
}

func (zl *ZapLogger) Warn(msg string, err error, keysAndValues ...interface{}) {
	allArgs := append(keysAndValues, "error", err)
	zl.sugaredLogger.Warnw(msg, allArgs...)
}

func (zl *ZapLogger) Error(msg string, err error, keysAndValues ...interface{}) {
	allArgs := append(keysAndValues, "error", err)
	zl.sugaredLogger.Errorw(msg, allArgs...)
}

func (zl *ZapLogger) Fatal(msg string, err error, keysAndValues ...interface{}) {
	allArgs := append(keysAndValues, "error", err)
	zl.sugaredLogger.Fatalw(msg, allArgs...)
}

func (zl *ZapLogger) With(keysAndValues ...interface{}) log.Logger {
	newSugaredLogger := zl.sugaredLogger.With(keysAndValues...)
	return &ZapLogger{sugaredLogger: newSugaredLogger}
}
