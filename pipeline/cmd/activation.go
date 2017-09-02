package main

import (
	"core/pipeline"
)

func main() {
	activationServer := pipeline.ActivationServer{}
	activationServer.Start()
}
