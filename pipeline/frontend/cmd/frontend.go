package main

import (
	"bytes"
	"runtime/debug"

	"core/pipeline/frontend"
	"core/utils"
)

func main() {
	defer func() {
		utils.SlackLog.Fatal("Process crash: ", recover(), bytes.NewBuffer(debug.Stack()).String())
	}()

	server := frontend.NewServer()
	server.Start()
}
