package zap

import (
	"github.com/octokerbs/chronocode-backend/internal/domain"
	"go.uber.org/zap"
)

type Logger struct {
	zapLogger *ZapLogger
}

func NewLogger() (*Logger, error) {
	logger, err := NewZapLogger()
	if err != nil {
		return nil, err
	}
	return &Logger{zapLogger: logger}, nil
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.zapLogger.sugaredLogger.Infow(msg, keysAndValues...)
}

func (l *Logger) Warn(msg string, err error, keysAndValues ...interface{}) {
	allArgs := append(keysAndValues, "error", err)
	l.zapLogger.sugaredLogger.Warnw(msg, allArgs...)
}

func (l *Logger) Error(msg string, err error, keysAndValues ...interface{}) {
	allArgs := append(keysAndValues, "error", err)
	l.zapLogger.sugaredLogger.Errorw(msg, allArgs...)
}

func (l *Logger) Fatal(msg string, err error, keysAndValues ...interface{}) {
	allArgs := append(keysAndValues, "error", err)
	l.zapLogger.sugaredLogger.Fatalw(msg, allArgs...)
}

func (l *Logger) With(keysAndValues ...interface{}) domain.Logger {
	return &Logger{zapLogger: l.zapLogger.With(keysAndValues...)}
}

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

func (zl *ZapLogger) With(keysAndValues ...interface{}) *ZapLogger {
	newSugaredLogger := zl.sugaredLogger.With(keysAndValues...)
	return &ZapLogger{sugaredLogger: newSugaredLogger}
}
