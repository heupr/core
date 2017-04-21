package onboarder

import (
	"github.com/google/go-github/github"

	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"
	"coralreefci/models/bhattacharya"
)

func (rs *RepoServer) AddModel(repo *github.Repository) error {
	repoID := *repo.ID
	context := &conflation.Context{}
	scenarios := []conflation.Scenario{&conflation.Scenario2{}}
	algos := []conflation.ConflationAlgorithm{
		&conflation.ComboAlgorithm{Context: context},
	}
	normalizer := conflation.Normalizer{Context: context}
	conflator := conflation.Conflator{
		Scenarios:            scenarios,
		ConflationAlgorithms: algos,
		Normalizer:           normalizer,
		Context:              context,
	}
	model := models.Model{Algorithm: &bhattacharya.NBModel{}}
	if rs.Repos == nil {
		rs.Repos = make(map[int]*ArchRepo)
	}
	if _, ok := rs.Repos[repoID]; !ok {
		rs.Repos[repoID] = &ArchRepo{Hive: &ArchHive{Blender: &Blender{}}}
	}
	rs.Repos[repoID].Hive.Blender.Models = append(rs.Repos[repoID].Hive.Blender.Models, &ArchModel{Model: &model, Conflator: &conflator})
	return nil
}
