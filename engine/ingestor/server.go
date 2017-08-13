package ingestor

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"coralreefci/engine/frontend"
	"coralreefci/utils"
)

type IngestorServer struct {
	Server          http.Server
	Database        Database
	RepoInitializer RepoInitializer
}

func (i *IngestorServer) activateHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != frontend.BackendSecret {
		utils.AppLog.Error("failed validating frontend-backend secret")
		return
	}
	repoInfo := r.FormValue("repos")
	// repoID, err := strconv.Atoi(string(repoInfo[0]))
	// if err != nil {
	// 	utils.AppLog.Error("converting repo ID: ", zap.Error(err))
	// 	http.Error(w, "failed converting repo ID", http.StatusForbidden)
	// 	return
	// }
	owner := string(repoInfo[1])
	repo := string(repoInfo[2])
	tokenString := r.FormValue("token")

	source := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tokenString})
	token := oauth2.NewClient(oauth2.NoContext, source)
	client := *github.NewClient(token)

	isssueOpts := github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	issues := []*github.Issue{}
	for {
		gotIssues, resp, err := client.Issues.ListByRepo(context.Background(), owner, repo, &isssueOpts)
		if err != nil {
			utils.AppLog.Error("failed issue pull down: ", zap.Error(err))
			http.Error(w, "failed issue pull down", http.StatusForbidden)
			return
		}
		issues = append(issues, gotIssues...)
		if resp.NextPage == 0 {
			break
		} else {
			isssueOpts.ListOptions.Page = resp.NextPage
		}
	}

	pullsOpts := github.PullRequestListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	pulls := []*github.PullRequest{}
	for {
		gotPulls, resp, err := client.PullRequests.List(context.Background(), owner, repo, &pullsOpts)
		if err != nil {
			utils.AppLog.Error("failed pull request pull down; ", zap.Error(err))
			http.Error(w, "failed pull request pull down", http.StatusForbidden)
			return
		}
		pulls = append(pulls, gotPulls...)
		if resp.NextPage == 0 {
			break
		} else {
			pullsOpts.ListOptions.Page = resp.NextPage
		}
	}

	i.Database.BulkInsertIssues(issues)
	i.Database.BulkInsertPullRequests(pulls)
}

func (i *IngestorServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/hook", collectorHandler())
	mux.HandleFunc("/activate-repos-ingestor", i.activateHandler)
	return mux
}

func (i *IngestorServer) Start() {
	bufferPool := NewPool()
	i.Database = Database{BufferPool: bufferPool}
	i.Database.Open()

	i.RepoInitializer = RepoInitializer{}
	i.Server = http.Server{Addr: "127.0.0.1:8030", Handler: i.routes()}
	err := i.Server.ListenAndServe()
	if err != nil {
		utils.AppLog.Error("ingestor server failed to start; ", zap.Error(err))
	}
}

func tokenizer(tokenByte []byte) (*github.Client, error) {
	token := oauth2.Token{}
	if err := json.Unmarshal(tokenByte, &token); err != nil {
		utils.AppLog.Error("converting tokens; ", zap.Error(err))
		return nil, err
	}

	source := oauth2.StaticTokenSource(&token)
	oaClient := oauth2.NewClient(oauth2.NoContext, source)
	client := github.NewClient(oaClient)
	return client, nil
}

const RESTART_QUERY = `
SELECT MAX(ghe.number)
FROM (
    SELECT repo_id, number
    FROM github_events
    WHERE is_pull = false
    ORDER BY repo_id, number
) ghe
WHERE repo_id = ?

UNION ALL

SELECT MAX(ghe.number)
FROM (
    SELECT repo_id, number
    FROM github_events
    WHERE is_pull = true
    ORDER BY repo_id, number
) ghe
WHERE repo_id = ?
`

