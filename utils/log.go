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
		CreateAppLogger()
		CreateModelLogger()
	})
}

func CreateAppLogger() {
	logPath := Config.LogPath
	AppLog = CreateLogger(logPath)
}

func CreateModelLogger() {
	logPath := Config.ModelSummaryPath
	ModelLog = CreateLogger(logPath)
}

func CreateLogger(logPath string) *zap.Logger {
	logConfig := zap.NewProductionConfig()
	logConfig.OutputPaths = []string{logPath}
	logger, err := logConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}
