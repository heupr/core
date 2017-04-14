package onboarder

import (
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"

	"coralreefci/engine/gateway/conflation"
)

type RepoServer struct {
	Server    http.Server
	Repos     map[int]*HeuprRepo
	Conflator conflation.Conflator
	Database  BoltDB
}

func (h *RepoServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)
	mux.HandleFunc("/login", githubLoginHandler)
	mux.HandleFunc("/github_oauth_cb", h.githubCallbackHandler)
	mux.HandleFunc("/setup_complete", completeHandle)
	// mux.Handle("/hook", collectorHandler())
	return mux
}

func (h *RepoServer) Start() {
	h.Server = http.Server{Addr: "127.0.0.1:8080", Handler: h.routes()}
	// TODO: Add in logging and remove print statement.
	err := h.Server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func (h *RepoServer) Stop() {
	// TODO: Closing the server down is a needed operation that will be added.
	// NOTE: Does the server need to be a pointer?
}

func (h *RepoServer) OpenDB() error {
	boltDB, err := bolt.Open("storage.db", 0644, nil)
	if err != nil {
		return err
	}
	h.Database = BoltDB{db: boltDB}
	return nil
}

func (h *RepoServer) CloseDB() {
	h.Database.db.Close()
}
