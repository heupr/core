// Sample language-quickstart uses the Google Cloud Natural API to analyze the
// sentiment of "Hello, world!".
package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"core/models/labelmaker"
	conf "core/pipeline/gateway/conflation"
	"core/pipeline/gateway"
	"core/tests/cmd/backtests/gateway/jira"
	"core/utils"

	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	// Imports the Google Cloud Natural Language API client package.
	language "cloud.google.com/go/language/apiv1"
	"golang.org/x/net/context"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
	jiraclient "github.com/andygrunwald/go-jira"
)

var client *language.Client
var ctx context.Context

func PrintLema(input string) {
	syntax, _ := client.AnalyzeSyntax(ctx, &languagepb.AnalyzeSyntaxRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: input,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})

	fmt.Println("\n\n", input)
	for i := 0; i < len(syntax.Tokens); i++ {
		//fmt.Println(syntax.Tokens[i].Text)
		//fmt.Println("Tag", syntax.Tokens[i].Text)
		fmt.Println(syntax.Tokens[i].DependencyEdge.Label, "(" + syntax.Tokens[syntax.Tokens[i].DependencyEdge.HeadTokenIndex].Text.Content + ", " +  syntax.Tokens[i].Text.Content + ")")
		/*
		if syntax.Tokens[i].PartOfSpeech.Aspect > 0 {
			fmt.Println("Aspect", syntax.Tokens[i].PartOfSpeech.Aspect)
		}
		if syntax.Tokens[i].PartOfSpeech.Case > 0 {
			fmt.Println("Case", syntax.Tokens[i].PartOfSpeech.Case)
		}
		if syntax.Tokens[i].PartOfSpeech.Mood > 0 {
			fmt.Println("Mood", syntax.Tokens[i].PartOfSpeech.Mood)
		}
		if syntax.Tokens[i].PartOfSpeech.Person > 0 {
			fmt.Println("Person", syntax.Tokens[i].PartOfSpeech.Person)
		}
		if syntax.Tokens[i].PartOfSpeech.Proper > 0 {
			fmt.Println("Proper", syntax.Tokens[i].PartOfSpeech.Proper)
		}
		if syntax.Tokens[i].PartOfSpeech.Reciprocity > 0 {
			fmt.Println("Reciprocity", syntax.Tokens[i].PartOfSpeech.Reciprocity)
		}
		if syntax.Tokens[i].PartOfSpeech.Tense > 0 {
			fmt.Println("Tense", syntax.Tokens[i].PartOfSpeech.Tense)
		}
		if syntax.Tokens[i].PartOfSpeech.Voice > 0 {
			fmt.Println("Voice", syntax.Tokens[i].PartOfSpeech.Voice)
		} */
	}

}

func GenerateLabel(input string, expectedOutput []string) bool {
	entities, err := client.AnalyzeEntities(ctx, &languagepb.AnalyzeEntitiesRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: input,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
	if err != nil {
		log.Fatalf("Failed to analyze text: %v", err)
	}

	results := entities.Entities

	for _, x := range results {
		fmt.Println("Entity ", x)
	}

	syntax, err := client.AnalyzeSyntax(ctx, &languagepb.AnalyzeSyntaxRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: input,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})

	for i := 0; i < len(syntax.Tokens); i++ {
		fmt.Println(syntax.Tokens[i])
		/*
		if syntax.Tokens[i].DependencyEdge.Label == languagepb.DependencyEdge_ROOT {
			fmt.Println(syntax.Tokens[i].Text.Content, syntax.Tokens[i].DependencyEdge.Label)
		}*/
	}

	fmt.Printf("Text: %v\n", input)
	return true
}

type LabelStats struct {
	Count int
	Total int
}

