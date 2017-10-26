package conflation

import "github.com/google/go-github/github"

type CRPullRequest struct {
	github.PullRequest
	RefIssueIds []int
	RefIssues   []CRIssue
}

type CRIssue struct {
	github.Issue
	RefPullIds []int
	RefPulls   []CRPullRequest
	Triaged    *bool
}

type ExpandedIssue struct {
	PullRequest CRPullRequest
	Issue       CRIssue
	Conflate    bool
	IsTrained   bool
}

func (cr *CRPullRequest) ReferencesIssues() bool {
	if len(cr.RefIssueIds) > 0 {
		return true
	} else {
		return false
	}
}
