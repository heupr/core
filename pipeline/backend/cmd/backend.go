package main

import (
	"core/pipeline/backend"
)

func main() {
	backendServer := backend.BackendServer{}
	backendServer.Repos = new(backend.ActiveRepos)
	backendServer.Repos.Actives = make(map[int]*backend.ArchRepo)
	backendServer.Start()
}
