package bhattacharya

import "testing"

var nonstemmedIssues = []Issue{
	Issue{Body: "Use the Force, Luke"},
	Issue{Body: "At last we will have our revenge"},
	Issue{Body: "No, I am your father"},
	Issue{Body: "Who's scruffy looking?"},
	Issue{Body: "I pledge myself to your teachings"},
}

var stemmedIssues = []Issue{
	Issue{Body: "use the force, luke"},
	Issue{Body: "at last we will have our reveng"},
	Issue{Body: "no, i am your father"},
	Issue{Body: "who scruffi looking?"},
	Issue{Body: "i pledg myself to your teach"},
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