func jiraBacktest() {
	ctx = context.Background()
	var err error
	client, err = language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	NlpGateway := labelmaker.CachedNlpGateway{NlpGateway: &labelmaker.NlpGateway{Client: client}, DiskCache: &labelmaker.HashedDiskCache{}}

	jiraClient, _ := jiraclient.NewClient(nil, "https://issues.apache.org/jira/")
	jiraGateway := jira.CachedGateway{Gateway: &jira.Gateway{Client: jiraClient}, DiskCache: &gateway.DiskCache{}}
	jiraIssues, _ := jiraGateway.GetIssues("JCR", "./jackrabbit_classification_vs_type.csv")
	jiraIssues2, _ := jiraGateway.GetIssues("LUCENE", "./lucene_classification_vs_type.csv")
	jiraIssues3,_ := jiraGateway.GetIssues("HTTPCLIENT", "./httpclient_classification_vs_type.csv")
	jiraIssues = append(jiraIssues, jiraIssues2...)
	jiraIssues = append(jiraIssues, jiraIssues3...)


	featureLabel := "RFE"
	bugLabel := "BUG"
	improvementLabel := "IMPROVEMENT"
	lbModel := labelmaker.LBModel{Classifier: &labelmaker.LBClassifier{Ctx: ctx, Client: client, Gateway: NlpGateway}, FeatureLabel: &featureLabel, BugLabel: &bugLabel, ImprovementLabel: &improvementLabel}
	lbModel.Learn([]string{"BUG", "IMPROVEMENT"})

	fmt.Println("Building Label Distribution Data Structure...")
	labelDist := make(map[string]LabelStats)

	fmt.Println("Title Input(Only) Score...")
	correct := 0
	adjustedTotal := 0
	total := 0
	apiTotal := 0
	for i := 0; i < len(jiraIssues); i++ {
		if jiraIssues[i].Labels != nil && jiraIssues[i].PullRequestLinks == nil && jiraIssues[i].Body != nil {
			bugFeature := false
			for j := 0; j < len(jiraIssues[i].Labels); j++ {
					if *jiraIssues[i].Labels[j].Name == "IMPROVEMENT" {
						bugFeature = true
						break
					}
			}
			if !bugFeature {
				continue
			}
			if total % 100 == 0 {
				time.Sleep(5 * time.Second)
			}
			apiTotal++
			prediction, _ := lbModel.BugOrFeature(conf.ExpandedIssue{Issue: conf.CRIssue{*jiraIssues[i], []int{}, []conf.CRPullRequest{}, github.Bool(false)}})
			var predictions []string
			if prediction != nil {
				predictions = []string{*prediction}
			}
			//predictions, _ := lbModel.Predict(conf.ExpandedIssue{Issue: conf.CRIssue{github.Issue{Title: githubIssues[i].Title }, []int{}, []conf.CRPullRequest{}, github.Bool(false)}})
			fmt.Println("Input:", *jiraIssues[i].Title)
			/*
			syntax,_ := NlpGateway.AnalyzeSyntax(*jiraIssues[i].Title)
			for i := 0; i < len(syntax.Tokens); i++ {
				fmt.Println(syntax.Tokens[i].DependencyEdge.Label, "(" + syntax.Tokens[syntax.Tokens[i].DependencyEdge.HeadTokenIndex].Text.Content + ", " +  syntax.Tokens[i].Text.Content + ")")
			} */
			for j := 0; j < len(predictions); j++ {
				fmt.Println("Predicted:", predictions[j])
			}
			for j := 0; j < len(jiraIssues[i].Labels); j++ {
				label := *jiraIssues[i].Labels[j].Name
				if val, ok := labelDist[label]; ok {
					val.Total++
					labelDist[label] = val
				} else {
					labelDist[label] = LabelStats{}
				}
				fmt.Println("Expected:", *jiraIssues[i].Labels[j].Name)
			}
			for j := 0; j < len(predictions); j++ {
				for k := 0; k < len(jiraIssues[i].Labels); k++ {
					if predictions[j] == *jiraIssues[i].Labels[k].Name {
						stats := labelDist[predictions[j]]
						stats.Count++
						labelDist[predictions[j]] = stats
						correct++
						break
					}
				}
			}
			fmt.Println("\n")
			if len(predictions) > 0 {
				adjustedTotal++
			}
			total++
		}
	}
	for label, stats := range labelDist {
		fmt.Println(label, "Ratio", float64(stats.Count) / float64(stats.Total), "Count", stats.Count, "Total", stats.Total)
	}

	fmt.Println("Correct", correct)
	fmt.Println("Adjusted Total", adjustedTotal)
	fmt.Println("Total", total)
	fmt.Println("Adjusted Accuracy", float64(correct) / float64(adjustedTotal))
	fmt.Println("Total Accuracy", float64(correct) / float64(total))
	fmt.Println("API Calls", apiTotal)
}


