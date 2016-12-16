package bhattacharya

import (
	"coralreefci/models/issues"
	"github.com/jbrukh/bayesian"
	"strings"
)

// DOC: NBClassifier is the specific type of classifier being utilized in this
//      particular back-end model implmentation.
type NBClassifier struct {
	classifier *bayesian.Classifier
	assignees  []bayesian.Class
	graph      *TossingGraph       // NOTE: Restructure to be called on Model
	Logger     *CoralReefLogger    // NOTE: Directly imported from new external logger package
}

func (c *NBClassifier) Learn(issues []issues.Issue) {
	RemoveStopWords(issues...)
	StemIssues(issues...)
	c.assignees = distinctAssignees(issues)
	c.classifier = bayesian.NewClassifierTfIdf(c.assignees...)
	c.graph = &TossingGraph{Assignees: convertClassToString(c.assignees), GraphDepth: 5, Logger: c.Logger}
	for i := 0; i < len(issues); i++ {
		c.classifier.Learn(strings.Split(issues[i].Body, " "), bayesian.Class(issues[i].Assignee))
	}
	c.classifier.ConvertTermsFreqToTfIdf()
}

func (c *NBClassifier) Predict(issue issues.Issue) []string {
	RemoveStopWordsSingle(&issue)
	StemIssuesSingle(&issue)
	scores, _, _ := c.classifier.LogScores(strings.Split(issue.Body, " "))
	names := []string{}
	indices := c.graph.Tossing(scores)
	for i := 0; i < len(indices); i++ {
		names = append(names, string(c.assignees[indices[i]]))
	}
	return names
}

// DOC: distinctAssignees is a helper function that may ultimately be
//      refactored out of this particular package.
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

// DOC: convertClassToString is a helper function designed to overcome the
//      limitations presented in the bayesian package where classes are stored
//      as an new type Class rather than a Go string type.
func convertClassToString(assignees []bayesian.Class) []string {
	result := []string{}
	for i := 0; i < len(assignees); i++ {
		result = append(result, string(assignees[i]))
	}
	return result
}
