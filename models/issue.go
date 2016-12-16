/*
The issues package creates an interal structure containing relavent GitHub object information.
This provides a universal object that CoralReefCI can operate on and pass around.
*/
package issues

import (
	"time"
)

// TODO: refactor out.
type Issue struct {
	RepoID   int
	IssueID  int
	URL      string
	Assignee string
	Body     string
	Resolved time.Time
	Labels   []string
}

// NOTE: due to the arguments for maintaining that only parental packages can
//       import from child inner packages, this would necessitate the Issue
//       struct being moved into models/ directly (rather than encapsulated
//       within issues/)
