package bhattacharya

import (
	"coralreefci/models/issues"
	"testing"
)

var nonstemmedIssues = []bhattacharyaIssue{
	bhattacharyaIssue{Body: "Use the Force, Luke"},
	bhattacharyaIssue{Body: "At last we will have our revenge"},
	bhattacharyaIssue{Body: "No, I am your father"},
	bhattacharyaIssue{Body: "Who's scruffy looking?"},
	bhattacharyaIssue{Body: "I pledge myself to your teachings"},
}

var stemmedIssues = []bhattacharyaIssue{
	bhattacharyaIssue{Body: "use the force, luke"},
	bhattacharyaIssue{Body: "at last we will have our reveng"},
	bhattacharyaIssue{Body: "no, i am your father"},
	bhattacharyaIssue{Body: "who scruffi looking?"},
	bhattacharyaIssue{Body: "i pledg myself to your teach"},
}

func TestStemIssues(t *testing.T) {
	StemIssues(nonstemmedIssues...)
	for i := 0; i < len(nonstemmedIssues); i++ {
		if nonstemmedIssues[i].Body != stemmedIssues[i].Body {
			t.Error(
				"\nINPUT STRING NOT PARSED",
				"\nEXPECTED:", stemmedIssues[i].Body,
				"\nRECEIVED:", nonstemmedIssues[i].Body,
			)
		}
	}
}
