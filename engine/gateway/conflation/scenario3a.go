package conflation

type Scenario3a struct {
}

// DOC: Scenario3a filters for "naked" pull requests.
//      These are pull requests without an associated issue.
func (s *Scenario3a) Filter(expandedIssue *ExpandedIssue) bool {
	if expandedIssue.PullRequest.RefIssueIds == nil {
		return true
	} else {
		return false
	}
}
