package conflation

type Scenario7 struct{}

// DOC: Scenario7 accounts for "naked pull requests" without an associated
//      issue; the logic here populates the RefIssueIds field with its own
//      number for the ComboAlgorithm conflation logic.
//      Note that this particular filter is crucial for the proper operation of
//      the Bhattacharya model (in conjunction with Scenario3).
func (s *Scenario7) Filter(expandedIssue *ExpandedIssue) bool {
	scenario3 := Scenario3{}
	if expandedIssue.PullRequest.Number != nil {
		if !scenario3.Filter(expandedIssue) {
			expandedIssue.PullRequest.RefIssueIds = []int{*expandedIssue.PullRequest.Number}
			return true
		}
	}
	return false
}
