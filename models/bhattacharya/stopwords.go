package bhattacharya

import (
	"coralreefci/models/issues"
	"github.com/bbalet/stopwords"
	"strings"
)

// DOC: RemoveStopWords is a helper function that clears stopwords from the
//      target issue body text.
func RemoveStopWords(issueList ...issues.Issue) []issues.Issue {
	issueOutput := []issues.Issue{}
	for i := 0; i < len(issueList); i++ {
		cleaned := strings.TrimSpace(stopwords.CleanString(issueList[i].Body, "en", false))
		issueOutput = append(issueOutput, issues.Issue{Body: cleaned})
	}
	return issueOutput
}
