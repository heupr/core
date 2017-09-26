package frontend

import (
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"go.uber.org/zap"

	"core/utils"
)

type FrontendServer struct {
	Server     http.Server
	httpClient http.Client
	Database   BoltDB
}

func (fs *FrontendServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", mainHandler)
	mux.HandleFunc("/login", githubLoginHandler)
	mux.HandleFunc("/github_oauth_cb", fs.githubCallbackHandler)
	return mux
}

func (fs *FrontendServer) Start() {
	fs.OpenBolt()
	if err := fs.Database.Initialize(); err != nil {
		utils.AppLog.Error("frontend server: ", zap.Error(err))
		panic(err)
	}
	fs.CloseBolt()
	fs.httpClient = http.Client{Timeout: time.Second * 10}
	fs.Server = http.Server{Addr: "10.142.1.0:80", Handler: fs.routes()}
	if err := fs.Server.ListenAndServe(); err != nil {
		utils.AppLog.Error("frontend server failed to start", zap.Error(err))
		panic(err)
	}
}

func (fs *FrontendServer) Stop() {
	// TODO: Implement graceful shutdown.
}

func (fs *FrontendServer) OpenBolt() error {
	boltDB, err := bolt.Open(utils.Config.BoltDBPath, 0644, nil)
	if err != nil {
		utils.AppLog.Error("failed opening bolt", zap.Error(err))
		return err
	}
	fs.Database = BoltDB{DB: boltDB}
	return nil
}

func (fs *FrontendServer) CloseBolt() {
	fs.Database.DB.Close()
}
