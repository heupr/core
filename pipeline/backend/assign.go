package backend

import (
	"context"

	"github.com/google/go-github/github"
)

func (bs *BackendServer) AssignContributor(assignee string, issue github.IssuesEvent) error {
	owner := *issue.Repo.Owner.Login
	repo := *issue.Repo.Name
	repoID := *issue.Repo.ID
	number := *issue.Issue.Number
	_, _, err := bs.Repos.Actives[repoID].Client.Issues.AddAssignees(context.Background(), owner, repo, number, []string{assignee})
	if err != nil {
		return err
	}
	return nil
}
