package backend

import (
	"testing"
)

const repoID = 66

var testBS = &Server{
	Repos: &ActiveRepos{Actives: make(map[int64]*ArchRepo)},
}

func TestNewArchRepo(t *testing.T) {
	hcs := HeuprConfigSettings{}
	testBS.NewArchRepo(repoID, hcs)
	if testBS.Repos.Actives[repoID] == nil {
		t.Error("failure generating new arch repo")
	}
}
