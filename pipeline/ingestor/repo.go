package ingestor

import (
	"github.com/google/go-github/github"
)

type AuthenticatedRepo struct {
	Repo   *github.Repository
	Client *github.Client
}
