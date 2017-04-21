package ingestor

import (
	"coralreefci/engine/gateway"
	"github.com/google/go-github/github"
	"net/url"
)

const localPath = "http://localhost:8000/"

type RepoInitializer struct {
	repos map[int]bool
}

func (r *RepoInitializer) LoadRepos() {
	repos := getTestRepos()
	for i := 0; i < len(repos); i++ {
		r.AddRepo(repos[i])
	}
}

func (r *RepoInitializer) AddRepo(repo AuthenticatedRepo) {
	db := Database{}
	db.Open()
	newGateway := gateway.Gateway{Client: repo.Client}
	githubIssues, _ := newGateway.GetIssues(*repo.Repo.Organization.Name, *repo.Repo.Name)
	//We have to add the Repo on the Issue Body (Bug in github api?)
	for i := 0; i < len(githubIssues); i++ {
		githubIssues[i].Repository = repo.Repo
	}
	db.BulkInsertIssues(githubIssues)
}

func getTestRepos() []AuthenticatedRepo {
	client := github.NewClient(nil)
	url, _ := url.Parse(localPath)
	client.BaseURL = url
	client.UploadURL = url
	return []AuthenticatedRepo{AuthenticatedRepo{Repo: &github.Repository{ID: github.Int(26295345), Organization: &github.Organization{Name: github.String("dotnet")}, Name: github.String("coreclr")}, Client: client}}
}
