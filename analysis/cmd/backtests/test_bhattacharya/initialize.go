package main

import (
	"coralreefci/engine/gateway"
	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"
	"coralreefci/models/bhattacharya"
	"fmt"
	. "github.com/ahmetalpbalkan/go-linq"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type TestContext struct {
	Model models.Model
}

type BackTestRunner struct {
	Context TestContext
}

func (t *BackTestRunner) Run() {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "23fc398670a80700b19b1ae1587825a16aa8ce57"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	newGateway := gateway.CachedGateway{Gateway: &gateway.Gateway{Client: client}, DiskCache: &gateway.DiskCache{}}

	githubIssues, _ := newGateway.GetIssues("dotnet", "corefx")
	githubPulls, _ := newGateway.GetPullRequests("dotnet", "corefx")

	context := &conflation.Context{}

	scenarios := []conflation.Scenario{&conflation.Scenario1b{}, &conflation.Scenario2b{}, &conflation.Scenario3a{}}
	conflationAlgorithms := []conflation.ConflationAlgorithm{&conflation.ComboAlgorithm{Context: context}}
	normalizer := conflation.Normalizer{Context: context}
	conflator := conflation.Conflator{Scenarios: scenarios, ConflationAlgorithms: conflationAlgorithms, Normalizer: normalizer, Context: context}

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

	conflator.Context.Issues = []conflation.ExpandedIssue{}
	conflator.SetIssueRequests(issuesCopy)
	conflator.SetPullRequests(pullsCopy)
	conflator.Conflate()

	trainingSet := []conflation.ExpandedIssue{}
	trainingSet = append(trainingSet, conflator.Context.Issues...)
	processedTrainingSet := []conflation.ExpandedIssue{}

	groupby := From(trainingSet).GroupBy(
		func(r interface{}) interface{} { return r.(bhattacharya.Issue).Assignee },
		func(r interface{}) interface{} { return r.(bhattacharya.Issue) })

	where := groupby.Where(func(groupby interface{}) bool {
		return len(groupby.(Group).Group) >= 10
	})

	orderby := where.OrderByDescending(func(where interface{}) interface{} {
		return len(where.(Group).Group)
	})

	orderby.SelectMany(func(orderby interface{}) Query {
		return From(orderby.(Group).Group)
	}).ToSlice(&processedTrainingSet)

	// orderby.Select(func(orderby interface{}) interface{} {
	// 	return orderby.(Group).Key
	// }).ToSlice(&assignees)

	Shuffle(processedTrainingSet, int64(5))

	fmt.Println("Train")
	fmt.Println(len(processedTrainingSet))

    // NOTE: should this be the processedTrainingSet instead of trainingSet?
	scoreJohn := t.Context.Model.JohnFold(trainingSet)
	fmt.Println("JOHN FOLD:", scoreJohn)

	scoreTen := t.Context.Model.TenFold(trainingSet)
	fmt.Println("TEN FOLD:", scoreTen)

	scoreTwo := t.Context.Model.TwoFold(trainingSet)
	fmt.Println("TWO FOLD:", scoreTwo)
}
