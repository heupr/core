package conflation

type ComboAlgorithm struct {
	Context *Context
}

//TODO: link up using Context reference data
func linkPullRequestsToIssue(issue *ExpandedIssue) {
	issue.Issue.Assignee = issue.PullRequest.Assignee
	if issue.Issue.Body != nil {
		// First step towards using additional information
		*issue.Issue.Body = "Test"
	}
}

// Accept a expanded "Issue" or "PR"
// PR's need to have reference information
func (c *ComboAlgorithm) Conflate(issue ExpandedIssue) bool {
	linkPullRequestsToIssue(&issue)
	return true
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
