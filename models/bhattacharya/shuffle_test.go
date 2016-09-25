package bhattacharya

import (
	"coralreef-ci/models/issues"
	"testing"
)

func generateIssues() []issues.Issue {
	list := []issues.Issue{}
	for i := 0; i < 10; i++ {
		list = append(list, issues.Issue{Body: string(i) + "random text"})
	}
	return list
}

const seed = 0

func TestShuffle(t *testing.T) {
	originalList := generateIssues()
	shuffledList := generateIssues()
	Shuffle(shuffledList, seed)

	for index, _ := range originalList {
		if originalList[index].Body == shuffledList[index].Body {
			t.Error("LISTS HAVE NOT BEEN SHUFFLED")
			break
		}
	}
}
