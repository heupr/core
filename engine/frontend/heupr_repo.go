package frontend

import (
	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"
	"github.com/google/go-github/github"
)

type HeuprRepo struct {
	Repo   github.Repository
	Model  ModelHive
	Client github.Client
}

func (h *HeuprRepo) BootstrapModel() {

}

func (h *HeuprRepo) GetModelBenchmark() {
	//TODO:Calculate AggregateBenchmark
}

func (h *HeuprRepo) Assign(issue github.Issue) {

}

func (h *HeuprRepo) Enable() {

}

func (h *HeuprRepo) Disable() {

}

func (h *HeuprRepo) Destroy() {

}

func (h *HeuprRepo) InactiveDevelopers() []string {
	//Even though the Github API will prevent it. Heupr could try to assign an issue to an inactive developer.
	//It might not be great for team morale if Heupr inadvertently exposed that in a UI/ Dashboard page.
}

type ModelHive struct {
	Models             []HeuprModel
	PilotModels        []HeuprModel
	ModelBlender       ModelBlender
	TossingGraph       TossingGraphAlgorithm
	StrategyParams     StrategyParams //Baseline with self evolving parameters (Tossing Graph?, Conflation Scenarios)
	AggregateBenchmark Benchmark
}

type HeuprModel struct {
	Model            models.Model
	Benchmark        Benchmark
	Scenarios        []conflation.Scenario
	PilotScenarios   []conflation.Scenario
	LearnedScenarios []conflation.Scenario
	StrategyParams   StrategyParams //Baseline with self evolving parameters (Tossing Graph?, Conflation Scenarios)
}

type Benchmark struct{}
type StrategyParams struct{}
type ModelBlender struct{}
type TossingGraphAlgorithm struct{}
