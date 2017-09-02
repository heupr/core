package ingestor

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"core/pipeline/frontend"
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
	if r.FormValue("state") != frontend.BackendSecret {
		errMsg := "failed validating frontend-backend secret"
		utils.AppLog.Error(errMsg)
		http.Error(w, errMsg, http.StatusUnauthorized)
		return
	}
	repoInfo := r.FormValue("repos")
	repoInfo = strings.Trim(repoInfo, "{")
	repoInfo = strings.Trim(repoInfo, "}")
	repoSlice := strings.Split(repoInfo, ", ")
	owner := string(repoSlice[1])
	name := string(repoSlice[2])
	tokenString := r.FormValue("token")

	client := NewClient(oauth2.Token{AccessToken: tokenString})

	// NOTE: This may ultimately be refactored out into a helper
	// method/function. Also see the similar code in the Restart method.
	repo, _, err := client.Repositories.Get(context.Background(), owner, name)
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
