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

// Server hosts the fronend, user-facing, website and associated logic.
type Server struct {
	Primary    http.Server
	Redirect   http.Server
	httpClient http.Client
}

// TODO: This needs to be updated as well to include additional handlers.
func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", mainHandler)
	mux.HandleFunc("/setup_complete", setupCompleteHandler)
	return mux
}

// LaunchServer spins up goroutines for primary and redirect listeners.
func (s *Server) LaunchServer(secure, unsecure, cert, key string) {
	if PROD {
		// Primary server with HTTPS.
		s.Primary = http.Server{
			Addr:    secure,
			Handler: routes(),
		}
		go func() {
			if err := s.Primary.ListenAndServeTLS(cert, key); err != nil {
				utils.AppLog.Error(
					"primary server failed to start",
					zap.Error(err),
				)
				panic(err)
			}
		}()

		// For redirection purposes only.
		s.Redirect = http.Server{
			Addr:    unsecure,
			Handler: http.HandlerFunc(httpRedirect),
		}
		go func() {
			if err := s.Redirect.ListenAndServe(); err != nil {
				utils.AppLog.Error(
					"redirect server failed to start",
					zap.Error(err),
				)
				panic(err)
			}
		}()
	} else {
		// Primary server with HTTP for testing only. Ngrok doesn't play well
		// with HTTPS.
		s.Primary = http.Server{
			Addr:    unsecure,
			Handler: routes(),
		}
		go func() {
			if err := s.Primary.ListenAndServe(); err != nil {
				utils.AppLog.Error(
					"primary server failed to start",
					zap.Error(err),
				)
				panic(err)
			}
		}()
	}
}

// Start provides live/test LaunchServer with necessary startup information.
func (s *Server) Start() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	if PROD {
		s.LaunchServer(
			"10.142.1.0:443",
			"10.142.1.0:80",
			"heupr_io.crt",
			"heupr.key",
		)
	} else {
		s.LaunchServer(
			"127.0.0.1:8081",
			"127.0.0.1:8080",
			"cert.pem",
			"key.pem",
		)
	}

	<-stop
	utils.AppLog.Info("keyboard interrupt received")
	s.Stop()
}

// Stop gracefully closes down all server instances.
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Primary.Shutdown(ctx)
	s.Redirect.Shutdown(ctx)
	utils.AppLog.Info("graceful frontend shutdown")
}
