package bhattacharya

import (
	"github.com/jbrukh/bayesian"
	"strings"
	"sort"
)

// DOC: NBClassifier is the specific type of classifier being utilized in this
//      particular back-end model implmentation.
type NBClassifier struct {
	classifier *bayesian.Classifier
	assignees  []bayesian.Class
}

type Result struct {
  id       int
  score    float64
}

type Results []Result

func (slice Results) Len() int {
    return len(slice)
}

func (slice Results) Less(i, j int) bool {
    return slice[i].score < slice[j].score;
}

func (slice Results) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

func (c *NBClassifier) Learn(issues []Issue) {
	RemoveStopWords(issues...)
	StemIssues(issues...)
	c.assignees = distinctAssignees(issues)
	c.classifier = bayesian.NewClassifierTfIdf(c.assignees...)
	for i := 0; i < len(issues); i++ {
		c.classifier.Learn(strings.Split(issues[i].Body, " "), bayesian.Class(issues[i].Assignee))
	}
	c.classifier.ConvertTermsFreqToTfIdf()
}

func (c *NBClassifier) Predict(issue Issue) []string {
	RemoveStopWordsSingle(&issue)
	StemIssuesSingle(&issue)
	scores, _, _ := c.classifier.LogScores(strings.Split(issue.Body, " "))

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
//			the unique occurances of the Class type (a "string" alias that is
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
func convertClassToString(assignees []bayesian.Class) []string {
	result := []string{}
	for i := 0; i < len(assignees); i++ {
		result = append(result, string(assignees[i]))
	}
	return result
}
