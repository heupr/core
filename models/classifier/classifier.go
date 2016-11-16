/*
The classifier package implmenets a basic interface to all backend models.
*/
package models

import "coralreefci/models/issues"

type Classifier interface {
	Learn(issues []issues.Issue)
	Predict(issue issues.Issue) []string
}
