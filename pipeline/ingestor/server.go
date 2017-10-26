package ingestor

import (
	"context"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"

	"go.uber.org/zap"

	"core/utils"
)

type IngestorServer struct {
	Server          http.Server
	Database        Database
	RepoInitializer RepoInitializer
}

// NewClient is a wrapper fo unit testing and stubbing out the client URLs.
var NewClient = func(appId int, installationId int) *github.Client {
	var key string
	if PROD {
		key = "heupr.2017-10-04.private-key.pem"
	} else {
		key = "heupr.2017-10-04.private-key.pem" //TODO: Change this after deployment to GCP
	}
	// Wrap the shared transport for use with the Github Installation.
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, appId, installationId, key)
	if err != nil {
		utils.AppLog.Error("could not obtain github installation key", zap.Error(err))
		return nil
	}
	client := github.NewClient(&http.Client{Transport: itr})
	return client
}

func (i *IngestorServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/webhook", collectorHandler())
	return mux
}

func (i *IngestorServer) Start() error {
	bufferPool := NewPool()
	i.Database = Database{BufferPool: bufferPool}
	defer i.Database.Close()
	i.Database.Open()

	i.RepoInitializer = RepoInitializer{Database: &i.Database, HttpClient: http.Client{Timeout: time.Second * 10}}

	dispatcher := Dispatcher{Database: &i.Database, RepoInitializer: &i.RepoInitializer}
	dispatcher.Start(5)

	i.Restart()
	i.Continuity()

	i.Server = http.Server{Addr: utils.Config.IngestorServerAddress, Handler: i.routes()}
	err := i.Server.ListenAndServe()
	if err != nil {
		utils.AppLog.Error("ingestor server failed to start; ", zap.Error(err))
		return err
	}
	return nil
}

func (i *IngestorServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	i.Server.Shutdown(ctx)
	utils.AppLog.Info("graceful ingestor shutdown")
}
