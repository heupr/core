package main

import (
	"github.com/pkg/profile"

	"core/models"
	"core/models/bhattacharya"
)

// Here is a collected list of repositories that can be plugged into the test:
// https://paper.dropbox.com/doc/Targeted-Repo-List-P22Hovh0G8ckJkanLE7nW
func main() {
	defer profile.Start().Stop()
	nbModel := bhattacharya.NBModel{}
	model := models.Model{Algorithm: &nbModel}
	testContext := TestContext{Model: model}
	testRunner := BackTestRunner{Context: testContext}
	testRunner.Run("dotnet/corefx")
}
