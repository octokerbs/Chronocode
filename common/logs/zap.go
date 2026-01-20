package logs

import "go.uber.org/zap"

func Init() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return logger
}
