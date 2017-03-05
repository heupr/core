package onboarder

import (
	"coralreefci/models"

	"github.com/google/go-github/github"
)

type ArchModel struct {
	Model *models.Model
	// Benchmark        Benchmark // TODO: Struct to build
	// Scenarios        []conflation.Scenario
	// PilotScenarios   []conflation.Scenario
	// LearnedScenarios []conflation.Scenario
	// StrategyParams   StrategyParams TODO: Baseline with self evolving
	//                                       parameters (Tossing Graph?,
	//                                       Conflation Scenarios, etc.)
}

type ArchHive struct {
	Models []*ArchModel
	// PilotModels        []ArchModel
	// ModelBlender       ModelBlender // TODO: Struct to build
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

// TODO: This method will place the desired model(s) for a given repository ID
//       into the respective field e.g.:
//       (h.Hive.Models = append(h.Hive.Models, exampleModel))
//       Additionally, this method will need to be able to properly call in the
//       necessary logic for defining how each model(s) is trained (although it
//       is quite likely that they will all train on all available information
//       which will be filtered down by the conflation scenarios).
func (rs *RepoServer) initModels(id int) {

}

/*
// TODO: This method will need to change substantially in the switch to gob.
func (rs *RepoServer) InitArchRepos(path ...string) {
	defer h.CloseDB()
	h.OpenDB()
	if path == nil {
		// tokens, err := h.Database.retrieveBulk()
		// if err != nil {
        //
		// }
	} else {

	}

	// TODO:
	// if string is nil
	// - default to opening the storage.db file
	// - load in info from file
	// - boot up models / repos via that information
	// else
	// - open the provided file
	// - boot up from those specifications
}
*/

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
