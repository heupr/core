package frontend

import (
	"github.com/google/go-github/github"

	"coralreefci/engine/gateway"
	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"
	"coralreefci/models/bhattacharya"
	// "coralreefci/utils"
)

const secretKey = "chalmun"

// TODO: Build so argument accepted is a slice (rather than a single value).
func (h *HeuprServer) NewHook(repo []*github.Repository, client *github.Client) error {
	for _, r := range repo {
		if check, err := h.hookExists(r, client); check {
			// handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			// TODO: Logic for handling an error here will be implemented; this
			//       will take the form of an exit from the parent NewHook method
			//       as well as a generation of an error/redirect page option to
			//       the end user of the Heupr application.
			// return handler, err
			// NOTE: Is this correct?
			return err
		}
		name := *r.Name
		owner := *r.Owner.Login
		url := "http://00ad0ac7.ngrok.io/hook"
		hook, _, err := client.Repositories.CreateHook(owner, name, &github.Hook{
			Name:   github.String("web"),
			Events: []string{"issues", "repository"},
			Active: github.Bool(true),
			Config: map[string]interface{}{
				"url":          url,
				"secret":       secretKey,
				"content_type": "json",
				"insecure_ssl": false,
			},
		})
		if err != nil {
			return err
		}
		err = h.Database.store(*r.ID, "hookID", *hook.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *HeuprServer) hookExists(repo *github.Repository, client *github.Client) (bool, error) {
	name, owner := "", ""
	if repo.Name != nil {
		name = *repo.Name
	}
	if repo.Owner.Login != nil {
		owner = *repo.Owner.Login
	}
	hookID, err := h.Database.retrieve(*repo.ID, "hookID")
	if err != nil {
		return false, err
	}

	_, _, err = client.Repositories.GetHook(owner, name, hookID.(int))
	if err != nil {
		return false, err
	}
	return true, nil
}

// NOTE: Possibly split out into separate file for handling model addition logic.
func (h *HeuprServer) AddModel(repo *github.Repository, client *github.Client) error {
	name := *repo.Name
	owner := *repo.Owner.Login
	repoID := *repo.ID
	//TODO: The comments field is not cached when using CachedGateway and will
	//      need to be fixed eventually.
	newGateway := gateway.Gateway{Client: client}
	githubIssues, err := newGateway.GetIssues(owner, name)
	if err != nil {
		// utils.Log.Error("Cannot get Issues from Gateway. ", err)
	}
	githubPulls, err := newGateway.GetPullRequests(owner, name)
	if err != nil {
		// utils.Log.Error("Cannot get PullRequests from Gateway. ", err)
	}

	context := &conflation.Context{}
	scenarios := []conflation.Scenario{&conflation.Scenario2{}}
	conflationAlgorithms := []conflation.ConflationAlgorithm{
		&conflation.ComboAlgorithm{Context: context},
	}
	normalizer := conflation.Normalizer{Context: context}
	conflator := conflation.Conflator{
		Scenarios:            scenarios,
		ConflationAlgorithms: conflationAlgorithms,
		Normalizer:           normalizer,
		Context:              context,
	}

	issuesCopy := make([]github.Issue, len(githubIssues))
	pullsCopy := make([]github.PullRequest, len(githubPulls))

	// TODO: Evaluate this particular snippet of code as it has potential
	//       performance optimization capabilities related to the hardware
	//       level. This may ultimately live in the actual gateway.go file to
	//	     improve the actual download operations.
	for i := 0; i < len(issuesCopy); i++ {
		issuesCopy[i] = *githubIssues[i]
	}
	for i := 0; i < len(pullsCopy); i++ {
		pullsCopy[i] = *githubPulls[i]
	}

	conflator.Context.Issues = []conflation.ExpandedIssue{}
	conflator.SetIssueRequests(issuesCopy)
	conflator.SetPullRequests(pullsCopy)
	conflator.Conflate()

	trainingSet := []conflation.ExpandedIssue{}

	for i := 0; i < len(conflator.Context.Issues); i++ {
		expandedIssue := conflator.Context.Issues[i]
		if expandedIssue.Conflate {
			if expandedIssue.Issue.Assignee == nil {
				continue
			} else {
				trainingSet = append(trainingSet, conflator.Context.Issues[i])
			}
		}
	}
	model := models.Model{Algorithm: &bhattacharya.NBModel{}}
	model.Algorithm.Learn(trainingSet)
	// TODO: Note that this logic is not complete; the "0" needs to be replaced
	//       with a dynamic means of allocating the model value.
	h.Repos[repoID].Hive.Models = append(h.Repos[repoID].Hive.Models, &HeuprModel{Model: &model})
	return nil
}
