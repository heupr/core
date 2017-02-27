package assignment

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"testing"
)

func TestAssignContributor(t *testing.T) {
	n := 12
	o := "forstmeier"
	r := "ihnil"
	a := "forstmeier"

	user := github.User{Login: &o}
	repo := github.Repository{Owner: &user, Name: &r}
	issue := github.Issue{Number: &n, Repository: &repo}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "b143196838e4d077475dc7ebc337d33de02c9ad3"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	err := AssignContributor(a, issue, client)
	if err != nil {
		t.Error(err)
	}
}
