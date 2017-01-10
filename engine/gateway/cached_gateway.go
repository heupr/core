package gateway

import (
	"coralreefci/utils"
	"github.com/google/go-github/github"
)

type CachedGateway struct {
	Gateway   *Gateway
	DiskCache *DiskCache
}

func (c *CachedGateway) GetPullRequests(org string, project string) (pulls []*github.PullRequest, err error) {
	key := "/" + org + project + "-pulls"
	cacheError := c.DiskCache.TryGet(key, &pulls)
	if cacheError != nil {
		utils.Log.Warning("CachedGateway: ", cacheError)
		utils.Log.Info("CachedGateway: Starting - Downloading Pulls from Github. Repo: ", org+project)
		pulls, err = c.Gateway.GetPullRequests(org, project)
		c.DiskCache.Set(key, pulls)
		utils.Log.Info("CachedGateway: Completed - Downloading Pulls from Github. Repo: ", org+project)
	}
	return pulls, err
}

func (c *CachedGateway) GetIssues(org string, project string) (issues []*github.Issue, err error) {
	key := "/" + org + project + "-issues"
	cacheError := c.DiskCache.TryGet(key, &issues)
	if cacheError != nil {
		utils.Log.Warning("CachedGateway: ", cacheError)
		utils.Log.Info("CachedGateway: Starting - Downloading Issues from Github. Repo: ", org+project)
		issues, err = c.Gateway.GetIssues(org, project)
		c.DiskCache.Set(key, issues)
		utils.Log.Info("CachedGateway: Completed - Downloading Issues from Github. Repo: ", org+project)
	}
	return issues, err
}
