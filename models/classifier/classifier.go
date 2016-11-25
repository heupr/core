package classifier

import "coralreefci/models/issues"

type Classifier interface {
	Learn(issues []issues.Issue)
	Predict(issue issues.Issue) []string
}
