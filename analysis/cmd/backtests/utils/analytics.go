package main

import (
	"coralreefci/engine/gateway"
	conf "coralreefci/engine/gateway/conflation"
	"coralreefci/models"
	"coralreefci/models/bhattacharya"
	"coralreefci/utils"
	"flag"
	. "github.com/ahmetalpbalkan/go-linq"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
)

//This program will allow us to reload a backtest trained model and query against the model
//TODO: Add Model DiskCache (See WriteTo in naive_bayes.go)
//TODO: Load Cached Model From Disk (See NewClassifierFromReader in naive_bayes.go)
func main() {
	model := flag.String("Model", "", "a string")
	issueNumber := flag.Int("IssueNumber", 12079, "a int")
	flag.Parse()

	if *model == "" {
		log.Fatal("Please specify a valid Model. Example ./analytics -Model JohnFold9.model -IssueId 12345")
	}

	nbModel := models.Model{Algorithm: &bhattacharya.NBModel{}}
	err := nbModel.RecoverModelFromFile(*model)
	if err != nil {
		log.Fatal(err)
	}

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
	conflationAlgorithms := []conf.ConflationAlgorithm{&conf.ComboAlgorithm{Context: context}}
	normalizer := conf.Normalizer{Context: context}
	conflator := conf.Conflator{Scenarios: scenarios, ConflationAlgorithms: conflationAlgorithms, Normalizer: normalizer, Context: context}

	issuesCopy := make([]github.Issue, len(githubIssues))
	pullsCopy := make([]github.PullRequest, len(githubPulls))
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
	processedTrainingSet := []conf.ExpandedIssue{}
	From(trainingSet).Where(func(include interface{}) bool {
		return *include.(conf.ExpandedIssue).Issue.Number == *issueNumber
	}).Select(func(i interface{}) interface{} {
		return i.(conf.ExpandedIssue)
	}).ToSlice(&processedTrainingSet)

	utils.ModelSummary.Info("# of Test Issues: ", len(processedTrainingSet))
	for i := 0; i < len(processedTrainingSet); i++ {
		nbModel.Predict(processedTrainingSet[i])
	}
}
