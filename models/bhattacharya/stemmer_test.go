package bhattacharya

import (
	"coralreefci/models/issues"
	"testing"
)

var nonstemmedIssues = []issues.Issue{
	issues.Issue{Body: "Use the Force, Luke"},
	issues.Issue{Body: "At last we will have our revenge"},
	issues.Issue{Body: "No, I am your father"},
	issues.Issue{Body: "Who's scruffy looking?"},
	issues.Issue{Body: "I pledge myself to your teachings"},
}

var stemmedIssues = []issues.Issue{
	issues.Issue{Body: "use the force, luke"},
	issues.Issue{Body: "at last we will have our reveng"},
	issues.Issue{Body: "no, i am your father"},
	issues.Issue{Body: "who scruffi looking?"},
	issues.Issue{Body: "i pledg myself to your teach"},
}

func TestStemIssues(t *testing.T) {
	stemmedOutput := StemIssues(nonstemmedIssues...)
	for i := 0; i < len(stemmedOutput); i++ {
		if stemmedOutput[i].Body != stemmedIssues[i].Body {
			t.Error(
				"\nINPUT STRING NOT PARSED",
				"\nEXPECTED:", stemmedIssues[i].Body,
				"\nRECEIVED:", stemmedOutput[i].Body,
			)
		}
	}
}
