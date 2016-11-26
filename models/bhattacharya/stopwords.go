package bhattacharya

import (
	"coralreefci/models/issues"
	"fmt"
	"github.com/bbalet/stopwords"
	"strings"
)

// DOC: RemoveStopWords is a helper function that clears stopwords from the
//      target issue body text.
func RemoveStopWords(issueList ...issues.Issue) {
	for i := 0; i < len(issueList); i++ {
		cleaned := strings.TrimSpace(stopwords.CleanString(issueList[i].Body, "en", false))
		issueList[0].Body = cleaned
		fmt.Println(issueList[0].Body) // TEMPORARY
	}
}
