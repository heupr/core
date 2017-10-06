package frontend

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	// "github.com/boltdb/bolt"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"core/utils"
)

type State struct {
	sync.Mutex
	Tokens map[int]*oauth2.Token
}

type FrontendServer struct {
	Primary    http.Server
	Redirect   http.Server
	httpClient http.Client
	Database   BoltDB
	state      State
}

func (fs *FrontendServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", mainHandler)
	mux.HandleFunc("/login", githubLoginHandler)
	mux.HandleFunc("/github_oauth_cb", fs.githubCallbackHandler)
	return mux
}

func (fs *FrontendServer) LaunchServer(secure, unsecure, cert, key string) {
	// Primary server with HTTPS.
	fs.Primary = http.Server{
		Addr:    secure,
		Handler: fs.routes(),
	}
	go func() {
		if err := fs.Primary.ListenAndServeTLS(cert, key); err != nil {
			utils.AppLog.Error("primary server failed to start", zap.Error(err))
			panic(err)
		}
	}()

	// For redirection purposes only.
	fs.Redirect = http.Server{
		Addr:    unsecure,
		Handler: http.HandlerFunc(httpRedirect),
	}
	if err := fs.Redirect.ListenAndServe(); err != nil {
		utils.AppLog.Error("redirect server failed to start", zap.Error(err))
		panic(err)
	}
}

func (fs *FrontendServer) Start() {
	// fs.OpenBolt()
	// if err := fs.Database.Initialize(); err != nil {
	// 	utils.AppLog.Error("frontend server: ", zap.Error(err))
	// 	panic(err)
	// }
	// fs.CloseBolt()
	//
	// fs.state = State{Tokens: make(map[int]*oauth2.Token)}
	// fs.httpClient = http.Client{
	// 	Timeout: time.Second * 10,
	// }

	stopper := make(chan os.Signal)
	signal.Notify(stopper, os.Interrupt)

	if PROD {
		fs.LaunchServer("10.142.1.0:443", "10.142.1.0:80", "heupr_io.crt", "heupr.key")
	} else {
		fs.LaunchServer("127.0.0.1:8081", "127.0.0.1:8080", "cert.pem", "key.pem")
	}

	<-stopper
	utils.AppLog.Info("keyboard interrupt received")
	fs.Stop()
}

func (fs *FrontendServer) Stop() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	fs.Primary.Shutdown(ctx)
	fs.Redirect.Shutdown(ctx)
	utils.AppLog.Info("graceful frontend shutdown")
}

// func (fs *FrontendServer) OpenBolt() error {
// 	boltDB, err := bolt.Open(utils.Config.BoltDBPath, 0644, nil)
// 	if err != nil {
// 		utils.AppLog.Error("failed opening bolt", zap.Error(err))
// 		return err
// 	}
// 	fs.Database = BoltDB{DB: boltDB}
// 	return nil
// }
//
// func (fs *FrontendServer) CloseBolt() {
// 	fs.Database.DB.Close()
// }
