package listener

import (
	"github.com/google/go-github/github"
)

func HeuprHook(owner, repo string, client *github.Client) error {
	name := "heupr-" + repo
	// config := make(map["url"]"delivery-url-address")
	hook := &github.Hook{Name: &name}
	_, _, err := client.Repositories.CreateHook(owner, repo, hook)
	if err != nil {
		return err
	}
	return nil
}
