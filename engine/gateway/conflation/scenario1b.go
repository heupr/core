package conflation

type Scenario1b struct {
}

// DOC: Scenario1b filters for issues that have comments attached to them.
func (s *Scenario1b) Filter(expandedIssue *ExpandedIssue) bool {
	return expandedIssue.Issue.Comments != nil
}
