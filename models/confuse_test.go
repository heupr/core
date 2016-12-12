package confuse

import (
	"coralreefci/models/issues"
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

func generateIssues(assignees []string) []issues.Issue {
	issueList := []issues.Issue{}
	for i := 0; i < len(assignees); i++ {
		issueList = append(issueList, issues.Issue{Assignee: assignees[i]})
	}
	return issueList
}

func TestBuildMatrix(t *testing.T) {
	exp := generateIssues(expectedList)
	pre := generateIssues(predictedList)

	matrix, _ := BuildMatrix(exp, pre)

	if len(matrix) == 0 {
		t.Error(
			"\nEMPTY MATRIX",
			"\nCONTENTS", matrix)
	}

	countTP := getClassTP("Mike", matrix)
	if metrics["MikeTP"] != countTP {
		t.Error(
			"\nCLASS TRUE POSITIVE MISCOUNT",
			"\nEXPECTED:  ", metrics["MikeTP"],
			"\nACTUAL:    ", countTP)
	}

	countNP := getClassTN("Mike", matrix)
	if metrics["MikeTN"] != countNP {
		t.Error(
			"\nCLASS TRUE NEGATIVE MISCOUNT",
			"\nEXPECTED:  ", metrics["MikeTN"],
			"\nACTUAL:    ", countNP)
	}

	countFP := getClassFP("Mike", matrix)
	if metrics["MikeFP"] != countFP {
		t.Error(
			"\nCLASS FALSE POSITIVE MISCOUNT",
			"\nEXPECTED:  ", metrics["MikeFP"],
			"\nACTUAL:    ", countFP)
	}

	countFN := getClassFN("Mike", matrix)
	if metrics["MikeFN"] != countFN {
		t.Error(
			"\nCLASS FALSE NEGATIVE MISCOUNT",
			"\nEXPECTED:  ", metrics["MikeFN"],
			"\nACTUAL:    ", countFN)
	}

	classPrecision := getPrecision("Mike", matrix)
	if metrics["Precision"] != classPrecision {
		t.Error(
			"\nCLASS PRECISION MISCALCULATED",
			"\nEXPECTED:  ", metrics["Precision"],
			"\nACTUAL:    ", classPrecision)
	}

	classRecall := getRecall("Mike", matrix)
	if metrics["Recall"] != classRecall {
		t.Error(
			"\nCLASS RECALL MISCALCULATED",
			"\nEXPECTED:  ", metrics["Recall"],
			"\nACTUAL:    ", classRecall)
	}

	fullAccuracy := getAccuracy(matrix)
	if metrics["Accuracy"] != fullAccuracy {
		t.Error(
			"\nALL TESTS INACCURATE",
			"\nEXPECTED:  ", metrics["Accuracy"],
			"\nACTUAL:    ", fullAccuracy)
	}

	fullCount := getTestCount(matrix)
	if metrics["FullCount"] != fullCount {
		t.Error(
			"\nALL TESTS MISCOUNT",
			"\nEXPECTED:  ", metrics["FullCount"],
			"\nACTUAL:    ", fullCount)
	}

	fullMatrix := fillMatrix(matrix)
	for key := range fullMatrix {
		if len(fullMatrix) != len(fullMatrix[key]) {
			t.Error(
				"\nMATRIX IS NOT EQUAL IN DIMENSIONS",
				"\nEXPECTED LENGTH:  ", len(fullMatrix),
				"\nACTUAL LENGTH:    ", len(fullMatrix[key]))
		}
	}

	classOutput := ClassSummary("John", fullMatrix)
	if classOutput == "" {
		t.Error(
			"\nNO OUTPUT STRING",
		)
	}

	fullOutput := FullSummary(fullMatrix)
	if fullOutput == "" {
		t.Error(
			"\nNO OUTPUT STRING",
		)
	}
}
