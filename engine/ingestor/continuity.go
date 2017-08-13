package ingestor

import (
	"context"
	"time"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	"go.uber.org/zap"

	"coralreefci/engine/frontend"
	"coralreefci/utils"
)

const CONTINUITY_QUERY = `
SELECT ghe.repo_id, ghe.num1, ghe.num2, ghe.is_pull
FROM (
    SELECT repo_id, number AS num1, LEAD(number) OVER(ORDER BY repo_id, number) AS num2
    FROM github_events
    WHERE is_pull = true
) ghe
WHERE (ghe.num2 - ghe.num1) > 1

UNION ALL

SELECT ghe.repo_id, ghe.num1, ghe.num2, ghe.is_pull
FROM (
    SELECT repo_id, number AS num1, LEAD(number) OVER(ORDER BY repo_id, number) AS num2
    FROM github_events
    WHERE is_pull = false
) ghe
WHERE (ghe.num2 - ghe.num1) > 1;
`

func (i *IngestorServer) continuityCheck(checker chan bool) ([]*github.Issue, []*github.PullRequest, error) {
	results, err := i.Database.db.Query(CONTINUITY_QUERY)
	if err != nil {
		utils.AppLog.Error("continuity check query; ", zap.Error(err))
		return nil, nil, err
	}
	defer results.Close()

	db, err := bolt.Open("frontend/storage.db", 0644, nil)
	if err != nil {
		utils.AppLog.Error("failed opening bolt on continuity check; ", zap.Error(err))
		return nil, nil, err
	}
	defer db.Close()

	boltDB := frontend.BoltDB{DB: db}

	issues := []*github.Issue{}
	pulls := []*github.PullRequest{}

	for results.Next() {
		repoID, num1, num2, is_pull := new(int), new(int), new(int), new(bool)
		if err := results.Scan(repoID, num1, num2, is_pull); err != nil {
			utils.AppLog.Error("continuity check row scan; ", zap.Error(err))
			return nil, nil, err
		}

		t, err := boltDB.Retrieve("token", *repoID)
		if err != nil {
			utils.AppLog.Error("retrieve token continuity check; ", zap.Error(err))
		}

		client, err := tokenizer(t)
		if err != nil {
			return nil, nil, err
		}

		repo, _, err := client.Repositories.GetByID(context.Background(), *repoID)
		if err != nil {
			utils.AppLog.Error("ingestor restart get by id; ", zap.Error(err))
			return nil, nil, err
		}

		owner := repo.Owner.Login
		name := repo.Name

		for i := 1; i < (*num2 - *num1); i++ {
			num := i + *num1
			if *is_pull {
				pull, _, err := client.PullRequests.Get(context.Background(), *owner, *name, num)
				if err != nil {
					return nil, nil, err
				}
				pulls = append(pulls, pull)
			} else {
				issue, _, err := client.Issues.Get(context.Background(), *owner, *name, num)
				if err != nil {
					return nil, nil, err
				}
				issues = append(issues, issue)
			}
		}
	}
	return issues, pulls, nil
}

func (i *IngestorServer) Continuity() {
	checker := make(chan bool)

	for {
		switch check := <-checker; check {
		case true:
			issues, pulls, err := i.continuityCheck(checker)
			if err != nil {
				utils.AppLog.Error("failure returning continuity check; ", zap.Error(err))
			}
			i.Database.BulkInsertIssues(issues)
			i.Database.BulkInsertPullRequests(pulls)
			continue
		case false:
			time.Sleep(10 * time.Second)
			continue
		}
	}
}
