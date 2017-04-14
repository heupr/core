package onboarder

import (
	"github.com/google/go-github/github"
)

const secretKey = "chalmun"

func (h *RepoServer) NewHook(repo []*github.Repository, client *github.Client) error {
	for _, r := range repo {
		if check, err := h.hookExists(r, client); check {
			// handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			// TODO: Logic for handling an error here will be implemented; this
			//       will take the form of an exit from the parent NewHook method
			//       as well as a generation of an error/redirect page option to
			//       the end user of the Heupr application.
			// return handler, err
			// NOTE: Is this correct?
			return err
		}
		name := *r.Name
		owner := *r.Owner.Login
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
		err = h.Database.store(*r.ID, "hookID", *hook.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *RepoServer) hookExists(repo *github.Repository, client *github.Client) (bool, error) {
	name, owner := "", ""
	if repo.Name != nil {
		name = *repo.Name
	}
	if repo.Owner.Login != nil {
		owner = *repo.Owner.Login
	}
	hookID, err := h.Database.retrieve(*repo.ID, "hookID")
	if err != nil {
		return false, err
	}

	_, _, err = client.Repositories.GetHook(owner, name, hookID.(int))
	if err != nil {
		return false, err
	}
	return true, nil
}
