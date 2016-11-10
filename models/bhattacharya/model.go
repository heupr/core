package bhattacharya

import (
	"coralreefci/models/classifier"
	"coralreefci/models/issues"
)

type Model struct {
	Classifier models.Classifier
	Logger *CoralReefLogger
}

func (model *Model) Learn(issues []issues.Issue) {
	removeStopWords(issues)
	//StemIssues(issues)
	model.Classifier.Learn(issues)
}

func (model *Model) Predict(issue issues.Issue) []string {
	issue.Body = RemoveStopWords(issue.Body)
	//StemIssues([]issues.Issue{issue})
	return model.Classifier.Predict(issue)
}

func removeStopWords(issues []issues.Issue) {
	for i := 0; i < len(issues); i++ {
		issues[i].Body = RemoveStopWords(issues[i].Body)
	}
}
