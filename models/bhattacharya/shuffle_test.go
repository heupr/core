package bhattacharya

import (
	"coralreef-ci/models/issues"
	"testing"
)

var letterRandomizers = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

func generateRandomIssues() []issues.Issue {
	list := []issues.Issue{}
	for _, value := range letterRandomizers {
		list = append(list, issues.Issue{Body: value + "random text"})
	}
	return list
}

const seed = 0

func TestShuffle(t *testing.T) {
	originalList := generateRandomIssues()
	shuffledList := generateRandomIssues()
	Shuffle(shuffledList, seed)

	for index, _ := range originalList {
		if originalList[index].Body == shuffledList[index].Body {
			t.Error("LISTS HAVE NOT BEEN SHUFFLED")
			break
		}
	}
}
