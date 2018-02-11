package labelmaker

import (
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

type CachedNlpGateway struct {
	NlpGateway *NlpGateway
	DiskCache  *HashedDiskCache
}

func (c *CachedNlpGateway) AnalyzeSentiment(input string) (sentiment *languagepb.AnalyzeSentimentResponse, err error) {
	key := input + "-AnalyzeSentiment"
	cacheError := c.DiskCache.TryGet(key, &sentiment)
	if cacheError != nil {
		sentiment, err = c.NlpGateway.AnalyzeSentiment(input)
		c.DiskCache.Set(key, sentiment)
	}
	return sentiment, err
}

func (c *CachedNlpGateway) AnalyzeSyntax(input string) (syntax *languagepb.AnalyzeSyntaxResponse, err error) {
	key := input + "-AnalyzeSyntax"
	cacheError := c.DiskCache.TryGet(key, &syntax)
	if cacheError != nil {
		syntax, err = c.NlpGateway.AnalyzeSyntax(input)
		c.DiskCache.Set(key, syntax)
	}
	return syntax, err
}
