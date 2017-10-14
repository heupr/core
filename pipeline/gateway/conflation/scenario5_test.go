package conflation

import (
	"strings"
	"testing"

    "github.com/google/go-github/github"
)

// These are the necessary initiation and test variables.
var testScenario5 = Scenario5{Words: 4}
var bodyText = "I am your father."
var wordCount = 4
var issue = github.Issue{Body: &bodyText}
var testExpandedIssue = &ExpandedIssue{Issue: CRIssue{issue, []int{}, []CRPullRequest{}, false}}

func TestFilter5(t *testing.T) {
	functionCount := strings.Count(*testExpandedIssue.Issue.Body, " ") + 1
	if functionCount != wordCount {
		t.Error(
			"\nWORD COUNT MISMATCH",
			"\nEXPECTED: ", wordCount,
			"\nACTUAL:   ", functionCount,
		)
	}
}
