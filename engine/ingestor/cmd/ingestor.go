package main

import (
	"coralreefci/engine/ingestor"
)

func main() {
	dispatcher := ingestor.Dispatcher{}
	dispatcher.Start(5)

	ingestorServer := ingestor.IngestorServer{}
	ingestorServer.Start()
}
