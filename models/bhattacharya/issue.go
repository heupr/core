package bhattacharya

import (
	"coralreefci/engine/gateway/conflation"
	"github.com/google/go-github/github"
	"time"
)

// DOC: Issue (within Bhattacharya) is a slimmed down version of
//      the more comprehensive ExpandedIssue type that is available to
//      the project.
type Issue struct {
	RepoID   int
	IssueID  int
	URL      string
	Assignee string
	Body     string
	Resolved time.Time
	Labels   []string
}

func (n *NBClassifier) bhattacharyaConverter(ei []conflation.ExpandedIssue) []Issue {
	output := []Issue{}
	for i := 0; i < len(ei); i++ {
		issue := Issue{
			RepoID:   *ei[i].Issue.ID,
			IssueID:  *ei[i].Issue.Number,
			URL:      *ei[i].Issue.URL,
			Assignee: *ei[i].Issue.User.Login,
			Body:     *ei[i].Issue.Body,
			Resolved: *ei[i].Issue.ClosedAt,
			Labels:   labelStrings(ei[i].Issue.Labels),
		}
		output = append(output, issue)
	}
	return output
}

// DOC: labelStrings is a helper function to generate the string values for
//      for the GitHub Label struct.
func labelStrings(labels []github.Label) []string {
	output := []string{}
	for i := 0; i < len(labels); i++ {
		output = append(output, *labels[i].Name)
	}
	return output
}
