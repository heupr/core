package models

import (
	"coralreefci/engine/gateway/conflation"
    "github.com/google/go-github/github"
	"testing"
)

var expectedList = []string{"John", "Mike", "Woz", "John", "Mike", "Woz", "John", "Mike", "Woz"}
var predictedList = []string{"John", "John", "Mike", "Woz", "Woz", "Mike", "Mike", "Mike", "John"}

var metrics = map[string]float64{
	"MikeTP":    1.0,
	"MikeTN":    3.0,
	"MikeFP":    3.0,
	"MikeFN":    2.0,
	"FullCount": 9.0,
	"Precision": 0.25,
	"Recall":    0.33,
	"Accuracy":  0.22,
}

func generateIssues(assignees []string) []conflation.ExpandedIssue {
	issueList := []conflation.ExpandedIssue{}
	for i := 0; i < len(assignees); i++ {
        name := assignees[i]
        assignee := github.User{Name: &name}
        githubIssue := github.Issue{Assignee: &assignee}
        crIssue := conflation.CRIssue{githubIssue, []int{}, []conflation.CRPullRequest{}}
		issueList = append(issueList, conflation.ExpandedIssue{Issue: crIssue})
	}
	return issueList
}

func TestBuildMatrix(t *testing.T) {
	exp := generateIssues(expectedList)
	pre := generateIssues(predictedList)

    nbModel := Model{}
	matrix, _ := nbModel.BuildMatrix(exp, pre)

	if len(matrix) == 0 {
		t.Error(
			"\nEMPTY MATRIX",
			"\nCONTENTS", matrix)
	}

	countTP := matrix.getClassTP("Mike")
	if metrics["MikeTP"] != countTP {
		t.Error(
			"\nCLASS TRUE POSITIVE MISCOUNT",
			"\nEXPECTED:  ", metrics["MikeTP"],
			"\nACTUAL:    ", countTP)
	}

	countNP := matrix.getClassTN("Mike")
	if metrics["MikeTN"] != countNP {
		t.Error(
			"\nCLASS TRUE NEGATIVE MISCOUNT",
			"\nEXPECTED:  ", metrics["MikeTN"],
			"\nACTUAL:    ", countNP)
	}

	countFP := matrix.getClassFP("Mike")
	if metrics["MikeFP"] != countFP {
		t.Error(
			"\nCLASS FALSE POSITIVE MISCOUNT",
			"\nEXPECTED:  ", metrics["MikeFP"],
			"\nACTUAL:    ", countFP)
	}

	countFN := matrix.getClassFN("Mike")
	if metrics["MikeFN"] != countFN {
		t.Error(
			"\nCLASS FALSE NEGATIVE MISCOUNT",
			"\nEXPECTED:  ", metrics["MikeFN"],
			"\nACTUAL:    ", countFN)
	}

	classPrecision := matrix.getPrecision("Mike")
	if metrics["Precision"] != classPrecision {
		t.Error(
			"\nCLASS PRECISION MISCALCULATED",
			"\nEXPECTED:  ", metrics["Precision"],
			"\nACTUAL:    ", classPrecision)
	}

	classRecall := matrix.getRecall("Mike")
	if metrics["Recall"] != classRecall {
		t.Error(
			"\nCLASS RECALL MISCALCULATED",
			"\nEXPECTED:  ", metrics["Recall"],
			"\nACTUAL:    ", classRecall)
	}

	fullAccuracy := matrix.getAccuracy()
	if metrics["Accuracy"] != fullAccuracy {
		t.Error(
			"\nALL TESTS INACCURATE",
			"\nEXPECTED:  ", metrics["Accuracy"],
			"\nACTUAL:    ", fullAccuracy)
	}

	fullCount := matrix.getTestCount()
	if metrics["FullCount"] != fullCount {
		t.Error(
			"\nALL TESTS MISCOUNT",
			"\nEXPECTED:  ", metrics["FullCount"],
			"\nACTUAL:    ", fullCount)
	}

	fullMatrix := matrix.fillMatrix()
	for key := range fullMatrix {
		if len(fullMatrix) != len(fullMatrix[key]) {
			t.Error(
				"\nMATRIX IS NOT EQUAL IN DIMENSIONS",
				"\nEXPECTED LENGTH:  ", len(fullMatrix),
				"\nACTUAL LENGTH:    ", len(fullMatrix[key]))
		}
	}

	classOutput := fullMatrix.ClassSummary("John")
	if classOutput == "" {
		t.Error(
			"\nNO OUTPUT STRING",
		)
	}

	fullOutput := fullMatrix.FullSummary()
	if fullOutput == "" {
		t.Error(
			"\nNO OUTPUT STRING",
		)
	}
}
