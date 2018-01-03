package labelmaker

import (
	"strings"

	"core/pipeline/gateway/conflation"

	language "cloud.google.com/go/language/apiv1"
	"golang.org/x/net/context"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

// DOC: LBClassifier is the struct implemented as the model algorithm.
type LBModel struct {
	Classifier *LBClassifier
	labels     []string
}

func (c *LBModel) IsBootstrapped() bool {
	return c.Classifier != nil
}

//Fast Path Straight into the "learn logic"
func (c *LBModel) Learn(labels []string) {
	//(Learning Logic)
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

type LBClassifier struct {
	Client  *language.Client
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
				continue;
			}
			if strings.Contains(c.classes[j].NormalizedText, normalizedEntities[i].NormalizedText) {
				//fmt.Println("[", c.classes[j].NormalizedText, "]", "==", "[", normalizedEntities[i].NormalizedText, "]") //TODO: Replace with logging
				duplicate := false
				for k := 0; k < len(results); k++ {
					if (results[k] == c.classes[j].Text) {
						duplicate = true
					}
				}
				if (!duplicate) {
					results = append(results, c.classes[j].Text)
					// if Title gives us nothing we take the best (top 4 Salience) we can get from Body
					if retry && len(results) > 4 {
						return results, nil
					}
				}
			} // end if
		} // end for
	}	// end for
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
		}// else if syntax.Tokens[i].PartOfSpeech.Tag == languagepb.PartOfSpeech_NOUN || syntax.Tokens[i].PartOfSpeech.Tag == languagepb.PartOfSpeech_ADJ {
		//	results[j].LinkedWords = append(results[j].LinkedWords, syntax.Tokens[i].Text.Content)
    //}
	}
  return results
}


func (c *LBClassifier) normalizeLabel(input string) *LBLabel {
	//TODO Remove Duplicate Code. Not necessary for MVP.
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
