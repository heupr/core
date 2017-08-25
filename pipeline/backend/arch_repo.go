package backend

import (
	"context"
	"core/pipeline/gateway/conflation"
	"core/models"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"strings"
	"sync"
)

type ArchModel struct {
	sync.Mutex
	Model *models.Model
}

type Blender struct {
	Models    []*ArchModel
	Conflator *conflation.Conflator // MVP: Moving the Conflator from ArchModel
} // to Blender. We might just need to circle back to this.

type ArchHive struct {
	Blender *Blender
}

type ArchRepo struct {
	sync.Mutex
	Hive   *ArchHive
	Client *github.Client
}

func (bs *BackendServer) NewArchRepo(repoID int) {
	bs.Repos.Lock()
	defer bs.Repos.Unlock()

	ctx := &conflation.Context{}
	scn := []conflation.Scenario{&conflation.Scenario3{}}
	algo := []conflation.ConflationAlgorithm{
		&conflation.ComboAlgorithm{
			Context: ctx,
		},
	}
	norm := conflation.Normalizer{Context: ctx}
	conf := conflation.Conflator{
		Scenarios:            scn,
		ConflationAlgorithms: algo,
		Normalizer:           norm,
		Context:              ctx,
	}
	model := ArchModel{}

	bs.Repos.Actives[repoID] = &ArchRepo{
		Hive: &ArchHive{
			Blender: &Blender{
				Models:    []*ArchModel{&model},
				Conflator: &conf,
			},
		},
	}
}

func (bs *BackendServer) NewClient(repoID int, token *oauth2.Token) {
	bs.Repos.Lock()
	defer bs.Repos.Unlock()

	tokenSource := oauth2.StaticTokenSource(token)
	authClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	githubClient := github.NewClient(authClient)

	bs.Repos.Actives[repoID].Client = githubClient
}

func (a *ArchRepo) TriageOpenIssues() {
	if !a.Hive.Blender.AllModelsBootstrapped() {
		//TODO: Add Logging
		return
	}
	openIssues := a.Hive.Blender.GetOpenIssues()
	for i := 0; i < len(openIssues); i++ {
		assignees := a.Hive.Blender.Predict(openIssues[i])
		openIssues[i].Issue.Triaged = true
		var name string
		if openIssues[i].Issue.Repository.FullName != nil {
			name = *openIssues[i].Issue.Repository.FullName
		} else {
			name = *openIssues[i].Issue.Repository.Name
		}
		r := strings.Split(name, "/")
		number := *openIssues[i].Issue.Number
		//HACK! (temp code)
		a.Client.Issues.AddAssignees(context.Background(), r[0], r[1], number, assignees) //[]string{assignees[0]})
	}
}

func (b *Blender) Predict(issue conflation.ExpandedIssue) []string {
	var assignees []string
	for i := 0; i < len(b.Models); i++ {
		assignees = b.Models[i].Model.Predict(issue)
	}
	return assignees
}

func (b *Blender) GetOpenIssues() []conflation.ExpandedIssue {
	openIssues := []conflation.ExpandedIssue{}
	issues := b.Conflator.Context.Issues
	for i := 0; i < len(issues); i++ {
		if issues[i].PullRequest.Number == nil && issues[i].Issue.ClosedAt == nil && !issues[i].Issue.Triaged {
			openIssues = append(openIssues, issues[i])
		}
	}
	return openIssues
}

func (b *Blender) GetClosedIssues() []conflation.ExpandedIssue {
	closedIssues := []conflation.ExpandedIssue{}
	issues := b.Conflator.Context.Issues
	for i := 0; i < len(issues); i++ {
		if issues[i].Issue.ClosedAt != nil && issues[i].Conflate {
			closedIssues = append(closedIssues, issues[i])
		}
	}
	return closedIssues
}

func (b *Blender) TrainModels() {
	closedIssues := b.GetClosedIssues()
	if len(closedIssues) == 0 {
		//TODO: Add Logging
		return
	}
	for i := 0; i < len(b.Models); i++ {
		if b.Models[i].Model.IsBootstrapped() {
			b.Models[i].Model.OnlineLearn(closedIssues)
		} else {
			b.Models[i].Model.Learn(closedIssues)
		}
	}
}

func (b *Blender) AllModelsBootstrapped() bool {
	for i := 0; i < len(b.Models); i++ {
		if !b.Models[i].Model.IsBootstrapped() {
			return false
		}
	}
	return true
}

/*
func (a *ArchRepo) TuneConflationScenarios() {
	//This method will iterate over all models in the hive and call either Learn or OnlineLearn
}

func (a *ArchRepo) GetModelBenchmark() {
 //Call this method before calling TriageOpenIssues.
}

//Every week we can generate a report that shows tossing graphs(pre) vs post signup
func (a *ArchRepo) PreHeuprTossingGraphDepth() {
	//For all Pre-Heupr assigned issues
	//Calcuate avg,min,max tossing graph depth for each developer
}

func (a *ArchRepo) TossingGraphDepth() {
	//For all Heupr assigned issues
	//Calcuate avg,min,max tossing graph depth for each developer
} */
