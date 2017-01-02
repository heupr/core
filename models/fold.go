package models

import (
	"coralreefci/engine/gateway/conflation"
	"fmt"
)

// DOC: JohnFold gradually increases the training data by increments of 1/10th.
func (m *Model) JohnFold(issues []conflation.ExpandedIssue) string {
	finalScore := 0.00
	// TODO: Double check the logic / math here in the loop.
	for i := 0.10; i < 0.90; i += 0.10 {
		split := int(Round(i * float64(len(issues))))
		score, _ := m.FoldImplementation(issues[:split], issues[split:])
		finalScore += score
	}
	return ToString(Round(finalScore / 9.00))
}

// DOC: TwoFold splits data in half - alternating training on each half.
func (m *Model) TwoFold(issues []conflation.ExpandedIssue) string {
	split := int(0.50 * float64(len(issues)))
	firstScore, _ := m.FoldImplementation(issues[:split], issues[split:])
	secondScore, _ := m.FoldImplementation(issues[split:], issues[:split])
	score := firstScore + secondScore
	return ToString(Round(score / 2.00))
}

// DOC: TenFold trains on a rolling 1/10th chunk of the input data.
func (m *Model) TenFold(issues []conflation.ExpandedIssue) string {
	length := len(issues)

	finalScore := 0.00
	start := 0
	for i := 0.10; i <= 1.00; i += 0.10 {
		end := int(Round(i * float64(length)))

		segment := issues[start:end]
		remainder := []conflation.ExpandedIssue{}
		remainder = append(issues[:start], issues[end:]...)

		score, _ := m.FoldImplementation(segment, remainder)

		finalScore += score
		start = end
	}
	return ToString(Round(finalScore / 10.00))
}

// DOC: FoldImplementation performs the learning / prediction operations on the
//      input data slices as determined by the "parent" fold method.
func (m *Model) FoldImplementation(train, test []conflation.ExpandedIssue) (float64, matrix) {
	expected := []string{}
	predicted := make([]string, len(test))

	// TODO: Possibly change to an indexing operation within the loop.
	for i := 0; i < len(test); i++ {
		expected = append(expected, *test[i].Issue.Assignee.Login)
	}

	m.Learn(train)
	correct := 0
	for i := 0; i < len(test); i++ {
		predictions := m.Predict(test[i])
		for j := 0; j < len(predictions); j++ {
			if predictions[j] == *test[i].Issue.Assignee.Login {
				predicted[i] = predictions[j]
				correct++
			} else {
				predicted[i] = predictions[0]
			}
		}
	}
	mat, err := m.BuildMatrix(expected, predicted)
	// TODO: Repair this error generation; design idiomatically.
	if err != nil {
		fmt.Println(err)
	}
	return float64(correct) / float64(len(test)), mat
}
