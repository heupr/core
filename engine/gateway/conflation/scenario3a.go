package conflation

import (
	"coralreefci/models/issues"
	"github.com/google/go-github/github"
)

// DOC: Scenario3a filters for "naked" pull requests.
//      These are pull requests without an associated issue.
type Scenario3a struct {
	// Algorithm Conflation
}

func (s *Scenario3a) Filter(pull github.PullRequest) issues.Issue {
	output := issues.Issue{}
	if *pull.IssueURL == "" {
		output = issues.Issue{
			// RepoID   int,           // NOTE: not sure if this is needed
			IssueID:  *pull.Number,
			Url:      *pull.URL,
			Assignee: *pull.Assignee.Name,
			Body:     *pull.Body,
			Resolved: *pull.MergedAt,  // NOTE: this is likely the correct time needed
			// Labels   []string,      // NOTE: not available on PRs
		}
	}
	return output
}
