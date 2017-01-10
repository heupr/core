package utils

import (
	log "github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"
	"github.com/rifflock/lfshook"
	"io/ioutil"
	"os"
	"sync"
)

var Log *log.Logger
var ModelSummary *log.Logger
var ModelDetails *log.Logger

var initOnceLog sync.Once

func init() {
	initOnceLog.Do(func() {
		InitLogger()
		InitModelSummary()
		InitModelDetails()
	})
}

func InitLogger() {
	logPath := Config.LogPath
	Log = CreateLogger(logPath)
}

func InitModelSummary() {
	logPath := Config.ModelSummaryPath
	ModelSummary = CreateLogger(logPath)
}

func InitModelDetails() {
	logPath := Config.ModelDetailsPath
	ModelDetails = CreateLogger(logPath)
	ModelDetails.Level = log.DebugLevel
	ModelDetails.Out = ioutil.Discard
}

func CreateLogger(logPath string) *log.Logger {
	logInstance := log.New()
	logInstance.Formatter = new(log.TextFormatter)
	logInstance.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
		log.DebugLevel: logPath,
		log.InfoLevel:  logPath,
		log.WarnLevel:  logPath,
		log.ErrorLevel: logPath,
	}))
	if isProdEnv := os.Getenv(LogEnv); isProdEnv == "PROD" {
		logInstance.Hooks.Add(&slackrus.SlackrusHook{
			HookURL:        Config.SlackHook,
			AcceptedLevels: slackrus.LevelThreshold(log.ErrorLevel),
			Channel:        Config.SlackChannel,
			Username:       Config.SlackUser,
		})
	}
	return logInstance
}
