package models

import "./issues"

type Classifier interface {
	Learn(issues []issues.Issue)
	Predict(issue issues.Issue) string
}
