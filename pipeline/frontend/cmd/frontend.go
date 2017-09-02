package main

import (
	"core/pipeline/frontend"
)

func main() {
	frontendServer := frontend.FrontendServer{}
	frontendServer.Start()
	frontendServer.OpenBolt()
}
