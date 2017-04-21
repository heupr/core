package onboarder

import (
	"fmt"
	"net/http"
	"time"

	"github.com/boltdb/bolt"

	"coralreefci/engine/onboarder/retriever"
)

type RepoServer struct {
	Server       http.Server
	Repos        map[int]*ArchRepo
	SQLDatabase  *retriever.MemSQL
	BoltDatabase BoltDB
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
}

func (rs *RepoServer) OpenSQL() {
	rs.SQLDatabase.Open()
}

func (rs *RepoServer) CloseSQL() {
	rs.SQLDatabase.Close()
}

func (rs *RepoServer) OpenBolt() error {
	boltDB, err := bolt.Open("storage.db", 0644, nil)
	if err != nil {
		return err
	}
	rs.BoltDatabase = BoltDB{db: boltDB}
	return nil
}

func (rs *RepoServer) CloseBolt() {
	rs.BoltDatabase.db.Close()
}

func (rs *RepoServer) Timer() {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for range ticker.C {
			// data, err := rs.SQLDatabase.Read()
			// if err != nil {
			//     // DO SOMETHING HERE
			// }
			// dispatcher := dispatcher.Dispatcher{
			//     Repos: rs.Repos,
			// }
            // dispatcher.Start(10)
			// Collector(data)
			// TODO: Implement the rest of the logic here.
		}
	}()
}

/*
func (rs *RepoServer) Conflate(issues map[int][]*github.Issue, pulls map[int][]*github.PullRequest, open map[int][]*github.Issue) {
    keys := []int{}
    for k := range issues {
        keys = append(keys, k)
    }
    for k := range pulls {
        if _, ok := issues[k]; !ok {
            keys = append(keys, k)
        }
    }

    for _, key := range keys {
        if _, ok := rs.Repos[key]; ok {
            // NOTE: The hard coded first element will need to change once
            // there are multiple models on the scene.
            if _, ok := issues[key]; ok {
                rs.Repos[key].Hive.Models[0].Conflator.SetIssueRequests(issues[key])
            }
            if _, ok := pulls[key]; ok {
                rs.Repos[key].Hive.Models[0].Conflator.SetPullRequests(pulls[key])
            }
        }
        rs.Repos[key].Hive.Models[0].Conflator.Conflate()
    }
}
*/