func main() {
	jiraBacktest()
  return
	ctx = context.Background()

	var err error
	// Creates a client.
	client, err = language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	NlpGateway := labelmaker.CachedNlpGateway{NlpGateway: &labelmaker.NlpGateway{Client: client}, DiskCache: &labelmaker.HashedDiskCache{}}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "c813d7dab123d3c4813618bf64503a7a1efa540f"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	gClient := github.NewClient(tc)
	newGateway := gateway.CachedGateway{Gateway: &gateway.Gateway{Client: gClient}, DiskCache: &gateway.DiskCache{}}

	//repo := "yarnpkg/yarn"
	repo := "rust-lang/rust"
	//repo := "eslint/eslint"
	//repo := "NixOS/nix"
	r := strings.Split(repo, "/")
	githubIssues, err := newGateway.GetClosedIssues(r[0], r[1])
	if err != nil {
		utils.AppLog.Error("Cannot get Issues from Github Gateway.", zap.Error(err))
	}

	gitLabels, _ := newGateway.Gateway.GetLabels(r[0], r[1])
	trainingSet := make([]string, len(gitLabels))
	for i := 0; i < len(gitLabels); i++ {
		trainingSet[i] = *gitLabels[i].Name
	}
	featureLabel := "C-feature-request"
	bugLabel := "C-bug"
	improvementLabel := "C-enhancement"
	lbModel := labelmaker.LBModel{Classifier: &labelmaker.LBClassifier{Ctx: ctx, Client: client, Gateway: NlpGateway}, FeatureLabel: &featureLabel, BugLabel: &bugLabel, ImprovementLabel: &improvementLabel}
	lbModel.Learn(trainingSet)

	fmt.Println("Building Label Distribution Data Structure...")
	labelDist := make(map[string]LabelStats)


	fmt.Println("Title Input(Only) Score...")
	correct := 0
	adjustedTotal := 0
	total := 0
	apiTotal := 0
	for i := 0; i < len(githubIssues); i++ {
		if githubIssues[i].Labels != nil && githubIssues[i].PullRequestLinks == nil && githubIssues[i].Body != nil {
			bugFeature := false
			for j := 0; j < len(githubIssues[i].Labels); j++ {
					if *githubIssues[i].Labels[j].Name == "C-feature-request" {//*githubIssues[i].Labels[j].Name == "cat-bug" || *githubIssues[i].Labels[j].Name == "cat-feature" || *githubIssues[i].Labels[j].Name == "C-bug" || *githubIssues[i].Labels[j].Name == "C-feature-request" ||  *githubIssues[i].Labels[j].Name == "C-enhancement" {
						bugFeature = true
						break
					}
			}
			if !bugFeature {
				continue
			}
			if total % 100 == 0 {
				time.Sleep(5 * time.Second)
			}
			apiTotal++
			prediction, _ := lbModel.BugOrFeature(conf.ExpandedIssue{Issue: conf.CRIssue{*githubIssues[i], []int{}, []conf.CRPullRequest{}, github.Bool(false)}})
			var predictions []string
			if prediction != nil {
				predictions = []string{*prediction}
			}
			//predictions, _ := lbModel.Predict(conf.ExpandedIssue{Issue: conf.CRIssue{github.Issue{Title: githubIssues[i].Title }, []int{}, []conf.CRPullRequest{}, github.Bool(false)}})
			fmt.Println("Input:", *githubIssues[i].Title)
			for j := 0; j < len(predictions); j++ {
				fmt.Println("Predicted:", predictions[j])
			}
			for j := 0; j < len(githubIssues[i].Labels); j++ {
				label := *githubIssues[i].Labels[j].Name
				if val, ok := labelDist[label]; ok {
	    		val.Total++
					labelDist[label] = val
				} else {
					labelDist[label] = LabelStats{}
				}
				fmt.Println("Expected:", *githubIssues[i].Labels[j].Name)
			}
			for j := 0; j < len(predictions); j++ {
				for k := 0; k < len(githubIssues[i].Labels); k++ {
					if predictions[j] == *githubIssues[i].Labels[k].Name {
						stats := labelDist[predictions[j]]
						stats.Count++
						labelDist[predictions[j]] = stats
						correct++
						break
					}
				}
			}
			fmt.Println("\n")
			if len(predictions) > 0 {
				adjustedTotal++
			}
			total++
		}
	}
	for label, stats := range labelDist {
		if label == "cat-bug" || label == "cat-feature" || label == "C-bug" || label == "C-feature-request" {
			fmt.Println(label, "Ratio", float64(stats.Count) / float64(stats.Total), "Count", stats.Count, "Total", stats.Total)
		}
		//fmt.Println(label, "Ratio", float64(stats.Count) / float64(stats.Total), "Count", stats.Count, "Total", stats.Total)
	}

	fmt.Println("Correct", correct)
	fmt.Println("Adjusted Total", adjustedTotal)
	fmt.Println("Total", total)
	fmt.Println("Adjusted Accuracy", float64(correct) / float64(adjustedTotal))
	fmt.Println("Total Accuracy", float64(correct) / float64(total))
	fmt.Println("API Calls", apiTotal)


	/*
	fmt.Println("Body Input Score")


	lbModel = labelmaker.LBModel{Classifier: &labelmaker.LBClassifier{Ctx: ctx, Client: client}}
	lbModel.Learn([]string{"Os-Windows", "Os-Linux"})

	// Sets the text to analyze.
	//text := "I ran a fresh yarn install of a React Native app and the download took much longer than usual (a little over a minute and thirty seconds). Previous install took roughly 20 seconds. That being said, it did download much faster than npm using v4.6.1 (roughly 5 minutes)."

	text2 := "I tried installing yarn using the installation script from https://yarnpkg.com/en/docs/install#alternatives-tab (I'm on Windows, but I don't have admin rights, so I can't use the Windows installer)."
	labels, _ := lbModel.Predict(conf.ExpandedIssue{Issue: conf.CRIssue{github.Issue{Body: github.String(text2)}, []int{}, []conf.CRPullRequest{}, github.Bool(false)}})
	for i := 0; i < len(labels); i++ {
		fmt.Println("Prediction", labels[i])
	}
	text2 = "Installation Problem: Arch Linux – invalid PGP signature"
	labels, _ = lbModel.Predict(conf.ExpandedIssue{Issue: conf.CRIssue{github.Issue{Body: github.String(text2)}, []int{}, []conf.CRPullRequest{}, github.Bool(false)}})
	for i := 0; i < len(labels); i++ {
		fmt.Println("Prediction", labels[i])
	}
	text := "Os-Windows"
	GenerateLabel(text, []string{"cat-bug", "os-windows"})
	PrintLema("Compiler")
	PrintLema("Compile")
	PrintLema("Compilation")
	PrintLema("go")
	PrintLema("goes")
	PrintLema("parser")
	PrintLema("parse")
	GenerateLabel("System.alloc returns unaligned pointer if align > size", nil)*/
	//GenerateLabel("Miscompilation/regression in nightly-2017-11-21-x86_64-apple-darwin", nil)
	_ = client.Close()
}
