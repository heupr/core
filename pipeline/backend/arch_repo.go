package backend

import (
	"context"
	"core/utils"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"go.uber.org/zap"

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
	//TODO: As an MVP Quick win you can slap the LabelMaker Model here. The Blender albeit underutilized is primarily good for bagging and or boosting. https://en.wikipedia.org/wiki/Bootstrap_aggregating
	Client              *github.Client
	Limit               time.Time
	AssigneeAllocations map[string]int
	EligibleAssignees   map[string]int
	Settings            HeuprConfigSettings
}

func (bs *BackendServer) NewArchRepo(repoID int64, settings HeuprConfigSettings) {
	bs.Repos.Lock()
	defer bs.Repos.Unlock()

	bs.Repos.Actives[repoID] = new(ArchRepo)
	bs.Repos.Actives[repoID].Hive = new(ArchHive)
	bs.Repos.Actives[repoID].Hive.Blender = new(Blender)

	bs.Repos.Actives[repoID].Settings = settings
}

func (bs *BackendServer) NewClient(repoID int64, appID, installationID int) {
	bs.Repos.Lock()
	defer bs.Repos.Unlock()

	var key string
	if PROD {
		key = "heupr.2017-10-04.private-key.pem"
	} else {
		key = "forstmeier-heupr.2017-11-18.private-key.pem"
	}
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, appID, installationID, key)
	if err != nil {
		utils.AppLog.Error("could not obtain github installation key", zap.Error(err))
		return
	}
	client := github.NewClient(&http.Client{Transport: itr})

	bs.Repos.Actives[repoID].Client = client
}

func (a *ArchRepo) ApplyLabelsOnOpenIssues() {
	openIssues := a.Hive.Blender.GetAllOpenIssues()
	utils.AppLog.Info("ApplyLabelsOnOpenIssues()", zap.Int("Total", len(openIssues)))
	if len(openIssues) == 0 {
		return
	}
	var name string
	if openIssues[0].Issue.Repository.FullName != nil {
		name = *openIssues[0].Issue.Repository.FullName
	} else {
		name = *openIssues[0].Issue.Repository.Name
	}
	r := strings.Split(name, "/")

	for i := 0; i < len(openIssues); i++ {
		if openIssues[i].Issue.CreatedAt.After(a.Settings.StartTime) {
			// TODO: Finish this...

			//Here is a contrived sample program example that might help
			/*
			featureLabel := "feature-request"
			bugLabel := "bug"
			lbModel := labelmaker.LBModel{Classifier: &labelmaker.LBClassifier{Ctx: ctx, Client: client, Gateway: NlpGateway}, FeatureLabel: &featureLabel, BugLabel: &bugLabel, ImprovementLabel: &improvementLabel}
			for i := 0; i < len(issues); i++ {
		    prediction, _ := lbModel.BugOrFeature(conf.ExpandedIssue{Issue: conf.CRIssue{*issues[i], []int{}, []conf.CRPullRequest{}, github.Bool(false), github.Bool(false)}})
		    if prediction != nil {
		      fmt.Println(*issues[i].Title)
		      fmt.Println(*prediction)
		    }
		  } */

		}
	}
}

func (a *ArchRepo) TriageOpenIssues() {
	if !a.Settings.EnableIssueAssignments {
		return
	}
	if !a.Hive.Blender.AllModelsBootstrapped() {
		utils.AppLog.Error("!AllModelsBootstrapped()")
		return
	}
	openIssues := a.Hive.Blender.GetOpenIssues()
	utils.AppLog.Info("TriageOpenIssues()", zap.Int("Total", len(openIssues)))
	if len(openIssues) == 0 {
		return
	}
	var name string
	if openIssues[0].Issue.Repository.FullName != nil {
		name = *openIssues[0].Issue.Repository.FullName
	} else {
		name = *openIssues[0].Issue.Repository.Name
	}
	r := strings.Split(name, "/")

	for i := 0; i < len(openIssues); i++ {
		if openIssues[i].Issue.CreatedAt.After(a.Settings.StartTime) {
			labelValid := true
			labels := openIssues[i].Issue.Labels
			for j := 0; j < len(labels); j++ {
				if _, ok := a.Settings.IgnoreLabels[*labels[j].Name]; ok {
					labelValid = false
					break
				}
			}
			if !labelValid {
				continue
			}
			*openIssues[i].Issue.Triaged = true
			assignees := a.Hive.Blender.Predict(openIssues[i])
			number := *openIssues[i].Issue.Number
			fallbackAssignee := ""
			assigned := false
			for i := 0; i < len(assignees); i++ {
				assignee := assignees[i]
				if _, ok := a.Settings.IgnoreUsers[assignee]; ok {
					continue
				}
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
								if _, ok := err.(*github.RateLimitError); ok {
									time.Sleep(15 * time.Minute)
									limits, _, _ := a.Client.RateLimits(context.Background())
									if limits != nil {
										limit := limits.Core.Limit
										remaining := limits.Core.Remaining
										utils.AppLog.Info("RateLimits()", zap.Int("Limit", limit), zap.Int("Remaining", remaining))
									}
									issue, _, err = a.Client.Issues.AddAssignees(context.Background(), r[0], r[1], number, []string{assignee})
									if err != nil {
										break
									}
								} else {
									break
								}
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
					utils.AppLog.Error("AddAssignees Failed. Fallback assignee not found.", zap.String("URL", *openIssues[i].Issue.URL), zap.Int64("IssueID", *openIssues[i].Issue.ID))
					break
				}
				_, _, err := a.Client.Issues.AddAssignees(context.Background(), r[0], r[1], number, []string{fallbackAssignee})
				if err != nil {
					utils.AppLog.Error("AddAssignees Failed", zap.Error(err))
					break
				}
				assigned = true
				a.AssigneeAllocations[fallbackAssignee]++
				utils.AppLog.Info("AddAssignees Success. Fallback assignee found.", zap.String("URL", *openIssues[i].Issue.URL), zap.Int64("IssueID", *openIssues[i].Issue.ID))
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
		if issues[i].PullRequest.Number == nil && issues[i].Issue.ClosedAt == nil && !*issues[i].Issue.Triaged && *issues[i].Issue.User.Login != "heupr" {
			if issues[i].Issue.Assignee == nil && issues[i].Issue.Assignees == nil { //MVP
				openIssues = append(openIssues, issues[i])
			}
		}
	}
	return openIssues
}

func (b *Blender) GetAllOpenIssues() []conflation.ExpandedIssue {
	openIssues := []conflation.ExpandedIssue{}
	issues := b.Conflator.Context.Issues
	for i := 0; i < len(issues); i++ {
		if issues[i].PullRequest.Number == nil && issues[i].Issue.ClosedAt == nil && !*issues[i].Issue.Labeled && *issues[i].Issue.User.Login != "heupr" {
			openIssues = append(openIssues, issues[i])
		}
	}
	return openIssues
}

func (b *Blender) GetClosedIssues() []conflation.ExpandedIssue {
	closedIssues := []conflation.ExpandedIssue{}
	issues := b.Conflator.Context.Issues
	for i := 0; i < len(issues); i++ {
		if issues[i].Issue.ClosedAt != nil && issues[i].Conflate && !issues[i].IsTrained && *issues[i].Issue.User.Login != "heupr" {
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
