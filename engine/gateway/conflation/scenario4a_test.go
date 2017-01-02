package conflation

import (
	"github.com/google/go-github/github"
	"strings"
	"testing"
)

// DOC: These are the necessary initiation and test variables.
var testScenario4a = Scenario4a{Words: 4}
var bodyText = "I am your father."
var wordCount = 4
var issue = github.Issue{Body: &bodyText}
var testExpandedIssue = &ExpandedIssue{Issue: CRIssue{issue, []int{}, []CRPullRequest{}}}

func TestFilter4a(t *testing.T) {
	functionCount := strings.Count(*testExpandedIssue.Issue.Body, " ") + 1
	if functionCount != wordCount {
		t.Error(
			"\nWORD COUNT MISMATCH",
			"\nEXPECTED: ", wordCount,
			"\nACTUAL:   ", functionCount,
		)
	}
}
