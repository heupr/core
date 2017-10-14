package conflation

import (
    "testing"
    
	"github.com/google/go-github/github"
)

var TestScenario6 = Scenario6{AssigneeCount: 1}

func TestFilter6(t *testing.T) {
	login := "Darth Plagueis"
	user := github.User{Login: &login}
	issue := github.Issue{Assignees: []*github.User{&user}}
	pull := github.PullRequest{}

	ei1 := &ExpandedIssue{Issue: CRIssue{issue, []int{}, []CRPullRequest{}, false}}
	ei2 := &ExpandedIssue{}
	ei3 := &ExpandedIssue{PullRequest: CRPullRequest{pull, []int{}, []CRIssue{}}}

	if !TestScenario6.Filter(ei1) {
		t.Error(
			"\nISSUE IMPROPERLY EXCLUDED",
			"\nASSIGNEE COUNT:", TestScenario6.AssigneeCount,
			"\nISSUE:", *ei1,
		)
	}
	if TestScenario6.Filter(ei2) {
		t.Error(
			"\nISSUE WITHOUT ASSIGNEE INCLUDED",
			"\nISSUE:", *ei2,
		)
	}
	if TestScenario6.Filter(ei3) {
		t.Error(
			"\nPULL REQUEST IMPROPERLY INCLUDED",
			"\nONLY ISSUES EXPECTED",
			"\nPULL REQUEST:", *ei3,
		)
	}
}
