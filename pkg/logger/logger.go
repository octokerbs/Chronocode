package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init(env string) error {
	var err error

	if env == "production" {
		Log, err = zap.NewProduction()
	} else {
		Log, err = zap.NewDevelopment()
	}

	if err != nil {
		return err
	}

	zap.ReplaceGlobals(Log)

	return nil
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
