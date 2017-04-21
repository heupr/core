package onboarder

import (
	"github.com/google/go-github/github"

	"testing"
)

func TestAddModel(t *testing.T) {
	reposerver := new(RepoServer)
	id := 7
	repo := &github.Repository{ID: &id}
	if err := reposerver.AddModel(repo); err != nil {
		t.Error("Error adding model to the RepoServer")
	}
	if len(reposerver.Repos[id].Hive.Blender.Models) == 0 {
		t.Error("Model not added to Models slice on RepoServer")
	}
}
