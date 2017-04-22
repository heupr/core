package ingestor

import (
	"coralreefci/engine/gateway"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"testing"
)

func TestInsert(t *testing.T) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "23fc398670a80700b19b1ae1587825a16aa8ce57"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	newGateway := gateway.CachedGateway{Gateway: &gateway.Gateway{Client: client}, DiskCache: &gateway.DiskCache{}}

	githubIssues, _ := newGateway.GetIssues("dotnet", "corefx")
	githubPulls, _ := newGateway.GetPullRequests("dotnet", "corefx")

	db := Database{}
	db.Open()

	db.BulkInsertIssues(githubIssues, 555)
	db.BulkInsertPullRequests(githubPulls, 555)
	db.EnableRepo(555)
}
