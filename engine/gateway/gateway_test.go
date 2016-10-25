package gateway

import (
	"github.com/google/go-github/github"
	"reflect"
	"testing"
)

func TestGateway(t *testing.T) {
	gateway := Gateway{Client: github.NewClient(nil)}

	pullRequests := gateway.GetPullRequests()
	issues := gateway.GetIssues()

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
