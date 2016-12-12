package conflation

import (
	"github.com/google/go-github/github"
)

type CRPullRequest struct {
	github.PullRequest
	RefIssueIds []int
	RefIssues   []CRIssue
}

type CRIssue struct {
	github.Issue
	RefPullIds []int
	RefPulls   []CRPullRequest
}

type ExpandedIssue struct {
	PullRequest CRPullRequest
	Issue       CRIssue
	Conflate    bool
}

func (cr *CRPullRequest) ReferencesIssues() bool {
	if len(cr.RefIssueIds) > 0 {
		return true
	} else {
		return false
	}
}

//TODO:
//Can we use these metrics to:
//A) create a more robust tossing graph
//B) question past issue assignment?
type Developer struct {
	//interests
	//stack overflow ranking or genuine areas of expertise
	//github repo count or github followers.
	//estimated hr/week
	//LOC/Mistake Ratio (LinesofCode, Mistake meaning subsequent fixes for the same issue)
	//PR/PR Rejected Ratio
}
