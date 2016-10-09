package models

import "coralreef-ci/models/issues"

type Classifier interface {
	Learn(issues []issues.Issue)
	Predict(issue issues.Issue) []string
}
