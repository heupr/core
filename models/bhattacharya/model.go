package bhattacharya

import (
	"coralreef-ci/models"
	"coralreef-ci/models/issues"
)

type Model struct {
	classifier models.Classifier
}

func (model *Model) Learn(issues []issues.Issue) {
	removeStopWords(issues)
	// TODO: implement the stemming functionality here
	// TODO: implement the shuffling functionality here
	model.classifier.Learn(issues)
}

func (model *Model) Predict(issue issues.Issue) string {
	return model.classifier.Predict(issue)
}

func removeStopWords(issues []issues.Issue) {
	for i := 0; i < len(issues); i++ {
		issues[i].Body = RemoveStopWords(issues[i].Body)
	}
}
