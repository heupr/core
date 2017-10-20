package conflation

import (
	"testing"

	"github.com/google/go-github/github"
)

var TestScenario2 = Scenario2{}

// This is the variable contruction for the dummy pure GitHub issue.
var (
	number    = 1187
	noNumber  = 0
	iTitle    = "Begun the clone war has."
	comments  = 3
	testIssue = github.Issue{Number: &number, Title: &iTitle, Comments: &comments}
	nonIssue  = github.Issue{Number: &noNumber}
)

// Below is the empty test GitHub pull request; no values are needed.
var (
	prTitle         = "We have been expecting you."
	testPullRequest = github.PullRequest{Title: &prTitle}
)

var TestWithIssue = &ExpandedIssue{
	Issue:       CRIssue{testIssue, []int{}, []CRPullRequest{}, false},
	PullRequest: CRPullRequest{testPullRequest, []int{}, []CRIssue{}},
}

var TestWithoutIssue = &ExpandedIssue{
	Issue:       CRIssue{nonIssue, []int{}, []CRPullRequest{}, false},
	PullRequest: CRPullRequest{testPullRequest, []int{}, []CRIssue{}},
}

func TestFilter2(t *testing.T) {
	firstOutput := TestScenario2.Filter(TestWithIssue)
	if firstOutput != true {
		t.Error(
			"\nISSUE WITH COMMENT INCORRECTLY FILTERED OUT",
			"\nISSUE NUMBER:          ", *TestWithIssue.Issue.Number,
			"\nBOOLEAN FILTER RETURN: ", firstOutput,
		)
	}
	secondOutput := TestScenario2.Filter(TestWithoutIssue)
	if secondOutput != false {
		t.Error(
			"\nNONEXISTENT ISSUE NOT FILTERED",
			"\nISSUE NUMBER:          ", *TestWithoutIssue.Issue.Number,
			"\nBOOLEAN FILTER RETURN: ", secondOutput,
		)
	}
}
