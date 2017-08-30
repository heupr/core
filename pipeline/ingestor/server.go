package ingestor

import (
	"context"
	"net/http"

	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"core/pipeline/frontend"
	"core/pipeline/gateway"
	"core/utils"
)

type IngestorServer struct {
	Server          http.Server
	Database        Database
	RepoInitializer RepoInitializer
}

// This is a global variable for unit testing and stubbing out the client URLs.
var makeClient = func(token oauth2.Token) *github.Client {
	source := oauth2.StaticTokenSource(&token)
	githubClient := github.NewClient(oauth2.NewClient(oauth2.NoContext, source))
	return githubClient
}

// var fetchGitHub = func(owner, name string, client github.Client) ([]*github.Issue, []*github.PullRequest, *github.Repository, error) {
// 	issues := []*github.Issue{}
// 	pulls := []*github.PullRequest{}
//
// 	isssueOpts := github.IssueListByRepoOptions{ListOptions: github.ListOptions{PerPage: 100}}
// 	pullsOpts := github.PullRequestListOptions{ListOptions: github.ListOptions{PerPage: 100}}
//
// 	for {
// 		gotIssues, resp, err := client.Issues.ListByRepo(context.Background(), owner, name, &isssueOpts)
// 		if err != nil {
// 			return nil, nil, nil, err
// 		}
// 		issues = append(issues, gotIssues...)
// 		if resp.NextPage == 0 {
// 			break
// 		} else {
// 			isssueOpts.ListOptions.Page = resp.NextPage
// 		}
// 	}
//
// 	for {
// 		gotPulls, resp, err := client.PullRequests.List(context.Background(), owner, name, &pullsOpts)
// 		if err != nil {
// 			return nil, nil, nil, err
// 		}
// 		pulls = append(pulls, gotPulls...)
// 		if resp.NextPage == 0 {
// 			break
// 		} else {
// 			pullsOpts.ListOptions.Page = resp.NextPage
// 		}
// 	}
//
// 	repo, _, err := client.Repositories.Get(context.Background(), owner, name)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
//
// 	return issues, pulls, repo, nil
// }

func (i *IngestorServer) activateHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != frontend.BackendSecret {
		errMsg := "failed validating frontend-backend secret"
		utils.AppLog.Error(errMsg)
		http.Error(w, errMsg, http.StatusUnauthorized)
		return
	}
	repoInfo := r.FormValue("repos")
	owner := string(repoInfo[1])
	repoName := string(repoInfo[2])
	tokenString := r.FormValue("token")

	client := makeClient(oauth2.Token{AccessToken: tokenString})

	gateway := gateway.Gateway{
		Client:      client,
		UnitTesting: false,
	}

	issues, err := gateway.GetIssues(owner, repoName)
	if err != nil {
		utils.AppLog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pulls, err := gateway.GetPullRequests(owner, repoName)
	if err != nil {
		utils.AppLog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// NOTE: This may ultimately be refactored out into a helper
	// method/function. Also see the similar code in the Restart method.
	repo, _, err := client.Repositories.Get(context.Background(), owner, repoName)
	if err != nil {
		utils.AppLog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		utils.AppLog.Error("ingestor github pulldown:", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	authRepo := AuthenticatedRepo{
		Repo:   repo,
		Client: client,
	}
	i.RepoInitializer = RepoInitializer{}
	i.RepoInitializer.AddRepo(authRepo)

	i.Database.BulkInsertIssues(issues)
	i.Database.BulkInsertPullRequests(pulls)
}

func (i *IngestorServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/hook", collectorHandler())
	mux.HandleFunc("/activate-ingestor-backend", i.activateHandler)
	return mux
}

func (i *IngestorServer) Start() error {
	bufferPool := NewPool()
	i.Database = Database{BufferPool: bufferPool}
	i.Database.Open()

	i.RepoInitializer = RepoInitializer{}
	i.Server = http.Server{Addr: "127.0.0.1:8030", Handler: i.routes()}
	err := i.Server.ListenAndServe()
	if err != nil {
		utils.AppLog.Error("ingestor server failed to start; ", zap.Error(err))
		return err
	}
	return nil
}

func (i *IngestorServer) Stop() {
	//TODO: Closing the server down is a needed operation that will be added.
}
