package main

import (
	"bytes"
	"runtime/debug"

	"core/pipeline/ingestor"
	"core/utils"
)

func main() {
	defer func() {
		utils.SlackLog.Fatal("Process crash: ", recover(), bytes.NewBuffer(debug.Stack()).String())
	}()

	ingestorServer := ingestor.IngestorServer{}
	ingestorServer.Start()
}
