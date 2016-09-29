package bhattacharya

import (
    "coralreef-ci/models/issues"
    "testing"
)

var issueTexts = []string{"Use the Force, Luke",
                          "At last we will have our revenge",
                          "No, I am your father",
                          "Who's scruffy looking?",
                          "I pledge myself to your teachings"}
var stemTexts = []string{"use the force, luke",
                         "at last we will have our reveng",
                         "no, i am your father",
                         "who scruffi looking?",
                         "i pledg myself to your teach"}

func generateStemIssues() []issues.Issue {
    list := []issues.Issue{}
    for _, value := range issueTexts {
        list = append(list, issues.Issue{Body: value})
    }
    return list
}

func TestStemIssues(t *testing.T) {
    issueList := generateStemIssues()
    StemIssues(issueList)
    for index, issue := range issueList {
        if issue.Body != stemTexts[index] {
            t.Error("\nINPUT STRING NOT PARSED",
                    "\nEXPECTED:", stemTexts[index],
                    "\nRECEIVED:", issue.Body)
        }
    }
}
