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
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}
	repoIDString := r.FormValue("repos")
	repoID, err := strconv.Atoi(repoIDString)
	if err != nil {
		utils.AppLog.Error("converting repo ID: ", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if bs.Repos.Actives[repoID] == nil {
		db, err := bolt.Open("../frontend/storage.db", 0644, nil)
		if err != nil {
			utils.AppLog.Error("backend storage access: ", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer db.Close()

		boltDB := frontend.BoltDB{DB: db}

		byteToken, err := boltDB.Retrieve("token", repoID)
		if err != nil {
			utils.AppLog.Error("retrieving bulk data: ", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		token := oauth2.Token{}
		if err := json.Unmarshal(byteToken, &token); err != nil {
			utils.AppLog.Error("converting stored tokens: ", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		bs.NewArchRepo(repoID)
		bs.NewClient(repoID, &token)
		bs.NewModel(repoID)
	}
}

func (bs *BackendServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/activate-repos", bs.activateHandler)
	bs.Server = http.Server{
		Addr:    "127.0.0.1:8080",
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
		panic(err)
		utils.AppLog.Error("frontend server: ", zap.Error(err))
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
}

func (bs *BackendServer) OpenSQL() {
	bs.Database.Open()
}

func (bs *BackendServer) CloseSQL() {
	bs.Database.Close()
}

func (bs *BackendServer) Timer() {
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for range ticker.C {
			data, err := bs.Database.Read()
			if err != nil {
				panic(err)
				utils.AppLog.Error("backend timer method: ", zap.Error(err))
			}
			bs.Dispatcher(10)
			Collector(data)
		}
	}()
}
