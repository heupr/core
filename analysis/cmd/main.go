package main

import (
	"coralreefci/analysis/backtests"
	"coralreefci/models/bhattacharya"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start().Stop()
	logger := bhattacharya.CreateLog("bhattacharya-backtest", true)
	nbModel := bhattacharya.Model{Classifier: &bhattacharya.NBClassifier{Logger: &logger}, Logger: &logger}
	testContext := backtests.TestContext{Model: nbModel}
	testRunner := backtests.BackTestRunner{Context: testContext}
	testRunner.Run()
	logger.Flush()
}
