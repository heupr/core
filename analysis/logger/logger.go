package logger

import (
  "coralreefci/analysis/config"
  log "github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"
	"github.com/rifflock/lfshook"
	"sync"
)

var Log *log.Logger
var ModelSummary *log.Logger

var initOnce sync.Once
var conf config.Config

//TODO: Relocate this package
//TODO: Decide how we want to make ModelSummary BackTest friendly
func init() {
	initOnce.Do(func() {
    conf = config.Load()
		InitLogger()
    InitModelSummary()
	})
}

func InitLogger() {
  logPath := conf.LogPath
  Log = CreateLogger(logPath)
}

func InitModelSummary() {
  logPath := conf.ModelSummaryPath
  ModelSummary = CreateLogger(logPath)
}

//TODO: pull out slack settings into config file
func CreateLogger(logPath string) *log.Logger {
  logInstance := log.New()
	logInstance.Formatter = new(log.TextFormatter)
	logInstance.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
		log.InfoLevel:  logPath,
		log.ErrorLevel: logPath,
	}))
	logInstance.Hooks.Add(&slackrus.SlackrusHook{
		HookURL:        "https://hooks.slack.com/services/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",
		AcceptedLevels: slackrus.LevelThreshold(log.ErrorLevel),
		Channel:        "#random",
		Username:       "coralreef-bot",
	})
  return logInstance
}
