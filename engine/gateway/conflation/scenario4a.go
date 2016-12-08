package conflation

import (
	"strings"
)

type Scenario4a struct {
	Words int
}

// DOC: Scenario4a is designed to filter for issues with a specified number of
//      words in the text body.
//      Note that this value is set in the Scenario4a Words struct value.
func (s *Scenario4a) Filter(expandedIssue *ExpandedIssue) bool {
	if strings.Count(*expandedIssue.Issue.Body, " ")+1 >= s.Words {
		return true
	} else {
		return false
	}
}
