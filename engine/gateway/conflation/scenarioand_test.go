package conflation

import (
	"github.com/google/go-github/github"
	"testing"
)

type TestSubScenario struct{}

func (s TestSubScenario) Filter(expandedIssue *ExpandedIssue) bool {
	return true
}

var TestScenarioAND = ScenarioAND{Scenarios: []Scenario{TestSubScenario{}}}

var (
	title        = "Let the Wookie win."
	andIssue     = github.Issue{Title: &title}
	andTestIssue = &ExpandedIssue{Issue: CRIssue{issue, []int{}, []CRPullRequest{}}}
)

func TestFilterAND(t *testing.T) {
	if !TestScenarioAND.Filter(andTestIssue) {
		t.Error(
			"\n\"AND\" LOGIC IS NOT OPERATING CORRECTLY",
			"\nSUB SCENARIOS PROVIDED:", len(TestScenarioAND.Scenarios),
		)
	}
}
