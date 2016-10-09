package bhattacharya

import (
	"coralreef-ci/models/issues"
	"errors"
	"math"
)

func (model *Model) Fold(issues []issues.Issue) (float64, error) {
	issueCount := len(issues)
	if issueCount < 10 {
		return 0.00, errors.New("LESS THAN 10 ISSUES SUBMITTED")
	}

	score := 0.00

	for i := 0.10; i < 1.00; i += 0.10 {
		correct := 0

		trainCount := int(Round(i * float64(issueCount)))
		testCount := issueCount - trainCount

		model.Learn(issues[0:trainCount])

		for j := trainCount + 1; j < testCount; j++ {
			//TODO: loop through assignee's
			assignees := model.Predict(issues[j])
			if assignees[0] == issues[j].Assignee {
				correct += 1
			} else {
				continue
			}
		}
		score += float64(correct) / float64(issueCount)
	}
	return score / 10.00, nil
}

func Round(input float64) float64 {
	return math.Floor(input + 0.5)
}
