package conflation

import (
	"github.com/google/go-github/github"
	"testing"
)

func TestConflater(t *testing.T) {
	context := &Context{}

	scenarios := []Scenario{&Scenario3{}}
	conflationAlgorithms := []ConflationAlgorithm{&ComboAlgorithm{Context: context}}
	normalizer := Normalizer{Context: context}
	conflator := Conflator{Scenarios: scenarios, ConflationAlgorithms: conflationAlgorithms, Normalizer: normalizer, Context: context}

	issueNumber := 3720
	issueTitle := "The odds"
	issueBody := "Sir, the possibility of successfully navigating an asteroid field is approximately 3,720 to 1!"
	issueLogin := "C-3P0"
	issueAssignee := github.User{Login: &issueLogin}
	issueAssignees := []*github.User{&issueAssignee}
	githubIssue := github.Issue{Number: &issueNumber, Title: &issueTitle, Body: &issueBody, Assignees: issueAssignees}
	issues := []github.Issue{githubIssue}

	pullNumber := 1
	pullBody := "Never tell me the odds!"
	pullTitle := "The response"
	pullLogin := "Han"
	pullAssignee := github.User{Name: &pullLogin}
	issueURL := "https://asteroid-belt.hoth/"
	githubPull := github.PullRequest{Number: &pullNumber, Title: &pullTitle, Body: &pullBody, IssueURL: &issueURL, User: &pullAssignee}
	pulls := []github.PullRequest{githubPull}

	conflator.Context.Issues = []ExpandedIssue{}
	conflator.SetIssueRequests(issues)
	conflator.SetPullRequests(pulls)

	conflator.Conflate()

	for i := 0; i < len(conflator.Context.Issues); i++ {
		if conflator.Context.Issues[i].Issue.Number != nil && *conflator.Context.Issues[i].Issue.Number == 12886 {
			Assert(t, *pullAssignee.Name, *conflator.Context.Issues[i].Issue.Assignee.Name, "C-3P0/Han Issue/Pull")
		}
	}
}

func Assert(t *testing.T, expected string, actual string, input string) {
	if actual != expected {
		t.Error(
			"\nFOR:       ", input,
			"\nEXPECTED:  ", expected,
			"\nACTUAL:    ", actual,
		)
	}
}
