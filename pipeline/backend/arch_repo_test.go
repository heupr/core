package backend

import (
	"testing"
)

const repoID = 66

var testBS = &Server{
	Repos: &ActiveRepos{
		Actives: make(map[int64]*ArchRepo),
	},
}

func TestNewArchRepo(t *testing.T) {
	settings := HeuprConfigSettings{}
	testBS.NewArchRepo(repoID, settings)
	if testBS.Repos.Actives[repoID] == nil {
		t.Error("failure generating new arch repo")
	}
}

func TestApplyLabelsOnOpenIssues(t *testing.T) {

}
