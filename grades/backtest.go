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

func readFile(filePath string) []issues.Issue {
	csvData, _ := os.Open(filePath)
	defer csvData.Close()
	reader := csv.NewReader(csvData)
	var repoIssues []issues.Issue
	fmt.Printf("\n\nLoading %s.......\n", filePath)
	for {
		rec, _ := reader.Read()
		if rec != nil {
			i := issues.Issue{Body: rec[3], Assignee: rec[4]}
			repoIssues = append(repoIssues, i)
		} else {
			break
		}
	}
	fmt.Println("Loading Complete")
	return repoIssues
}

func(t *BackTestRunner) Run() {
  filePath := t.Context.File
  trainingSet := readFile(filePath)
  testSet := readFile(filePath)//make([]issues.Issue, len(trainingSet))
  copy(testSet, trainingSet)
  t.Context.Model.Learn(trainingSet)
  correctCount := 0
  for i := 0; i < len(trainingSet); i++ {
    first, second, third := t.Context.Model.Predict(testSet[i])
    if testSet[i].Assignee == first {
      testSet[i].Assignee = first
      correctCount++
    } else if testSet[i].Assignee == second {
      testSet[i].Assignee = second
      correctCount++
    } else if testSet[i].Assignee == third {
      testSet[i].Assignee = third
      correctCount++
    } else {
      testSet[i].Assignee = first
    }

  }

  matrix,_ := bhattacharya.BuildMatrix(trainingSet, testSet)
  bhattacharya.FullSummary(matrix)

  fmt.Printf("\nAccuracy: %f", float64(correctCount)/float64(len(testSet)))
  fmt.Printf("\nCorrect Count: %d Total Count: %d", correctCount, len(testSet))
}
