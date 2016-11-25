package bhattacharya

import (
	"coralreefci/models/classifier"
	"coralreefci/models/issues"
	// "coralreefci/models/model"  // NOTE: eventual import after refactoring
    // "fmt"
)

// TODO: refactor into separate directory
type Model struct {
	Classifier classifier.Classifier
	Logger     *CoralReefLogger
}

func (model *Model) Learn(issues []issues.Issue) {
    // fmt.Println(issues)
	// stopwordsOutput := RemoveStopWords(issues...)    // NOTE: eventual implementation
    // fmt.Println(stopwordsOutput)
	// stemmerOutput := StemIssues(stopwordsOutput...)  // NOTE: eventual implementation
	// model.Classifier.Learn(stemmerOutput)            // NOTE: eventual implementation
    // model.Classifier.Learn(issues)
    // model.Classifier.Learn(stopwordsOutput)         // NOTE: temporary skipping stemmer results
	model.Classifier.Learn(issues)              // TODO: REMOVE
}

func (model *Model) Predict(issue issues.Issue) []string {
	// issue.Body = RemoveStopWords(issue)[0].Body     // NOTE: eventual implementation
	// issue.Body = StemIssues(issue)[0].Body          // NOTE: eventual implementation
    // return model.Classifier.Predict(issue)          // NOTE: eventual implementation
	return model.Classifier.Predict(issue)      // TODO: REMOVE
}
