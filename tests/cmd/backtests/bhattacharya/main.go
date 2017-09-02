package main

import (
	"core/models"
	"core/models/bhattacharya"
	"core/pipeline/ingestor"
	"core/tests/cmd/endtoendtests/replay"
	"github.com/pkg/profile"
)

func main() {
	bufferPool := ingestor.NewPool()
	db := ingestor.Database{BufferPool: bufferPool}
	db.Open()
	bs := replay.BacktestServer{DB: &db}
	go bs.Start()

	defer profile.Start().Stop()
	nbModel := models.Model{Algorithm: &bhattacharya.NBModel{}}
	testContext := TestContext{Model: nbModel}
	testRunner := BackTestRunner{Context: testContext}
	//testRunner.Run("golang/go")
	//testRunner.Run("docker/docker")
	//testRunner.Run("dotnet/roslyn")
	//testRunner.Run("openSUSE/osem")
	testRunner.Run("dotnet/corefx")
	//testRunner.Run("dotnet/coreclr")

	//testRunner.Run("fabric8io/fabric8")
	//testRunner.Run("systemd/systemd")
	//testRunner.Run("checkstyle/checkstyle")
	//testRunner.Run("twosigma/beaker-notebook")
	//testRunner.Run("HabitRPG/habitrpg")
	//testRunner.Run("grafana/grafana")
}
