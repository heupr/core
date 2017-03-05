package onboarder

import (
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"

	"coralreefci/engine/gateway/conflation"
)

type RepoServer struct {
	Server    http.Server
	Repos     map[int]*ArchRepo
	Conflator conflation.Conflator
	Database  BoltDB
}

func (rs *RepoServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)
	mux.HandleFunc("/login", githubLoginHandler)
	mux.HandleFunc("/github_oauth_cb", rs.githubCallbackHandler)
	mux.HandleFunc("/setup_complete", completeHandle)
	return mux
}

func (rs *RepoServer) Start() {
	rs.Server = http.Server{Addr: "127.0.0.1:8080", Handler: rs.routes()}
	// TODO: Add in logging and remove print statement.
	err := rs.Server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func (rs *RepoServer) Stop() {
	// TODO: Closing the server down is a needed operation that will be added.
	// NOTE: Does the server need to be a pointer?
}

func (rs *RepoServer) OpenDB() error {
	boltDB, err := bolt.Open("storage.db", 0644, nil)
	if err != nil {
		return err
	}
	rs.Database = BoltDB{db: boltDB}
	return nil
}

func (rs *RepoServer) CloseDB() {
	rs.Database.db.Close()
}
