package frontend

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"core/utils"
)

type FrontendServer struct {
	Primary    http.Server
	Redirect   http.Server
	httpClient http.Client
}

var mainHandler = http.StripPrefix("/", http.FileServer(http.Dir("../website/")))

func (fs *FrontendServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", mainHandler)
	mux.HandleFunc("/setup_complete", fs.setupCompleteHandler)
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
	go func() {
		if err := fs.Redirect.ListenAndServe(); err != nil {
			utils.AppLog.Error("redirect server failed to start", zap.Error(err))
			panic(err)
		}
	}()
}

func (fs *FrontendServer) Start() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	if PROD {
		fs.LaunchServer("10.142.1.0:443", "10.142.1.0:80", "heupr_io.crt", "heupr.key")
	} else {
		fs.LaunchServer("127.0.0.1:8081", "127.0.0.1:8080", "cert.pem", "key.pem")
	}
	<-stop
	utils.AppLog.Info("keyboard interrupt received")
	fs.Stop()
}

func (fs *FrontendServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fs.Primary.Shutdown(ctx)
	fs.Redirect.Shutdown(ctx)
	utils.AppLog.Info("graceful frontend shutdown")
}
