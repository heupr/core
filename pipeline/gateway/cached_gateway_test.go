package gateway

import (
	"testing"

	"github.com/google/go-github/github"
)

func TestCachedGateway(t *testing.T) {
	client := github.NewClient(nil)
	gateway := CachedGateway{
		Gateway: &Gateway{
			Client:      client,
			UnitTesting: true,
		},
		DiskCache: &DiskCache{},
	}

	_, err := gateway.GetClosedPulls("dotnet", "corefx")
	if err != nil {
		t.Errorf("failed cached gateway pull request fetch: %v", err)
	}
}
