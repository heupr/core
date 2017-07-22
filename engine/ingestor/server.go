package ingestor

import (
	"context"
	"net/http"
	// "strconv"

	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"coralreefci/engine/frontend"
	"coralreefci/utils"
)

type IngestorServer struct {
	Server          http.Server
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
			utils.AppLog.Error("failed pull request pull down: ", zap.Error(err))
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

	bufferPool := NewPool()
	db := Database{BufferPool: bufferPool}
	db.Open()

	db.BulkInsertIssues(issues)
	db.BulkInsertPullRequests(pulls)
}

func (i *IngestorServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/hook", collectorHandler())
	mux.HandleFunc("/activate-repos-ingestor", i.activateHandler)
	return mux
}

func (i *IngestorServer) Start() {
	i.RepoInitializer = RepoInitializer{}
	i.Server = http.Server{Addr: "127.0.0.1:8030", Handler: i.routes()}
	err := i.Server.ListenAndServe()
	if err != nil {
		utils.AppLog.Error("ingestor server failed to start", zap.Error(err))
	}
}

func (i *IngestorServer) Restart() {}

func (i *IngestorServer) ContinuityCheck() {}

func (i *IngestorServer) Stop() {
	//TODO: Closing the server down is a needed operation that will be added.
}
