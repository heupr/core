package issues

import (
    "time"
)

type Issue struct {
	RepoID         int
	IssueID        int
  Url            string
	Assignee       string
	Body           string
  Resolved       time.Time
  Labels         []string
}
