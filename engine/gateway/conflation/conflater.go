package conflation

import (
	"github.com/google/go-github/github"
	"regexp"
)

var digitRegexp = regexp.MustCompile("[0-9]+")

type Conflator struct {
	Scenarios            []Scenario
	ConflationAlgorithms []ConflationAlgorithm
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
	isValid := false
	issue := ExpandedIssue{}

	for i := 0; i < len(c.Context.Issues); i++ {
		issue = c.Context.Issues[i]
		for j := 0; j < len(c.Scenarios); j++ {
			isValid = c.Scenarios[j].Filter(issue)
			// DOC: combination conflation logic loop
			if isValid == true {
				for k := 0; k < len(c.ConflationAlgorithms); k++ {
					isConflated := c.ConflationAlgorithms[k].Conflate(issue)
					if isConflated == true {
						c.Context.Issues[i] = issue
						break
					}
				}
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
