package bhattacharya

import (
	"coralreef-ci/models/issues"
	"fmt"
	"testing"
	// "time"
)

var testIssues = []issues.Issue{
	{Assignee: "John", Resolved: time.Date(2016, time.October, 9, 0, 0, 0, 0, time.UTC), Labels: []string{"Blue", "White"}},
	{Assignee: "Mike", Resolved: time.Date(2016, time.October, 9, 0, 0, 0, 0, time.UTC), Labels: []string{"Blue", "Gold"}},
	{Assignee: "John", Resolved: time.Date(2016, time.October, 8, 0, 0, 0, 0, time.UTC), Labels: []string{"Pompous", "Ass"}},
}

func TestBuildProfiles(t *testing.T) {
	output := BuildProfiles(testIssues)
    fmt.Println(output)
}
