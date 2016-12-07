package conflation

import (
	"strings"
)

type Scenario4a struct {
	Words int
}

func (s *Scenario4a) Filter(expandedIssue *ExpandedIssue) bool {
	if strings.Count(*expandedIssue.Issue.Body, " ")+1 >= s.Words {
		return true
	} else {
		return false
	}
}
