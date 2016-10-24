package gateway

import (
	"github.com/google/go-github/github"
)

var variables = []string{"dotnet", "corefx"}

type Gateway struct {
	Client *github.Client
}

func (c *Gateway) GetPullRequests() []*github.PullRequest {
	pullRequests, _, _ := c.Client.PullRequests.List(variables[0], variables[1], nil)
	return pullRequests
}

func (c *Gateway) GetIssues() []*github.Issue {
	filteredIssues := []*github.Issue{}
	issues, _, _ := c.Client.Issues.ListByRepo(variables[0], variables[1], nil)
	for i := 0; i < len(issues); i++ {
		if issues[i].PullRequestLinks == nil {
			filteredIssues = append(filteredIssues, issues[i])
		}
	}
	return filteredIssues
}

// TODO: this may not be needed if a better mapping alternative is found
func (c *Gateway) GetPullEvents() []github.PullRequestEvent {
	pullEvents := []github.PullRequestEvent{}
	return pullEvents
}

// TODO: this may not be needed if a better mapping alternative is found
func (c *Gateway) GetIssueEvents() []*github.Event {
	issuesEvents, _, _ := c.Client.Activity.ListIssueEventsForRepository(variables[0], variables[1], nil)
	return issuesEvents
}
