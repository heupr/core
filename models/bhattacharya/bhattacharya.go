package bhattacharya

import (
	"core/pipeline/gateway/conflation"
	"core/utils"
	"go.uber.org/zap"
	"os"
	"sort"
	"strconv"
	"strings"
)

// DOC: NBClassifier is the struct implemented as the model algorithm.
type NBModel struct {
	classifier *NBClassifier
	assignees  []NBClass
}

// TODO: remove assets into separate file
type Result struct {
	id    int
	score float64
}

type Results []Result

func (slice Results) Len() int {
	return len(slice)
}

func (slice Results) Less(i, j int) bool {
	return slice[i].score < slice[j].score
}

func (slice Results) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (c *NBModel) IsBootstrapped() bool {
	return c.classifier != nil
}

func (c *NBModel) Learn(input []conflation.ExpandedIssue) {
	adjusted := c.converter(input...)

	removeStopWords(adjusted...)
	stemIssues(adjusted...)
	c.assignees = distinctAssignees(adjusted)
	if len(c.assignees) < 2 {
		//TODO: Add logging
		return
	}
	utils.ModelLog.Info("Bhattacharya Learn", zap.Int("AssigneesCount", len(c.assignees))) //, zap.String("Repository", repo))
	c.classifier = NewNBClassifierTfIdf(c.assignees...)
	for i := 0; i < len(input); i++ {
		c.classifier.Learn(strings.Split(adjusted[i].Body, " "), NBClass(adjusted[i].Assignees[0])) // NOTE: First position is a workaround
	}
	c.classifier.ConvertTermsFreqToTfIdf()
	//TODO: Fix later (logging related)
	/*
		for _, class := range c.assignees {
			utils.ModelDetails.Debug("Class: " + string(class))
			wordcount := c.classifier.WordsByClass(NBClass(class))
			utils.ModelDetails.Debug(wordcount)
		} */
}

func (c *NBModel) OnlineLearn(input []conflation.ExpandedIssue) {
	adjusted := c.converter(input...)
	removeStopWords(adjusted...)
	stemIssues(adjusted...)
	for i := 0; i < len(input); i++ {
		c.classifier.OnlineLearn(strings.Split(adjusted[i].Body, " "), NBClass(adjusted[i].Assignees[0])) // NOTE: First position is a workaround
	}
}

func (c *NBModel) Predict(input conflation.ExpandedIssue) []string {
	adjusted := c.converter(input)
	removeStopWordsSingle(&adjusted[0])
	stemIssuesSingle(&adjusted[0])
	scores, _, _ := c.classifier.LogScores(strings.Split(adjusted[0].Body, " "))

	results := Results{}
	for i := 0; i < len(scores); i++ {
		results = append(results, Result{id: i, score: scores[i]})
	}

	sort.Sort(sort.Reverse(results))

	names := []string{}
	for i := 0; i < len(results); i++ {
		names = append(names, string(c.assignees[results[i].id]))
	}

	//TODO: Improve logging
	utils.ModelLog.Info("\n")
	if input.Issue.Assignee != nil { //TODO: confirm why this is nil when processing a webhook
		utils.ModelLog.Info("", zap.String("Assignee", *input.Issue.Assignee.Login))
	}
	if input.Issue.URL != nil { //TODO: confirm why this is nil when processing a webhook
		utils.ModelLog.Info("", zap.String("URL", *input.Issue.URL))
	}
	for i := 0; i < len(names); i++ {
		utils.ModelLog.Info("", zap.String("Class", strconv.Itoa(i)+": "+names[i]+", Score: "+strconv.FormatFloat(results[i].score, 'f', -1, 64)))
	}
	return names
}

func (c *NBModel) GenerateRecoveryFile(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return c.classifier.WriteTo(file)
}

func (c *NBModel) RecoverModelFromFile(path string) error {
	NBClassifier, err := NewNBClassifierFromFile(path)
	if err != nil {
		return err
	}
	c.classifier = NBClassifier
	c.assignees = NBClassifier.Classes
	return nil
}

func distinctAssignees(issues []Issue) []NBClass {
	result := []NBClass{}
	j := 0
	for i := 0; i < len(issues); i++ {
		for j = 0; j < len(result); j++ {
			if issues[i].Assignees[0] == string(result[j]) { // NOTE: First position is a workaround
				break
			}
		}
		if j == len(result) {
			result = append(result, NBClass(issues[i].Assignees[0])) // NOTE: First position is a workaround
		}
	}
	return result
}

func convertClassToString(assignees []NBClass) []string {
	result := []string{}
	for i := 0; i < len(assignees); i++ {
		result = append(result, string(assignees[i]))
	}
	return result
}
