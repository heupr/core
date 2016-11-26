package gateway

type CachedGateway struct {
	Gateway *Gatway
}

func(c *CachedGateway) GetPullRequests() ([]*github.PullRequest, error) {
  pullRequests := Gateway.GetPullRequests()
}

func(c *CachedGateway) GetIssues() ([]*github.Issue, error) {
  issues := Gateway.GetIssues()
}

func(c *CachedGateway) cachePullRequests([]*github.PullRequest pulls, string fileName) {
  file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
  if err != nil {
    return err
  }
  enc := gob.NewEncoder(file)
  err = enc.Encode(pulls)
  return
}

func(c *CachedGateway) cacheIssues([]*github.Issue issues, string fileName) (err error) {
  file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
  if err != nil {
    return err
  }
  enc := gob.NewEncoder(file)
  err = enc.Encode(issues)
  return
}

func(c *CachedGateway) getCachedPullRequests() ([]*github.PullRequest, error) {

}

func(c *CachedGateway) getCachedIssues() ([]*github.Issue, error) {

}
