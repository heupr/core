package conflation

import (
	"github.com/google/go-github/github"
	"testing"
)

var TestScenario1b = Scenario1b{}

// DOC: this is the variable contruction for the dummy pure GitHub issue.
var (
	number    = 1187
	noNumber  = 0
	iTitle    = "Begun the clone war has."
	comments  = 3
	testIssue = github.Issue{Number: &number, Title: &iTitle, Comments: &comments}
	nonIssue  = github.Issue{Number: &noNumber}
)

// DOC: below is the empty test GitHub pull request; no values are needed.
var (
	prTitle         = "We have been expecting you."
	testPullRequest = github.PullRequest{Title: &prTitle}
)

var TestWithIssue = &ExpandedIssue{
	Issue:       CrIssue{testIssue, []int{}, []CrPullRequest{}},
	PullRequest: CrPullRequest{testPullRequest, []int{}, []CrIssue{}},
}

var TestWithoutIssue = &ExpandedIssue{
	Issue:       CrIssue{nonIssue, []int{}, []CrPullRequest{}},
	PullRequest: CrPullRequest{testPullRequest, []int{}, []CrIssue{}},
}

func TestFilter1b(t *testing.T) {
	firstOutput := TestScenario1b.Filter(TestWithIssue)
	if firstOutput != true {
		t.Error(
			"\nISSUE WITH COMMENT INCORRECTLY FILTERED OUT",
			"\nISSUE NUMBER:          ", *TestWithIssue.Issue.Number,
			"\nBOOLEAN FILTER RETURN: ", firstOutput,
		)
	}
	secondOutput := TestScenario1b.Filter(TestWithoutIssue)
	if secondOutput != false {
		t.Error(
			"\nNONEXISTENT ISSUE NOT FILTERED",
			"\nISSUE NUMBER:          ", *TestWithoutIssue.Issue.Number,
			"\nBOOLEAN FILTER RETURN: ", secondOutput,
		)
	}
}
