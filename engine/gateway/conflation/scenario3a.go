package conflation

import (
	"github.com/google/go-github/github"
)

// DOC: Scenario3a filters for "naked" pull requests.
//      These are pull requests without an associated issue.
type Scenario3a struct {
	// Algorithm Conflation
}

func (s *Scenario3a) Filter(pull github.PullRequest) bool {
	result := false
	if *pull.IssueURL == "" {
		result = true
	}
	return result
}
