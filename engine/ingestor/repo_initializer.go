package ingestor

import (
	"coralreefci/engine/gateway"
)

type RepoInitializer struct {
	repos map[int]bool
}

func (r *RepoInitializer) LoadRepos() {

}

func (r *RepoInitializer) AddRepo(repo AuthenticatedRepo) {
	bufferPool := NewPool()
	db := Database{BufferPool: bufferPool}
	db.Open()
	newGateway := gateway.Gateway{Client: repo.Client}
	githubIssues, _ := newGateway.GetIssues(*repo.Repo.Organization.Name, *repo.Repo.Name)
	//We have to add the Repo on the Issue Body (Bug in github api?)
	for i := 0; i < len(githubIssues); i++ {
		githubIssues[i].Repository = repo.Repo
	}
	githubPulls, _ := newGateway.GetPullRequests(*repo.Repo.Organization.Name, *repo.Repo.Name)
	db.BulkInsertIssues(githubIssues)
	db.BulkInsertPullRequests(githubPulls)
}
