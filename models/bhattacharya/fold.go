package bhattacharya

import (
	"coralreefci/models/issues"
	"errors"
	"fmt"
	"math"
)

func Round(input float64) float64 {
	return math.Floor(input + 0.5)
}

// NOTE: Arguments:
// trainIssues - this is the fold-defined length to train on (e.g. 10%)
// testIssues - this is the fold-defined length to test on (e.g. 90%)
func (m *Model) FoldImplementation(trainIssues, testIssues []issues.Issue) (float64, matrix) {
	testCount := len(testIssues) // TODO: put into one line (see #19)
	correct := 0                 // TODO: put into one line (see #18)
	expected := []issues.Issue{} // TODO: put into one line (see #22)
	expected = append(expected, testIssues...)
	predicted := []issues.Issue{} // TODO: put into one line (see #20)
	predicted = append(predicted, testIssues...)
	m.Learn(trainIssues)
	// TODO: there will need to be a logger living within this loop so that
	// issue assignees and URLs can be captured correctly
	// EXAMPLE:
	// 	m.Logger.Log(issues[j].Url)
	// 	m.Logger.Log(issues[j].Assignee)

	for i := 0; i < len(testIssues); i++ {
		assignees := m.Predict(testIssues[i])
		for j := 0; j < len(assignees); j++ {
			if assignees[j] == testIssues[i].Assignee {
				correct++
				predicted[i].Assignee = assignees[j]
				break
			} else {
				predicted[i].Assignee = assignees[0]
			}
		}
	}

	/*
		for i := 0; i < len(testIssues); i++ {
			assignees := m.Predict(testIssues[i])
			switch {
			case assignees[0] == testIssues[i].Assignee:
				correct++
				predicted[i].Assignee = assignees[0]
			case assignees[1] == testIssues[i].Assignee && graphDepth == 2:
				correct++
				predicted[i].Assignee = assignees[1]
			case assignees[2] == testIssues[i].Assignee && graphDepth == 3:
				correct++
				predicted[i].Assignee = assignees[2]
			case assignees[3] == testIssues[i].Assignee && graphDepth == 4:
				correct++
				predicted[i].Assignee = assignees[3]
			case assignees[4] == testIssues[i].Assignee && graphDepth == 5:
				correct++
				predicted[i].Assignee = assignees[4]
			default:
				predicted[i].Assignee = assignees[0]
				continue
			}
		}
	*/

	mat, err := BuildMatrix(expected, predicted)
	if err != nil {
		fmt.Println(err)
	}
	return float64(correct) / float64(testCount), mat
}

func (m *Model) JohnFold(issues []issues.Issue) (float64, error) {
	m.Logger.Log("--- John's fold ---") // TODO: possibly remove
	issueCount := len(issues)
	if issueCount < 10 {
		return 0.00, errors.New("LESS THAN 10 ISSUES SUBMITTED - JOHN FOLD")
	}

	finalScore := 0.00
	for i := 0.10; i < 0.90; i += 0.10 { // TODO: double check the logic / math here
		trainCount := int(Round(i * float64(issueCount)))

		// TODO: add in logging here for the output matrix on each loop run
		score, _ := m.FoldImplementation(issues[:trainCount], issues[trainCount:]) // TODO: specific logs for matrices

		finalScore += score
	}
	return Round(finalScore / 9.00), nil
}

func (m *Model) TwoFold(issueList []issues.Issue) (float64, error) {
	m.Logger.Log("--- Two fold ---")
	issueCount := len(issueList)
	if issueCount < 10 {
		return 0.00, errors.New("LESS THAN 10 ISSUES SUBMITTED - TWO FOLD")
	}
	trainEndPos := int(0.50 * float64(issueCount))
	trainIssues := []issues.Issue{} // TODO: put into one line (see #84)
	testIssues := []issues.Issue{}  // TODO: put into one line (see #83)
	trainIssues = append(trainIssues, issueList[0:trainEndPos]...)
	testIssues = append(testIssues, issueList[trainEndPos+1:]...)

	firstScore, _ := m.FoldImplementation(trainIssues, testIssues)
	secondScore, _ := m.FoldImplementation(testIssues, trainIssues)

	score := firstScore + secondScore

	return Round(score / 2.00), nil
}

func (m *Model) TenFold(issueList []issues.Issue) (float64, error) {
	m.Logger.Log("--- Ten fold ---")
	issueCount := len(issueList)
	if issueCount < 10 {
		return 0.00, errors.New("LESS THAN 10 ISSUES SUBMITTED - TEN FOLD")
	}

	finalScore := 0.00
	start := 0
	for i := 0.10; i <= 1.00; i += 0.10 {
		end := int(Round(i * float64(issueCount)))

		segment := issueList[start:end]
		remainder := []issues.Issue{}
		remainder = append(issueList[:start], issueList[end:]...)

		score, _ := m.FoldImplementation(segment, remainder) // TODO: specific logs for matrices

		finalScore += score
		start = end
	}
	return Round(finalScore / 10.00), nil
}
