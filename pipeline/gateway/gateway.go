package gateway

import (
	"context"

	"github.com/google/go-github/github"
)

type Gateway struct {
	Client      *github.Client
	UnitTesting bool
}

func (g *Gateway) getPulls(owner, repo, state string) ([]*github.PullRequest, error) {
	opt := &github.PullRequestListOptions{
		State: state,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	output := []*github.PullRequest{}

	for {
		pulls, resp, err := g.Client.PullRequests.List(context.Background(), owner, repo, opt)
		if err != nil {
			return nil, err
		}
		output = append(output, pulls...)

		if resp.NextPage == 0 || g.UnitTesting {
			break
		} else {
			opt.ListOptions.Page = resp.NextPage
		}
	}
	return output, nil
}

func (g *Gateway) getIssues(owner, repo, state string) ([]*github.Issue, error) {
	opt := &github.IssueListByRepoOptions{
		State: state,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	output := []*github.Issue{}
	for {
		issues, resp, err := g.Client.Issues.ListByRepo(context.Background(), owner, repo, opt)
		if err != nil {
			return nil, err
		}
		output = append(output, issues...)
		if resp.NextPage == 0 || g.UnitTesting {
			break
		} else {
			opt.ListOptions.Page = resp.NextPage
		}
	}
	return output, nil
}

func (g *Gateway) GetOpenPulls(owner, repo string) ([]*github.PullRequest, error) {
	return g.getPulls(owner, repo, "open")
}

func (g *Gateway) GetClosedPulls(owner, repo string) ([]*github.PullRequest, error) {
	return g.getPulls(owner, repo, "closed")
}

func (g *Gateway) GetOpenIssues(owner, repo string) ([]*github.Issue, error) {
	return g.getIssues(owner, repo, "open")
}

func (g *Gateway) GetClosedIssues(owner, repo string) ([]*github.Issue, error) {
	return g.getIssues(owner, repo, "closed")
}
