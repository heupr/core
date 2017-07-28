package utils

import (
	"go.uber.org/zap"
	"log"
	"sync"
)

var AppLog *zap.Logger
var ModelLog *zap.Logger

var initOnceLog sync.Once

func init() {
	initOnceLog.Do(func() {
		AppLog = IntializeLog(Config.AppLogPath)
		ModelLog = IntializeLog(Config.ModelLogPath)
	})
}

func IntializeLog(logPath string) *zap.Logger {
	logConfig := zap.NewProductionConfig()
	logConfig.OutputPaths = []string{logPath}
	logConfig.Sampling = nil
	logger, err := logConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}
