package gateway

import (
	"fmt"
	"github.com/google/go-github/github"
	"testing"
)

func TestGateway(t *testing.T) {
	gateway := Gateway{Client: github.NewClient(nil)}

	pullRequests := gateway.GetPullRequests()
	issues := gateway.GetIssues()

	fmt.Println(pullRequests)
	fmt.Println(issues)
}
