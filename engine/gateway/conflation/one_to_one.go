package conflation

import (
	"github.com/google/go-github/github"
	"strconv"
	"strings"
)

type OneToOne struct {
	Context *Context
}

// TODO: evaluate optimization
// There could be a better way to handle this logic. Once our unit testing is
// robust @taylor will play around (if needed for performance).
func (c *OneToOne) extractIssueID(pull *github.PullRequest) int {
	fixIdx := strings.LastIndex(*pull.Body, "Fixes")
	if fixIdx == -1 {
		return -1
	}
	body := string(*pull.Body)
	body = body[fixIdx:]
	digit := digitRegexp.Find([]byte(body))
	s, _ := strconv.ParseInt(string(digit), 10, 32)
	return int(s)
}

// TODO: consider a more robust solution
func (c *OneToOne) linkPullRequestsToIssues() {
	pulls := c.Context.Pulls
	for i := 0; i < len(pulls); i++ {
		if pulls[i].Body != nil {
			pull := &pulls[i]
			issueId := c.extractIssueID(&pulls[i])
			if issueId != -1 {
				issue := c.Context.Issues[issueId]
				if issue.Number != nil {
					//TODO: (Check for err/nil lookup)
					issue.Assignee = pull.Assignee
				}
				if issue.Body != nil {
					// First step towards using additional information
					*issue.Body = *issue.Body + *pulls[i].Body
				}
				c.Context.Issues[issueId] = issue
			}
		}
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
