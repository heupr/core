package models

import "coralreefci/engine/gateway/conflation"

type Model struct {
	Algorithm Algorithm
}

type Algorithm interface {
	IsBootstrapped() bool
	Learn(input []conflation.ExpandedIssue)
	OnlineLearn(input []conflation.ExpandedIssue)
	Predict(input conflation.ExpandedIssue) []string
	GenerateRecoveryFile(path string) error
	RecoverModelFromFile(path string) error
}

func (m *Model) IsBootstrapped() bool {
	return m.Algorithm.IsBootstrapped()
}

func (m *Model) Learn(input []conflation.ExpandedIssue) {
	m.Algorithm.Learn(input)
}

func (m *Model) OnlineLearn(input []conflation.ExpandedIssue) {
	m.Algorithm.OnlineLearn(input)
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
