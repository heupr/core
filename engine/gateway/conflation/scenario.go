package conflation

type Scenario interface {
	Filter(input ExpandedIssue) bool
}
