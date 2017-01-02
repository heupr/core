package models

import (
	"coralreefci/engine/gateway/conflation"
    "fmt"
)

func (m *Model) FoldImplementation(trainIssues, testIssues []conflation.ExpandedIssue) (float64, matrix) {
	expected, predicted := []conflation.ExpandedIssue{}, []conflation.ExpandedIssue{}
	expected = append(expected, testIssues...)
	predicted = append(predicted, testIssues...)

	m.Learn(trainIssues)

	correct := 0
	for i := 0; i < len(testIssues); i++ {
		assignees := m.Predict(testIssues[i])
		for j := 0; j < len(assignees); j++ {
			if assignees[j] == *predicted[i].Issue.Assignee.Login {
				correct++
				*predicted[i].Issue.Assignee.Login = assignees[j]
				break
			} else {
				*predicted[i].Issue.Assignee.Login = assignees[0]
			}
		}
	}
	mat, err := m.BuildMatrix(expected, predicted)
	// TODO: repair this error generation; design idiomatically
	if err != nil {
		fmt.Println(err)
	}
	return float64(correct) / float64(len(testIssues)), mat
}

func (m *Model) JohnFold(issues []conflation.ExpandedIssue) string {
	finalScore := 0.00
	for i := 0.10; i < 0.90; i += 0.10 { // TODO: double check the logic / math here
		split := int(Round(i * float64(len(issues))))
		score, _ := m.FoldImplementation(issues[:split], issues[split:])
		finalScore += score
	}
	return ToString(Round(finalScore / 9.00))
}

func (m *Model) TwoFold(issues []conflation.ExpandedIssue) string {
	length := len(issues)
	split := int(0.50 * float64(length))

	firstHalf, secondHalf := []conflation.ExpandedIssue{}, []conflation.ExpandedIssue{}
	firstHalf = append(firstHalf, issues[:split]...)
	secondHalf = append(secondHalf, issues[split:]...)

	firstScore, _ := m.FoldImplementation(firstHalf, secondHalf)
	secondScore, _ := m.FoldImplementation(secondHalf, firstHalf)

	score := firstScore + secondScore

	return ToString(Round(score / 2.00))
}

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
