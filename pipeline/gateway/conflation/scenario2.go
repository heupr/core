package conflation

type Scenario2 struct {
}

// DOC: Scenario2 filters for issues that have comments attached to them.
func (s *Scenario2) Filter(expandedIssue *ExpandedIssue) bool {
	return expandedIssue.Issue.Comments != nil
}
