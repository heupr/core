package backend

import (
	"core/models"
	"core/models/bhattacharya"
	"core/pipeline/gateway/conflation"
)

func (bs *BackendServer) NewModel(repoID int64) error {
	bs.Repos.Lock()
	defer bs.Repos.Unlock()
	context := &conflation.Context{}
	scenarios := []conflation.Scenario{&conflation.Scenario2{}, &conflation.Scenario3{}, &conflation.Scenario7{}}
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
	bs.Repos.Actives[repoID].Hive.Blender.Conflator = &conflator
	bs.Repos.Actives[repoID].Hive.Blender.Models = append(bs.Repos.Actives[repoID].Hive.Blender.Models, &ArchModel{Model: &model})
	return nil
}
