package bhattacharya

import (
	"github.com/bbalet/stopwords"
	"strings"
)

// DOC: RemoveStopWords is a helper function that clears stopwords from the
//      target issue body text.
func RemoveStopWords(issueList ...Issue) {
	for i := 0; i < len(issueList); i++ {
		cleaned := strings.TrimSpace(stopwords.CleanString(issueList[i].Body, "en", false))
		issueList[i].Body = cleaned
	}
}

// TODO: This is a workaround that needs to be fixed
func RemoveStopWordsSingle(issue *Issue) {
	issue.Body = strings.TrimSpace(stopwords.CleanString(issue.Body, "en", false))
}
