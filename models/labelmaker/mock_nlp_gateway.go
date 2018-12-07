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
	sentences := []*languagepb.Sentence{
		*languagepb.Sentence{1, *languagepb.TextSpan{0,"Google, headquartered in Mountain View, unveiled the new Android phone at the Consumer Electronic Show."}},
		*languagepb.Sentence{1, *languagepb.TextSpan{105, "Sundar Pichai said in his keynote that users love their new Android phones."}}
	}
	tokens := []*languagepb.Token{
		*languagepb.Token{*languagepb.TextSpan{0, "Google"}, *languagepb.PartOfSpeech{
			*languagepb.PartOfSpeech_Tag.PartOfSpeech_NOUN,
			*languagepb.PartOfSpeech_Aspect.PartOfSpeech_ASPECT_UNKNOWN,
			*languagepb.PartOfSpeech_Case.PartOfSpeech_CASE_UNKNOWN,
			*languagepb.PartOfSpeech_Form.PartOfSpeech_FORM_UNKNOWN,
			*languagepb.PartOfSpeech_Gender.PartOfSpeech_GENDER_UNKNOWN,
			*languagepb.PartOfSpeech_Mood.PartOfSpeech_MOOD_UNKNOWN,
			*languagepb.PartOfSpeech_Number.PartOfSpeech_SINGULAR,
			*languagepb.PartOfSpeech_Person.PartOfSpeech_PERSON_UNKNOWN,
			*languagepb.PartOfSpeech_Proper.PartOfSpeech_PROPER,
			*languagepb.PartOfSpeech_Reciprocity.PartOfSpeech_RECIPROCITY_UNKNOWN,
			*languagepb.PartOfSpeech_Tense.PartOfSpeech_TENSE_UNKNOWN,
			*languagepb.PartOfSpeech_Voice.PartOfSpeech_VOICE_UNKNOWN,
		}, *languagepb.DependencyEdge{7, "NSUBJ"}, "Google"},
		*languagepb.Token{*languagepb.TextSpan{179, "."}, *languagepb.PartOfSpeech{
			*languagepb.PartOfSpeech_Tag.PartOfSpeech_PUNCT,
			*languagepb.PartOfSpeech_Aspect.PartOfSpeech_ASPECT_UNKNOWN,
			*languagepb.PartOfSpeech_Case.PartOfSpeech_CASE_UNKNOWN,
			*languagepb.PartOfSpeech_Form.PartOfSpeech_FORM_UNKNOWN,
			*languagepb.PartOfSpeech_Gender.PartOfSpeech_GENDER_UNKNOWN,
			*languagepb.PartOfSpeech_Mood.PartOfSpeech_MOOD_UNKNOWN,
			*languagepb.PartOfSpeech_Number.PartOfSpeech_UNKNOWN,
			*languagepb.PartOfSpeech_Person.PartOfSpeech_PERSON_UNKNOWN,
			*languagepb.PartOfSpeech_Proper.PartOfSpeech_PROPER_UNKNOWN,
			*languagepb.PartOfSpeech_Reciprocity.PartOfSpeech_RECIPROCITY_UNKNOWN,
			*languagepb.PartOfSpeech_Tense.PartOfSpeech_TENSE_UNKNOWN,
			*languagepb.PartOfSpeech_Voice.PartOfSpeech_VOICE_UNKNOWN,
		}, *languagepb.DependencyEdge{20, "P"}, "."}
	}
	return *languagepb.AnalyzeSyntaxResponse{
		sentences,
		tokens,
		"en",
	}
}
