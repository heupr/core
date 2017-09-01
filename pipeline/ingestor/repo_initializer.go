package ingestor

import (
	"fmt"

	"core/pipeline/gateway"
	"core/utils"
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
	newGateway := gateway.Gateway{
		Client:      authRepo.Client,
		UnitTesting: false,
	}
	githubIssues, err := newGateway.GetIssues(*authRepo.Repo.Owner.Login, *authRepo.Repo.Name)
	if err != nil {
		// TODO: Proper error handling should be evaluated for this method;
		// possibly adjust to return an error variable.
		utils.AppLog.Error(err.Error())
		fmt.Println(err)
	}
	// The Repo struct needs to be added to the Issue struct body - this is
	// possibly a bug in the GitHub API.
	for i := 0; i < len(githubIssues); i++ {
		githubIssues[i].Repository = authRepo.Repo
	}
	githubPulls, err := newGateway.GetPullRequests(*authRepo.Repo.Owner.Login, *authRepo.Repo.Name)
	if err != nil {
		utils.AppLog.Error(err.Error())
		fmt.Println(err)
	}
	db.BulkInsertIssues(githubIssues)
	db.BulkInsertPullRequests(githubPulls)
}
