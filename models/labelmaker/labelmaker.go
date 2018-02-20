package labelmaker

import (
	"strings"

	language "cloud.google.com/go/language/apiv1"
	"golang.org/x/net/context"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"

	"core/pipeline/gateway/conflation"
)

const (
	Unknown = iota
	Bug
	Feature
	Improvement
)

type LBModel struct {
	// LBClassifier is the struct implemented as the model algorithm.
	Classifier       *LBClassifier
	labels           []string
	FeatureLabel     *string
	BugLabel         *string
	ImprovementLabel *string
}

func (c *LBModel) IsBootstrapped() bool {
	return c.Classifier != nil
}

// Learn is a "fast path straight into the "learn logic".
func (c *LBModel) Learn(labels []string) {
	// Learning logic is accessed here.
	c.Classifier.Learn(labels)
}

func (c *LBModel) OnlineLearn(input []conflation.ExpandedIssue) {

}

func (c *LBModel) Predict(input conflation.ExpandedIssue) ([]string, error) {
	results, err := c.Classifier.Predict(*input.Issue.Title, false)
	return results, err
}

func (c *LBModel) ExperimentalPredict(input conflation.ExpandedIssue) ([]string, error) {
	results, err := c.Classifier.Predict(*input.Issue.Title, false)
	if results == nil && input.Issue.Body != nil {
		results, err = c.Classifier.Predict(*input.Issue.Body, true)
	}
	return results, err
}

func (c *LBModel) BugOrFeature(input conflation.ExpandedIssue) (*string, error) {
	result, err := c.Classifier.BugOrFeature(input)
	if err != nil {
		return nil, err
	}
	switch result {
	case Bug:
		return c.BugLabel, nil
	case Feature:
		return c.FeatureLabel, nil
	case Improvement:
		return c.ImprovementLabel, nil
	}
	return nil, nil
}

type LBClassifier struct {
	Client  *language.Client
	Gateway CachedNlpGateway
	Ctx     context.Context
	classes []LBLabel
}

type LBLabel struct {
	Text           string
	NormalizedText string
	LinkedWords    []string
}

func (c *LBClassifier) Learn(labels []string) {
	for i := 0; i < len(labels); i++ {
		label := c.normalizeLabel(labels[i])
		if label != nil {
			c.classes = append(c.classes, *label)
		}
	}
}

func (c *LBClassifier) BugOrFeature(input conflation.ExpandedIssue) (int, error) {
	label, err := c.BugOrFeatureTitle(*input.Issue.Title, *input.Issue.Body)
	if err != nil {
		return Unknown, err
	}
	/*
		if label == nil {
			label, err =  c.BugOrFeatureBody(*input.Issue.Body)
		}*/
	return label, err
}

func (c *LBClassifier) BugOrFeatureTitle(input string, body string) (int, error) {
	sentiment, err := c.Gateway.AnalyzeSentiment(input)
	if err != nil {
		return Unknown, err
	}

	if c.IsImprovement(input) {
		return Improvement, nil
	}

	if sentiment.DocumentSentiment.Magnitude > 0.5 && sentiment.DocumentSentiment.Score <= -0.8 {
		return Bug, nil
	}
	if sentiment.DocumentSentiment.Magnitude > 0.5 && sentiment.DocumentSentiment.Score >= 0.6 {
		return Feature, nil
	}
	//max(0.2, 0.4)
	if c.isVerb(input) && sentiment.DocumentSentiment.Magnitude > 0.2 && sentiment.DocumentSentiment.Score >= 0.4 {
		return Feature, nil
	}
	return Unknown, nil
}

func (c *LBClassifier) BugOrFeatureBody(input string) (int, error) {
	sentiment, err := c.Gateway.AnalyzeSentiment(input)
	if err != nil {
		return Unknown, err
	}
	length := 2
	if len(sentiment.Sentences) < length {
		length = len(sentiment.Sentences)
	}
	for i := 0; i < length; i++ {
		if sentiment.Sentences[i].Sentiment.Magnitude > 0.5 && sentiment.Sentences[i].Sentiment.Score <= -0.8 {
			return Bug, nil
		}
		if sentiment.Sentences[i].Sentiment.Magnitude > 0.5 && sentiment.Sentences[i].Sentiment.Score >= 0.6 {
			return Feature, nil
		}
	}
	return Unknown, nil
}

func (c *LBClassifier) RelaxedBugOrFeatureBody(input string) (int, error) {
	sentiment, err := c.Gateway.AnalyzeSentiment(input)
	if err != nil {
		return Unknown, err
	}
	length := 2
	if len(sentiment.Sentences) < length {
		length = len(sentiment.Sentences)
	}
	for i := 0; i < length; i++ {
		if sentiment.Sentences[i].Sentiment.Magnitude > 0.5 && sentiment.Sentences[i].Sentiment.Score <= -0.4 {
			return Bug, nil
		}
		if sentiment.Sentences[i].Sentiment.Magnitude > 0.5 && sentiment.Sentences[i].Sentiment.Score >= 0.5 {
			return Feature, nil
		}
	}
	return Unknown, nil
}

