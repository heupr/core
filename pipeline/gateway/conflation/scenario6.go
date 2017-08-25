package conflation

type Scenario6 struct {
	AssigneeCount int
}

// DOC: Scenario6 filters only for issues (specifically) with a user-defined
//      count available as the AssigneeCount field.
//      Note that this particular filter is crucial to the clean operation of
//      the Bhattacharya model because the FoldImplementation will require
//      the Assignees field to operate.
func (s *Scenario6) Filter(expandedIssue *ExpandedIssue) bool {
	if len(expandedIssue.Issue.Assignees) >= s.AssigneeCount {
		return true
	}
	return false
}
