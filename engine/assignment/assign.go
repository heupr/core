package assignment

import (
	"github.com/google/go-github/github"
)

func AssignContributor(assignee string, issue github.Issue, client *github.Client) error {
	// TODO: Logic to handle multiple assignees presented for the given issue.
	// TODO: Error avoidance when an issue already has an assignee on it.
	// HACK: Hardcoded owner and repo
	// TODO: issue.Repository.Owner.Login and issue.Repository.Name are missing from the issue
	owner := "heupr" //*issue.Repository.Owner.Login
	repo := "test"   //*issue.Repository.Name
	number := *issue.Number
	_, _, err := client.Issues.AddAssignees(owner, repo, number, []string{assignee})
	if err != nil {
		return err
	}
	return nil
}
