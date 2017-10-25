package gateway

import (
	"github.com/google/go-github/github"
	"go.uber.org/zap"

	"core/utils"
)

type CachedGateway struct {
	Gateway   *Gateway
	DiskCache *DiskCache
}

func (c *CachedGateway) getPulls(owner, repo, state string) (pulls []*github.PullRequest, err error) {
	key := "/" + owner + "-" + repo + state + "-pulls"
	cacheError := c.DiskCache.TryGet(key, &pulls)
	if cacheError != nil {
		utils.AppLog.Warn("CachedGateway: ", zap.Error(cacheError))
		utils.AppLog.Info("CachedGateway: Starting - Downloading Pulls from Github.",
			zap.String("repo", owner+repo),
		)
		if state == "open" {
			pulls, err = c.Gateway.GetOpenPulls(owner, repo)
		} else {
			pulls, err = c.Gateway.GetClosedPulls(owner, repo)
		}
		c.DiskCache.Set(key, pulls)
		utils.AppLog.Info("CachedGateway: Completed - Downloading Pulls from Github.",
			zap.String("repo", owner+repo),
		)
	}
	return pulls, err
}

func (c *CachedGateway) getIssues(owner, repo, state string) (issues []*github.Issue, err error) {
	key := "/" + owner + "-" + repo + state + "-issues"
	cacheError := c.DiskCache.TryGet(key, &issues)
	if cacheError != nil {
		utils.AppLog.Warn("CachedGateway: ", zap.Error(cacheError))
		utils.AppLog.Info("CachedGateway: Starting - Downloading Issues from Github.",
			zap.String("repo", owner+repo),
		)
		if state == "open" {
			issues, err = c.Gateway.GetOpenIssues(owner, repo)
		} else {
			issues, err = c.Gateway.GetClosedIssues(owner, repo)
		}
		c.DiskCache.Set(key, issues)
		utils.AppLog.Info("CachedGateway: Completed - Downloading Issues from Github.",
			zap.String("repo", owner+repo),
		)
	}
	return issues, err
}

func (c *CachedGateway) GetOpenPulls(owner, repo string) ([]*github.PullRequest, error) {
	return c.getPulls(owner, repo, "open")
}

func (c *CachedGateway) GetClosedPulls(owner, repo string) ([]*github.PullRequest, error) {
	return c.getPulls(owner, repo, "closed")
}

func (c *CachedGateway) GetOpenIssues(owner, repo string) ([]*github.Issue, error) {
	return c.getIssues(owner, repo, "open")
}

func (c *CachedGateway) GetClosedIssues(owner, repo string) ([]*github.Issue, error) {
	return c.getIssues(owner, repo, "closed")
}
