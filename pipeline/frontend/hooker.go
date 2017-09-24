package frontend

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

const secretKey = "figrin-dan-and-the-modal-nodes"

func (fs *FrontendServer) NewHook(repo *github.Repository, client *github.Client) error {
	if check, err := fs.hookExists(repo, client); check {
		return errors.Wrap(err, "hook exists")
	}
	name := *repo.Name
	owner := *repo.Owner.Login
	url := "http://35.196.33.44/hook"
	hook, _, err := client.Repositories.CreateHook(context.Background(), owner, name, &github.Hook{
		Name:   github.String("web"),
		Events: []string{"issues", "pull_request"},
		Active: github.Bool(true),
		Config: map[string]interface{}{
			"url":          url,
			"secret":       secretKey,
			"content_type": "json",
			"insecure_ssl": false,
		},
	})
	if err != nil {
		return errors.Wrap(err, "adding new hook to repo")
	}
	if err = fs.Database.Store("hook", *repo.ID, []byte(strconv.Itoa(*hook.ID))); err != nil {
		return errors.Wrap(err, "error storing hook info")
	}
	return nil
}

func (fs *FrontendServer) hookExists(repo *github.Repository, client *github.Client) (bool, error) {
	name, owner := "", ""
	if repo.Name != nil && repo.Owner.Login != nil {
		name = *repo.Name
		owner = *repo.Owner.Login
	}
	hook, err := fs.Database.Retrieve("hook", *repo.ID)
	if err != nil {
		return false, errors.Wrap(err, "error retrieving hook info")
	}

	hookID, err := strconv.Atoi(string(hook))
	if err != nil {
		return false, errors.Wrap(err, "failed hook string conversion")
	}
	_, _, err = client.Repositories.GetHook(context.Background(), owner, name, hookID)
	if err != nil {
		return false, errors.Wrap(err, "error getting GitHub hook info")
	}
	return true, fmt.Errorf("hook previously established for repo %v", name)
}
