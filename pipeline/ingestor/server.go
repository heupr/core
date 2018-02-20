package ingestor

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/go-fsnotify/fsnotify"
	"github.com/google/go-github/github"
	"go.uber.org/zap"

	"core/utils"
)

type IngestorServer struct {
	Server          http.Server
	Database        DataAccess
	RepoInitializer RepoInitializer
}

type storage struct {
	RepoID   int64
	FullName string   `schema:"FullName"`
	Labels   []string `schema:"Labels"`
	Buckets  map[string][]label
}

type label struct {
	Name     string
	Selected bool
}

// NewClient is a wrapper fo unit testing and stubbing out the client URLs.
var NewClient = func(appID int, installationID int) *github.Client {
	var key string
	if PROD {
		key = "heupr.2017-10-04.private-key.pem"
	} else {
		key = "mikeheuprtest.2017-11-16.private-key.pem"
	}
	// Wrap the shared transport for use with the Github Installation.
	itr, err := ghinstallation.NewKeyFromFile(
		http.DefaultTransport,
		appID,
		installationID,
		key,
	)
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
	i.Database = &Database{BufferPool: bufferPool}
	defer i.Database.Close()
	i.Database.open()

	i.RepoInitializer = RepoInitializer{
		Database: i.Database,
		HTTPClient: http.Client{
			Timeout: time.Second * 10,
		},
	}
	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					utils.AppLog.Info("filewatcher modified file:", zap.String("Event", event.Name))
					f, err := os.Open(event.Name)
					if err != nil {
						f.Close()
						utils.AppLog.Error("error opening user settings", zap.Error(err))
						continue
					}
					decoder := gob.NewDecoder(f)
					s := storage{}
					err = decoder.Decode(&s)
					f.Close()
					if err != nil {
						utils.AppLog.Error("error decoding user settings", zap.Error(err))
						continue
					}
					i.Database.InsertGobLabelSettings(s)
				}
			case err := <-watcher.Errors:
				utils.AppLog.Error("filewatcher error", zap.Error(err))
			}
		}
	}()

	err = watcher.Add(utils.Config.IngestorGobs)
	if err != nil {
		log.Fatal(err)
	}

	dispatcher := Dispatcher{Database: i.Database, RepoInitializer: &i.RepoInitializer}
	dispatcher.Start(5)

	//i.Restart()
	//i.Continuity()

	i.Server = http.Server{Addr: utils.Config.IngestorServerAddress, Handler: i.routes()}
	err = i.Server.ListenAndServe()
	if err != nil {
		utils.AppLog.Error("ingestor server failed to start", zap.Error(err))
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
