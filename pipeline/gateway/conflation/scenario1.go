package conflation

type Scenario1 struct{}

// Only Issue objects are collected without any additional criteria; this
// excludes Pull Requests.
func (s *Scenario1) Filter(expandedIssue *ExpandedIssue) bool {
	return expandedIssue.Issue.Number != nil
}
