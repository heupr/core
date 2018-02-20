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
	"core/models/labelmaker"
	"core/pipeline/gateway/conflation"
)

type ArchModel struct {
	sync.Mutex
	Model *models.Model
}

type Blender struct {
	// NOTE: It might make more sense to make this a map[string]*ArchModel
	// so that we can identify models based on a string name rather than
	// a slice index. Just a thought.
	Models    []*ArchModel
	Conflator *conflation.Conflator
	// MVP: Moving the Conflator from ArchModel to Blender. We might just
	// need to circle back to this.
}
type ArchHive struct {
	Blender *Blender
}

type ArchRepo struct {
	sync.Mutex
	Hive                     *ArchHive
	Labelmaker               *labelmaker.LBModel
	Labels                   []string
	Client                   *github.Client
	Limit                    time.Time
	AssigneeAllocations      map[string]int
	EligibleAssignees        map[string]int
	Settings                 HeuprConfigSettings
	TriagedLabelEnabledCheck bool          //TEMPORARY FIX
	TriagedLabel             *github.Label //TEMPORARY FIX
	TriagedLabelEnabled      bool          //TEMPORARY FIX
}

func (s *Server) NewArchRepo(repoID int64, settings HeuprConfigSettings) {
	s.Repos.Lock()
	defer s.Repos.Unlock()

	s.Repos.Actives[repoID] = new(ArchRepo)
	s.Repos.Actives[repoID].Hive = new(ArchHive)
	s.Repos.Actives[repoID].Hive.Blender = new(Blender)
	s.Repos.Actives[repoID].Settings = settings
}

func (s *Server) NewClient(repoID int64, appID int, installationID int64) {
	s.Repos.Lock()
	defer s.Repos.Unlock()

	var key string
	if PROD {
		key = "heupr.2017-10-04.private-key.pem"
	} else {
		key = "mikeheuprtest.2017-11-16.private-key.pem"
	}
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, appID, int(installationID), key)
	if err != nil {
		utils.AppLog.Error("could not obtain github installation key", zap.Error(err))
		return
	}
	client := github.NewClient(&http.Client{Transport: itr})

	s.Repos.Actives[repoID].Client = client
}

func (a *ArchRepo) ApplyLabelsOnOpenIssues() {
	openIssues := a.Hive.Blender.GetAllOpenIssues()
	utils.AppLog.Info("ApplyLabelsOnOpenIssues()", zap.Int("Total", len(openIssues)))
	if len(openIssues) == 0 {
		return
	}
	name := ""
	if openIssues[0].Issue.Repository.FullName != nil {
		name = *openIssues[0].Issue.Repository.FullName
	} else {
		name = *openIssues[0].Issue.Repository.Name
	}
	repo := strings.Split(name, "/")

	skip := false
	for i := 0; i < len(openIssues); i++ {
		if openIssues[i].Issue.CreatedAt.After(a.Settings.StartTime) {
			*openIssues[i].Issue.Labeled = true 
			label, err := a.Labelmaker.BugOrFeature(openIssues[i])
			if err != nil {
				utils.AppLog.Error("AddLabelsToIssue Failed", zap.Error(err))
				break
			}
			for j := 0; j < len(openIssues[i].Issue.Labels); j++ {
				if *openIssues[i].Issue.Labels[j].Name != "triaged" {
					skip = true
				}
			}
			if skip == true {
				skip = false
				continue
			}
			if label != nil && *label != "" {
				_, _, err := a.Client.Issues.AddLabelsToIssue(context.Background(), repo[0], repo[1], *openIssues[i].Issue.Number, []string{*label})
				if err != nil {
					utils.AppLog.Error("AddLabelsToIssue Failed", zap.Error(err))
					break
				}
			}
		}
	}
}

func (a *ArchRepo) TriageOpenIssues() {
	if !a.Settings.EnableTriager {
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

	//TEMPORARY FIX
	var label *github.Label
	if a.TriagedLabelEnabledCheck == false {
		a.TriagedLabelEnabledCheck = true
		lbl, _, err := a.Client.Issues.GetLabel(context.Background(), r[0], r[1], "triaged")
		if err != nil {
			utils.AppLog.Error("could not get triaged label", zap.String("RepoName", r[0]+"/"+r[1]))
		}
		a.TriagedLabel = lbl
	}
	label = a.TriagedLabel

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

							if label != nil {
								if *label.Name == "triaged" {
									_, _, err := a.Client.Issues.AddLabelsToIssue(context.Background(), r[0], r[1], number, []string{*label.Name})
									if err != nil {
										utils.AppLog.Error("AddLabelsToIssue failed for primary assignee", zap.Error(err))
									}
								}
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
				if label != nil {
					if *label.Name == "triaged" {
						_, _, err := a.Client.Issues.AddLabelsToIssue(context.Background(), r[0], r[1], number, []string{*label.Name})
						if err != nil {
							utils.AppLog.Error("AddLabelsToIssue failed for fallback assignee", zap.Error(err))
						}
					}
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
