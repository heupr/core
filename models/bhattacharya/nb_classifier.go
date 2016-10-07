package bhattacharya

import (
	"coralreef-ci/models/issues"
	"github.com/jbrukh/bayesian"
	"strings"
)

type NBClassifier struct {
	classifier *bayesian.Classifier
	assignees  []bayesian.Class
}

func (c *NBClassifier) Learn(issues []issues.Issue) {
	c.assignees = distinctAssignees(issues)
	c.classifier = bayesian.NewClassifierTfIdf(c.assignees...)
	for i := 0; i < len(issues); i++ {
		c.classifier.Learn(strings.Split(issues[i].Body, " "), bayesian.Class(issues[i].Assignee))
	}
	c.classifier.ConvertTermsFreqToTfIdf()
}

func (c *NBClassifier) Predict(issue issues.Issue) string {
	return string(c.assignees[0])
}

func distinctAssignees(issues []issues.Issue) []bayesian.Class {
	result := []bayesian.Class{}
	j := 0
	for i := 0; i < len(issues); i++ {
		for j = 0; j < len(result); j++ {
			if issues[i].Assignee == string(result[j]) {
				break
			}
		}
		if j == len(result) {
			result = append(result, bayesian.Class(issues[i].Assignee))
		}
	}
	return result
}
