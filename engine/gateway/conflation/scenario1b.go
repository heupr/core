package conflation

import (

)

type Scenario1b struct {
    // TODO: Add specifics to this struct
}

func (s *Scenario1b) Filter(issue ExpandedIssue) bool {
    return true
}



/*
type CrPullRequest struct {
	github.PullRequest
	RefIssueIds []int
	RefIssues   []CrIssue
}

type CrIssue struct {
	github.Issue
	RefPullIds []int
	RefPulls   []CrPullRequest
}
type ExpandedIssue struct {
Issue       CrIssue
PullRequest CrPullRequest
}


func (cr *CrPullRequest) ReferencesIssues() bool {
	if len(cr.RefIssueIds) > 0 {
		return true
	} else {
		return false
	}
}
*/
