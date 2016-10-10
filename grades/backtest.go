package grades

import (
  "encoding/csv"
  "fmt"
  "os"
	"coralreef-ci/models/bhattacharya"
  "coralreef-ci/models/issues"
)

type TestContext struct {
  File string
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
	fmt.Printf("\n\nLoading %s.......\n", filePath)
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
      if (!skipRecord) {
        i := issues.Issue{Body: rec[3], Assignee: rec[4]}
        repoIssues = append(repoIssues, i)
      }
		} else {
			break
		}
	}
	fmt.Println("Loading Complete")
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

func(t *BackTestRunner) Run() {
  filePath := t.Context.File
  trainingSet := readFile(filePath, []string{"dotnet-bot", "dotnet-mc-bot", "00101010b", "stephentoub"})
  testComparsionSet := make([]issues.Issue, len(trainingSet))
  testSet := make([]issues.Issue, len(trainingSet))
  copy(testSet, trainingSet)
  copy(testComparsionSet, trainingSet)
  t.Context.Model.Learn(trainingSet)

  score,_ := t.Context.Model.TwoFold(trainingSet, 1)
  fmt.Println("Graph Length", 1)
  fmt.Println("Two Fold Weighted Accuracy:", score)
  score,_ = t.Context.Model.TwoFold(trainingSet, 2)
  fmt.Println("Graph Length:", 2)
  fmt.Println("Two Fold Weighted Accuracy:", score)
  score,_ = t.Context.Model.TwoFold(trainingSet, 3)
  fmt.Println("Graph Length:", 3)
  fmt.Println("Two Fold Weighted Accuracy:", score)
}