func (i *IngestorServer) Restart() error {
	bufferPool := NewPool()
	i.Database = Database{BufferPool: bufferPool}
	i.Database.Open()
	defer i.Database.Close()

	db, err := bolt.Open("frontend/storage.db", 0644, nil)
	if err != nil {
		utils.AppLog.Error("failed opening bolt on ingestor restart; ", zap.Error(err))
		return err
	}
	defer db.Close()

	boltDB := frontend.BoltDB{DB: db}

	repos, tokens, err := boltDB.RetrieveBulk("token")
	if err != nil {
		utils.AppLog.Error("retrieve bulk tokens on ingestor restart; ", zap.Error(err))
	}

	for key := range tokens {
		client, err := tokenizer(tokens[key])
		if err != nil {
			return err
		}

		repoID, err := strconv.Atoi(string(repos[key]))
		if err != nil {
			utils.AppLog.Error("repo id int conversion; ", zap.Error(err))
			return err
		}

		repo, _, err := client.Repositories.GetByID(context.Background(), repoID)
		if err != nil {
			utils.AppLog.Error("ingestor restart get by id; ", zap.Error(err))
			return err
		}

		owner := repo.Owner.Login
		name := repo.Name

		iOldest, pOldest, iNewest, pNewest := new(int), new(int), new(int), new(int)
		result := i.Database.db.QueryRow(RESTART_QUERY, repoID, repoID).Scan(&iOldest, &pOldest)
		switch {
		case result == sql.ErrNoRows:
			utils.AppLog.Error("no rows in restart query; ", zap.Error(result))
			break
		case result != nil:
			utils.AppLog.Error("restart query; ", zap.Error(result))
		default:
			continue
		}

		if iOldest == nil && pOldest == nil {
			authRepo := AuthenticatedRepo{
				Repo:   repo,
				Client: client,
			}
			i.RepoInitializer = RepoInitializer{}
			i.RepoInitializer.AddRepo(authRepo)
		}
		if iOldest == nil {
			iOldest = iNewest
		}
		if pOldest == nil {
			pOldest = pNewest
		}

		issue, _, err := client.Issues.ListByRepo(context.Background(), *owner, *name, &github.IssueListByRepoOptions{
			ListOptions: github.ListOptions{
				PerPage: 1,
			},
		})
		if err != nil {
			utils.AppLog.Error("newest issue retrival; ", zap.Error(err))
		} else {
			iNewest = issue[0].Number
		}

		iDiff := *iNewest - *iOldest
		missingIssues := []*github.Issue{}
		for iDiff > 1 {
			opts := github.IssueListByRepoOptions{
				ListOptions: github.ListOptions{},
			}
			switch {
			case iDiff > 1 && iDiff <= 100:
				opts.ListOptions.PerPage = iDiff
				iDiff = 0
			case iDiff > 100:
				opts.ListOptions.PerPage = 100
				iDiff = iDiff - 100
			}
			issues, resp, err := client.Issues.ListByRepo(context.Background(), *owner, *name, &opts)
			if err != nil {
				utils.AppLog.Error("newest issue retrival; ", zap.Error(err))
			}
			missingIssues = append(missingIssues, issues...)
			if resp.NextPage == 0 {
				break
			} else {
				opts.ListOptions.Page = resp.NextPage
			}
		}

		pull, _, err := client.PullRequests.List(context.Background(), *owner, *name, &github.PullRequestListOptions{
			ListOptions: github.ListOptions{
				PerPage: 1,
			},
		})
		if err != nil {
			utils.AppLog.Error("newest pull request retrival; ", zap.Error(err))
		} else {
			pNewest = pull[0].Number
		}

		pDiff := *pNewest - *pOldest
		missingPulls := []*github.PullRequest{}
		for pDiff > 1 {
			opts := github.PullRequestListOptions{
				ListOptions: github.ListOptions{},
			}
			switch {
			case pDiff > 1 && pDiff <= 100:
				opts.ListOptions.PerPage = pDiff
				pDiff = 0
			case pDiff > 100:
				opts.ListOptions.PerPage = 100
				pDiff = pDiff - 100
			}
			pulls, resp, err := client.PullRequests.List(context.Background(), *owner, *name, &opts)
			if err != nil {
				utils.AppLog.Error("newest pull request retrival; ", zap.Error(err))
			}
			missingPulls = append(missingPulls, pulls...)
			if resp.NextPage == 0 {
				break
			} else {
				opts.ListOptions.Page = resp.NextPage
			}
		}
		i.Database.BulkInsertIssues(missingIssues)
		i.Database.BulkInsertPullRequests(missingPulls)
	}
	return nil
}

func (i *IngestorServer) Stop() {
	//TODO: Closing the server down is a needed operation that will be added.
}
