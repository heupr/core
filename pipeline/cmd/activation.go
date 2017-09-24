package main

import (
	"bytes"
	"runtime/debug"

	"core/pipeline"
	"core/utils"
)

func main() {
	defer func() {
		utils.SlackLog.Fatal("Process crash: ", recover(), bytes.NewBuffer(debug.Stack()).String())
	}()

	activationServer := pipeline.ActivationServer{}
	activationServer.Start()
}
