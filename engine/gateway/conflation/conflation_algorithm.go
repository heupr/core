package conflation

type ConflationAlgorithm interface {
	Conflate(ExpandedIssue) bool
}
