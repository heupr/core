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
	graph      *TossingGraph
	Logger     *CoralReefLogger
}

func (c *NBClassifier) Learn(issues []issues.Issue) {
	c.assignees = distinctAssignees(issues)
	c.classifier = bayesian.NewClassifierTfIdf(c.assignees...)
	c.graph = &TossingGraph{Assignees: convertClassToString(c.assignees), GraphDepth: 5, Logger: c.Logger}
	for i := 0; i < len(issues); i++ {
		c.classifier.Learn(strings.Split(issues[i].Body, " "), bayesian.Class(issues[i].Assignee))
	}
	c.classifier.ConvertTermsFreqToTfIdf()
}

func (c *NBClassifier) Predict(issue issues.Issue) []string {
	scores, _, _, err := c.classifier.SafeProbScores(strings.Split(issue.Body, " "))
	if err != nil {
		scores, _, _ = c.classifier.LogScores(strings.Split(issue.Body, " "))
	}
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

// DOC: findMax locates the maximum value in a score set.
//      If strict is set to true (single unique maximum value); the default
//      is true.
func findMax(scores []float64) (inx int, strict bool) {
	inx = 0
	strict = true
	for i := 1; i < len(scores); i++ {
		if scores[inx] < scores[i] {
			inx = i
			strict = true
		} else if scores[inx] == scores[i] {
			strict = false
		}
	}
	return
}

func topThree(scores []float64) (first int, second int, third int, strict bool) {
	first = 0
	second = 0
	third = 0
	strict = true
	for i := 1; i < len(scores); i++ {
		if scores[i] > scores[first] {
			third = second
			second = first
			first = i
		} else if scores[i] > scores[second] && scores[i] < scores[first] {
			second = i
		} else if scores[i] > scores[third] && scores[i] < scores[second] {
			third = i
		}
	}
	return
}
