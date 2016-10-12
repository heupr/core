package bhattacharya

import (
	"coralreefci/models/issues"
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

func (c *NBClassifier) Predict(issue issues.Issue) []string {
	scores, _, _ := c.classifier.LogScores(strings.Split(issue.Body, " "))
	first, second, third, _ := topThree(scores)
	return []string{string(c.assignees[first]), string(c.assignees[second]), string(c.assignees[third])}
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

// findMax finds the maximum of a set of scores; if the
// maximum is strict -- that is, it is the single unique
// maximum from the set -- then strict has return value
// true. Otherwise it is false.
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
		if (scores[i] > scores[first]) {
			third = second
			second = first
			first = i
		} else if (scores[i] > scores[second] && scores[i] < scores[first]) {
			second = i
		} else if (scores[i] > scores[third] && scores[i] < scores[second]) {
			third = i
		}
	}
	return
}
