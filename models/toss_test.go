package models

import (
	"coralreefci/models/issues"
	"testing"
	"time"
)

var testIssues = []issues.Issue{
	{Assignee: "John", Resolved: time.Date(2016, time.October, 9, 0, 0, 0, 0, time.UTC), Labels: []string{"Blue", "White"}},
	{Assignee: "Mike", Resolved: time.Date(2016, time.October, 9, 0, 0, 0, 0, time.UTC), Labels: []string{"Blue", "Gold"}},
	{Assignee: "John", Resolved: time.Date(2016, time.October, 8, 0, 0, 0, 0, time.UTC), Labels: []string{"Blue", "Lion"}},
	{Assignee: "Woz", Resolved: time.Date(1990, time.October, 8, 0, 0, 0, 0, time.UTC), Labels: []string{"Blue", "Gold"}},
}

var logScores = []float64{1.00, 5.00, 2.00, 7.00, 6.00, 8.00, 3.00, 4.00, 9.00}
var topSelection = 3
var topIndex = []int{8, 5, 3}

func TestBuildProfiles(t *testing.T) {
	output := BuildProfiles(testIssues)
	if len(output) > 3 {
		t.Error(
			"\nTOO MANY RESULTS",
			"\nONLY 3 EXPECTED",
		)
	}

	for index, _ := range output {
		if len(output[index].Profile) > 3 {
			t.Error(
				"\nDUPLICATE PROFILE LABELS",
				"\nNAME:    ", output[index].Name,
				"\nPROFILE: ", output[index].Profile,
			)
		}
		if output[index].Contributions > 2 {
			t.Error(
				"\nCONTRIBUTIONS MISCOUNTED",
				"\nNAME:          ", output[index].Name,
				"\nCONTRIBUTIONS: ", output[index].Contributions,
			)
		}
	}
}

func TestTossing(t *testing.T) {
	testTossingGraph := TossingGraph{
		Assignees:  []string{"John", "Mike", "John", "Mike", "John", "Mike", "John", "Mike", "John", "Mike"},
		GraphDepth: topSelection,
	}

	output := testTossingGraph.Tossing(logScores)
	if len(output) != topSelection {
		t.Error(
			"\nINCORRECT NUMBER OF OUTPUTS",
			"\nEXPECTED: ", topSelection,
			"\nACTUAL:   ", len(output),
		)
	}
	for i := 0; i < len(output); i++ {
		if output[i] != topIndex[i] {
			t.Error(
				"\nSORTING ERROR",
				"\nEXPECTED: ", output[i], ",", "ACTUAL: ", topIndex[i],
			)
		}
	}
}
