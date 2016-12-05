package bhattacharya

import (
	"coralreefci/models/issues"
	"github.com/bbalet/stopwords"
	"strings"
)

// DOC: RemoveStopWords is a helper function that clears stopwords from the
//      target issue body text.
func RemoveStopWords(issueList ...issues.Issue) {
	for i := 0; i < len(issueList); i++ {
		cleaned := strings.TrimSpace(stopwords.CleanString(issueList[i].Body, "en", false))
		issueList[i].Body = cleaned
	}
}

// Workaround
// TODO: come up with a fix
func RemoveStopWordsSingle(issue *issues.Issue) {
	issue.Body = strings.TrimSpace(stopwords.CleanString(issue.Body, "en", false))
}
