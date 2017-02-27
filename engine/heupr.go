package main

import (
	"coralreefci/engine/frontend"
	"coralreefci/models"
)

func main() {
	server := createServer()
	server.Start()
}

func createServer() frontend.HeuprServer {
	//Models might be read from a configuration
	modelMap := make(map[int]models.Model)
	//TODO: pointer? &modelMap
	dispatcher := frontend.Dispatcher{Models: modelMap}
	dispatcher.Start(5)
	//TODO: This case is only valid during recovery
	//modelMap[getRepoId("corefx")] = models.Model{Algorithm: &bhattacharya.NBModel{}}
	return frontend.HeuprServer{Models: modelMap}
}

//TODO: flesh out this logic
// func getRepoId(repo string) int {
// 	if repo == "corefx" {
// 		return 555
// 	} else {
// 		return -1
// 	}
// }
