package backend

import (
	"context"

	language "cloud.google.com/go/language/apiv1"

	"core/models"
	"core/models/bhattacharya"
	"core/models/labelmaker"
	"core/pipeline/gateway/conflation"
	"core/utils"

	"go.uber.org/zap"
)

func (s *Server) NewModel(repoID int64) error {
	s.Repos.Lock()
	defer s.Repos.Unlock()
	confCxt := &conflation.Context{}
	scenarios := []conflation.Scenario{&conflation.Scenario2{}, &conflation.Scenario3{}, &conflation.Scenario7{}}
	algos := []conflation.ConflationAlgorithm{
		&conflation.ComboAlgorithm{Context: confCxt},
	}
	normalizer := conflation.Normalizer{Context: confCxt}
	conflator := conflation.Conflator{
		Scenarios:            scenarios,
		ConflationAlgorithms: algos,
		Normalizer:           normalizer,
		Context:              confCxt,
	}
	s.Repos.Actives[repoID].Hive.Blender.Conflator = &conflator
	model := models.Model{Algorithm: &bhattacharya.NBModel{}}
	s.Repos.Actives[repoID].Hive.Blender.Models = append(
		s.Repos.Actives[repoID].Hive.Blender.Models,
		&ArchModel{Model: &model},
	)

	ctx := context.Background()
	client, err := language.NewClient(ctx)
	if err != nil {
		utils.AppLog.Error("NewModel() language client", zap.Error(err))
		return err
	}

	s.Repos.Actives[repoID].Labelmaker = &labelmaker.LBModel{
		Classifier: &labelmaker.LBClassifier{
			Client: client,
			Gateway: labelmaker.CachedNlpGateway{
				NlpGateway: &labelmaker.NlpGateway{
					Client: client,
				},
			},
			Ctx: ctx,
		},
		BugLabel:         s.Repos.Actives[repoID].Settings.Bug,
		ImprovementLabel: s.Repos.Actives[repoID].Settings.Improvement,
		FeatureLabel:     s.Repos.Actives[repoID].Settings.Feature,
	}
	return nil
}
