package bhattacharya

import (
	"coralreefci/models/issues"
	"github.com/kljensen/snowball"
	"strings"
)

// DOC: StemIssues runs the snowball stemmer on object body text.
//      This helper function relies on the robustness of the string parsing
//      and of the third party library for performance.
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
