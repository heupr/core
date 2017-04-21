package ingestor

import (
	"fmt"
	"net/http"
)

type IngestorServer struct {
	Server          http.Server
	RepoInitializer RepoInitializer
}

func (i *IngestorServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/hook", collectorHandler())
	return mux
}

func (i *IngestorServer) Start() {
	i.RepoInitializer = RepoInitializer{}
	i.RepoInitializer.LoadRepos()

	i.Server = http.Server{Addr: "127.0.0.1:8080", Handler: i.routes()}
	// TODO: Add in logging and remove print statement.
	err := i.Server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func (i *IngestorServer) Stop() {
	//TODO: Closing the server down is a needed operation that will be added.
}
