package bhattacharya

import (
	"coralreefci/engine/gateway/conflation"
	"github.com/google/go-github/github"
	"time"
)

// DOC: Issue (within Bhattacharya) is a slimmed down version of ExpandedIssue.
type Issue struct {
	// RepoID      int     // TODO: Evaluate if this field is necessary.
	IssueNumber int
	URL         string
	Assignees   []string
	Body        string
	Resolved    time.Time
	Labels      []string
}

func (n *NBModel) converter(expandedIssue ...conflation.ExpandedIssue) []Issue {
	output := []Issue{}
	for i := 0; i < len(expandedIssue); i++ {
		issue := Issue{
			IssueNumber: *expandedIssue[i].Issue.Number,
			URL:         *expandedIssue[i].Issue.URL,
			Labels:      labelStrings(expandedIssue[i].Issue.Labels),
		}

		if expandedIssue[i].Issue.ClosedAt != nil {
			issue.Resolved = *expandedIssue[i].Issue.ClosedAt
		}
		if expandedIssue[i].Issue.Body == nil {
			issue.Body = "no body"
		} else {
			issue.Body = *expandedIssue[i].Issue.Body
		}
		if expandedIssue[i].Issue.Assignees == nil {
			issue.Assignees = append(issue.Assignees, "no issue assignee")
		} else {
			for j := 0; j < len(expandedIssue[i].Issue.Assignees); j++ {
				issue.Assignees = append(issue.Assignees, *expandedIssue[i].Issue.Assignees[j].Login)
			}
		}
		output = append(output, issue)
	}
	return output
}

func labelStrings(labels []github.Label) []string {
	output := []string{}
	if len(labels) != 0 {
		for i := 0; i < len(labels); i++ {
			output = append(output, *labels[i].Name)
		}
	}
	return output
}
