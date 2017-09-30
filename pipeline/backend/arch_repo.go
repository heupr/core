package backend

import (
	"context"
	"core/utils"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"core/models"
	"core/pipeline/gateway/conflation"
)

type ArchModel struct {
	sync.Mutex
	Model *models.Model
}

type Blender struct {
	Models    []*ArchModel
	Conflator *conflation.Conflator
	// MVP: Moving the Conflator from ArchModel
	// to Blender. We might just need to circle back to this.
}
type ArchHive struct {
	Blender *Blender
}

type ArchRepo struct {
	sync.Mutex
	Hive                *ArchHive
	Client              *github.Client
	Limit               time.Time
	AssigneeAllocations map[string]int
	EligibleAssignees   map[string]int
}

func (bs *BackendServer) NewArchRepo(repoID int, limit time.Time) {
	bs.Repos.Lock()
	defer bs.Repos.Unlock()

	bs.Repos.Actives[repoID] = new(ArchRepo)
	bs.Repos.Actives[repoID].Hive = new(ArchHive)
	bs.Repos.Actives[repoID].Hive.Blender = new(Blender)

	bs.Repos.Actives[repoID].Limit = limit
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
		utils.AppLog.Error("!AllModelsBootstrapped()")
		return
	}

	openIssues := a.Hive.Blender.GetOpenIssues()
	utils.AppLog.Info("TriageOpenIssues()", zap.Int("Total", len(openIssues)))
	for i := 0; i < len(openIssues); i++ {
		if openIssues[i].Issue.CreatedAt.After(a.Limit) {
			assignees := a.Hive.Blender.Predict(openIssues[i])
			var name string
			if openIssues[i].Issue.Repository.FullName != nil {
				name = *openIssues[i].Issue.Repository.FullName
			} else {
				name = *openIssues[i].Issue.Repository.Name
			}
			r := strings.Split(name, "/")
			number := *openIssues[i].Issue.Number
			fallbackAssignee := ""
			assigned := false
			for i := 0; i < len(assignees); i++ {
				assignee := assignees[i]
				if assignmentsCap, ok := a.EligibleAssignees[assignee]; ok {
					if fallbackAssignee == "" {
						fallbackAssignee = assignee
					}

					if _, ok := a.AssigneeAllocations[assignee]; !ok {
						a.AssigneeAllocations[assignee] = 0
					}
					if assignmentsCount, ok := a.AssigneeAllocations[assignee]; ok {
						if assignmentsCount < assignmentsCap {
							issue, _, err := a.Client.Issues.AddAssignees(context.Background(), r[0], r[1], number, []string{assignee})
							if err != nil {
								utils.AppLog.Error("AddAssignees Failed", zap.Error(err))
								break
							}
							if issue.Assignees == nil || len(issue.Assignees) == 0 {
								if fallbackAssignee == assignee {
									fallbackAssignee = ""
								}
								continue
							}
							assigned = true
							assignmentsCount++
							a.AssigneeAllocations[assignee] = assignmentsCount
							break
						}
					}
				}
			}
			if !assigned {
				if fallbackAssignee == "" {
					utils.AppLog.Error("AddAssignees Failed. Fallback assignee not found.", zap.String("URL", *openIssues[i].Issue.URL), zap.Int("IssueID", *openIssues[i].Issue.ID))
					break
				}
				_, _, err := a.Client.Issues.AddAssignees(context.Background(), r[0], r[1], number, []string{fallbackAssignee})
				if err != nil {
					utils.AppLog.Error("AddAssignees Failed", zap.Error(err))
					break
				}
				assigned = true
				a.AssigneeAllocations[fallbackAssignee]++
				utils.AppLog.Info("AddAssignees Success. Fallback assignee found.", zap.String("URL", *openIssues[i].Issue.URL), zap.Int("IssueID", *openIssues[i].Issue.ID))
			}
		}
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
			if issues[i].Issue.Assignee == nil && issues[i].Issue.Assignees == nil { //MVP
				openIssues = append(openIssues, issues[i])
				issues[i].Issue.Triaged = true
			}
		}
	}
	return openIssues
}

func (b *Blender) GetClosedIssues() []conflation.ExpandedIssue {
	closedIssues := []conflation.ExpandedIssue{}
	issues := b.Conflator.Context.Issues
	for i := 0; i < len(issues); i++ {
		if issues[i].Issue.ClosedAt != nil && issues[i].Conflate && !issues[i].IsTrained {
			closedIssues = append(closedIssues, issues[i])
			issues[i].IsTrained = true
		}
	}
	return closedIssues
}

func (b *Blender) TrainModels() {
	closedIssues := b.GetClosedIssues()
	utils.AppLog.Info("TrainModels() ", zap.Int("Total", len(closedIssues)))
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
