package bhattacharya

import (
	"../issues"
	"../../models"
)

// Model...
type Model struct {
	classifier models.Classifier
}

func (model *Model) Learn(issues []issues.Issue) {
	//TODO: Implement Learn Unit Test
	//TODO: Implement Learn
	//TODO: Add shuffle (argument seed)
	//TODO: Add Logging
	removeStopWords(issues)
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
