package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	LogPath string
	ModelSummaryPath string
}

func Load() Config {
	file, _ := os.Open("./config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("config error:", err)
	}
	return configuration
}
