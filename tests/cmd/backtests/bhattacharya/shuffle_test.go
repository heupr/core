package main

import (
	"core/pipeline/gateway/conflation"
	"github.com/google/go-github/github"
	"testing"
)

var letters = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

func generateRandomIssues() []conflation.ExpandedIssue {
	list := []conflation.ExpandedIssue{}
	for _, letter := range letters {
		githubIssue := github.Issue{Body: &letter}
		crIssue := conflation.CRIssue{githubIssue, []int{}, []conflation.CRPullRequest{}, false}
		list = append(list, conflation.ExpandedIssue{Issue: crIssue})
	}
	return list
}

const seed = 0

func TestShuffle(t *testing.T) {
	originalList := generateRandomIssues()
	shuffledList := generateRandomIssues()
	Shuffle(shuffledList, seed)

	for i, _ := range originalList {
		if originalList[i].Issue.Body == shuffledList[i].Issue.Body {
			t.Error(
				"LISTS HAVE NOT BEEN SHUFFLED",
				"\n", "ORIGINAL:", originalList,
				"\n", "SHUFFLED:", shuffledList,
			)
			break
		}
	}
}
