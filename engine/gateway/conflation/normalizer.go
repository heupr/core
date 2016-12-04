package conflation

type Normalizer struct {
	Context *Context
}

//One level deep (simple)
//TODO: handle deeper webs
func (n *Normalizer) Normalize() {
	expandedIssues := n.Context.Issues
	for i := 0; i < len(expandedIssues); i++ {
		refIssueIds := expandedIssues[i].PullRequest.RefIssueIds
		if refIssueIds == nil {
			continue
		} else {
			for j := 0; j < len(refIssueIds); j++ {
				for k := 0; k < len(expandedIssues); k++ {
					if expandedIssues[k].Issue.Number != nil && *expandedIssues[k].Issue.Number == refIssueIds[j] {
						expandedIssues[i].PullRequest.RefIssues = append(expandedIssues[i].PullRequest.RefIssues, expandedIssues[k].Issue)
						expandedIssues[k].Issue.RefPulls = append(expandedIssues[k].Issue.RefPulls, expandedIssues[i].PullRequest)
						if expandedIssues[i].Conflate {
							expandedIssues[k].Conflate = true
							expandedIssues[i].Conflate = false
						}
					}
				}
			}
		}
	}
}
