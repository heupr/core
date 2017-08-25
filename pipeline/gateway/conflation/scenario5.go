package conflation

import (
	"strings"
)

type Scenario5 struct {
	Words int
}

// DOC: Scenario5 is designed to filter for issues with a specified number of
//      words in the text body.
//      Note that this value is set in the Scenario4a Words struct value.
func (s *Scenario5) Filter(expandedIssue *ExpandedIssue) bool {
	if strings.Count(*expandedIssue.Issue.Body, " ")+1 >= s.Words {
		return true
	} else {
		return false
	}
}
