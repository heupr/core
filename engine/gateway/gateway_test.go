package gateway

import (
	"github.com/google/go-github/github"
	"reflect"
	"testing"
)

func TestGateway(t *testing.T) {
	client := github.NewClient(nil)
	gateway := Gateway{Client: client, UnitTesting: true}

	pullRequests, _ := gateway.GetPullRequests("dotnet", "corefx")
	issues, _ := gateway.GetIssues("dotnet", "corefx")

	if pullRequests == nil {
		t.Error(
			"\nEMPTY PULL REQUEST SLICE",
			"\nEXPECTED: ", reflect.TypeOf(new(github.PullRequest)),
			"\nRECEIVED: ", reflect.TypeOf(pullRequests[0]),
		)
	}

	if issues == nil {
		t.Error(
			"\nEMPTY ISSUES SLICE",
			"\nEXPECTED: ", reflect.TypeOf(new(github.Issue)),
			"\nRECEIVED: ", reflect.TypeOf(issues[0]),
		)
	}
}

func TestCachedGateway(t *testing.T) {
	client := github.NewClient(nil)
	gateway := CachedGateway{Gateway: &Gateway{Client: client, UnitTesting: true}, DiskCache: &DiskCache{}}

	pullRequests, _ := gateway.GetPullRequests("dotnet", "corefx")
	issues, _ := gateway.GetIssues("dotnet", "corefx")

	if pullRequests == nil {
		t.Error(
			"\nEMPTY PULL REQUEST SLICE",
			"\nEXPECTED: ", reflect.TypeOf(new(github.PullRequest)),
			"\nRECEIVED: ", reflect.TypeOf(pullRequests[0]),
		)
	}

	if issues == nil {
		t.Error(
			"\nEMPTY ISSUES SLICE",
			"\nEXPECTED: ", reflect.TypeOf(new(github.Issue)),
			"\nRECEIVED: ", reflect.TypeOf(issues[0]),
		)
	}
}
