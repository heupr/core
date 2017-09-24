package main

import (
	"bytes"
	"runtime/debug"

	"core/pipeline/backend"
	"core/utils"
)

func main() {
	defer func() {
		utils.SlackLog.Fatal("Process crash: ", recover(), bytes.NewBuffer(debug.Stack()).String())
	}()

	backendServer := backend.BackendServer{}
	backendServer.Repos = new(backend.ActiveRepos)
	backendServer.Repos.Actives = make(map[int]*backend.ArchRepo)
	backendServer.Start()
}
