package frontend

import (
	"net/http"

	"github.com/boltdb/bolt"
	"go.uber.org/zap"

	"coralreefci/utils"
)

type FrontendServer struct {
	Server   http.Server
	Database BoltDB
}

func (fs *FrontendServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", mainHandler)
	mux.HandleFunc("/login", githubLoginHandler)
	mux.HandleFunc("/github_oauth_cb", fs.githubCallbackHandler)
	mux.HandleFunc("/setup_complete", completeHandle)
	return mux
}

func (fs *FrontendServer) Start() {
	fs.Server = http.Server{Addr: "127.0.0.1:8080", Handler: fs.routes()}
	if err := fs.Server.ListenAndServe(); err != nil {
		utils.AppLog.Error("frontend server failed to start", zap.Error(err))
		panic(err)
	}
}

func (fs *FrontendServer) Stop() {
	// TODO: Implement graceful shutdown.
}

func (fs *FrontendServer) OpenBolt() error {
	boltDB, err := bolt.Open("storage.db", 0644, nil)
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
