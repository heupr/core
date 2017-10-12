package ingestor

import (
	"context"

	"github.com/google/go-github/github"
	"go.uber.org/zap"

	"core/utils"
)

const RESTART_QUERY = `
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


func (i *IngestorServer) Restart() error {

	integrations, err := i.Database.ReadIntegrations()
	if err != nil {
		utils.AppLog.Error("Retrieve bulk tokens on ingestor restart", zap.Error(err))
	}

	for _, integration := range integrations {
		client := NewClient(integration.AppId, integration.InstallationId)

		repo, _, err := client.Repositories.GetByID(context.Background(), integration.RepoId)
		if err != nil {
			utils.AppLog.Error("Ingestor restart get by id", zap.Error(err))
			return err
		}

		owner := repo.Owner.Login
		name := repo.Name

		iOld, pOld, iNew, pNew := new(int), new(int), new(int), new(int)
		rows, err := i.Database.db.Query(RESTART_QUERY, integration.RepoId, integration.RepoId)
		if err != nil {
			utils.AppLog.Error("restart query: ", zap.Error(err))
		}
		defer rows.Close()
		for rows.Next() {
			number := new(int)
			is_pull := new(bool)
			if err := rows.Scan(number, is_pull); err != nil {
				utils.AppLog.Error("restart scan: ", zap.Error(err))
			}
			switch *is_pull {
			case false:
				*iOld = *number
			case true:
				*pOld = *number
			}
		}

		if *iOld == 0 && *pOld == 0 {
			authRepo := AuthenticatedRepo{
				Repo:   repo,
				Client: client,
			}
			i.RepoInitializer.AddRepo(authRepo)
			return nil
		}

		issue, _, err := client.Issues.ListByRepo(context.Background(), *owner, *name, &github.IssueListByRepoOptions{
			State: "all",
			ListOptions: github.ListOptions{
				PerPage: 1,
			},
		})
		if err != nil {
			utils.AppLog.Error("newest issue retrival", zap.Error(err))
			return err
		}
		if len(issue) > 0 {
			iNew = issue[0].Number
		}

		iDiff := *iNew - *iOld
		missingIssues := []*github.Issue{}
		for iDiff > 0 {
			opts := github.IssueListByRepoOptions{
				State:       "all",
				ListOptions: github.ListOptions{},
			}
			switch {
			case iDiff > 0 && iDiff <= 100:
				opts.ListOptions.PerPage = iDiff
				iDiff = 0
			case iDiff > 100:
				opts.ListOptions.PerPage = 100
				iDiff = iDiff - 100
			}
			issues, resp, err := client.Issues.ListByRepo(context.Background(), *owner, *name, &opts)
			if err != nil {
				utils.AppLog.Error("newest issue retrival", zap.Error(err))
			}
			missingIssues = append(missingIssues, issues...)
			if resp.NextPage == 0 {
				break
			} else {
				opts.ListOptions.Page = resp.NextPage
			}
		}

		pull, _, err := client.PullRequests.List(context.Background(), *owner, *name, &github.PullRequestListOptions{
			State: "all",
			ListOptions: github.ListOptions{
				PerPage: 1,
			},
		})
		if err != nil {
			utils.AppLog.Error("newest pull request retrival", zap.Error(err))
		}

		if len(pull) > 0 {
			pNew = pull[0].Number
		}

		pDiff := *pNew - *pOld
		missingPulls := []*github.PullRequest{}
		for pDiff > 0 {
			opts := github.PullRequestListOptions{
				State:       "all",
				ListOptions: github.ListOptions{},
			}
			switch {
			case pDiff > 0 && pDiff <= 100:
				opts.ListOptions.PerPage = pDiff
				pDiff = 0
			case pDiff > 100:
				opts.ListOptions.PerPage = 100
				pDiff = pDiff - 100
			}
			pulls, resp, err := client.PullRequests.List(context.Background(), *owner, *name, &opts)
			if err != nil {
				utils.AppLog.Error("newest pull request retrival", zap.Error(err))
			}
			missingPulls = append(missingPulls, pulls...)
			if resp.NextPage == 0 {
				break
			} else {
				opts.ListOptions.Page = resp.NextPage
			}
		}

		for j := 0; j < len(missingIssues); j++ {
			missingIssues[j].Repository = repo
		}
		i.Database.BulkInsertIssuesPullRequests(missingIssues, missingPulls)
	}
	return nil
}
