package main

import (
	"core/pipeline/ingestor"
)

func main() {
	dispatcher := ingestor.Dispatcher{}
	dispatcher.Start(5)

	ingestorServer := ingestor.IngestorServer{}
	ingestorServer.Start()
}
