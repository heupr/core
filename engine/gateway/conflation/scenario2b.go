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
func extractIssueId(crPull CrPullRequest) int {
	fixIdx := 0
	for i := 0; i < len(keywords); i++ {
		fixIdx = strings.LastIndex(*crPull.Body, keywords[i])
		if fixIdx != -1 {
			break
		}
	}
	if fixIdx == -1 {
		return -1
	}
	body := string(*crPull.Body)
	body = body[fixIdx:]
	digit := digitRegexp.Find([]byte(body))
	issueId, _ := strconv.ParseInt(string(digit), 10, 32)
	crPull.RefIssueIds = []int{int(issueId)}
	return int(issueId)
}

//TODO: Finish
func (s *Scenario2b) Filter(input ExpandedIssue) bool {
	/*
	  crPullRequest := input.(CRPullRequest)
	  issueId := extractIssueId(crPullRequest)
	  if issueId != -1 {
	    return true
	  } else {
	    return false
	  }*/

	return false
}
