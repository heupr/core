package models

import "coralreefci/engine/gateway/conflation"

type Model struct {
    // NOTE: possibly truncate to an embedded field
	Algorithm Algorithm
}

type Algorithm interface {
	Learn(input []conflation.ExpandedIssue)
	Predict(input conflation.ExpandedIssue) []string
}

func (m *Model) Learn(input []conflation.ExpandedIssue) {
	m.Algorithm.Learn(input)
}

func (m *Model) Predict(input conflation.ExpandedIssue) []string {
	return m.Algorithm.Predict(input)
}
