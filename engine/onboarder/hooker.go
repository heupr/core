package onboarder

import (
	"github.com/google/go-github/github"
)

const secretKey = "figrin-dan-and-the-modal-nodes"

func (rs *RepoServer) NewHook(repo *github.Repository, client *github.Client) error {
	if check, err := rs.hookExists(repo, client); check {
		return err
	}
	name := *repo.Name
	owner := *repo.Owner.Login
	// TODO: This URL will change to a config parameter.
	url := "http://00ad0ac7.ngrok.io/hook"
	hook, _, err := client.Repositories.CreateHook(owner, name, &github.Hook{
		Name:   github.String("web"),
		Events: []string{"issues", "repository"},
		Active: github.Bool(true),
		Config: map[string]interface{}{
			"url":          url,
			"secret":       secretKey,
			"content_type": "json",
			"insecure_ssl": false,
		},
	})
	if err != nil {
		return err
	}
	if err = rs.Database.store(*repo.ID, "hookID", *hook.ID); err != nil {
		return err
	}
	return nil
}

func (rs *RepoServer) hookExists(repo *github.Repository, client *github.Client) (bool, error) {
	name, owner := "", ""
	if repo.Name != nil && repo.Owner.Login != nil {
		name = *repo.Name
		owner = *repo.Owner.Login
	}
	hookID, err := rs.Database.retrieve(*repo.ID, "hookID")
	if err != nil {
		return false, err
	}

	_, _, err = client.Repositories.GetHook(owner, name, hookID.(int))
	if err != nil {
		return false, err
	}
	return true, nil
}
