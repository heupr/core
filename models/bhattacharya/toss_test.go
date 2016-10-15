package bhattacharya

import (
	"coralreefci/models/issues"
	"testing"
	"time"
)

var testIssues = []issues.Issue{
	{Assignee: "John", Resolved: time.Date(2016, time.October, 9, 0, 0, 0, 0, time.UTC), Labels: []string{"Blue", "White"}},
	{Assignee: "Mike", Resolved: time.Date(2016, time.October, 9, 0, 0, 0, 0, time.UTC), Labels: []string{"Blue", "Gold"}},
	{Assignee: "John", Resolved: time.Date(2016, time.October, 8, 0, 0, 0, 0, time.UTC), Labels: []string{"Blue", "Lion"}},
}

func TestBuildProfiles(t *testing.T) {
	output := BuildProfiles(testIssues)
	if len(output) > 2 {
		t.Error(
			"\nTOO MANY RESULTS",
			"\nONLY 2 EXPECTED",
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
	}
}
