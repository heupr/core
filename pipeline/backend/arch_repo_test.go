package backend

import (
	"testing"
	"time"
)

const repoID = 66

var testBS = &BackendServer{
	Repos: &ActiveRepos{Actives: make(map[int]*ArchRepo)},
}

func TestNewArchRepo(t *testing.T) {
	testTime := time.Time{}
	testBS.NewArchRepo(repoID, testTime)
	if testBS.Repos.Actives[repoID] == nil {
		t.Error("failure generating new arch repo")
	}
}
