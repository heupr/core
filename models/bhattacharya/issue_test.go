package bhattacharya

import (
	"coralreefci/engine/gateway/conflation"
	"github.com/google/go-github/github"
	"reflect"
	"testing"
	"time"
)

var (
	id          = 1137
	number      = 1
	url         = "http://www.death-star-party.com"
	title       = "Detention Block AA-23"
	login       = "Han Solo"
	assignee    = github.User{Login: &login}
	body        = "Uh, we're fine, we're all fine here now, thank you...how are you?"
	resolved    = time.Time{}
	name        = "garbage-shute"
	labels      = []github.Label{github.Label{Name: &name}}
	githubIssue = github.Issue{ID: &id, Number: &number, URL: &url, Title: &title, Assignee: &assignee, Body: &body, ClosedAt: &resolved, Labels: labels}
	crIssue     = conflation.CRIssue{githubIssue, []int{}, []conflation.CRPullRequest{}}
	testIssue   = conflation.ExpandedIssue{Issue: crIssue}
)

func TestConverter(t *testing.T) {
	nbc := NBClassifier{}
	convertedIssue := nbc.converter(testIssue)
	if reflect.DeepEqual(convertedIssue, testIssue) {
		t.Error(
			"\nNOT RETURNING ISSUE STRUCT",
			"\nOUTPUT:   ", convertedIssue,
		)
	}
}
