package classifier

import "coralreefci/models/issues"

// DOC: Classifier serves as the "contract" that all models must abide by in
//      that they should provide Learn and Predict methods.
type Classifier interface {
	Learn(issues []issues.Issue)
	Predict(issue issues.Issue) []string
}
