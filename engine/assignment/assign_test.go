package assignment

import (
    "testing"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func TestAssignContributor(t *testing.T) {
	n := 12
	o := "forstmeier"
	r := "ihnil"
	a := "forstmeier"

	user := github.User{Login: &o}
	repo := github.Repository{Owner: &user, Name: &r}
	issue := github.Issue{Number: &n}
	issuesEvent := github.IssuesEvent{Issue: &issue, Repo: &repo}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "b143196838e4d077475dc7ebc337d33de02c9ad3"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	err := AssignContributor(a, issuesEvent, client)
	if err != nil {
		t.Error(err)
	}
}
