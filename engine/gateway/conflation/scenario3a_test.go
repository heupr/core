package conflation

import (
	"github.com/google/go-github/github"
	"testing"
)

var TestScenario3a = Scenario3a{}

var url = "https://www.rule-of-two.com/"
var pullRequest = github.PullRequest{IssueURL: &url}
var TestWithPullRequest = &ExpandedIssue{PullRequest: CrPullRequest{pullRequest, []int{}, []CrIssue{}}}
var TestWithoutPullRequest = &ExpandedIssue{}

func TestFilter3a(t *testing.T) {
	withURL := TestScenario3a.Filter(TestWithPullRequest)
	if withURL != false {
		t.Error(
			"PULL REQUEST WITH ASSOCIATED ISSUES INCLUDED",
		)
	}
	withoutURL := TestScenario3a.Filter(TestWithoutPullRequest)
	if withoutURL != true {
		t.Error(
			"PULL REQUEST WITHOUT ISSUES EXCLUDED",
		)
	}
}
