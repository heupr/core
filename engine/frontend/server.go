package frontend

import (
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	// "github.com/google/go-github/github"
	// "golang.org/x/oauth2"

	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"
)

type HeuprServer struct {
	Server    http.Server
	Models    map[int]models.Model
	Conflator conflation.Conflator
	Database  BoltDB
}

func (h *HeuprServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandle)
	mux.HandleFunc("/login", githubLoginHandle)
	mux.HandleFunc("/github_oauth_cb", h.githubCallbackHandle)
	mux.HandleFunc("/setup_complete", completeHandle)
	// mux.Handle("/hook", collectorHandler())
	return mux
}

func (h *HeuprServer) openDB() error {
	boltDB, err := bolt.Open("storage.db", 0644, nil)
	if err != nil {
		return err
	}
	h.Database = BoltDB{db: boltDB}
	return nil
}

func (h *HeuprServer) closeDB() {
	h.Database.db.Close()
}

func (h *HeuprServer) Start() {
	h.Server = http.Server{Addr: "127.0.0.1:8080", Handler: h.routes()}
	// TODO: Add in logging and remove print statement.
	err := h.Server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func (h *HeuprServer) Stop() {
	//TODO: Closing the server down is a needed operation that will be added.
}
