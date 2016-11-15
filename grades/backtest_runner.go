package grades

import (
	"coralreefci/models/bhattacharya"
	"coralreefci/models/issues"
	"encoding/csv"
	"fmt"
	. "github.com/ahmetalpbalkan/go-linq"
	"os"
)

type TestContext struct {
	File  string
	Model bhattacharya.Model
}

type BackTestRunner struct {
	Context TestContext
}

func readFile(filePath string, exclude []string) []issues.Issue {
	csvData, _ := os.Open(filePath)
	defer csvData.Close()
	reader := csv.NewReader(csvData)
	var repoIssues []issues.Issue
	fmt.Printf("LOADING: %s.......\n", filePath)
	for {
		rec, _ := reader.Read()
		if rec != nil {
			skipRecord := false
			for i := 0; i < len(exclude); i++ {
				if rec[4] == exclude[i] {
					skipRecord = true
					break
				}
			}
			if !skipRecord {
				i := issues.Issue{Url: rec[0], Body: rec[3], Assignee: rec[4]}
				repoIssues = append(repoIssues, i)
			}
		} else {
			break
		}
	}
	fmt.Println("LOADING COMPLETE")
	return repoIssues
}

func distinctAssignees(issues []issues.Issue) []string {
	result := []string{}
	j := 0
	for i := 0; i < len(issues); i++ {
		for j = 0; j < len(result); j++ {
			if issues[i].Assignee == result[j] {
				break
			}
		}
		if j == len(result) {
			result = append(result, issues[i].Assignee)
		}
	}
	return result
}

func (t *BackTestRunner) Run() {
	filePath := t.Context.File

	data := HistoricalData{}
	data.Download()
	// TODO: remove this workaround eventually
	// BOTS: dotnet-bot, dotnet-mc-bot, 00101010b
	// PROJECT MANAGERS: stephentoub
	excludeAssignees := []string{"dotnet-bot", "dotnet-mc-bot", "00101010b", "stephentoub"}
	fileData := readFile(filePath, excludeAssignees)

	trainingSet := []issues.Issue{}
	assignees := []string{}

	groupby := From(fileData).GroupBy(
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

	orderby.Select(func(orderby interface{}) interface{} {
		return orderby.(Group).Key
	}).ToSlice(&assignees)

	bhattacharya.Shuffle(trainingSet, int64(5))

	logger := bhattacharya.CreateLog("backtest-summary")
	logger.Log("NUMBER OF ASSIGNEES:" + string(len(distinctAssignees(trainingSet))))
	logger.Log("NUMBER OF ISSUES:" + string(len(trainingSet)))

    scoreJohn, _ := t.Context.Model.JohnFold(trainingSet)
    logger.Log("JOHN FOLD: " + scoreJohn)
    scoreTwo, _ := t.Context.Model.TwoFold(trainingSet)
    logger.Log("TWO FOLD: " + scoreTwo)
    scoreTen, _ := t.Context.Model.TenFold(trainingSet)
    logger.Log("TEN FOLD: " + scoreTen)

	logger.Flush()
}
