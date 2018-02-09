package backend

import "testing"

func TestNewModel(t *testing.T) {
	repoID := int64(7)
	testBS := new(BackendServer)
	testBS.Repos = new(ActiveRepos)
	testBS.Repos.Actives = make(map[int64]*ArchRepo)
	testBS.Repos.Actives[repoID] = new(ArchRepo)
	testBS.Repos.Actives[repoID].Hive = new(ArchHive)
	testBS.Repos.Actives[repoID].Hive.Blender = new(Blender)

	if err := testBS.NewModel(repoID); err != nil {
		t.Errorf("error adding model to test backendserver: %v", err)
	}
	if len(testBS.Repos.Actives[repoID].Hive.Blender.Models) == 0 {
		t.Error("model not added to slice test backendserver")
	}
}
