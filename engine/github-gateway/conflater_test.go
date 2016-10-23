package github-gateway

import (
  "testing"
//  "github.com/google/go-github/github"
)



func TestConflater(t *testing.T) {
  conflator := Conflator{}

  conflator.SetPullEvents(nil)
  conflator.SetIssueEvents(nil)

  conflator.SetPullRequests(nil)
  conflator.SetIssueRequests(nil)
  conflator.Conflate()
}
