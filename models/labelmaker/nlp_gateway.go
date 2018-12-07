package labelmaker

import (
	"context"

	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

type NlpGateway struct {
	Client *language.Client
}

func (g *NlpGateway) AnalyzeSentiment(input string) (*languagepb.AnalyzeSentimentResponse, error) {
	return g.Client.AnalyzeSentiment(context.Background(), &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: input,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}

func (g *NlpGateway) AnalyzeSyntax(input string) (syntax *languagepb.AnalyzeSyntaxResponse, err error) {
	return g.Client.AnalyzeSyntax(context.Background(), &languagepb.AnalyzeSyntaxRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: input,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}
