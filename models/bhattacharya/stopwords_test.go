package bhattacharya

import (
    "strings"
    "testing"
)

var withStopwords = "This is a sentence with stopwords in it, the best ever."
var withoutStopwords = "sentence stopwords"

func TestRemoveStopWords(t *testing.T) {
    output := RemoveStopWords(withStopwords)
    if strings.Contains(output, "this") || strings.Contains(output, "with") || strings.Contains(output, "the") {
        t.Error(
            "\nSTOPWORDS NOT REMOVED",
            "\nEXPECTED: ", withoutStopwords,
            "\nACTUAL:   ", output,
        )
    }
}
