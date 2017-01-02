package bhattacharya

import (
	"coralreefci/engine/gateway/conflation"
	"github.com/jbrukh/bayesian"
	"sort"
	"strings"
)

// DOC: NBClassifier is the struct implemented as the model algorithm.
type NBClassifier struct {
	classifier *bayesian.Classifier
	assignees  []bayesian.Class
}

// TODO: remove assets into separate file
type Result struct {
	id    int
	score float64
}

type Results []Result

func (slice Results) Len() int {
	return len(slice)
}

func (slice Results) Less(i, j int) bool {
	return slice[i].score < slice[j].score
}

func (slice Results) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (c *NBClassifier) Learn(input []conflation.ExpandedIssue) {
	adjusted := c.converter(input...)
	RemoveStopWords(adjusted...)
	StemIssues(adjusted...)
	c.assignees = distinctAssignees(adjusted)
	c.classifier = bayesian.NewClassifierTfIdf(c.assignees...)
	for i := 0; i < len(input); i++ {
		c.classifier.Learn(strings.Split(adjusted[i].Body, " "), bayesian.Class(adjusted[i].Assignee))
	}
	c.classifier.ConvertTermsFreqToTfIdf()
}

func (c *NBClassifier) Predict(input conflation.ExpandedIssue) []string {
	adjusted := c.converter(input)
	RemoveStopWordsSingle(&adjusted[0])
	StemIssuesSingle(&adjusted[0])
	scores, _, _ := c.classifier.LogScores(strings.Split(adjusted[0].Body, " "))

	results := Results{}
	for i := 0; i < len(scores); i++ {
		results = append(results, Result{id: i, score: scores[i]})
	}

	sort.Reverse(results)

	names := []string{}
	for i := 0; i < len(results); i++ {
		names = append(names, string(c.assignees[results[i].id]))
	}
	return names
}

// DOC: distinctAssignees is a helper function that finds and collects all of
//	    the unique occurances of the Class type (a "string" alias that is
//      unique to the bayesian third-party package.
func distinctAssignees(issues []Issue) []bayesian.Class {
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
//      Return - string of the feature used in the Bayesian classifier.
func convertClassToString(assignees []bayesian.Class) []string {
	result := []string{}
	for i := 0; i < len(assignees); i++ {
		result = append(result, string(assignees[i]))
	}
	return result
}
