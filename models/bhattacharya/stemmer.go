package bhattacharya

import (
	"github.com/kljensen/snowball"
	"strings"
)

// DOC: This applies the snowball stemmer formula to the target words.
func stemIssues(issueList ...Issue) {
	// TODO: change the parameter types to "...*Issue"; this could
	//       possibly resolve the copy value errors and allow for the removal
	//       of the StemIssuesSingle workaround (also see StopWords for the
	//       same solution)
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

func stemIssuesSingle(issue *Issue) {
	wordList := []string{}
	words := strings.Split(issue.Body, " ")
	for _, word := range words {
		stem, _ := snowball.Stem(word, "english", true)
		wordList = append(wordList, stem)
	}
	wordString := strings.Join(wordList, " ")
	issue.Body = wordString
}
