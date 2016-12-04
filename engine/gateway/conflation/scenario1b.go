package conflation

type Scenario1b struct {
}

// DOC: Scenario1b filters for issues that have comments attached to them.
func (s *Scenario1b) Filter(expandedIssue ExpandedIssue) bool {
	result := false
	if expandedIssue.Issue.Number != nil {
		if expandedIssue.Issue.Comments != nil {
			result = true
		}
	}
	return result
}
