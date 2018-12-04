package labelmaker

import languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"

type MockLanguageClient struct {
}

type MockNlpGateway struct {
	Client MockLanguageClient
}

func (client *MockLanguageClient) AnalyzeSentiment() (*languagepb.AnalyzeSentimentResponse, error) {
	documentSentiment := *languagepb.Sentiment{10.0, 0.5}
	return *languagepb.AnalyzeSentimentResponse{documentSentiment, "en"}, nil
}

func (client *MockLanguageClient) AnalyzeSyntax() {
}
