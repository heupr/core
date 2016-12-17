package models

import (
    "coralreefci/engine/gateway/conflation"
    "errors"
	"fmt"
)

// DOC: argument inputs into FoldImplementation:
//      trainIssues - this is the fold-defined length to train on (e.g. 10%)
//      testIssues - this is the fold-defined length to test on (e.g. 90%)
func (m *Model) FoldImplementation(trainIssues, testIssues []conflation.ExpandedIssue) (float64, matrix) {
	testCount, correct := len(testIssues), 0
	expected, predicted := []conflation.ExpandedIssue{}, []conflation.ExpandedIssue{}
	expected = append(expected, testIssues...)
	predicted = append(predicted, testIssues...)
	m.Learn(trainIssues)

	for i := 0; i < len(testIssues); i++ {
		assignees := m.Predict(testIssues[i])
		for j := 0; j < len(assignees); j++ {
			// NOTE: Implement a log that records:
            // - issue URL
            // - issue assignee(s)
			if assignees[j] == *testIssues[i].Issue.Assignee.Name {
				correct++
				*predicted[i].Issue.Assignee.Name = assignees[j]
				break
			} else {
				*predicted[i].Issue.Assignee.Name = assignees[0]
			}
		}
	}
	mat, err := m.BuildMatrix(expected, predicted)
	if err != nil {
		fmt.Println(err)
	}
	return float64(correct) / float64(testCount), mat
}

func (m *Model) JohnFold(issues []conflation.ExpandedIssue) (string, error) {
	issueCount := len(issues)
	if issueCount < 10 {
		return "", errors.New("LESS THAN 10 ISSUES SUBMITTED - JOHN FOLD")
	}

	finalScore := 0.00
	for i := 0.10; i < 0.90; i += 0.10 { // TODO: double check the logic / math here
		trainCount := int(Round(i * float64(issueCount)))
		score, _ := m.FoldImplementation(issues[:trainCount], issues[trainCount:])
		finalScore += score
	}
	return ToString(Round(finalScore / 9.00)), nil
}

func (m *Model) TwoFold(issueList []conflation.ExpandedIssue) (string, error) {
	issueCount := len(issueList)
	if issueCount < 10 {
		return "", errors.New("LESS THAN 10 ISSUES SUBMITTED - TWO FOLD")
	}
	trainEndPos := int(0.50 * float64(issueCount))
	trainIssues, testIssues := []conflation.ExpandedIssue{}, []conflation.ExpandedIssue{}
	trainIssues = append(trainIssues, issueList[0:trainEndPos]...)
	testIssues = append(testIssues, issueList[trainEndPos+1:]...)

	firstScore, _ := m.FoldImplementation(trainIssues, testIssues)
	secondScore, _ := m.FoldImplementation(testIssues, trainIssues)

	score := firstScore + secondScore

	return ToString(Round(score / 2.00)), nil
}

func (m *Model) TenFold(issueList []conflation.ExpandedIssue) (string, error) {
	issueCount := len(issueList)
	if issueCount < 10 {
		return "", errors.New("LESS THAN 10 ISSUES SUBMITTED - TEN FOLD")
	}

	finalScore := 0.00
	start := 0
	for i := 0.10; i <= 1.00; i += 0.10 {
		end := int(Round(i * float64(issueCount)))

		segment := issueList[start:end]
		remainder := []conflation.ExpandedIssue{}
		remainder = append(issueList[:start], issueList[end:]...)

		score, _ := m.FoldImplementation(segment, remainder) // TODO: specific logs for matrices

		finalScore += score
		start = end
	}
	return ToString(Round(finalScore / 10.00)), nil
}
