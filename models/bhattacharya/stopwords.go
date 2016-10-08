package bhattacharya

import (
	"github.com/bbalet/stopwords"
	"strings"
)

func RemoveStopWords(body string) string {
	return strings.TrimSpace(stopwords.CleanString(body, "en", false))
}
