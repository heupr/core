package bhattacharya

import (
	"coralreefci/models/classifier"
	"coralreefci/models/issues"
	"fmt"
)

// TODO: refactor into separate directory
type Model struct {
	Classifier classifier.Classifier
	Logger     *CoralReefLogger
}

func (model *Model) Learn(issues []issues.Issue) {
	RemoveStopWords(issues...)
	//StemIssues(issues...)
	model.Classifier.Learn(issues)
}

func (model *Model) Predict(issue issues.Issue) []string {
	RemoveStopWords(issue)
	fmt.Println("Predict Method (Model): ", issue.Body)
	//StemIssues(issue)
	return model.Classifier.Predict(issue)
}
