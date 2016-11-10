package gateway

import (
	"github.com/google/go-github/github"
)

var variables = []string{"dotnet", "corefx"}

type Gateway struct {
	Client *github.Client
}

func (c *Gateway) GetPullRequests() ([]*github.PullRequest, error) {
	pullsOpt := &github.PullRequestListOptions{
		State: "closed",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	pulls := []*github.PullRequest{}
	for {
		pullRequests, resp, err := c.Client.PullRequests.List("dotnet", "corefx", pullsOpt)
		if err != nil {
			return nil, err
		}
		pulls = append(pulls, pullRequests...)

		if resp.NextPage == 0 {
			break
		} else {
			pullsOpt.ListOptions.Page = resp.NextPage
		}
	}
	return pulls, nil
}

func (c *Gateway) GetIssues() ([]*github.Issue, error) {
	// TODO: Handle opened/closed
	issuesOpt := &github.IssueListByRepoOptions{
		State: "closed",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	filteredIssues := []*github.Issue{}
	for {
		issues, resp, err := c.Client.Issues.ListByRepo("dotnet", "corefx", issuesOpt)
		for i := 0; i < len(issues); i++ {
			if err != nil {
				return nil, err
			}
			if issues[i].PullRequestLinks == nil {
				filteredIssues = append(filteredIssues, issues[i])
			}
		}
		if resp.NextPage == 0 {
			break
		} else {
			issuesOpt.ListOptions.Page = resp.NextPage
		}
	}
	return filteredIssues, nil
}

// TODO: this may not be needed if a better mapping alternative is found
func (c *Gateway) GetPullEvents() ([]*github.PullRequestEvent, error) {
	pullEvents := []*github.PullRequestEvent{}
	return pullEvents, nil
}

// TODO: this may not be needed if a better mapping alternative is found
func (c *Gateway) GetIssueEvents() ([]*github.Event, error) {
	issuesEvents, _, _ := c.Client.Activity.ListIssueEventsForRepository("dotnet", "corefx", nil)
	return issuesEvents, nil
}
