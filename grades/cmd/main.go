package main

import (
	"coralreefci/grades"
	"coralreefci/models/bhattacharya"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start().Stop()
	logger := bhattacharya.CreateLog("bhattacharya-backtest", true)
	nbModel := bhattacharya.Model{Classifier: &bhattacharya.NBClassifier{Logger: &logger}, Logger: &logger}
	testContext := grades.TestContext{Model: nbModel}
	testRunner := grades.BackTestRunner{Context: testContext}
	testRunner.Run()
	logger.Flush()
}
