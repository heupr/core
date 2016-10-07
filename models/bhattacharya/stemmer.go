package bhattacharya

import (
	"coralreef-ci/models/issues"
	"github.com/kljensen/snowball"
	"strings"
)

func StemIssues(issueList []issues.Issue) {
	for i := 0; i < len(issueList); i++ {
		wordList := []string{}
		words := strings.Split(issueList[i].Body, " ")
		for _, word := range words {
			stem, _ := snowball.Stem(word, "english", true)
			wordList = append(wordList, stem)
		}
		wordString := strings.Join(wordList, " ")
		issueList[i].Body = wordString
	}
}
