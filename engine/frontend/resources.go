package frontend

import "github.com/google/go-github/github"

func listRepositories(client *github.Client) ([]*github.Repository, error) {
	opts := &github.RepositoryListOptions{
		Type: "all",
	}
	// TODO: Evaluate ^ to ensure that this is an exhaustive list (e.g. does
	//       not need pagination).
	repos, _, err := client.Repositories.List("", opts)
	if err != nil {
		return nil, err
	}
	return repos, nil
}
