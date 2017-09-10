package ingestor

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"core/utils"
)

type IngestorServer struct {
	Server          http.Server
	Database        Database
	RepoInitializer RepoInitializer
}

// This is a global variable for unit testing and stubbing out the client URLs.
var NewClient = func(token oauth2.Token) *github.Client {
	source := oauth2.StaticTokenSource(&token)
	githubClient := github.NewClient(oauth2.NewClient(oauth2.NoContext, source))
	return githubClient
}

func (i *IngestorServer) activateHandler(w http.ResponseWriter, r *http.Request) {
	var activationParams struct {
		Repo  github.Repository `json:"repo"`
		Token *oauth2.Token     `json:"token"`
	}
	err := json.NewDecoder(r.Body).Decode(&activationParams)
	if err != nil {
		utils.AppLog.Error("unable to decode json message. ", zap.Error(err))
	}
	client := NewClient(*activationParams.Token)

	// NOTE: This may ultimately be refactored out into a helper
	// method/function. Also see the similar code in the Restart method.
	repo, _, err := client.Repositories.Get(context.Background(), *activationParams.Repo.Owner.Login, *activationParams.Repo.Name)
	if err != nil {
		utils.AppLog.Error("ingestor github pulldown:", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authRepo := AuthenticatedRepo{
		Repo:   repo,
		Client: client,
	}
	i.RepoInitializer = RepoInitializer{}
	i.RepoInitializer.AddRepo(authRepo)
}

func (i *IngestorServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/hook", collectorHandler())
	mux.HandleFunc("/activate-ingestor-backend", i.activateHandler)
	return mux
}

func (i *IngestorServer) Start() error {
	bufferPool := NewPool()
	i.Database = Database{BufferPool: bufferPool}
	i.Database.Open()

	i.RepoInitializer = RepoInitializer{}
	i.Restart()
	i.Continuity()
	i.Server = http.Server{Addr: "127.0.0.1:8030", Handler: i.routes()}
	err := i.Server.ListenAndServe()
	if err != nil {
		utils.AppLog.Error("ingestor server failed to start; ", zap.Error(err))
		return err
	}
	return nil
}

func (i *IngestorServer) Stop() {
	//TODO: Closing the server down is a needed operation that will be added.
}
