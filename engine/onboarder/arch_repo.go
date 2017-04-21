package onboarder

import (
	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"

	"github.com/google/go-github/github"
)

type ArchModel struct {
	Model     *models.Model
	Conflator *conflation.Conflator
	// Benchmark        Benchmark // TODO: Struct to build
	// Scenarios        []conflation.Scenario
	// PilotScenarios   []conflation.Scenario
	// LearnedScenarios []conflation.Scenario
	// StrategyParams   StrategyParams TODO: Baseline with self evolving
	//                                       parameters (Tossing Graph?,
	//                                       Conflation Scenarios, etc.)
}

type Blender struct {
	Models []*ArchModel
	// PilotModels []*ArchModel
}

type ArchHive struct {
	Blender *Blender
	// TossingGraph       TossingGraphAlgorithm // TODO: Struct to build
	// StrategyParams     StrategyParams // TODO: Struct to build
	// AggregateBenchmark Benchmark
}

type ArchRepo struct {
	Repo   *github.Repository
	Hive   *ArchHive
	Client *github.Client
}

func (rs *RepoServer) NewArchRepo(repo *github.Repository, client *github.Client) {
	rs.Repos[*repo.ID] = &ArchRepo{
		Repo:   repo,
		Client: client,
	}
}

// TODO: Instantiate the Conflator struct on the ArchRepo.

// TODO:
// Below are several potential helper methods for the ArchRepo:
// BootstrapModel() - performs preliminary training / assignments / startup
// GetModelBenchmark() TODO: Calculate AggregateBenchmark for this method
// Assign(issue github.Issue) - assign newly raised issue to contributor
// Enable()
// Disable()
// Destroy()
// InactiveDevelopers() []string
// Even though the Github API will prevent it. Heupr could try to assign an
// issue to an inactive developer. It might not be great for team morale if
// Heupr inadvertently exposed that in a UI/dashboard page.
