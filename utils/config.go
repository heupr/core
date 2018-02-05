package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	strftime "github.com/lestrrat/go-strftime"
	"github.com/spf13/viper"
)

type configuration struct {
	AppLogPath                 string
	ModelLogPath               string
	DataCachesPath             string
	IngestorGobs               string
	ActivationServerAddress    string
	ActivationServiceEndpoint  string
	IngestorServerAddress      string
	IngestorActivationEndpoint string
	BackendServerAddress       string
	BackendActivationEndpoint  string
}

var initOnceCnf sync.Once

// Config is a public variable that allows for package-level settings changes
var Config configuration

func init() {
	initOnceCnf.Do(func() {
		viper.SetConfigName("config")                                             // name of the config file
		viper.AddConfigPath("$GOPATH/src/core/tests/cmd/backtests/bhattacharya/") // look for config in the working directory
		viper.AddConfigPath(".")                                                  // optionally for config in the working directory

		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			panic(fmt.Errorf("fatal config file read: %s", err))
		}

		viper.Unmarshal(&Config)
		Config.AppLogPath = replaceEnvVariable(fmtTimestamp(Config.AppLogPath))
		Config.ModelLogPath = replaceEnvVariable(fmtTimestamp(Config.ModelLogPath))
		Config.DataCachesPath = replaceEnvVariable(fmtTimestamp(Config.DataCachesPath))
		Config.IngestorGobs = replaceEnvVariable(Config.IngestorGobs)

		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			var cnf configuration
			viper.Unmarshal(&cnf)
			Config = cnf
			fmt.Println("Config file changed:", e.Name)
		})
	})
}

func replaceEnvVariable(path string) string {
	if strings.HasPrefix(path, "$GOPATH") {
		goPath := os.Getenv("GOPATH")
		return filepath.Join(goPath, path[7:])
	}
	return path
}

func fmtTimestamp(path string) string {
	f, err := strftime.New(path)
	if err != nil {
		panic(fmt.Errorf("new config strftime error", err))
	}
	return f.FormatString(time.Now())
}
