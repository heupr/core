package conflation

type ComboAlgorithm struct {
	Context *Context
}

// DOC: linkPullRequestsToIssue allows for conflating individual pull requests
//	    to the respective issue objects.
func linkPullRequestsToIssue(issue *ExpandedIssue) {
	// DOC: This operation takes the values from the related go-github
	//	    PullRequest struct ("User") and places it into the go-github Issue
	//		struct field "Assignees" (which is a slice of "User" struct).
	// DOC: The singular "Assignee" field is necessary for linq sorting.
	issue.Issue.Assignee = issue.Issue.RefPulls[0].User
	for i := 0; i < len(issue.Issue.RefPulls); i++ {
		issue.Issue.Assignees = append(issue.Issue.Assignees, issue.Issue.RefPulls[i].User)
	}
	if issue.Issue.Body != nil && issue.Issue.RefPulls[0].Body != nil {
		*issue.Issue.Body = *issue.Issue.Body + " " + *issue.Issue.RefPulls[0].Body
	} else if issue.Issue.Body == nil && issue.Issue.RefPulls[0].Body != nil {
		issue.Issue.Body = issue.Issue.RefPulls[0].Body
	}
}

func linkTitleToBody(issue *ExpandedIssue) {
	if issue.PullRequest.Body != nil && issue.PullRequest.Title != nil {
		*issue.PullRequest.Body = *issue.PullRequest.Body + " " + *issue.PullRequest.Title
	} else if issue.Issue.Body != nil && issue.Issue.Title != nil {
		*issue.Issue.Body = *issue.Issue.Body + " " + *issue.Issue.Title
	}
}

// Accept a expanded "Issue" or "PR"
// PR's need to have reference information
func (c *ComboAlgorithm) Conflate(issue *ExpandedIssue) bool {
	if len(issue.Issue.RefPulls) > 0 {
		linkPullRequestsToIssue(issue)
	} else {
		linkTitleToBody(issue)
	}
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
