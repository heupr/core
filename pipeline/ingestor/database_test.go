package ingestor

import (
	"core/pipeline/gateway"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"runtime"
	"testing"
)

func TestInsert(t *testing.T) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "Enter A Token"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	newGateway := gateway.CachedGateway{Gateway: &gateway.Gateway{Client: client}, DiskCache: &gateway.DiskCache{}}

	githubIssues, _ := newGateway.GetIssues("dotnet", "corefx")
	githubPulls, _ := newGateway.GetPullRequests("dotnet", "corefx")

	bufferPool := NewPool()
	db := Database{BufferPool: bufferPool}
	db.Open()

	repo := &github.Repository{ID: github.Int(26295345), Organization: &github.Organization{Name: github.String("dotnet")}, Name: github.String("coreclr")}
	for i := 0; i < len(githubIssues); i++ {
		githubIssues[i].Repository = repo
	}

	db.BulkInsertIssues(githubIssues)
	runtime.GC()

	issues, _ := db.ReadIssuesTest()
	fmt.Println(issues[0].Repository)

	db.BulkInsertPullRequests(githubPulls)
	runtime.GC()

	//db.InsertIssue(*githubIssues[0])
	//db.InsertPullRequest(*githubPulls[0])
	//db.EnableRepo(555)
}
