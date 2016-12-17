package bhattacharya

import (
	"coralreefci/engine/gateway/conflation"
	"github.com/google/go-github/github"
	"testing"
)

var (
	number      = 1
	title       = "Detention Block AA-23"
	body        = "Uh, we're fine, we're all fine here now, thank you...how are you?"
	name        = "Han Solo"
	assignee    = github.User{Name: &name}
    githubIssue = github.Issue{Number: &number, Title: &title, Body: &body, Assignee: &assignee}
    crIssue     = conflation.CRIssue{githubIssue, []int{}, []conflation.CRPullRequest{}}
	testIssue   = []conflation.ExpandedIssue{conflation.ExpandedIssue{Issue: crIssue}}
)

func TestBhattacharyaConverter(t *testing.T) {
    nbc := &NBClassifier{}
	convertedIssues := nbc.bhattacharyaConverter(testIssue)
	t.Error(
        "ERROR",
        convertedIssues,
    )
}
