package conflation

import (
	"github.com/google/go-github/github"
	"strconv"
)

type OneToOne struct {
	Context *Context
}

func (c *OneToOne) extractIssueId(pull *github.PullRequest) int {
	//TODO: error handling
	digit, _ := strconv.ParseInt(string(digitRegexp.Find([]byte(*pull.IssueURL))), 10, 32)
	return int(digit)
}

//TODO: flesh this out
func (c *OneToOne) linkPullRequestsToIssues() {
	pulls := c.Context.Pulls

	for i := 0; i < len(pulls); i++ {
		pull := &pulls[i]
		issueId := c.extractIssueId(&pulls[i])
		issue := c.Context.Issues[issueId]
		//TODO: (Check for err/nil lookup)
		issue.Assignee = pull.Assignee
		c.Context.Issues[issueId] = issue
	}
}

func (c *OneToOne) Conflate() {
	c.linkPullRequestsToIssues()
}

// 1:1 Algorithm (Naive) (We may need to exclude 1:M issues)
// (Ideal? Approach 1)  We should be able to just use the closed indicator in corefx/pulls/12923
// Query https://api.github.com/repos/dotnet/corefx/pulls/12923 (pull request)
// Query https://api.github.com/repos/dotnet/corefx/issues/12886 (issue)
//
// (Alternative Approach 2) We can also use the event id
// (step 1)https://api.github.com/repos/dotnet/corefx/issues/12886/events
//                "id": 832840421,
//                "url": "https://api.github.com/repos/dotnet/corefx/issues/events/832840421",
//                "actor": {
//                "login": "stephentoub",
// (optional step 2) https://api.github.com/repos/dotnet/corefx/issues/events/832840421
// Next steps: Implement Approach 1
// Test Approach 1 with unit Test
