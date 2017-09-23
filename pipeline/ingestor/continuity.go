package ingestor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"core/pipeline/frontend"
	"core/utils"
)

const CONTINUITY_QUERY = `
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
	results, err := i.Database.db.Query(CONTINUITY_QUERY)
	if err != nil {
		utils.AppLog.Error("continuity check query", zap.Error(err))
		return nil, nil, err
	}
	defer results.Close()

	db, err := bolt.Open(utils.Config.BoltDBPath, 0644, nil)
	if err != nil {
		utils.AppLog.Error("failed opening bolt on continuity check", zap.Error(err))
		return nil, nil, err
	}
	defer db.Close()
	boltDB := frontend.BoltDB{DB: db}

	issues := []*github.Issue{}
	pulls := []*github.PullRequest{}

	for results.Next() {
		repoID, startNum, endNum, is_pull := new(int), new(int), new(int), new(bool)
		if err := results.Scan(repoID, startNum, endNum, is_pull); err != nil {
			utils.AppLog.Error("continuity check row scan", zap.Error(err))
			return nil, nil, err
		}
		tokenByte, err := boltDB.Retrieve("token", *repoID)
		if err != nil {
			utils.AppLog.Error("retrieve token continuity check", zap.Error(err))
		}

		token := oauth2.Token{}
		if err := json.Unmarshal(tokenByte, &token); err != nil {
			utils.AppLog.Error("converting tokens", zap.Error(err))
			return nil, nil, err
		}
		client := NewClient(token)

		repo, _, err := client.Repositories.GetByID(context.Background(), *repoID)
		if err != nil {
			utils.AppLog.Error("ingestor restart get by id", zap.Error(err))
			return nil, nil, err
		}
		owner := repo.Owner.Login
		name := repo.Name

		for j := *startNum + 1; j < *endNum; j++ {
			if *is_pull {
				pull, _, err := client.PullRequests.Get(context.Background(), *owner, *name, j)
				if err != nil {
					return nil, nil, err
				}
				pulls = append(pulls, pull)
			} else {
				issue, _, err := client.Issues.Get(context.Background(), *owner, *name, j)
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

// Periodically ensure that data contained in MemSQL is contiguous.
func (i *IngestorServer) Continuity() {
	ticker := time.NewTicker(time.Second * 30) // TEMPORARY
	// This chan is being kept as a means for thread-safe graceful shutdowns
	// and could be eventually passed as an argument into Continuity().
	ender := make(chan bool)
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
}
