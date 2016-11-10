package bhattacharya

import (
	"coralreefci/models/issues"
	"errors"
	"fmt" // currently necessary until logging is ready
	"math"
)

func (model *Model) Fold(issues []issues.Issue, graphDepth int) (float64, error) {
	issueCount := len(issues)
	if issueCount < 10 {
		return 0.00, errors.New("LESS THAN 10 ISSUES SUBMITTED")
	}

	score := 0.00

	for i := 0.10; i < 0.90; i += 0.10 { // TODO: double check the logic / math here
		correct := 0

		trainCount := int(Round(i * float64(issueCount)))
		testCount := issueCount - trainCount

		model.Learn(issues[0:trainCount])

		for j := trainCount + 1; j < issueCount; j++ {
			model.Logger.Log(issues[j].Url)
			model.Logger.Log(issues[j].Assignee)
			assignees := model.Predict(issues[j])
			for k := 0; k < graphDepth; k++ {
				if assignees[k] == issues[j].Assignee {
					correct += 1
				} else {
					continue
				}
			}
		}
		fmt.Println("Fold ", Round(i*10))
		fmt.Println("Total ", testCount)
		fmt.Println("Correct ", correct)
		fmt.Println("Accuracy ", float64(correct)/float64(testCount))
		score += float64(correct) / float64(testCount)
	}
	return score / 9.00, nil
}

func AppendCopy(slice []issues.Issue, elements ...issues.Issue) []issues.Issue {
	n := len(slice)
	total := len(slice) + len(elements)
	newSlice := make([]issues.Issue, total)
	if total > cap(slice) {
		// Reallocate. Grow to 1.5 times the new size, so we can still grow.
		newSize := total*3/2 + 1
		newSlice = make([]issues.Issue, total, newSize)
	}
	copy(newSlice, slice)
	copy(newSlice[n:], elements)
	return newSlice
}

func Append(slice []issues.Issue, elements ...issues.Issue) []issues.Issue {
	n := len(slice)
	total := len(slice) + len(elements)
	if total > cap(slice) {
		// Reallocate. Grow to 1.5 times the new size, so we can still grow.
		newSize := total*3/2 + 1
		newSlice := make([]issues.Issue, total, newSize)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[:total]
	copy(slice[n:], elements)
	return slice
}

func (model *Model) FoldImpl(train []issues.Issue, test []issues.Issue, tossingGraphLength int) (float64, matrix) {
	testCount, correct := len(test), 0
	predicted := AppendCopy(test[0:])
	model.Learn(train)
	for j := 0; j < len(test); j++ {
		assignees := model.Predict(test[j])
		if assignees[0] == test[j].Assignee {
			correct += 1
			predicted[j].Assignee = assignees[0]
		} else if assignees[1] == test[j].Assignee && tossingGraphLength > 1 {
			correct += 1
			predicted[j].Assignee = assignees[1]
		} else if assignees[2] == test[j].Assignee && tossingGraphLength == 3 {
			correct += 1
			predicted[j].Assignee = assignees[2]
		} else if assignees[3] == test[j].Assignee && tossingGraphLength == 4 {
			correct += 1
			predicted[j].Assignee = assignees[3]
		} else if assignees[4] == test[j].Assignee && tossingGraphLength == 5 {
			correct += 1
			predicted[j].Assignee = assignees[4]
		} else {
			predicted[j].Assignee = assignees[0]
			continue
		}
	}
	mat, err := BuildMatrix(AppendCopy(test[0:]), predicted)
	if err != nil {
		fmt.Println(err)
	}
	return float64(correct) / float64(testCount), mat
}

func (model *Model) TwoFold(issues []issues.Issue, tossingGraphLength int) (float64, []matrix, error) {
	length := len(issues)
	trainEndPos := int(0.50 * float64(length))
	trainIssues := AppendCopy(issues[0:trainEndPos])
	testIssues := AppendCopy(issues[trainEndPos+1 : length])

	score1, matrix1 := model.FoldImpl(trainIssues, testIssues, tossingGraphLength)
	score2, matrix2 := model.FoldImpl(testIssues, trainIssues, tossingGraphLength)
	score := score1 + score2

	return score / 2.00, []matrix{matrix1, matrix2}, nil
}

func (model *Model) TenFold(issues []issues.Issue) (float64, error) {
	pStart, testStartPos, testEndPos, testCount := 0, 0, 0, 0
	score := 0.00

	length := len(issues)
	for i := 0.10; i < 0.90; i += 0.10 {
		correct := 0
		testStartPos = int(i * float64(length))
		testEndPos = int((i + 0.10) * float64(length))
		trainIssues := AppendCopy(issues[pStart:testStartPos-1], issues[testEndPos+1:length]...)
		trainIssuesLength := len(trainIssues)
		testIssues := AppendCopy(issues[testStartPos:testEndPos])
		testCount = len(testIssues)

		model.Learn(trainIssues)
		for j := 0; j < len(testIssues); j++ {
			assignees := model.Predict(testIssues[j])
			if assignees[0] == issues[j].Assignee || assignees[1] == issues[j].Assignee || assignees[2] == issues[j].Assignee {
				correct += 1
			} else {
				continue
			}
		}
		fmt.Println("Fold ", Round(i*10))
		fmt.Println("Accuracy ", float64(correct)/float64(testCount))
		fmt.Println("Correct", float64(correct))
		fmt.Println("Train Count", float64(trainIssuesLength))
		fmt.Println("Test Count", float64(testCount))
		score += float64(correct) / float64(testCount)
	}
	return score / 9.00, nil
}

func Round(input float64) float64 {
	return math.Floor(input + 0.5)
}
