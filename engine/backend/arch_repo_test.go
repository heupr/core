package backend

import (
	"testing"

	"golang.org/x/oauth2"
)

const repoID = 66

var testBS = &BackendServer{
	Repos: &ActiveRepos{Actives: make(map[int]*ArchRepo)},
}

func TestNewArchRepo(t *testing.T) {
	testBS.NewArchRepo(repoID)
	if testBS.Repos.Actives[repoID] == nil {
		t.Error("failure generating new arch repo")
	}
}

func TestNewClient(t *testing.T) {
	testToken := oauth2.Token{AccessToken: "test-token"}
	testBS.NewClient(repoID, &testToken)
	if testBS.Repos.Actives[repoID].Client == nil {
		t.Error("failure generating client for arch repo")
	}
}
