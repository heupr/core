package utils

import (
	"encoding/json"
	strftime "github.com/lestrrat/go-strftime"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Configuration struct {
	LogPath          string
	ModelSummaryPath string
	ModelDetailsPath string
	DataCachesPath   string
	SlackUser        string
	SlackChannel     string
	SlackHook        string
}

var initOnceConfig sync.Once
var Config Configuration

func init() {
	initOnceConfig.Do(func() {
		Config = loadConfig()
	})
}

func loadConfig() Configuration {
	file, err := os.Open("./config.json")
	if err != nil {
		if configEnvPath := os.Getenv(ConfigEnv); configEnvPath != "" {
			file, err = os.Open(configEnvPath)
			if err != nil {
				log.Fatal("config error:", err)
			}
		} else {
			gitConfigPath := "$GOPATH/src/coralreefci/analysis/cmd/backtests/bhattacharya/config.json"
			gitConfigPath = expand(gitConfigPath)
			file, err = os.Open(gitConfigPath)
			if err != nil {
				log.Fatal("config error:", err)
			}
		}
	}
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal("config error:", err)
	}
	configuration.expandRelativePaths()
	return configuration
}

func (c *Configuration) expandRelativePaths() {
	c.LogPath = timestamp(expand(c.LogPath))
	c.ModelSummaryPath = timestamp(expand(c.ModelSummaryPath))
	c.ModelDetailsPath = timestamp(expand(c.ModelDetailsPath))
	c.DataCachesPath = expand(c.DataCachesPath)
}

func expand(path string) string {
	if strings.HasPrefix(path, "$GOPATH") {
		goPath := os.Getenv("GOPATH")
		return filepath.Join(goPath, path[7:])
	} else {
		return path
	}
}

func timestamp(path string) string {
	f, err := strftime.New(path)
	if err != nil {
		log.Fatal("config error:", err)
	}
	return f.FormatString(time.Now())
}
