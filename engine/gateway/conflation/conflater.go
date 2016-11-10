package conflation

import (
	"github.com/google/go-github/github"
	"regexp"
)

var digitRegexp = regexp.MustCompile("[0-9]+")

// NOTE: possibly needed for 1:M conflation
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

/*
func (c *Conflator) SetPullEvents(events []github.PullRequestEvent) {
    // stuff goes here
}

//TODO: Determine if we can just use SetPullEvents? Remove if that is the case
func (c *Conflator) SetIssueEvents(events []github.IssuesEvent) {
    // stuff goes here
}
*/

func (c *Conflator) SetPullRequests(pulls []github.PullRequest) {
	c.Context.Pulls = pulls
}

func (c *Conflator) SetIssueRequests(issues []github.Issue) {
	for i := 0; i < len(issues); i++ {
		issueNumber := *issues[i].Number
		c.Context.Issues[issueNumber] = issues[i]
	}
}

// NOTE: this could possibly be expanded
// Although the current focus is simply on conflating pull requests to raised
// issues, there will likely eventually be an expanded logic that accounts
// for a variety of other aspects to GitHub issues (e.g. participants,
// reference numbers, etc.) that will require their own logic.

func (c *Conflator) Conflate() {
	c.Algorithm.Conflate()
}

// NOTE: handle or ignore
// Pull requests that were intentionally not merged (as it pertains to the
// tossing graph: should that be a thumbs down for a developer in our model?)
//       "merged": false,
//       "mergeable": false,
// example: https://api.github.com/repos/dotnet/corefx/pulls/9460
// Next steps: ensure merged flag is populated
// Possibly create a "thin ice" list for our model's tossing graph
