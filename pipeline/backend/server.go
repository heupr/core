package backend

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"core/pipeline/frontend"
	"core/utils"
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
	var activationParams struct {
		Repo  github.Repository `json:"repo"`
		Token *oauth2.Token     `json:"token"`
		Limit int               `json:"limit"`
	}
	err := json.NewDecoder(r.Body).Decode(&activationParams)
	if err != nil {
		utils.AppLog.Error("unable to decode json message. ", zap.Error(err))
	}
	repoID := *activationParams.Repo.ID
	token := activationParams.Token
	limit := activationParams.Limit
	if bs.Repos.Actives[repoID] == nil {
		bs.NewArchRepo(repoID, limit)
		bs.NewClient(repoID, token)
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
	bs.OpenSQL()
	defer bs.CloseSQL()

	bs.Repos = &ActiveRepos{Actives: make(map[int]*ArchRepo)}

	db, err := bolt.Open(utils.Config.BoltDBPath, 0644, nil)
	boltDB := frontend.BoltDB{DB: db}

	if err := boltDB.Initialize(); err != nil {
		utils.AppLog.Error("backend server: ", zap.Error(err))
		panic(err)
	}

	keys, tokens, err := boltDB.RetrieveBulk("token")
	if err != nil {
		utils.AppLog.Error("backend server: ", zap.Error(err))
		panic(err)
	}

	for i := 0; i < len(keys); i++ {
		key, err := strconv.Atoi(string(keys[i]))
		if err != nil {
			utils.AppLog.Error("backend server: ", zap.Error(err))
			panic(err)
		}
		token := oauth2.Token{}
		if err := json.Unmarshal(tokens[i], &token); err != nil {
			utils.AppLog.Error("backend server: ", zap.Error(err))
			panic(err)
		}
		if _, ok := bs.Repos.Actives[key]; !ok {
			bs.NewArchRepo(key, 0)
			bs.NewClient(key, &token)
			bs.NewModel(key)
		}
	}
	db.Close()
	// Keeping this channel to implement graceful shutdowns if needed.
	wiggin := make(chan bool)
	bs.Timer(wiggin)
	bs.Server.ListenAndServe()
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

	bs.Dispatcher(10)

	go func() {
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
	}()
}
