package conflation

import (
	"strconv"
	"strings"
)

// Scenario4 provides a filter to identify pull requests that have closed
// specific issues on GitHub.
type Scenario3 struct{}

var keywords = []string{"Close #", "Closes #", "Closed #", "Fix #", "Fixes #", "Fixed #", "Resolve #", "Resolves #", "Resolved #"}

func extractIssueID(expandedIssue *ExpandedIssue) int {
	fixIdx := 0
	for i := 0; i < len(keywords); i++ {
		fixIdx = strings.LastIndex(*expandedIssue.PullRequest.Body, keywords[i])
		if fixIdx != -1 {
			break
		}
	}
	if fixIdx == -1 {
		return -1
	}
	body := string(*expandedIssue.PullRequest.Body)
	body = body[fixIdx:]
	digit := digitRegexp.Find([]byte(body))
	issueID, _ := strconv.ParseInt(string(digit), 10, 32)
	return int(issueID)
}

func (s *Scenario3) ResolveIssueID(expandedIssue *ExpandedIssue) bool {
	issueID := extractIssueID(expandedIssue)
	if issueID != -1 {
		expandedIssue.PullRequest.RefIssueIds = []int{issueID}
		return true
	} else {
		return false
	}
}

func (s *Scenario3) Filter(expandedIssue *ExpandedIssue) bool {
	if expandedIssue.PullRequest.Body != nil {
		return s.ResolveIssueID(expandedIssue)
	} else {
		return false
	}
}
