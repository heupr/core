package jira

import (
  "core/pipeline/gateway"

	"github.com/google/go-github/github"
)

type CachedGateway struct {
	Gateway   *Gateway
	DiskCache *gateway.DiskCache
}

func (c *CachedGateway) GetIssues(repo, correctedFile string) (issues []*github.Issue, err error) {
	key := "/" + repo + "jira-issues"
	cacheError := c.DiskCache.TryGet(key, &issues)
	if cacheError != nil {
		issues, err = c.Gateway.GetIssues(repo, correctedFile)
		c.DiskCache.Set(key, issues)
	}
	return issues, err
}
