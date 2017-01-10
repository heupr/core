package bhattacharya

import (
	"coralreefci/engine/gateway/conflation"
	"github.com/google/go-github/github"
	"reflect"
	"testing"
	"time"
)

func TestConverter(t *testing.T) {
	number := 1
	url := "http://www.death-star-party.com"
	login := "Han Solo"
	user := github.User{Login: &login}
	assignees := []*github.User{&user}
	body := "Uh, we're fine, we're all fine here now, thank you...how are you?"
	resolved := time.Time{}
	name := "garbage-shute"
	labels := []github.Label{github.Label{Name: &name}}
	nbc := NBModel{}

	githubIssue := github.Issue{
		Number:    &number,
		URL:       &url,
		Assignees: assignees,
		Body:      &body,
		ClosedAt:  &resolved,
		Labels:    labels,
	}
	crIssue := conflation.CRIssue{githubIssue, []int{}, []conflation.CRPullRequest{}}
	testIssue := conflation.ExpandedIssue{Issue: crIssue}

	convertedIssue := nbc.converter(testIssue)
	if reflect.DeepEqual(convertedIssue, testIssue) {
		t.Error(
			"\nNOT RETURNING ISSUE STRUCT",
			"\nOUTPUT:   ", convertedIssue,
		)
	}
	if convertedIssue[0].IssueNumber != number {
		t.Error(
			"\nISSUE FIELDS NOT POPULATING",
			"\nISSUE.REPOID:", convertedIssue[0].IssueNumber,
		)
	}
	if convertedIssue[0].Body != body {
		t.Error("\nISSUE BODY FIELD NOT POPULATING")
	}
	if convertedIssue[0].Resolved != resolved {
		t.Error("\nISSUE RESOLVED FIELD NOT POPULATING")
	}
	if convertedIssue[0].Labels[0] != name {
		t.Error("\nISSUE LABELS FIELD NOT POPULATING")
	}

}
