package conflation

type Scenario4 struct {
}

// DOC: Scenario4 filters for "naked" pull requests.
//      These are pull requests without an associated issue.
func (s *Scenario4) Filter(expandedIssue *ExpandedIssue) bool {
	if expandedIssue.PullRequest.RefIssueIds == nil {
		return true
	} else {
		return false
	}
}
