package backtests

import (
	"coralreefci/engine/gateway"
	"coralreefci/engine/gateway/conflation"
	"coralreefci/models/bhattacharya"
	"coralreefci/models/issues"
	"fmt"
	. "github.com/ahmetalpbalkan/go-linq"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type TestContext struct {
	Model bhattacharya.Model
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

	trainingSet := []issues.Issue{}
	processedTrainingSet := []issues.Issue{}

	for i := 0; i < len(conflator.Context.Issues); i++ {
		expandedIssue := conflator.Context.Issues[i]
		if expandedIssue.PullRequest.Number != nil && expandedIssue.Conflate {
			truncatedIssue := issues.Issue{
				RepoID:   *expandedIssue.PullRequest.ID,
				IssueID:  *expandedIssue.PullRequest.Number,
				URL:      *expandedIssue.PullRequest.URL,
				Assignee: *expandedIssue.PullRequest.User.Login,
			}
			if expandedIssue.PullRequest.Body != nil {
				truncatedIssue.Body = *expandedIssue.PullRequest.Body
			}
			trainingSet = append(trainingSet, truncatedIssue)
		} else if expandedIssue.Issue.Number != nil && expandedIssue.Conflate {
			truncatedIssue := issues.Issue{
				RepoID:   *expandedIssue.Issue.ID,
				IssueID:  *expandedIssue.Issue.Number,
				URL:      *expandedIssue.Issue.URL,
				Assignee: *expandedIssue.Issue.User.Login,
				Resolved: *expandedIssue.Issue.ClosedAt,
			}
			if expandedIssue.Issue.Body != nil {
				truncatedIssue.Body = *expandedIssue.Issue.Body
			}
			trainingSet = append(trainingSet, truncatedIssue)
		}
	}
	// TODO: remove this workaround eventually
	// BOTS: dotnet-bot, dotnet-mc-bot, 00101010b
	// PROJECT MANAGERS: stephentoub
	// excludeAssignees := []string{"dotnet-bot", "dotnet-mc-bot", "00101010b", "stephentoub"}
	// fileData := readFile(filePath, excludeAssignees)
	// TODO: implement a linq statement to remove the excludeAssignees
	// referecne: excludeAssignees := []string{"dotnet-bot", "dotnet-mc-bot", "00101010b", "stephentoub"}

	groupby := From(trainingSet).GroupBy(
		func(r interface{}) interface{} { return r.(issues.Issue).Assignee },
		func(r interface{}) interface{} { return r.(issues.Issue) })

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

	bhattacharya.Shuffle(processedTrainingSet, int64(5))

	// logger := bhattacharya.CreateLog("backtest-summary")
	// logger.Log("NUMBER OF ASSIGNEES:" + string(len(distinctAssignees(trainingSet))))
	// logger.Log("NUMBER OF ISSUES:" + string(len(trainingSet)))

	fmt.Println("Train")
	fmt.Println(len(processedTrainingSet))

	scoreJohn, _ := t.Context.Model.JohnFold(trainingSet)
	fmt.Println("JOHN FOLD:", scoreJohn)

	scoreTen, _ := t.Context.Model.TenFold(trainingSet)
	fmt.Println("TEN FOLD:", scoreTen)

	scoreTwo, _ := t.Context.Model.TwoFold(trainingSet)
	fmt.Println("TWO FOLD:", scoreTwo)
}
