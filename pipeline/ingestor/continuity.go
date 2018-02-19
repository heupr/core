package ingestor

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	"go.uber.org/zap"

	"core/utils"
)

const continuityQuery = `
SELECT *
FROM (
    SELECT repo_id, number AS start_num, LEAD(number, 1) OVER (PARTITION BY repo_id ORDER BY number) AS end_num, is_pull
    FROM github_events
  ) t
WHERE t.end_num - t.start_num > 1 AND is_pull = FALSE

UNION ALL

SELECT *
FROM (
    SELECT repo_id, number AS start_num, LEAD(number, 1) OVER (PARTITION BY repo_id ORDER BY number) AS end_num, is_pull
    FROM github_events
  ) t
WHERE t.end_num - t.start_num > 1 AND is_pull = TRUE;
`

func (i *IngestorServer) continuityCheck() ([]*github.Issue, []*github.PullRequest, error) {
	gaps, err := i.Database.continuityCheck(continuityQuery)
	if err != nil {
		utils.AppLog.Error("continuity check query", zap.Error(err))
		return nil, nil, err
	}

	issues := []*github.Issue{}
	pulls := []*github.PullRequest{}
	for j := range gaps {
		repoID := gaps[j][0].(int64)
		ctx := context.Background()
		integration, err := i.Database.ReadIntegrationByRepoID(repoID)
		if err != nil {
			utils.AppLog.Error(
				"retrieve token continuity check",
				zap.Error(err),
			)
		}

		client := NewClient(integration.AppID, integration.InstallationID)
		repo, _, err := client.Repositories.GetByID(ctx, repoID)
		if err != nil {
			utils.AppLog.Error("ingestor restart get by id", zap.Error(err))
			return nil, nil, err
		}

		owner := *repo.Owner.Login
		name := *repo.Name
		startNum, endNum := gaps[j][1].(int), gaps[j][2].(int)
		isPull := gaps[j][3].(bool)
		for k := startNum + 1; k < endNum; k++ {
			if isPull {
				pull, _, err := client.PullRequests.Get(ctx, owner, name, k)
				if err != nil {
					return nil, nil, err
				}
				pulls = append(pulls, pull)
			} else {
				issue, _, err := client.Issues.Get(ctx, owner, name, k)
				if err != nil {
					return nil, nil, err
				}
				// This is a patch in what may be an error in the GitHub API.
				issue.Repository = repo
				issues = append(issues, issue)
			}
		}
	}
	return issues, pulls, nil
}

// Continuity periodically ensures that data contained in MemSQL is contiguous.
func (i *IngestorServer) Continuity() {
	ticker := time.NewTicker(time.Second * 300) // TEMPORARY
	// This chan is being kept as a means for thread-safe graceful shutdowns
	// and could be eventually passed as an argument into Continuity().
	ender := make(chan bool)
	go func() {
		defer close(ender)
		for {
			select {
			case <-ticker.C:
				issues, pulls, err := i.continuityCheck()
				if err != nil {
					utils.AppLog.Error("continuity check", zap.Error(err))
				}
				i.Database.BulkInsertIssuesPullRequests(issues, pulls)
			case <-ender:
				ticker.Stop()
				return
			}
		}
	}()
}
