package main

import (
	"coralreefci/grades"
	"coralreefci/models/bhattacharya"
)

func main() {
	nbModel := bhattacharya.Model{Classifier: &bhattacharya.NBClassifier{}}
	testContext := grades.TestContext{File: "./trainingset_corefx", Model: nbModel}
	testRunner := grades.BackTestRunner{Context: testContext}
	testRunner.Run()
}
