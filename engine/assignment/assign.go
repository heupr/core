package assignment

import (
	"github.com/google/go-github/github"
)

func AssignContributor(assignee string, issue github.IssuesEvent, client *github.Client) error {
	// TODO: Logic to handle multiple assignees presented for the given issue.
	// TODO: Error avoidance when an issue already has an assignee on it.
	owner := *issue.Repo.Owner.Login
	repo := *issue.Repo.Name
	number := *issue.Issue.Number
	_, _, err := client.Issues.AddAssignees(owner, repo, number, []string{assignee})
	if err != nil {
		return err
	}
	return nil
}
