package main

import (
	"bytes"
	"runtime/debug"

	"core/pipeline/backend"
	"core/utils"
)

func main() {
	defer func() {
		if backend.PROD {
			utils.SlackLog.Fatal("Process crash: ", recover(), bytes.NewBuffer(debug.Stack()).String())
		}
	}()

	backendServer := backend.Server{}
	backendServer.Repos = new(backend.ActiveRepos)
	backendServer.Repos.Actives = make(map[int64]*backend.ArchRepo)
	backendServer.Start()
}
