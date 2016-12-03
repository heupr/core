package grades

import (
	"coralreefci/engine/gateway"
	"coralreefci/engine/gateway/conflation"
	"coralreefci/models/bhattacharya"
	"coralreefci/models/issues"
	// "encoding/csv"
	"fmt"
	. "github.com/ahmetalpbalkan/go-linq"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	// "os"
)

type TestContext struct {
	Model bhattacharya.Model
}

type BackTestRunner struct {
	Context TestContext
}

// func readFile(filePath string, exclude []string) []issues.Issue {
// 	csvData, _ := os.Open(filePath)
// 	defer csvData.Close()
// 	reader := csv.NewReader(csvData)
// 	var repoIssues []issues.Issue
// 	fmt.Printf("LOADING: %s.......\n", filePath)
// 	for {
// 		rec, _ := reader.Read()
// 		if rec != nil {
// 			skipRecord := false
// 			for i := 0; i < len(exclude); i++ {
// 				if rec[4] == exclude[i] {
// 					skipRecord = true
// 					break
// 				}
// 			}
// 			if !skipRecord {
// 				i := issues.Issue{Url: rec[0], Body: rec[3], Assignee: rec[4]}
// 				repoIssues = append(repoIssues, i)
// 			}
// 		} else {
// 			break
// 		}
// 	}
// 	fmt.Println("LOADING COMPLETE")
// 	return repoIssues
// }

// func distinctAssignees(issues []issues.Issue) []string {
// 	result := []string{}
// 	j := 0
// 	for i := 0; i < len(issues); i++ {
// 		for j = 0; j < len(result); j++ {
// 			if issues[i].Assignee == result[j] {
// 				break
// 			}
// 		}
// 		if j == len(result) {
// 			result = append(result, issues[i].Assignee)
// 		}
// 	}
// 	return result
// }

func (t *BackTestRunner) Run() {
	// filePath := t.Context.File

	// data := HistoricalData{}
	// data.Download()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "23fc398670a80700b19b1ae1587825a16aa8ce57"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	// client := github.NewClient(nil)
	newGateway := gateway.CachedGateway{Gateway: &gateway.Gateway{Client: client, UnitTesting: true}, DiskCache: &gateway.DiskCache{}}

	githubIssues, _ := newGateway.GetIssues("dotnet", "corefx")
	githubPulls, _ := newGateway.GetPullRequests("dotnet", "corefx")

	context := &conflation.Context{}

	scenarios := []conflation.Scenario{&conflation.Scenario2b{}}
	conflationAlgorithms := []conflation.ConflationAlgorithm{&conflation.ComboAlgorithm{Context: context}}
	normalizer := conflation.Normalizer{Context: context}
	conflator := conflation.Conflator{Scenarios: scenarios, ConflationAlgorithms: conflationAlgorithms, Normalizer: normalizer, Context: context}

	issuesCopy := make([]github.Issue, len(githubIssues))
	pullsCopy := make([]github.PullRequest, len(githubPulls))

	// Workaround
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

	// TODO: Transform ExpandedIssue to Issue type
	// TODO: the "fileData" will be a list of the Issue type

	trainingSet := []issues.Issue{}

	for i := 0; i < len(conflator.Context.Issues); i++ {
		expandedIssue := conflator.Context.Issues[i]
		if expandedIssue.PullRequest.Number != nil {
			truncatedIssue := issues.Issue{
				RepoID:   *expandedIssue.PullRequest.ID,
				IssueID:  *expandedIssue.PullRequest.Number,
				Url:      *expandedIssue.PullRequest.URL,
				Assignee: *expandedIssue.PullRequest.User.Login,
				Body:     *expandedIssue.PullRequest.Body,
				// Resolved: *expandedIssue.PullRequest.MergedAt,
			}
			trainingSet = append(trainingSet, truncatedIssue)
		} else {
			truncatedIssue := issues.Issue{
				RepoID:   *expandedIssue.Issue.ID,
				IssueID:  *expandedIssue.Issue.Number,
				Url:      *expandedIssue.Issue.URL,
				Assignee: *expandedIssue.Issue.User.Login,
				Body:     *expandedIssue.Issue.Body,
				Resolved: *expandedIssue.Issue.ClosedAt,
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
	}).ToSlice(&trainingSet)

	// orderby.Select(func(orderby interface{}) interface{} {
	// 	return orderby.(Group).Key
	// }).ToSlice(&assignees)

	bhattacharya.Shuffle(trainingSet, int64(5))

	// logger := bhattacharya.CreateLog("backtest-summary")
	// logger.Log("NUMBER OF ASSIGNEES:" + string(len(distinctAssignees(trainingSet))))
	// logger.Log("NUMBER OF ISSUES:" + string(len(trainingSet)))

	scoreJohn, _ := t.Context.Model.JohnFold(trainingSet)
	// logger.Log("JOHN FOLD: " + scoreJohn)
	scoreTwo, _ := t.Context.Model.TwoFold(trainingSet)
	// logger.Log("TWO FOLD: " + scoreTwo)
	scoreTen, _ := t.Context.Model.TenFold(trainingSet)
	// logger.Log("TEN FOLD: " + scoreTen)

	fmt.Println(scoreJohn)
	fmt.Println(scoreTwo)
	fmt.Println(scoreTen)

	// logger.Flush()
}
