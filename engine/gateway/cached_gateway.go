package gateway

import (
	"github.com/google/go-github/github"
)

type CachedGateway struct {
	Gateway   *Gateway
	DiskCache *DiskCache
}

func (c *CachedGateway) GetPullRequests(org string, project string) (pulls []*github.PullRequest, err error) {
	key := "./" + org + project + "_pulls"
	err = c.DiskCache.TryGet(key, &pulls)
	if err != nil {
		pulls, err = c.Gateway.GetPullRequests(org, project)
		c.DiskCache.Set(key, pulls)
	}
	return pulls, err
}

func (c *CachedGateway) GetIssues(org string, project string) (issues []*github.Issue, err error) {
	key := "./" + org + project + "_issues"
	err = c.DiskCache.TryGet(key, &issues)
	if err != nil {
		issues, err = c.Gateway.GetIssues(org, project)
		c.DiskCache.Set(key, issues)
	}
	return issues, err
}
