package conflation

import (
	"github.com/google/go-github/github"
	"regexp"
)

var digitRegexp = regexp.MustCompile("[0-9]+")

type Conflator struct {
	Scenarios            []Scenario
	ConflationAlgorithms []ConflationAlgorithm
	Normalizer           Normalizer
	Context              *Context
}

type Context struct {
	Issues []ExpandedIssue
}

func (c *Conflator) SetPullRequests(pulls []github.PullRequest) {
	for i := 0; i < len(pulls); i++ {
		c.Context.Issues = append(c.Context.Issues, ExpandedIssue{PullRequest: CrPullRequest{pulls[i], []int{}, []CrIssue{}}})
	}
}

func (c *Conflator) SetIssueRequests(issues []github.Issue) {
	for i := 0; i < len(issues); i++ {
		c.Context.Issues = append(c.Context.Issues, ExpandedIssue{Issue: CrIssue{issues[i], []int{}, []CrPullRequest{}}})
	}
}

// NOTE: this could possibly be expanded
// Although the current focus is simply on conflating pull requests to raised
// issues, there will likely eventually be an expanded logic that accounts
// for a variety of other aspects to GitHub issues (e.g. participants,
// reference numbers, etc.) that will require their own logic.
func (c *Conflator) Conflate() {
	c.filter()
	c.normalize()
	c.executeConflation()
}

func (c *Conflator) filter() {
	isValid := false
	issue := ExpandedIssue{}
	for i := 0; i < len(c.Context.Issues); i++ {
		issue = c.Context.Issues[i]
		for j := 0; j < len(c.Scenarios); j++ {
			isValid = c.Scenarios[j].Filter(&issue)
			issue.Conflate = isValid
			c.Context.Issues[i] = issue
		}
	}
}

func (c *Conflator) normalize() {
	c.Normalizer.Normalize()
}

func (c *Conflator) executeConflation() {
	issue := ExpandedIssue{}
	for i := 0; i < len(c.Context.Issues); i++ {
		issue = c.Context.Issues[i]
		if issue.Conflate {
			for j := 0; j < len(c.ConflationAlgorithms); j++ {
				c.ConflationAlgorithms[j].Conflate(&issue)
			}
		}
	}
}

// NOTE: handle or ignore
// Pull requests that were intentionally not merged (as it pertains to the
// tossing graph: should that be a thumbs down for a developer in our model?)
//       "merged": false,
//       "mergeable": false,
// example: https://api.github.com/repos/dotnet/corefx/pulls/9460
// Next steps: ensure merged flag is populated
// Possibly create a "thin ice" list for our model's tossing graph
