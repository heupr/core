package ingestor

import (
	"core/pipeline/gateway"
)

type RepoInitializer struct {
	repos map[int]bool
}

func (r *RepoInitializer) LoadRepos() {

}

func (r *RepoInitializer) AddRepo(authRepo AuthenticatedRepo) {
	bufferPool := NewPool()
	db := Database{BufferPool: bufferPool}
	db.Open()
	newGateway := gateway.Gateway{Client: authRepo.Client}
	githubIssues, _ := newGateway.GetIssues(*authRepo.Repo.Organization.Name, *authRepo.Repo.Name)
	// The Repo struct needs to be added to the Issue struct body - possibly a bug in the GitHub API.
	for i := 0; i < len(githubIssues); i++ {
		githubIssues[i].Repository = authRepo.Repo
	}
	githubPulls, _ := newGateway.GetPullRequests(*authRepo.Repo.Organization.Name, *authRepo.Repo.Name)
	db.BulkInsertIssues(githubIssues)
	db.BulkInsertPullRequests(githubPulls)
}
