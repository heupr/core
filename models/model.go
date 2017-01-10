package models

import "coralreefci/engine/gateway/conflation"

type Model struct {
	Algorithm Algorithm
}

type Algorithm interface {
	Learn(input []conflation.ExpandedIssue)
	Predict(input conflation.ExpandedIssue) []string
	GenerateRecoveryFile(path string) error
	RecoverModelFromFile(path string) error
}

func (m *Model) Learn(input []conflation.ExpandedIssue) {
	m.Algorithm.Learn(input)
}

func (m *Model) Predict(input conflation.ExpandedIssue) []string {
	return m.Algorithm.Predict(input)
}

func (m *Model) GenerateRecoveryFile(path string) error {
	return m.Algorithm.GenerateRecoveryFile(path)
}

func (m *Model) RecoverModelFromFile(path string) error {
	return m.Algorithm.RecoverModelFromFile(path)
}
