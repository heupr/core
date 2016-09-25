package issues

type Issue struct {
	RepoId         int
	IssueId        int
	Assignee       string
	Body           string
	ImportantWords []string
}
