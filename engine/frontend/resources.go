package frontend

import (
	"context"

	"github.com/google/go-github/github"
)

func listRepositories(client *github.Client) ([]*github.Repository, error) {
	opts := &github.RepositoryListOptions{
		Type: "all",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	allRepos := []*github.Repository{}
	for {
		repos, resp, err := client.Repositories.List(context.Background(), "", opts)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		} else {
			opts.ListOptions.Page = resp.NextPage
		}
	}
	return allRepos, nil
}
