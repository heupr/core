package grades

import (
	"coralreef-ci/models/bhattacharya"
	"coralreef-ci/models/issues"
	"encoding/csv"
	"fmt"
	"os"
)

type TestContext struct {
	File  string
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

func (t *BackTestRunner) Run() {
	filePath := t.Context.File
	trainingSet := readFile(filePath)
	testComparsionSet := make([]issues.Issue, len(trainingSet))
	testSet := make([]issues.Issue, len(trainingSet))
	copy(testSet, trainingSet)
	copy(testComparsionSet, trainingSet)
	t.Context.Model.Learn(trainingSet)
	correctCount := 0
	for i := 0; i < len(trainingSet); i++ {
		assignees := t.Context.Model.Predict(testSet[i])
		/*  testSet[i].Assignee = assignees[0]
		    fmt.Println("BackTest Assignee: ", testSet[i].Assignee)
		    fmt.Println("Expected Assignee: ", testComparsionSet[i].Assignee)
		    if (testComparsionSet[i].Assignee == assignees[0]) {
		      correctCount++
		      //break
		    } */

		for j := 0; j < len(assignees); j++ {
			testSet[i].Assignee = assignees[j]
			if testComparsionSet[i].Assignee == assignees[j] {
				correctCount++
				break
			}
		}
	}

	matrix, _ := bhattacharya.BuildMatrix(trainingSet, testSet)
	bhattacharya.FullSummary(matrix)

	fmt.Printf("\nACCURACY: %f\n", float64(correctCount)/float64(len(testSet)))
	fmt.Printf("\nCORRECT COUNT: %d\n", correctCount)
	fmt.Printf("\nTOTAL COUNT: %d\n", len(testSet))
}
