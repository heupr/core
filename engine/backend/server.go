package backend

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"coralreefci/engine/frontend"
	"coralreefci/utils"
)

type ActiveRepos struct {
	sync.RWMutex
	Actives map[int]*ArchRepo
}

type BackendServer struct {
	Server   http.Server
	Database MemSQL
	Repos    *ActiveRepos
}

func (bs *BackendServer) activateHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != frontend.BackendSecret {
		utils.AppLog.Error("failed validating frontend-backend secret")
		return
	}
	repoInfo := r.FormValue("repos")
	repoID, err := strconv.Atoi(string(repoInfo[0]))
	if err != nil {
		utils.AppLog.Error("converting repo ID: ", zap.Error(err))
	}

	tokenString := r.FormValue("token")
	if bs.Repos.Actives[repoID] == nil {
		token := oauth2.Token{}
		if err := json.Unmarshal([]byte(tokenString), &token); err != nil {
			utils.AppLog.Error("converting tokens: ", zap.Error(err))
		}
		bs.NewArchRepo(repoID)
		bs.NewClient(repoID, &token)
		bs.NewModel(repoID)
	}
}

func (bs *BackendServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/activate-ingestor-backend", bs.activateHandler)
	bs.Server = http.Server{
		Addr:    "127.0.0.1:8020",
		Handler: mux,
	}
	bs.Server.ListenAndServe()

	bs.OpenSQL()
	defer bs.CloseSQL()

	bs.Repos = &ActiveRepos{Actives: make(map[int]*ArchRepo)}

	db, err := bolt.Open("storage.db", 0644, nil)
	defer db.Close()
	boltDB := frontend.BoltDB{DB: db}

	if err := boltDB.Initialize(); err != nil {
		utils.AppLog.Error("frontend server: ", zap.Error(err))
		panic(err)
	}

	keys, tokens, err := boltDB.RetrieveBulk("tokens")
	if err != nil {
		utils.AppLog.Error("frontend server: ", zap.Error(err))
		panic(err)
	}

	for i := 0; i < len(keys); i++ {
		key, err := strconv.Atoi(string(keys[i]))
		if err != nil {
			utils.AppLog.Error("frontend server: ", zap.Error(err))
			panic(err)
		}
		token := oauth2.Token{}
		if err := json.Unmarshal(tokens[i], &token); err != nil {
			utils.AppLog.Error("frontend server: ", zap.Error(err))
			panic(err)
		}
		if _, ok := bs.Repos.Actives[key]; !ok {
			bs.NewArchRepo(key)
			bs.NewClient(key, &token)
			bs.NewModel(key)
		}
	}

	// Keeping this channel to implement graceful shutdowns if needed.
	wiggin := make(chan bool)
	bs.Timer(wiggin)
}

func (bs *BackendServer) OpenSQL() {
	bs.Database.Open()
}

func (bs *BackendServer) CloseSQL() {
	bs.Database.Close()
}

// Periodically conducts pulldowns from the MemSQL database for processing.
func (bs *BackendServer) Timer(ender chan bool) {
	ticker := time.NewTicker(time.Second * 30)
	defer close(ender)

	bs.Dispatcher(10)

	for {
		select {
		case <-ticker.C:
			data, err := bs.Database.Read()
			if err != nil {
				utils.AppLog.Error("backend timer: ", zap.Error(err))
			}
			Collector(data)
		case <-ender:
			ticker.Stop()
			close(ender)
			return
		}
	}
}
