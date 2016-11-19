package conflation

import (
	"github.com/google/go-github/github"
	"testing"
)

var TestScenario3a = Scenario3a{}

var (
	number         = 1
	body           = "At last we will reveal ourselves to the Jedi; at last we will have our revenge."
	title          = "Sith Apprentice"
	assignee       = "Darth Maul"
	githubAssignee = github.User{Name: &assignee}
	url            = "https://www.rule-of-two.com/"
	pullRequest    = github.PullRequest{Number: &number, Title: &title, Body: &body, IssueURL: &url, Assignee: &githubAssignee}
)

func TestFilter(t *testing.T) {
	withURL := TestScenario3a.Filter(pullRequest)
	if withURL != false {
		t.Error(
			"PULL REQUEST WITH URL NOT FILTERED",
		)
	}
	noURL := ""
	pullRequest.IssueURL = &noURL
	withoutURL := TestScenario3a.Filter(pullRequest)
	if withoutURL != true {
		t.Error(
			"PULL REQUEST WITHOUT URL NOT INCLUDED",
		)
	}
}
