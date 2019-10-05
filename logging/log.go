package logging

import "go.uber.org/zap"

var logger *zap.Logger

func init() {
	var err error

	logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	logger.Info("initialize logger is successful")
}

func Log() *zap.Logger {
	return logger
}
