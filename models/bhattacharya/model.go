package bhattacharya

import (
	"coralreefci/models/classifier"
	"coralreefci/models/issues"
)

// TODO: refactor into separate directory
type Model struct {
	Classifier classifier.Classifier
	Logger     *CoralReefLogger
}

func (model *Model) Learn(issues []issues.Issue) {
	RemoveStopWords(issues...)
    StemIssues(issues...)
	model.Classifier.Learn(issues)
}

func (model *Model) Predict(issue issues.Issue) []string {
	RemoveStopWords(issue)
    StemIssues(issue)
	return model.Classifier.Predict(issue)
}
