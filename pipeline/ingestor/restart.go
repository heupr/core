package ingestor

import (
	"context"

	"github.com/google/go-github/github"
	"go.uber.org/zap"

	"core/utils"
)

const restartQuery = `
SELECT ifnull(g.number, 0) as number, g.is_pull as is_pull
FROM (
    SELECT MAX(number) as number, ifnull(is_pull, FALSE) as is_pull
    FROM github_events
    WHERE repo_id = ? AND is_pull = FALSE

    UNION ALL

    SELECT MAX(number) as number, ifnull(is_pull, TRUE) as is_pull
    FROM github_events
    WHERE repo_id = ? AND is_pull = TRUE
) g
ORDER BY g.is_pull`

var issueGaps = func(client *github.Client, owner, name string, dbIssueNum int) ([]*github.Issue, error) {
	ctx := context.Background()
	opts := github.IssueListByRepoOptions{
		State: "all",
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	}
	issue, _, err := client.Issues.ListByRepo(ctx, owner, name, &opts)
	if err != nil {
		utils.AppLog.Error("newest GitHub issue retrival", zap.Error(err))
		return nil, err
	}

	githubIssueNum := 0
	if len(issue) > 0 {
		githubIssueNum = *issue[0].Number
	}

	diff := githubIssueNum - dbIssueNum
	missingIssues := []*github.Issue{}
	for diff > 0 {
		switch {
		case diff > 0 && diff <= 100:
			opts.ListOptions.PerPage = diff
			diff = 0
		case diff > 100:
			opts.ListOptions.PerPage = 100
			diff = diff - 100
		}
		issues, resp, err := client.Issues.ListByRepo(ctx, owner, name, &opts)
		if err != nil {
			utils.AppLog.Error("missing GitHub issue retrival", zap.Error(err))
			return nil, err
		}
		missingIssues = append(missingIssues, issues...)
		if resp.NextPage == 0 {
			break
		} else {
			opts.ListOptions.Page = resp.NextPage
		}
	}
	return missingIssues, nil
}

var pullGaps = func(client *github.Client, owner, name string, dbPullNum int) ([]*github.PullRequest, error) {
	ctx := context.Background()
	pull, _, err := client.PullRequests.List(
		ctx,
		owner,
		name,
		&github.PullRequestListOptions{
			State: "all",
			ListOptions: github.ListOptions{
				PerPage: 1,
			},
		},
	)
	if err != nil {
		utils.AppLog.Error(
			"newest GitHub pull request retrival",
			zap.Error(err),
		)
		return nil, err
	}
	githubPullNum := 0
	if len(pull) > 0 {
		githubPullNum = *pull[0].Number
	}

	diff := githubPullNum - dbPullNum
	missingPulls := []*github.PullRequest{}
	for diff > 0 {
		opts := github.PullRequestListOptions{
			State:       "all",
			ListOptions: github.ListOptions{},
		}
		switch {
		case diff > 0 && diff <= 100:
			opts.ListOptions.PerPage = diff
			diff = 0
		case diff > 100:
			opts.ListOptions.PerPage = 100
			diff = diff - 100
		}
		pulls, resp, err := client.PullRequests.List(ctx, owner, name, &opts)
		if err != nil {
			utils.AppLog.Error(
				"missing GitHub pull request retrival",
				zap.Error(err),
			)
			return nil, err
		}
		missingPulls = append(missingPulls, pulls...)
		if resp.NextPage == 0 {
			break
		} else {
			opts.ListOptions.Page = resp.NextPage
		}
	}
	return missingPulls, nil
}

// Restart starts the server and looks for new objects missed durring downtime.
func (i *IngestorServer) Restart() error {
	integrations, err := i.Database.ReadIntegrations()
	if err != nil {
		return err
	}

	for _, integration := range integrations {
		client := NewClient(integration.AppID, integration.InstallationID)
		repo, _, err := client.Repositories.GetByID(context.Background(), integration.RepoID)
		if err != nil {
			utils.AppLog.Error("restart get repo by id", zap.Error(err))
			return err
		}

		owner := *repo.Owner.Login
		name := *repo.Name

		dbIssueNum, dbPullNum, err := i.Database.restartCheck(
			restartQuery,
			integration.RepoID,
		)
		if err != nil {
			return err
		}

		if dbIssueNum == 0 && dbPullNum == 0 {
			authRepo := AuthenticatedRepo{
				Repo:   repo,
				Client: client,
			}
			i.RepoInitializer.AddRepo(authRepo)
			utils.AppLog.Info("initializing new repo", zap.String("repo name", *repo.Name))
			return nil
		}

		missingIssues, err := issueGaps(client, owner, name, dbIssueNum)
		if err != nil {
			return err
		}
		missingPulls, err := pullGaps(client, owner, name, dbPullNum)
		if err != nil {
			return err
		}

		// This is a fix for a deficiency in the GitHub API.
		for j := 0; j < len(missingIssues); j++ {
			missingIssues[j].Repository = repo
		}
		i.Database.BulkInsertIssuesPullRequests(missingIssues, missingPulls)
	}
	utils.AppLog.Info("successful ingestor server restart")
	return nil
}
