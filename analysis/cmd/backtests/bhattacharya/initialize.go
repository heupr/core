package main

import (
	//"bytes"
	"coralreefci/engine/gateway"
	conf "coralreefci/engine/gateway/conflation"
	"coralreefci/models"
	"coralreefci/utils"
	. "github.com/ahmetalpbalkan/go-linq"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	//"runtime/debug"
)

type TestContext struct {
	Model models.Model
}

type BackTestRunner struct {
	Context TestContext
}

func (t *BackTestRunner) Run() {

	defer func() {
		//utils.Log.Error("Panic Recovered: ", recover(), bytes.NewBuffer(debug.Stack()).String())
	}()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "23fc398670a80700b19b1ae1587825a16aa8ce57"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	newGateway := gateway.CachedGateway{Gateway: &gateway.Gateway{Client: client}, DiskCache: &gateway.DiskCache{}}

	githubIssues, err := newGateway.GetIssues("dotnet", "corefx")
	if err != nil {
		utils.Log.Error("Cannot get Issues from Github Gateway. ", err)
	}
	githubPulls, err := newGateway.GetPullRequests("dotnet", "corefx")
	if err != nil {
		utils.Log.Error("Cannot get PullRequests from Github Gateway. ", err)
	}

	context := &conf.Context{}

	scenarios := []conf.Scenario{&conf.Scenario3{}}
	//scenarios := []conf.Scenario{&conf.ScenarioAND{Scenarios: []conf.Scenario{&conf.Scenario3{}}}}

	conflationAlgorithms := []conf.ConflationAlgorithm{&conf.ComboAlgorithm{Context: context}}
	normalizer := conf.Normalizer{Context: context}
	conflator := conf.Conflator{Scenarios: scenarios, ConflationAlgorithms: conflationAlgorithms, Normalizer: normalizer, Context: context}

	issuesCopy := make([]github.Issue, len(githubIssues))
	pullsCopy := make([]github.PullRequest, len(githubPulls))

	// TODO: Evaluate this particular snippet of code as it has potential
	//       performance optimization capabilities related to the hardware
	//       level. This may ultimately live in the actual gateway.go file to
	//	     improve the actual download operations.
	for i := 0; i < len(issuesCopy); i++ {
		issuesCopy[i] = *githubIssues[i]
	}
	for i := 0; i < len(pullsCopy); i++ {
		pullsCopy[i] = *githubPulls[i]
	}

	conflator.Context.Issues = []conf.ExpandedIssue{}
	conflator.SetIssueRequests(issuesCopy)
	conflator.SetPullRequests(pullsCopy)
	conflator.Conflate()

	trainingSet := []conf.ExpandedIssue{}

	for i := 0; i < len(conflator.Context.Issues); i++ {
		expandedIssue := conflator.Context.Issues[i]
		if expandedIssue.Conflate {
			if expandedIssue.Issue.Assignee == nil {
				continue
			} else {
				trainingSet = append(trainingSet, conflator.Context.Issues[i])
			}
		}
	}
	utils.ModelSummary.Info("Training set size (before Linq): ", len(trainingSet))
	processedTrainingSet := []conf.ExpandedIssue{}

	excludeAssignees := From(trainingSet).Where(func(exclude interface{}) bool {
		assignee := *exclude.(conf.ExpandedIssue).Issue.Assignee.Login
		return assignee != "dotnet-bot" && assignee != "dotnet-mc-bot" && assignee != "00101010b" && assignee != "stephentoub"
	})

	groupby := excludeAssignees.GroupBy(
		func(r interface{}) interface{} {
			return *r.(conf.ExpandedIssue).Issue.Assignee.ID
		}, func(r interface{}) interface{} {
			return r.(conf.ExpandedIssue)
		})

	where := groupby.Where(func(groupby interface{}) bool {
		return len(groupby.(Group).Group) >= 30
	})

	orderby := where.OrderByDescending(func(where interface{}) interface{} {
		return len(where.(Group).Group)
	}).ThenBy(func(where interface{}) interface{} {
		return where.(Group).Key
	})

	orderby.SelectMany(func(orderby interface{}) Query {
		return From(orderby.(Group).Group).OrderBy(
			func(where interface{}) interface{} {
				return *where.(conf.ExpandedIssue).Issue.ID
			}).Query
	}).ToSlice(&processedTrainingSet)

	Shuffle(processedTrainingSet, int64(5))

	utils.ModelSummary.Info("Backtest model training...")
	utils.ModelSummary.Info("Training set size: ", len(processedTrainingSet))

	//scoreTwo := t.Context.Model.TwoFold(processedTrainingSet)
	//utils.ModelSummary.Info("TWO FOLD:", scoreTwo)

	//scoreTen := t.Context.Model.TenFold(processedTrainingSet)
	//utils.ModelSummary.Info"TEN FOLD:", scoreTen)

	scoreJohn := t.Context.Model.JohnFold(processedTrainingSet)
	utils.ModelSummary.Info("John Fold score: ", scoreJohn)

}
