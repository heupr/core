package bhattacharya

import (
	//"time"
	"coralreef-ci/models/classifier"
	"coralreef-ci/models/issues"
)

type Model struct {
	Classifier models.Classifier
}

func (model *Model) Learn(issues []issues.Issue) {
	removeStopWords(issues)
	//Shuffle(issues, int64(time.Now().Nanosecond()))
	//StemIssues(issues)
	model.Classifier.Learn(issues)
}

func (model *Model) Predict(issue issues.Issue) (string, string, string) {
	//StemIssues([]issues.Issue{issue})
	return model.Classifier.Predict(issue)
}

func removeStopWords(issues []issues.Issue) {
	for i := 0; i < len(issues); i++ {
		issues[i].Body = RemoveStopWords(issues[i].Body)
	}
}
