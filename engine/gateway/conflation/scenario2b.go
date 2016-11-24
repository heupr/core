package conflation

import (
	"strconv"
	"strings"
)

type Scenario2b struct {
}

var keywords = []string{"Close #", "Closes #", "Closed #", "Fix #", "Fixes #", "Fixed #", "Resolve #", "Resolves #", "Resolved #"}

// TODO: evaluate optimization
// There could be a better way to handle this logic. Once our unit testing is
// robust @taylor will play around (if needed for performance).
func extractIssueId(expandedIssue *ExpandedIssue) int {
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
	issueId, _ := strconv.ParseInt(string(digit), 10, 32)
	return int(issueId)
}

func (s *Scenario2b) ResolveIssueId(expandedIssue *ExpandedIssue) bool {
	issueId := extractIssueId(expandedIssue)
	if issueId != -1 {
		expandedIssue.PullRequest.RefIssueIds = []int{issueId}
		return true
	} else {
		return false
	}
}

func (s *Scenario2b) Filter(expandedIssue *ExpandedIssue) bool {
	if expandedIssue.PullRequest.Number != nil {
		return s.ResolveIssueId(expandedIssue)
	} else {
		return false
	}
}
