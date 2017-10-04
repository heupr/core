package conflation

type Scenario1 struct{}

// - "Basic" issues
// - Issue objects without any additional criteria
// - Only issues (pull requests excluded from filtering)
func (s *Scenario1) Filter(expandedIssue *ExpandedIssue) bool {
  return expandedIssue.Issue.Number != nil
}