func (c *LBClassifier) Predict(input string, retry bool) ([]string, error) {
	var results []string
	entities, err := c.Client.AnalyzeEntities(c.Ctx, &languagepb.AnalyzeEntitiesRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: input,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
	if err != nil {
		return nil, err
	}

	entityString := ""
	for i := 0; i < len(entities.Entities); i++ {
		entityString += entities.Entities[i].Name + ". "
	}

	normalizedEntities := c.normalizeLabels(entityString)
	for i := 0; i < len(normalizedEntities); i++ {
		for j := 0; j < len(c.classes); j++ {
			//Single Letter Entities Create False Positives and are Skipped.
			if len(normalizedEntities[i].NormalizedText) == 1 {
				continue
			}
			if strings.Contains(c.classes[j].NormalizedText, normalizedEntities[i].NormalizedText) {
				//fmt.Println("[", c.classes[j].NormalizedText, "]", "==", "[", normalizedEntities[i].NormalizedText, "]") //TODO: Replace with logging
				duplicate := false
				for k := 0; k < len(results); k++ {
					if results[k] == c.classes[j].Text {
						duplicate = true
					}
				}
				if !duplicate {
					results = append(results, c.classes[j].Text)
					// if Title gives us nothing we take the best (top 4 Salience) we can get from Body
					if retry && len(results) > 4 {
						return results, nil
					}
				}
			} // end if
		} // end for
	} // end for
	return results, nil
}

func (c *LBClassifier) normalizeLabels(input string) []LBLabel {
	//TODO Fix Linked Words Logic. (Currently panics). Not necessary for MVP.
	syntax, _ := c.Client.AnalyzeSyntax(c.Ctx, &languagepb.AnalyzeSyntaxRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: input,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})

	if len(syntax.Tokens) == 0 {
		return nil
	}

	results := []LBLabel{}
	for i := 0; i < len(syntax.Sentences); i++ {
		label := LBLabel{Text: strings.TrimSuffix(syntax.Sentences[i].Text.Content, "."), NormalizedText: "", LinkedWords: []string{}}
		results = append(results, label)
	}

	j := -1
	for i := 0; i < len(syntax.Tokens); i++ {
		if syntax.Tokens[i].DependencyEdge.Label == languagepb.DependencyEdge_ROOT {
			j++
			rootWord := syntax.Tokens[i].Text.Content
			//fmt.Println(rootWord)
			results[j].NormalizedText = rootWord
		} // else if syntax.Tokens[i].PartOfSpeech.Tag == languagepb.PartOfSpeech_NOUN || syntax.Tokens[i].PartOfSpeech.Tag == languagepb.PartOfSpeech_ADJ {
		//	results[j].LinkedWords = append(results[j].LinkedWords, syntax.Tokens[i].Text.Content)
		//}
	}
	return results
}

func (c *LBClassifier) isVerb(input string) bool {
	syntax, _ := c.Gateway.AnalyzeSyntax(input)
	if len(syntax.Tokens) == 0 {
		return false
	}
	return syntax.Tokens[0].PartOfSpeech.Tag == languagepb.PartOfSpeech_VERB
}

func (c *LBClassifier) normalizeLabel(input string) *LBLabel {
	syntax, _ := c.Gateway.AnalyzeSyntax(input)

	if len(syntax.Tokens) == 0 {
		return nil
	}

	linkedWords := []string{}
	var rootWord string
	for i := 0; i < len(syntax.Tokens); i++ {
		if syntax.Tokens[i].DependencyEdge.Label == languagepb.DependencyEdge_ROOT {
			rootWord = syntax.Tokens[i].Text.Content
		} else if syntax.Tokens[i].PartOfSpeech.Tag == languagepb.PartOfSpeech_NOUN || syntax.Tokens[i].PartOfSpeech.Tag == languagepb.PartOfSpeech_ADJ {
			linkedWords = append(linkedWords, syntax.Tokens[i].Text.Content)
		}
	}
	return &LBLabel{Text: input, NormalizedText: rootWord, LinkedWords: linkedWords}
}

//IMO I think an "Improvement" is an artificial construct. Whereas Bug/Feature are innate.
//This particular artificial construct has very well defined boundaries. We are utilizing the boundaries defined in this research.
//Unfortunately we need this logic because otherwise Improvements are "misclassified" as Bug/Features.
//Perhaps hidden beneath these specifics lies some innate generalitiy.
//TODO: Generalize as much as possible. Automatically Verify boundaries for each repo that signs up.
func (c *LBClassifier) IsImprovement(input string) bool {
	if strings.HasPrefix(input, "Improve") {
		return true
	}

	if strings.Contains(input, "dependency") {
		if strings.HasPrefix(input, "Replace") || strings.HasPrefix(input, "Upgrade") || strings.HasPrefix(input, "Remove") {
			return true
		}
	}

	if strings.Contains(input, "memory usage") {
		if strings.Contains(input, "Reduce") || strings.Contains(input, "reduce") || strings.Contains(input, "High") || strings.Contains(input, "high") {
			return true
		}
	}

	if strings.Contains(input, "exception message") {
		return true
	}

	if strings.Contains(input, "Performance improvement") || strings.Contains(input, "performance improvement") {
		return true
	}

	if strings.Contains(input, "Poor performance") || strings.Contains(input, "poor performance") {
		return true
	}

	if strings.HasPrefix(input, "Optimize") || strings.HasPrefix(input, "optimize") || strings.HasPrefix(input, "Optimization") || strings.HasPrefix(input, "optimization") {
		return true
	}

	if strings.HasPrefix(input, "Log") {
		return true
	}

	/*
		if (strings.HasPrefix(input, "Limit")) {
			return true
		}*/

	return false
}
