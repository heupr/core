package models

import (
    "coralreefci/models/bhattacharya"
	"coralreefci/engine/gateway/conflation"
    "github.com/google/go-github/github"
	"strconv"
	"testing"
)

func buildTestIssues() []conflation.ExpandedIssue {
    issues := []conflation.ExpandedIssue{}
    for i := 1; i < 31; i ++ {
        if i % 2 == 0 {
            name := "JOHN"
            assignee := github.User{Name: &name}
        } else {
            name := "MIKE"
            assignee := github.User{Name: &name}
        }
        githubIssue := github.Issue{Assignee: assignee}
        crIssue := conflation.CRIssue{githubIssue}
        issues = append(issues, conflation.ExpandedIssue{crIssue})
    }
    return issues
}

func TestFold(t *testing.T) {
	nbModel := Model{Algorithm: &bhattacharya.NBClassifier{}}
    testingIssues := buildTestIssues()
	result, _ := nbModel.JohnFold(testingIssues)
	number, _ := strconv.ParseFloat(result, 64)
	if number < 0.00 && number > 1.00 {
		t.Error(
			"\nRESULT IS OUTSIDE ACCEPTABLE RANGE - JOHN FOLD",
			"\nEXPECTED BETWEEN 0.00 AND 1.00",
			"\nACTUAL: %f", number,
		)
	}
}
