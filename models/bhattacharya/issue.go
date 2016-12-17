package bhattacharya

import (
	"coralreefci/engine/gateway/conflation"
	"github.com/google/go-github/github"
	"time"
)

// DOC: bhattacharyaIssue (within Bhattacharya) is a slimmed down version of
//      the more comprehensive ExpandedIssue type that is available to
//      the project.
type bhattacharyaIssue struct {
	RepoID   int
	IssueID  int
	URL      string
	Assignee string
	Body     string
	Resolved time.Time
	Labels   []string
}

func (n *NBClassifier) issueConverter(ei []conflation.ExpandedIssue) []bhattacharyaIssue {
	output := []bhattacharyaIssue{}
	for i := 0; i < len(ei); i++ {
		issue := bhattacharyaIssue{
			RepoID:   ei.Issue.ID,
			IssueID:  ei.Issue.Number,
			URL:      ei.Issue.URL,
			Assignee: ei.Issue.User.Login,
			Body:     ei.Issue.Body,
			Resolved: ei.Issue.ClosedAt,
			Labels:   labelStrings([]ei.Issue.Label),
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
		output = append(output, labels.Name)
	}
	return output
}
