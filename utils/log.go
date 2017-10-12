package utils

import (
	"log"
	"sync"

	"github.com/johntdyer/slackrus"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

var AppLog *zap.Logger
var ModelLog *zap.Logger
var SlackLog *logrus.Logger

var initOnceLog sync.Once

func init() {
	initOnceLog.Do(func() {
		AppLog = IntializeLog(Config.AppLogPath)
		ModelLog = IntializeLog(Config.ModelLogPath)
		SlackLog = InitializeSlackLog()
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

func InitializeSlackLog() *logrus.Logger {
	logInstance := logrus.New()
	logInstance.Formatter = new(logrus.TextFormatter)
	logInstance.Hooks.Add(&slackrus.SlackrusHook{
		HookURL:        "https://hooks.slack.com/services/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",
		AcceptedLevels: slackrus.LevelThreshold(logrus.InfoLevel),
		Channel:        "#status",
		Username:       "status-update",
	})
	return logInstance
}
