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

type Configuration struct {
	AppLogPath     string
	ModelLogPath   string
	DataCachesPath string
	BoltDBPath 		 string
}

var initOnceCnf sync.Once

var Config Configuration

func init() {
	initOnceCnf.Do(func() {
		viper.SetConfigName("config")                                                       // name of the config file
		viper.AddConfigPath(".")                                                            // look for config in the working directory
		viper.AddConfigPath("$GOPATH/src/core/tests/cmd/backtests/bhattacharya/") // optionally look here

		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

		viper.Unmarshal(&Config)
		Config.AppLogPath = replaceEnvVariable(fmtTimestamp(Config.AppLogPath))
		Config.ModelLogPath = replaceEnvVariable(fmtTimestamp(Config.ModelLogPath))
		Config.DataCachesPath = replaceEnvVariable(fmtTimestamp(Config.DataCachesPath))
		Config.BoltDBPath = replaceEnvVariable(Config.BoltDBPath)

		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			var cnf Configuration
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
	} else {
		return path
	}
}

func fmtTimestamp(path string) string {
	f, err := strftime.New(path)
	if err != nil {
		panic(fmt.Errorf("config error:", err))
	}
	return f.FormatString(time.Now())
}
