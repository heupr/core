package conflation

import (
	"github.com/google/go-github/github"
	"regexp"
)

var digitRegexp = regexp.MustCompile("[0-9]+")

type SubTask struct {
	Assignee string
	Body     string
}

type Conflator struct {
	Algorithm Conflation
	Context   *Context
}

type Context struct {
	Issues   map[int]github.Issue
	SubTasks map[int][]SubTask
	Pulls    []github.PullRequest
	Test     []string
}

func (c *Conflator) SetPullEvents(events []github.PullRequestEvent) {

}

//TODO: Determine if we can just use SetPullEvents? Remove if that is the case
func (c *Conflator) SetIssueEvents(events []github.IssuesEvent) {

}

func (c *Conflator) SetPullRequests(pulls []github.PullRequest) {
	c.Context.Pulls = pulls
}

func (c *Conflator) SetIssueRequests(issues []github.Issue) {
	for i := 0; i < len(issues); i++ {
		issueNumber := *issues[i].Number
		c.Context.Issues[issueNumber] = issues[i]
	}
}

func (c *Conflator) Conflate() {
	c.Algorithm.Conflate()
}

// Just a Thought
// Handle or Ignore? pull requests that were intentionally not merged (Tossing graph: should that be a thumbs down for a developer in our model? )
//       "merged": false,
//       "mergeable": false,
// example: https://api.github.com/repos/dotnet/corefx/pulls/9460
// Next steps: Ensure merged flag is populated
// Possibly create a "thin ice" list for our model's tossing graph
