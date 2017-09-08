package pipeline

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"

	"core/utils"
)

var (
	destinationBase  = "http://127.0.0.1"
	destinationPorts = []string{":8020", ":8030"}
	destinationEnd   = "/activate-ingestor-backend"
)

type ActivationServer struct {
	Server     http.Server
	httpClient http.Client
}

func (as *ActivationServer) activationServerHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.AppLog.Error("failed to read payload:", zap.Error(err))
		return
	}
	for i := range destinationPorts {
		req, err := http.NewRequest("POST", destinationBase+destinationPorts[i]+destinationEnd, bytes.NewBuffer(payload))
		if err != nil {
			utils.AppLog.Error("failed to create http request:", zap.Error(err))
			continue
		}
		resp, err := as.httpClient.Do(req)
		if err != nil {
			utils.AppLog.Error("failed internal post call:", zap.Error(err))
			http.Error(w, "failed internal post call", http.StatusForbidden)
			return
		} else {
			defer resp.Body.Close()
		}
	}
}

func (as *ActivationServer) Start() {
	as.httpClient = http.Client{Timeout: time.Second * 10}
	mux := http.NewServeMux()
	mux.HandleFunc("/activate-service", as.activationServerHandler)

	as.Server = http.Server{
		Addr:    "127.0.0.1:8010",
		Handler: mux,
	}
	as.Server.ListenAndServe()
}
