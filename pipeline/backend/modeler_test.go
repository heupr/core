package backend

import (
	"context"
	"testing"

	language "cloud.google.com/go/language/apiv1"
)

func TestNewModel(t *testing.T) {
	repoID := int64(7)
	testBS := new(Server)
	testBS.Repos = new(ActiveRepos)
	testBS.Repos.Actives = make(map[int64]*ArchRepo)
	testBS.Repos.Actives[repoID] = new(ArchRepo)
	testBS.Repos.Actives[repoID].Hive = new(ArchHive)
	testBS.Repos.Actives[repoID].Hive.Blender = new(Blender)

	NewLanguageClient = func(ctx context.Context) (*language.Client, error) {
		return &language.Client{}, nil
	}

	testBS.NewModel(repoID)
	if len(testBS.Repos.Actives[repoID].Hive.Blender.Models) == 0 {
		t.Error("model not added to slice test backendserver")
	}
}
