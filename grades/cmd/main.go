package main

import (
  "coralreef-ci/grades"
  "coralreef-ci/models/bhattacharya"
)

func main() {
  nbModel := bhattacharya.Model{Classifier: &bhattacharya.NBClassifier{}}
  testContext := grades.TestContext{File: "/home/michael/golang/src/coralreef-ci/data/training/static/trainingset_corefx", Model: nbModel}
  testRunner := grades.BackTestRunner{Context: testContext}
  testRunner.Run()
}
