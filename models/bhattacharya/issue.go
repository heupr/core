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

func convertIssue(expandedIssue conflation.ExpandedIssue) Issue {
	issue := Issue{
		IssueNumber: *expandedIssue.Issue.Number,
		URL:         *expandedIssue.Issue.URL,
		Labels:      labelStrings(expandedIssue.Issue.Labels),
	}

	if expandedIssue.Issue.ClosedAt != nil {
		issue.Resolved = *expandedIssue.Issue.ClosedAt
	}
	if expandedIssue.Issue.Body == nil {
		issue.Body = "no body"
	} else {
		issue.Body = *expandedIssue.Issue.Body
	}

	if expandedIssue.Issue.Assignees != nil {
		for j := 0; j < len(expandedIssue.Issue.Assignees); j++ {
			issue.Assignees = append(issue.Assignees, *expandedIssue.Issue.Assignees[j].Login)
		}
	} else if expandedIssue.Issue.Assignee != nil {
		issue.Assignees = append(issue.Assignees, *expandedIssue.Issue.Assignee.Login)
	} else {
		issue.Assignees = append(issue.Assignees, "no issue assignee")
	}
	return issue
}

func convertPull(expandedIssue conflation.ExpandedIssue) Issue {
	issue := Issue{
		IssueNumber: *expandedIssue.PullRequest.Number,
		URL:         *expandedIssue.PullRequest.URL,
	}

	if expandedIssue.PullRequest.ClosedAt != nil {
		issue.Resolved = *expandedIssue.PullRequest.ClosedAt
	}
	if expandedIssue.PullRequest.Body == nil {
		issue.Body = "no body"
	} else {
		issue.Body = *expandedIssue.PullRequest.Body
	}
	if expandedIssue.PullRequest.User == nil {
		issue.Assignees = append(issue.Assignees, "no issue assignee")
	} else {
		issue.Assignees = append(issue.Assignees, *expandedIssue.PullRequest.User.Login)
	}
	return issue
}

func (n *NBModel) converter(expandedIssue ...conflation.ExpandedIssue) []Issue {
	output := []Issue{}
	for i := 0; i < len(expandedIssue); i++ {
		var issue Issue
		if expandedIssue[i].Issue.Number != nil {
			issue = convertIssue(expandedIssue[i])
		} else {
			issue = convertPull(expandedIssue[i])
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
