package main

import (
	"coralreefci/models"
	"coralreefci/models/bhattacharya"
	// "github.com/pkg/profile"
)

func main() {
	// defer profile.Start().Stop()
	nbModel := models.Model{Algorithm: &bhattacharya.NBModel{}}
	testContext := TestContext{Model: nbModel}
	testRunner := BackTestRunner{Context: testContext}
	testRunner.Run()
}
