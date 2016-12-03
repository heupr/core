package main

import (
	"coralreefci/grades"
	"coralreefci/models/bhattacharya"
)

func main() {
	logger := bhattacharya.CreateLog("bhattacharya-backtest")
	nbModel := bhattacharya.Model{Classifier: &bhattacharya.NBClassifier{Logger: &logger}, Logger: &logger}
	testContext := grades.TestContext{Model: nbModel}
	testRunner := grades.BackTestRunner{Context: testContext}
	testRunner.Run()
	logger.Flush()
}
