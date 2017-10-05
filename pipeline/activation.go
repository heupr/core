package pipeline

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"

	"core/utils"
)

type ActivationServer struct {
	Server     http.Server
	httpClient http.Client
}

func (as *ActivationServer) activationServerHandler(w http.ResponseWriter, r *http.Request) {
	activationEndpoints := []string{utils.Config.IngestorActivationEndpoint, utils.Config.BackendActivationEndpoint}
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.AppLog.Error("failed to read payload", zap.Error(err))
		return
	}
	for i := range activationEndpoints {
		req, err := http.NewRequest("POST", activationEndpoints[i], bytes.NewBuffer(payload))
		if err != nil {
			utils.AppLog.Error("failed to create http request", zap.Error(err))
			continue
		}
		resp, err := as.httpClient.Do(req)
		if err != nil {
			errMsg := "failed internal post"
			utils.AppLog.Error(errMsg, zap.Error(err))
			http.Error(w, errMsg, http.StatusForbidden)
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
		Addr:    utils.Config.ActivationServerAddress,
		Handler: mux,
	}
	as.Server.ListenAndServe()
}
