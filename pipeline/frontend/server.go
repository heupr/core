package frontend

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/csrf"
	"go.uber.org/zap"

	"core/utils"
)

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", render("../templates/home.html"))
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/repos", repos)
	mux.HandleFunc("/console", console)
	mux.HandleFunc("/complete", complete)
	mux.HandleFunc("/docs", render("../templates/docs.html"))
	mux.HandleFunc("/privacy", render("../templates/privacy.html"))
	mux.HandleFunc("/terms", render("../templates/terms.html"))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	return mux
}

// NOTE: This needs to be fleshed out (I assume).
func csrfError(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	fmt.Println("error in csrf: %v", csrf.FailureReason(r))
	return
}

func redirect(target string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	}
}

// Server hosts the frontend, user-facing website and associated logic.
type Server struct {
	Redirect http.Server
	Primary  http.Server
	Cert     string
	Key      string
}

// NewServer generates a frontend server while accounting for PROD/DEV tags.
func NewServer() *Server {
	secure, unsecure, redirectURL, cert, key := "", "", "", "", ""
	if PROD {
		secure, unsecure, redirectURL, cert, key = "10.142.1.0:443", "10.142.1.0:80", "https://heupr.io", "heupr_io.crt", "heupr.key"
	} else {
		secure, unsecure, redirectURL, cert, key = "127.0.0.1:8081", "127.0.0.1:8080", "https://127.0.0.1:8081", "cert.pem", "key.pem"
	}

	csrfProtection := csrf.Protect(
		[]byte("32-byte-long-auth-key"),
		csrf.Secure(false),
		csrf.ErrorHandler(http.HandlerFunc(csrfError)),
	)

	return &Server{
		Redirect: http.Server{
			Addr:    unsecure,
			Handler: http.HandlerFunc(redirect(redirectURL)),
		},
		Primary: http.Server{
			Addr:    secure,
			Handler: csrfProtection(routes()),
		},
		Cert: cert,
		Key:  key,
	}
}

// Start provides live/test LaunchServer with necessary startup information.
func (s *Server) Start() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := s.Redirect.ListenAndServe(); err != nil {
			utils.AppLog.Error("redirect server failed to start", zap.Error(err))
			panic(err)
		}
	}()

	go func() {
		if err := s.Primary.ListenAndServeTLS(s.Cert, s.Key); err != nil {
			utils.AppLog.Error("primary server failed to start", zap.Error(err))
			panic(err)
		}
	}()

	<-stop
	s.Stop()
	utils.AppLog.Info("keyboard interrupt received")
}

// Stop gracefully closes down all server instances.
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Redirect.Shutdown(ctx)
	s.Primary.Shutdown(ctx)
	utils.AppLog.Info("graceful frontend shutdown")
}
