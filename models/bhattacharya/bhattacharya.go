package bhattacharya

import (
	"coralreefci/models/classifier"
	"coralreefci/models/issues"
)

// TODO: refactor into separate directory
// NOTE: this struct will no longer be necessary and will be subsumed by the
//       "NBClassifier" struct available in naive_bayes.go
type Model struct {
	Classifier classifier.Classifier
	Logger     *CoralReefLogger
}

// TODO: delete entirely
func (model *Model) Learn(issues []issues.Issue) {
	RemoveStopWords(issues...)     // TODO: move to naive_bayes.go
	StemIssues(issues...)          // TODO: move to naive_bayes.go
	model.Classifier.Learn(issues)
}

// TODO: delete entirely
func (model *Model) Predict(issue issues.Issue) []string {
	RemoveStopWordsSingle(&issue)  // TODO: move to naive_bayes.go
	StemIssuesSingle(&issue)       // TODO: move to naive_bayes.go
	return model.Classifier.Predict(issue)
}
